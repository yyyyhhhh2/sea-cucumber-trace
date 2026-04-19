<script setup lang="ts">
import { RouterLink } from "vue-router";
</script>

<template>
  <section class="relative mx-auto max-w-6xl px-6 pb-28 pt-14 md:pt-20">
    <div class="grid gap-14 md:grid-cols-[1.05fr_0.95fr] md:items-center">
      <div class="animate-fade-up">
        <p
          class="mb-4 inline-flex items-center gap-2 rounded-full border border-cyan-400/25 bg-gradient-to-r from-cyan-400/15 to-teal-500/10 px-3 py-1 text-xs font-medium text-cyan-100 shadow-[0_0_24px_-4px_rgba(34,211,238,0.35)]"
        >
          <span class="h-1.5 w-1.5 animate-pulse rounded-full bg-cyan-300 shadow-[0_0_12px_rgba(34,211,238,0.9)]" />
          区块链可信溯源
        </p>
        <h1 class="font-display text-4xl font-bold leading-[1.15] text-white md:text-5xl lg:text-[3.25rem]">
          从海域到餐桌，
          <span
            class="bg-gradient-to-r from-cyan-200 via-teal-200 to-emerald-300 bg-clip-text text-transparent"
            >每一根海参都可验证</span
          >
        </h1>
        <p class="mt-6 max-w-xl text-base leading-relaxed text-slate-400">
          以 Hyperledger Fabric 联盟链为信任底座，关键业务数据在链下高效存储，
          事件哈希与元数据在链上不可篡改地锚定，让消费者扫码即可核对全链路。
        </p>
        <div class="mt-9 flex flex-wrap gap-3">
          <RouterLink
            :to="{ name: 'trace', params: { batchNo: 'HSC-2025-DL-PREMIUM-001' } }"
            class="btn-primary px-7"
          >
            体验溯源
          </RouterLink>
          <RouterLink to="/login" class="btn-ghost"> 企业登录 </RouterLink>
        </div>
        <dl class="mt-14 grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div
            class="glass noise rounded-2xl bg-card-shine p-5 transition hover:border-cyan-400/20"
          >
            <dt class="text-xs font-medium uppercase tracking-wider text-slate-500">链上锚定</dt>
            <dd class="mt-2 font-display text-lg font-semibold text-white">事件级哈希</dd>
            <p class="mt-1 text-xs text-slate-500">SHA-256 指纹与交易绑定</p>
          </div>
          <div
            class="glass noise rounded-2xl bg-card-shine p-5 transition hover:border-cyan-400/20"
          >
            <dt class="text-xs font-medium uppercase tracking-wider text-slate-500">环节覆盖</dt>
            <dd class="mt-2 font-display text-lg font-semibold text-white">养殖 · 加工 · 物流</dd>
            <p class="mt-1 text-xs text-slate-500">时间线可追溯、可核验</p>
          </div>
          <div
            class="glass noise rounded-2xl bg-card-shine p-5 transition hover:border-cyan-400/20"
          >
            <dt class="text-xs font-medium uppercase tracking-wider text-slate-500">开放验证</dt>
            <dd class="mt-2 font-display text-lg font-semibold text-white">公开 API</dd>
            <p class="mt-1 text-xs text-slate-500">消费者与监管同源视图</p>
          </div>
        </dl>
      </div>

      <div class="relative animate-fade-up [animation-delay:120ms]">
        <div
          class="pointer-events-none absolute -left-8 top-1/4 h-40 w-40 rounded-full bg-cyan-500/15 blur-3xl"
        />
        <div
          class="glass noise relative overflow-hidden rounded-[1.75rem] p-7 shadow-card ring-1 ring-white/[0.07]"
        >
          <div
            class="absolute inset-0 bg-gradient-to-br from-cyan-500/15 via-transparent to-teal-600/10"
          />
          <div class="relative space-y-5">
            <div class="flex items-center justify-between gap-3">
              <p class="text-sm font-semibold text-white">批次时间线 · 演示</p>
              <span
                class="rounded-full bg-emerald-400/15 px-2.5 py-1 text-[11px] font-medium text-emerald-200 ring-1 ring-emerald-400/25"
                >已上链</span
              >
            </div>
            <div class="space-y-3 text-sm">
              <div
                class="flex gap-3 rounded-2xl border border-white/[0.08] bg-white/[0.04] p-3.5 transition hover:border-cyan-400/20"
              >
                <div
                  class="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-cyan-400/20 text-xs font-semibold text-cyan-100"
                >
                  养
                </div>
                <div>
                  <p class="font-medium text-slate-100">底播养殖 · 水质抽检</p>
                  <p class="text-xs text-slate-500">哈希写入 Fabric 世界状态</p>
                </div>
              </div>
              <div
                class="flex gap-3 rounded-2xl border border-white/[0.08] bg-white/[0.04] p-3.5 transition hover:border-cyan-400/20"
              >
                <div
                  class="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-teal-400/20 text-xs font-semibold text-teal-100"
                >
                  工
                </div>
                <div>
                  <p class="font-medium text-slate-100">低温加工 · 规格分级</p>
                  <p class="text-xs text-slate-500">与批次号强绑定，可独立校验</p>
                </div>
              </div>
              <div
                class="flex gap-3 rounded-2xl border border-white/[0.08] bg-white/[0.04] p-3.5 transition hover:border-cyan-400/20"
              >
                <div
                  class="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-sky-400/20 text-xs font-semibold text-sky-100"
                >
                  流
                </div>
                <div>
                  <p class="font-medium text-slate-100">冷链运输 · 节点回传</p>
                  <p class="text-xs text-slate-500">监管与消费者同源数据视图</p>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div
          class="pointer-events-none absolute -right-4 -top-8 h-28 w-28 animate-float rounded-full bg-teal-400/20 blur-2xl"
        />
      </div>
    </div>
  </section>
</template>
