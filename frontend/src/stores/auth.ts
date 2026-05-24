import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "../api/client";

export const useAuthStore = defineStore("auth", () => {
  const token = ref<string | null>(localStorage.getItem("token"));
  const user = ref<{
    id?: number;
    displayName?: string;
    username?: string;
    role?: string;
    orgId?: number;
  } | null>(
    JSON.parse(localStorage.getItem("user") || "null"),
  );

  const isAuthed = computed(() => !!token.value);

  function setSession(t: string, u: typeof user.value) {
    token.value = t;
    user.value = u || null;
    localStorage.setItem("token", t);
    localStorage.setItem("user", JSON.stringify(u || null));
    api.defaults.headers.common.Authorization = `Bearer ${t}`;
  }

  function clear() {
    token.value = null;
    user.value = null;
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    delete api.defaults.headers.common.Authorization;
  }

  if (token.value) {
    api.defaults.headers.common.Authorization = `Bearer ${token.value}`;
  }

  return { token, user, isAuthed, setSession, clear };
});
