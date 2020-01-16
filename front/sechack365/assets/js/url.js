// const makeUrl = (path) => {
//   return process.env.serverOrigin + "/" + path;
// }

// export const makeQuery = (params) => {
//   let query = "?";
//   for (const key in params) {
//     const value = params[key];
//     query += `${key}=${value}&`
//   }
//   return query;
// }

// export const urlLogin = {
//   base: makeUrl("login")
// }

// export const urlRegister = {
//   base: makeUrl("register")
// }

// export const urlUsers

const routeWithID = (base, id) => {
  const url = base;
  if (!url.endsWith("/")) {
    url += "/";
  }
  return url + id;
}

export const urlDownloads = "downloads";
export const urlFiles = "files";
export const urlFolders = "folders";
export const urlFollows = "follows/";
export const urlLessons = "lessons/";
export const urlMaterials = "materials/"
export const urlLogin = "login";
export const urlLogout = "logout";
export const urlRegister = "register";
export const urlUser = "user";
export const urlUsers = "users/";

export class Url {
  constructor(type) {
    this.base = process.env.serverOrigin + "/" + type;
    this.update = (id) => {
      return routeWithID(this.base, id);
    };
    this.delete = (id) => {
      return routeWithID(this.base, id);
    }
    switch (type) {
      case urlUsers:
        this.search = this.base + "search";
        this.follow = (id) => {
          return routeWithID(this.base, id) + "/follow"
        };
        this.follower = (id) => {
          return routeWithID(this.base, id) + "/follower"
        };
    }
  }
}
