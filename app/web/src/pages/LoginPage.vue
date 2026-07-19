<template>
    <div class="min-h-screen flex items-center justify-center p-4">
        <div class="card w-full max-w-md bg-base-100 shadow-xl">
            <div class="card-body">
                <h2 class="card-title text-2xl font-bold justify-center mb-4">
                    Expense Tracker
                </h2>
                <p class="text-center text-base-content/70 mb-6">Sign in to your account</p>

                <form @submit.prevent="handleLogin" class="space-y-4">
                    <!-- Email Input -->
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text">Email</span>
                        </label>
                        <input v-model="form.email" type="email" placeholder="you@example.com"
                            class="input input-bordered w-full" required />
                    </div>

                    <!-- Password Input -->
                    <div class="form-control">
                        <label class="label">
                            <span class="label-text">Password</span>
                        </label>
                        <input v-model="form.password" type="password" placeholder="••••••••"
                            class="input input-bordered w-full" required />
                    </div>

                    <!-- Error Message -->
                    <div v-if="error" class="alert alert-error text-sm py-2">
                        <span>{{ error }}</span>
                    </div>

                    <!-- Submit Button -->
                    <div class="form-control mt-6">
                        <button type="submit" class="btn btn-primary w-full" :class="{ 'loading': isLoading }"
                            :disabled="isLoading">
                            {{ isLoading ? 'Signing in...' : 'Sign In' }}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import type { LoginRequest } from '../types';

const router = useRouter();
const authStore = useAuthStore();

const form = reactive<LoginRequest>({
    email: '',
    password: ''
});

const isLoading = ref(false);
const error = ref<string | null>(null);

const handleLogin = async () => {
    isLoading.value = true;
    error.value = null;

    try {
        await authStore.login(form);
        router.push({ name: 'Dashboard' });
    } catch (err: any) {
        error.value = err.response?.data?.error || 'Invalid email or password';
    } finally {
        isLoading.value = false;
    }
};
</script>