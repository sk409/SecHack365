import Axios from "axios";

class Ajax {

  get(url, data, config) {
    url += "?";
    for (const key in data) {
      url += `${key}=${data[key]}&`;
    }
    return Axios.get(url, config)
  }

  post(url, data, config) {
    const params = new URLSearchParams();
    if (data) {
      for (const key in data) {
        params.append(key, data[key]);
      }
    }
    return Axios.post(url, params, config);
  }
}

const ajax = new Ajax();
export default ajax;
