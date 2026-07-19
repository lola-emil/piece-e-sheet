<template>
    <dialog :id="modalId" class="modal">
        <div class="modal-box">
            <h3 class="font-bold text-lg">{{ isEdit ? 'Edit Category' : 'Add Category' }}</h3>

            <form @submit.prevent="handleSubmit" class="space-y-4 mt-4">
                <!-- Name Input -->
                <div class="form-control">
                    <label class="label"><span class="label-text">Category Name</span></label>
                    <input v-model="name" type="text" placeholder="e.g., Groceries, Rent, Salary"
                        class="input input-bordered w-full" required />
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
import type { Category } from '../types';

const props = defineProps<{
    modalId: string;
    category: Category | null;
    isSaving: boolean;
}>();

const emit = defineEmits<{
    (e: 'save', name: string): void;
    (e: 'close'): void;
}>();

const isEdit = computed(() => !!props.category);
const name = ref('');

// Reset form when modal opens or category changes
watch(() => props.category, (newCat) => {
    name.value = newCat ? newCat.name : '';
}, { immediate: true });

const handleSubmit = () => {
    emit('save', name.value);
};

const closeModal = () => {
    const modal = document.getElementById(props.modalId) as HTMLDialogElement;
    if (modal) modal.close();
    emit('close');
};
</script>