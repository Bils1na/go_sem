package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []int    `json:"friends"`
}

var (
	users = make(map[int]*User)
	usersMutex sync.Mutex
	nextID = 1
)

func main() {
		r := chi.NewRouter()
		r.Use(middleware.Logger)

		r.Post("/create", createUserHandler)
		r.Post("/make_friends", makeFriendsHandler)
		r.Delete("/user", deleteUserHandler)
		r.Get("/friends/{userID}", getFriendsHandler)
		r.Put("/{userID}", updateUserAgeHandler)

		http.ListenAndServe(":8080", r)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	user.ID = nextID
	nextID++
	users[user.ID] = &user
	usersMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": user.ID})
}

func makeFriendsHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SourceID int json:"source_id"
		TargetID int json:"target_id"
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
 	}

	usersMutex.Lock()
	sourceUser, sourceExists := users[request.SourceID]
	targetUser, targetExists := users[request.TargetID]
	if sourceExists && targetExists {
		sourceUser.Friends = append(sourceUser.Friends, request.TargetID)
		targetUser.Friends = append(targetUser.Friends, request.SourceID)
	}
	usersMutex.Unlock()

	if !sourceExists || !targetExists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		TargetID int json:"target_id"
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()
	user, exists := users[request.TargetID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	for _, friendID := range user.Friends {
		friend := users[friendID]
		newFriends := []int{}
		for _, id := range friend.Friends {
			if id != user.ID {
				newFriends = append(newFriends, id)
			}
		}
		friend.Friends = newFriends
	}

	delete(users, request.TargetID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"name": user.Name})
}

func getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	usersMutex.Lock()
	defer usersMutex.Unlock()

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	friends := []User{}
	for _, friendID := range user.Friends {
		if friend, exists := users[friendID]; exists {
			friends = append(friends, *friend)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)
}

func updateUserAgeHandler(w http.ResponseWriter, r *http.Request) {
 	userID := chi.URLParam(r, "userID")

	var request struct {
	NewAge int json:"new age"
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Age = request.NewAge

	w.WriteHeader(http.StatusOK)
	message := "возраст пользователя успешно обновлён"
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}