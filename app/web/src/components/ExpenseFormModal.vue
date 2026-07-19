<template>
    <dialog :id="modalId" class="modal">
        <div class="modal-box">
            <h3 class="font-bold text-lg">{{ isEdit ? 'Edit Expense' : 'Add Expense' }}</h3>

            <form @submit.prevent="handleSubmit" class="space-y-4 mt-4">
                <!-- Description -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Description</span></label>
                    <input v-model="form.description" type="text" class="input input-bordered w-full" required />
                </div>

                <!-- Amount -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Amount</span></label>
                    <input v-model.number="form.amount" type="number" step="0.01" class="input input-bordered w-full"
                        required />
                </div>

                <!-- Category -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Category</span></label>
                    <select v-model="form.category_id" class="select select-bordered w-full">
                        <option :value="null">Uncategorized</option>
                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                            {{ cat.name }}
                        </option>
                    </select>
                </div>

                <!-- Date -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Date</span></label>
                    <input v-model="form.date" type="date" class="input input-bordered w-full" required />
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
import type { Expense, Category, CreateExpenseRequest } from '../types';

const props = defineProps<{
    modalId: string;
    expense: Expense | null;
    categories: Category[];
    isSaving: boolean;
}>();

const emit = defineEmits<{
    (e: 'save', payload: CreateExpenseRequest): void;
    (e: 'close'): void;
}>();

const isEdit = computed(() => !!props.expense);

const form = ref({
    description: '',
    amount: 0,
    category_id: null as string | null,
    date: new Date().toISOString().split('T')[0]
});

// Reset form when modal opens or expense changes
watch(() => props.expense, (newExp) => {
    if (newExp) {
        form.value = {
            description: newExp.description,
            amount: newExp.amount,
            category_id: newExp.category_id,
            date: new Date(newExp.occurred_at).toISOString().split('T')[0]
        };
    } else {
        form.value = {
            description: '',
            amount: 0,
            category_id: null,
            date: new Date().toISOString().split('T')[0]
        };
    }
}, { immediate: true });

const handleSubmit = () => {
    // Format date to ISO string for API
    const payload: CreateExpenseRequest = {
        description: form.value.description,
        amount: form.value.amount,
        category_id: form.value.category_id,
        occurred_at: `${form.value.date}T12:00:00Z`
    };
    emit('save', payload);
};

const closeModal = () => {
    const modal = document.getElementById(props.modalId) as HTMLDialogElement;
    if (modal) modal.close();
    emit('close');
};
</script>