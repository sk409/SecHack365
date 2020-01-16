<template>
  <div>
    <div v-if="error">{{error}}</div>
    <div else></div>
  </div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlMaterials } from "@/assets/js/url.js";
export default {
  data() {
    return {
      error: null,
      material: null
    };
  },
  created() {
    this.$fetchUser(async () => {
      const materialId = this.$route.params.id;
      const url = new Url(urlMaterials);
      const data = {
        id: materialId
      };
      ajax
        .get(url.base, data)
        .then(async response => {
          const material = response.data[0];
          const url = new Url(urlMaterials);
          const data = {
            user_id: this.$user.ID
          };
          const r = await ajax.get(url.base, data);
          const notFound = -1;
          return r.data.findIndex(m => m.ID === material.ID) !== notFound;
        })
        .then(allowed => {
          if (!allowed) {
            this.error = "この教材をダウンロードしていません";
          }
        });
    });
  }
};
</script>