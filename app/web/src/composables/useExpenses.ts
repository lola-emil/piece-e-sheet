import { ref } from 'vue';
import api from '../services/api';
import type { Expense, Category, ExpenseFilter, CreateExpenseRequest } from '../types';

export function useExpenses() {
    const expenses = ref<Expense[]>([]);
    const categories = ref<Category[]>([]);
    const isLoading = ref(false);
    const isSaving = ref(false);

    const fetchExpenses = async (filter: ExpenseFilter = {}) => {
        isLoading.value = true;
        try {
            const { data } = await api.get('/api/expenses', { params: filter });
            expenses.value = data.data ?? [];
        } catch (error) {
            console.error('Failed to fetch expenses', error);
        } finally {
            isLoading.value = false;
        }
    };

    const fetchCategories = async () => {
        try {
            const { data } = await api.get('/api/categories');
            categories.value = data.data;
        } catch (error) {
            console.error('Failed to fetch categories', error);
        }
    };

    const saveExpense = async (payload: CreateExpenseRequest, id?: string) => {
        isSaving.value = true;
        try {
            if (id) {
                await api.put(`/api/expenses/${id}`, payload);
            } else {
                await api.post('/api/expenses', payload);
            }
            await fetchExpenses(); // Refresh list
            return true;
        } catch (error) {
            console.error('Failed to save expense', error);
            return false;
        } finally {
            isSaving.value = false;
        }
    };

    const deleteExpense = async (id: string) => {
        if (!confirm('Are you sure you want to delete this expense?')) return;

        try {
            await api.delete(`/api/expenses/${id}`);
            await fetchExpenses(); // Refresh list
        } catch (error) {
            console.error('Failed to delete expense', error);
        }
    };

    return {
        expenses,
        categories,
        isLoading,
        isSaving,
        fetchExpenses,
        fetchCategories,
        saveExpense,
        deleteExpense
    };
}