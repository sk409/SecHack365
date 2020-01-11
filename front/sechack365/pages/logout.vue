<template>
  <div></div>
</template>

<script>
import ajax from "@/assets/js/ajax.js";
import Mutations from "@/assets/js/mutations.js";
import { Url, urlLogout } from "@/assets/js/url.js";
import { mapMutations } from "vuex";
export default {
  created() {
    const url = new Url(urlLogout);
    ajax.post(url.base, {}, { withCredentials: true }).then(response => {
      if (response.status === 200) {
        this.setUser(null);
        this.$router.push("/");
      }
    });
  },
  methods: {
    ...mapMutations({
      setUser: Mutations.users().setUser
    })
  }
};
</script>
