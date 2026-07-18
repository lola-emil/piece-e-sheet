import { Command } from 'commander';
import api from '../api';
import chalk from 'chalk';

export const totalCommand = new Command('total')
    .description('Calculate total expenses')
    .action(async () => {
        try {
            const res = await api.get('/expenses');
            const expenses = res.data.data;
            
            const total = expenses.reduce((sum: number, exp: any) => sum + exp.amount, 0);
            
            console.log(chalk.bold(`Total Expenses: ${total.toFixed(2)}`));
        } catch (error: any) {
            console.error(chalk.red('Error calculating total:'), error.message);
        }
    });