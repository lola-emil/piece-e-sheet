import { Command } from 'commander';
import api from '../api';
import Table from 'cli-table3';
import chalk from 'chalk';
import { Category } from '../types/category';

export const summaryCommand = new Command('summary')
    .description('Summary of expenses by category')
    .action(async () => {
        try {
            // Fetch expenses and categories to map IDs to Names
            const [expRes, catRes] = await Promise.all([
                api.get('/expenses'),
                api.get('/categories')
            ]);

            const expenses = expRes.data.data;
            const categories = catRes.data.data;
            const catMap = new Map<number, string>(categories.map((c: Category) => [c.id, c.name]));

            const summary: Record<string, number> = {};

            expenses.forEach((exp: any) => {
                const catName = exp.category_id ? catMap.get(exp.category_id) || 'Unknown' : 'Uncategorized';
                summary[catName] = (summary[catName] || 0) + exp.amount;
            });

            const table = new Table({
                head: [chalk.cyan('Category'), chalk.cyan('Total')],
            });

            for (const [cat, total] of Object.entries(summary)) {
                table.push([cat, `${total.toFixed(2)}`]);
            }

            console.log(table.toString());
        } catch (error: any) {
            console.error(chalk.red('Error generating summary:'), error.message);
        }
    });