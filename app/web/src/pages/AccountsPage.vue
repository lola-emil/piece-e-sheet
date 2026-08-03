<template>
    <div class="space-y-6">
        <div class="flex justify-between items-center">
            <h1 class="text-3xl font-bold">Accounts</h1>
            <button class="btn btn-primary" @click="openModal()">+ Add Account</button>
        </div>

        <div v-if="isLoading" class="flex justify-center py-12">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>

        <div v-else-if="accounts?.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <div v-for="acc in accounts" :key="acc.id"
                class="card bg-base-100 shadow-md hover:shadow-lg transition-shadow">
                <div class="card-body p-4">
                    <div class="flex justify-between items-start">
                        <div class="flex items-center gap-3">
                            <div class="p-2 rounded-lg bg-primary/10 text-primary">
                                <Banknote v-if="acc.type === 'cash'" />
                                <CreditCard v-else />
                            </div>
                            <div>
                                <h2 class="card-title text-base">{{ acc.name }}</h2>
                                <p class="text-xs text-base-content/60 capitalize">{{ acc.type.replace('_', ' ') }}</p>
                            </div>
                        </div>

                        <div class="dropdown dropdown-end">
                            <label tabindex="0" class="btn btn-ghost btn-xs btn-circle">
                                <EllipsisVertical :size="18" />
                            </label>
                            <ul tabindex="0"
                                class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-32 z-10 border border-base-200">
                                <li><a @click="openModal(acc)">Edit</a></li>
                                <li><a class="text-error" @click="deleteAccount(acc.id)">Delete</a></li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div v-else class="text-center py-12 bg-base-100 rounded-lg shadow-sm">
            <p class="text-base-content/60">No accounts found. Add one to start tracking!</p>
        </div>

        <AccountFormModal ref="account-modal" :account="selectedAccount" :is-saving="isSaving" @save="handleSave"
            @close="selectedAccount = null" />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, useTemplateRef } from 'vue';
import { useAccounts } from '../composables/useAccounts';
import AccountFormModal from '../components/AccountFormModal.vue';
import type { Account, CreateAccountRequest } from '../types';
import { Banknote, CreditCard, EllipsisVertical } from '@lucide/vue';

const { accounts, isLoading, isSaving, fetchAccounts, saveAccount, deleteAccount } = useAccounts();
const selectedAccount = ref<Account | null>(null);

const accountModal = useTemplateRef("account-modal");

const openModal = (acc: Account | null = null) => {
    selectedAccount.value = acc;
    accountModal.value?.openModal();
};

const handleSave = async (payload: CreateAccountRequest) => {
    const success = await saveAccount(payload, selectedAccount.value?.id);
    if (success) {
        accountModal.value?.closeModal();
    }
};

onMounted(() => { fetchAccounts(); });
</script>