package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	err = app.loggerReq("auth", fmt.Sprintf("%s logged in", user.Email))

	if err != nil {
		app.errorJSON(w, err)
	}

	payload := jsonResponse{
		Error:   false,
		Data:    user,
		Message: fmt.Sprintf("Logged in User %s", user.Email),
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) loggerReq(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Data = data
	entry.Name = name
	logUrl := "http://logger-service/log"
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	request, err := http.NewRequest("POST", logUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}

	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}
