import { ref } from 'vue';
import api from '../services/api';
import type { Category } from '../types';

export function useCategories() {
    const categories = ref<Category[]>([]);
    const isLoading = ref(false);
    const isSaving = ref(false);

    const fetchCategories = async () => {
        isLoading.value = true;
        try {
            const { data } = await api.get('/api/categories');
            categories.value = data.data;
        } catch (error) {
            console.error('Failed to fetch categories', error);
        } finally {
            isLoading.value = false;
        }
    };

    const saveCategory = async (name: string, id?: string) => {
        isSaving.value = true;
        try {
            if (id) {
                await api.put(`/api/categories/${id}`, { name });
            } else {
                await api.post('/api/categories', { name });
            }
            await fetchCategories(); // Refresh list
            return true;
        } catch (error) {
            console.error('Failed to save category', error);
            return false;
        } finally {
            isSaving.value = false;
        }
    };

    const deleteCategory = async (id: string) => {
        if (!confirm('Are you sure? This will not delete expenses, but they will become uncategorized.')) return;

        try {
            await api.delete(`/api/categories/${id}`);
            await fetchCategories();
        } catch (error) {
            console.error('Failed to delete category', error);
        }
    };

    return {
        categories,
        isLoading,
        isSaving,
        fetchCategories,
        saveCategory,
        deleteCategory
    };
}