<template>
  <div>
    <v-container>
      <v-row>
        <DataTableToolbar :sort="sortButtons" @search="setSearchMaterial" @sort="sortMaterials"></DataTableToolbar>
      </v-row>
      <v-row>
        <v-col>
          <v-data-table
            :custom-filter="broadMatchFilter"
            :headers="tableHeaders"
            item-key="ID"
            :items="materials"
            :search="searchMaterial"
            :sort-by="sortKey"
            :sort-desc="sortDesc"
            @click:row="clickRow"
          >
            <template v-slot:item.action="{item}">
              <v-btn icon @click="deleteMaterial(item)">
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </template>
          </v-data-table>
        </v-col>
      </v-row>
      <v-btn color="accent" fab fixed right bottom :to="$routes.materials.create">
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-container>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import DataTableToolbar from "@/components/DataTableToolbar.vue";
import { Url, urlMaterials, urlUsers } from "@/assets/js/url.js";
export default {
  layout: "dashboard",
  components: {
    DataTableToolbar
  },
  data() {
    return {
      materials: [],
      searchMaterial: "",
      sortButtons: [
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
      sortDesc: true,
      sortKey: "CreatedAt",
      tableHeaders: [
        { text: "タイトル", value: "Title" },
        { text: "作成者", value: "author" },
        { text: "レッスン数", value: "length" },
        { text: "取得日", value: "CreatedAt" },
        { text: "アクション", value: "action", sortable: false }
      ]
    };
  },
  created() {
    this.$fetchUser(() => {
      const url = new Url(urlUsers);
      ajax
        .get(url.base)
        .then(response => {
          return response.data;
        })
        .then(users => {
          const url = new Url(urlMaterials);
          const data = {
            user_id: this.$user.ID,
            downloaded: 1
          };
          ajax.get(url.base, data).then(response => {
            // console.log(response);
            if (response.status !== 200) {
              return;
            }
            response.data.forEach(material => {
              const user = users.find(user => user.ID === material.UserID);
              if (!user) {
                return;
              }
              material.author = user.Name;
              material.length = material.lessons.length;
            });
            this.materials = response.data;
          });
        });
    });
  },
  methods: {
    broadMatchFilter(value, serach, item) {
      if (!value) {
        return;
      }
      return typeof value === "string" && value.includes(this.searchMaterial);
    },
    clickRow(material) {
      this.$router.push(this.$routes.materials.showDownloaded(material.ID));
    },
    deleteMaterial(material) {
      console.log(material);
    },
    setSearchMaterial(searchMaterial) {
      this.searchMaterial = searchMaterial;
    },
    sortMaterials(key, desc) {
      this.sortKey = key;
      this.sortDesc = desc;
    }
  }
};
</script>
