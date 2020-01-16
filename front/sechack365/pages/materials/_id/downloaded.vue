<template>
  <div>
    <div v-if="error">{{error}}</div>
    <div else>
      <v-container v-if="material">
        <v-row>
          <v-col cols="3">
            <v-img :src="$serverUrl(material.ThumbnailPath)"></v-img>
          </v-col>
          <v-col cols="9">
            <div class="title">{{material.Title}}</div>
            <div class="subtitle-1">{{material.author}}</div>
            <div class="body">{{material.Description}}</div>
          </v-col>
        </v-row>
        <v-divider></v-divider>
        <div v-for="lesson in material.lessons" :key="lesson.id" class="mx-auto my-4 lesson">
          <div>
            <v-row align="center">
              <v-col cols="2" style="border-left:5px solid green;">
                <v-avatar size="64">
                  <v-img :src="$serverUrl(lesson.ThumbnailPath)"></v-img>
                </v-avatar>
              </v-col>
              <v-col cols="3">
                <div class="caption">タイトル</div>
                <div>{{lesson.Title}}</div>
              </v-col>
              <v-col cols="2" offset="2">
                <v-btn color="secondary" @click="$router.push($routes.lessons.ide(lesson.ID))">学習</v-btn>
              </v-col>
              <!-- <v-col cols="2" offset="1">
                <v-checkbox></v-checkbox>
              </v-col>-->
            </v-row>
            <v-divider></v-divider>
          </div>
        </div>
      </v-container>
    </div>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlMaterials, urlUsers } from "@/assets/js/url.js";
export default {
  data() {
    return {
      error: null,
      material: null
    };
  },
  created() {
    this.$fetchUser(async () => {
      const materialId = this.$route.params.id;
      const url = new Url(urlMaterials);
      const data = {
        id: materialId
      };
      ajax
        .get(url.base, data)
        .then(async response => {
          const material = response.data[0];
          const url = new Url(urlMaterials);
          const data = {
            user_id: this.$user.ID
          };
          const r = await ajax.get(url.base, data);
          const notFound = -1;
          const allowed =
            r.data.findIndex(m => m.ID === material.ID) !== notFound;
          if (!allowed) {
            this.error = "この教材を取得していません";
            return;
          }
          return material;
        })
        .then(async material => {
          if (!material) {
            return;
          }
          const url = new Url(urlUsers);
          const data = {
            id: material.AuthorUserID
          };
          const response = await ajax.get(url.base, data);
          material.author = response.data[0].Name;
          this.material = material;
        });
    });
  }
};
</script>

<style>
.lesson {
  width: 60%;
}
</style>