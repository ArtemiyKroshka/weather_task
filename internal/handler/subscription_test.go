package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"weather_task/internal/service"

	"github.com/stretchr/testify/assert"
)

type mockSubscriptionSvc struct {
	subscribeErr   error
	confirmErr     error
	unsubscribeErr error
}

func (m *mockSubscriptionSvc) Subscribe(_ context.Context, _, _, _ string) error {
	return m.subscribeErr
}

func (m *mockSubscriptionSvc) Confirm(_ context.Context, _ string) error {
	return m.confirmErr
}

func (m *mockSubscriptionSvc) Unsubscribe(_ context.Context, _ string) error {
	return m.unsubscribeErr
}

func formBody(fields map[string]string) *strings.Reader {
	v := url.Values{}
	for k, val := range fields {
		v.Set(k, val)
	}
	return strings.NewReader(v.Encode())
}

func subscribeRequest(fields map[string]string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/subscriptions", formBody(fields))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestSubscribe_Success(t *testing.T) {
	h := NewSubscriptionHandler(&mockSubscriptionSvc{})
	w := httptest.NewRecorder()
	h.Subscribe(w, subscribeRequest(map[string]string{
		"email": "user@example.com", "city": "Kyiv", "frequency": "daily",
	}))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSubscribe_MissingFields(t *testing.T) {
	h := NewSubscriptionHandler(&mockSubscriptionSvc{})
	w := httptest.NewRecorder()
	h.Subscribe(w, subscribeRequest(map[string]string{"email": "user@example.com"}))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assertErrorBody(t, w, "required")
}

func TestSubscribe_InvalidFrequency(t *testing.T) {
	h := NewSubscriptionHandler(&mockSubscriptionSvc{})
	w := httptest.NewRecorder()
	h.Subscribe(w, subscribeRequest(map[string]string{
		"email": "user@example.com", "city": "Kyiv", "frequency": "weekly",
	}))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assertErrorBody(t, w, "frequency")
}

func TestSubscribe_AlreadySubscribed(t *testing.T) {
	h := NewSubscriptionHandler(&mockSubscriptionSvc{subscribeErr: service.ErrAlreadySubscribed})
	w := httptest.NewRecorder()
	h.Subscribe(w, subscribeRequest(map[string]string{
		"email": "user@example.com", "city": "Kyiv", "frequency": "daily",
	}))
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestSubscribe_InternalError(t *testing.T) {
	h := NewSubscriptionHandler(&mockSubscriptionSvc{subscribeErr: assert.AnError})
	w := httptest.NewRecorder()
	h.Subscribe(w, subscribeRequest(map[string]string{
		"email": "user@example.com", "city": "Kyiv", "frequency": "daily",
	}))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func confirmRequest(token string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/api/subscriptions/confirm/"+token, nil)
	// simulate path value set by net/http mux
	req.SetPathValue("token", token)
	return req
}

func TestConfirm_Success(t *testing.T) {
	h := NewSubscriptionHandler(&mockSubscriptionSvc{})
	w := httptest.NewRecorder()
	h.Confirm(w, confirmRequest("valid-token"))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestConfirm_TokenNotFound(t *testing.T) {
	h := NewSubscriptionHandler(&mockSubscriptionSvc{confirmErr: service.ErrTokenNotFound})
	w := httptest.NewRecorder()
	h.Confirm(w, confirmRequest("bad-token"))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func unsubscribeRequest(token string) *http.Request {
	req := httptest.NewRequest(http.MethodDelete, "/api/subscriptions/"+token, nil)
	req.SetPathValue("token", token)
	return req
}

func TestUnsubscribe_Success(t *testing.T) {
	h := NewSubscriptionHandler(&mockSubscriptionSvc{})
	w := httptest.NewRecorder()
	h.Unsubscribe(w, unsubscribeRequest("valid-token"))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUnsubscribe_TokenNotFound(t *testing.T) {
	h := NewSubscriptionHandler(&mockSubscriptionSvc{unsubscribeErr: service.ErrTokenNotFound})
	w := httptest.NewRecorder()
	h.Unsubscribe(w, unsubscribeRequest("bad-token"))
	assert.Equal(t, http.StatusNotFound, w.Code)
}
