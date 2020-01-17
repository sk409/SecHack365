<template>
  <nav>
    <v-app-bar app dark color="primary">
      <v-btn text icon @click="drawer = !drawer">
        <v-icon>mdi-menu</v-icon>
      </v-btn>
      <v-spacer></v-spacer>
    </v-app-bar>
    <v-navigation-drawer v-model="drawer" app class="grey lighten-4">
      <v-list>
        <v-list-item-group>
          <v-list-item v-for="navLink in navLinks" :key="navLink.title" :to="navLink.route">
            <v-list-item-icon>
              <v-icon>{{ navLink.icon }}</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>{{ navLink.title }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>
  </nav>
</template>

<script>
export default {
  data() {
    return {
      drawer: false,
      navLinks: [
        {
          title: "ホーム",
          icon: "mdi-home",
          route: this.$routes.home
        },
        {
          title: "ダッシュボード",
          icon: "mdi-view-dashboard",
          route: this.$routes.dashboard.materials.downloaded
        }
      ]
    };
  }
};
</script>
