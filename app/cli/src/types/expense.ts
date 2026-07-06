export enum ExpenseProp {
  TRANS_NO,
  DESCRIPTION,
  TRANS_DATE,
  AMOUNT,
}

export interface Expense {
  transactionId: string;
  description: string;
  transactionDate: string;
  amount: number;
  category: string;
}
