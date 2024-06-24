package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/api/apiv1"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/pb"
	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/pkg/v1"
)

type mockUsersClient struct{}

func (m *mockUsersClient) ListUsers(ctx context.Context, in *pb.Empty) (*pb.ListUsersResponse, error) {
	return &pb.ListUsersResponse{Users: []*pb.User{}}, nil
}

func (m *mockUsersClient) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
	return &pb.User{Id: in.Id}, nil
}

func (m *mockUsersClient) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (m *mockUsersClient) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	return &pb.User{Id: in.Id}, nil
}

func (m *mockUsersClient) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.User, error) {
	return &pb.User{Id: in.Id}, nil
}

func TestPostUsers(t *testing.T) {
	handler := v1.NewUsersHandler(&mockUsersClient{})

	user := apiv1.UserCreate{
		Id:       "test-id",
		Username: "test-username",
		Password: "test-password",
	}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.PostUsers(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestDeleteUsersId(t *testing.T) {
	handler := v1.NewUsersHandler(&mockUsersClient{})

	req := httptest.NewRequest(http.MethodDelete, "/users/test-id", nil)
	w := httptest.NewRecorder()

	handler.DeleteUsersId(w, req, "test-id")

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}
