package restaurants

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// DishReaderServiceRESTAdaptor exposes dish read operations over a REST API.
type DishReaderServiceRESTAdaptor struct {
	reader DishReaderService
}

// NewDishReaderServiceRESTAdaptor returns a new DishReaderServiceRESTAdaptor.
func NewDishReaderServiceRESTAdaptor(reader DishReaderService) *DishReaderServiceRESTAdaptor {
	return &DishReaderServiceRESTAdaptor{reader: reader}
}

// GetDish

type GetDishRESTResponse struct {
	Dish Dish `json:"dish"`
}

func (a *DishReaderServiceRESTAdaptor) GetDish(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	resp, err := a.reader.GetDish(r.Context(), GetDishRequest{ID: id})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to get dish")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetDishRESTResponse{Dish: resp.Dish})
}

// ListDishes

type ListDishesRESTResponse struct {
	Dishes []Dish `json:"dishes"`
	Total  int64  `json:"total"`
}

func (a *DishReaderServiceRESTAdaptor) ListDishes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	restaurantID := query.Get("restaurant_id")
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if limit == 0 {
		limit = 20
	}

	resp, err := a.reader.ListDishes(r.Context(), ListDishesRequest{
		RestaurantID: restaurantID,
		Offset:       offset,
		Limit:        limit,
	})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to list dishes")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListDishesRESTResponse{
		Dishes: resp.Dishes,
		Total:  resp.Total,
	})
}

// SearchDishes

type SearchDishesRESTResponse struct {
	Dishes []Dish `json:"dishes"`
	Total  int64  `json:"total"`
}

func (a *DishReaderServiceRESTAdaptor) SearchDishes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q := query.Get("q")
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if limit == 0 {
		limit = 20
	}

	resp, err := a.reader.SearchDishes(
		r.Context(),
		SearchDishesRequest{
			Query:  q,
			Offset: offset,
			Limit:  limit,
		})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to search dishes")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SearchDishesRESTResponse{
		Dishes: resp.Dishes,
		Total:  resp.Total,
	})
}
