package restaurants

import (
	"encoding/json"
	"net/http"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/authentication"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	"github.com/rs/zerolog/log"
)

// RestaurantRegistrationRESTAdaptor exposes restaurant registration over a REST API.
type RestaurantRegistrationRESTAdaptor struct {
	registrar RestaurantRegistrationService
}

// NewRestaurantRegistrationRESTAdaptor returns a new RestaurantRegistrationRESTAdaptor.
func NewRestaurantRegistrationRESTAdaptor(registrar RestaurantRegistrationService) *RestaurantRegistrationRESTAdaptor {
	return &RestaurantRegistrationRESTAdaptor{registrar: registrar}
}

type RegisterRestaurantRESTRequest struct {
	Name  string `json:"name"`
	City  string `json:"city"`
	Image string `json:"image"`
}

type RegisterRestaurantRESTResponse struct {
	Restaurant Restaurant `json:"restaurant"`
}

func (a *RestaurantRegistrationRESTAdaptor) RegisterRestaurant(w http.ResponseWriter, r *http.Request) {
	claim, ok := authentication.LoginClaimFromContext(r.Context())
	if !ok {
		log.Ctx(r.Context()).Warn().Msg("no login claim in context")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	var request RegisterRestaurantRESTRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to decode register restaurant request")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	resp, err := a.registrar.RegisterRestaurant(r.Context(), RegisterRestaurantRequest{
		UserID: claim.UserID,
		Name:   request.Name,
		City:   request.City,
		Image:  request.Image,
	})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to register restaurant")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterRestaurantRESTResponse{Restaurant: resp.Restaurant})
}
