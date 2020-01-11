<template>
  <v-app class="grey lighten-4">
    <Navbar></Navbar>
    <v-content class="white">
      <nuxt></nuxt>
    </v-content>
  </v-app>
</template>

<script>
import Navbar from "@/components/Navbar.vue";
export default {
  components: {
    Navbar
  }
};
</script>
