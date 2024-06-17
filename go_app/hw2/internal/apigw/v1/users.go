package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

func newUsersHandler(usersClient usersClient) *usersHandler {
	return &usersHandler{client: usersClient}
}

type usersHandler struct {
	client usersClient
}

func (h *usersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.GetUsers(context.Background(), &pb.GetUsersRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *usersHandler) PostUsers(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.client.CreateUser(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *usersHandler) DeleteUsersId(w http.ResponseWriter, r *http.Request, id string) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp, err := h.client.DeleteUser(context.Background(), &pb.DeleteUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *usersHandler) GetUsersId(w http.ResponseWriter, r *http.Request, id string) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp, err := h.client.GetUser(context.Background(), &pb.GetUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *usersHandler) PutUsersId(w http.ResponseWriter, r *http.Request, id string) {
	var req pb.UpdateUserRequest
	vars := mux.Vars(r)
	req.Id = vars["id"]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.client.UpdateUser(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}
