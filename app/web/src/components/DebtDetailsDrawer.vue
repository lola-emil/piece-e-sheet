<template>
    <div class="drawer drawer-end">
        <input ref="drawer-toggle" :id="drawerId" type="checkbox" class="drawer-toggle" />
        <div class="drawer-side z-30">
            <label :for="drawerId" class="drawer-overlay"></label>

            <div class="menu p-4 w-80 min-h-full bg-base-200 text-base-content overflow-y-auto">
                <div v-if="debt" class="space-y-6">
                    <div class="flex justify-between items-center">
                        <h3 class="font-bold text-xl">Debt Details</h3>
                        <label :for="drawerId" class="btn btn-sm btn-circle btn-ghost">✕</label>
                    </div>

                    <div class="card bg-base-100 shadow-sm">
                        <div class="card-body p-4">
                            <h4 class="font-bold text-lg">{{ debt.person_name }}</h4>
                            <p class="text-sm text-base-content/60">{{ debt.description }}</p>
                            <div class="divider my-2"></div>
                            <div class="flex justify-between">
                                <span class="text-sm">Total:</span>
                                <span class="font-bold">{{ formatCurrency(debt.amount) }}</span>
                            </div>
                            <div class="flex justify-between">
                                <span class="text-sm">Remaining:</span>
                                <span class="font-bold" :class="debt.type === 'i_owe' ? 'text-error' : 'text-success'">
                                    {{ formatCurrency(debt.remaining_amount) }}
                                </span>
                            </div>
                            <div class="badge mt-2" :class="getStatusBadge(debt.status)">{{ debt.status }}</div>
                        </div>
                    </div>

                    <div v-if="debt.status !== 'paid'" class="card bg-base-100 shadow-sm">
                        <div class="card-body p-4">
                            <h4 class="font-bold">Record Payment</h4>
                            <form @submit.prevent="handlePayment" class="space-y-3 mt-2">
                                <input v-model.number="paymentForm.amount" type="number" step="0.01"
                                    placeholder="Amount" class="input input-bordered w-full text-sm" required />
                                <input v-model="paymentForm.note" type="text" placeholder="Note (optional)"
                                    class="input input-bordered w-full text-sm" />
                                <button type="submit" class="btn btn-primary btn-sm w-full"
                                    :class="{ 'loading': isSaving }" :disabled="isSaving">
                                    Add Payment
                                </button>
                            </form>
                        </div>
                    </div>

                    <div>
                        <h4 class="font-bold mb-3">Payment History</h4>
                        <div v-if="isLoading" class="text-center py-4"><span
                                class="loading loading-spinner loading-sm"></span></div>
                        <div v-else-if="payments.length === 0" class="text-center text-sm text-base-content/50 py-4">No
                            payments yet.</div>
                        <ul v-else class="space-y-2">
                            <li v-for="p in payments" :key="p.id"
                                class="bg-base-100 p-3 rounded-lg shadow-sm flex justify-between items-center">
                                <div>
                                    <div class="font-bold text-sm">{{ formatCurrency(p.amount) }}</div>
                                    <div class="text-xs text-base-content/50">{{ formatDate(p.paid_at) }}</div>
                                </div>
                                <div class="text-xs text-right max-w-25 truncate" :title="p.note">{{ p.note || '-'
                                    }}</div>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, useTemplateRef, watch } from 'vue';
import { usePayments } from '../composables/usePayments';
import type { Debt, CreatePaymentRequest } from '../types';
import { formatCurrency, formatDate } from '@/utils/helpers';

const props = defineProps<{
    drawerId: string;
    debt: Debt | null;
}>();

const paymentForm = ref({ amount: 0, note: '' });

const debtIdRef = ref(props.debt?.id || '');
watch(() => props.debt?.id, (newId) => { debtIdRef.value = newId || ''; });

const { payments, isLoading, isSaving, fetchPayments, recordPayment } = usePayments(debtIdRef.value);

watch(() => props.debt, (newDebt) => {
    if (newDebt) {
        debtIdRef.value = newDebt.id;
        fetchPayments();
        paymentForm.value = { amount: 0, note: '' };
    }
});

const handlePayment = async () => {
    if (!props.debt) return;
    const payload: CreatePaymentRequest = {
        debt_id: props.debt.id,
        amount: paymentForm.value.amount,
        note: paymentForm.value.note
    };

    const success = await recordPayment(payload);
    if (success) {
        paymentForm.value = { amount: 0, note: '' };
    }
};

const getStatusBadge = (status: string) => {
    if (status === 'paid') return 'badge-success text-white';
    if (status === 'partial') return 'badge-warning text-white';
    return 'badge-ghost';
};

const drawer = useTemplateRef('drawer-toggle');

const openDrawer = () => {
    drawer.value!.checked = true;
}

defineExpose({
    openDrawer,
})
</script>