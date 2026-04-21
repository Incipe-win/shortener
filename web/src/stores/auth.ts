import { create } from 'zustand';
import { login as apiLogin, logout as apiLogout, getMe } from '@/lib/api';

interface AuthState {
  isAuthenticated: boolean;
  user: { username: string } | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

export const useAuth = create<AuthState>((set) => ({
  isAuthenticated: false,
  user: null,
  loading: true,

  login: async (username: string, password: string) => {
    await apiLogin(username, password);
    set({ isAuthenticated: true, user: { username } });
  },

  logout: async () => {
    try {
      await apiLogout();
    } catch {
      // ignore
    }
    set({ isAuthenticated: false, user: null });
  },

  checkAuth: async () => {
    try {
      const user = await getMe();
      set({ isAuthenticated: true, user, loading: false });
    } catch {
      set({ isAuthenticated: false, user: null, loading: false });
    }
  },
}));
