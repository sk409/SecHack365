<template>
  <div>
    <div class="pa-4">レッスン一覧</div>
    <v-divider></v-divider>
    <v-container>
      <v-row align="center">
        <DataTableToolbar
          :search="searchLesson"
          :sort="sortButtons"
          @search="setSearchLesson"
          @sort="sortLessonTable"
        ></DataTableToolbar>
      </v-row>
      <v-row>
        <v-col>
          <v-data-table
            :custom-filter="broadMatchFilter"
            :headers="tableHeaders"
            item-key="ID"
            :items="lessons"
            :loading="lessonLoading"
            loading-text="レッスンを取得しています"
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
          </v-data-table>
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
import DataTableToolbar from "@/components/DataTableToolbar.vue";
import { Url, urlLessons } from "@/assets/js/url.js";
export default {
  layout: "dashboard",
  components: {
    DataTableToolbar
  },
  data() {
    return {
      lessonLoading: true,
      lessons: [],
      searchLesson: "",
      sortButtons: [
        { title: "タイトル(昇順)", key: "Title", desc: false },
        { title: "タイトル(降順)", key: "Title", desc: true },
        { title: "作成日(昇順)", key: "CreatedAt", desc: false },
        { title: "作成日(降順)", key: "CreatedAt", desc: true }
      ],
      tableHeaders: [
        { text: "タイトル", value: "Title" },
        { text: "作成日", value: "CreatedAt" },
        { text: "アクション", value: "action", sortable: false }
      ],
      tableSortDesc: false,
      tableSortKey: null
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
        this.lessonLoading = false;
        this.lessons = response.data;
      });
    });
  },
  methods: {
    broadMatchFilter(value, search, item) {
      if (!value) {
        return;
      }
      return typeof value === "string" && value.includes(search);
    },
    clickRow(lesson) {
      this.$router.push(this.$routes.lessons.edit(lesson.ID));
    },
    deleteLesson(lesson) {
      console.log(lesson);
    },
    setSearchLesson(searchLesson) {
      this.searchLesson = searchLesson;
    },
    sortLessonTable(key, desc) {
      this.tableSortKey = key;
      this.tableSortDesc = desc;
    }
  }
};
</script>
