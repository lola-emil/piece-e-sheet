<template>
    <dialog ref="account-modal" class="modal">
        <div class="modal-box">
            <h3 class="font-bold text-lg">{{ isEdit ? 'Edit Account' : 'Add Account' }}</h3>

            <form @submit.prevent="handleSubmit" class="space-y-4 mt-4">
                <!-- Name -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Account Name</span></label>
                    <input v-model="form.name" type="text" placeholder="e.g., Chase Checking, Cash"
                        class="input input-bordered w-full" required />
                </div>

                <!-- Type -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Account Type</span></label>
                    <select v-model="form.type" class="select select-bordered w-full">
                        <option value="cash">Cash</option>
                        <option value="checking">Checking</option>
                        <option value="savings">Savings</option>
                        <option value="credit_card">Credit Card</option>
                    </select>
                </div>

                <!-- Actions -->
                <div class="modal-action">
                    <button type="button" class="btn" @click="closeModal">Cancel</button>
                    <button type="submit" class="btn btn-primary" :class="{ 'loading': isSaving }" :disabled="isSaving">
                        {{ isEdit ? 'Update' : 'Save' }}
                    </button>
                </div>
            </form>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button @click="closeModal">close</button>
        </form>
    </dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed, useTemplateRef } from 'vue';
import type { Account, CreateAccountRequest } from '../types';

const props = defineProps<{
    account: Account | null;
    isSaving: boolean;
}>();

const emit = defineEmits<{
    (e: 'save', payload: CreateAccountRequest): void;
    (e: 'close'): void;
}>();

const isEdit = computed(() => !!props.account);

const form = ref({
    name: '',
    type: 'checking'
});

watch(() => props.account, (newAcc) => {
    if (newAcc) {
        form.value = { name: newAcc.name, type: newAcc.type };
    } else {
        form.value = { name: '', type: 'checking' };
    }
}, { immediate: true });

const handleSubmit = () => {
    emit('save', form.value);
};

const modal = useTemplateRef("account-modal");

const closeModal = () => {
    modal.value?.close();
    emit('close');
};

const openModal = () => {
    modal.value?.showModal();
}

defineExpose({
    openModal,
    closeModal,
})
</script>