<template>
  <div class="w-100 h-100">
    <mavon-editor v-model="book" language="ja" class="w-100 h-100 editor" @imgAdd="imgAdd" />
  </div>
</template>

<script>
import _ from "lodash";
import ajax from "@/assets/js/ajax.js";
import { Url, urlLessons } from "@/assets/js/url.js";
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
  watch: {
    book(value) {
      this.lesson.Book = this.book;
      this.updateQueue.then(() => {
        this.delayedUpdate();
      });
    }
  },
  methods: {
    imgAdd(pos, file) {
      // var formdata = new FormData();
      // formdata.append("image", $file);
      // axios({
      //   url: "server url",
      //   method: "post",
      //   data: formdata,
      //   headers: { "Content-Type": "multipart/form-data" }
      // }).then(url => {
      //   // step 2. replace url ![...](0) -> ![...](url)
      //   $vm.$img2Url(pos, url);
      // });
      const url = new Url(urlLessons);
      const data = {
        image: file
      };
      const config = {
        headers: {
          "content-type": "multipart/form-data"
        }
      };
      ajax.post(url.book, data, config).then(response => {
        const regex = new RegExp(`!\\[${file.name}\\]\\([0-9]+\\)`);
        const u = process.env.serverOrigin + response.data;
        this.book = this.book.replace(regex, `![${file.name}](${u})`);
      });
    },
    update() {
      const url = new Url(urlLessons);
      const data = {
        book: this.lesson.Book
      };
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
