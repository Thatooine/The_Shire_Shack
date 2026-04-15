package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// UserReaderServiceRESTAdaptor exposes user read operations over a REST API.
type UserReaderServiceRESTAdaptor struct {
	reader UserReaderService
}

// NewUserReaderServiceRESTAdaptor returns a new UserReaderServiceRESTAdaptor.
func NewUserReaderServiceRESTAdaptor(reader UserReaderService) *UserReaderServiceRESTAdaptor {
	return &UserReaderServiceRESTAdaptor{reader: reader}
}

// GetUser

type GetUserRESTResponse struct {
	User User `json:"user"`
}

func (a *UserReaderServiceRESTAdaptor) GetUser(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]

	resp, err := a.reader.GetUser(r.Context(), GetUserRequest{Email: email})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to get user")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetUserRESTResponse{User: resp.User})
}

// ListUsers

type ListUsersRESTResponse struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
}

func (a *UserReaderServiceRESTAdaptor) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if limit == 0 {
		limit = 20
	}

	resp, err := a.reader.ListUsers(r.Context(), ListUsersRequest{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to list users")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListUsersRESTResponse{
		Users: resp.Users,
		Total: resp.Total,
	})
}

// SearchUsers

type SearchUsersRESTResponse struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
}

func (a *UserReaderServiceRESTAdaptor) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q := query.Get("q")
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if limit == 0 {
		limit = 20
	}

	resp, err := a.reader.SearchUsers(
		r.Context(),
		SearchUsersRequest{
			Query:  q,
			Offset: offset,
			Limit:  limit,
		})
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to search users")
		errs.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SearchUsersRESTResponse{
		Users: resp.Users,
		Total: resp.Total,
	})
}
