import Vue from "vue";
Vue.prototype.$routes = {
  dashboard: {
    base: "/dashboard",
    account: "/dashboard/account",
    follow: "/dashboard/follow",
    follower: "/dashboard/follower",
    lessons: "/dashboard/lessons",
    materials: {
      created: "/dashboard/materials/created",
      downloaded: "/dashboard/materials/downloaded"
    }
  },
  home: "/",
  lessons: {
    base: "/lessons",
    create: {
      inheritance: "/lessons/create/inheritance",
      new: "/lessons/create/new",
    },
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
  materials: {
    create: "/materials/create",
    show(id) {
      return "/materials/" + id;
    },
    showDownloaded(id) {
      return "/materials/" + id + "/downloaded";
    }
  },
  register: {
    base: "/register"
  },
  users: {
    edit(id) {
      return "/users/" + id + "/edit";
    },
    show(id) {
      return "/users/" + id;
    }
  }
};
