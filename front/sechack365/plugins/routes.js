import Vue from "vue";
Vue.prototype.$routes = {
  dashboard: {
    base: "/dashboard",
    lessons: "/dashboard/lessons",
    materials: {
      created: "/dashboard/materials/created",
      downloaded: "/dashboard/materials/downloaded"
    }
  },
  lessons: {
    base: "/lessons",
    create: "/lessons/create",
    edit(id) {
      return "/lessons/" + id + "/edit";
    },
    ide(id) {
      return "/lessons/" + id + "/ide";
    },
    write(id) {
      return "/lessons/" + id + "/write";
    }
  },
  login: {
    base: "/login"
  },
  logout: {
    base: "/logout"
  },
  register: {
    base: "/register"
  }
};
