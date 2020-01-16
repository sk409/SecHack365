import Axios from "axios";

const makeData = (data, config) => {
  let params = new URLSearchParams();
  if (config && config.headers && config.headers["content-type"] === "multipart/form-data") {
    params = new FormData();
  }
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
  return params;
}

class Ajax {
  get(url, data, config) {
    url += "?";
    for (const key in data) {
      url += `${key}=${data[key]}&`;
    }
    return Axios.get(url, config);
  }

  post(url, data, config) {
    return Axios.post(url, makeData(data, config), config);
  }

  put(url, data, config) {
    return Axios.put(url, makeData(data, config), config);
  }

  delete(url, data, config) {
    return Axios.delete(url, makeData(data, config), config);
  }
}

const ajax = new Ajax();
export default ajax;
