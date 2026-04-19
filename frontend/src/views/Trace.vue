<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import api from "../api/client";

const STAGE_LABEL: Record<string, string> = {
  breeding: "养殖",
  harvest: "采捕",
  processing: "加工",
  packaging: "包装",
  logistics: "物流",
  retail: "零售",
};

const route = useRoute();
const router = useRouter();
const batchNoInput = ref("");
const loading = ref(false);
const error = ref("");
const copiedEventId = ref<number | null>(null);
const batch = ref<Record<string, unknown> | null>(null);
const items = ref<
  {
    event: Record<string, unknown>;
    chain?: { txId?: string; status?: string; blockNumber?: number };
  }[]
>([]);

function formatBlockHeight(n: unknown): string {
  if (typeof n !== "number" || !Number.isFinite(n)) return "";
  return `#${n.toLocaleString("en-US")}`;
}

function stageLabel(stage: unknown): string {
  const s = String(stage);
  return STAGE_LABEL[s] || s;
}

function orgName(b: Record<string, unknown>): string {
  const o = b.org as { name?: string } | undefined;
  return o?.name ?? "";
}

const batchNo = computed(() => (route.params.batchNo as string) || "");

async function load(no: string) {
  if (!no) return;
  loading.value = true;
  error.value = "";
  try {
    const { data } = await api.get(`/trace/${encodeURIComponent(no)}`);
    batch.value = data.batch;
    items.value = data.items || [];
  } catch {
    error.value = "未找到该批次，或无法连接后端（请确认服务已启动）。";
    batch.value = null;
    items.value = [];
  } finally {
    loading.value = false;
  }
}

function search() {
  const q = batchNoInput.value.trim();
  if (!q) return;
  router.push({ name: "trace", params: { batchNo: q } });
}

async function copyHash(eventId: number, hash: unknown) {
  const h = String(hash);
  try {
    await navigator.clipboard.writeText(h);
    copiedEventId.value = eventId;
    window.setTimeout(() => {
      if (copiedEventId.value === eventId) copiedEventId.value = null;
    }, 2000);
  } catch {
    /* ignore */
  }
}

onMounted(() => {
  batchNoInput.value = batchNo.value || "HSC-2025-DL-PREMIUM-001";
  if (batchNo.value) load(batchNo.value);
});

watch(
  () => route.params.batchNo,
  (v) => {
    if (typeof v === "string" && v) load(v);
  },
);
</script>

<template>
  <section class="mx-auto max-w-4xl px-6 py-12 md:py-16">
    <div class="text-center animate-fade-up">
      <h1 class="font-display text-3xl font-semibold tracking-tight text-white md:text-4xl">
        公开溯源查询
      </h1>
      <p class="mx-auto mt-3 max-w-lg text-sm leading-relaxed text-slate-400">
        输入批次号，查看链下时间线与链上锚定摘要（交易号、区块高度）。
      </p>
    </div>

    <div
      class="mx-auto mt-10 flex max-w-xl flex-col gap-3 sm:flex-row sm:items-center sm:gap-2"
    >
      <input
        v-model="batchNoInput"
        class="input-field rounded-2xl sm:flex-1"
        placeholder="例如 HSC-2025-DL-PREMIUM-001"
        @keyup.enter="search"
      />
      <button type="button" class="btn-primary shrink-0 px-8" @click="search">查询</button>
    </div>

    <div v-if="loading" class="mt-12 space-y-4">
      <div class="glass noise rounded-3xl p-6">
        <div class="skeleton mb-4 h-6 w-1/2" />
        <div class="skeleton mb-2 h-4 w-full" />
        <div class="skeleton h-4 w-2/3" />
      </div>
      <div class="skeleton h-24 rounded-2xl" />
    </div>

    <p
      v-else-if="error"
      class="mt-10 rounded-2xl border border-rose-500/20 bg-rose-500/10 px-4 py-3 text-center text-sm text-rose-100"
    >
      {{ error }}
    </p>

    <div v-else-if="batch" class="mt-12 space-y-8">
      <div
        class="glass noise relative overflow-hidden rounded-[1.75rem] p-6 shadow-card ring-1 ring-white/[0.06] md:p-8"
      >
        <div
          class="pointer-events-none absolute -right-20 -top-20 h-56 w-56 rounded-full bg-cyan-500/10 blur-3xl"
        />
        <div class="relative flex flex-wrap items-start justify-between gap-4">
          <div>
            <p class="text-[11px] font-medium uppercase tracking-[0.25em] text-slate-500">批次号</p>
            <p class="mt-1 font-mono text-xl text-cyan-100 md:text-2xl">{{ batch.batchNo }}</p>
            <p v-if="batch.productName" class="mt-2 text-base font-medium text-white">
              {{ batch.productName }}
            </p>
          </div>
          <div class="text-right text-xs text-slate-400">
            <p class="text-[11px] uppercase tracking-wider text-slate-500">养殖主体</p>
            <p class="mt-1 text-sm text-slate-200">{{ orgName(batch) }}</p>
          </div>
        </div>
        <dl class="relative mt-6 grid gap-4 text-sm sm:grid-cols-2 lg:grid-cols-3">
          <div v-if="batch.farmBase" class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
            <dt class="text-xs text-slate-500">养殖基地</dt>
            <dd class="mt-0.5 text-slate-100">{{ batch.farmBase }}</dd>
          </div>
          <div v-if="batch.catchDate" class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
            <dt class="text-xs text-slate-500">采捕日期</dt>
            <dd class="mt-0.5 text-slate-100">
              {{ new Date(String(batch.catchDate)).toLocaleDateString() }}
            </dd>
          </div>
          <div v-if="batch.quality" class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
            <dt class="text-xs text-slate-500">质量等级</dt>
            <dd class="mt-0.5 text-slate-100">{{ batch.quality }}</dd>
          </div>
          <div class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
            <dt class="text-xs text-slate-500">海域</dt>
            <dd class="mt-0.5 text-slate-100">{{ batch.breedArea }}</dd>
          </div>
          <div class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
            <dt class="text-xs text-slate-500">规格</dt>
            <dd class="mt-0.5 text-slate-100">{{ batch.spec }}</dd>
          </div>
          <div class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
            <dt class="text-xs text-slate-500">数量</dt>
            <dd class="mt-0.5 text-slate-100">{{ batch.quantity }}</dd>
          </div>
        </dl>
      </div>

      <div class="space-y-5">
        <h2 class="text-base font-semibold text-white">溯源时间线</h2>
        <p v-if="!items.length" class="text-sm text-slate-500">
          暂无事件。企业用户可在工作台通过 API 登记事件后刷新查看。
        </p>
        <ol v-else class="relative space-y-0 border-l border-white/10 pl-6">
          <li v-for="(it, idx) in items" :key="idx" class="relative pb-10 last:pb-0">
            <span
              class="absolute -left-[29px] top-1.5 flex h-[11px] w-[11px] rounded-full border-2 border-abyss-950 bg-cyan-400 shadow-[0_0_14px_rgba(34,211,238,0.85)]"
            />
            <div class="glass rounded-2xl p-5 ring-1 ring-white/[0.05]">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <p class="text-sm font-medium text-white">{{ it.event.title }}</p>
                <span
                  class="rounded-full bg-white/5 px-2.5 py-0.5 text-[11px] text-cyan-100/90 ring-1 ring-white/10"
                >
                  {{ stageLabel(it.event.stage) }}
                </span>
              </div>
              <p class="mt-1.5 text-xs text-slate-500">
                {{ new Date(String(it.event.occurredAt)).toLocaleString() }} · {{ it.event.location }}
              </p>
              <div class="mt-3 flex flex-wrap items-center gap-2">
                <p class="text-xs text-slate-400">
                  数据哈希
                  <span class="ml-1 font-mono text-[11px] text-cyan-100/90">{{ String(it.event.dataHash).slice(0, 20) }}…</span>
                </p>
                <button
                  type="button"
                  class="rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-[11px] text-slate-300 hover:border-cyan-400/30 hover:text-white"
                  @click="copyHash(Number(it.event.id), it.event.dataHash)"
                >
                  {{ copiedEventId === Number(it.event.id) ? "已复制" : "复制全文" }}
                </button>
              </div>
              <p v-if="it.chain?.txId" class="mt-3 text-xs text-emerald-200/95">
                链上交易 · <span class="break-all font-mono text-[11px]">{{ it.chain.txId }}</span>
                <span v-if="it.chain.status" class="text-slate-500">（{{ it.chain.status }}）</span>
              </p>
              <p
                v-if="it.chain && formatBlockHeight(it.chain.blockNumber)"
                class="mt-1.5 text-xs font-medium text-cyan-200/95"
              >
                区块高度 · {{ formatBlockHeight(it.chain.blockNumber) }}
              </p>
            </div>
          </li>
        </ol>
      </div>
    </div>
  </section>
</template>
