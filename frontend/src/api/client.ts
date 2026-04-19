import axios, { type AxiosError } from "axios";

const api = axios.create({ baseURL: "/api", timeout: 30000 });

api.interceptors.response.use(
  (res) => res,
  (err: AxiosError<{ error?: string }>) => {
    const msg = err.response?.data?.error;
    if (typeof msg === "string" && msg.length > 0) {
      err.message = msg;
    }
    return Promise.reject(err);
  },
);

export default api;
