package model

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

// TestComment_Validate covers both the happy path and the error path.
func TestComment_Validate(t *testing.T) {
	tests := []struct {
		name     string
		comment  Comment
		expected bool
	}{
		{
			name: "Valid comment",
			comment: Comment{
				Body:      "This is a valid comment",
				UserID:    1,
				ArticleID: 1,
			},
			expected: true,
		},
		{
			name: "Empty body",
			comment: Comment{
				UserID:    1,
				ArticleID: 1,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.comment.Validate()
			assert.Equal(t, tt.expected, err == nil)
		})
	}
}

// TestComment_ProtoComment verifies that the proto conversion uses the embedded ID and timestamps.
func TestComment_ProtoComment(t *testing.T) {
	createdAt := time.Date(2025, time.April, 28, 12, 34, 56, 0, time.UTC)
	updatedAt := createdAt.Add(5 * time.Minute)

	// Embed a gorm.Model so ID, CreatedAt, UpdatedAt are set
	comment := Comment{
		Model:     gorm.Model{ID: 42, CreatedAt: createdAt, UpdatedAt: updatedAt},
		Body:      "Test comment",
		UserID:    7,
		ArticleID: 13,
	}

	protoComment := comment.ProtoComment()

	expected := &pb.Comment{
		Id:        "42",
		Body:      "Test comment",
		CreatedAt: createdAt.Format(ISO8601),
		UpdatedAt: updatedAt.Format(ISO8601),
	}

	assert.Equal(t, expected, protoComment)
}
