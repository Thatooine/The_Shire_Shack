import {
  createContext,
  useCallback,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import {
  signInWithEmailAndPassword,
  signInWithPopup,
  signOut,
  onAuthStateChanged,
  GoogleAuthProvider,
  type User as FirebaseUser,
} from "firebase/auth";
import { auth } from "../lib/firebase";
import {
  authenticateWithFirebase,
  registerWithEmailAndPassword,
  registerWithFirebaseToken,
} from "../api/auth";

export interface AuthUser {
  userID: string;
  email: string;
}

interface AuthContextValue {
  user: AuthUser | null;
  loading: boolean;
  loginWithEmail: (email: string, password: string) => Promise<void>;
  loginWithGoogle: () => Promise<void>;
  registerWithEmail: (name: string, email: string, password: string) => Promise<void>;
  registerWithGoogle: (name: string) => Promise<void>;
  logout: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextValue | null>(null);

const googleProvider = new GoogleAuthProvider();

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [loading, setLoading] = useState(true);

  // Exchange a Firebase user's ID token with our backend so it sets the
  // access_token cookie and returns the application-level user identity.
  const exchangeToken = useCallback(async (firebaseUser: FirebaseUser) => {
    const idToken = await firebaseUser.getIdToken();
    const res = await authenticateWithFirebase(idToken);
    setUser({ userID: res.userID, email: res.email });
  }, []);

  // Listen to Firebase auth state so we can re-exchange on page reload if
  // the Firebase session is still alive.
  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (firebaseUser) => {
      if (firebaseUser) {
        try {
          await exchangeToken(firebaseUser);
        } catch {
          setUser(null);
        }
      } else {
        setUser(null);
      }
      setLoading(false);
    });

    return unsubscribe;
  }, [exchangeToken]);

  const loginWithEmail = useCallback(
    async (email: string, password: string) => {
      const credential = await signInWithEmailAndPassword(
        auth,
        email,
        password,
      );
      await exchangeToken(credential.user);
    },
    [exchangeToken],
  );

  const loginWithGoogle = useCallback(async () => {
    const credential = await signInWithPopup(auth, googleProvider);
    const idToken = await credential.user.getIdToken();
    const displayName = credential.user.displayName ?? "";
    const res = await registerWithFirebaseToken(displayName, idToken);
    setUser({ userID: res.userID, email: res.email });
  }, []);

  const registerWithEmail = useCallback(
    async (name: string, email: string, password: string) => {
      const res = await registerWithEmailAndPassword(name, email, password);
      setUser({ userID: res.userID, email: res.email });
    },
    [],
  );

  const registerWithGoogle = useCallback(
    async (name: string) => {
      const credential = await signInWithPopup(auth, googleProvider);
      const idToken = await credential.user.getIdToken();
      const res = await registerWithFirebaseToken(name, idToken);
      setUser({ userID: res.userID, email: res.email });
    },
    [],
  );

  const logout = useCallback(async () => {
    await signOut(auth);
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider
      value={{ user, loading, loginWithEmail, loginWithGoogle, registerWithEmail, registerWithGoogle, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
}
