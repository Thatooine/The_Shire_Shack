import { useEffect, useState, type FormEvent } from "react";
import { getMyRestaurant, registerRestaurant, type Restaurant } from "../api/restaurants";

type Step = "details" | "plan" | "confirm";

interface Plan {
  id: string;
  name: string;
  price: string;
  period: string;
  features: string[];
  tag?: string;
}

const PLANS: Plan[] = [
  {
    id: "trial",
    name: "Free Trial",
    price: "R0",
    period: "14 days",
    tag: "No card required",
    features: [
      "Up to 10 menu items",
      "Basic analytics dashboard",
      "Standard listing placement",
      "Email support",
    ],
  },
  {
    id: "basic",
    name: "Basic",
    price: "R299",
    period: "/month",
    tag: "Most popular",
    features: [
      "Up to 50 menu items",
      "Advanced analytics & reports",
      "Priority listing placement",
      "Rating response tools",
      "Email & chat support",
    ],
  },
  {
    id: "premium",
    name: "Premium",
    price: "R799",
    period: "/month",
    features: [
      "Unlimited menu items",
      "Real-time analytics & insights",
      "Featured listing & promotions",
      "Bulk dish management",
      "Custom branding",
      "Dedicated account manager",
    ],
  },
];

interface Props {
  onClose: () => void;
}

export default function OnboardRestaurantModal({ onClose }: Props) {
  const [step, setStep] = useState<Step>("details");
  const [name, setName] = useState("");
  const [city, setCity] = useState("");
  const [image, setImage] = useState("");
  const [selectedPlan, setSelectedPlan] = useState<Plan>(PLANS[0]);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [existing, setExisting] = useState<Restaurant | null>(null);
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    getMyRestaurant()
      .then((res) => setExisting(res.restaurant))
      .catch(() => {})
      .finally(() => setChecking(false));
  }, []);

  function handleDetailsSubmit(e: FormEvent) {
    e.preventDefault();
    if (!name.trim() || !city.trim()) {
      setError("Name and city are required.");
      return;
    }
    setError("");
    setStep("plan");
  }

  function handlePlanContinue() {
    setStep("confirm");
  }

  async function handleConfirm() {
    setSubmitting(true);
    setError("");
    try {
      await registerRestaurant(name.trim(), city.trim(), image.trim());
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
      setSubmitting(false);
    }
  }

  const trialEnd = new Date();
  trialEnd.setDate(trialEnd.getDate() + 14);
  const trialEndStr = trialEnd.toLocaleDateString("en-ZA", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });

  const stepTitles: Record<Step, string> = {
    details: "Register Your Restaurant",
    plan: "Choose Your Plan",
    confirm: "You're All Set!",
  };

  return (
    <div className="fullscreen-modal-overlay" onClick={onClose}>
      <div
        className={`fullscreen-modal ${step === "plan" ? "fullscreen-modal--wide" : ""}`}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="fullscreen-modal-header">
          {!checking && !existing && (
            <div className="stepper-header">
              <h2>{stepTitles[step]}</h2>
              {step !== "confirm" && (
                <div className="step-indicators">
                  <span className={`step-dot ${step === "details" ? "active" : "done"}`}>1</span>
                  <span className="step-line" />
                  <span className={`step-dot ${step === "plan" ? "active" : ""}`}>2</span>
                  <span className="step-line" />
                  <span className="step-dot">3</span>
                </div>
              )}
            </div>
          )}
          <button className="modal-close" onClick={onClose}>
            &times;
          </button>
        </div>

        <div className="fullscreen-modal-body">
          {checking ? (
            <p className="no-results">Loading...</p>
          ) : existing ? (
            <div className="confirm-step">
              <div className="confirm-icon-ring">
                <span className="confirm-icon">&#10003;</span>
              </div>
              <h3>You already have a registered restaurant</h3>
              <div className="confirm-details">
                <div className="confirm-row">
                  <span className="confirm-label">Restaurant</span>
                  <span className="confirm-value">{existing.name}</span>
                </div>
                <div className="confirm-row">
                  <span className="confirm-label">Location</span>
                  <span className="confirm-value">{existing.city}</span>
                </div>
              </div>
              <div className="confirm-actions">
                <button className="auth-btn" onClick={onClose}>
                  Close
                </button>
              </div>
            </div>
          ) : (
            <>
              {step === "details" && (
                <form onSubmit={handleDetailsSubmit} className="onboard-form">
                  <div className="form-group">
                    <label htmlFor="restaurant-name">Restaurant Name</label>
                    <input
                      id="restaurant-name"
                      type="text"
                      placeholder="e.g. The Prancing Pony"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                    />
                  </div>

                  <div className="form-group">
                    <label htmlFor="restaurant-city">City</label>
                    <input
                      id="restaurant-city"
                      type="text"
                      placeholder="e.g. Bree"
                      value={city}
                      onChange={(e) => setCity(e.target.value)}
                    />
                  </div>

                  <div className="form-group">
                    <label htmlFor="restaurant-image">Image URL</label>
                    <input
                      id="restaurant-image"
                      type="text"
                      placeholder="https://example.com/image.jpg"
                      value={image}
                      onChange={(e) => setImage(e.target.value)}
                    />
                  </div>

                  {error && <p className="auth-error">{error}</p>}

                  <button type="submit" className="auth-btn">
                    Continue to Plans
                  </button>
                </form>
              )}

              {step === "plan" && (
            <div className="plan-step">
              <p className="plan-subtitle">
                Start with a 14-day free trial. Upgrade anytime.
              </p>
              <div className="plan-cards">
                {PLANS.map((plan) => (
                  <div
                    key={plan.id}
                    className={`plan-card ${selectedPlan.id === plan.id ? "selected" : ""}`}
                    onClick={() => setSelectedPlan(plan)}
                  >
                    {plan.tag && <span className="plan-tag">{plan.tag}</span>}
                    <h3 className="plan-name">{plan.name}</h3>
                    <div className="plan-pricing">
                      <span className="plan-price">{plan.price}</span>
                      <span className="plan-period">{plan.period}</span>
                    </div>
                    <ul className="plan-features">
                      {plan.features.map((f) => (
                        <li key={f}>
                          <span className="feature-check">&#10003;</span>
                          {f}
                        </li>
                      ))}
                    </ul>
                    <button
                      className={`plan-select-btn ${selectedPlan.id === plan.id ? "active" : ""}`}
                      onClick={(e) => {
                        e.stopPropagation();
                        setSelectedPlan(plan);
                      }}
                    >
                      {selectedPlan.id === plan.id ? "Selected" : "Select"}
                    </button>
                  </div>
                ))}
              </div>

              <div className="plan-actions">
                <button className="back-btn" onClick={() => setStep("details")}>
                  &larr; Back
                </button>
                <button className="auth-btn" onClick={handlePlanContinue}>
                  Continue with {selectedPlan.name}
                </button>
              </div>
            </div>
          )}

          {step === "confirm" && (
            <div className="confirm-step">
              <div className="confirm-icon-ring">
                <span className="confirm-icon">&#10003;</span>
              </div>

              <div className="confirm-details">
                <div className="confirm-row">
                  <span className="confirm-label">Restaurant</span>
                  <span className="confirm-value">{name}</span>
                </div>
                <div className="confirm-row">
                  <span className="confirm-label">Location</span>
                  <span className="confirm-value">{city}</span>
                </div>
                <div className="confirm-row">
                  <span className="confirm-label">Plan</span>
                  <span className="confirm-value">{selectedPlan.name}</span>
                </div>
                {selectedPlan.id === "trial" && (
                  <div className="confirm-row">
                    <span className="confirm-label">Trial ends</span>
                    <span className="confirm-value">{trialEndStr}</span>
                  </div>
                )}
                {selectedPlan.id !== "trial" && (
                  <div className="confirm-row">
                    <span className="confirm-label">Billing</span>
                    <span className="confirm-value">
                      {selectedPlan.price}{selectedPlan.period} &mdash; starts after 14-day trial
                    </span>
                  </div>
                )}
              </div>

              {selectedPlan.id === "trial" && (
                <div className="trial-banner">
                  <span className="trial-banner-icon">&#9201;</span>
                  <div>
                    <strong>Your 14-day free trial starts now</strong>
                    <p>
                      Full access to {selectedPlan.name} features until{" "}
                      {trialEndStr}. No credit card needed.
                    </p>
                  </div>
                </div>
              )}

              {selectedPlan.id !== "trial" && (
                <div className="trial-banner">
                  <span className="trial-banner-icon">&#9201;</span>
                  <div>
                    <strong>14-day free trial included</strong>
                    <p>
                      You won't be charged until {trialEndStr}. Cancel anytime
                      before then.
                    </p>
                  </div>
                </div>
              )}

              {error && <p className="auth-error">{error}</p>}

              <div className="confirm-actions">
                {!submitting && (
                  <button className="back-btn" onClick={() => setStep("plan")}>
                    &larr; Change plan
                  </button>
                )}
                <button className="auth-btn" onClick={handleConfirm} disabled={submitting}>
                  {submitting ? "Setting up..." : "Start Trial & Open Dashboard"}
                </button>
              </div>
            </div>
          )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}
