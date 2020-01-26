<template>
  <div>
    <div class="pa-4">教材作成</div>
    <v-divider></v-divider>
    <v-container class="mt-4">
      <v-card>
        <v-card-text>
          <v-form ref="form">
            <v-file-input v-model="thumbnail" label="サムネイル画像"></v-file-input>
            <v-text-field v-model="title" :rules="titleRules" label="タイトル"></v-text-field>
            <v-text-field v-model="description" :rules="descriptionRules" label="概要"></v-text-field>
          </v-form>
          <div class="mt-5 subtitle-1 text--grey">選択したレッスン</div>
          <v-divider class="mb-3"></v-divider>
          <v-card
            v-for="lesson in lessonsSelected"
            :key="lesson.ID"
            width="512"
            class="mb-3 mx-auto"
          >
            <v-container>
              <v-row>
                <v-col cols="3">
                  <v-img v-if="lesson" :src="$serverUrl(lesson.ThumbnailPath)" height="100" contain></v-img>
                </v-col>
                <v-col>
                  <v-card-title>{{lesson.Title | truncate(20)}}</v-card-title>
                  <v-card-text>{{lesson.Description | truncate(30)}}</v-card-text>
                </v-col>
              </v-row>
            </v-container>
          </v-card>
          <div class="text-center mt-5 mb-5">
            <v-btn color="primary" @click="lessonDialog = true">レッスンを追加</v-btn>
          </div>
          <div class="text-center">
            <v-btn color="accent" :loading="loading" @click="submit">作成</v-btn>
          </div>
        </v-card-text>
      </v-card>
    </v-container>
    <v-dialog v-model="lessonDialog" max-width="256">
      <v-card>
        <v-list subheader>
          <v-subheader>レッスン一覧</v-subheader>
          <v-list-item v-for="lesson in lessons" :key="lesson.ID">
            <v-list-item-avatar size="64" class="grey lighten-2">
              <v-img v-if="lesson" :src="$serverUrl(lesson.ThumbnailPath)" contain></v-img>
            </v-list-item-avatar>
            <v-list-item-content class="title">{{lesson.Title}}</v-list-item-content>
            <v-list-item-action>
              <v-checkbox :input-value="lesson.selected" @change="clickLessonCheckbox(lesson)"></v-checkbox>
            </v-list-item-action>
          </v-list-item>
        </v-list>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlLessons, urlMaterials } from "@/assets/js/url.js";
export default {
  data() {
    return {
      description: "",
      descriptionRules: [
        v => !!v || "概要を入力してください",
        v => v.length <= 1024 || "1024文字以内で入力してください"
      ],
      lessonDialog: false,
      lessons: [],
      lessonsSelected: [],
      loading: false,
      thumbnail: null,
      title: "",
      titleRules: [
        v => !!v || "タイトルを入力してください",
        v => v.length <= 128 || "128文字以内で入力してください"
      ]
    };
  },
  created() {
    this.$fetchUser(() => {
      const url = new Url(urlLessons);
      const data = {
        user_id: this.$user.ID,
        downloaded: 0
      };
      ajax.get(url.base, data).then(response => {
        response.data.forEach(lesson => (lesson.selected = false));
        this.lessons = response.data;
      });
    });
  },
  methods: {
    clickLessonCheckbox(lesson) {
      lesson.selected = !lesson.selected;
      if (lesson.selected) {
        this.lessonsSelected.push(lesson);
      } else {
        this.lessonsSelected = this.lessonsSelected.filter(l => l !== lesson);
      }
    },
    submit() {
      if (!this.$refs.form.validate()) {
        return;
      }
      const data = {
        thumbnail: this.thumbnail,
        title: this.title,
        description: this.description,
        "lessonIDs[]": this.lessonsSelected.map(lesson => lesson.ID),
        userID: this.$user.ID
      };
      const config = {
        headers: {
          "content-type": "multipart/form-data"
        }
      };
      const url = new Url(urlMaterials);
      this.loading = true;
      ajax.post(url.base, data, config).then(response => {
        this.loading = false;
        if (response.status !== 200) {
          return;
        }
        this.$router.push(this.$routes.dashboard.materials.created);
      });
    }
  }
};
</script>

<style>
.lesson-card-thumbnail {
  object-fit: contain;
  width: 100%;
  height: 128px;
}
</style>