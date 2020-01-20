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
            @delete-item="showDialog"
          ></DataTable>
        </v-card>
      </v-col>
    </v-row>
    <v-btn color="accent" fab fixed right bottom :to="$routes.materials.create">
      <v-icon>mdi-plus</v-icon>
    </v-btn>
    <DialogDeletingItem :deleting="deleting" :visible.sync="dialog" @delete-item="deleteMaterial"></DialogDeletingItem>
    <v-snackbar v-model="snackbar" :timeout="2000" top>
      <span>{{notification}}</span>
      <v-btn icon @click="snackbar = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-snackbar>
  </v-container>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import DataTable from "@/components/DataTable.vue";
import DialogDeletingItem from "@/components/DialogDeletingItem.vue";
import { Url, urlMaterials } from "@/assets/js/url.js";
import { defaultDateFormatter } from "@/assets/js/utils.js";
export default {
  layout: "dashboard",
  components: {
    DataTable,
    DialogDeletingItem
  },
  data() {
    return {
      deleting: false,
      dialog: false,
      materials: [],
      notification: "",
      selectedMaterial: null,
      snackbar: false,
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
          material.CreatedAt = defaultDateFormatter(material.CreatedAt);
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
    deleteMaterial() {
      if (!this.selectedMaterial) {
        return;
      }
      this.deleting = true;
      const url = new Url(urlMaterials);
      ajax.delete(url.delete(this.selectedMaterial.ID)).then(response => {
        this.deleting = false;
        this.dialog = false;
        this.snackbar = true;
        if (response.status === 200) {
          this.notification = this.selectedMaterial.Title + "を削除しました";
          this.materials = this.materials.filter(
            material => material.ID !== this.selectedMaterial.ID
          );
        } else {
          this.notification = "削除に失敗しました";
        }
        this.selectedMaterial = null;
      });
    },
    clickRow(material) {
      this.$router.push(this.$routes.materials.showDownloaded(material.ID));
    },
    showDialog(material) {
      this.selectedMaterial = material;
      this.dialog = true;
    }
  }
};
</script>