/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Plugins
import vuetify from './vuetify'
import pinia from '../stores'
import router from '../router'
import axios from 'axios';
import VueAxios from "vue-axios";


// axios.defaults.baseURL = "http://127.0.0.1:8080/";

axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL

// axios.defaults.baseURL = "http://89.117.37.208:8080/";
// axios.defaults.baseURL = "/api/";

axios.defaults.withCredentials = true;

// Types
import type { App } from 'vue'

export function registerPlugins (app: App) {
  app
    .use(vuetify)
    .use(router)
    .use(pinia)
    .use(VueAxios, axios)
}
