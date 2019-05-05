package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/smapig/secret-server/helpers"
	"github.com/smapig/secret-server/models"
)

var AddSecret = func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	secretText := r.FormValue("secret")
	expireAfterViews, err1 := strconv.Atoi(r.FormValue("expireAfterViews"))
	expireAfter, err2 := strconv.Atoi(r.FormValue("expireAfter"))

	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	secret := &models.Secret{
		SecretText:     secretText,
		RemainingViews: expireAfterViews,
	}

	if expireAfter == 0 {
		secret.ExpiresAt = &time.Time{}
	} else {
		expiresAt := time.Now().Add(time.Minute * time.Duration(expireAfter))
		secret.ExpiresAt = &expiresAt
	}

	err := secret.Create()
	contentType := r.Header.Get("accept")

	var rv map[string]interface{}

	if err != nil {
		rv = helpers.Message(false, "Something went wrong!")
	} else {
		rv := helpers.Message(true, "Add secret successfully!")
		rv["data"] = secret
	}

	helpers.Respond(w, rv, contentType, err)
}

var GetSecret = func(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("accept")

	params := mux.Vars(r)
	hash := params["hash"]
	if hash == "" {
		helpers.Respond(w, helpers.Message(false, "No hash provide."), contentType, nil)
		return
	}

	var rv map[string]interface{}
	var err error
	secret := models.GetSecretByHash(hash)

	if secret.IsAlive() {
		err = secret.DecreaseRemainingViews()

		if err != nil {
			rv = helpers.Message(false, "Something went wrong.")
		} else {
			rv = helpers.Message(true, "Get secret successfully!")
			rv["data"] = secret
		}
	} else {
		rv = helpers.Message(false, "A secret makes a woman, woman.")
	}

	helpers.Respond(w, rv, contentType, err)
}
