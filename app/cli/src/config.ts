import fs from 'fs';
import path from 'path';
import os from 'os';

const CONFIG_DIR = path.join(os.homedir(), '.tota-cli');
const CONFIG_FILE = path.join(CONFIG_DIR, 'config.json');

export interface Config {
    token?: string;
    apiUrl: string;
}

export const loadConfig = (): Config => {
    if (!fs.existsSync(CONFIG_DIR)) fs.mkdirSync(CONFIG_DIR);
    if (!fs.existsSync(CONFIG_FILE)) {
        const defaultConfig: Config = { apiUrl: 'http://localhost:8080/api' };
        fs.writeFileSync(CONFIG_FILE, JSON.stringify(defaultConfig));
        return defaultConfig;
    }
    return JSON.parse(fs.readFileSync(CONFIG_FILE, 'utf-8'));
};

export const saveToken = (token: string) => {
    const config = loadConfig();
    config.token = token;
    fs.writeFileSync(CONFIG_FILE, JSON.stringify(config));
};