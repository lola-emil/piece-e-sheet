<template>
    <div class="space-y-8">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
            <div>
                <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">
                    Expenses
                </h1>

                <p class="mt-1 text-sm text-base-content/60">
                    Manage and review your spending history.
                </p>
            </div>

            <button class="btn btn-primary btn-sm sm:btn-md" @click="openAddModal">
                <Plus :size="17" />
                Add Expense
            </button>
        </div>

        <section class="rounded-xl border border-base-300/60 bg-base-100">
            <!-- Primary filters -->
            <div class="p-4">
                <div class="flex flex-col gap-3 lg:flex-row lg:items-center">
                    <!-- Search -->
                    <label class="input input-md flex w-full items-center gap-2 bg-base-100 lg:flex-1">
                        <Search :size="16" class="shrink-0 text-base-content/40" />

                        <input v-model="filters.search" type="text" placeholder="Search expenses..." class="grow" />

                        <button v-if="filters.search" type="button" class="btn btn-ghost btn-xs btn-circle"
                            aria-label="Clear search" @click="filters.search = ''">
                            <X :size="14" />
                        </button>
                    </label>

                    <!-- Category -->
                    <select v-model="filters.category_id" class="select select-md w-full bg-base-100 lg:w-52">
                        <option :value="undefined">
                            All Categories
                        </option>

                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                            {{ cat.name }}
                        </option>
                    </select>

                    <!-- Filters -->
                    <button type="button" class="btn btn-md lg:w-auto" :class="showAdvancedFilters
                        ? 'btn-active'
                        : 'btn-ghost bg-base-100'
                        " @click="
                            showAdvancedFilters =
                            !showAdvancedFilters
                            ">
                        <SlidersHorizontal :size="16" />

                        Filters

                        <span v-if="activeFilterCount > 0" class="badge badge-primary badge-xs">
                            {{ activeFilterCount }}
                        </span>
                    </button>

                    <!-- Clear -->
                    <button v-if="activeFilterCount > 0" type="button" class="btn btn-ghost btn-sm text-base-content/60"
                        @click="clearFilters">
                        Clear
                    </button>
                </div>
            </div>

            <!-- Advanced filters -->
            <div v-if="showAdvancedFilters" class="border-t border-base-300/50 px-4 py-4">
                <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
                    <!-- Amount -->
                    <div>
                        <label class="mb-1.5 block text-xs font-medium text-base-content/60">
                            Amount range
                        </label>

                        <div class="flex items-center gap-2">
                            <input v-model.number="filters.min_amount" type="number" placeholder="Min"
                                class="input input-md w-full bg-base-100" min="0" step="0.01" />

                            <span class="text-xs text-base-content/30">
                                —
                            </span>

                            <input v-model.number="filters.max_amount" type="number" placeholder="Max"
                                class="input input-md w-full bg-base-100" min="0" step="0.01" />
                        </div>
                    </div>

                    <!-- Date range -->
                    <div>
                        <label class="mb-1.5 block text-xs font-medium text-base-content/60">
                            Date range
                        </label>

                        <div class="flex items-center gap-2">
                            <input v-model="filters.start_date" type="date" class="input input-md w-full bg-base-100"
                                title="Start date" />

                            <span class="text-xs text-base-content/30">
                                —
                            </span>

                            <input v-model="filters.end_date" type="date" class="input input-md w-full bg-base-100"
                                title="End date" />
                        </div>
                    </div>

                    <!-- Sort -->
                    <div>
                        <label class="mb-1.5 block text-xs font-medium text-base-content/60">
                            Sort by
                        </label>

                        <select v-model="filters.sort_by" class="select select-md w-full bg-base-100">
                            <option value="date_desc">
                                Newest first
                            </option>

                            <option value="date_asc">
                                Oldest first
                            </option>

                            <option value="amount_desc">
                                Highest amount
                            </option>

                            <option value="amount_asc">
                                Lowest amount
                            </option>

                            <option value="description_asc">
                                Description (A–Z)
                            </option>
                        </select>
                    </div>

                    <!-- Date preset -->
                    <div>
                        <label class="mb-1.5 block text-xs font-medium text-base-content/60">
                            Quick date range
                        </label>

                        <select v-model="datePreset" class="select select-md w-full bg-base-100"
                            @change="applyDatePreset">
                            <option value="">
                                Custom date range
                            </option>

                            <option value="today">
                                Today
                            </option>

                            <option value="week">
                                This week
                            </option>

                            <option value="month">
                                This month
                            </option>

                            <option value="year">
                                This year
                            </option>

                            <option value="last_month">
                                Last month
                            </option>
                        </select>
                    </div>
                </div>

                <!-- Active filters -->
                <div v-if="activeFilterCount > 0" class="mt-4 flex flex-wrap items-center gap-2">
                    <span class="text-xs font-medium text-base-content/50">
                        Active filters:
                    </span>

                    <span v-if="filters.search" class="badge badge-sm border-0 bg-base-200 text-base-content/70">
                        Search: {{ filters.search }}
                    </span>

                    <span v-if="filters.category_id !== undefined"
                        class="badge badge-sm border-0 bg-base-200 text-base-content/70">
                        Category:
                        {{
                            getCategoryName(
                                filters.category_id,
                            )
                        }}
                    </span>

                    <span v-if="
                        filters.min_amount !==
                        undefined &&
                        filters.min_amount !== null
                    " class="badge badge-sm border-0 bg-base-200 text-base-content/70">
                        Min:
                        {{
                            formatCurrency(
                                filters.min_amount,
                            )
                        }}
                    </span>

                    <span v-if="
                        filters.max_amount !==
                        undefined &&
                        filters.max_amount !== null
                    " class="badge badge-sm border-0 bg-base-200 text-base-content/70">
                        Max:
                        {{
                            formatCurrency(
                                filters.max_amount,
                            )
                        }}
                    </span>

                    <span v-if="filters.start_date" class="badge badge-sm border-0 bg-base-200 text-base-content/70">
                        From:
                        {{ filters.start_date }}
                    </span>

                    <span v-if="filters.end_date" class="badge badge-sm border-0 bg-base-200 text-base-content/70">
                        To:
                        {{ filters.end_date }}
                    </span>
                </div>
            </div>
        </section>

        <section class="overflow-hidden rounded-xl border border-base-300/60 bg-base-100">
            <!-- Summary header -->
            <div
                class="flex flex-col gap-1 border-b border-base-300/50 px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                    <h2 class="text-sm font-semibold">
                        Filtered Summary
                    </h2>

                    <p class="mt-0.5 text-xs text-base-content/50">
                        Summary of the expenses currently loaded
                    </p>
                </div>

                <div class="text-xs text-base-content/50">
                    {{ totalCount }}
                    {{
                        totalCount === 1
                            ? 'matching expense'
                            : 'matching expenses'
                    }}
                </div>
            </div>

            <!-- Summary metrics -->
            <div class="grid grid-cols-2 divide-x divide-y divide-base-300/50 lg:grid-cols-4 lg:divide-y-0">
                <!-- Total -->
                <div class="p-5">
                    <div class="flex items-center gap-2 text-xs font-medium text-base-content/50">
                        <Wallet :size="14" />

                        Total
                    </div>

                    <p class="mt-2 text-xl font-bold tracking-tight text-primary">
                        {{ formatCurrency(expenseSummary?.total ?? 0) }}
                    </p>
                </div>

                <!-- Transactions -->
                <div class="p-5">
                    <div class="flex items-center gap-2 text-xs font-medium text-base-content/50">
                        <Receipt :size="14" />

                        Transactions
                    </div>

                    <p class="mt-2 text-xl font-bold tracking-tight">
                        {{ expenseSummary?.count }}
                    </p>
                </div>

                <!-- Average -->
                <div class="p-5">
                    <div class="flex items-center gap-2 text-xs font-medium text-base-content/50">
                        <Calculator :size="14" />

                        Average
                    </div>

                    <p class="mt-2 text-xl font-bold tracking-tight">
                        {{ formatCurrency(expenseSummary?.average ?? 0) }}
                    </p>
                </div>

                <!-- Largest -->
                <div class="p-5">
                    <div class="flex items-center gap-2 text-xs font-medium text-base-content/50">
                        <ArrowUp :size="14" />

                        Largest
                    </div>

                    <p class="mt-2 truncate text-xl font-bold tracking-tight"
                        :title="expenseSummary?.largest_description ?? ''">
                        {{
                            formatCurrency(
                                expenseSummary?.largest_amount ?? 0
                            )
                        }}
                    </p>

                    <p v-if="
                        expenseSummary?.largest_description
                    " class="mt-1 truncate text-xs text-base-content/40" :title="expenseSummary.largest_description
                        ">
                        {{
                            expenseSummary.largest_description
                        }}
                    </p>
                </div>
            </div>
        </section>

        <section class="overflow-hidden rounded-xl border border-base-300/60 bg-base-100">
            <div class="overflow-x-auto">
                <table class="table w-full">
                    <thead>
                        <tr class="border-base-300/50 text-xs uppercase tracking-wider text-base-content/40">
                            <th class="px-5 py-4 font-medium">
                                Date
                            </th>

                            <th class="px-5 py-4 font-medium">
                                Description
                            </th>

                            <th class="px-5 py-4 font-medium">
                                Category
                            </th>

                            <th class="px-5 py-4 text-right font-medium">
                                Amount
                            </th>

                            <th class="px-5 py-4 text-right font-medium">
                                Actions
                            </th>
                        </tr>
                    </thead>

                    <tbody>
                        <!-- Loading -->
                        <tr v-if="isLoading">
                            <td colspan="5" class="px-5 py-16">
                                <div class="flex flex-col items-center justify-center">
                                    <span class="loading loading-spinner loading-md text-primary"></span>

                                    <p class="mt-3 text-sm text-base-content/50">
                                        Loading expenses...
                                    </p>
                                </div>
                            </td>
                        </tr>

                        <!-- Empty -->
                        <tr v-else-if="
                            expenses.length === 0
                        ">
                            <td colspan="5" class="px-5 py-16">
                                <div class="flex flex-col items-center justify-center text-center">
                                    <div
                                        class="flex h-12 w-12 items-center justify-center rounded-full bg-base-200 text-base-content/40">
                                        <Receipt :size="20" />
                                    </div>

                                    <p class="mt-4 text-sm font-semibold">
                                        No expenses found
                                    </p>

                                    <p class="mt-1 max-w-sm text-xs text-base-content/50">
                                        Try adjusting your
                                        filters or add a new
                                        expense.
                                    </p>

                                    <button v-if="
                                        activeFilterCount >
                                        0
                                    " class="btn btn-ghost btn-sm mt-4" @click="
                                        clearFilters
                                    ">
                                        Clear filters
                                    </button>
                                </div>
                            </td>
                        </tr>

                        <!-- Expenses -->
                        <tr v-for="exp in expenses" v-else :key="exp.id"
                            class="border-base-300/40 transition-colors hover:bg-base-200/30">
                            <!-- Date -->
                            <td class="whitespace-nowrap px-5 py-4">
                                <span class="text-sm text-base-content/60">
                                    {{
                                        formatDate(
                                            exp.occurred_at,
                                        )
                                    }}
                                </span>
                            </td>

                            <!-- Description -->
                            <td class="px-5 py-4">
                                <span class="font-medium text-base-content">
                                    {{ exp.description }}
                                </span>
                            </td>

                            <!-- Category -->
                            <td class="px-5 py-4">
                                <span v-if="
                                    exp.category_id
                                " class="badge badge-sm border-0 bg-primary/10 text-primary">
                                    {{
                                        getCategoryName(
                                            exp.category_id,
                                        )
                                    }}
                                </span>

                                <span v-else class="badge badge-sm border-0 bg-base-200 text-base-content/50">
                                    Uncategorized
                                </span>
                            </td>

                            <!-- Amount -->
                            <td class="whitespace-nowrap px-5 py-4 text-right">
                                <span class="font-semibold tabular-nums" :class="exp.type ===
                                    'income'
                                    ? 'text-success'
                                    : 'text-error'
                                    ">
                                    {{
                                        exp.type ===
                                            'income'
                                            ? '+'
                                            : '-'
                                    }}{{
                                        formatCurrency(
                                            exp.amount,
                                        )
                                    }}
                                </span>
                            </td>

                            <!-- Actions -->
                            <td class="whitespace-nowrap px-5 py-4 text-right">
                                <div class="flex items-center justify-end gap-1">
                                    <button type="button"
                                        class="btn btn-ghost btn-xs text-base-content/60 hover:text-base-content"
                                        title="Edit expense" @click="
                                            openEditModal(
                                                exp,
                                            )
                                            ">
                                        <Pencil :size="14" />

                                        <span class="sr-only">
                                            Edit
                                        </span>
                                    </button>

                                    <button type="button"
                                        class="btn btn-ghost btn-xs text-error/70 hover:bg-error/10 hover:text-error"
                                        title="Delete expense" @click="
                                            deleteExpense(
                                                exp.id,
                                            )
                                            ">
                                        <Trash2 :size="14" />

                                        <span class="sr-only">
                                            Delete
                                        </span>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- =====================================================
                 PAGINATION
            ====================================================== -->
            <div v-if="totalCount > 0"
                class="flex flex-col gap-4 border-t border-base-300/50 px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
                <div class="flex flex-wrap items-center gap-4 text-xs text-base-content/50">
                    <span>
                        Showing
                        <span class="font-medium text-base-content/70">
                            {{ rangeStart }}–{{ rangeEnd }}
                        </span>
                        of
                        <span class="font-medium text-base-content/70">
                            {{ totalCount }}
                        </span>
                    </span>

                    <label class="flex items-center gap-2">
                        <span>Rows:</span>

                        <select v-model.number="pageSize"
                            class="select select-bordered select-xs border-base-300 bg-base-100" @change="goToPage(1)">
                            <option :value="10">
                                10
                            </option>

                            <option :value="20">
                                20
                            </option>

                            <option :value="50">
                                50
                            </option>

                            <option :value="100">
                                100
                            </option>
                        </select>
                    </label>
                </div>

                <div class="join">
                    <button class="join-item btn btn-sm btn-ghost" :disabled="currentPage === 1
                        " aria-label="First page" @click="goToPage(1)">
                        «
                    </button>

                    <button class="join-item btn btn-sm btn-ghost" :disabled="currentPage === 1
                        " aria-label="Previous page" @click="
                            goToPage(
                                currentPage - 1,
                            )
                            ">
                        ‹
                    </button>

                    <button v-for="page in visiblePages" :key="page" class="join-item btn btn-sm" :class="page === currentPage
                        ? 'btn-primary'
                        : 'btn-ghost'
                        " :disabled="page === '...'
                            " @click="
                                page !== '...' &&
                                goToPage(page)
                                ">
                        {{ page }}
                    </button>

                    <button class="join-item btn btn-sm btn-ghost" :disabled="currentPage ===
                        totalPages
                        " aria-label="Next page" @click="
                            goToPage(
                                currentPage + 1,
                            )
                            ">
                        ›
                    </button>

                    <button class="join-item btn btn-sm btn-ghost" :disabled="currentPage ===
                        totalPages
                        " aria-label="Last page" @click="
                            goToPage(totalPages)
                            ">
                        »
                    </button>
                </div>
            </div>
        </section>

        <ExpenseFormModal ref="expense-modal" :expense="selectedExpense" :categories="categories" :accounts="accounts"
            :is-saving="isSaving" @save="handleSave" @close="selectedExpense = null" />
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
    useExpenses,
} from '../composables/useExpenses';

import ExpenseFormModal from '../components/ExpenseFormModal.vue';

import type {
    Expense,
    CreateExpenseRequest,
} from '../types/index';

import {
    formatCurrency,
    formatDate,
} from '@/utils/helpers';

import {
    Plus,
    Search,
    SlidersHorizontal,
    X,
    Receipt,
    Pencil,
    Trash2,
    Wallet,
    Calculator,
    ArrowUp,
} from '@lucide/vue';

const {
    expenses,
    expenseSummary,
    categories,
    accounts,

    isLoading,
    isSaving,

    totalCount,
    currentPage,
    pageSize,

    filters,

    fetchAccounts,
    fetchExpenses,
    fetchFilteredSummary,
    fetchCategories,
    saveExpense,
    deleteExpense,
} = useExpenses();

const selectedExpense =
    ref<Expense | null>(null);

const datePreset = ref('');

const showAdvancedFilters =
    ref(true);

const activeFilterCount =
    computed(() => {
        let count = 0;

        if (filters.value.search) {
            count++;
        }

        if (
            filters.value.category_id !==
            undefined
        ) {
            count++;
        }

        if (
            filters.value.min_amount !==
            undefined &&
            filters.value.min_amount !== null
        ) {
            count++;
        }

        if (
            filters.value.max_amount !==
            undefined &&
            filters.value.max_amount !== null
        ) {
            count++;
        }

        if (filters.value.start_date) {
            count++;
        }

        if (filters.value.end_date) {
            count++;
        }

        return count;
    });


const summary = computed(() => {
    const records = expenses.value;

    if (records.length === 0) {
        return {
            total: 0,
            count: 0,
            average: 0,
            largestAmount: 0,
            largestDescription: '',
        };
    }

    const total = records.reduce(
        (sum, expense) => {
            return sum + expense.amount;
        },
        0,
    );

    const largest = records.reduce(
        (current, expense) => {
            return expense.amount >
                (current?.amount ?? 0)
                ? expense
                : current;
        },
        records[0],
    );

    return {
        total,

        count: records.length,

        average:
            records.length > 0
                ? total / records.length
                : 0,

        largestAmount:
            largest?.amount ?? 0,

        largestDescription:
            largest?.description ?? '',
    };
});

const getCategoryName = (
    id: string,
) => {
    const cat =
        categories.value.find(
            c => c.id === id,
        );

    return cat
        ? cat.name
        : 'Unknown';
};

const applyDatePreset = () => {
    const now = new Date();

    const toISO = (d: Date) =>
        d.toISOString().split('T')[0];

    switch (datePreset.value) {
        case 'today': {
            filters.value.start_date =
                toISO(now);

            filters.value.end_date =
                toISO(now);

            break;
        }

        case 'week': {
            const start =
                new Date(now);

            start.setDate(
                now.getDate() -
                now.getDay(),
            );

            filters.value.start_date =
                toISO(start);

            filters.value.end_date =
                toISO(now);

            break;
        }

        case 'month': {
            filters.value.start_date =
                `${now.getFullYear()}-${String(
                    now.getMonth() + 1,
                ).padStart(
                    2,
                    '0',
                )}-01`;

            filters.value.end_date =
                toISO(now);

            break;
        }

        case 'year': {
            filters.value.start_date =
                `${now.getFullYear()}-01-01`;

            filters.value.end_date =
                toISO(now);

            break;
        }

        case 'last_month': {
            const firstOfThisMonth =
                new Date(
                    now.getFullYear(),
                    now.getMonth(),
                    1,
                );

            const lastMonthEnd =
                new Date(
                    firstOfThisMonth,
                );

            lastMonthEnd.setDate(
                lastMonthEnd.getDate() -
                1,
            );

            const lastMonthStart =
                new Date(
                    lastMonthEnd.getFullYear(),
                    lastMonthEnd.getMonth(),
                    1,
                );

            filters.value.start_date =
                toISO(lastMonthStart);

            filters.value.end_date =
                toISO(lastMonthEnd);

            break;
        }

        default:
            break;
    }
};

const expenseModal = useTemplateRef<InstanceType<typeof ExpenseFormModal>>('expense-modal');

const openAddModal = () => {
    selectedExpense.value = null;
    expenseModal.value?.openModal();
};

const openEditModal = (exp: Expense,) => {
    selectedExpense.value = exp;
    expenseModal.value?.openModal();
};

const clearFilters = () => {
    datePreset.value = '';
    filters.value = {
        sort_by: 'date_desc',
    };
    currentPage.value = 1;
};

const handleSave = async (payload: CreateExpenseRequest) => {
    const success = await saveExpense(payload, selectedExpense.value?.id);

    if (success) {
        expenseModal.value?.closeModal();

        selectedExpense.value =
            null;
    }
};

const totalPages = computed(() => Math.max(
    1,
    Math.ceil(
        totalCount.value /
        pageSize.value,
    ),
),
);

const rangeStart = computed(() =>
    totalCount.value === 0
        ? 0
        : (currentPage.value - 1) *
        pageSize.value +
        1,
);

const rangeEnd = computed(() =>
    Math.min(
        currentPage.value *
        pageSize.value,
        totalCount.value,
    ),
);

function goToPage(page: number | string): void {
    if (typeof page !== 'number') {
        return;
    }

    if (
        page < 1 ||
        page > totalPages.value ||
        page === currentPage.value
    ) {
        return;
    }

    currentPage.value = page;
}

const visiblePages = computed<(number | string)[]>(() => {
    const total = totalPages.value;
    const current = currentPage.value;
    const delta = 1;

    if (total <= 7) {
        return Array.from(
            {
                length: total,
            },
            (_, i) => i + 1,
        );
    }

    const pages: (
        | number
        | string
    )[] = [];

    const left = Math.max(2, current - delta,);

    const right = Math.min(
        total - 1,
        current + delta,
    );

    pages.push(1);

    if (left > 2) {
        pages.push('...');
    }

    for (let i = left; i <= right; i++) {
        pages.push(i);
    }

    if (right < total - 1) {
        pages.push('...');
    }

    pages.push(total);

    return pages;
});

watch(
    [
        currentPage,
        pageSize,
        filters,
    ],
    (
        newValues,
        oldValues,
    ) => {
        const [
            newPage,
            newPageSize,
            newFilters,
        ] = newValues;

        const [
            oldPage,
            oldPageSize,
            oldFilters,
        ] = oldValues;

        const filtersChanged =
            JSON.stringify(
                newFilters,
            ) !==
            JSON.stringify(
                oldFilters,
            );

        if (filtersChanged && newPage !== 1) {
            currentPage.value = 1;
            return;
        }

        filters.value.limit = newPageSize;

        filters.value.offset = (newPage - 1) * newPageSize;

        fetchExpenses();
    },
    {
        deep: true,
    },
);

onMounted(() => {
    filters.value = {
        sort_by: 'date_desc',
    };

    fetchCategories();
    fetchExpenses();
    fetchAccounts();
});
</script>