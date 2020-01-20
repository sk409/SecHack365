<template>
  <div>
    <v-container>
      <v-row>
        <v-col>
          <v-card>
            <DataTable
              :headers="tableHeaders"
              :items="lessons"
              :sort-list="sortList"
              @click:row="clickRow"
              @delete-item="deleteLesson"
            ></DataTable>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    <v-snackbar v-model="snackbar" top>
      <span>{{notification}}</span>
      <v-btn icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
    <v-dialog v-model="modeDialog" width="500">
      <v-card class="pa-3">
        <v-card-title>レッスン作成方法の選択</v-card-title>
        <v-card-text>
          <div>新規作成する場合には新規作成を選択してください。</div>
          <div>既存のレッスンを継承する場合には継承を選択してください。</div>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            class="ml-auto mr-4"
            @click="$router.push($routes.lessons.create.inheritance)"
          >継承</v-btn>
          <v-btn
            color="primary"
            class="mr-auto"
            @click="$router.push($routes.lessons.create.new)"
          >新規作成</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-btn color="accent" fab fixed right bottom @click="showModalOrTransition">
      <v-icon>mdi-plus</v-icon>
    </v-btn>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import DataTable from "@/components/DataTable.vue";
import { Url, urlLessons } from "@/assets/js/url.js";
import { defaultDateFormatter } from "@/assets/js/utils.js";
export default {
  layout: "dashboard",
  components: {
    DataTable
  },
  data() {
    return {
      lessons: [],
      modeDialog: false,
      notification: "",
      snackbar: false,
      sortList: [
        { title: "タイトル(昇順)", key: "Title", desc: false },
        { title: "タイトル(降順)", key: "Title", desc: true },
        { title: "作成日(昇順)", key: "CreatedAt", desc: false },
        { title: "作成日(降順)", key: "CreatedAt", desc: true }
      ],
      tableHeaders: [
        { text: "タイトル", value: "Title" },
        { text: "作成日", value: "CreatedAt" },
        { text: "アクション", value: "action", sortable: false }
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
        response.data.forEach(lesson => {
          lesson.CreatedAt = defaultDateFormatter(lesson.CreatedAt);
        });
        this.lessons = response.data;
      });
    });
  },
  mounted() {
    this.$nuxt.$emit("setTitle", "作成したレッスン一覧");
  },
  methods: {
    clickRow(lesson) {
      this.$router.push(this.$routes.lessons.ide(lesson.ID));
    },
    deleteLesson(lesson) {
      const url = new Url(urlLessons);
      ajax.delete(url.delete(lesson.ID)).then(response => {
        this.snackbar = true;
        if (response.status !== 200) {
          this.notification = lesson.Title + "の削除に失敗しました";
          return;
        }
        this.notification = lesson.Title + "を削除しました";
        this.lessons = this.lessons.filter(l => l.ID !== lesson.ID);
      });
    },
    showModalOrTransition() {
      if (this.lessons.length !== 0) {
        this.modeDialog = true;
      } else {
        this.$router.push(this.$routes.lessons.create.new);
      }
    }
  }
};
</script>
