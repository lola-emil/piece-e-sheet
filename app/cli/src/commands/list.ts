import { Command } from 'commander';
import api from '../api';
import Table from 'cli-table3';
import chalk from 'chalk';

export const listCommand = new Command('list')
    .description('List all expenses')
    .action(async () => {
        try {
            const res = await api.get('/expenses');
            const expenses = res.data.data;

            const table = new Table({
                head: [chalk.cyan('ID'), chalk.cyan('Date'), chalk.cyan('Description'), chalk.cyan('Amount')],
                colWidths: [10, 12, 30, 10]
            });

            expenses.forEach((exp: any) => {
                table.push([
                    exp.id.substring(0, 8) + '...',
                    new Date(exp.occurred_at).toLocaleDateString(),
                    exp.description,
                    `${exp.amount}`
                ]);
            });

            console.log(table.toString());
        } catch (error: any) {
            console.error(chalk.red('Error fetching expenses:'), error.message);
        }
    });