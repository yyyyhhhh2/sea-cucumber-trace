<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";

const metrics = [
  { value: "99.98%", label: "数据完整率", note: "关键节点自动计算哈希并留存审计线索" },
  { value: "12", label: "协同机构", note: "覆盖养殖、加工、仓配、渠道与品牌运营" },
  { value: "T+0", label: "追溯响应", note: "公开查询页秒级返回，支持验收现场演示" },
];

const scenes = [
  {
    eyebrow: "Supply Chain",
    title: "业务系统和可信存证在同一条流水线上协作",
    body: "企业侧保留高效录入与管理体验，链上侧仅锚定关键摘要，保障可追溯与可核验。",
  },
  {
    eyebrow: "Compliance",
    title: "面向监管与审计的证据结构",
    body: "批次、事件、时间戳、操作人、地点与交易摘要形成完整证据链，方便复核。",
  },
  {
    eyebrow: "Brand Trust",
    title: "面向消费者的公开查询体验",
    body: "从“能查到”升级到“看得懂”，用清晰时间线呈现产品流转与质量背书。",
  },
];

const heroRef = ref<HTMLElement | null>(null);
const pointerX = ref(50);
const pointerY = ref(30);

function onPointerMove(event: PointerEvent) {
  const el = heroRef.value;
  if (!el) return;
  const rect = el.getBoundingClientRect();
  pointerX.value = ((event.clientX - rect.left) / rect.width) * 100;
  pointerY.value = ((event.clientY - rect.top) / rect.height) * 100;
  el.style.setProperty("--spotlight-x", `${pointerX.value}%`);
  el.style.setProperty("--spotlight-y", `${pointerY.value}%`);
}

function onPointerLeave() {
  const el = heroRef.value;
  if (!el) return;
  pointerX.value = 50;
  pointerY.value = 30;
  el.style.setProperty("--spotlight-x", "50%");
  el.style.setProperty("--spotlight-y", "30%");
}

onMounted(() => {
  onPointerLeave();
  heroRef.value?.addEventListener("pointermove", onPointerMove);
  heroRef.value?.addEventListener("pointerleave", onPointerLeave);
});

onBeforeUnmount(() => {
  heroRef.value?.removeEventListener("pointermove", onPointerMove);
  heroRef.value?.removeEventListener("pointerleave", onPointerLeave);
});
</script>

<template>
  <section class="relative mx-auto max-w-6xl px-6 pb-24 pt-10 md:pt-16">
    <div
      ref="heroRef"
      class="spotlight-card relative overflow-hidden rounded-[2.35rem] border border-white/[0.06] bg-white/[0.03] px-6 py-10 shadow-card backdrop-blur-xl md:px-10 md:py-12"
    >
      <div class="absolute inset-0 bg-[linear-gradient(145deg,rgba(255,255,255,0.08),transparent_28%,rgba(34,211,238,0.08)_64%,rgba(250,204,21,0.04))]" />
      <div
        class="absolute right-[-8rem] top-[-7rem] h-96 w-96 rounded-full bg-cyan-300/10 blur-3xl animate-drift"
        :style="{ transform: `translate(${(pointerX - 50) * 0.12}px, ${(pointerY - 40) * 0.08}px)` }"
      />
      <div
        class="absolute bottom-[-10rem] left-[-8rem] h-80 w-80 rounded-full bg-teal-300/12 blur-3xl animate-sway"
        :style="{ transform: `translate(${(50 - pointerX) * 0.1}px, ${(45 - pointerY) * 0.08}px)` }"
      />

      <div class="relative grid gap-10 lg:grid-cols-[1.05fr_0.95fr] lg:items-center">
        <div class="animate-fade-up">
          <p class="inline-flex items-center gap-2 rounded-full border border-cyan-300/20 bg-cyan-300/10 px-3 py-1 text-xs font-medium tracking-[0.18em] text-cyan-100">
            <span class="h-1.5 w-1.5 rounded-full bg-cyan-200 shadow-[0_0_12px_rgba(165,243,252,0.9)] animate-pulse" />
            ENTERPRISE TRACEABILITY SUITE
          </p>

          <h1 class="mt-6 max-w-3xl font-display text-4xl font-semibold leading-[1.01] text-white md:text-6xl">
            把供应链数据
            <span class="text-gradient-sea animate-shimmer">变成可验证的品牌信任资产</span>
          </h1>

          <p class="mt-6 max-w-2xl text-base leading-relaxed text-slate-300/90 md:text-lg">
            SeaTrace Cloud 面向海参产业链提供“批次管理 + 事件上链 + 公开查询”一体化能力。
            平台既能满足企业日常运营效率，也能在验收、路演和客户拜访场景中，清晰展示可信追溯价值。
          </p>

          <div class="mt-8 flex flex-wrap gap-3">
            <RouterLink
              :to="{ name: 'trace', params: { batchNo: 'HSC-2025-DL-PREMIUM-001' } }"
              class="btn-primary px-7"
            >
              查看示例追溯链
            </RouterLink>
            <RouterLink to="/login" class="btn-ghost">进入企业工作台</RouterLink>
          </div>

          <div class="mt-10 grid max-w-2xl gap-4 sm:grid-cols-3">
            <div
              v-for="metric in metrics"
              :key="metric.label"
              class="rounded-2xl border border-white/[0.08] bg-white/[0.04] p-4 transition hover:border-cyan-300/20 hover:bg-white/[0.06]"
            >
              <p class="text-[11px] uppercase tracking-[0.22em] text-slate-500">{{ metric.label }}</p>
              <p class="mt-3 font-display text-3xl text-white">{{ metric.value }}</p>
              <p class="mt-2 text-xs leading-relaxed text-slate-400">{{ metric.note }}</p>
            </div>
          </div>
        </div>

        <div class="relative animate-fade-up [animation-delay:160ms]">
          <div
            class="aurora-border relative rounded-[2rem] bg-slate-950/45 p-5 shadow-card backdrop-blur-2xl"
            :style="{ transform: `perspective(1400px) rotateX(${(50 - pointerY) * 0.05}deg) rotateY(${(pointerX - 50) * 0.06}deg)` }"
          >
            <div class="grid gap-4 md:grid-cols-2">
              <div class="rounded-[1.6rem] border border-white/[0.08] bg-white/[0.04] p-5">
                <p class="text-[11px] uppercase tracking-[0.24em] text-slate-500">Batch Snapshot</p>
                <div class="relative mt-4 overflow-hidden rounded-[1.35rem] border border-white/[0.08] bg-[linear-gradient(180deg,rgba(126,211,255,0.18),rgba(8,47,73,0.55)_55%,rgba(4,10,24,0.9))] p-4">
                  <div class="absolute inset-x-0 bottom-0 h-20 bg-[radial-gradient(circle_at_top,rgba(103,232,249,0.14),transparent_55%)]" />
                  <div class="absolute inset-x-0 top-0 h-20 bg-[linear-gradient(180deg,rgba(255,214,102,0.18),transparent)]" />
                  <div class="relative">
                    <p class="text-[11px] uppercase tracking-[0.22em] text-cyan-100/70">Premium Sea Cucumber</p>
                    <p class="mt-8 text-lg font-semibold text-white">批次号 HSC-2025-DL-PREMIUM-001</p>
                    <p class="mt-2 text-xs leading-relaxed text-slate-300">
                      养殖、采捕、加工、仓配、门店全流程同步留痕，支持按批次一键回溯。
                    </p>
                  </div>
                </div>
              </div>

              <div class="rounded-[1.6rem] border border-white/[0.08] bg-gradient-to-br from-white/[0.05] to-white/[0.02] p-5">
                <p class="text-[11px] uppercase tracking-[0.24em] text-slate-500">Live Integrity</p>
                <div class="mt-4 rounded-2xl border border-white/[0.08] bg-white/[0.03] p-4">
                  <p class="font-mono text-sm text-cyan-100">Anchored: 5 / 5</p>
                  <p class="mt-3 text-xs text-slate-500">关键节点实时登记并展示链上锚定状态，现场可直接验签核对。</p>
                  <div class="mt-4 space-y-3">
                    <div class="h-2 rounded-full bg-white/[0.06]">
                      <div class="h-2 w-[88%] rounded-full bg-gradient-to-r from-cyan-300 to-teal-300 animate-shimmer" />
                    </div>
                    <div class="h-2 rounded-full bg-white/[0.06]">
                      <div class="h-2 w-[74%] rounded-full bg-gradient-to-r from-sky-300/90 to-cyan-200/90 animate-shimmer [animation-delay:1s]" />
                    </div>
                    <div class="h-2 rounded-full bg-white/[0.06]">
                      <div class="h-2 w-[92%] rounded-full bg-gradient-to-r from-emerald-300/90 to-cyan-200/90 animate-shimmer [animation-delay:2s]" />
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="mt-4 rounded-[1.65rem] border border-white/[0.08] bg-white/[0.04] p-5">
              <div class="flex items-center justify-between gap-4">
                <div>
                  <p class="text-[11px] uppercase tracking-[0.24em] text-slate-500">Product Narrative</p>
                  <p class="mt-2 text-base font-semibold text-white">从数据采集到品牌背书，一条链讲清价值闭环</p>
                </div>
                <span class="rounded-full bg-emerald-400/15 px-3 py-1 text-[11px] font-medium text-emerald-200 ring-1 ring-emerald-400/20">
                  Demo Ready
                </span>
              </div>

              <div class="mt-5 grid gap-3">
                <div
                  v-for="(scene, index) in scenes"
                  :key="scene.title"
                  class="spotlight-card relative overflow-hidden rounded-2xl border border-white/[0.08] bg-white/[0.03] p-4 transition hover:-translate-y-0.5 hover:border-cyan-300/25 hover:bg-white/[0.05]"
                >
                  <div class="absolute inset-y-0 left-0 w-px bg-gradient-to-b from-transparent via-cyan-300/70 to-transparent" />
                  <div class="flex gap-4">
                    <div class="mt-0.5 flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl bg-cyan-300/12 text-sm font-semibold text-cyan-100">
                      0{{ index + 1 }}
                    </div>
                    <div>
                      <p class="text-[11px] uppercase tracking-[0.24em] text-slate-500">{{ scene.eyebrow }}</p>
                      <p class="mt-1 text-sm font-semibold text-white">{{ scene.title }}</p>
                      <p class="mt-1 text-xs leading-relaxed text-slate-400">{{ scene.body }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-16 grid gap-6 lg:grid-cols-[0.9fr_1.1fr]">
      <div class="glass rounded-[1.8rem] p-6 shadow-card ring-1 ring-white/[0.05] md:p-8">
        <p class="text-xs uppercase tracking-[0.28em] text-slate-500">Why It Works</p>
        <h2 class="mt-3 font-display text-2xl font-semibold text-white">既能对内运营，也能对外展示</h2>
        <p class="mt-4 text-sm leading-relaxed text-slate-400">
          这版首页重点强化了品牌化表达与价值可视化，不再只是技术说明页。用于验收阶段时，
          能先说明产品定位，再落到业务流程和可信机制，叙事更完整。
        </p>
      </div>

      <div class="grid gap-4 md:grid-cols-3">
        <div class="glass rounded-[1.5rem] p-5 ring-1 ring-white/[0.05]">
          <p class="text-[11px] uppercase tracking-[0.24em] text-slate-500">产品定位</p>
          <p class="mt-3 text-sm leading-relaxed text-slate-300">突出“企业级可信追溯平台”，弱化课程作业表达。</p>
        </div>
        <div class="glass rounded-[1.5rem] p-5 ring-1 ring-white/[0.05]">
          <p class="text-[11px] uppercase tracking-[0.24em] text-slate-500">业务价值</p>
          <p class="mt-3 text-sm leading-relaxed text-slate-300">强调协同效率、合规核验、消费者信任三个核心收益。</p>
        </div>
        <div class="glass rounded-[1.5rem] p-5 ring-1 ring-white/[0.05]">
          <p class="text-[11px] uppercase tracking-[0.24em] text-slate-500">验收友好</p>
          <p class="mt-3 text-sm leading-relaxed text-slate-300">关键按钮、示例批次与状态数据可直接用于现场演示。</p>
        </div>
      </div>
    </div>
  </section>
</template>
