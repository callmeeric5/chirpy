package main

import (
	"encoding/json"
	"go_server/internal/auth"
	"go_server/internal/database"
	"log"
	"net/http"
	"sort"

	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		post database.Post
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("GetBearerToken error: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't get token", err)
		return
	}
	user_uuid, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("ValidateJWT error: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Can't valide token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleaned := cleanProfanity(params.Body)
	post, err := cfg.db.CreatePost(r.Context(), database.CreatePostParams{
		Body:   cleaned,
		UserID: user_uuid,
	})
	if err != nil {
		log.Fatalf("Can't create post of User_id: %d", params.User_id)
		return
	}
	respondWithJSON(w, 201, post)
}

func (cfg *apiConfig) handlerGetPost(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid post ID", err)
		return
	}

	post, err := cfg.db.GetPost(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get post", err)
		return
	}
	respondWithJSON(w, http.StatusOK, post)
}

func (cfg *apiConfig) handlerGetPosts(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")
	if sortParam == "" {
		sortParam = "asc"
	}

	var posts []database.Post
	var err error

	if user_id != "" {
		user_uuid, _ := uuid.Parse(user_id)
		posts, err = cfg.db.GetPostsByUserId(r.Context(), user_uuid)
	} else {
		posts, err = cfg.db.GetPosts(r.Context())

	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't get posts", err)
	}
	sort.Slice(posts, func(i, j int) bool {
		if sortParam == "desc" {
			return posts[i].CreatedAt.After(posts[j].CreatedAt)
		}
		return posts[i].CreatedAt.Before(posts[j].CreatedAt)
	})
	respondWithJSON(w, http.StatusOK, posts)
}

func cleanProfanity(text string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(text, " ")

	for i, word := range words {
		lowercaseWord := strings.ToLower(word)
		for _, profane := range profaneWords {
			if lowercaseWord == profane {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerDeletePost(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	idStr := r.PathValue("id")
	postID, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid post ID", err)
		return
	}

	post, err := cfg.db.GetPost(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Post not found", err)
		return
	}

	if userID != post.UserID {
		respondWithError(w, http.StatusForbidden, "You are not the author of this post", err)
		return
	}

	_, err = cfg.db.DeletePost(r.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete post", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
