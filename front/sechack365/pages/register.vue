<template>
  <div>
    <v-container class="mt-3">
      <v-row justify="center">
        <v-col>
          <AuthForm
            title="登録"
            button-icon="mdi-account-plus"
            :usernames="usernames"
            @submit="register"
          ></AuthForm>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col sm="4">
          <v-btn color="primary" :to="$routes.login.base" block class="mx-auto">アカウントをお持ちの方</v-btn>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import AuthForm from "@/components/AuthForm.vue";
import { Url, urlRegister, urlUsers } from "@/assets/js/url.js";
export default {
  components: {
    AuthForm
  },
  data() {
    return {
      usernames: []
    };
  },
  created() {
    const url = new Url(urlUsers);
    ajax.get(url.base).then(response => {
      this.usernames = response.data.map(user => user.Name);
    });
  },
  methods: {
    register(username, password) {
      const url = new Url(urlRegister);
      const data = {
        username,
        password
      };
      const config = {
        withCredentials: true
      };
      ajax.post(url.base, data, config).then(response => {
        if (response.status === 200) {
          this.$router.push(this.$routes.dashboard.materials.downloaded);
        }
      });
    }
  }
};
</script>
