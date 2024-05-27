package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID      int      
	Name    string   
	Age     int      
	Friends []int    
}

func main() {
	initDB()
	defer closeDB()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/create", createUserHandler)
	r.Post("/make_friends", makeFriendsHandler)
	r.Delete("/user", deleteUserHandler)
	r.Get("/friends/{userID}", getFriendsHandler)
	r.Put("/{userID}", updateUserAgeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on port", port)
	http.ListenAndServe(":" + port, r)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	friends := strings.Join(intSliceToStringSlice(user.Friends), ",")
	result, err := db.Exec("INSERT INTO users (name, age, friends) VALUES (?, ?, ?)", user.Name, user.Age, friends)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func makeFriendsHandler(w http.ResponseWriter, r *http.Request) {
 	var request struct {
		SourceID int 
		TargetID int 
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sourceUser, err := getUserByID(tx, request.SourceID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	targetUser, err := getUserByID(tx, request.TargetID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	sourceUser.Friends = append(sourceUser.Friends, request.TargetID)
	targetUser.Friends = append(targetUser.Friends, request.SourceID)

	if err := updateUserFriends(tx, sourceUser.ID, sourceUser.Friends); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := updateUserFriends(tx, targetUser.ID, targetUser.Friends); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	message := fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		TargetID int
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := getUserByID(tx, request.TargetID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	for _, friendID := range user.Friends {
	friend, err := getUserByID(tx, friendID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		friend.Friends = removeFriend(friend.Friends, user.ID)
		if err := updateUserFriends(tx, friend.ID, friend.Friends); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if _, err := tx.Exec("DELETE FROM users WHERE id = ?", user.ID); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"name": user.Name})
}

func getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
   
	user, err := getUserByID(db, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
   
	friends := []User{}
	for _, friendID := range user.Friends {
		friend, err := getUserByID(db, friendID)
		if err == nil {
			friends = append(friends, *friend)
		}
	}
   
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)
}
   
func updateUserAgeHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	var request struct {
	 	NewAge int
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
   
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
   
	if _, err := db.Exec("UPDATE users SET age = ? WHERE id = ?", request.NewAge, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
   
	w.WriteHeader(http.StatusOK)
	message := "возраст пользователя успешно обновлён"
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
   
func getUserByID(db *sql.DB, id int) (*User, error) {
	var user User
	var friends string
	row := db.QueryRow("SELECT id, name, age, friends FROM users WHERE id = ?", id)
	if err := row.Scan(&user.ID, &user.Name, &user.Age, &friends); err != nil {
		return nil, err
	}
	user.Friends = stringSliceToIntSlice(strings.Split(friends, ","))
	return &user, nil
}
   
func updateUserFriends(db *sql.DB, id int, friends []int) error {
	friendsStr := strings.Join(intSliceToStringSlice(friends), ",")
	_, err := db.Exec("UPDATE users SET friends = ? WHERE id = ?", friendsStr, id)
	return err
}
   
func intSliceToStringSlice(ints []int) []string {
	strs := make([]string, len(ints))
	for i, v := range ints {
	 	strs[i] = strconv.Itoa(v)
	}
	return strs
}
   
func stringSliceToIntSlice(strs []string) []int {
	ints := make([]int, len(strs))
	for i, v := range strs {
		if val, err := strconv.Atoi(v); err == nil {
			ints[i] = val
		}
	}
	return ints
}
   
func removeFriend(friends []int, friendID int) []int {
	newFriends := []int{}
	for _, id := range friends {
		if id != friendID {
			newFriends = append(newFriends, id)
		}
	}
	return newFriends
}