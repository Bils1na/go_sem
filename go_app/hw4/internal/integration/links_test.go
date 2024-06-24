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

type mockLinksClient struct{}

func (m *mockLinksClient) ListLinks(ctx context.Context, in *pb.Empty) (*pb.ListLinksResponse, error) {
	return &pb.ListLinksResponse{Links: []*pb.Link{}}, nil
}

func (m *mockLinksClient) CreateLink(ctx context.Context, in *pb.CreateLinkRequest) (*pb.Link, error) {
	return &pb.Link{Id: in.Id}, nil
}

func (m *mockLinksClient) DeleteLink(ctx context.Context, in *pb.DeleteLinkRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (m *mockLinksClient) GetLink(ctx context.Context, in *pb.GetLinkRequest) (*pb.Link, error) {
	return &pb.Link{Id: in.Id}, nil
}

func (m *mockLinksClient) UpdateLink(ctx context.Context, in *pb.UpdateLinkRequest) (*pb.Link, error) {
	return &pb.Link{Id: in.Id}, nil
}

func TestPostLinks(t *testing.T) {
	handler := v1.NewLinksHandler(&mockLinksClient{})

	link := apiv1.LinkCreate{
		Id:     "test-id",
		Title:  "test-title",
		Url:    "http://test-url",
		Images: []string{"http://test-image"},
		Tags:   []string{"test-tag"},
		UserId: "test-user",
	}
	body, _ := json.Marshal(link)

	req := httptest.NewRequest(http.MethodPost, "/links", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.PostLinks(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestDeleteLinksId(t *testing.T) {
	handler := v1.NewLinksHandler(&mockLinksClient{})

	req := httptest.NewRequest(http.MethodDelete, "/links/test-id", nil)
	w := httptest.NewRecorder()

	handler.DeleteLinksId(w, req, "test-id")

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}
