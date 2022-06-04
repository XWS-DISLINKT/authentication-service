package startup

import (
	"authentication_service/infrastructure/services"
	"authentication_service/startup/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-DISLINKT/dislinkt/common/proto/profile-service"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Id       string `json:"id"`
}

type Claims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	server := &Server{
		config: config,
	}
	server.initHandlers()
	return server
}

func (server *Server) initHandlers() {
	http.HandleFunc("/login", server.LogIn)
	http.HandleFunc("/refresh", server.Refresh)
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), nil))
}

func (server *Server) LogIn(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	profileEndpoint := fmt.Sprintf("%s:%s", server.config.ProfileHost, server.config.ProfilePort)
	profileClient := services.NewProfileClient(profileEndpoint)
	response, err := profileClient.GetCredentials(context.TODO(), &profile.GetCredentialsRequest{Username: credentials.Username})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err != nil || response.Username != credentials.Username || response.Password != credentials.Password {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := &Claims{
		Id:       response.Id,
		Username: response.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	res, err := json.Marshal(claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (server *Server) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Minute {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(30 * time.Second)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
