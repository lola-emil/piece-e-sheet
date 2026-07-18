export interface Category {
    id: string;
    user_id: string;
    name: string;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
    revision: number;
}

export interface CreateCategoryRequest {
    name: string;
}

export interface UpdateCategoryRequest {
    name: string;
}