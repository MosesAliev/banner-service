package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserBannerThanOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/userbanner?feature_id=0&use_last_revision=true&limit=1&offset=1&tag_id=0", nil)
	req.Header.Add("token", "user_token")
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getUserBannerHandler)
	handler.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestGetUserBannerThanBadRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/userbanner?use_last_revision=true&limit=1&offset=1&tag_id=0", nil)
	req.Header.Add("token", "user_token")
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getUserBannerHandler)
	handler.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestGetUserBannerThanUnauthorized(t *testing.T) {
	req := httptest.NewRequest("GET", "/userbanner?feature_id=0&use_last_revision=true&limit=1&offset=1&tag_id=0", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getUserBannerHandler)
	handler.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestGetUserBannerThanForbidden(t *testing.T) {
	req := httptest.NewRequest("GET", "/userbanner?feature_id=0&use_last_revision=true&limit=1&offset=1&tag_id=0", nil)
	req.Header.Add("token", "token")
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getUserBannerHandler)
	handler.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusForbidden, responseRecorder.Code)
}
