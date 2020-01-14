<template>
  <div style="overflow:scroll;height:100px;">
    <div v-for="item in items" :key="item.name">
      <FileTreeItem :item="item" @click="clickItem"></FileTreeItem>
    </div>
  </div>
</template>

<script>
import FileTreeItem from "@/components/FileTreeItem.vue";
export default {
  components: {
    FileTreeItem
  },
  props: {
    items: {
      type: Array,
      default: []
    },
    loadChildren: {
      type: Function
    }
  },
  methods: {
    clickItem(e, item) {
      if (
        Array.isArray(item.children) &&
        item.children.length == 0 &&
        this.loadChildren
      ) {
        this.loadChildren(item);
      }
      this.$emit("click-item", item);
    }
  }
};
</script>