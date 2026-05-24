<script setup lang="ts">
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import axios from "axios";
import api from "../api/client";
import { useAuthStore } from "../stores/auth";

const username = ref("orguser");
const password = ref("org123");
const loading = ref(false);
const error = ref("");

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();

function fillDemo(u: string, p: string) {
  username.value = u;
  password.value = p;
  error.value = "";
}

async function submit() {
  error.value = "";
  loading.value = true;
  try {
    const { data } = await api.post("/auth/login", {
      username: username.value.trim(),
      password: password.value,
    });
    auth.setSession(data.token, data.user);
    const redirect = (route.query.redirect as string) || "/dashboard";
    await router.push(redirect);
  } catch (e: unknown) {
    if (axios.isAxiosError(e) && typeof e.message === "string") {
      error.value = e.message;
    } else {
      error.value = "登录失败，请检查账号或后端服务状态。";
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="mx-auto max-w-md px-6 py-16 md:py-20">
    <div
      class="relative overflow-hidden rounded-[1.75rem] border border-white/[0.08] bg-gradient-to-b from-white/[0.07] to-white/[0.02] p-8 shadow-card backdrop-blur-xl"
    >
      <div class="pointer-events-none absolute -right-16 -top-16 h-40 w-40 rounded-full bg-cyan-400/20 blur-3xl" />
      <div class="relative">
        <h1 class="font-display text-2xl font-semibold tracking-tight text-white">企业控制台登录</h1>
        <p class="mt-2 text-sm leading-relaxed text-slate-400">
          登录后可管理批次信息、登记追溯事件，并同步查看链上锚定状态。
        </p>

        <div class="mt-6 flex flex-wrap gap-2">
          <button
            type="button"
            class="rounded-full border border-white/10 bg-white/5 px-3 py-1.5 text-xs text-slate-300 transition hover:border-cyan-400/35 hover:text-white"
            @click="fillDemo('orguser', 'org123')"
          >
            企业账号（演示）
          </button>
          <button
            type="button"
            class="rounded-full border border-white/10 bg-white/5 px-3 py-1.5 text-xs text-slate-300 transition hover:border-cyan-400/35 hover:text-white"
            @click="fillDemo('admin', 'admin123')"
          >
            管理员账号（演示）
          </button>
        </div>

        <form class="mt-8 space-y-4" @submit.prevent="submit">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">账号</label>
            <input
              v-model="username"
              class="input-field"
              autocomplete="username"
              placeholder="请输入账号"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">密码</label>
            <input
              v-model="password"
              type="password"
              class="input-field"
              autocomplete="current-password"
              placeholder="请输入密码"
            />
          </div>
          <p
            v-if="error"
            class="rounded-xl border border-rose-500/25 bg-rose-500/10 px-3 py-2 text-sm text-rose-200"
          >
            {{ error }}
          </p>
          <button type="submit" class="btn-primary mt-2 w-full" :disabled="loading">
            {{ loading ? "登录中..." : "登录并进入工作台" }}
          </button>
        </form>
      </div>
    </div>
  </section>
</template>
