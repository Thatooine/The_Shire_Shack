import { get, post, put } from "./client";

export interface Dish {
  id: string;
  restaurant_id: string;
  name: string;
  description: string;
  price: number;
  image: string;
}

interface ListDishesResponse {
  dishes: Dish[];
  total: number;
}

export async function listDishes(): Promise<ListDishesResponse> {
  return get<ListDishesResponse>("/dishes");
}

export async function listDishesByRestaurant(restaurantId: string): Promise<ListDishesResponse> {
  return get<ListDishesResponse>(`/dishes?restaurant_id=${encodeURIComponent(restaurantId)}`);
}

interface CreateDishResponse {
  dish: Dish;
}

export async function createDish(data: {
  name: string;
  description: string;
  price: number;
  image: string;
  restaurant_id: string;
}): Promise<CreateDishResponse> {
  return post<CreateDishResponse>("/dishes", data);
}

interface UpdateDishResponse {
  dish: Dish;
}

export async function updateDish(
  id: string,
  data: { name: string; description: string; price: number; image: string },
): Promise<UpdateDishResponse> {
  return put<UpdateDishResponse>(`/dishes/${id}`, data);
}

interface SubmitRatingRequest {
  dish_id: string;
  score: number;
  review: string;
}

export interface Rating {
  id: string;
  dish_id: string;
  user_id: string;
  score: number;
  review: string;
  created_at: string;
}

interface SubmitRatingResponse {
  rating: Rating;
}

interface ListRatingsResponse {
  ratings: Rating[];
  total: number;
}

export async function listRatings(dishId: string): Promise<ListRatingsResponse> {
  return get<ListRatingsResponse>(`/dishes/${dishId}/ratings`);
}

export async function submitRating(dishId: string, score: number, review: string): Promise<SubmitRatingResponse> {
  return post<SubmitRatingResponse>(`/dishes/${dishId}/ratings`, {
    dish_id: dishId,
    score,
    review,
  } satisfies SubmitRatingRequest);
}
