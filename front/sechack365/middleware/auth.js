import axios from "axios";

export default async function ({
  redirect
}) {
  const url = `${process.env.serverOrigin}/auth`
  const response = await axios.get(url);
  if (response.data.authenticated === false) {
    redirect("/login");
  }
}
