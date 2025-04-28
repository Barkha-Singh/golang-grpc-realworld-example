package handler

import (
	"context"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

// --- Define a local context key type for testing
type userIDKey struct{}

// helper to inject userID into context
func contextWithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func TestHandler_ShowProfile(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	fooUser := model.User{
		Username: "foo",
		Email:    "foo@example.com",
		Password: "secret",
	}

	barUser := model.User{
		Username: "bar",
		Email:    "bar@example.com",
		Password: "secret",
	}

	// Create foo user
	if err := h.us.Create(&fooUser); err != nil {
		t.Fatalf("failed to create initial foo user: %v", err)
	}

	// Create bar user
	if err := h.us.Create(&barUser); err != nil {
		t.Fatalf("failed to create initial bar user: %v", err)
	}

	// Foo follows Bar
	if err := h.us.Follow(&fooUser, &barUser); err != nil {
		t.Fatalf("failed to create follow relationship: %v", err)
	}

	t.Run("ShowProfile_Successful", func(t *testing.T) {
		ctx := contextWithUserID(context.Background(), fooUser.ID)
		req := &proto.ShowProfileRequest{Username: "bar"}

		resp, err := h.ShowProfile(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "bar", resp.Profile.Username)
		assert.True(t, resp.Profile.Following)
	})

	t.Run("ShowProfile_UserNotFound", func(t *testing.T) {
		ctx := contextWithUserID(context.Background(), fooUser.ID)
		req := &proto.ShowProfileRequest{Username: "nonexistent_user"}

		resp, err := h.ShowProfile(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("ShowProfile_Unauthenticated", func(t *testing.T) {
		ctx := context.Background() // No user ID injected
		req := &proto.ShowProfileRequest{Username: "bar"}

		resp, err := h.ShowProfile(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
