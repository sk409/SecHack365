<template>
  <div>
    <v-card class="mx-auto pa-3">
      <v-card-title>{{ title }}</v-card-title>
      <v-card-text>
        <v-form ref="form">
          <v-text-field
            v-model="username"
            :rules="usernameRules"
            prepend-icon="mdi-account"
            label="名前"
          ></v-text-field>
          <v-text-field
            v-model="password"
            :rules="passwordRules"
            type="password"
            prepend-icon="mdi-lock"
            label="パスワード"
          ></v-text-field>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-btn depressed color="primary" class="mx-auto" @click="submit">
          <v-icon left>{{ buttonIcon }}</v-icon>
          {{ title }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
export default {
  props: {
    buttonIcon: {
      type: String,
      required: true
    },
    title: {
      type: String,
      required: true
    },
    usernames: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      password: "",
      passwordRules: [
        v => !!v || "パスワードを入力してください",
        v => v.length <= 128 || "128文字以内で入力してください"
      ],
      username: "",
      usernameRules: [
        v => !!v || "名前を入力してください",
        v => v.length <= 128 || "128文字以内で入力してください",
        v => !this.usernames.includes(v) || "このユーザ名は既に使われています"
      ]
    };
  },
  methods: {
    submit() {
      if (this.$refs.form.validate()) {
        this.$emit("submit", this.username, this.password);
      }
    }
  }
};
</script>
