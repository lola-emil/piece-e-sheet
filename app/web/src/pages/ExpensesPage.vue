<template>
    <div class="space-y-6">
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <h1 class="text-3xl font-bold">Expenses</h1>
            <button class="btn btn-primary" @click="openAddModal">
                + Add Expense
            </button>
        </div>

        <div class="card bg-base-100 shadow-sm">
            <div class="card-body p-4">
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                    <input v-model="searchQuery" type="text" placeholder="Search description..."
                        class="input input-bordered w-full" />

                    <select v-model="filters.category_id" class="select select-bordered w-full">
                        <option :value="undefined">All Categories</option>
                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                            {{ cat.name }}
                        </option>
                    </select>

                    <div class="flex gap-2">
                        <input v-model.number="filters.min_amount" type="number" placeholder="Min $"
                            class="input input-bordered w-full" min="0" step="0.01" />
                        <input v-model.number="filters.max_amount" type="number" placeholder="Max $"
                            class="input input-bordered w-full" min="0" step="0.01" />
                    </div>

                    <div class="flex gap-2">
                        <input v-model="filters.start_date" type="date" class="input input-bordered w-full"
                            title="Start date" />
                        <input v-model="filters.end_date" type="date" class="input input-bordered w-full"
                            title="End date" />
                    </div>

                    <select v-model="filters.sort_by" class="select select-bordered w-full">
                        <option value="date_desc">Newest first</option>
                        <option value="date_asc">Oldest first</option>
                        <option value="amount_desc">Highest amount</option>
                        <option value="amount_asc">Lowest amount</option>
                        <option value="description_asc">Description (A–Z)</option>
                    </select>

                    <select v-model="datePreset" class="select select-bordered w-full" @change="applyDatePreset">
                        <option value="">Custom date range</option>
                        <option value="today">Today</option>
                        <option value="week">This week</option>
                        <option value="month">This month</option>
                        <option value="year">This year</option>
                        <option value="last_month">Last month</option>
                    </select>

                    <div class="flex items-center gap-2">
                        <button class="btn btn-ghost btn-sm" @click="clearFilters">
                            Clear filters
                        </button>
                        <span v-if="activeFilterCount > 0" class="badge badge-neutral">
                            {{ activeFilterCount }} active
                        </span>
                    </div>
                </div>
            </div>
        </div>

        <div class="card bg-base-100 shadow-xl overflow-hidden">
            <div class="overflow-x-auto">
                <table class="table table-zebra w-full">
                    <thead>
                        <tr>
                            <th>Date</th>
                            <th>Description</th>
                            <th>Category</th>
                            <th class="text-right">Amount</th>
                            <th class="text-center">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="isLoading">
                            <td colspan="5" class="text-center py-8">
                                <span class="loading loading-spinner loading-lg"></span>
                            </td>
                        </tr>
                        <tr v-else-if="filteredExpenses.length === 0">
                            <td colspan="5" class="text-center py-8 text-base-content/50">
                                No expenses found.
                            </td>
                        </tr>
                        <tr v-for="exp in filteredExpenses" :key="exp.id">
                            <td class="whitespace-nowrap">{{ formatDate(exp.occurred_at) }}</td>
                            <td>{{ exp.description }}</td>
                            <td>
                                <span v-if="exp.category_id" class="badge badge-primary badge-outline">
                                    {{ getCategoryName(exp.category_id) }}
                                </span>
                                <span v-else class="badge badge-ghost">Uncategorized</span>
                            </td>
                            <td class="text-right font-bold text-error whitespace-nowrap">
                                -{{ formatCurrency(exp.amount) }}
                            </td>
                            <td class="text-center whitespace-nowrap">
                                <button class="btn btn-xs btn-ghost" @click="openEditModal(exp)">Edit</button>
                                <button class="btn btn-xs btn-ghost text-error"
                                    @click="deleteExpense(exp.id)">Delete</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <ExpenseFormModal ref="expense-modal" :expense="selectedExpense" :categories="categories" :is-saving="isSaving"
            @save="handleSave" @close="selectedExpense = null" />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, useTemplateRef } from 'vue';
import { useExpenses } from '../composables/useExpenses';
import ExpenseFormModal from '../components/ExpenseFormModal.vue';
import type { Expense, ExpenseFilter, CreateExpenseRequest } from '../types/index';
import { formatCurrency, formatDate } from '@/utils/helpers';

const {
    expenses, categories, isLoading, isSaving,
    fetchExpenses, fetchCategories, saveExpense, deleteExpense
} = useExpenses();

const selectedExpense = ref<Expense | null>(null);
const searchQuery = ref('');
const datePreset = ref('');
const filters = ref<ExpenseFilter & {
    min_amount?: number;
    max_amount?: number;
    sort_by?: string;
}>({
    sort_by: 'date_desc',
});

const activeFilterCount = computed(() => {
    let count = 0;
    if (searchQuery.value) count++;
    if (filters.value.category_id !== undefined) count++;
    if (filters.value.min_amount !== undefined && filters.value.min_amount !== null) count++;
    if (filters.value.max_amount !== undefined && filters.value.max_amount !== null) count++;
    if (filters.value.start_date) count++;
    if (filters.value.end_date) count++;
    return count;
});

const filteredExpenses = computed(() => {
    const q = searchQuery.value.toLowerCase().trim();
    const { category_id, start_date, end_date, min_amount, max_amount, sort_by } = filters.value;

    let result = expenses.value.filter(exp => {
        if (q && !exp.description.toLowerCase().includes(q)) return false;
        if (category_id !== undefined && exp.category_id !== category_id) return false;
        if (start_date && new Date(exp.occurred_at) < new Date(start_date)) return false;
        if (end_date) {
            const end = new Date(end_date);
            end.setHours(23, 59, 59, 999);
            if (new Date(exp.occurred_at) > end) return false;
        }
        if (min_amount !== undefined && min_amount !== null && exp.amount < min_amount) return false;
        if (max_amount !== undefined && max_amount !== null && exp.amount > max_amount) return false;
        return true;
    });

    switch (sort_by) {
        case 'date_asc':
            result.sort((a, b) => new Date(a.occurred_at).getTime() - new Date(b.occurred_at).getTime());
            break;
        case 'date_desc':
            result.sort((a, b) => new Date(b.occurred_at).getTime() - new Date(a.occurred_at).getTime());
            break;
        case 'amount_asc':
            result.sort((a, b) => a.amount - b.amount);
            break;
        case 'amount_desc':
            result.sort((a, b) => b.amount - a.amount);
            break;
        case 'description_asc':
            result.sort((a, b) => a.description.localeCompare(b.description));
            break;
    }

    return result;
});

const applyDatePreset = () => {
    const now = new Date();
    const toISO = (d: Date) => d.toISOString().split('T')[0];

    switch (datePreset.value) {
        case 'today':
            filters.value.start_date = toISO(now);
            filters.value.end_date = toISO(now);
            break;
        case 'week': {
            const start = new Date(now);
            start.setDate(now.getDate() - now.getDay());
            filters.value.start_date = toISO(start);
            filters.value.end_date = toISO(now);
            break;
        }
        case 'month':
            filters.value.start_date = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`;
            filters.value.end_date = toISO(now);
            break;
        case 'year':
            filters.value.start_date = `${now.getFullYear()}-01-01`;
            filters.value.end_date = toISO(now);
            break;
        case 'last_month': {
            const firstOfThisMonth = new Date(now.getFullYear(), now.getMonth(), 1);
            const lastMonthEnd = new Date(firstOfThisMonth);
            lastMonthEnd.setDate(lastMonthEnd.getDate() - 1);
            const lastMonthStart = new Date(lastMonthEnd.getFullYear(), lastMonthEnd.getMonth(), 1);
            filters.value.start_date = toISO(lastMonthStart);
            filters.value.end_date = toISO(lastMonthEnd);
            break;
        }
        default:
            break;
    }
};

const expenseModal = useTemplateRef<InstanceType<typeof ExpenseFormModal>>('expense-modal');

const clearFilters = () => {
    searchQuery.value = '';
    datePreset.value = '';
    filters.value = { sort_by: 'date_desc' };
};

const openAddModal = () => {
    selectedExpense.value = null;
    expenseModal.value?.openModal();
};

const openEditModal = (exp: Expense) => {
    selectedExpense.value = exp;
    expenseModal.value?.openModal();
};

const handleSave = async (payload: CreateExpenseRequest) => {
    const success = await saveExpense(payload, selectedExpense.value?.id);
    if (success) {
        expenseModal.value?.closeModal();
    }
};

const getCategoryName = (id: string) => {
    const cat = categories.value.find(c => c.id === id);
    return cat ? cat.name : 'Unknown';
};

onMounted(() => {
    fetchCategories();
    fetchExpenses();
});
</script>