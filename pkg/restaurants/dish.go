package restaurants

// Dish represents a menu item offered by a restaurant.
type Dish struct {
	// Unique identifier for this dish
	ID string `json:"id" bson:"id"`
	// ID of the restaurant that offers this dish
	RestaurantID string `json:"restaurant_id" bson:"restaurant_id"`
	// Name of the dish
	Name string `json:"name" bson:"name"`
	// Description of the dish
	Description string `json:"description" bson:"description"`
	// Price in the restaurant's local currency
	Price float64 `json:"price" bson:"price"`
	// URL or path to the dish image
	Image string `json:"image" bson:"image"`
}
