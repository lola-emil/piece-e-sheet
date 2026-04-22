# Piece-E-Sheet 🧾

**PizZ-O-Sheet** is a simple command-line tool to track personal expenses using Google Sheets. It allows you to quickly add and retrieve records via the terminal.

> Built with Node.js, Google Sheets API, Commander.js, and Chalk.

---

## 🚀 Features

- Add expense records with a UUID, timestamp, description, and amount
- Retrieve (peek) all records from the spreadsheet
- Simple and friendly CLI interface
- All data stored in Google Sheets

---

## 🛠️ Setup

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/pizz-o-sheet.git
cd pizz-o-sheet
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Configure Environment Variables
```.env
GOOGLE_SHEET_ID=your_google_sheet_id_here
GOOGLE_APPLICATION_CREDENTIALS=./service-key.json
```

### 4. Google API Setup

1. Go to the Google Cloud Console

2. Create a project

3. Enable the Google Sheets API

4. Create credentials (OAuth Client or Service Account)

5. Download the credentials and set them up properly (Use the service account; download the service key as JSON)

## 📦 Usage

### 1. Install the as GLOBAL
```bash
npm install -g
```

### 2. --help for the instructions
```bash
psheet --help
```
