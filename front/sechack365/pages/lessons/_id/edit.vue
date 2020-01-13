<template>
  <div>
    <v-container>
      <v-row>
        <v-col cols="3">
          <v-btn
            fixed
            right
            bottom
            color="accent"
            @click="$router.push($routes.lessons.ide(lesson.ID))"
            >IDEで編集</v-btn
          >
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlLessons } from "@/assets/js/url.js";
export default {
  layout: "lesson",
  data() {
    return {
      lesson: null
    };
  },
  created() {
    const url = new Url(urlLessons);
    const lessonID = this.$route.params.id;
    const data = {
      id: lessonID
    };
    ajax.get(url.base, data).then(response => {
      this.lesson = response.data[0];
    });
  }
};
</script>
