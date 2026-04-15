package restaurants

// Restaurant represents a dining establishment.
type Restaurant struct {
	// Unique identifier for this restaurant
	ID string `json:"id" bson:"id"`

	// ID of the user who owns this restaurant
	OwnerID string `json:"ownerID" bson:"ownerID"`

	// Name of the restaurant
	Name string `json:"name" bson:"name"`

	// URL or path to the dish image
	Image string `json:"image" bson:"image"`

	// City where the restaurant is located
	City string `json:"city" bson:"city"`
}
