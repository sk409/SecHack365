<template>
  <v-app class="grey lighten-4">
    <NavbarIde v-if="isIde" :lesson="lesson"></NavbarIde>
    <Navbar v-else></Navbar>
    <v-content class="white">
      <div v-if="isLoading">LOADING</div>
      <div v-else-if="isEditing" class="h-100">
        <nuxt></nuxt>
      </div>
      <div v-else-if="isPermissionError">PERMISSION ERROR</div>
    </v-content>
  </v-app>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import Navbar from "@/components/Navbar.vue";
import NavbarIde from "@/components/NavbarIde.vue";
import { Url, urlLessons } from "@/assets/js/url.js";
const statusLoading = "statusLoading";
const statusEditing = "statusEditing";
const statusPermissionError = "statusPermissionError";
export default {
  middleware: "auth",
  components: {
    Navbar,
    NavbarIde
  },
  data() {
    return {
      status: statusLoading,
      lesson: null
    };
  },
  computed: {
    isIde() {
      return this.$route.path.includes("ide");
    },
    isLoading() {
      return this.status === statusLoading;
    },
    isEditing() {
      return this.status === statusEditing;
    },
    isPermissionError() {
      return this.status === statusPermissionError;
    },
    thumbnailURL() {
      if (!this.lesson) {
        return "";
      }
      return process.env.serverOrigin + this.lesson.ThumbnailPath;
    }
  },
  created() {
    const lessonID = this.$route.params.id;
    const url = new Url(urlLessons);
    const data = {
      user_id: this.$store.state.users.user.ID
    };
    ajax.get(url.base, data).then(response => {
      const lessons = response.data;
      const index = lessons.findIndex(lesson => lesson.ID == lessonID);
      const notFound = -1;
      this.status = index === notFound ? statusPermissionError : statusEditing;
      if (this.isEditing) {
        this.lesson = lessons[index];
      }
    });
  }
};
</script>
