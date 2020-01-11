<template>
  <div>
    <div class="pa-4">レッスン作成</div>
    <v-divider></v-divider>
    <v-container class="mt-3">
      <v-card class="pa-4">
        <v-form ref="form">
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
          <v-text-field v-model="ports" label="公開するポート"></v-text-field>
          <!-- <v-text-field type="number" label="コンソール用のポート"></v-text-field> -->
        </v-form>
        <v-card-actions>
          <v-btn depressed color="primary" class="mx-auto" @click="submit"
            >作成</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-container>
  </div>
</template>

<script>
export default {
  data() {
    return {
      description: "",
      descriptionRules: [
        v => !!v || "説明文を入力してください",
        v => v.length <= 2048 || "2048文字以内で入力してください"
      ],
      os: "",
      osList: ["CentOS 7"],
      osRules: [v => !!v || "osを選択してください"],
      ports: "",
      title: "",
      titleRules: [
        v => !!v || "タイトルを入力してください",
        v => v.length <= 512 || "512文字以内で入力してください"
      ]
    };
  },
  methods: {
    submit() {
      if (this.$refs.form.validate()) {
      }
    }
  }
};
</script>
