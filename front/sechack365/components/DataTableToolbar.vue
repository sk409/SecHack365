<template>
  <div class="d-flex w-100">
    <div>
      <v-text-field
        single-line
        hide-details
        prepend-icon="mdi-magnify"
        :value="search"
        @input="inputSearch"
      ></v-text-field>
    </div>
    <div class="ml-auto mr-3">
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
              v-for="sortButton in sort"
              :key="sortButton.title"
              @click="clickSortButton(sortButton)"
            >
              <v-list-item-title>{{ sortButton.title }}</v-list-item-title>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-menu>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    search: {
      type: String,
      default: ""
    },
    sort: {
      type: Array,
      default: []
    }
  },
  methods: {
    inputSearch(search) {
      this.$emit("search", search);
    },
    clickSortButton(sortButton) {
      this.$emit("sort", sortButton.key, sortButton.desc);
    }
  }
};
</script>