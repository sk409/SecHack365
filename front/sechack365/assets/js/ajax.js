import Axios from "axios";

class Ajax {
  get(url, data, config) {
    url += "?";
    for (const key in data) {
      url += `${key}=${data[key]}&`;
    }
    return Axios.get(url, config);
  }

  post(url, data, config) {
    const params = new FormData();
    if (data) {
      for (const key in data) {
        const value = data[key];
        if (key.endsWith("[]")) {
          value.forEach(entry => {
            params.append(key, entry);
          });
        } else {
          params.append(key, value);
        }
      }
    }
    return Axios.post(url, params, config);
  }
}

const ajax = new Ajax();
export default ajax;
