<template>
  <v-container>
    <v-row>
      <v-col cols="3">
        <img :src="profileImagePath" class="thumbnail-3 thumbnail-top" />
      </v-col>
      <v-col cols="9">
        <div class="title">{{user ? user.Name : ''}}</div>
        <div class="body mt-2">バイオテキスト</div>
      </v-col>
    </v-row>
    <v-divider></v-divider>
    <v-subheader>教材一覧</v-subheader>
    <div class="materials">
      <v-card
        v-for="material in materials"
        :key="material.ID"
        class="ma-2"
        @click="$router.push($routes.materials.show(material.ID))"
      >
        <img v-if="material" :src="$serverUrl(material.ThumbnailPath)" class="thumbnail-2" />
        <v-divider></v-divider>
        <v-card-title>{{material.Title | truncate(12)}}</v-card-title>
        <v-card-subtitle>全{{material.lessons.length}}レッスン</v-card-subtitle>
        <v-card-text>{{material.Description | truncate(30)}}</v-card-text>
      </v-card>
    </div>
    <v-snackbar v-model="snackbar" :timeout="2000" top>
      <span>{{notification}}</span>
      <v-btn color="red" icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
    <v-btn v-if="!isMe" color="accent" fab fixed right bottom @click="clickHeartButton">
      <v-icon>{{buttonIcon}}</v-icon>
    </v-btn>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlFollows, urlMaterials, urlUsers } from "@/assets/js/url.js";
export default {
  data() {
    return {
      following: false,
      followModel: null,
      materials: [],
      notification: "",
      snackbar: false,
      user: null
    };
  },
  computed: {
    buttonIcon() {
      return this.following ? "mdi-heart-broken" : "mdi-heart";
    },
    isMe() {
      if (!this.$user || !this.user) {
        return false;
      }
      return this.$user.ID === this.user.ID;
    },
    profileImagePath() {
      return this.user
        ? process.env.serverOrigin + this.user.ProfileImagePath
        : "";
    }
  },
  created() {
    this.$fetchUser(() => {
      const url = new Url(urlUsers);
      const userID = this.$route.params.id;
      const data = {
        id: userID
      };
      ajax
        .get(url.base, data)
        .then(response => {
          this.user = response.data[0];
        })
        .then(() => {
          const url = new Url(urlMaterials);
          const data = {
            user_id: this.user.ID,
            downloaded: 0
          };
          ajax.get(url.base, data).then(response => {
            this.materials = response.data;
          });
        })
        .then(() => {
          const url = new Url(urlFollows);
          const data = {
            following_user_id: this.$user.ID
          };
          ajax.get(url.base, data).then(response => {
            const index = response.data.findIndex(
              follow => follow.FollowedUserID === this.user.ID
            );
            const notFound = -1;
            this.following = index !== notFound;
            if (index !== notFound) {
              this.followModel = response.data[index];
            }
          });
        });
    });
  },
  methods: {
    clickHeartButton() {
      if (this.following) {
        this.unfollow();
      } else {
        this.follow();
      }
    },
    follow() {
      const url = new Url(urlFollows);
      const data = {
        followingUserID: this.$user.ID,
        followedUserID: this.user.ID
      };
      ajax.post(url.base, data).then(response => {
        this.notification =
          response.status === 200
            ? "フォローしました"
            : "フォローに失敗しました";
        this.snackbar = true;
        this.following = response.status === 200;
        if (response.status === 200) {
          this.following = true;
          this.followModel = response.data;
        }
      });
    },
    unfollow() {
      const url = new Url(urlFollows);
      ajax.delete(url.delete(this.followModel.ID)).then(response => {
        this.notification =
          response.status === 200
            ? "フォローを解除しました"
            : "フォロー解除に失敗しました";
        this.snackbar = true;
        if (response.status === 200) {
          this.following = false;
          this.followModel = null;
        }
      });
    }
  }
};
</script>

<style>
@media screen and (min-width: 960px) {
  .materials {
    display: grid;
    grid-template-columns: repeat(5, 20%);
  }
}
</style>