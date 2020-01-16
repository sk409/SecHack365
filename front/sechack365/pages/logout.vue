<template>
  <div></div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import { Url, urlLogout } from "@/assets/js/url.js";
export default {
  middleware: "auth",
  created() {
    const url = new Url(urlLogout);
    ajax.post(url.base, {}, { withCredentials: true }).then(response => {
      if (response.status === 200) {
        this.$router.push("/");
      }
    });
  }
};
</script>
