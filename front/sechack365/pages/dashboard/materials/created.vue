<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card>
          <DataTable
            :headers="tableHeaders"
            :items="materials"
            :sort-list="sortList"
            @click:row="clickRow"
          ></DataTable>
        </v-card>
      </v-col>
    </v-row>
    <v-btn color="accent" fab fixed right bottom :to="$routes.materials.create">
      <v-icon>mdi-plus</v-icon>
    </v-btn>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import DataTable from "@/components/DataTable.vue";
import { Url, urlMaterials } from "@/assets/js/url.js";
export default {
  layout: "dashboard",
  components: {
    DataTable
  },
  data() {
    return {
      materials: [],
      sortList: [
        {
          title: "タイトル(昇順)",
          key: "Title",
          desc: false
        },
        {
          title: "タイトル(降順)",
          key: "Title",
          desc: true
        },
        {
          title: "作成日(昇順)",
          key: "CreatedAt",
          desc: false
        },
        {
          title: "作成日(降順)",
          key: "CreatedAt",
          desc: true
        }
      ],
      tableHeaders: [
        { text: "タイトル", value: "Title" },
        { text: "レッスン数", value: "length" },
        { text: "作成日", value: "CreatedAt" },
        { text: "アクション", value: "action", sortable: false }
      ]
    };
  },
  created() {
    this.$fetchUser(() => {
      const url = new Url(urlMaterials);
      const data = {
        user_id: this.$user.ID,
        downloaded: 0
      };
      ajax.get(url.base, data).then(response => {
        if (response.status !== 200) {
          return;
        }
        response.data.forEach(material => {
          material.length = material.lessons.length;
        });
        this.materials = response.data;
      });
    });
  },
  mounted() {
    this.$nuxt.$emit("setTitle", "作成した教材一覧");
  },
  methods: {
    clickRow(material) {
      this.$router.push(this.$routes.materials.showDownloaded(material.ID));
    }
  }
};
</script>