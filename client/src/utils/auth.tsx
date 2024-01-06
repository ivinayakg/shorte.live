import fetch from "./axios";
import { setInLocalStorage } from "./localstorage";

export default async function (token: string) {
  if (!token) return;
  const res = await fetch.get("/user/self", {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (res.status !== 200) {
    setInLocalStorage("userToken", null);
    return;
  } else {
    setInLocalStorage("userToken", token);
  }
  return res.data;
}
