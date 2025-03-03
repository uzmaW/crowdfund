import { defineStore } from 'pinia';

interface User {
  id: number;
  username: string;
  email: string;
  // Add other user properties as needed
}

interface UserState {
  token: string | null;
  user: User | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: null,
    user: null,
  }),
  actions: {
    setToken(token: string): void {
      this.token = token;
    },
    setUser(user: User): void {
      this.user = user;
    },
    clearUser(): void {
      this.token = null;
      this.user = null;
    },
  },
  persist: true,
});