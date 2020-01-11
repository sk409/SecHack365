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

export const urlLogin = "login";
export const urlLogout = "logout";
export const urlRegister = "register";
export const urlUser = "user"
export const urlUsers = "users";

export class Url {

  constructor(type) {
    switch (type) {
      default:
        this.base = process.env.serverOrigin + "/" + type;
    }
  }

}
