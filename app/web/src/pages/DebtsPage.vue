<template>
    <div class="space-y-6">
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <h1 class="text-3xl font-bold">Debts</h1>
            <button class="btn btn-primary" @click="openModal()"><Plus :size="16"/> Add Debt</button>
        </div>

        <div role="tablist" class="tabs tabs-boxed bg-base-100 p-2 w-fit">
            <a role="tab" class="tab" :class="{ 'tab-active': filter === 'all' }" @click="filter = 'all'">All</a>
            <a role="tab" class="tab" :class="{ 'tab-active': filter === 'i_owe' }" @click="filter = 'i_owe'">I Owe</a>
            <a role="tab" class="tab" :class="{ 'tab-active': filter === 'owed_to_me' }"
                @click="filter = 'owed_to_me'">Owed to Me</a>
        </div>

        <div v-if="isLoading" class="flex justify-center py-12">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>

        <div v-else-if="filteredDebts?.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div v-for="debt in filteredDebts" :key="debt.id"
                class="card bg-base-100 shadow-md hover:shadow-lg transition-shadow border-l-4 cursor-pointer"
                :class="debt.type === 'i_owe' ? 'border-error' : 'border-success'" @click="openDrawer(debt)">
                <div class="card-body p-4">
                    <div class="flex justify-between items-start">
                        <div>
                            <h2 class="card-title text-lg">{{ debt.person_name }}</h2>
                            <p class="text-sm text-base-content/60">{{ debt.description }}</p>
                        </div>
                        <div class="dropdown dropdown-end">
                            <label tabindex="0" class="btn btn-ghost btn-xs btn-circle">
                                <EllipsisVertical :size="16" />
                            </label>
                            <ul tabindex="0"
                                class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-32 z-10 border border-base-200">
                                <li><a @click="openModal(debt)">Edit</a></li>
                                <li><a class="text-error" @click="deleteDebt(debt.id)">Delete</a></li>
                            </ul>
                        </div>
                    </div>

                    <div class="mt-4 flex justify-between items-end">
                        <div>
                            <div class="text-xs text-base-content/50">Remaining</div>
                            <div class="text-xl font-bold"
                                :class="debt.type === 'i_owe' ? 'text-error' : 'text-success'">
                                {{ formatCurrency(debt.remaining_amount) }}
                            </div>
                        </div>

                        <div class="text-right">
                            <div class="badge" :class="getStatusBadge(debt.status)">{{ debt.status }}</div>
                            <div v-if="debt.due_date" class="text-xs mt-1 text-base-content/50">
                                Due: {{ formatDate(debt.due_date) }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div v-else class="text-center py-12 bg-base-100 rounded-lg shadow-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-base-content/30 mb-4" fill="none"
                viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <p class="text-base-content/60">No debts found.</p>
        </div>

        <DebtFormModal ref="debt-modal" modalId="debt_modal" :debt="selectedDebt" :is-saving="isSaving"
            @save="handleSave" @close="selectedDebt = null" />

        <DebtDetailsDrawer ref="debt-details-drawer" drawerId="debt_details_drawer" :debt="selectedDebtForDrawer" />

    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, useTemplateRef } from 'vue';
import { useDebts } from '../composables/useDebts';
import DebtFormModal from '../components/DebtFormModal.vue';
import type { Debt, CreateDebtRequest, UpdateDebtRequest } from '../types';
import { formatCurrency, formatDate } from '@/utils/helpers.ts';
import DebtDetailsDrawer from '../components/DebtDetailsDrawer.vue';
import { EllipsisVertical, Plus } from '@lucide/vue';

const selectedDebtForDrawer = ref<Debt | null>(null);

const { debts, isLoading, isSaving, fetchDebts, saveDebt, deleteDebt } = useDebts();

const selectedDebt = ref<Debt | null>(null);
const filter = ref<'all' | 'i_owe' | 'owed_to_me'>('all');

const filteredDebts = computed(() => {
    if (filter.value === 'all') return debts.value;
    return debts.value.filter(d => d.type === filter.value);
});

const getStatusBadge = (status: string) => {
    if (status === 'paid') return 'badge-success text-white';
    if (status === 'partial') return 'badge-warning text-white';
    return 'badge-ghost';
};

const debtModal = useTemplateRef<InstanceType<typeof DebtFormModal>>('debt-modal')

const openModal = (debt: Debt | null = null) => {
    selectedDebt.value = debt;
    debtModal.value?.openModal();
};

const handleSave = async (payload: CreateDebtRequest | UpdateDebtRequest) => {
    const success = await saveDebt(payload, selectedDebt.value?.id);
    if (success) {
        debtModal.value?.closeModal();
    }
};

const drawer = useTemplateRef<InstanceType<typeof DebtDetailsDrawer>>('debt-details-drawer');

const openDrawer = (debt: Debt) => {
    selectedDebtForDrawer.value = debt;
    drawer.value?.openDrawer();
};

onMounted(() => {
    fetchDebts();
});
</script>