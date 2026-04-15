import { useEffect, useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { listDishes, listRatings, type Dish, type Rating } from "../api/dishes";
import { useAuth } from "../hooks/useAuth";
import DishCard from "../components/DishCard";
import ReviewModal from "../components/ReviewModal";
import RestaurantDropdown from "../components/RestaurantDropdown";
import OnboardRestaurantModal from "../components/OnboardRestaurantModal";
import ManageRestaurantModal from "../components/ManageRestaurantModal";

export default function HomePage() {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const [dishes, setDishes] = useState<Dish[]>([]);
  const [ratingsMap, setRatingsMap] = useState<Record<string, Rating[]>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [reviewDish, setReviewDish] = useState<Dish | null>(null);
  const [search, setSearch] = useState("");
  const [showOnboard, setShowOnboard] = useState(false);
  const [showManage, setShowManage] = useState(false);

  const fetchDishes = useCallback(() => {
    setLoading(true);
    listDishes()
      .then(async (res) => {
        setDishes(res.dishes);

        const ratingsResults = await Promise.all(
          res.dishes.map((dish) => listRatings(dish.id)),
        );
        const map: Record<string, Rating[]> = {};
        res.dishes.forEach((dish, i) => {
          map[dish.id] = ratingsResults[i].ratings;
        });
        setRatingsMap(map);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    fetchDishes();
  }, [fetchDishes]);

  const filtered = dishes.filter(
    (d) =>
      d.name.toLowerCase().includes(search.toLowerCase()) ||
      d.description.toLowerCase().includes(search.toLowerCase()),
  );

  function handleRatingSubmitted(rating: Rating) {
    setRatingsMap((prev) => ({
      ...prev,
      [rating.dish_id]: [...(prev[rating.dish_id] ?? []), rating],
    }));
  }

  async function handleSignOut() {
    await logout();
    navigate("/");
  }

  return (
    <div className="home-page">
      <header className="home-header">
        <div className="header-left">
          <span className="tavern-icon-sm">&#9876;</span>
          <h1>The Orc Shack</h1>
        </div>
        <div className="header-right">
          <RestaurantDropdown
            onRegister={() => setShowOnboard(true)}
            onManage={() => setShowManage(true)}
          />
          <button className="sign-out-btn" onClick={handleSignOut}>
            Sign Out
          </button>
        </div>
      </header>

      <div className="home-banner">
        <h2>Our Menu</h2>
        <p>Traditional dishes from the heart of Middle Earth</p>
        <input
          type="text"
          className="search-input"
          placeholder="Search dishes..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      <main className="dish-grid">
        {loading ? (
          <p className="no-results">Loading dishes...</p>
        ) : error ? (
          <p className="no-results">Failed to load dishes: {error}</p>
        ) : filtered.length > 0 ? (
          filtered.map((dish) => (
            <DishCard
              key={dish.id}
              dish={dish}
              ratings={ratingsMap[dish.id] ?? []}
              onReview={setReviewDish}
            />
          ))
        ) : (
          <p className="no-results">
            No dishes found. Even the Elves couldn't find that one.
          </p>
        )}
      </main>

      {reviewDish && (
        <ReviewModal
          dish={reviewDish}
          ratings={ratingsMap[reviewDish.id] ?? []}
          onClose={() => setReviewDish(null)}
          onSubmitted={handleRatingSubmitted}
        />
      )}

      {showOnboard && (
        <OnboardRestaurantModal onClose={() => setShowOnboard(false)} />
      )}

      {showManage && (
        <ManageRestaurantModal 
          onClose={() => setShowManage(false)} 
          onRegister={() => {
            setShowManage(false);
            setShowOnboard(true);
          }}
          onDishChanged={fetchDishes}
        />
      )}
    </div>
  );
}
