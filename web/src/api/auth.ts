import { post } from "./client";

export interface AuthResponse {
  userID: string;
  email: string;
}

export interface LoginResponse extends AuthResponse {
  token: string;
}

export function authenticateWithFirebase(
  firebaseToken: string,
): Promise<AuthResponse> {
  return post<AuthResponse>("/auth/firebase", { firebaseToken });
}

export function registerWithEmailAndPassword(
  name: string,
  email: string,
  password: string,
): Promise<AuthResponse> {
  return post<AuthResponse>("/auth/register", { name, email, password });
}

export function registerWithFirebaseToken(
  name: string,
  firebaseToken: string,
): Promise<AuthResponse> {
  return post<AuthResponse>("/auth/register/firebase", {
    name,
    firebaseToken,
  });
}
