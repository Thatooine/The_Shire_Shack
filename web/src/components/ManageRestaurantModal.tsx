import { useEffect, useState, type FormEvent } from "react";
import { getMyRestaurant } from "../api/restaurants";
import {
  createDish,
  listDishesByRestaurant,
  updateDish,
  type Dish,
} from "../api/dishes";

type View = "list" | "edit" | "add";

interface Props {
  onClose: () => void;
  onRegister?: () => void;
  onDishChanged?: () => void;
}

export default function ManageRestaurantModal({ onClose, onRegister, onDishChanged }: Props) {
  const [restaurantId, setRestaurantId] = useState("");
  const [dishes, setDishes] = useState<Dish[]>([]);
  const [editing, setEditing] = useState<Dish | null>(null);
  const [view, setView] = useState<View>("list");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [noRestaurant, setNoRestaurant] = useState(false);

  useEffect(() => {
    getMyRestaurant()
      .then((res) => {
        if (!res || !res.restaurant) {
          throw new Error("No restaurant");
        }
        setRestaurantId(res.restaurant.id);
        return listDishesByRestaurant(res.restaurant.id);
      })
      .then((res) => setDishes(res.dishes ?? []))
      .catch((err) => {
        const msg = err.message ? err.message.toLowerCase() : "";
        if (msg.includes("not found") || msg === "no restaurant" || msg.includes("failed")) {
          setNoRestaurant(true);
        } else {
          setError(err.message);
        }
      })
      .finally(() => setLoading(false));
  }, []);

  async function handleUpdate(updated: Dish) {
    try {
      const res = await updateDish(updated.id, {
        name: updated.name,
        description: updated.description,
        price: updated.price,
        image: updated.image,
      });
      setDishes((prev) =>
        prev.map((d) => (d.id === res.dish.id ? res.dish : d)),
      );
      setEditing(null);
      setView("list");
      onDishChanged?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save dish");
    }
  }

  async function handleCreate(data: {
    name: string;
    description: string;
    price: number;
    image: string;
  }) {
    try {
      const res = await createDish({ ...data, restaurant_id: restaurantId });
      setDishes((prev) => [...prev, res.dish]);
      setView("list");
      onDishChanged?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to add dish");
    }
  }

  function goBack() {
    setError("");
    setEditing(null);
    setView("list");
  }

  return (
    <div className="fullscreen-modal-overlay" onClick={onClose}>
      <div
        className="fullscreen-modal fullscreen-modal--wide"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="fullscreen-modal-header">
          <h2>My Restaurant &mdash; Dishes</h2>
          <button className="modal-close" onClick={onClose}>
            &times;
          </button>
        </div>

        <div className="fullscreen-modal-body">
          {loading ? (
            <p className="no-results">Loading dishes...</p>
          ) : noRestaurant ? (
            <div className="onboard-success">
              <span className="success-icon">&#9876;</span>
              <h3>No Restaurant Found</h3>
              <p>You haven't registered a restaurant yet. Register one to start managing your dishes and serving customers.</p>
              <button 
                className="auth-btn" 
                onClick={() => {
                  onClose();
                  if (onRegister) onRegister();
                }}
              >
                Register a Restaurant
              </button>
            </div>
          ) : error && view === "list" ? (
            <p className="auth-error">{error}</p>
          ) : view === "edit" && editing ? (
            <DishForm
              title={`Edit: ${editing.name}`}
              initial={editing}
              error={error}
              submitLabel="Save Changes"
              onSubmit={(data) =>
                handleUpdate({ ...editing, ...data })
              }
              onCancel={goBack}
            />
          ) : view === "add" ? (
            <DishForm
              title="Add New Dish"
              error={error}
              submitLabel="Add Dish"
              onSubmit={handleCreate}
              onCancel={goBack}
            />
          ) : (
            <>
              <div className="manage-actions">
                <button
                  className="auth-btn"
                  onClick={() => {
                    setError("");
                    setView("add");
                  }}
                >
                  + Add Dish
                </button>
              </div>
              {dishes.length === 0 ? (
                <p className="no-results">
                  No dishes yet. Add dishes to your restaurant to see them here.
                </p>
              ) : (
                <div className="manage-table-wrapper">
                  <table className="manage-table">
                    <thead>
                      <tr>
                        <th></th>
                        <th>Name</th>
                        <th>Description</th>
                        <th>Price</th>
                        <th></th>
                      </tr>
                    </thead>
                    <tbody>
                      {dishes.map((dish) => (
                        <tr
                          key={dish.id}
                          onClick={() => {
                            setEditing(dish);
                            setView("edit");
                          }}
                        >
                          <td>
                            <img
                              src={dish.image}
                              alt={dish.name}
                              className="manage-dish-thumb"
                            />
                          </td>
                          <td className="manage-dish-name">{dish.name}</td>
                          <td className="manage-dish-desc">
                            {dish.description}
                          </td>
                          <td className="manage-dish-price">
                            R{dish.price.toFixed(2)}
                          </td>
                          <td>
                            <button
                              className="manage-edit-btn"
                              onClick={(e) => {
                                e.stopPropagation();
                                setEditing(dish);
                                setView("edit");
                              }}
                            >
                              Edit
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}

function DishForm({
  title,
  initial,
  error,
  submitLabel,
  onSubmit,
  onCancel,
}: {
  title: string;
  initial?: { name: string; description: string; price: number; image: string };
  error: string;
  submitLabel: string;
  onSubmit: (data: {
    name: string;
    description: string;
    price: number;
    image: string;
  }) => void;
  onCancel: () => void;
}) {
  const [name, setName] = useState(initial?.name ?? "");
  const [description, setDescription] = useState(initial?.description ?? "");
  const [price, setPrice] = useState(initial?.price.toString() ?? "");
  const [image, setImage] = useState(initial?.image ?? "");
  const [localError, setLocalError] = useState("");

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!name.trim()) {
      setLocalError("Dish name is required.");
      return;
    }
    const parsed = parseFloat(price);
    if (isNaN(parsed) || parsed <= 0) {
      setLocalError("Enter a valid price.");
      return;
    }
    setLocalError("");
    onSubmit({
      name: name.trim(),
      description: description.trim(),
      price: parsed,
      image: image.trim(),
    });
  }

  const displayError = localError || error;

  return (
    <div className="edit-dish-form-wrapper">
      <button className="back-btn" onClick={onCancel}>
        &larr; Back to dishes
      </button>
      <h3>{title}</h3>
      <form onSubmit={handleSubmit} className="onboard-form">
        <div className="form-group">
          <label htmlFor="dish-name">Dish Name</label>
          <input
            id="dish-name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label htmlFor="dish-desc">Description</label>
          <textarea
            id="dish-desc"
            rows={3}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label htmlFor="dish-price">Price (R)</label>
          <input
            id="dish-price"
            type="number"
            step="0.01"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
          />
        </div>

        <div className="form-group">
          <label htmlFor="dish-image">Image URL</label>
          <input
            id="dish-image"
            type="text"
            value={image}
            onChange={(e) => setImage(e.target.value)}
          />
        </div>

        {displayError && <p className="auth-error">{displayError}</p>}

        <div className="edit-actions">
          <button type="button" className="cancel-btn" onClick={onCancel}>
            Cancel
          </button>
          <button type="submit" className="auth-btn">
            {submitLabel}
          </button>
        </div>
      </form>
    </div>
  );
}
