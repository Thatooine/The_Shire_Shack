package restaurants

import (
	"encoding/json"
	"net/http"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/authentication"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	"github.com/rs/zerolog/log"
)

// RatingSubmitterServiceRESTAdaptor exposes rating submission over a REST API.
type RatingSubmitterServiceRESTAdaptor struct {
	submitter RatingSubmitterService
}

// NewRatingSubmitterServiceRESTAdaptor returns a new RatingSubmitterServiceRESTAdaptor.
func NewRatingSubmitterServiceRESTAdaptor(submitter RatingSubmitterService) *RatingSubmitterServiceRESTAdaptor {
	return &RatingSubmitterServiceRESTAdaptor{submitter: submitter}
}

// SubmitRating

type SubmitRatingRESTRequest struct {
	DishID string `json:"dish_id"`
	Score  int    `json:"score"`
	Review string `json:"review"`
}

type SubmitRatingRESTResponse struct {
	Rating Rating `json:"rating"`
}

func (a *RatingSubmitterServiceRESTAdaptor) SubmitRating(w http.ResponseWriter, r *http.Request) {
	var request SubmitRatingRESTRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to decode submit rating request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	claim, ok := authentication.LoginClaimFromContext(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	resp, err := a.submitter.SubmitRating(r.Context(), SubmitRatingRequest{
		DishID: request.DishID,
		UserID: claim.UserID,
		Score:  request.Score,
		Review: request.Review,
	})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to submit rating")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SubmitRatingRESTResponse{Rating: resp.Rating})
}
