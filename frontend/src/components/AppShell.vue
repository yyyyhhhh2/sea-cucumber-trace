<script setup lang="ts">
import { RouterLink, useRoute } from "vue-router";
import { useAuthStore } from "../stores/auth";
import WaveBackground from "./WaveBackground.vue";

const auth = useAuthStore();
const route = useRoute();

const demoBatch = "HSC-2025-DL-PREMIUM-001";

function navClass(name: string | string[]) {
  const names = Array.isArray(name) ? name : [name];
  const on = names.includes(String(route.name));
  return [
    "relative rounded-full px-3 py-1.5 text-sm transition",
    on
      ? "bg-white/10 text-white shadow-[inset_0_0_0_1px_rgba(34,211,238,0.25)]"
      : "text-slate-400 hover:bg-white/5 hover:text-slate-200",
  ];
}
</script>

<template>
  <div class="relative min-h-screen overflow-x-hidden bg-abyss-950">
    <WaveBackground />
    <header class="relative z-20 border-b border-white/[0.06] bg-abyss-950/75 backdrop-blur-xl">
      <div class="mx-auto flex max-w-6xl items-center justify-between gap-4 px-6 py-4">
        <RouterLink to="/" class="group flex min-w-0 items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl bg-gradient-to-br from-cyan-400 to-teal-600 shadow-glow ring-1 ring-white/10 transition group-hover:ring-cyan-300/40"
          >
            <span class="text-lg font-bold text-slate-950">S</span>
          </div>
          <div class="min-w-0">
            <p class="font-display text-sm font-semibold tracking-wide text-white">SeaTrace Cloud</p>
            <p class="truncate text-xs text-slate-500">海参供应链可信追溯平台</p>
          </div>
        </RouterLink>

        <nav class="hidden items-center gap-1 md:flex">
          <RouterLink :class="navClass('home')" to="/">产品首页</RouterLink>
          <RouterLink :class="navClass('trace')" :to="{ name: 'trace', params: { batchNo: demoBatch } }">
            公开查询
          </RouterLink>
          <RouterLink v-if="auth.isAuthed" :class="navClass('dashboard')" to="/dashboard">
            企业工作台
          </RouterLink>
        </nav>

        <div class="flex shrink-0 items-center gap-2">
          <RouterLink
            v-if="!auth.isAuthed"
            to="/login"
            class="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-white transition hover:border-cyan-400/40 hover:bg-white/10"
          >
            企业登录
          </RouterLink>
          <button
            v-else
            type="button"
            class="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-white transition hover:border-cyan-400/40"
            @click="auth.clear(); $router.push('/')"
          >
            退出账号
          </button>
        </div>
      </div>
    </header>

    <main class="relative z-10 min-h-[calc(100vh-8rem)]">
      <slot />
    </main>

    <footer class="relative z-10 border-t border-white/[0.06] bg-gradient-to-t from-abyss-950 to-transparent py-12">
      <div class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 text-xs text-slate-500 md:flex-row">
        <p class="text-center md:text-left">
          SeaTrace Cloud · 面向品牌商、加工企业与渠道商的一体化可信追溯服务
        </p>
        <div class="flex flex-wrap items-center justify-center gap-4">
          <RouterLink class="link-subtle" :to="{ name: 'trace', params: { batchNo: demoBatch } }">
            查看示例批次
          </RouterLink>
          <span class="hidden text-slate-700 sm:inline">|</span>
          <span class="text-slate-600">
            API
            <code class="rounded bg-white/5 px-1.5 py-0.5 font-mono text-[11px] text-slate-400">/api/health</code>
          </span>
        </div>
      </div>
    </footer>
  </div>
</template>
