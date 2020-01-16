<template>
  <nav>
    <v-app-bar app dark color="primary">
      <v-btn text icon @click="drawer = !drawer">
        <v-icon>mdi-menu</v-icon>
      </v-btn>
      <v-spacer></v-spacer>
      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <!-- <v-avatar v-on="on">
            <img src="@/static/v.png" alt="" />
          </v-avatar>-->
          <v-btn v-on="on" icon>
            <v-icon>mdi-account</v-icon>
          </v-btn>
        </template>
        <v-list-item-group class="white">
          <v-list-item
            v-for="accountLink in accountLinks"
            :key="accountLink.title"
            :to="accountLink.route"
          >
            <v-list-item-icon>
              <v-icon>{{ accountLink.icon }}</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>{{ accountLink.title }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-menu>
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
  props: {
    user: {
      required: true,
      validator: user => typeof user === "object" || user === null
    }
  },
  data() {
    return {
      drawer: false,
      navLinks: [
        {
          title: "取得した教材",
          icon: "mdi-book-outline",
          route: "/materials/downloaded"
        },
        {
          title: "作成した教材",
          icon: "mdi-book-open-page-variant",
          route: "/materials/created"
        },
        {
          title: "作成したレッスン",
          icon: "mdi-pencil",
          route: "/lessons"
        }
      ]
    };
  },
  computed: {
    accountLinks() {
      if (this.user) {
        return [
          {
            title: "ダッシュボード",
            icon: "mdi-view-dashboard",
            route: this.$routes.dashboard.materials.downloaded
          },
          {
            title: "ログアウト",
            icon: "mdi-logout",
            route: this.$routes.logout.base
          }
        ];
      } else {
        return [
          {
            title: "ログイン",
            icon: "mdi-login",
            route: this.$routes.login.base
          },
          {
            title: "登録",
            icon: "mdi-account-plus",
            route: this.$routes.register.base
          }
        ];
      }
    }
  }
};
</script>
