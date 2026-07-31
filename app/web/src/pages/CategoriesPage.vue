<template>
    <div class="space-y-6">
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
            <h1 class="text-3xl font-bold">Categories</h1>
            <button class="btn btn-primary" @click="openAddModal">
                + Add Category
            </button>
        </div>

        <div class="form-control w-full max-w-xs">
            <input v-model="searchQuery" type="text" placeholder="Search categories..."
                class="input input-bordered w-full" />
        </div>

        <div v-if="isLoading" class="flex justify-center py-12">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>

        <div v-else-if="filteredCategories?.length > 0"
            class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            <div v-for="cat in filteredCategories" :key="cat.id"
                class="card bg-base-100 shadow-md hover:shadow-lg transition-shadow">
                <div class="card-body p-4">
                    <div class="flex justify-between items-start">
                        <h2 class="card-title text-base truncate" :title="cat.name">
                            <Tag :size="16" />
                            {{ cat.name }}
                        </h2>

                        <div class="dropdown dropdown-end">
                            <label tabindex="0" class="btn btn-ghost btn-xs btn-circle">
                                <EllipsisVertical :size="16" />
                            </label>
                            <ul tabindex="0"
                                class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-32 z-10 border border-base-200">
                                <li><a @click="openEditModal(cat)">Edit</a></li>
                                <li><a class="text-error" @click="deleteCategory(cat.id)">Delete</a></li>
                            </ul>
                        </div>
                    </div>

                    <div class="text-xs text-base-content/50 mt-2">
                        Created {{ formatDate(cat.created_at) }}
                    </div>
                </div>
            </div>
        </div>

        <div v-else class="text-center py-12 bg-base-100 rounded-lg shadow-sm flex flex-col items-center gap-5">

            <div class="mx-auto text-base-content/30">
                <Tag :size="48" />
            </div>
            <p class="text-base-content/60">No categories found.</p>
        </div>

        <CategoryFormModal modalId="category_modal" :category="selectedCategory" :is-saving="isSaving"
            @save="handleSave" @close="selectedCategory = null" />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useCategories } from '../composables/useCategories';
import CategoryFormModal from '../components/CategoryFormModal.vue';
import type { Category } from '../types';
import { EllipsisVertical, Tag } from '@lucide/vue';

const { categories, isLoading, isSaving, fetchCategories, saveCategory, deleteCategory } = useCategories();

const selectedCategory = ref<Category | null>(null);
const searchQuery = ref('');

const filteredCategories = computed(() => {
    if (!searchQuery.value) return categories.value;
    return categories.value.filter(cat =>
        cat.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
});

// Modal Handlers
const openAddModal = () => {
    selectedCategory.value = null;
    const modal = document.getElementById('category_modal') as HTMLDialogElement;
    if (modal) modal.showModal();
};

const openEditModal = (cat: Category) => {
    selectedCategory.value = cat;
    const modal = document.getElementById('category_modal') as HTMLDialogElement;
    if (modal) modal.showModal();
};

const handleSave = async (name: string) => {
    const success = await saveCategory(name, selectedCategory.value?.id);
    if (success) {
        const modal = document.getElementById('category_modal') as HTMLDialogElement;
        if (modal) modal.close();
    }
};

// Helpers
const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString();
};

// Load data
onMounted(() => {
    fetchCategories();
});
</script>