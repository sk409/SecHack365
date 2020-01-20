<template>
  <v-dialog :value="visible" width="400" @click:outside="$emit('update:visible', false)">
    <v-card class="pa-3">
      <v-card-title>本当に削除しますか?</v-card-title>
      <v-card-text>削除すると元に戻すことができません</v-card-text>
      <v-card-actions>
        <v-btn
          color="error"
          :loading="deleting"
          class="ml-auto mr-4"
          @click="$emit('delete-item')"
        >削除</v-btn>
        <v-btn color="primary" class="mr-auto" @click="$emit('update:visible', false)">キャンセル</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  props: {
    deleting: {
      type: Boolean,
      default: false
    },
    visible: {
      type: Boolean,
      default: false
    }
  }
};
</script>