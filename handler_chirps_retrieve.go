package main

import "net/http"

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.RetrieveAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}
	apiChirps := make([]Chirp, 0, len(chirps))
	for _, c := range chirps {
		apiChirps = append(apiChirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, apiChirps)
}
