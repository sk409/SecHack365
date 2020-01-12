<template>
  <div>
    <div v-if="isLoading">LOADING</div>
    <div v-else-if="isEditing">EDITAING</div>
    <div v-else-if="isPermissionError">PERMISSION ERROR</div>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlLessons } from "@/assets/js/url.js";
const statusLoading = "statusLoading";
const statusEditing = "statusEditing";
const statusPermissionError = "statusPermissionError";
export default {
  data() {
    return {
      status: statusLoading,
      lesson: null
    };
  },
  computed: {
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
    });
  }
};
</script>
