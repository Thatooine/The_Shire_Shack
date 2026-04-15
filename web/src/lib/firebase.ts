import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";

const firebaseConfig = {
  apiKey: "AIzaSyBuyyp4nhrFhHWR44mgCfV4hx5TJ8_HQvw",
  authDomain: "bash-interview-project.firebaseapp.com",
  projectId: "bash-interview-project",
  storageBucket: "bash-interview-project.firebasestorage.app",
  messagingSenderId: "890373082925",
  appId: "1:890373082925:web:91b4ac3c5f368e20bc4da4",
  measurementId: "G-Z514KTED5C",
};

const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);
