export interface User {
    id: string;
    email: string;
    display_name: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface AuthResponse {
    user: User;
    token: string;
}


export interface Expense {
  id: string;
  user_id: string;
  category_id: string | null;
  description: string;
  amount: number;
  occurred_at: string;
  created_at: string;
  updated_at: string;
  revision: number;
}

export interface ExpenseFilter {
  category_id?: string;
  start_date?: string;
  end_date?: string;
}

export interface CreateExpenseRequest {
  category_id: string | null;
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

export interface ExpenseFilter {
    category_id?: string;
    start_date?: string;
    end_date?: string;
}


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