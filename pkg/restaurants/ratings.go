package restaurants

import "time"

// Rating represents a user's review and score for a specific dish.
type Rating struct {
	// Unique identifier for this rating
	ID string `json:"id" bson:"id"`
	// ID of the dish being rated
	DishID string `json:"dish_id" bson:"dish_id"`
	// ID of the user who submitted the rating
	UserID string `json:"user_id" bson:"user_id"`

	// Score given to the dish (e.g. 1-5)
	Score int `json:"score" bson:"score"`
	// Optional written review
	Review string `json:"review" bson:"review"`

	// Timestamp of when the rating was submitted
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
