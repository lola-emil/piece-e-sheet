<template>
    <div class="space-y-6">
        <h1 class="text-3xl font-bold">Dashboard</h1>

        <!-- 1. Stats Section -->
        <div class="stats stats-vertical lg:stats-horizontal shadow w-full bg-base-100">
            <div class="stat">
                <div class="stat-figure text-primary">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                        class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z">
                        </path>
                    </svg>
                </div>
                <div class="stat-title">Total Expenses (This Month)</div>
                <div class="stat-value text-primary">{{ formatCurrency(stats.thisMonth) }}</div>
                <div class="stat-desc">Since {{ stats.thisMonthStart }}</div>
            </div>

            <div class="stat">
                <div class="stat-figure text-secondary">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                        class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M13 10V3L4 14h7v7l9-11h-7z"></path>
                    </svg>
                </div>
                <div class="stat-title">Last Month</div>
                <div class="stat-value text-secondary">{{ formatCurrency(stats.lastMonth) }}</div>
                <div class="stat-desc">Previous 30 days</div>
            </div>

            <div class="stat">
                <div class="stat-title">Transactions</div>
                <div class="stat-value">{{ expenses.length }}</div>
                <div class="stat-desc">Total records</div>
            </div>
        </div>

        <!-- 2. Charts & Recent Transactions Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

            <!-- Recent Transactions (Takes up 2 columns on desktop) -->
            <div class="card bg-base-100 shadow-xl lg:col-span-2">
                <div class="card-body">
                    <h2 class="card-title">Recent Transactions</h2>
                    <div class="overflow-x-auto">
                        <table class="table table-zebra">
                            <thead>
                                <tr>
                                    <th>Date</th>
                                    <th>Description</th>
                                    <th>Category</th>
                                    <th class="text-right">Amount</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="exp in recentExpenses" :key="exp.id">
                                    <td>{{ formatDate(exp.occurred_at) }}</td>
                                    <td>{{ exp.description }}</td>
                                    <td><span class="badge badge-ghost badge-sm">{{ getCategoryName(exp.category_id)
                                            }}</span></td>
                                    <td class="text-right font-bold text-error">-{{ formatCurrency(exp.amount) }}</td>
                                </tr>
                                <tr v-if="recentExpenses.length === 0">
                                    <td colspan="4" class="text-center text-base-content/50">No expenses found.</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

            <!-- Category Breakdown (Takes up 1 column) -->
            <div class="card bg-base-100 shadow-xl">
                <div class="card-body">
                    <h2 class="card-title">Category Breakdown</h2>
                    <div v-if="Object.keys(categoryTotals).length > 0" class="space-y-4 mt-4">
                        <div v-for="(total, catName) in categoryTotals" :key="catName">
                            <div class="flex justify-between text-sm mb-1">
                                <span>{{ catName }}</span>
                                <span class="font-bold">{{ formatCurrency(total) }}</span>
                            </div>
                            <progress class="progress progress-primary w-full" :value="(total / stats.thisMonth) * 100"
                                max="100"></progress>
                        </div>
                    </div>
                    <div v-else class="flex items-center justify-center h-40 text-base-content/50">
                        No data for this month
                    </div>
                </div>
            </div>

        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import api from '../services/api';
import type { Expense, Category } from '../types';

const expenses = ref<Expense[]>([]);
const categories = ref<Category[]>([]);
const isLoading = ref(true);

// Fetch data
onMounted(async () => {
    try {
        const [expRes, catRes] = await Promise.all([
            api.get('/api/expenses'),
            api.get('/api/categories')
        ]);
        expenses.value = expRes.data.data ?? [];
        categories.value = catRes.data.data ?? [];
    } catch (error) {
        console.error('Failed to load dashboard data', error);
    } finally {
        isLoading.value = false;
    }
});

// Computed Stats
const stats = computed(() => {
    const now = new Date();
    const thisMonth = now.getMonth();
    const thisYear = now.getFullYear();

    let thisMonthTotal = 0;
    let lastMonthTotal = 0;

    expenses.value.forEach(exp => {
        const expDate = new Date(exp.occurred_at);
        if (expDate.getMonth() === thisMonth && expDate.getFullYear() === thisYear) {
            thisMonthTotal += exp.amount;
        }
        // Simple last month check
        if (expDate.getMonth() === thisMonth - 1 && expDate.getFullYear() === thisYear) {
            lastMonthTotal += exp.amount;
        }
    });

    return {
        thisMonth: thisMonthTotal,
        lastMonth: lastMonthTotal,
        thisMonthStart: now.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    };
});

const recentExpenses = computed(() => {
    return [...expenses.value]
        .sort((a, b) => new Date(b.occurred_at).getTime() - new Date(a.occurred_at).getTime())
        .slice(0, 5); // Get top 5
});

const categoryTotals = computed(() => {
    const totals: Record<string, number> = {};
    const now = new Date();

    expenses.value.forEach(exp => {
        const expDate = new Date(exp.occurred_at);
        if (expDate.getMonth() === now.getMonth() && expDate.getFullYear() === now.getFullYear()) {
            const catName = getCategoryName(exp.category_id);
            totals[catName] = (totals[catName] || 0) + exp.amount;
        }
    });
    return totals;
});

// Helpers
const getCategoryName = (id: string | null) => {
    if (!id) return 'Uncategorized';
    const cat = categories.value.find(c => c.id === id);
    return cat ? cat.name : 'Unknown';
};

const formatCurrency = (val: number) => {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);
};

const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString();
};
</script>