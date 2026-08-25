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


// Ortama göre .env(.local) içindeki VITE_API_BASE_URL belirler.
// (Buradaki yoruma alınmış sabit adresler silindi: biri artık kullanılmayan
// eski sunucuyu gösteriyordu ve yanlış yere yönlendirme riski taşıyordu.)
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL

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
