package main

import (
	"awesomeProject1/src/app/application"
	"awesomeProject1/src/app/domain"
	"awesomeProject1/src/app/infrastructure/Repositories"
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"net/http"
	"strings"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"pass"`
}

type Message struct {
	Text string `json:"text"`
}

var connection *websocket.Conn

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/new", "/api/user/login", "/ws", "/messages"}
		requestPath := r.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, "Missing auth token", http.StatusBadRequest)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, "Invalid/Malformed auth token", http.StatusBadRequest)
			return
		}

		tokenPart := splitted[1]
		tk := &domain.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("token_password"), nil
		})

		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, "Malformed authentication token", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, "Token is not valid.", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()

	router.Use(JwtAuthentication)

	router.HandleFunc("/ws", wsHandler)
	router.HandleFunc("/", createMessage).Methods("POST")
	router.HandleFunc("/api/user/new", CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login", Authenticate).Methods("POST")
	router.HandleFunc("/messages", GetMessages).Methods("GET")

	handler := cors.Default().Handler(router)
	panic(http.ListenAndServe(":8080", handler))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	connection = conn
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mes Message
	_ = json.NewDecoder(r.Body).Decode(&mes)
	userId := r.Context().Value("user").(string)
	service := application.NewMessageService()
	service.CreateMessage(mes.Text, connection, userId)
	w.WriteHeader(200)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	json.NewDecoder(r.Body).Decode(user)

	userRepo := Repositories.NewUserRepo()
	userRepo.CreateUser(user.Name, user.Password)

	w.WriteHeader(200)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	json.NewDecoder(r.Body).Decode(user)

	userRepo := Repositories.NewUserRepo()
	resp := userRepo.Login(user.Name, user.Password)

	json.NewEncoder(w).Encode(resp)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	result := application.NewMessageService().GetLastMessages()

	_ = json.NewEncoder(w).Encode(result)
	w.WriteHeader(200)
}
