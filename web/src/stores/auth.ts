import { create } from 'zustand';
import { login as apiLogin, register as apiRegister, logout as apiLogout, getMe } from '@/lib/api';

interface AuthState {
  isAuthenticated: boolean;
  user: { username: string; user_id: number } | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;
}

export const useAuth = create<AuthState>((set) => ({
  isAuthenticated: false,
  user: null,
  loading: true,

  login: async (username: string, password: string) => {
    const res = await apiLogin(username, password);
    set({ isAuthenticated: true, user: { username: res.username, user_id: res.user_id } });
  },

  register: async (username: string, password: string) => {
    const res = await apiRegister(username, password);
    set({ isAuthenticated: true, user: { username: res.username, user_id: res.user_id } });
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
