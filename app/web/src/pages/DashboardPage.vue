<template>
    <div class="space-y-6">
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <h1 class="text-3xl font-bold">Dashboard</h1>
            <div class="flex gap-2">
                <button class="btn btn-primary btn-sm" @click="$router.push('/expenses/new')">
                    + Add Expense
                </button>
                <button class="btn btn-outline btn-sm" @click="$router.push('/expenses')">
                    View All
                </button>
            </div>
        </div>

        <template v-if="isLoading">
            <div class="stats stats-vertical lg:stats-horizontal shadow w-full bg-base-100">
                <div class="stat">
                    <div class="skeleton h-10 w-20"></div>
                    <div class="skeleton h-4 w-28 mt-2"></div>
                </div>
                <div class="stat">
                    <div class="skeleton h-10 w-20"></div>
                    <div class="skeleton h-4 w-28 mt-2"></div>
                </div>
                <div class="stat">
                    <div class="skeleton h-10 w-20"></div>
                    <div class="skeleton h-4 w-28 mt-2"></div>
                </div>
                <div class="stat">
                    <div class="skeleton h-10 w-20"></div>
                    <div class="skeleton h-4 w-28 mt-2"></div>
                </div>
            </div>
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div class="card bg-base-100 shadow-xl lg:col-span-2">
                    <div class="card-body">
                        <div class="skeleton h-64 w-full"></div>
                    </div>
                </div>
                <div class="card bg-base-100 shadow-xl">
                    <div class="card-body">
                        <div class="skeleton h-64 w-full"></div>
                    </div>
                </div>
            </div>
        </template>

        <template v-else>
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
                    <div class="stat-desc" :class="stats.trend >= 0 ? 'text-error' : 'text-success'">
                        <span>{{ stats.trend >= 0 ? '↑' : '↓' }} {{ Math.abs(stats.trend).toFixed(1) }}%</span> vs last
                        month
                    </div>
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
                    <div class="stat-figure text-accent">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                            class="inline-block w-8 h-8 stroke-current">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z">
                            </path>
                        </svg>
                    </div>
                    <div class="stat-title">Daily Average</div>
                    <div class="stat-value text-accent">{{ formatCurrency(stats.dailyAverage) }}</div>
                    <div class="stat-desc">This month</div>
                </div>

                <div class="stat">
                    <div class="stat-title">Top Category</div>
                    <div class="stat-value text-lg">{{ stats.topCategoryName }}</div>
                    <div class="stat-desc">{{ formatCurrency(stats.topCategoryAmount) }}</div>
                </div>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div class="card bg-base-100 shadow-xl lg:col-span-2">
                    <div class="card-body">
                        <h2 class="card-title">Spending Trend</h2>
                        <div class="h-64 w-full">
                            <Line v-if="lineChartData" :data="lineChartData" :options="lineChartOptions"
                                class="w-full h-full!" />
                            <div v-else class="flex items-center justify-center h-full text-base-content/50">
                                Not enough data for trend
                            </div>
                        </div>
                    </div>
                </div>

                <div class="card bg-base-100 shadow-xl">
                    <div class="card-body">
                        <h2 class="card-title">Category Breakdown</h2>
                        <div class="h-64 w-full flex items-center justify-center">
                            <Doughnut v-if="doughnutChartData" :data="doughnutChartData" :options="doughnutOptions"
                                class="max-h-full!" />
                            <div v-else class="text-base-content/50 text-center">
                                No data for this month
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="card bg-base-100 shadow-xl">
                <div class="card-body">
                    <div class="flex justify-between items-center mb-4">
                        <h2 class="card-title">Recent Transactions</h2>
                        <button class="btn btn-ghost btn-sm" @click="$router.push('/expenses')">View All →</button>
                    </div>
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
                                    <td colspan="4" class="text-center text-base-content/50 py-8">No expenses found.
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { Line, Doughnut } from 'vue-chartjs';
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    LineElement,
    CategoryScale,
    LinearScale,
    PointElement,
    ArcElement,
    Filler
} from 'chart.js';
import api from '../services/api';
import type { Expense, Category } from '../types';
import { formatCurrency, formatDate } from '@/utils/helpers';

ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement, ArcElement, Filler);

const expenses = ref<Expense[]>([]);
const categories = ref<Category[]>([]);
const isLoading = ref(true);

const getThemeColor = (varName: string, opacity = 1): string => {
    const raw = getComputedStyle(document.documentElement)
        .getPropertyValue(varName)
        .trim();
    if (!raw) return '#888888';
    if (opacity === 1) return raw;
    return `color-mix(in oklch, ${raw} ${Math.round(opacity * 100)}%, transparent)`;
};

const themeColors = () => ({
    primary: getThemeColor('--color-primary'),
    primaryFaded: getThemeColor('--color-primary', 0.15),
    secondary: getThemeColor('--color-secondary'),
    secondaryFaded: getThemeColor('--color-secondary', 0.15),
    accent: getThemeColor('--color-accent'),
    info: getThemeColor('--color-info'),
    success: getThemeColor('--color-success'),
    warning: getThemeColor('--color-warning'),
    error: getThemeColor('--color-error'),
    baseContent: getThemeColor('--color-base-content'),
    baseContentMuted: getThemeColor('--color-base-content', 0.5),
    baseContentFaint: getThemeColor('--color-base-content', 0.1),
});

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

const stats = computed(() => {
    const now = new Date();
    const thisMonth = now.getMonth();
    const thisYear = now.getFullYear();
    const daysPassed = now.getDate();

    let thisMonthTotal = 0;
    let lastMonthTotal = 0;
    const categorySums: Record<string, number> = {};

    expenses.value.forEach(exp => {
        const expDate = new Date(exp.occurred_at);
        const expMonth = expDate.getMonth();
        const expYear = expDate.getFullYear();

        if (expMonth === thisMonth && expYear === thisYear) {
            thisMonthTotal += exp.amount;
            const catName = getCategoryName(exp.category_id);
            categorySums[catName] = (categorySums[catName] || 0) + exp.amount;
        }

        const lastMonthDate = new Date(thisYear, thisMonth - 1, 1);
        if (expMonth === lastMonthDate.getMonth() && expYear === lastMonthDate.getFullYear()) {
            lastMonthTotal += exp.amount;
        }
    });

    const trend = lastMonthTotal > 0
        ? ((thisMonthTotal - lastMonthTotal) / lastMonthTotal) * 100
        : (thisMonthTotal > 0 ? 100 : 0);

    let topCategoryName = 'None';
    let topCategoryAmount = 0;
    for (const [name, amount] of Object.entries(categorySums)) {
        if (amount > topCategoryAmount) {
            topCategoryName = name;
            topCategoryAmount = amount;
        }
    }

    return {
        thisMonth: thisMonthTotal,
        lastMonth: lastMonthTotal,
        dailyAverage: daysPassed > 0 ? thisMonthTotal / daysPassed : 0,
        trend,
        topCategoryName,
        topCategoryAmount
    };
});

const recentExpenses = computed(() => {
    return [...expenses.value]
        .sort((a, b) => new Date(b.occurred_at).getTime() - new Date(a.occurred_at).getTime())
        .slice(0, 5);
});

const lineChartData = computed(() => {
    const t = themeColors();
    const now = new Date();
    const thisMonth = now.getMonth();
    const thisYear = now.getFullYear();
    const daysInMonth = new Date(thisYear, thisMonth + 1, 0).getDate();

    const thisMonthData = new Array(daysInMonth).fill(0);

    const lastMonthDate = new Date(thisYear, thisMonth - 1, 1);
    const lastMonth = lastMonthDate.getMonth();
    const lastMonthYear = lastMonthDate.getFullYear();
    const daysInLastMonth = new Date(lastMonthYear, lastMonth + 1, 0).getDate();
    const lastMonthData = new Array(daysInLastMonth).fill(0);

    expenses.value.forEach(exp => {
        const d = new Date(exp.occurred_at);
        const day = d.getDate() - 1;
        if (d.getMonth() === thisMonth && d.getFullYear() === thisYear && day < daysInMonth) {
            thisMonthData[day] += exp.amount;
        } else if (d.getMonth() === lastMonth && d.getFullYear() === lastMonthYear && day < daysInLastMonth) {
            lastMonthData[day] += exp.amount;
        }
    });

    const currentDay = now.getDate();
    const labels = Array.from({ length: currentDay }, (_, i) => i + 1);
    const thisMonthSliced = thisMonthData.slice(0, currentDay);
    const lastMonthSliced = lastMonthData.slice(0, Math.min(currentDay, daysInLastMonth));
    while (lastMonthSliced.length < currentDay) lastMonthSliced.push(0);

    return {
        labels,
        datasets: [
            {
                label: 'This Month',
                data: thisMonthSliced,
                borderColor: t.primary,
                backgroundColor: t.primaryFaded,
                fill: true,
                tension: 0.4,
                pointRadius: 0,
                pointHoverRadius: 5,
                pointHoverBackgroundColor: t.primary,
            },
            {
                label: 'Last Month',
                data: lastMonthSliced,
                borderColor: t.secondary,
                backgroundColor: 'transparent',
                borderDash: [5, 5],
                fill: false,
                tension: 0.4,
                pointRadius: 0,
                pointHoverRadius: 5,
                pointHoverBackgroundColor: t.secondary,
            }
        ]
    };
});

const lineChartOptions = computed(() => {
    const t = themeColors();
    return {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: true,
                position: 'top' as const,
                labels: { color: t.baseContent, boxWidth: 12, padding: 15 }
            },
            tooltip: {
                backgroundColor: t.baseContent,
                titleColor: getThemeColor('--color-base-100'),
                bodyColor: getThemeColor('--color-base-100'),
                borderColor: t.baseContentFaint,
                borderWidth: 1,
                callbacks: {
                    label: (ctx: any) => `${ctx.dataset.label}: ${formatCurrency(ctx.raw)}`
                }
            }
        },
        scales: {
            x: {
                grid: { color: t.baseContentFaint },
                ticks: { color: t.baseContentMuted }
            },
            y: {
                beginAtZero: true,
                grid: { color: t.baseContentFaint },
                ticks: {
                    color: t.baseContentMuted,
                    callback: (val: any) => formatCurrency(val)
                }
            }
        }
    };
});

const doughnutChartData = computed(() => {
    const t = themeColors();
    const totals: Record<string, number> = {};
    const now = new Date();

    expenses.value.forEach(exp => {
        const expDate = new Date(exp.occurred_at);
        if (expDate.getMonth() === now.getMonth() && expDate.getFullYear() === now.getFullYear()) {
            const catName = getCategoryName(exp.category_id);
            totals[catName] = (totals[catName] || 0) + exp.amount;
        }
    });

    if (Object.keys(totals).length === 0) return null;

    const palette = [
        t.primary, t.secondary, t.accent,
        t.info, t.success, t.warning, t.error
    ];

    return {
        labels: Object.keys(totals),
        datasets: [{
            data: Object.values(totals),
            backgroundColor: Object.keys(totals).map((_, i) => palette[i % palette.length]),
            borderColor: getThemeColor('--color-base-100'),
            borderWidth: 2,
            hoverOffset: 8
        }]
    };
});

const doughnutOptions = computed(() => {
    const t = themeColors();
    return {
        responsive: true,
        maintainAspectRatio: false,
        cutout: '65%',
        plugins: {
            legend: {
                display: true,
                position: 'bottom' as const,
                labels: {
                    color: t.baseContent,
                    boxWidth: 12,
                    padding: 15
                }
            },
            tooltip: {
                backgroundColor: t.baseContent,
                titleColor: getThemeColor('--color-base-100'),
                bodyColor: getThemeColor('--color-base-100'),
                callbacks: {
                    label: (ctx: any) => `${ctx.label}: ${formatCurrency(ctx.raw)}`
                }
            }
        }
    };
});

const getCategoryName = (id: string | null) => {
    if (!id) return 'Uncategorized';
    const cat = categories.value.find(c => c.id === id);
    return cat ? cat.name : 'Unknown';
};
</script>

<style scoped>
/* Ensure chart.js canvas respects container height in Vue */
.h-full\! {
    height: 100% !important;
}

.max-h-full\! {
    max-height: 100% !important;
}
</style>