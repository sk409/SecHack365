<template>
  <div>
    <div class="pa-4">レッスン作成</div>
    <v-divider></v-divider>
    <v-container class="mt-3">
      <v-card class="pa-4">
        <v-form ref="form">
          <v-file-input
            ref="fileInput"
            v-model="thumbnail"
            label="サムネイル画像"
          ></v-file-input>
          <v-text-field
            v-model="title"
            :rules="titleRules"
            required
            label="タイトル"
          ></v-text-field>
          <v-textarea
            v-model="description"
            :rules="descriptionRules"
            required
            label="説明"
          ></v-textarea>
          <v-select
            v-model="os"
            :items="osList"
            :rules="osRules"
            required
            label="OS"
          ></v-select>
          <v-text-field
            v-model="consolePort"
            type="number"
            label="コンソール用のポート"
          ></v-text-field>
          <v-text-field v-model="ports" label="公開するポート"></v-text-field>
        </v-form>
        <v-card-actions>
          <v-btn
            depressed
            :loading="loading"
            color="primary"
            class="mx-auto"
            @click="submit"
            >作成</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-container>
  </div>
</template>

<script>
import axios from "axios";
import ajax from "@/assets/js/ajax.js";
import { Url, urlLessons } from "@/assets/js/url.js";
export default {
  middleware: "auth",
  data() {
    return {
      consolePort: "",
      consolePortRules: [
        v => !!v || "コンソール用のポートを入力してください",
        v => 1024 <= v || "ウェルノウンポートは指定できません",
        v => v <= 65535 || "65535以下の値を入力してください"
      ],
      description: "",
      descriptionRules: [
        v => !!v || "説明文を入力してください",
        v => v.length <= 2048 || "2048文字以内で入力してください"
      ],
      loading: false,
      os: "",
      osList: ["centos:7"],
      osRules: [v => !!v || "osを選択してください"],
      ports: "",
      thumbnail: null,
      title: "",
      titleRules: [
        v => !!v || "タイトルを入力してください",
        v => v.length <= 512 || "512文字以内で入力してください"
      ]
    };
  },
  methods: {
    submit() {
      // const url = new Url(urlLessons);
      // const data = new FormData();
      // data.append("title", "abc");
      // const config = {
      //   headers: {
      //     "content-type": "multipart/form-data"
      //   }
      // };
      // axios.post(url.base, data, config).then(response => {
      //   console.log(response);
      // });
      if (this.$refs.form.validate()) {
        const url = new Url(urlLessons);
        const data = {
          thumbnail: this.thumbnail,
          title: this.title,
          description: this.description,
          os: this.os,
          consolePort: this.consolePort,
          userID: this.$store.state.users.user.ID,
          "ports[]": this.ports.split(" ")
        };
        const config = {
          headers: {
            "content-type": "multipart/form-data"
          }
        };
        // const form = new FormData();
        // form.append("title", data.title);
        // form.append("description", data.description);
        // form.append("os", data.os);
        // form.append("consolePort", data.consolePort);
        // form.append("ports[]", 80);
        // form.append("ports[]", 3306);
        // form.append("thumbnail", this.thumbnail);
        // form.append("userID", 1);
        // form.append("title", data.title);
        // axios.post(url.base, form, config).then(response => {
        //   console.log(response);
        // });
        // this.loading = true;
        ajax.post(url.base, data, config).then(response => {
          this.loading = false;
          this.$router.push(this.$routes.dashboard.lessons);
        });
      }
    }
  }
};
</script>
