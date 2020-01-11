<template>
  <div>
    <div class="pa-4">レッスン一覧</div>
    <v-divider></v-divider>
    <v-container>
      <v-row class="d-flex">
        <v-col cols="12" md="4">
          <v-text-field
            prepend-icon="mdi-magnify"
            single-line
            hide-details
          ></v-text-field>
        </v-col>
        <v-col cols="12" md="2" class="ml-auto">
          <v-menu offset-y>
            <template v-slot:activator="{ on }">
              <v-btn v-on="on">
                <v-icon>mdi-sort</v-icon>
                <span>並べ替え</span>
              </v-btn>
            </template>
            <v-list>
              <v-list-item-group>
                <v-list-item>
                  <v-list-item-title>名前</v-list-item-title>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>作成日</v-list-item-title>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-menu>
        </v-col>
      </v-row>
    </v-container>
    <v-btn
      color="accent"
      fab
      large
      fixed
      right
      bottom
      @click="$router.push($routes.lessons.create)"
    >
      <v-icon>mdi-plus</v-icon>
    </v-btn>
  </div>
</template>

<script>
export default {
  layout: "dashboard"
};
</script>
