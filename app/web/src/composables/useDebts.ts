import { ref } from 'vue';
import api from '../services/api';
import type { Debt, CreateDebtRequest, UpdateDebtRequest } from '../types';

export function useDebts() {
    const debts = ref<Debt[]>([]);
    const isLoading = ref(false);
    const isSaving = ref(false);

    const fetchDebts = async () => {
        isLoading.value = true;
        try {
            const { data } = await api.get('/api/debts');
            debts.value = data.data;
        } catch (error) {
            console.error('Failed to fetch debts', error);
        } finally {
            isLoading.value = false;
        }
    };

    const saveDebt = async (payload: CreateDebtRequest | UpdateDebtRequest, id?: string) => {
        isSaving.value = true;
        try {
            if (id) {
                await api.put(`/api/debts/${id}`, payload);
            } else {
                await api.post('/api/debts', payload);
            }
            await fetchDebts();
            return true;
        } catch (error) {
            console.error('Failed to save debt', error);
            return false;
        } finally {
            isSaving.value = false;
        }
    };

    const deleteDebt = async (id: string) => {
        if (!confirm('Are you sure you want to delete this debt?')) return;
        try {
            await api.delete(`/api/debts/${id}`);
            await fetchDebts();
        } catch (error) {
            console.error('Failed to delete debt', error);
        }
    };

    return { debts, isLoading, isSaving, fetchDebts, saveDebt, deleteDebt };
}