import { ref } from 'vue';
import api from '../services/api';
import type { Expense, Category, ExpenseFilter, CreateExpenseRequest, Account } from '../types';

export function useExpenses() {
    const expenses = ref<Expense[]>([]);
    const categories = ref<Category[]>([]);
    const accounts = ref<Account[]>([]);
    const isLoading = ref(false);
    const isSaving = ref(false);
    const totalCount = ref(0);
    const currentPage = ref(1);
    const pageSize = ref(20)


    const fetchExpenses = async (filter: Partial<ExpenseFilter> = {}) => {
        isLoading.value = true;
        try {
            filter.limit = pageSize.value;
            filter.offset = ((currentPage.value - 1) * pageSize.value);

            const { data } = await api.get('/api/expenses', { params: filter });

            expenses.value = data.data ?? [];
            totalCount.value = data.count ?? 0;
        } catch (error) {
            console.error('Failed to fetch expenses', error);
        } finally {
            isLoading.value = false;
        }
    };

    const fetchAccounts = async () => {
        try {
            const { data } = await api.get('/api/accounts');
            accounts.value = data.data;
        } catch (error) {
            console.error('Failed to fetch accounts', error);
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
            await fetchExpenses();
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
            await fetchExpenses();
        } catch (error) {
            console.error('Failed to delete expense', error);
        }
    };

    return {
        expenses,
        accounts,
        categories,
        isLoading,
        isSaving,
        totalCount,
        currentPage,
        pageSize,
        fetchExpenses,
        fetchCategories,
        fetchAccounts,
        saveExpense,
        deleteExpense
    };
}