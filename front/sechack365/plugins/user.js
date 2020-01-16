import ajax from "@/assets/js/ajax.js";
import Vue from "vue";
import {
  Url,
  urlUser
} from "@/assets/js/url.js";

Vue.prototype.$user = null;

Vue.prototype.$fetchUser = (completion) => {
  const url = new Url(urlUser);
  ajax.get(url.base, {}, {
    withCredentials: true
  }).then(response => {
    Vue.prototype.$user = response.data;
    if (completion) {
      completion();
    }
  })
};
