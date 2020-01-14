<template>
  <div>
    <div class="pa-2" @click.stop="click">
      <v-icon>{{item.icon}}</v-icon>
      <span>{{item.name}}</span>
    </div>
    <div v-if="item.children && item.open">
      <div v-for="child in item.children" :key="child.name">
        <FileTreeItem :item="child" class="ml-2" @click="clickChild"></FileTreeItem>
      </div>
    </div>
  </div>
</template>

<script>
import FileTreeItem from "@/components/FileTreeItem.vue";
export default {
  name: "FileTreeItem",
  components: {
    FileTreeItem
  },
  props: {
    item: {
      type: Object,
      required: true
    }
  },
  methods: {
    click(e) {
      this.$emit("click", e, this.item);
    },
    clickChild(e, item) {
      this.$emit("click", e, item);
    }
  }
};
</script>