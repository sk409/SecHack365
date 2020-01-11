export default class Mutations {
  static users() {
    const namespace = "users/";
    return {
      setUser: namespace + "setUser"
    }
  }
}
