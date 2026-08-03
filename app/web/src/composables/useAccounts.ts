import { ref } from 'vue';
import api from '../services/api';
import type { Account, CreateAccountRequest } from '../types';

export function useAccounts() {
    const accounts = ref<Account[]>([]);
    const isLoading = ref(false);
    const isSaving = ref(false);

    const fetchAccounts = async () => {
        isLoading.value = true;
        try {
            const { data } = await api.get('/api/accounts');
            accounts.value = data.data;
        } catch (error) {
            console.error('Failed to fetch accounts', error);
        } finally {
            isLoading.value = false;
        }
    };

    const saveAccount = async (payload: CreateAccountRequest, id?: string) => {
        isSaving.value = true;
        try {
            if (id) {
                await api.put(`/api/accounts/${id}`, payload);
            } else {
                await api.post('/api/accounts', payload);
            }
            await fetchAccounts();
            return true;
        } catch (error) {
            console.error('Failed to save account', error);
            return false;
        } finally {
            isSaving.value = false;
        }
    };

    const deleteAccount = async (id: string) => {
        if (!confirm('Are you sure? Expenses linked to this account will become uncategorized.')) return;
        try {
            await api.delete(`/api/accounts/${id}`);
            await fetchAccounts();
        } catch (error) {
            console.error('Failed to delete account', error);
        }
    };

    return { accounts, isLoading, isSaving, fetchAccounts, saveAccount, deleteAccount };
}