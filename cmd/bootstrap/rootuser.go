package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bash/the-dancing-pony-v2-rnyfbr/pkg/users"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func createRootUser(ctx context.Context, client *mongo.Client) {
	collection := client.Database("shire_shack").Collection("users")

	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		log.Fatalf("failed to clear users collection: %v", err)
	}
	fmt.Println("cleared users collection")

	rootUser := users.User{
		ID:    "00000000-0000-0000-0000-000000000000",
		Name:  "Root User",
		Email: "root+user@gmail.com",
		Roles: []users.Role{users.RoleAdmin, users.RoleRestaurantOwner},
	}

	if _, err := collection.InsertOne(ctx, rootUser); err != nil {
		log.Fatalf("failed to insert root user: %v", err)
	}

	fmt.Printf("root user created: %s (%s)\n", rootUser.Email, rootUser.ID)
}
