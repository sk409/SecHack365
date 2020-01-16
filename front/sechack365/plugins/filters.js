import Vue from "vue";

Vue.filter("truncate", (string, length) => {
  if (string.length <= length) {
    return string;
  }
  return string.substr(string, length) + "...";
});

Vue.filter("date", (string) => {
  const d = new Date(string);
  const year = d.getFullYear();
  const month = d.getMonth() + 1;
  const date = d.getDate();
  return `${year}年 ${month}月${date}日`
})
