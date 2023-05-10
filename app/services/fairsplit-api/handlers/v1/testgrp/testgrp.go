package testgrp

import (
	"context"
	"net/http"

	"github.com/shekosk1/webservice-kit/foundation/web"
)

// Status represents a test handler for now.
func Status(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	data := struct {
		Status string
	}{
		Status: "TEST_OK",
	}

	return web.Respond(ctx, w, data, http.StatusOK)
}
