#! /usr/bin/env node

import { config } from "dotenv";
import { program } from "commander";
import chalk from "chalk";
import { google } from "googleapis";

import { v4 } from "uuid";
import { ExpenseProp } from "./types/expense";


config({ quiet: true });

program
  .name("Piece 'e Sheet")
  .description("CLI Google Sheet personal expense tracker")
  .version("1.0.1");

program
  .command("peek")
  .description("Retrieves all the records from the spreadsheet")
  .action(async () => {
    const googleAuth = await google.auth.getClient({
      scopes: ["https://www.googleapis.com/auth/spreadsheets"],
    });

    const sheets = google.sheets({ version: "v4", auth: googleAuth });

    const response = await sheets.spreadsheets.values.get({
      spreadsheetId: process.env["GOOGLE_SHEET_ID"],
      range: "Sheet1",
    });

    if (response.status == 200) {
      console.table(response.data.values);
    }
  });

program
  .command("add")
  .description("Add expense record to the spreadsheet")
  .option("-a, --amount <number>", "Expense amount")
  .option("-d, --description <string>", "Expense description", "Unspecified")
  .action(async (options) => {
    const opt: {
      amount?: number;
      description?: string;
    } = options;

    if (!opt.amount) {
      console.log(
        `${chalk.redBright("Error")}: option \"-a, --amount\" must be specified`
      );
      process.exit(1);
    }

    const googleAuth = await google.auth.getClient({
      scopes: ["https://www.googleapis.com/auth/spreadsheets"],
    });

    const sheets = google.sheets({ version: "v4", auth: googleAuth });

    let newExpense: any[] = [];

    let today = new Date();

    newExpense[ExpenseProp.TRANS_NO] = v4();
    newExpense[ExpenseProp.DESCRIPTION] = opt.description;
    newExpense[ExpenseProp.AMOUNT] = opt.amount;
    newExpense[ExpenseProp.TRANS_DATE] = today.toLocaleString();

    const response = await sheets.spreadsheets.values.append({
      spreadsheetId: process.env["GOOGLE_SHEET_ID"],
      range: "Sheet1!A1",
      valueInputOption: "USER_ENTERED",
      insertDataOption: "INSERT_ROWS",
      requestBody: {
        values: [newExpense],
      },
    });

    if (response.status == 200) console.log(`${chalk.greenBright("Ok")}`);
  });

program.parse(process.argv);
