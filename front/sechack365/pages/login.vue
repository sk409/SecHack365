<template>
  <div>
    <v-container class="mt-3">
      <v-row justify="center">
        <v-col>
          <AuthForm button-icon="mdi-login" title="ログイン" @submit="login"></AuthForm>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col sm="4">
          <v-btn color="primary" :to="$routes.register.base" block class="mx-auto">アカウントをお持ちでない方</v-btn>
        </v-col>
      </v-row>
    </v-container>
    <v-snackbar v-model="snackbar" top>
      <span>{{ errorMessage }}</span>
      <v-btn icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import AuthForm from "@/components/AuthForm.vue";
import { Url, urlLogin, urlUsers } from "@/assets/js/url.js";
export default {
  components: {
    AuthForm
  },
  data() {
    return {
      errorMessage: "",
      snackbar: false
    };
  },
  methods: {
    login(username, password) {
      const url = new Url(urlLogin);
      const data = {
        username,
        password
      };
      const config = {
        withCredentials: true
      };
      ajax.post(url.base, data, config).then(response => {
        if (response.data === "USERNAME_DOES_NOT_EXIST") {
          this.errorMessage = "ユーザ名が登録されていません";
          this.snackbar = true;
        } else if (response.data === "PASSWORD_DOEN_NOT_MATCH") {
          this.errorMessage = "パスワードが正しくありません";
          this.snackbar = true;
        } else {
          this.$router.push(this.$routes.dashboard.materials.downloaded);
        }
      });
    }
  }
};
</script>
