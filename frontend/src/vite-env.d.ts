/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  // import Vue from "@vitejs/plugin-vue";
  const component: DefineComponent<{}, {}, any>
  // Vue.config.devtools = true;
  export default component
}
