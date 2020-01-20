<template>
  <div>
    <v-toolbar color="secondary" flat>
      <v-text-field
        v-model="search"
        color="white"
        dark
        hide-details
        prepend-icon="mdi-magnify"
        single-line
      ></v-text-field>
      <v-spacer></v-spacer>
      <v-menu offset-y>
        <template v-slot:activator="{on}">
          <v-btn color="white" icon v-on="on">
            <v-icon>mdi-sort</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item
            v-for="sortButton in sortList"
            :key="sortButton.title"
            @click="sort(sortButton.key, sortButton.desc)"
          >
            <v-list-item-title>{{sortButton.title}}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-toolbar>
    <v-data-table
      :custom-filter="broadMatchFilter"
      :headers="headers"
      :hide-default-header="hideDefaultHeader"
      item-key="ID"
      :items="items"
      no-data-text="データがありません"
      :search="search"
      :sort-by="sortKey"
      :sort-desc="sortDesc"
      @click:row="clickRow"
    >
      <template v-slot:header.name="{header}">{{header.text}}</template>
      <template v-slot:item.action="{ item }">
        <v-btn icon @click.stop="showDeletingDialog(item)">
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </template>
    </v-data-table>
    <v-dialog v-model="deletingDialog" width="400">
      <v-card class="pa-3">
        <v-card-title class>本当に削除しますか?</v-card-title>
        <v-card-text>削除すると元に戻すことができません</v-card-text>
        <v-card-actions>
          <v-btn color="error" class="ml-auto mr-4" @click="deleteItem">削除</v-btn>
          <v-btn color="primary" class="mr-auto" @click="deletingDialog=false">キャンセル</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
  props: {
    headers: {
      type: Array,
      required: true
    },
    items: {
      type: Array,
      required: true
    },
    sortList: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      deletingDialog: false,
      hideDefaultHeader: true,
      search: "",
      selectedItem: null,
      sortDesc: false,
      sortKey: null
    };
  },
  created() {
    this.hideDefaultHeader = window.innerWidth < 600;
    onresize = () => {
      this.hideDefaultHeader = window.innerWidth < 600;
    };
  },
  methods: {
    broadMatchFilter(value, search, item) {
      if (!value) {
        return;
      }
      return typeof value === "string" && value.includes(search);
    },
    clickRow(item) {
      this.$emit("click:row", item);
    },
    deleteItem() {
      this.$emit("delete-item", this.selectedItem);
      this.deletingDialog = false;
    },
    showDeletingDialog(item) {
      this.deletingDialog = true;
      this.selectedItem = item;
    },
    sort(key, desc) {
      this.sortKey = key;
      this.sortDesc = desc;
    }
  }
};
</script>