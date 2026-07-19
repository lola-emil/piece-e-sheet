import "@/assets/style.css"

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import { useAuthStore } from "./stores/auth.ts"

const app = createApp(App)

app.use(createPinia())
app.use(router)

const store = useAuthStore();
store.initAuth();

app.mount('#app')
