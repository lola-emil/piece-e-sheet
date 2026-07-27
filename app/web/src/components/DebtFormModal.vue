<template>
    <dialog :id="modalId" class="modal">
        <div class="modal-box">
            <h3 class="font-bold text-lg">{{ isEdit ? 'Edit Debt' : 'Add Debt' }}</h3>

            <form @submit.prevent="handleSubmit" class="space-y-4 mt-4">
                <!-- Person Name -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Person / Entity</span></label>
                    <input v-model="form.person_name" type="text" class="input input-bordered w-full" required />
                </div>

                <!-- Description -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Description</span></label>
                    <input v-model="form.description" type="text" class="input input-bordered w-full" />
                </div>

                <div class="grid grid-cols-2 gap-4">
                    <!-- Amount -->
                    <div class="form-control">
                        <label class="label"><span class="label-text">Amount</span></label>
                        <input v-model.number="form.amount" type="number" step="0.01"
                            class="input input-bordered w-full" required />
                    </div>

                    <!-- Type -->
                    <div class="form-control">
                        <label class="label"><span class="label-text">Type</span></label>
                        <select v-model="form.type" class="select select-bordered w-full">
                            <option value="i_owe">I Owe (Liability)</option>
                            <option value="owed_to_me">Owed to Me (Asset)</option>
                        </select>
                    </div>
                </div>

                <!-- Due Date -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Due Date (Optional)</span></label>
                    <input v-model="form.due_date" type="date" class="input input-bordered w-full" />
                </div>

                <!-- Edit Only: Status & Remaining -->
                <div v-if="isEdit" class="divider my-2"></div>

                <div v-if="isEdit" class="grid grid-cols-2 gap-4">
                    <div class="form-control">
                        <label class="label"><span class="label-text">Status</span></label>
                        <select v-model="form.status" class="select select-bordered w-full">
                            <option value="pending">Pending</option>
                            <option value="partial">Partial</option>
                            <option value="paid">Paid</option>
                        </select>
                    </div>
                    <div class="form-control">
                        <label class="label"><span class="label-text">Remaining Amount</span></label>
                        <input v-model.number="form.remaining_amount" type="number" step="0.01"
                            class="input input-bordered w-full" required />
                    </div>
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
import { ref, watch, computed } from 'vue';
import type { Debt, CreateDebtRequest, UpdateDebtRequest } from '../types';

const props = defineProps<{
    modalId: string;
    debt: Debt | null;
    isSaving: boolean;
}>();

const emit = defineEmits<{
    (e: 'save', payload: CreateDebtRequest | UpdateDebtRequest): void;
    (e: 'close'): void;
}>();

const isEdit = computed(() => !!props.debt);

const form = ref({
    person_name: '',
    description: '',
    amount: 0,
    type: 'i_owe' as 'i_owe' | 'owed_to_me',
    due_date: '',
    status: 'pending' as 'pending' | 'partial' | 'paid',
    remaining_amount: 0
});

watch(() => props.debt, (newDebt) => {
    if (newDebt) {
        form.value = {
            person_name: newDebt.person_name,
            description: newDebt.description || '',
            amount: newDebt.amount,
            type: newDebt.type,
            due_date: newDebt.due_date ? new Date(newDebt.due_date).toISOString().slice(0, 10) : '',
            status: newDebt.status,
            remaining_amount: newDebt.remaining_amount
        };
    } else {
        form.value = {
            person_name: '',
            description: '',
            amount: 0,
            type: 'i_owe',
            due_date: '',
            status: 'pending',
            remaining_amount: 0
        };
    }
}, { immediate: true });

const handleSubmit = () => {
    const payload: any = {
        person_name: form.value.person_name,
        description: form.value.description,
        amount: form.value.amount,
        type: form.value.type,
        due_date: form.value.due_date ? `${form.value.due_date}T12:00:00Z` : null,
    };

    if (isEdit.value) {
        payload.status = form.value.status;
        payload.remaining_amount = form.value.remaining_amount;
    }

    emit('save', payload);
};

const closeModal = () => {
    const modal = document.getElementById(props.modalId) as HTMLDialogElement;
    if (modal) modal.close();
    emit('close');
};
</script>