<template>
  <v-app class="grey lighten-4">
    <NavbarIde v-if="isIde" :lesson="lesson"></NavbarIde>
    <Navbar v-else :user="$user"></Navbar>
    <v-content class="white">
      <div v-if="isLoading"></div>
      <div v-else-if="isEditing" class="h-100">
        <nuxt></nuxt>
      </div>
      <div v-else-if="isPermissionError">PERMISSION ERROR</div>
    </v-content>
  </v-app>
  <!-- <div class="h-100">
    <div v-if="isLoading">LOADING</div>
    <div v-else-if="isEditing" class="h-100">
      <nuxt></nuxt>
    </div>
    <div v-else-if="isPermissionError">PERMISSION ERROR</div>
  </div>-->
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import Navbar from "@/components/Navbar.vue";
import NavbarIde from "@/components/NavbarIde.vue";
import { Url, urlLessons, urlUser } from "@/assets/js/url.js";
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
      // user: null
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
    }
  },
  created() {
    this.$fetchUser(() => {
      const lessonID = this.$route.params.id;
      const urlL = new Url(urlLessons);
      const data = {
        user_id: this.$user.ID
      };
      ajax.get(urlL.base, data).then(response => {
        const lessons = response.data;
        const index = lessons.findIndex(lesson => lesson.ID == lessonID);
        const notFound = -1;
        this.status =
          index === notFound ? statusPermissionError : statusEditing;
        if (this.isEditing) {
          this.lesson = lessons[index];
        }
      });
    });
  }
};
</script>
