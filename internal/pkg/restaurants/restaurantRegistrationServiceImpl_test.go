package restaurants

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/errs"
	pkgMongo "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/mongo"
	pkgRestaurants "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/restaurants"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/users"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func newTestRestaurantRegistrationService(restaurantStore pkgMongo.Storer, userStore pkgMongo.Storer) *RestaurantRegistrationServiceImpl {
	return &RestaurantRegistrationServiceImpl{
		restaurantStore: restaurantStore,
		userStore:       userStore,
	}
}

func TestRegisterRestaurant_ValidationFails(t *testing.T) {
	svc := newTestRestaurantRegistrationService(nil, nil)

	// Missing name and city
	_, err := svc.RegisterRestaurant(context.Background(), pkgRestaurants.RegisterRestaurantRequest{
		UserID: "user-1",
	})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestRegisterRestaurant_UserNotFound(t *testing.T) {
	userStore := &pkgMongo.StoreMock{
		GetFn: func(_ context.Context, id string, _ interface{}) error {
			return fmt.Errorf("user not found")
		},
	}

	svc := newTestRestaurantRegistrationService(nil, userStore)

	_, err := svc.RegisterRestaurant(context.Background(), pkgRestaurants.RegisterRestaurantRequest{
		UserID: "user-1",
		Name:   "The Prancing Pony",
		City:   "Bree",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRegisterRestaurant_AlreadyFullyRegistered(t *testing.T) {
	userStore := &pkgMongo.StoreMock{
		GetFn: func(_ context.Context, _ string, result interface{}) error {
			u := result.(*users.User)
			*u = users.User{
				ID:    "user-1",
				Name:  "Aragorn",
				Email: "aragorn@gondor.com",
				Roles: []users.Role{users.RoleCustomer, users.RoleRestaurantOwner},
			}
			return nil
		},
	}

	restaurantStore := &pkgMongo.StoreMock{
		ListFn: func(_ context.Context, _ bson.M, _ int64, _ int64, results interface{}) (int64, error) {
			r := results.(*[]pkgRestaurants.Restaurant)
			*r = []pkgRestaurants.Restaurant{
				{ID: "rest-1", OwnerID: "user-1", Name: "Existing Restaurant", City: "Bree"},
			}
			return 1, nil
		},
	}

	svc := newTestRestaurantRegistrationService(restaurantStore, userStore)

	_, err := svc.RegisterRestaurant(context.Background(), pkgRestaurants.RegisterRestaurantRequest{
		UserID: "user-1",
		Name:   "New Restaurant",
		City:   "Rohan",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errs.ErrConflict) {
		t.Errorf("expected ErrConflict, got: %v", err)
	}
}

func TestRegisterRestaurant_RecoveryFromPartialFailure(t *testing.T) {
	// Restaurant exists but user lacks the RestaurantOwner role
	userStore := &pkgMongo.StoreMock{
		GetFn: func(_ context.Context, _ string, result interface{}) error {
			u := result.(*users.User)
			*u = users.User{
				ID:    "user-1",
				Name:  "Aragorn",
				Email: "aragorn@gondor.com",
				Roles: []users.Role{users.RoleCustomer}, // Missing RestaurantOwner
			}
			return nil
		},
		UpdateFn: func(_ context.Context, id string, update bson.M, result interface{}) error {
			if id != "user-1" {
				t.Fatalf("expected update for user-1, got %s", id)
			}
			roles, ok := update["roles"].([]users.Role)
			if !ok {
				t.Fatal("expected roles in update")
			}
			hasOwnerRole := false
			for _, r := range roles {
				if r == users.RoleRestaurantOwner {
					hasOwnerRole = true
				}
			}
			if !hasOwnerRole {
				t.Fatal("expected RestaurantOwner role in update")
			}
			return nil
		},
	}

	existingRestaurant := pkgRestaurants.Restaurant{
		ID: "rest-1", OwnerID: "user-1", Name: "The Prancing Pony", City: "Bree",
	}

	restaurantStore := &pkgMongo.StoreMock{
		ListFn: func(_ context.Context, _ bson.M, _ int64, _ int64, results interface{}) (int64, error) {
			r := results.(*[]pkgRestaurants.Restaurant)
			*r = []pkgRestaurants.Restaurant{existingRestaurant}
			return 1, nil
		},
	}

	svc := newTestRestaurantRegistrationService(restaurantStore, userStore)

	resp, err := svc.RegisterRestaurant(context.Background(), pkgRestaurants.RegisterRestaurantRequest{
		UserID: "user-1",
		Name:   "Ignored Name",
		City:   "Ignored City",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should return the existing restaurant, not create a new one
	if resp.Restaurant.ID != "rest-1" {
		t.Errorf("expected existing restaurant ID rest-1, got %s", resp.Restaurant.ID)
	}
	if resp.Restaurant.Name != "The Prancing Pony" {
		t.Errorf("expected existing restaurant name, got %s", resp.Restaurant.Name)
	}
}

func TestRegisterRestaurant_NewRegistration(t *testing.T) {
	var storedRestaurant pkgRestaurants.Restaurant
	roleUpdated := false

	userStore := &pkgMongo.StoreMock{
		GetFn: func(_ context.Context, _ string, result interface{}) error {
			u := result.(*users.User)
			*u = users.User{
				ID:    "user-1",
				Name:  "Aragorn",
				Email: "aragorn@gondor.com",
				Roles: []users.Role{users.RoleCustomer},
			}
			return nil
		},
		UpdateFn: func(_ context.Context, id string, update bson.M, _ interface{}) error {
			if id != "user-1" {
				t.Fatalf("expected update for user-1, got %s", id)
			}
			roleUpdated = true
			return nil
		},
	}

	restaurantStore := &pkgMongo.StoreMock{
		ListFn: func(_ context.Context, _ bson.M, _ int64, _ int64, results interface{}) (int64, error) {
			// No existing restaurants
			r := results.(*[]pkgRestaurants.Restaurant)
			*r = []pkgRestaurants.Restaurant{}
			return 0, nil
		},
		PutFn: func(_ context.Context, document interface{}) (string, error) {
			restaurant, ok := document.(pkgRestaurants.Restaurant)
			if !ok {
				t.Fatal("expected Restaurant document")
			}
			storedRestaurant = restaurant
			return restaurant.ID, nil
		},
	}

	svc := newTestRestaurantRegistrationService(restaurantStore, userStore)

	resp, err := svc.RegisterRestaurant(context.Background(), pkgRestaurants.RegisterRestaurantRequest{
		UserID: "user-1",
		Name:   "The Green Dragon",
		City:   "Hobbiton",
		Image:  "dragon.jpg",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify restaurant was created with correct fields
	if resp.Restaurant.Name != "The Green Dragon" {
		t.Errorf("expected name The Green Dragon, got %s", resp.Restaurant.Name)
	}
	if resp.Restaurant.City != "Hobbiton" {
		t.Errorf("expected city Hobbiton, got %s", resp.Restaurant.City)
	}
	if resp.Restaurant.OwnerID != "user-1" {
		t.Errorf("expected ownerID user-1, got %s", resp.Restaurant.OwnerID)
	}
	if resp.Restaurant.ID == "" {
		t.Error("expected a generated restaurant ID")
	}

	// Verify the restaurant was stored
	if !reflect.DeepEqual(storedRestaurant, resp.Restaurant) {
		t.Error("stored restaurant does not match response")
	}

	// Verify user role was updated
	if !roleUpdated {
		t.Error("expected user role to be updated")
	}
}

func TestRegisterRestaurant_StorePutFails(t *testing.T) {
	userStore := &pkgMongo.StoreMock{
		GetFn: func(_ context.Context, _ string, result interface{}) error {
			u := result.(*users.User)
			*u = users.User{ID: "user-1", Roles: []users.Role{users.RoleCustomer}}
			return nil
		},
	}

	restaurantStore := &pkgMongo.StoreMock{
		ListFn: func(_ context.Context, _ bson.M, _ int64, _ int64, results interface{}) (int64, error) {
			r := results.(*[]pkgRestaurants.Restaurant)
			*r = []pkgRestaurants.Restaurant{}
			return 0, nil
		},
		PutFn: func(_ context.Context, _ interface{}) (string, error) {
			return "", fmt.Errorf("database error")
		},
	}

	svc := newTestRestaurantRegistrationService(restaurantStore, userStore)

	_, err := svc.RegisterRestaurant(context.Background(), pkgRestaurants.RegisterRestaurantRequest{
		UserID: "user-1",
		Name:   "Test",
		City:   "Test",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRegisterRestaurant_RoleUpdateFails(t *testing.T) {
	userStore := &pkgMongo.StoreMock{
		GetFn: func(_ context.Context, _ string, result interface{}) error {
			u := result.(*users.User)
			*u = users.User{ID: "user-1", Roles: []users.Role{users.RoleCustomer}}
			return nil
		},
		UpdateFn: func(_ context.Context, _ string, _ bson.M, _ interface{}) error {
			return fmt.Errorf("update failed")
		},
	}

	restaurantStore := &pkgMongo.StoreMock{
		ListFn: func(_ context.Context, _ bson.M, _ int64, _ int64, results interface{}) (int64, error) {
			r := results.(*[]pkgRestaurants.Restaurant)
			*r = []pkgRestaurants.Restaurant{}
			return 0, nil
		},
		PutFn: func(_ context.Context, doc interface{}) (string, error) {
			return "rest-1", nil
		},
	}

	svc := newTestRestaurantRegistrationService(restaurantStore, userStore)

	_, err := svc.RegisterRestaurant(context.Background(), pkgRestaurants.RegisterRestaurantRequest{
		UserID: "user-1",
		Name:   "Test",
		City:   "Test",
	})
	if err == nil {
		t.Fatal("expected error when role update fails, got nil")
	}
}
