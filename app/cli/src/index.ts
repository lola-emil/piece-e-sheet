import { Command } from 'commander';
import { addCommand } from './commands/add';
import { listCommand } from './commands/list';
import { totalCommand } from './commands/total';
import { summaryCommand } from './commands/summary';
import api from './api';
import { saveToken } from './config';
import chalk from 'chalk';

const program = new Command();

program
    .command('login')
    .description('Login to the expense tracker')
    .argument('<email>', 'Your email')
    .argument('<password>', 'Your password')
    .action(async (email, password) => {
        try {
            const res = await api.post('http://localhost:8080/auth/login', { email, password });
            saveToken(res.data.token);
            console.log(chalk.green('Logged in successfully!'));
        } catch (error: any) {
            console.error(chalk.red('Login failed:'), error.response?.data?.error || error.message);
        }
    });

program.addCommand(addCommand);
program.addCommand(listCommand);
program.addCommand(totalCommand);
program.addCommand(summaryCommand);

program.parse(process.argv);