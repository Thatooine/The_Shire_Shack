package restaurants

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// RatingReaderServiceRESTAdaptor exposes rating read operations over a REST API.
type RatingReaderServiceRESTAdaptor struct {
	reader RatingReaderService
}

// NewRatingReaderServiceRESTAdaptor returns a new RatingReaderServiceRESTAdaptor.
func NewRatingReaderServiceRESTAdaptor(reader RatingReaderService) *RatingReaderServiceRESTAdaptor {
	return &RatingReaderServiceRESTAdaptor{reader: reader}
}

// ListRatings

type ListRatingsRESTResponse struct {
	Ratings []Rating `json:"ratings"`
	Total   int64    `json:"total"`
}

func (a *RatingReaderServiceRESTAdaptor) ListRatings(w http.ResponseWriter, r *http.Request) {
	dishID := mux.Vars(r)["id"]
	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if limit == 0 {
		limit = 20
	}

	resp, err := a.reader.ListRatings(r.Context(), ListRatingsRequest{
		DishID: dishID,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to list ratings")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListRatingsRESTResponse{
		Ratings: resp.Ratings,
		Total:   resp.Total,
	})
}
