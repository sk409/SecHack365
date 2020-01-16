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
    <v-divider></v-divider>
    <v-btn v-if="$user" color="accent" fab fixed right bottom @click="showAuthForm">
      <v-icon>mdi-account-edit</v-icon>
    </v-btn>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
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
  methods: {
    showAuthForm() {}
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