package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/raulcv/goapiws/models"
	"github.com/raulcv/goapiws/repository"
	"github.com/raulcv/goapiws/server"
	"github.com/segmentio/ksuid"
)

type AddPostRequest struct {
	Content string `json:"content"`
}
type AddPostResponse struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}
type UpdatePostRequest struct {
	Content string `json:"content"`
}
type UpdatePostResponse struct {
	Message string `json:"message"`
}

type PaginationPostResponse struct {
	Posts []*models.Post `json:"data"`
	Total int64          `json:"total"`
}

func AddPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			// Creating a post
			var postRequest = AddPostRequest{}
			err := json.NewDecoder(r.Body).Decode(&postRequest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			generatedID, err := ksuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var post = models.Post{
				Id:      generatedID.String(),
				Content: postRequest.Content,
				UserId:  claims.UserId, //UserId
			}
			err = repository.AddPost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// After created a data in db create a message throught websocket
			var postMessage = models.WebsocketMessage{
				Type:    "Post_Created",
				Payload: post,
			}
			s.Hub().BroadCast(postMessage, nil)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(AddPostResponse{
				Id:      post.Id,
				Content: post.Content,
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if post.Id == "" {
			http.Error(w, "Post was deleted", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		postRequest := UpdatePostRequest{}
		err := json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {

			var newPost = models.Post{
				Id:        params["id"],
				Content:   postRequest.Content,
				UserId:    claims.UserId,
				UpdatedAt: time.Now(),
			}
			err := repository.UpdatePost(r.Context(), &newPost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			postResponse := UpdatePostResponse{Message: "Post Updated Successfully"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(postResponse)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeletePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {

			params := mux.Vars(r)

			n, err := repository.DeletePost(r.Context(), params["id"], claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if n <= 0 {
				http.Error(w, "Post not found, It might was deleted already", http.StatusUnauthorized)
				return
			}
			deletePostResponse := UpdatePostResponse{Message: "Post was deleted successfully"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(deletePostResponse)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ActivatePostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		err := repository.ActivatePost(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		deletePostResponse := UpdatePostResponse{Message: "Post was Activated successfully"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deletePostResponse)
	}
}

func ListPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		pageStr := r.URL.Query().Get("page") //content?params=1&limit=10
		var page = uint64(0)
		if pageStr != "" {
			page, err = strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		fmt.Println("pageStr: ", pageStr)
		posts, err := repository.ListPost(r.Context(), page)
		fmt.Println("pageStr: ", posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func ListPostTwoHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		pageStr := r.URL.Query().Get("page") //content?params=1&limit=10
		var page = uint64(0)
		if pageStr != "" {
			page, err = strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		fmt.Println("pageStr: ", pageStr)
		posts, total, err := repository.ListPostTwo(r.Context(), page)
		fmt.Println("pageStr: ", posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		postResponse := &PaginationPostResponse{posts, total}
		//  total := 122

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(postResponse)
	}
}
