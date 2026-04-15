export interface Dish {
  id: string;
  name: string;
  description: string;
  price: number;
  image: string;
  reviews: Review[];
}

export interface Review {
  id: string;
  author: string;
  rating: number;
  comment: string;
  date: string;
}

export const dishes: Dish[] = [
  {
    id: "1",
    name: "Lembas Bread",
    description:
      "Elven waybread wrapped in mallorn leaves. One bite is enough to fill the stomach of a grown man.",
    price: 12.99,
    image: "https://images.unsplash.com/photo-1509440159596-0249088772ff?w=400&h=300&fit=crop",
    reviews: [
      {
        id: "r1",
        author: "Samwise G.",
        rating: 5,
        comment: "Kept me going all the way to Mordor!",
        date: "3019-03-25",
      },
    ],
  },
  {
    id: "2",
    name: "Shire Mushroom Stew",
    description:
      "A hearty stew made with the finest mushrooms from Farmer Maggot's fields. Seasoned with herbs from the Shire.",
    price: 9.49,
    image: "https://images.unsplash.com/photo-1547592166-23ac45744acd?w=400&h=300&fit=crop",
    reviews: [
      {
        id: "r2",
        author: "Frodo B.",
        rating: 4,
        comment: "Reminds me of home. Could use a bit more salt.",
        date: "3019-01-12",
      },
    ],
  },
  {
    id: "3",
    name: "Second Breakfast Platter",
    description:
      "Eggs, bacon, sausages, toast, tomatoes, and nice crispy hash browns. Elevenses not included.",
    price: 14.99,
    image: "https://images.unsplash.com/photo-1533089860892-a7c6f0a88666?w=400&h=300&fit=crop",
    reviews: [
      {
        id: "r3",
        author: "Pippin T.",
        rating: 5,
        comment: "What about second breakfast? This is it!",
        date: "3019-02-01",
      },
    ],
  },
  {
    id: "4",
    name: "Roasted Coney",
    description:
      "Tender rabbit roasted over an open fire. Sam's signature recipe with taters on the side.",
    price: 16.99,
    image: "https://images.unsplash.com/photo-1544025162-d76694265947?w=400&h=300&fit=crop",
    reviews: [],
  },
  {
    id: "5",
    name: "Mirkwood Honey Cake",
    description:
      "Sweet honey cake baked in the tradition of the Woodland Elves. Drizzled with golden forest honey.",
    price: 8.99,
    image: "https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=400&h=300&fit=crop",
    reviews: [
      {
        id: "r5",
        author: "Legolas G.",
        rating: 4,
        comment: "Tastes like the forests of my homeland.",
        date: "3019-03-01",
      },
    ],
  },
  {
    id: "6",
    name: "Dwarvish Ale & Meat Pie",
    description:
      "A robust meat pie filled with beef and root vegetables, paired with a tankard of strong Dwarvish ale.",
    price: 18.49,
    image: "https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?w=400&h=300&fit=crop",
    reviews: [
      {
        id: "r6",
        author: "Gimli S.",
        rating: 5,
        comment: "Now that's a meal fit for a Dwarf lord!",
        date: "3019-03-15",
      },
    ],
  },
];
