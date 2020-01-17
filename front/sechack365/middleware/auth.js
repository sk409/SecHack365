import axios from "axios";

export default async function ({
  redirect
}) {
  const url = `${process.env.serverOrigin}/auth`;
  const config = {
    withCredentials: true
  };
  const response = await axios.get(url, config);
  if (response.data.authenticated === false) {
    redirect("/login");
  }
}
