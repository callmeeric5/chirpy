package main

import (
	"go_server/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	dbToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not found", err)
		return
	}

	if time.Now().UTC().After(dbToken.ExpiredAt) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", err)
		return
	}

	if dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", err)
		return
	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	type response struct{}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "can't get refresh_token", err)
		return
	}
	err = cfg.db.UpdateRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "can't update refresh_token", err)
		return
	}
	respondWithJSON(w, 204, response{})
}
