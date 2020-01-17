<template>
  <div>
    <UserList :users="users"></UserList>
    <v-dialog v-model="dialog">
      <div class="mx-auto pa-3 white">
        <v-subheader>ユーザを検索</v-subheader>
        <v-text-field
          v-model="searchKeyword"
          prepend-icon="mdi-magnify"
          label="キーワード"
          @input="inputKeyword"
        ></v-text-field>
        <div class="searched-users-container">
          <v-card
            v-for="user in searchedUsers"
            :key="user.ID"
            class="ma-2 searched-user-card"
            @click="$router.push($routes.users.show(user.ID))"
          >
            <v-avatar size="128" class="mx-auto d-block">
              <v-img :src="profileImageURL(user)"></v-img>
            </v-avatar>
            <v-divider></v-divider>
            <v-card-title>{{user.Name}}</v-card-title>
          </v-card>
        </div>
      </div>
    </v-dialog>
    <v-btn color="accent" fab fixed right bottom @click="dialog = true">
      <v-icon>mdi-plus</v-icon>
    </v-btn>
  </div>
</template>

<script>
import _ from "lodash";
import ajax from "@/assets/js/ajax.js";
import UserList from "@/components/UserList.vue";
import { Url, urlUsers } from "@/assets/js/url.js";
export default {
  layout: "dashboard",
  components: {
    UserList
  },
  data() {
    return {
      delayedSearchUsers: _.debounce(this.searchUsers, 1000),
      dialog: false,
      searchedUsers: [],
      searchKeyword: "",
      users: []
    };
  },
  created() {
    this.$fetchUser(() => {
      const url = new Url(urlUsers);
      ajax.get(url.follow(this.$user.ID)).then(response => {
        this.users = response.data.map(user => {
          return {
            id: user.ID,
            name: user.Name,
            biotext: "バイオテキスト",
            avatar: process.env.serverOrigin + "/" + user.ProfileImagePath
          };
        });
      });
    });
  },
  mounted() {
    this.$nuxt.$emit("setTitle", "フォローしているユーザ一覧");
  },
  methods: {
    inputKeyword() {
      this.delayedSearchUsers();
    },
    profileImageURL(user) {
      return process.env.serverOrigin + "/" + user.ProfileImagePath;
    },
    searchUsers() {
      const url = new Url(urlUsers);
      const data = {
        keyword: this.searchKeyword
      };
      ajax.get(url.search, data).then(response => {
        this.searchedUsers = response.data;
      });
    }
  }
};
</script>

<style>
.searched-user-card {
  cursor: pointer;
}
.searched-users-container {
  display: grid;
  grid-template-columns: repeat(5, 20%);
}
</style>