<template>
  <div>
    <div class="pa-4">レッスン一覧</div>
    <v-divider></v-divider>
    <v-container>
      <v-row align="center">
        <v-col cols="6" md="4">
          <v-text-field
            v-model="searchLesson"
            single-line
            hide-details
            prepend-icon="mdi-magnify"
          ></v-text-field>
        </v-col>
        <v-col cols="6" md="2" class="ml-auto">
          <v-menu offset-y>
            <template v-slot:activator="{ on }">
              <v-btn v-on="on">
                <v-icon>mdi-sort</v-icon>
                <span>並べ替え</span>
              </v-btn>
            </template>
            <v-list>
              <v-list-item-group>
                <v-list-item
                  v-for="sortButton in sortButtons"
                  :key="sortButton.title"
                  @click="sortLessonTable(sortButton.key, sortButton.desc)"
                >
                  <v-list-item-title>{{ sortButton.title }}</v-list-item-title>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-menu>
        </v-col>
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
          >
            <template v-slot:item.action="{ item }">
              <v-btn
                icon
                class="mr-2"
                @click="$router.push($routes.lessons.edit(item.ID))"
              >
                <v-icon center>mdi-pencil</v-icon>
              </v-btn>
              <v-btn icon @click="deleteLesson(item)">
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </template>
          </v-data-table>
        </v-col>
      </v-row>
    </v-container>
    <v-btn
      color="accent"
      fab
      large
      fixed
      right
      bottom
      @click="$router.push($routes.lessons.create)"
    >
      <v-icon>mdi-plus</v-icon>
    </v-btn>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlLessons } from "@/assets/js/url.js";
export default {
  layout: "dashboard",
  data() {
    return {
      lessonLoading: true,
      lessons: [],
      tableHeaders: [
        { text: "タイトル", value: "Title" },
        { text: "作成日", value: "CreatedAt" },
        { text: "アクション", value: "action", sortable: false }
      ],
      tableSortDesc: false,
      tableSortKey: null,
      searchLesson: "",
      sortButtons: [
        { title: "タイトル(昇順)", key: "Title", desc: false },
        { title: "タイトル(降順)", key: "Title", desc: true },
        { title: "作成日(昇順)", key: "CreatedAt", desc: false },
        { title: "作成日(降順)", key: "CreatedAt", desc: true }
      ]
    };
  },
  created() {
    const url = new Url(urlLessons);
    const data = {
      user_id: this.$store.state.users.user.ID
    };
    ajax.get(url.base, data).then(response => {
      this.lessonLoading = false;
      this.lessons = response.data;
    });
  },
  methods: {
    broadMatchFilter(value, search, item) {
      if (!value) {
        return;
      }
      return typeof value === "string" && value.includes(search);
    },
    deleteLesson(lesson) {
      console.log(lesson);
    },
    editLesson(lesson) {},
    sortLessonTable(key, desc) {
      this.tableSortKey = key;
      this.tableSortDesc = desc;
    }
  }
};
</script>
