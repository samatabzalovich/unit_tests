package main

import (
	"net/http"
	)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
		"environment": app.config.env,
		"version": version,
		"user_name": app.contextGetUser(r).Name,
		},
	}
	
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w,r,err)
	}
}
