import { Command } from 'commander';
import inquirer from 'inquirer';
import api from '../api';
import chalk from 'chalk';

export const addCommand = new Command('add')
    .description('Add a new expense')
    .action(async () => {
        try {
            const answers = await inquirer.prompt([
                { type: 'input', name: 'description', message: 'Description:' },
                { type: 'number', name: 'amount', message: 'Amount:' },
                { type: 'input', name: 'category_id', message: 'Category ID (UUID):', default: '' },
                { type: 'input', name: 'occurred_at', message: 'Date (YYYY-MM-DD):', default: new Date().toISOString().split('T')[0] }
            ]);

            // Format date for API
            const occurred_at = answers.occurred_at.includes('T') 
                ? answers.occurred_at 
                : `${answers.occurred_at}T12:00:00Z`;

            const payload = {
                description: answers.description,
                amount: answers.amount,
                occurred_at,
                category_id: answers.category_id || null
            };

            const res = await api.post('/expenses', payload);
            console.log(chalk.green('Expense added successfully!'));
            console.log(res.data);
        } catch (error: any) {
            console.error(chalk.red('Error adding expense:'), error.response?.data || error.message);
        }
    });