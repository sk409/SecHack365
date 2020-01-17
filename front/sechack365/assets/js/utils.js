export function uuid() {
  const chars = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".split("");
  for (let i = 0, len = chars.length; i < len; ++i) {
    switch (chars[i]) {
      case "x":
        chars[i] = Math.floor(Math.random() * 16).toString(16);
        break;
      case "y":
        chars[i] = (Math.floor(Math.random() * 4) + 8).toString(16);
        break;
    }
  }
  return chars.join("");
}

export function defaultDateFormatter(dateString) {
  const date = new Date(dateString);
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const d = date.getDate();
  const hours = date.getHours();
  const minutes = date.getMinutes();
  return `${year}/${month}/${d} ${hours}:${minutes}`;
}
