<script lang="ts" setup>

</script>

<template>
  <div class="background-wrapper">
    <div class="background-fixed"></div>
    <v-app class="app-container">
      <v-main>
        <router-view />
      </v-main>
    </v-app>
  </div>
</template>
<style>
.background-wrapper {
  position: relative;
  min-height: 100vh;
}

.background-fixed {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  background: linear-gradient(to left, rgba(0, 0, 0, 0.7), rgba(0, 0, 0, 0.4)), url('/cami.jpeg') no-repeat center center;
  background-size: cover;
  z-index: -1;
}

.app-container {
  position: relative;
  min-height: 100vh;
  background: transparent !important;
  z-index: 1;
}

/* Make v-main transparent */
.app-container :deep(.v-main) {
  background: transparent !important;
}

/* iOS Safari fix for fixed background */
@supports (-webkit-touch-callout: none) {
  .background-fixed {
    position: fixed;
    height: 100vh;
    width: 100vw;
  }
}

/* Prevent background scrolling on mobile */
@media (max-width: 768px) {
  .background-fixed {
    position: fixed;
    height: 100vh;
    background-attachment: scroll;
  }
}
</style>
