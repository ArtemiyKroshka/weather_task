package handler

import (
	"errors"
	"net/http"

	"weather_task/internal/service"
)

// SubscriptionHandler handles subscription lifecycle endpoints.
type SubscriptionHandler struct {
	svc service.SubscriptionService
}

// NewSubscriptionHandler creates a SubscriptionHandler with the given service.
func NewSubscriptionHandler(svc service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

// Subscribe godoc
//
//	@Summary		Subscribe to weather updates
//	@Tags			subscriptions
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			email		formData	string	true	"Email address"
//	@Param			city		formData	string	true	"City name"
//	@Param			frequency	formData	string	true	"hourly or daily"	Enums(hourly, daily)
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	map[string]string
//	@Failure		409			{object}	map[string]string
//	@Router			/api/subscriptions [post]
func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	email := r.FormValue("email")
	city := r.FormValue("city")
	freq := r.FormValue("frequency")

	if email == "" || city == "" || freq == "" {
		writeError(w, http.StatusBadRequest, "email, city, and frequency are required")
		return
	}
	if freq != "hourly" && freq != "daily" {
		writeError(w, http.StatusBadRequest, "frequency must be 'hourly' or 'daily'")
		return
	}

	if err := h.svc.Subscribe(r.Context(), email, city, freq); err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadySubscribed):
			writeError(w, http.StatusConflict, "email already subscribed")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "subscription created, please check your email for confirmation",
	})
}

// Confirm godoc
//
//	@Summary		Confirm a subscription
//	@Tags			subscriptions
//	@Produce		json
//	@Param			token	path		string	true	"Confirmation token"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/api/subscriptions/confirm/{token} [post]
func (h *SubscriptionHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}

	if err := h.svc.Confirm(r.Context(), token); err != nil {
		switch {
		case errors.Is(err, service.ErrTokenNotFound):
			writeError(w, http.StatusNotFound, "token not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "subscription confirmed"})
}

// Unsubscribe godoc
//
//	@Summary		Unsubscribe from weather updates
//	@Tags			subscriptions
//	@Produce		json
//	@Param			token	path		string	true	"Unsubscribe token"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/api/subscriptions/{token} [delete]
func (h *SubscriptionHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}

	if err := h.svc.Unsubscribe(r.Context(), token); err != nil {
		switch {
		case errors.Is(err, service.ErrTokenNotFound):
			writeError(w, http.StatusNotFound, "token not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "unsubscribed successfully"})
}
