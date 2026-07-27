import { ref } from 'vue';
import api from '../services/api';
import type { Payment, CreatePaymentRequest } from '../types';

export function usePayments(debtId: string) {
    const payments = ref<Payment[]>([]);
    const isLoading = ref(false);
    const isSaving = ref(false);

    const fetchPayments = async () => {
        if (!debtId) return;
        isLoading.value = true;
        try {
            const { data } = await api.get(`/api/payments/debt/${debtId}`);
            payments.value = data.data;
        } catch (error) {
            console.error('Failed to fetch payments', error);
        } finally {
            isLoading.value = false;
        }
    };

    const recordPayment = async (payload: CreatePaymentRequest) => {
        isSaving.value = true;
        try {
            await api.post('/api/payments', payload);
            await fetchPayments();
            return true;
        } catch (error) {
            console.error('Failed to record payment', error);
            return false;
        } finally {
            isSaving.value = false;
        }
    };

    return { payments, isLoading, isSaving, fetchPayments, recordPayment };
}