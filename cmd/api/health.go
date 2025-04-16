package main

import (
	"net/http"
)

// @Summary		Health Check
// @Description	Check if the server is running
// @Tags			ops
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]string	"Example: {\"status\":\"ok\",\"env\":\"development\",\"version\":\"0.0.2\"}"
// @Failure		500	{object}	nil
// @Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
