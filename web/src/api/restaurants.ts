import { get, post } from "./client";

export interface Restaurant {
  id: string;
  ownerID: string;
  name: string;
  image: string;
  city: string;
}

interface RegisterRestaurantResponse {
  restaurant: Restaurant;
}

interface GetMyRestaurantResponse {
  restaurant: Restaurant;
}

export function getMyRestaurant(): Promise<GetMyRestaurantResponse> {
  return get<GetMyRestaurantResponse>("/restaurants/mine");
}

export function registerRestaurant(
  name: string,
  city: string,
  image: string,
): Promise<RegisterRestaurantResponse> {
  return post<RegisterRestaurantResponse>("/restaurants/register", {
    name,
    city,
    image,
  });
}
