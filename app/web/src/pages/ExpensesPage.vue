<template>
    <div class="space-y-6">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <h1 class="text-3xl font-bold">Expenses</h1>
            <button class="btn btn-primary" @click="openAddModal">
                + Add Expense
            </button>
        </div>

        <!-- Filters -->
        <div class="card bg-base-100 shadow-sm">
            <div class="card-body p-4">
                <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <!-- Search -->
                    <input v-model="searchQuery" type="text" placeholder="Search description..."
                        class="input input-bordered w-full" />

                    <!-- Category Filter -->
                    <select v-model="filters.category_id" class="select select-bordered w-full">
                        <option :value="undefined">All Categories</option>
                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                            {{ cat.name }}
                        </option>
                    </select>

                    <!-- Date Range -->
                    <input v-model="filters.start_date" type="date" class="input input-bordered w-full" />
                    <input v-model="filters.end_date" type="date" class="input input-bordered w-full" />
                </div>
            </div>
        </div>

        <!-- Table -->
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

        <!-- Add/Edit Modal -->
        <ExpenseFormModal modalId="expense_modal" :expense="selectedExpense" :categories="categories"
            :is-saving="isSaving" @save="handleSave" @close="selectedExpense = null" />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useExpenses } from '../composables/useExpenses.ts';
import ExpenseFormModal from '../components/ExpenseFormModal.vue';
import type { Expense, ExpenseFilter, CreateExpenseRequest } from '../types/index.ts';

const {
    expenses, categories, isLoading, isSaving,
    fetchExpenses, fetchCategories, saveExpense, deleteExpense
} = useExpenses();

const selectedExpense = ref<Expense | null>(null);
const searchQuery = ref('');
const filters = ref<ExpenseFilter>({});

// Apply filters and search
const filteredExpenses = computed(() => {
    return expenses.value.filter(exp => {
        const matchesSearch = exp.description.toLowerCase().includes(searchQuery.value.toLowerCase());
        return matchesSearch;
    });
});

// Modal Handlers
const openAddModal = () => {
    selectedExpense.value = null;
    const modal = document.getElementById('expense_modal') as HTMLDialogElement;
    if (modal) modal.showModal();
};

const openEditModal = (exp: Expense) => {
    selectedExpense.value = exp;
    const modal = document.getElementById('expense_modal') as HTMLDialogElement;
    if (modal) modal.showModal();
};

const handleSave = async (payload: CreateExpenseRequest) => {
    const success = await saveExpense(payload, selectedExpense.value?.id);
    if (success) {
        const modal = document.getElementById('expense_modal') as HTMLDialogElement;
        if (modal) modal.close();
    }
};

// Helpers
const getCategoryName = (id: string) => {
    const cat = categories.value.find(c => c.id === id);
    return cat ? cat.name : 'Unknown';
};

const formatCurrency = (val: number) => {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);
};

const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString();
};

// Load data
onMounted(() => {
    fetchCategories();
    fetchExpenses();
});
</script>