<template>
  <div>
    <v-container>
      <v-card>
        <v-card-text>
          <v-form>
            <v-file-input v-model="profileImage"></v-file-input>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-btn color="accent" class="mx-auto" @click="submit">
            <v-icon left>mdi-update</v-icon>更新
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-container>
    <v-snackbar v-model="snackbar" :timeout="2000" top>
      <span>{{notification}}</span>
      <v-btn icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlUsers } from "@/assets/js/url.js";
export default {
  data() {
    return {
      notification: "",
      profileImage: null,
      snackbar: false
    };
  },
  created() {
    this.$fetchUser();
  },
  methods: {
    submit() {
      const data = {
        profileImage: this.profileImage
      };
      const config = {
        headers: {
          "content-type": "multipart/form-data"
        }
      };
      const url = new Url(urlUsers);
      ajax.put(url.update(this.$user.ID), data, config).then(response => {
        this.snackbar = true;
        if (response.status === 200) {
          this.notification = "更新しました。";
        } else {
          this.notification = "更新しました。";
        }
      });
    }
  }
};
</script>