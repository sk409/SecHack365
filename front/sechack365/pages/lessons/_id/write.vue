<template>
  <div class="w-100 h-100">
    <mavon-editor v-model="book" language="ja" class="w-100 h-100" />
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import mavonEditor from "mavon-editor";
import "mavon-editor/dist/css/index.css";
import Vue from "vue";
import { Url, urlLessons } from "@/assets/js/url.js";
Vue.use(mavonEditor);
export default {
  data() {
    return {
      book: "",
      lesson: null
    };
  },
  created() {
    const lessonID = this.$route.params.id;
    const data = {
      id: lessonID
    };
    const url = new Url(urlLessons);
    ajax.get(url.base, data).then(response => {
      this.lesson = response.data[0];
      this.book = this.lesson.Book;
    });
  }
};
</script>
