import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

const flavourTexts = [
  "Lighting the hearth...",
  "Gathering ingredients from the Shire...",
  "Consulting the palantir...",
  "Convincing the cook to wake up...",
  "Polishing the tankards...",
  "Rolling out the Lembas dough...",
  "Stoking the embers...",
];

export default function LoadingPage() {
  const navigate = useNavigate();
  const [textIndex, setTextIndex] = useState(0);
  const [fadeClass, setFadeClass] = useState("visible");

  // Cycle through flavour text
  useEffect(() => {
    const interval = setInterval(() => {
      setFadeClass("hidden");
      setTimeout(() => {
        setTextIndex((i) => (i + 1) % flavourTexts.length);
        setFadeClass("visible");
      }, 400);
    }, 2200);
    return () => clearInterval(interval);
  }, []);

  // Navigate to home after the "loading" finishes
  useEffect(() => {
    const timeout = setTimeout(() => navigate("/home", { replace: true }), 4000);
    return () => clearTimeout(timeout);
  }, [navigate]);

  return (
    <div className="loading-page">
      {/* Floating ember particles */}
      <div className="embers">
        {Array.from({ length: 18 }).map((_, i) => (
          <span key={i} className="ember" />
        ))}
      </div>

      <div className="loading-content">
        <div className="ring-container">
          <div className="ring-outer" />
          <div className="ring-inner" />
          <span className="ring-icon">&#9876;</span>
        </div>

        <h1 className="loading-title">The Orc Shack</h1>

        <p className={`loading-flavour ${fadeClass}`}>
          {flavourTexts[textIndex]}
        </p>

        <div className="loading-bar-track">
          <div className="loading-bar-fill" />
        </div>
      </div>
    </div>
  );
}
