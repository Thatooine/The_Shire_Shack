package restaurants

import (
	"encoding/json"
	"net/http"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/authentication"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// DishWriterServiceRESTAdaptor exposes dish write operations over a REST API.
type DishWriterServiceRESTAdaptor struct {
	writer DishWriterService
}

// NewDishWriterServiceRESTAdaptor returns a new DishWriterServiceRESTAdaptor.
func NewDishWriterServiceRESTAdaptor(writer DishWriterService) *DishWriterServiceRESTAdaptor {
	return &DishWriterServiceRESTAdaptor{writer: writer}
}

// CreateDish

type CreateDishRESTRequest struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	RestaurantID string  `json:"restaurant_id"`
	Image        string  `json:"image"`
}

type CreateDishRESTResponse struct {
	Dish Dish `json:"dish"`
}

func (a *DishWriterServiceRESTAdaptor) CreateDish(w http.ResponseWriter, r *http.Request) {
	claim, ok := authentication.LoginClaimFromContext(r.Context())
	if !ok {
		log.Ctx(r.Context()).Warn().Msg("no login claim in context")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	var request CreateDishRESTRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to decode create dish request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	resp, err := a.writer.CreateDish(r.Context(), CreateDishRequest{
		UserID:       claim.UserID,
		Name:         request.Name,
		Description:  request.Description,
		Price:        request.Price,
		RestaurantID: request.RestaurantID,
		Image:        request.Image,
	})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to create dish")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateDishRESTResponse{Dish: resp.Dish})
}

// UpdateDish

type UpdateDishRESTRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
}

type UpdateDishRESTResponse struct {
	Dish Dish `json:"dish"`
}

func (a *DishWriterServiceRESTAdaptor) UpdateDish(w http.ResponseWriter, r *http.Request) {
	claim, ok := authentication.LoginClaimFromContext(r.Context())
	if !ok {
		log.Ctx(r.Context()).Warn().Msg("no login claim in context")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	id := mux.Vars(r)["id"]

	var request UpdateDishRESTRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to decode update dish request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	resp, err := a.writer.UpdateDish(r.Context(), UpdateDishRequest{
		UserID:      claim.UserID,
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Image:       request.Image,
	})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to update dish")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UpdateDishRESTResponse{Dish: resp.Dish})
}

// DeleteDish

func (a *DishWriterServiceRESTAdaptor) DeleteDish(w http.ResponseWriter, r *http.Request) {
	claim, ok := authentication.LoginClaimFromContext(r.Context())
	if !ok {
		log.Ctx(r.Context()).Warn().Msg("no login claim in context")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	id := mux.Vars(r)["id"]

	if err := a.writer.DeleteDish(r.Context(), DeleteDishRequest{UserID: claim.UserID, ID: id}); err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to delete dish")
		errs.WriteHTTPError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
