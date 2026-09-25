<template>
    <div class="space-y-8">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
            <div>
                <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">
                    Dashboard
                </h1>

                <p class="mt-1 text-sm text-base-content/60">
                    Your spending overview for {{ currentMonthLabel }}
                </p>
            </div>

            <div class="flex flex-wrap items-center gap-2">
                <select v-model="selectedAccount" class="select select-sm w-auto min-w-36 border-base-300 bg-base-100">
                    <option value="">
                        All Accounts
                    </option>

                    <option value="-1">
                        Not Specified
                    </option>

                    <option v-for="value in accounts" :key="value.id" :value="value.id">
                        {{ value.name }}
                    </option>
                </select>

                <button class="btn btn-primary btn-sm" @click="expenseModal?.openModal">
                    <Plus :size="16" />
                    Add Expense
                </button>
            </div>
        </div>

        <template v-if="isLoading">
            <!-- Stats skeleton -->
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
                <div v-for="index in 4" :key="index" class="rounded-xl border border-base-300/60 bg-base-100 p-5">
                    <div class="flex items-start justify-between">
                        <div class="space-y-3">
                            <div class="skeleton h-4 w-24"></div>
                            <div class="skeleton h-8 w-32"></div>
                            <div class="skeleton h-3 w-24"></div>
                        </div>

                        <div class="skeleton h-10 w-10 rounded-lg"></div>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-1 gap-5">
                <div class="rounded-xl border border-base-300/60 bg-base-100 p-5 ">
                    <div class="mb-5 space-y-2">
                        <div class="skeleton h-5 w-32"></div>
                        <div class="skeleton h-3 w-56"></div>
                    </div>

                    <div class="skeleton h-80 w-full"></div>
                </div>

            </div>

            <div class="grid grid-cols-1 gap-5 xl:grid-cols-3">

                <div class="rounded-xl border border-base-300/60 bg-base-100 p-5">
                    <div class="mb-5 space-y-2">
                        <div class="skeleton h-5 w-40"></div>
                        <div class="skeleton h-3 w-48"></div>
                    </div>

                    <div class="skeleton mx-auto h-52 w-52 rounded-full"></div>
                </div>


                <div class="rounded-xl border border-base-300/60 bg-base-100 p-5 xl:col-span-2">
                    <div class="mb-5 flex items-center justify-between">
                        <div class="space-y-2">
                            <div class="skeleton h-5 w-40"></div>
                            <div class="skeleton h-3 w-52"></div>
                        </div>

                        <div class="skeleton h-8 w-20"></div>
                    </div>

                    <div class="space-y-1">
                        <div v-for="index in 5" :key="index"
                            class="flex items-center gap-4 border-b border-base-300/40 py-4 last:border-0">
                            <div class="skeleton h-10 w-10 shrink-0 rounded-lg"></div>

                            <div class="min-w-0 flex-1 space-y-2">
                                <div class="skeleton h-4 w-28"></div>
                                <div class="skeleton h-3 w-20"></div>
                            </div>

                            <div class="skeleton h-4 w-20"></div>
                        </div>
                    </div>
                </div>
            </div>

        </template>

        <template v-else>
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
                <!-- Total Spending -->
                <div
                    class="rounded-xl border border-base-300/60 bg-base-100 p-5 transition-colors hover:border-base-300">
                    <div class="flex items-start justify-between gap-4">
                        <div class="min-w-0">
                            <p class="text-sm font-medium text-base-content/60">
                                Total spending
                            </p>

                            <p class="mt-2 truncate text-2xl font-bold tracking-tight text-primary">
                                {{ formatCurrency(stats.thisMonth) }}
                            </p>

                            <p class="mt-1 text-xs font-medium" :class="stats.trend >= 0
                                ? 'text-error'
                                : 'text-success'
                                ">
                                {{ stats.trend >= 0 ? '↑' : '↓' }}
                                {{ Math.abs(stats.trend).toFixed(1) }}%

                                <span class="font-normal text-base-content/50">
                                    vs last month
                                </span>
                            </p>
                        </div>

                        <div
                            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
                            <Receipt :size="19" />
                        </div>
                    </div>
                </div>

                <!-- Daily Average -->
                <div
                    class="rounded-xl border border-base-300/60 bg-base-100 p-5 transition-colors hover:border-base-300">
                    <div class="flex items-start justify-between gap-4">
                        <div class="min-w-0">
                            <p class="text-sm font-medium text-base-content/60">
                                Daily average
                            </p>

                            <p class="mt-2 truncate text-2xl font-bold tracking-tight">
                                {{ formatCurrency(stats.dailyAverage) }}
                            </p>

                            <p class="mt-1 text-xs text-base-content/50">
                                This month
                            </p>
                        </div>

                        <div
                            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-base-200 text-base-content/60">
                            <Calendar :size="19" />
                        </div>
                    </div>
                </div>

                <!-- Last Month -->
                <div
                    class="rounded-xl border border-base-300/60 bg-base-100 p-5 transition-colors hover:border-base-300">
                    <div class="flex items-start justify-between gap-4">
                        <div class="min-w-0">
                            <p class="text-sm font-medium text-base-content/60">
                                Last month
                            </p>

                            <p class="mt-2 truncate text-2xl font-bold tracking-tight">
                                {{ formatCurrency(stats.lastMonth) }}
                            </p>

                            <p class="mt-1 text-xs text-base-content/50">
                                Previous month
                            </p>
                        </div>

                        <div
                            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-base-200 text-base-content/60">
                            <Zap :size="19" />
                        </div>
                    </div>
                </div>

                <!-- Top Category -->
                <div
                    class="rounded-xl border border-base-300/60 bg-base-100 p-5 transition-colors hover:border-base-300">
                    <div class="flex items-start justify-between gap-4">
                        <div class="min-w-0">
                            <p class="text-sm font-medium text-base-content/60">
                                Top category
                            </p>

                            <p class="mt-2 truncate text-xl font-bold tracking-tight" :title="stats.topCategoryName">
                                {{ stats.topCategoryName }}
                            </p>

                            <p class="mt-1 text-xs text-base-content/50">
                                {{ formatCurrency(stats.topCategoryAmount) }}
                            </p>
                        </div>

                        <div
                            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-warning/10 text-warning">
                            <Tag :size="19" />
                        </div>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-1 gap-5 xl:grid-cols-1">
                <!-- Spending Trend -->
                <section class="rounded-xl border border-base-300/60 bg-base-100 p-5">
                    <div class="mb-5 flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                        <div>
                            <h2 class="text-base font-semibold">
                                Spending Trend
                            </h2>

                            <p class="mt-1 text-xs text-base-content/50">
                                Daily spending compared with last month
                            </p>
                        </div>

                        <div class="flex items-center gap-4 text-xs text-base-content/60">
                            <div class="flex items-center gap-2">
                                <span class="h-2 w-2 rounded-full bg-primary"></span>
                                This month
                            </div>

                            <div class="flex items-center gap-2">
                                <span class="h-2 w-2 rounded-full bg-secondary/60"></span>
                                Last month
                            </div>
                        </div>
                    </div>

                    <div class="h-80 w-full">
                        <Line v-if="lineChartData" :data="lineChartData" :options="lineChartOptions"
                            class="h-full! w-full!" />

                        <div v-else class="flex h-full items-center justify-center text-sm text-base-content/40">
                            Not enough data for trend
                        </div>
                    </div>
                </section>

            </div>

            <div class="grid grid-cols-1 gap-5 xl:grid-cols-3">

                <!-- Category Breakdown -->
                <section class="rounded-xl border border-base-300/60 bg-base-100 p-5">
                    <div class="mb-5">
                        <h2 class="text-base font-semibold">
                            Category Breakdown
                        </h2>

                        <p class="mt-1 text-xs text-base-content/50">
                            Where your money is going this month
                        </p>
                    </div>

                    <div v-if="doughnutChartData" class="flex flex-col gap-5">
                        <!-- Donut -->
                        <div class="relative mx-auto h-52 w-52">
                            <Doughnut :data="doughnutChartData" :options="doughnutOptions" class="h-full! w-full!" />

                            <div class="pointer-events-none absolute inset-0 flex flex-col items-center justify-center">
                                <span class="text-xs text-base-content/50">
                                    This month
                                </span>

                                <span class="mt-1 text-xl font-bold tracking-tight">
                                    {{ formatCurrency(stats.thisMonth) }}
                                </span>
                            </div>
                        </div>

                        <!-- Category List -->
                        <div class="space-y-3">
                            <div v-for="category in categoryBreakdown" :key="category.name"
                                class="flex items-center gap-3">
                                <span class="h-2.5 w-2.5 shrink-0 rounded-full" :style="{
                                    backgroundColor: category.color,
                                }"></span>

                                <div class="min-w-0 flex-1">
                                    <div class="flex items-center justify-between gap-3">
                                        <span class="truncate text-sm text-base-content/80" :title="category.name">
                                            {{ category.name }}
                                        </span>

                                        <span class="shrink-0 text-sm font-medium tabular-nums">
                                            {{
                                                formatCurrency(
                                                    category.amount
                                                )
                                            }}
                                        </span>
                                    </div>

                                    <div class="mt-1 flex items-center justify-between gap-3">
                                        <div class="h-1 flex-1 overflow-hidden rounded-full bg-base-300/50">
                                            <div class="h-full rounded-full transition-all" :style="{
                                                width: `${category.percentage}%`,
                                                backgroundColor:
                                                    category.color,
                                            }"></div>
                                        </div>

                                        <span class="w-10 text-right text-[11px] text-base-content/40">
                                            {{ category.percentage.toFixed(0) }}%
                                        </span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div v-else class="flex min-h-80 items-center justify-center text-sm text-base-content/40">
                        No data for this month
                    </div>
                </section>

                <section class="rounded-xl border border-base-300/60 bg-base-100 p-5 xl:col-span-2">
                    <div class="mb-4 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                        <div>
                            <h2 class="text-base font-semibold">
                                Recent Transactions
                            </h2>

                            <p class="mt-1 text-xs text-base-content/50">
                                Your latest expenses
                            </p>
                        </div>

                        <button
                            class="btn btn-ghost btn-sm self-start text-base-content/60 hover:text-base-content sm:self-auto"
                            @click="$router.push('/expenses')">
                            View all
                            <span aria-hidden="true">→</span>
                        </button>
                    </div>

                    <div class="overflow-x-auto">
                        <table class="table">
                            <thead>
                                <tr class="border-base-300/50 text-xs uppercase tracking-wider text-base-content/40">
                                    <th class="font-medium">
                                        Date
                                    </th>

                                    <th class="font-medium">
                                        Description
                                    </th>

                                    <th class="font-medium">
                                        Category
                                    </th>

                                    <th class="text-right font-medium">
                                        Amount
                                    </th>
                                </tr>
                            </thead>

                            <tbody>
                                <tr v-for="exp in recentExpenses" :key="exp.id"
                                    class="border-base-300/40 transition-colors hover:bg-base-200/30">
                                    <td class="whitespace-nowrap text-sm text-base-content/50">
                                        {{ formatDate(exp.occurred_at) }}
                                    </td>

                                    <td>
                                        <span class="font-medium">
                                            {{ exp.description }}
                                        </span>
                                    </td>

                                    <td>
                                        <span class="badge badge-sm border-0 bg-base-200 text-base-content/70">
                                            {{ getCategoryName(exp.category_id) }}
                                        </span>
                                    </td>

                                    <td class="text-right font-semibold tabular-nums text-error">
                                        -{{ formatCurrency(exp.amount) }}
                                    </td>
                                </tr>

                                <tr v-if="recentExpenses.length === 0">
                                    <td colspan="4" class="py-12 text-center">
                                        <div class="flex flex-col items-center justify-center">
                                            <div
                                                class="mb-3 flex h-10 w-10 items-center justify-center rounded-full bg-base-200 text-base-content/40">
                                                <Receipt :size="18" />
                                            </div>

                                            <p class="text-sm font-medium">
                                                No expenses found
                                            </p>

                                            <p class="mt-1 text-xs text-base-content/50">
                                                Your recent transactions will appear
                                                here.
                                            </p>
                                        </div>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </section>
            </div>
        </template>

        <ExpenseFormModal ref="expense-modal" :expense="null" :categories="categories" :accounts="accounts"
            :is-saving="isSaving" @save="handleSave" />
    </div>
</template>

<script setup lang="ts">
import {
    ref,
    computed,
    onMounted,
    useTemplateRef,
    watch,
} from 'vue';

import {
    Line,
    Doughnut,
} from 'vue-chartjs';

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
    Filler,
} from 'chart.js';

import type { CreateExpenseRequest } from '../types';

import {
    formatCurrency,
    formatDate,
} from '@/utils/helpers';

import {
    Calendar,
    Plus,
    Receipt,
    Tag,
    Zap,
} from '@lucide/vue';

import ExpenseFormModal from '@/components/ExpenseFormModal.vue';
import { useExpenses } from '@/composables/useExpenses';

ChartJS.register(
    Title,
    Tooltip,
    Legend,
    LineElement,
    CategoryScale,
    LinearScale,
    PointElement,
    ArcElement,
    Filler,
);

const isLoading = ref(true);

const {
    expenses,
    categories,
    accounts,
    isSaving,
    filters,
    saveExpense,
    fetchExpenses,
    fetchCategories,
    fetchAccounts,
} = useExpenses();

const expenseModal = useTemplateRef('expense-modal');

const selectedAccount = ref('');

const getThemeColor = (
    varName: string,
    opacity = 1,
): string => {
    const raw = getComputedStyle(document.documentElement)
        .getPropertyValue(varName)
        .trim();

    if (!raw) {
        return '#888888';
    }

    if (opacity === 1) {
        return raw;
    }

    return `color-mix(in oklch, ${raw} ${Math.round(
        opacity * 100,
    )}%, transparent)`;
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
    baseContentMuted: getThemeColor(
        '--color-base-content',
        0.5,
    ),
    baseContentFaint: getThemeColor(
        '--color-base-content',
        0.06,
    ),
});

const currentMonthLabel = computed(() => {
    return new Intl.DateTimeFormat('en-US', {
        month: 'long',
        year: 'numeric',
    }).format(new Date());
});

watch(selectedAccount, () => {
    fetchExpenses();
});

onMounted(async () => {
    try {
        filters.value = {
            no_pagination: true,
        };

        await Promise.all([
            fetchExpenses(),
            fetchCategories(),
            fetchAccounts(),
        ]);
    } catch (error) {
        console.error(
            'Failed to load dashboard data',
            error,
        );
    } finally {
        isLoading.value = false;
    }
});

const getCategoryName = (
    id: string | null,
) => {
    if (!id) {
        return 'Uncategorized';
    }

    const cat = categories.value.find(
        c => c.id === id,
    );

    return cat
        ? cat.name
        : 'Unknown';
};


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

        if (
            expMonth === thisMonth &&
            expYear === thisYear
        ) {
            thisMonthTotal += exp.amount;

            const catName = getCategoryName(
                exp.category_id,
            );

            categorySums[catName] =
                (categorySums[catName] || 0) +
                exp.amount;
        }

        const lastMonthDate = new Date(
            thisYear,
            thisMonth - 1,
            1,
        );

        if (
            expMonth === lastMonthDate.getMonth() &&
            expYear === lastMonthDate.getFullYear()
        ) {
            lastMonthTotal += exp.amount;
        }
    });

    const trend =
        lastMonthTotal > 0
            ? ((thisMonthTotal - lastMonthTotal) /
                lastMonthTotal) *
            100
            : thisMonthTotal > 0
                ? 100
                : 0;

    let topCategoryName = 'None';
    let topCategoryAmount = 0;

    for (const [
        name,
        amount,
    ] of Object.entries(categorySums)) {
        if (amount > topCategoryAmount) {
            topCategoryName = name;
            topCategoryAmount = amount;
        }
    }

    return {
        thisMonth: thisMonthTotal,
        lastMonth: lastMonthTotal,

        dailyAverage:
            daysPassed > 0
                ? thisMonthTotal / daysPassed
                : 0,

        trend,

        topCategoryName,
        topCategoryAmount,
    };
});

const recentExpenses = computed(() => {
    return [...expenses.value]
        .sort(
            (a, b) =>
                new Date(
                    b.occurred_at,
                ).getTime() -
                new Date(
                    a.occurred_at,
                ).getTime(),
        )
        .slice(0, 10);
});

const categoryPalette = computed(() => {
    const t = themeColors();

    return [
        t.primary,
        t.secondary,
        t.accent,
        t.info,
        t.success,
        t.warning,
        t.error,
    ];
});

const categoryBreakdown = computed(() => {
    const totals: Record<string, number> = {};

    const now = new Date();

    expenses.value.forEach(exp => {
        const expDate = new Date(
            exp.occurred_at,
        );

        if (
            expDate.getMonth() !== now.getMonth() ||
            expDate.getFullYear() !==
            now.getFullYear()
        ) {
            return;
        }

        const catName = getCategoryName(
            exp.category_id,
        );

        totals[catName] =
            (totals[catName] || 0) +
            exp.amount;
    });

    const total = Object.values(totals).reduce(
        (sum, value) => sum + value,
        0,
    );

    const palette = categoryPalette.value;

    return Object.entries(totals)
        .sort(
            ([, a], [, b]) => b - a,
        )
        .map(([name, amount], index) => ({
            name,
            amount,
            percentage:
                total > 0
                    ? (amount / total) *
                    100
                    : 0,
            color:
                palette[
                index %
                palette.length
                ],
        }),
        );
});

const lineChartData = computed(() => {
    const t = themeColors();

    const now = new Date();

    const thisMonth = now.getMonth();
    const thisYear = now.getFullYear();

    const daysInMonth = new Date(
        thisYear,
        thisMonth + 1,
        0,
    ).getDate();

    const thisMonthData = new Array(
        daysInMonth,
    ).fill(0);

    const lastMonthDate = new Date(
        thisYear,
        thisMonth - 1,
        1,
    );

    const lastMonth = lastMonthDate.getMonth();

    const lastMonthYear = lastMonthDate.getFullYear();

    const daysInLastMonth =
        new Date(
            lastMonthYear,
            lastMonth + 1,
            0,
        ).getDate();

    const lastMonthData = new Array(
        daysInLastMonth,
    ).fill(0);

    expenses.value.forEach(exp => {
        const d = new Date(
            exp.occurred_at,
        );

        const day = d.getDate() - 1;

        if (
            d.getMonth() === thisMonth &&
            d.getFullYear() === thisYear &&
            day < daysInMonth
        ) {
            thisMonthData[day] +=
                exp.amount;
        }

        else if (
            d.getMonth() === lastMonth &&
            d.getFullYear() ===
            lastMonthYear &&
            day < daysInLastMonth
        ) {
            lastMonthData[day] +=
                exp.amount;
        }
    });

    const currentDay = now.getDate();

    const labels = Array.from(
        {
            length: currentDay,
        },
        (_, i) => i + 1,
    );

    const thisMonthSliced = thisMonthData.slice(
        0,
        currentDay,
    );

    const lastMonthSliced = lastMonthData.slice(
        0,
        Math.min(
            currentDay,
            daysInLastMonth,
        ),
    );

    while (
        lastMonthSliced.length <
        currentDay
    ) {
        lastMonthSliced.push(0);
    }

    return {
        labels,

        datasets: [
            {
                label: 'This Month',
                data: thisMonthSliced,
                borderColor: t.primary,
                backgroundColor:
                    t.primaryFaded,
                fill: true,
                tension: 0.4,
                pointRadius: 0,
                pointHoverRadius: 5,
                pointHoverBackgroundColor:
                    t.primary,
                pointHoverBorderColor:
                    getThemeColor(
                        '--color-base-100',
                    ),

                pointHoverBorderWidth: 2,
                borderWidth: 2.5,
            },
            {
                label: 'Last Month',
                data: lastMonthSliced,
                borderColor:
                    t.secondary,
                backgroundColor:
                    'transparent',
                borderDash: [5, 5],
                fill: false,
                tension: 0.4,
                pointRadius: 0,
                pointHoverRadius: 5,
                pointHoverBackgroundColor:
                    t.secondary,
                pointHoverBorderColor:
                    getThemeColor(
                        '--color-base-100',
                    ),
                pointHoverBorderWidth: 2,
                borderWidth: 2,
            },
        ],
    };
});

const lineChartOptions = computed(() => {
    const t = themeColors();

    return {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
            intersect: false,
            mode: 'index' as const,
        },
        plugins: {
            legend: {
                display: false,
            },
            tooltip: {
                backgroundColor:
                    t.baseContent,
                titleColor:
                    getThemeColor(
                        '--color-base-100',
                    ),
                bodyColor:
                    getThemeColor(
                        '--color-base-100',
                    ),
                borderColor:
                    t.baseContentFaint,
                borderWidth: 1,
                padding: 10,
                displayColors: true,
                callbacks: {
                    label: (ctx: any) =>
                        `${ctx.dataset.label}: ${formatCurrency(
                            ctx.raw,
                        )}`,
                },
            },
        },

        scales: {
            x: {
                border: {
                    display: false,
                },
                grid: {
                    display: false,
                },
                ticks: {
                    color: t.baseContentMuted,
                    maxTicksLimit: 8,
                    padding: 8,
                },
            },
            y: {
                beginAtZero: true,
                border: {
                    display: false,
                },
                grid: {
                    color: t.baseContentFaint,
                    drawTicks: false,
                },
                ticks: {
                    color: t.baseContentMuted,
                    padding: 8,
                    maxTicksLimit: 6,
                    callback: (val: any) =>
                        formatCurrency(val),
                },
            },
        },
    };
});

const doughnutChartData = computed(() => {
    const breakdown = categoryBreakdown.value;

    if (!breakdown.length) {
        return null;
    }

    return {
        labels: breakdown.map(category => category.name),

        datasets: [
            {
                data: breakdown.map(
                    category => category.amount,
                ),
                backgroundColor: breakdown.map(category => category.color,),
                borderColor: getThemeColor('--color-base-100',),
                borderWidth: 3,
                hoverOffset: 6,
            },
        ],
    };
});

const doughnutOptions = computed(() => {
    const t = themeColors();

    return {
        responsive: true,
        maintainAspectRatio: false,
        cutout: '68%',
        plugins: {
            legend: {
                display: false,
            },
            tooltip: {
                backgroundColor:
                    t.baseContent,
                titleColor: getThemeColor('--color-base-100',),
                bodyColor: getThemeColor('--color-base-100',),
                borderColor: t.baseContentFaint,
                borderWidth: 1,
                padding: 10,
                callbacks: {
                    label: (ctx: any) =>
                        `${ctx.label}: ${formatCurrency(
                            ctx.raw,
                        )}`,
                },
            },
        },
    };
});

const handleSave = async (
    payload: CreateExpenseRequest,
) => {
    const success = await saveExpense(payload);
    if (success) {
        expenseModal.value?.closeModal();
    }
};
</script>

<style scoped>
.h-full\! {
    height: 100% !important;
}

.w-full\! {
    width: 100% !important;
}
</style>