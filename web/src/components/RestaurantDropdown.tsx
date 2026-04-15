import { useState, useRef, useEffect } from "react";

interface Props {
  onRegister: () => void;
  onManage: () => void;
}

export default function RestaurantDropdown({ onRegister, onManage }: Props) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="restaurant-dropdown" ref={ref}>
      <button
        className="restaurant-dropdown-btn"
        onClick={() => setOpen((prev) => !prev)}
      >
        Restaurants
        <span className={`dropdown-arrow ${open ? "open" : ""}`}>&#9662;</span>
      </button>
      {open && (
        <div className="restaurant-dropdown-menu">
          <button
            className="dropdown-item"
            onClick={() => {
              setOpen(false);
              onRegister();
            }}
          >
            <span className="dropdown-item-icon">&#x2726;</span>
            Register Restaurant
          </button>
          <button
            className="dropdown-item"
            onClick={() => {
              setOpen(false);
              onManage();
            }}
          >
            <span className="dropdown-item-icon">&#x2699;</span>
            My Restaurant
          </button>
        </div>
      )}
    </div>
  );
}
