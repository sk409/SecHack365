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
            ></DataTable>
            <!-- <v-toolbar color="secondary" flat>
              <v-text-field color="white" dark hide-details prepend-icon="mdi-magnify" single-line></v-text-field>
              <v-spacer></v-spacer>
              <v-btn color="white" icon>
                <v-icon>mdi-sort</v-icon>
              </v-btn>
            </v-toolbar>
            <v-data-table
              :custom-filter="broadMatchFilter"
              :headers="tableHeaders"
              hide-default-header
              hide-default-footer
              item-key="ID"
              :items="lessons"
              :loading="lessonLoading"
              loading-text="レッスンを読み込んでいます"
              no-data-text="作成したレッスンがありません"
              :search="searchLesson"
              :sort-by="tableSortKey"
              :sort-desc="tableSortDesc"
              @click:row="clickRow"
            >
              <template v-slot:item.action="{ item }">
                <v-btn icon @click="deleteLesson(item)">
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </template>
            </v-data-table>-->
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    <v-btn color="accent" fab fixed right bottom @click="$router.push($routes.lessons.create)">
      <v-icon>mdi-plus</v-icon>
    </v-btn>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import DataTable from "@/components/DataTable.vue";
import { Url, urlLessons } from "@/assets/js/url.js";
export default {
  layout: "dashboard",
  components: {
    DataTable
  },
  data() {
    return {
      lessons: [],
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
    }
  }
};
</script>
