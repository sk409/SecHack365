<template>
  <v-app class="grey lighten-4">
    <Navbar></Navbar>
    <v-content class="white">
      <div class="h-100">
        <v-container fluid class="h-100">
          <v-row class="h-100">
            <v-col md="3" class="d-sm-none d-md-block separator">
              <v-list>
                <v-list-item-group>
                  <v-list-item
                    color="secondary"
                    v-for="sidebarLink in sidebarLinks"
                    :key="sidebarLink.title"
                    :to="sidebarLink.route"
                  >
                    <v-list-item-icon>
                      <v-icon>{{ sidebarLink.icon }}</v-icon>
                    </v-list-item-icon>
                    <v-list-item-content>
                      <v-list-item-title>
                        {{ sidebarLink.title }}
                      </v-list-item-title>
                    </v-list-item-content>
                  </v-list-item>
                </v-list-item-group>
              </v-list>
            </v-col>
            <v-col sm="12" md="9">
              <nuxt></nuxt>
            </v-col>
          </v-row>
        </v-container>
      </div>
    </v-content>
  </v-app>
</template>

<script>
import Navbar from "@/components/Navbar.vue";
export default {
  components: {
    Navbar
  },
  data() {
    return {
      sidebarLinks: [
        {
          title: "取得した教材",
          icon: "mdi-book-outline",
          route: this.$routes.dashboard.materials.downloaded
        },
        {
          title: "作成した教材",
          icon: "mdi-book-open-page-variant",
          route: this.$routes.dashboard.materials.created
        },
        {
          title: "作成したレッスン",
          icon: "mdi-pencil",
          route: this.$routes.dashboard.lessons
        }
      ]
    };
  }
};
</script>

<style>
.separator {
  border-right: 1px solid rgba(0, 0, 0, 0.4);
}
</style>
