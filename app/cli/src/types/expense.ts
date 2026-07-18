export interface Expense {
    id: string;
    user_id: string;
    category_id: string | null;
    description: string;
    amount: number;
    occurred_at: string;
    created_at: string;
    updated_at: string;
    deleted_at: string | null;
    revision: number;
}

export interface CreateExpenseRequest {
    category_id?: string | null;
    description: string;
    amount: number;
    occurred_at: string;
}

export interface UpdateExpenseRequest {
    category_id?: string | null;
    description: string;
    amount: number;
    occurred_at: string;
}

// Used for query parameters in GET /expenses
export interface ExpenseFilter {
    category_id?: string;
    start_date?: string;
    end_date?: string;
}