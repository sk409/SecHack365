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
      hide-default-footer
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
        <v-btn icon @click="deleteItem(item)">
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </template>
    </v-data-table>
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
      hideDefaultHeader: true,
      search: "",
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
    deleteItem(item) {},
    sort(key, desc) {
      this.sortKey = key;
      this.sortDesc = desc;
    }
  }
};
</script>