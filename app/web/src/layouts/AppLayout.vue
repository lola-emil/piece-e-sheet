<template>
    <div class="drawer lg:drawer-open">
        <input id="my-drawer-2" type="checkbox" class="drawer-toggle" />

        <!-- Page Content -->
        <div class="drawer-content flex flex-col bg-base-200 min-h-screen">
            <!-- Navbar (Mobile only) -->
            <div class="navbar bg-base-100 shadow-sm lg:hidden">
                <div class="flex-none">
                    <label for="my-drawer-2" class="btn btn-square btn-ghost drawer-button">
                        <Menu :size="16" />
                    </label>
                </div>
                <div class="flex-1">
                    <span class="text-xl font-bold px-4">Expense Tracker</span>
                </div>
            </div>

            <!-- Main Content Area -->
            <main class="flex-1 p-4 lg:p-8 overflow-y-auto">
                <router-view />
            </main>
        </div>

        <!-- Sidebar -->
        <div class="drawer-side z-20">
            <label for="my-drawer-2" aria-label="close sidebar" class="drawer-overlay"></label>
            <ul class="menu p-4 w-64 min-h-full bg-base-100 text-base-content border-r border-base-300">
                <!-- Sidebar Header -->
                <li class="mb-6 px-4">
                    <h2 class="text-2xl font-bold text-primary">ExpenseTracker</h2>
                </li>

                <!-- Navigation Links -->
                <li>
                    <router-link to="/" active-class="active" class="flex items-center gap-3">
                        <layout-panel-left :size="18" />
                        Dashboard
                    </router-link>
                </li>
                <li>
                    <router-link to="/expenses" active-class="active" class="flex items-center gap-3">
                        <receipt :size="18" />
                        Expenses
                    </router-link>
                </li>
                <li>
                    <router-link to="/debts" active-class="active" class="flex items-center gap-3">
                        <banknote-arrow-down :size="18" />
                        Debts
                    </router-link>
                </li>
                <li>
                    <router-link to="/categories" active-class="active" class="flex items-center gap-3">
                        <chart-bar-stacked :size="18" />
                        Categories
                    </router-link>
                </li>

                <!-- Logout at bottom -->
                <li class="mt-auto pt-6 border-t border-base-300">
                    <button @click="authStore.logout()" class="flex items-center gap-3 text-error hover:bg-error/10">
                        <log-out :size="18" />
                        Logout
                    </button>
                </li>
            </ul>
        </div>
    </div>
</template>

<script setup lang="ts">
import { BanknoteArrowDown, ChartBarStacked, LayoutPanelLeft, LogOut, Menu, Receipt } from '@lucide/vue';
import { useAuthStore } from '../stores/auth';
const authStore = useAuthStore();
</script>