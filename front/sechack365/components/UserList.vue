<template>
  <v-list>
    <template v-for="item in items">
      <v-divider v-if="item.divider" :key="item.id" :inset="item.inset"></v-divider>
      <v-list-item v-else :key="item.id" :to="$routes.users.show(item.id)">
        <v-list-item-avatar size="64">
          <v-img :src="item.avatar" class="thumbnail-1 thumbnail-left"></v-img>
        </v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>{{item.name}}</v-list-item-title>
          <!-- <v-list-item-subtitle>{{item.biotext}}</v-list-item-subtitle> -->
        </v-list-item-content>
      </v-list-item>
    </template>
  </v-list>
</template>

<script>
import { uuid } from "@/assets/js/utils.js";
export default {
  props: {
    users: {
      type: Array,
      default: []
    }
  },
  computed: {
    items() {
      const i = [];
      for (const user of this.users) {
        i.push(user);
        i.push({
          id: uuid(),
          divider: true,
          inset: true
        });
      }
      return i;
    }
  }
};
</script>