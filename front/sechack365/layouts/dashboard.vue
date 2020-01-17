<template>
  <v-app class="grey lighten-4">
    <Navbar :user="$user"></Navbar>
    <v-content class="white">
      <div class="h-100">
        <div fluid class="h-100">
          <div class="d-flex h-100">
            <div class="d-none d-md-inline separator">
              <v-list>
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
                    <v-list-item-title>{{ sidebarLink.title }}</v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list>
            </div>
            <div class="w-100">
              <div class="d-flex justify-space-between align-center">
                <div class="pa-4 title">{{title}}</div>
                <div class="d-inline d-md-none mr-3">
                  <v-menu offset-y>
                    <template v-slot:activator="{on}">
                      <v-btn icon large v-on="on">
                        <v-icon>mdi-dots-horizontal-circle-outline</v-icon>
                      </v-btn>
                    </template>
                    <v-list>
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
                          <v-list-item-title>{{ sidebarLink.title }}</v-list-item-title>
                        </v-list-item-content>
                      </v-list-item>
                    </v-list>
                  </v-menu>
                </div>
              </div>
              <v-divider></v-divider>
              <nuxt></nuxt>
            </div>
          </div>
        </div>
      </div>
    </v-content>
  </v-app>
</template>

<script>
import Navbar from "@/components/Navbar.vue";
export default {
  middleware: "auth",
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
          icon: "mdi-monitor-clean",
          route: this.$routes.dashboard.lessons
        },
        {
          title: "フォロー",
          icon: "mdi-account-check",
          route: this.$routes.dashboard.follow
        },
        {
          title: "フォロワー",
          icon: "mdi-account-star",
          route: this.$routes.dashboard.follower
        },
        {
          title: "アカウント情報",
          icon: "mdi-account-badge-horizontal",
          route: this.$routes.dashboard.account
        }
      ],
      title: ""
    };
  },
  created() {
    this.$fetchUser();
    this.$nuxt.$on("setTitle", this.setTitle);
  },
  methods: {
    setTitle(title) {
      this.title = title || "";
    }
  }
};
</script>

<style>
.separator {
  border-right: 1px solid rgba(0, 0, 0, 0.4);
}
</style>
