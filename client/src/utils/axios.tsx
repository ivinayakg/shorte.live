import axios from "axios";

export const baseURL = import.meta.env.VITE_BASE_URL;

const fetch = axios.create({
  baseURL: baseURL,
});
export default fetch;
