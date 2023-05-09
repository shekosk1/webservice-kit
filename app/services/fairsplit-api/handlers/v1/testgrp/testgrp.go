package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	v1 "github.com/shekosk1/webservice-kit/business/web/v1"
	"github.com/shekosk1/webservice-kit/foundation/web"
)

// Status represents a test handler for now.
func Status(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return v1.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}

	data := struct {
		Status string
	}{
		Status: "TEST_OK",
	}

	return web.Respond(ctx, w, data, http.StatusOK)
}

// Status represents a test handler for now.
func Empty(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return errors.New("NON trusted error")
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
