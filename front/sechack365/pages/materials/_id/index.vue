<template>
  <v-container>
    <v-row>
      <v-col cols="3">
        <img
          v-if="material"
          :src="$serverUrl(material.ThumbnailPath)"
          class="thumbnail-2 thumbnail-top"
        />
      </v-col>
      <v-col cols="9">
        <div class="title">{{material ? material.Title : ""}}</div>
        <!-- <div class="body">バイオテキスト</div> -->
      </v-col>
    </v-row>
    <v-divider></v-divider>
    <div class="subtitle-1 my-2">この教材に含まれているレッスン一覧</div>
    <div v-if="material" class="lesson-card-container">
      <v-card v-for="lesson in material.lessons" :key="lesson.ID" class="ma-4">
        <img v-if="lesson" contain :src="$serverUrl(lesson.ThumbnailPath)" class="thumbnail-2" />
        <v-divider></v-divider>
        <v-card-title>{{lesson.Title | truncate(20)}}</v-card-title>
        <v-card-text>{{lesson.Description | truncate(30)}}</v-card-text>
      </v-card>
    </div>
    <v-snackbar v-model="snackbar" :timeout="2000" top>
      <span>{{notification}}</span>
      <v-btn icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
    <v-btn
      v-if="downloaded === false"
      color="accent"
      :loading="loading"
      fab
      fixed
      right
      bottom
      @click="download"
    >
      <v-icon>mdi-download</v-icon>
    </v-btn>
    <v-btn v-else-if="downloaded === true" color="primary" fixed right bottom disabled>取得済み</v-btn>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlDownloads, urlMaterials } from "@/assets/js/url.js";
export default {
  data() {
    return {
      downloaded: null,
      loading: false,
      material: null,
      notification: "",
      snackbar: false
    };
  },
  created() {
    this.$fetchUser(() => {
      const url = new Url(urlMaterials);
      const data = {
        id: this.$route.params.id
      };
      ajax
        .get(url.base, data)
        .then(response => {
          this.material = response.data[0];
        })
        .then(() => {
          const url = new Url(urlDownloads);
          const data = {
            user_id: this.$user.ID
          };
          ajax.get(url.base, data).then(response => {
            const index = response.data.findIndex(
              download => download.MaterialID === this.material.ID
            );
            const notFound = -1;
            this.downloaded = index !== notFound;
          });
        });
    });
  },
  methods: {
    download() {
      const url = new Url(urlDownloads);
      const data = {
        userID: this.$user.ID,
        materialID: this.material.ID
      };
      this.loading = true;
      ajax.post(url.base, data).then(response => {
        this.loading = false;
        if (response.status === 200) {
          this.downloaded = true;
          this.notification = "教材を取得しました";
        } else {
          this.notification = "教材の取得に失敗しました";
        }
        this.snackbar = true;
      });
    }
  }
};
</script>

<style>
/*
:TODO
*/
@media screen and (min-width: 960px) {
  .lesson-card-container {
    display: grid;
    grid-template-columns: repeat(5, 20%);
  }
}
</style>