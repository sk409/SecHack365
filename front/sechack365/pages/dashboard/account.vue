<template>
  <v-container fluid>
    <div class="d-flex">
      <div class="profile-image-container">
        <img :src="profileImagePath" class="profile-image" />
      </div>
      <div>
        <div>
          <span class="headline">{{$user ? $user.Name : ""}}</span>
        </div>
        <div class="mt-3 body">バイオテキスト</div>
      </div>
    </div>
    <v-divider class="my-3"></v-divider>
    <div class="text-center">
      <v-btn color="secondary" @click="logout">ログアウト</v-btn>
    </div>
    <v-btn v-if="$user" color="accent" fab fixed right bottom>
      <v-icon>mdi-account-edit</v-icon>
    </v-btn>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlLogout } from "@/assets/js/url.js";
export default {
  layout: "dashboard",
  computed: {
    profileImagePath() {
      return this.$user
        ? process.env.serverOrigin + this.$user.ProfileImagePath
        : "";
    }
  },
  created() {
    this.$fetchUser();
  },
  mounted() {
    this.$nuxt.$emit("setTitle", "アカウント情報");
  },
  methods: {
    logout() {
      const url = new Url(urlLogout);
      ajax.post(url.base, {}, { withCredentials: true }).then(response => {
        if (response.status === 200) {
          this.$router.push("/");
        }
      });
    }
  }
};
</script>

<style>
.profile-image-container {
  width: 200px;
}
.profile-image {
  width: 100%;
  object-fit: contain;
}
</style>