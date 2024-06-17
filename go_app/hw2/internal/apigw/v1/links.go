package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

func newLinksHandler(linksClient linksClient) *linksHandler {
	return &linksHandler{client: linksClient}
}

type linksHandler struct {
	client linksClient
}

func (h *linksHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetLinks(context.Background(), &pb.GetLinksRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *linksHandler) PostLinks(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.client.CreateLink(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *linksHandler) DeleteLinksId(w http.ResponseWriter, r *http.Request, id string) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp, err := h.client.DeleteLink(context.Background(), &pb.DeleteLinkRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h Handler) GetLinksId(w http.ResponseWriter, r *http.Request, id string) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp, err := h.client.GetLink(context.Background(), &pb.GetLinkRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *linksHandler) PutLinksId(w http.ResponseWriter, r *http.Request, id string) {
	var req pb.UpdateLinkRequest
	vars := mux.Vars(r)
	req.Id = vars["id"]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.client.UpdateLink(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *linksHandler) GetLinksUserUserID(w http.ResponseWriter, r *http.Request, userID string) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	resp, err := h.client.GetLinksByUser(context.Background(), &pb.GetLinksByUserRequest{UserId: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}
