<!-- src/components/FilterPanel.vue -->
<template>
    <div class="card bg-base-100 shadow-sm">
        <div class="card-body p-4 space-y-4">
            <!-- Search -->
            <div class="form-control">
                <label class="label"><span class="label-text">Search</span></label>
                <input v-model="searchQuery" type="text" placeholder="Description..."
                    class="input input-bordered input-sm w-full" />
            </div>

            <!-- Category -->
            <div class="form-control">
                <label class="label"><span class="label-text">Category</span></label>
                <select v-model="filters.category_id" class="select select-bordered select-sm w-full">
                    <option :value="undefined">All Categories</option>
                    <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                        {{ cat.name }}
                    </option>
                </select>
            </div>

            <!-- Amount Range -->
            <div class="form-control">
                <label class="label"><span class="label-text">Amount Range</span></label>
                <div class="flex gap-2">
                    <input v-model.number="filters.min_amount" type="number" placeholder="Min"
                        class="input input-bordered input-sm w-full" min="0" step="0.01" />
                    <input v-model.number="filters.max_amount" type="number" placeholder="Max"
                        class="input input-bordered input-sm w-full" min="0" step="0.01" />
                </div>
            </div>

            <!-- Date Range -->
            <div class="form-control">
                <label class="label"><span class="label-text">Date Range</span></label>
                <input v-model="filters.start_date" type="date" class="input input-bordered input-sm w-full mb-2" />
                <input v-model="filters.end_date" type="date" class="input input-bordered input-sm w-full" />
            </div>

            <!-- Quick Presets -->
            <div class="form-control">
                <label class="label"><span class="label-text">Quick Preset</span></label>
                <select v-model="datePreset" class="select select-bordered select-sm w-full" @change="applyDatePreset">
                    <option value="">Custom</option>
                    <option value="today">Today</option>
                    <option value="week">This week</option>
                    <option value="month">This month</option>
                    <option value="year">This year</option>
                    <option value="last_month">Last month</option>
                </select>
            </div>

            <!-- Sort -->
            <div class="form-control">
                <label class="label"><span class="label-text">Sort By</span></label>
                <select v-model="filters.sort_by" class="select select-bordered select-sm w-full">
                    <option value="date_desc">Newest first</option>
                    <option value="date_asc">Oldest first</option>
                    <option value="amount_desc">Highest amount</option>
                    <option value="amount_asc">Lowest amount</option>
                    <option value="description_asc">Description (A–Z)</option>
                </select>
            </div>

            <!-- Actions -->
            <div class="divider my-0"></div>
            <div class="flex items-center justify-between">
                <button class="btn btn-ghost btn-sm" @click="clearFilters">
                    Clear filters
                </button>
                <span v-if="activeFilterCount > 0" class="badge badge-primary">
                    {{ activeFilterCount }} active
                </span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { Category } from '../types';

const props = defineProps<{
    categories: Category[];
}>();

const searchQuery = defineModel<string>('searchQuery', { default: '' });
const datePreset = defineModel<string>('datePreset', { default: '' });
const filters = defineModel<any>('filters', {
    default: () => ({ sort_by: 'date_desc' })
});
const activeFilterCount = defineModel<number>('activeFilterCount', { default: 0 });

const emit = defineEmits<{
    (e: 'clear'): void;
    (e: 'applyPreset'): void;
}>();

const clearFilters = () => emit('clear');
const applyDatePreset = () => emit('applyPreset');
</script>