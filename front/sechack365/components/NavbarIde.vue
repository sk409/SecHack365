<template>
  <nav>
    <v-app-bar app dark color="primary">
      <v-btn text icon @click="drawer = !drawer">
        <v-icon>mdi-menu</v-icon>
      </v-btn>
      <v-spacer></v-spacer>
    </v-app-bar>
    <v-navigation-drawer v-model="drawer" app class="grey lighten-4">
      <div class="pa-2">
        <img v-if="lesson" :src="$serverUrl(lesson.ThumbnailPath)" class="thumbnail-2" />
      </div>
      <div class="text-center title">{{ lesson ? lesson.Title : "" }}</div>
      <v-divider class="my-2"></v-divider>
      <v-list>
        <v-list-item :to="$routes.lessons.write(lesson ? lesson.ID : '')" target="_blank">
          <v-list-item-icon>
            <v-icon>mdi-pencil-box-multiple-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>説明文</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-group prepend-icon="mdi-wan">
          <template v-slot:activator>
            <v-list-item-title>ポート</v-list-item-title>
          </template>
          <v-list-item
            v-for="port in ports"
            :key="port.ID"
            :href="serverURL(port.HostPort)"
            target="_blank"
            link
            class="ml-5"
          >
            <v-list-item-icon>
              <v-icon>mdi-lan-connect</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>{{ port.Port }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-group>
        <v-list-item link :to="$routes.dashboard.materials.downloaded">
          <v-list-item-icon>
            <v-icon>mdi-view-dashboard</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>ダッシュボード</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
  </nav>
</template>

<script>
export default {
  props: {
    lesson: {
      required: true
    }
  },
  data() {
    return {
      drawer: false,
      navLinks: [
        {
          title: "説明文",
          icon: "mdi-pencil-box-multiple-outline",
          route: "/materials/downloaded"
        },
        {
          title: "ポート",
          icon: "mdi-wan",
          route: "/materials/created"
        },
        {
          title: "ダッシュボード",
          icon: "mdi-view-dashboard",
          route: this.$routes.dashboard.materials.downloaded
        }
      ]
    };
  },
  computed: {
    ports() {
      if (!this.lesson) {
        return [];
      }
      return this.lesson.Ports.filter(
        port => port.Port !== this.lesson.ConsolePort
      );
    }
  },
  methods: {
    serverURL(port) {
      return process.env.serverHost + ":" + port;
    }
  }
};
</script>

<style scoped>
.navigation-drawer {
  z-index: 100;
}
</style>