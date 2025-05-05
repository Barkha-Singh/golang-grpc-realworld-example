package handler

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
   ---------------------------------------------------------------------------
   Mock definitions – generated-style (kept local to avoid external code-gen)
   ---------------------------------------------------------------------------
*/

// MockUserService is a gomock compatible mock for the user service used by
// the handler.  Only the methods needed by ShowProfile are mocked.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

type MockUserServiceMockRecorder struct{ mock *MockUserService }

func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	m := &MockUserService{ctrl: ctrl}
	m.recorder = &MockUserServiceMockRecorder{mock: m}
	return m
}

func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder { return m.recorder }

// --- mocked methods -----------------------------------------------------------------

func (m *MockUserService) GetByID(id uint) (*model.User, error) {
	ret := m.ctrl.Call(m, "GetByID", id)
	user, _ := ret[0].(*model.User)
	err, _ := ret[1].(error)
	return user, err
}

func (mr *MockUserServiceMockRecorder) GetByID(id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID",
		reflect.TypeOf((*MockUserService)(nil).GetByID), id)
}

func (m *MockUserService) GetByUsername(username string) (*model.User, error) {
	ret := m.ctrl.Call(m, "GetByUsername", username)
	user, _ := ret[0].(*model.User)
	err, _ := ret[1].(error)
	return user, err
}

func (mr *MockUserServiceMockRecorder) GetByUsername(username interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername",
		reflect.TypeOf((*MockUserService)(nil).GetByUsername), username)
}

func (m *MockUserService) IsFollowing(currentUser, targetUser *model.User) (bool, error) {
	ret := m.ctrl.Call(m, "IsFollowing", currentUser, targetUser)
	following, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return following, err
}

func (mr *MockUserServiceMockRecorder) IsFollowing(current, target interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFollowing",
		reflect.TypeOf((*MockUserService)(nil).IsFollowing), current, target)
}

/*
   ---------------------------------------------------------------------------
   Unit tests for Handler.ShowProfile
   ---------------------------------------------------------------------------
*/

func TestHandler_ShowProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zerolog.New(io.Discard)

	// shared fixtures
	current := &model.User{ID: 1, Username: "current", Email: "cur@example.com"}
	target := &model.User{ID: 2, Username: "target", Email: "tar@example.com"}

	type want struct {
		code      codes.Code
		following bool
	}

	tests := []struct {
		name    string
		ctx     context.Context
		prepare func(us *MockUserService)
		req     *pb.ShowProfileRequest
		want    want
	}{
		{
			name: "success – following",
			ctx:  auth.NewContextWithUserID(context.Background(), current.ID),
			prepare: func(us *MockUserService) {
				us.EXPECT().GetByID(current.ID).Return(current, nil)
				us.EXPECT().GetByUsername(target.Username).Return(target, nil)
				us.EXPECT().IsFollowing(current, target).Return(true, nil)
			},
			req:  &pb.ShowProfileRequest{Username: target.Username},
			want: want{code: codes.OK, following: true},
		},
		{
			name: "unauthenticated when context lacks user id",
			ctx:  context.Background(),
			prepare: func(us *MockUserService) {
				// no interaction expected
			},
			req:  &pb.ShowProfileRequest{Username: target.Username},
			want: want{code: codes.Unauthenticated},
		},
		{
			name: "current user not found",
			ctx:  auth.NewContextWithUserID(context.Background(), current.ID),
			prepare: func(us *MockUserService) {
				us.EXPECT().GetByID(current.ID).Return(nil, errors.New("not found"))
			},
			req:  &pb.ShowProfileRequest{Username: target.Username},
			want: want{code: codes.NotFound},
		},
		{
			name: "requested user not found",
			ctx:  auth.NewContextWithUserID(context.Background(), current.ID),
			prepare: func(us *MockUserService) {
				us.EXPECT().GetByID(current.ID).Return(current, nil)
				us.EXPECT().GetByUsername(target.Username).Return(nil, errors.New("not found"))
			},
			req:  &pb.ShowProfileRequest{Username: target.Username},
			want: want{code: codes.NotFound},
		},
		{
			name: "error while checking following relationship",
			ctx:  auth.NewContextWithUserID(context.Background(), current.ID),
			prepare: func(us *MockUserService) {
				us.EXPECT().GetByID(current.ID).Return(current, nil)
				us.EXPECT().GetByUsername(target.Username).Return(target, nil)
				us.EXPECT().IsFollowing(current, target).Return(false, errors.New("db error"))
			},
			req:  &pb.ShowProfileRequest{Username: target.Username},
			want: want{code: codes.NotFound}, // handler maps to NotFound w/ “internal server error”
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			usMock := NewMockUserService(ctrl)
			tt.prepare(usMock)

			h := &Handler{
				us:     usMock, // mock implementation
				logger: logger, // discard logger
			}

			resp, err := h.ShowProfile(tt.ctx, tt.req)

			if tt.want.code == codes.OK {
				// success path
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want.following, resp.Profile.Following)
				assert.Equal(t, tt.req.Username, resp.Profile.Username)
			} else {
				// error path
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				assert.True(t, ok, "error should be a grpc status")
				assert.Equal(t, tt.want.code, st.Code())
			}
		})
	}
}
