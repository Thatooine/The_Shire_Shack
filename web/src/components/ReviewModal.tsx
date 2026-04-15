import { useState, type FormEvent } from "react";
import { submitRating, type Dish, type Rating } from "../api/dishes";

interface Props {
  dish: Dish;
  ratings: Rating[];
  onClose: () => void;
  onSubmitted: (rating: Rating) => void;
}

export default function ReviewModal({ dish, ratings, onClose, onSubmitted }: Props) {
  const [rating, setRating] = useState(5);
  const [comment, setComment] = useState("");
  const [hoveredStar, setHoveredStar] = useState(0);
  const [submitting, setSubmitting] = useState(false);
  const [submitted, setSubmitted] = useState<Rating | null>(null);
  const [error, setError] = useState("");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!comment.trim()) return;

    setSubmitting(true);
    setError("");
    try {
      const res = await submitRating(dish.id, rating, comment.trim());
      setSubmitted(res.rating);
      onSubmitted(res.rating);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to submit review");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose}>
          &times;
        </button>
        <h2>Review: {dish.name}</h2>

        {submitted ? (
          <div className="review-success">
            <p className="success-text">Review submitted!</p>
            <div className="review-item">
              <div className="review-meta">
                <span className="review-stars">
                  {"★".repeat(submitted.score)}
                  {"☆".repeat(5 - submitted.score)}
                </span>
              </div>
              <p>{submitted.review}</p>
            </div>
            <button className="auth-btn" onClick={onClose}>
              Close
            </button>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="review-form">
            <label>Rating</label>
            <div className="star-picker">
              {[1, 2, 3, 4, 5].map((star) => (
                <button
                  key={star}
                  type="button"
                  className={`star ${star <= (hoveredStar || rating) ? "filled" : ""}`}
                  onClick={() => setRating(star)}
                  onMouseEnter={() => setHoveredStar(star)}
                  onMouseLeave={() => setHoveredStar(0)}
                >
                  &#9733;
                </button>
              ))}
            </div>

            <label htmlFor="review-comment">Your Review</label>
            <textarea
              id="review-comment"
              rows={4}
              placeholder="Share your thoughts on this dish..."
              value={comment}
              onChange={(e) => setComment(e.target.value)}
            />

            {error && <p className="error-text">{error}</p>}
            <button type="submit" className="auth-btn" disabled={submitting}>
              {submitting ? "Submitting..." : "Submit Review"}
            </button>
          </form>
        )}

        {ratings.length > 0 && (
          <div className="existing-reviews">
            <h3>Reviews</h3>
            {ratings.map((r) => (
              <div key={r.id} className="review-item">
                <div className="review-meta">
                  <span className="review-stars">
                    {"★".repeat(r.score)}
                    {"☆".repeat(5 - r.score)}
                  </span>
                </div>
                <p>{r.review}</p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
