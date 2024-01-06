import axios from "axios";

const baseURL = "";

const fetch = axios.create({
  baseURL: baseURL,
});

export default fetch;
