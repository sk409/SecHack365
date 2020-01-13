<template>
  <div class="w-100 h-100">
    <mavon-editor v-model="book" language="ja" class="w-100 h-100 editor" @change="change" />
  </div>
</template>

<script>
import _ from "lodash";
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
      lesson: null,
      delayedUpdate: _.debounce(this.update, 1000),
      updateQueue: Promise.resolve()
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
  },
  methods: {
    change() {
      this.lesson.Book = this.book;
      this.updateQueue.then(() => {
        this.delayedUpdate();
      });
    },
    update() {
      const url = new Url(urlLessons);
      const data = {
        book: this.lesson.Book
      };
      console.log(url.update(this.lesson.ID));
      ajax.put(url.update(this.lesson.ID), data);
    }
  }
};
</script>

<style scoped>
.editor {
  z-index: 0;
}
</style>
