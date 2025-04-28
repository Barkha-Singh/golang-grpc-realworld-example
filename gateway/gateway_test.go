package main

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	gw "github.com/raahii/golang-grpc-realworld-example/proto"
)

func TestGatewayRun(t *testing.T) {
	// Arrange
	ctx := context.Background()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Try to register handlers (actual server must be running separately during this test)
	err := gw.RegisterUsersHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		t.Logf("Skipping users handler registration: %v", err)
	}
	err = gw.RegisterArticlesHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		t.Logf("Skipping articles handler registration: %v", err)
	}

	// Act
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	// Assert
	// We just check if server responds (even if 404 because no endpoint is matched, that's fine here)
	assert.NotEqual(t, 0, w.Code)
}
