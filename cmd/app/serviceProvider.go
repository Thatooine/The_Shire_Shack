package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	firebase "firebase.google.com/go/v4"
	authenticationImpl "github.com/bash/the-dancing-pony-v2-rnyfbr/internal/pkg/authentication"
	rateLimitingImpl "github.com/bash/the-dancing-pony-v2-rnyfbr/internal/pkg/rateLimiting"
	restaurantsImpl "github.com/bash/the-dancing-pony-v2-rnyfbr/internal/pkg/restaurants"
	usersImpl "github.com/bash/the-dancing-pony-v2-rnyfbr/internal/pkg/users"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/authentication"
	pkgMongo "github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/mongo"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/rateLimiting"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/restaurants"
	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/users"
	"github.com/go-jose/go-jose/v4"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"
)

type ServiceProviders struct {
	UserCreator                          users.UserCreatorService
	UserReaderService                    users.UserReaderService
	DishWriterService                    restaurants.DishWriterService
	DishReaderService                    restaurants.DishReaderService
	RatingSubmitterService               restaurants.RatingSubmitterService
	RatingReaderService                  restaurants.RatingReaderService
	RestaurantReaderService              restaurants.RestaurantReaderService
	RestaurantRegistrationService        restaurants.RestaurantRegistrationService
	FirebaseAuthenticatorService         authentication.FirebaseAuthenticatorService
	EmailAndPasswordAuthenticatorService authentication.EmailAndPasswordAuthenticatorService
	UserRegistrationService              users.UserRegistrationService
	AccessTokenCreatorService            authentication.AccessTokenCreatorService
	AccessTokenValidatorService          authentication.AccessTokenValidatorService
	RateLimiter                          rateLimiting.RedisTokenBucketRateLimiter
}

func NewServiceProviders(ctx context.Context, conf *Config, secureConf *SecureConfig) (*ServiceProviders, error) {
	mongoClient, err := pkgMongo.NewClient(ctx, secureConf.MongoURI)
	if err != nil {
		return nil, err
	}

	const DatabaseName = "shire_shack"

	dishStore := pkgMongo.NewStore(mongoClient, DatabaseName, "dishes")
	ratingsStore := pkgMongo.NewStore(mongoClient, DatabaseName, "ratings")
	restaurantStore := pkgMongo.NewStore(mongoClient, DatabaseName, "restaurants")
	usersStore := pkgMongo.NewStore(mongoClient, DatabaseName, "users")

	userNewCreator := usersImpl.NewUserCreatorServiceImpl(usersStore)
	userReader := usersImpl.NewUserReaderServiceImpl(usersStore)
	dishWriter := restaurantsImpl.NewDishWriterServiceImpl(dishStore, usersStore, restaurantStore)
	dishReader := restaurantsImpl.NewDishReaderServiceImpl(dishStore)
	ratingSubmitter := restaurantsImpl.NewRatingSubmitterServiceImpl(ratingsStore)
	ratingReader := restaurantsImpl.NewRatingReaderServiceImpl(ratingsStore)
	restaurantReader := restaurantsImpl.NewRestaurantReaderServiceImpl(restaurantStore)
	restaurantRegistrar := restaurantsImpl.NewRestaurantRegistrationServiceImpl(restaurantStore, usersStore)

	block, _ := pem.Decode([]byte(secureConf.JWTPrivateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode JWT private key PEM")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT private key: %w", err)
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("JWT private key is not RSA")
	}

	tokenSigner, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: privateKey},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create token signer: %w", err)
	}

	// Initialize Firebase app using service account credentials file.
	firebaseApp, err := firebase.NewApp(ctx,
		&firebase.Config{ProjectID: "bash-interview-project"},
		option.WithCredentialsFile(conf.FirebaseServiceAccountPath),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	// redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: secureConf.RedisURI,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	rateLimiter := rateLimitingImpl.NewRedisRateLimiterImpl(redisClient)

	// authentication services
	accessTokenCreator := authenticationImpl.NewAccessTokenCreatorServiceImpl(tokenSigner)
	accessTokenValidator := authenticationImpl.NewAccessTokenValidatorServiceImpl(&privateKey.PublicKey)
	registrationService := usersImpl.NewUserRegistrationServiceImpl(firebaseApp, accessTokenCreator, userNewCreator, userReader, secureConf.FirebaseWebAPIKey)
	firebaseAuthenticator := authenticationImpl.NewFirebaseAuthenticatorService(firebaseApp, accessTokenCreator, userReader)
	emailPasswordAuthenticator := authenticationImpl.NewEmailAndPasswordAuthenticatorService(accessTokenCreator, userReader, secureConf.FirebaseWebAPIKey)

	return &ServiceProviders{
		UserCreator:                          userNewCreator,
		UserReaderService:                    userReader,
		DishWriterService:                    dishWriter,
		DishReaderService:                    dishReader,
		RatingSubmitterService:               ratingSubmitter,
		RatingReaderService:                  ratingReader,
		RestaurantReaderService:              restaurantReader,
		RestaurantRegistrationService:        restaurantRegistrar,
		UserRegistrationService:              registrationService,
		FirebaseAuthenticatorService:         firebaseAuthenticator,
		EmailAndPasswordAuthenticatorService: emailPasswordAuthenticator,
		AccessTokenCreatorService:            accessTokenCreator,
		AccessTokenValidatorService:          accessTokenValidator,
		RateLimiter:                          rateLimiter,
	}, nil
}
