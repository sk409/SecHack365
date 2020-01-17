<template>
  <div>
    <UserList :users="users"></UserList>
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
      users: []
    };
  },
  created() {
    this.$fetchUser(() => {
      const url = new Url(urlUsers);
      ajax.get(url.follower(this.$user.ID)).then(response => {
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
    this.$nuxt.$emit("setTitle", "フォローされているユーザ一覧");
  }
};
</script>