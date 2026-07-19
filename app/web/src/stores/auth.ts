import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import api from '../services/api';
import type { User, LoginRequest, AuthResponse } from '../types';

export const useAuthStore = defineStore('auth', () => {
    const user = ref<User | null>(null);
    const token = ref<string | null>(localStorage.getItem('token'));

    const isAuthenticated = computed(() => !!token.value);

    const initAuth = () => {
        const storedUser = localStorage.getItem('user');
        if (token.value && storedUser) {
            user.value = JSON.parse(storedUser);
        }
    };

    const login = async (credentials: LoginRequest) => {
        try {
            const { data } = await api.post<AuthResponse>('/auth/login', credentials);

            token.value = data.token;
            user.value = data.user;

            localStorage.setItem('token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));
        } catch (error) {
            throw error;
        }
    };

    const logout = () => {
        token.value = null;
        user.value = null;
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        location.reload();
    };

    return { user, token, isAuthenticated, initAuth, login, logout };
});