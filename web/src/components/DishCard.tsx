import type { Dish, Rating } from "../api/dishes";

interface Props {
  dish: Dish;
  ratings: Rating[];
  onReview: (dish: Dish) => void;
}

export default function DishCard({ dish, ratings, onReview }: Props) {
  const avgScore =
    ratings.length > 0
      ? ratings.reduce((sum, r) => sum + r.score, 0) / ratings.length
      : 0;

  return (
    <div className="dish-card">
      <div className="dish-img-wrapper">
        <img src={dish.image} alt={dish.name} className="dish-img" />
      </div>
      <div className="dish-body">
        <div className="dish-header">
          <h3 className="dish-name">{dish.name}</h3>
          <span className="dish-price">R{dish.price.toFixed(2)}</span>
        </div>
        <p className="dish-desc">{dish.description}</p>
        <div className="dish-footer">
          <span className="dish-rating">
            {avgScore > 0 ? (
              <>
                <span className="star-display">
                  {"★".repeat(Math.round(avgScore))}
                  {"☆".repeat(5 - Math.round(avgScore))}
                </span>
                <span className="review-count">
                  ({ratings.length})
                </span>
              </>
            ) : (
              <span className="no-reviews">No reviews yet</span>
            )}
          </span>
          <button className="review-btn" onClick={() => onReview(dish)}>
            Write a Review
          </button>
        </div>
      </div>
    </div>
  );
}
