<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import api from "../api/client";
import {
  type DetailView,
  buildStructuredDetailView,
  stageLabel,
} from "../lib/traceSchemas";

type BatchData = {
  batchNo: string;
  productName?: string;
  farmBase?: string;
  catchDate?: string;
  quality?: string;
  breedArea?: string;
  spec?: string;
  quantity?: string;
  extraJson?: string;
  org?: { name?: string };
};

type TimelineApiItem = {
  event: {
    id: number;
    stage: string;
    title: string;
    detailJson?: string;
    location?: string;
    operatorName?: string;
    evidenceUrls?: string;
    occurredAt: string;
    dataHash: string;
  };
  chain?: { txId?: string; status?: string; blockNumber?: number };
};

type TimelineItem = TimelineApiItem & {
  detailView: DetailView | null;
};

const route = useRoute();
const router = useRouter();
const batchNoInput = ref("");
const loading = ref(false);
const error = ref("");
const copiedEventId = ref<number | null>(null);
const batch = ref<BatchData | null>(null);
const items = ref<TimelineItem[]>([]);

const batchNo = computed(() => (route.params.batchNo as string) || "");
const anchoredCount = computed(() => items.value.filter((it) => !!it.chain?.txId).length);
const stageCount = computed(() => new Set(items.value.map((it) => String(it.event.stage))).size);
const batchDetailView = computed(() => buildStructuredDetailView(batch.value?.extraJson));

function enrichTimelineItem(item: TimelineApiItem): TimelineItem {
  return {
    ...item,
    detailView: buildStructuredDetailView(item.event.detailJson, item.event.evidenceUrls),
  };
}

function formatBlockHeight(n: unknown): string {
  if (typeof n !== "number" || !Number.isFinite(n)) return "";
  return `#${n.toLocaleString("en-US")}`;
}

function orgName(record: BatchData): string {
  return record.org?.name ?? "";
}

async function load(no: string) {
  if (!no) return;
  loading.value = true;
  error.value = "";
  try {
    const { data } = await api.get(`/trace/${encodeURIComponent(no)}`);
    batch.value = data.batch || null;
    items.value = (data.items || []).map(enrichTimelineItem);
  } catch {
    error.value = "未找到该批次，或后端服务不可用。";
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
  const text = String(hash);
  try {
    await navigator.clipboard.writeText(text);
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
  (value) => {
    if (typeof value === "string" && value) {
      batchNoInput.value = value;
      load(value);
    }
  },
);
</script>

<template>
  <section class="mx-auto max-w-5xl px-6 py-12 md:py-16">
    <div
      class="relative overflow-hidden rounded-[2.2rem] border border-white/[0.06] bg-white/[0.03] px-6 py-10 shadow-card backdrop-blur-xl md:px-8 md:py-12"
    >
      <div class="absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(34,211,238,0.16),transparent_30%),linear-gradient(180deg,rgba(255,255,255,0.03),transparent_30%,rgba(8,47,73,0.22))]" />
      <div class="absolute right-[-4rem] top-[-3rem] h-56 w-56 rounded-full bg-cyan-400/10 blur-3xl animate-drift" />
      <div class="absolute left-[-3rem] top-[35%] h-48 w-48 rounded-full bg-teal-300/10 blur-3xl animate-float" />
      <div class="noise absolute inset-0 opacity-35" />

      <div class="relative text-center animate-fade-up">
        <p class="text-xs uppercase tracking-[0.28em] text-cyan-100/70">Public Trace Portal</p>
        <h1 class="mt-4 font-display text-3xl font-semibold tracking-tight text-white md:text-5xl">
          一键查询批次全链路追溯信息
        </h1>
        <p class="mx-auto mt-4 max-w-2xl text-sm leading-relaxed text-slate-300/90 md:text-base">
          输入批次号后，可查看批次档案、标准化事件明细、链上交易摘要与数据哈希。页面可直接用于验收、渠道演示和消费者查询。
        </p>
      </div>

      <div class="relative mx-auto mt-10 flex max-w-2xl flex-col gap-3 sm:flex-row sm:items-center">
        <input
          v-model="batchNoInput"
          class="input-field rounded-2xl sm:flex-1"
          placeholder="例如 HSC-2025-DL-PREMIUM-001"
          @keyup.enter="search"
        />
        <button type="button" class="btn-primary shrink-0 px-8" @click="search">查询批次</button>
      </div>

      <div class="relative mt-8 grid gap-4 sm:grid-cols-3">
        <div class="rounded-2xl border border-white/[0.08] bg-white/[0.04] p-4 backdrop-blur-xl">
          <p class="text-[11px] uppercase tracking-[0.22em] text-slate-500">链上锚定</p>
          <p class="mt-2 font-display text-2xl text-white">{{ anchoredCount }}</p>
          <p class="mt-1 text-xs text-slate-400">包含交易摘要的事件节点数</p>
        </div>
        <div class="rounded-2xl border border-white/[0.08] bg-white/[0.04] p-4 backdrop-blur-xl">
          <p class="text-[11px] uppercase tracking-[0.22em] text-slate-500">覆盖阶段</p>
          <p class="mt-2 font-display text-2xl text-white">{{ stageCount }}</p>
          <p class="mt-1 text-xs text-slate-400">已登记的业务阶段总数</p>
        </div>
        <div class="rounded-2xl border border-white/[0.08] bg-white/[0.04] p-4 backdrop-blur-xl">
          <p class="text-[11px] uppercase tracking-[0.22em] text-slate-500">查询模式</p>
          <p class="mt-2 font-display text-2xl text-white">公开</p>
          <p class="mt-1 text-xs text-slate-400">无需登录即可查看标准化追溯结果</p>
        </div>
      </div>
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
      <div class="grid gap-5 lg:grid-cols-[1.05fr_0.95fr]">
        <div class="glass noise relative overflow-hidden rounded-[1.75rem] p-6 shadow-card ring-1 ring-white/[0.06] md:p-8">
          <div class="pointer-events-none absolute -right-20 -top-20 h-56 w-56 rounded-full bg-cyan-500/10 blur-3xl animate-drift" />
          <div class="relative flex flex-wrap items-start justify-between gap-4">
            <div>
              <p class="text-[11px] font-medium uppercase tracking-[0.25em] text-slate-500">批次号</p>
              <p class="mt-1 font-mono text-xl text-cyan-100 md:text-2xl">{{ batch.batchNo }}</p>
              <p v-if="batch.productName" class="mt-2 text-base font-medium text-white">{{ batch.productName }}</p>
            </div>
            <div class="text-right text-xs text-slate-400">
              <p class="text-[11px] uppercase tracking-wider text-slate-500">主体机构</p>
              <p class="mt-1 text-sm text-slate-200">{{ orgName(batch) || "-" }}</p>
            </div>
          </div>
          <dl class="relative mt-6 grid gap-4 text-sm sm:grid-cols-2 lg:grid-cols-3">
            <div v-if="batch.farmBase" class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
              <dt class="text-xs text-slate-500">养殖基地</dt>
              <dd class="mt-0.5 text-slate-100">{{ batch.farmBase }}</dd>
            </div>
            <div v-if="batch.catchDate" class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
              <dt class="text-xs text-slate-500">采捕日期</dt>
              <dd class="mt-0.5 text-slate-100">{{ new Date(String(batch.catchDate)).toLocaleDateString() }}</dd>
            </div>
            <div v-if="batch.quality" class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
              <dt class="text-xs text-slate-500">质量等级</dt>
              <dd class="mt-0.5 text-slate-100">{{ batch.quality }}</dd>
            </div>
            <div class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
              <dt class="text-xs text-slate-500">养殖海域</dt>
              <dd class="mt-0.5 text-slate-100">{{ batch.breedArea || "-" }}</dd>
            </div>
            <div class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
              <dt class="text-xs text-slate-500">产品规格</dt>
              <dd class="mt-0.5 text-slate-100">{{ batch.spec || "-" }}</dd>
            </div>
            <div class="rounded-xl bg-white/[0.03] px-3 py-2 ring-1 ring-white/[0.06]">
              <dt class="text-xs text-slate-500">数量</dt>
              <dd class="mt-0.5 text-slate-100">{{ batch.quantity || "-" }}</dd>
            </div>
          </dl>
        </div>

        <div class="glass relative overflow-hidden rounded-[1.75rem] p-6 shadow-card ring-1 ring-white/[0.06]">
          <div class="relative">
            <p class="text-sm font-semibold text-white">档案与可信摘要</p>
            <div class="mt-5 space-y-3">
              <div class="rounded-2xl border border-white/[0.08] bg-white/[0.04] p-4">
                <p class="text-[11px] uppercase tracking-wider text-slate-500">时间线节点</p>
                <p class="mt-2 font-display text-2xl text-white">{{ items.length }}</p>
              </div>
              <div class="rounded-2xl border border-white/[0.08] bg-white/[0.04] p-4">
                <p class="text-[11px] uppercase tracking-wider text-slate-500">链上交易</p>
                <p class="mt-2 font-display text-2xl text-white">{{ anchoredCount }}</p>
              </div>
              <div v-if="batchDetailView" class="rounded-2xl border border-white/[0.08] bg-white/[0.04] p-4">
                <p class="text-[11px] uppercase tracking-wider text-slate-500">扩展档案字段</p>
                <div class="mt-3 space-y-3">
                  <div
                    v-for="section in batchDetailView.sections"
                    :key="section.title"
                    class="rounded-xl border border-white/[0.06] bg-white/[0.03] p-3"
                  >
                    <p class="text-xs font-medium uppercase tracking-[0.16em] text-slate-500">{{ section.title }}</p>
                    <dl class="mt-2 space-y-2">
                      <div v-for="detail in section.items" :key="detail.label" class="grid gap-1 text-sm">
                        <dt class="text-xs text-slate-500">{{ detail.label }}</dt>
                        <dd class="break-all text-slate-200">{{ detail.value }}</dd>
                      </div>
                    </dl>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="space-y-5">
        <div class="flex items-baseline justify-between gap-4">
          <h2 class="text-base font-semibold text-white">追溯时间线</h2>
          <span class="text-xs text-slate-500">{{ items.length }} 个登记事件</span>
        </div>
        <p v-if="!items.length" class="text-sm text-slate-500">当前批次暂无事件记录。</p>
        <ol v-else class="relative space-y-0 border-l border-white/10 pl-6">
          <li v-for="(item, index) in items" :key="index" class="relative pb-10 last:pb-0">
            <span
              class="absolute -left-[29px] top-1.5 flex h-[11px] w-[11px] rounded-full border-2 border-abyss-950 bg-cyan-400 shadow-[0_0_14px_rgba(34,211,238,0.85)]"
            />
            <div class="spotlight-card glass rounded-2xl p-5 ring-1 ring-white/[0.05] transition hover:-translate-y-0.5 hover:ring-cyan-300/20">
              <div class="flex flex-wrap items-center justify-between gap-2">
                <p class="text-sm font-medium text-white">{{ item.event.title }}</p>
                <span class="rounded-full bg-white/5 px-2.5 py-0.5 text-[11px] text-cyan-100/90 ring-1 ring-white/10">
                  {{ stageLabel(item.event.stage) }}
                </span>
              </div>
              <p class="mt-1.5 text-xs text-slate-500">
                {{ new Date(String(item.event.occurredAt)).toLocaleString() }}
                <span v-if="item.event.location"> | {{ item.event.location }}</span>
                <span v-if="item.event.operatorName"> | {{ item.event.operatorName }}</span>
              </p>

              <div v-if="item.detailView" class="mt-4 space-y-4">
                <div v-if="item.detailView.summary.length" class="flex flex-wrap gap-2">
                  <span
                    v-for="summary in item.detailView.summary"
                    :key="`${item.event.id}-${summary.label}`"
                    class="rounded-full bg-cyan-400/10 px-3 py-1 text-[11px] text-cyan-100 ring-1 ring-cyan-400/20"
                  >
                    {{ summary.label }}：{{ summary.value }}
                  </span>
                </div>

                <div class="grid gap-3 lg:grid-cols-2">
                  <div
                    v-for="section in item.detailView.sections"
                    :key="`${item.event.id}-${section.title}`"
                    class="rounded-2xl border border-white/[0.08] bg-white/[0.03] p-3"
                  >
                    <p class="text-xs font-medium uppercase tracking-[0.16em] text-slate-500">{{ section.title }}</p>
                    <dl class="mt-3 space-y-2">
                      <div v-for="detail in section.items" :key="detail.label" class="grid gap-1 text-sm">
                        <dt class="text-xs text-slate-500">{{ detail.label }}</dt>
                        <dd class="break-all text-slate-200">{{ detail.value }}</dd>
                      </div>
                    </dl>
                  </div>
                </div>
              </div>

              <div class="mt-4 flex flex-wrap items-center gap-2">
                <p class="text-xs text-slate-400">
                  数据哈希
                  <span class="ml-1 font-mono text-[11px] text-cyan-100/90">{{ String(item.event.dataHash).slice(0, 20) }}...</span>
                </p>
                <button
                  type="button"
                  class="rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-[11px] text-slate-300 hover:border-cyan-400/30 hover:text-white"
                  @click="copyHash(Number(item.event.id), item.event.dataHash)"
                >
                  {{ copiedEventId === Number(item.event.id) ? "已复制" : "复制完整哈希" }}
                </button>
              </div>

              <p v-if="item.chain?.txId" class="mt-3 text-xs text-emerald-200/95">
                链上交易 | <span class="break-all font-mono text-[11px]">{{ item.chain.txId }}</span>
                <span v-if="item.chain.status" class="text-slate-500">（{{ item.chain.status }}）</span>
              </p>
              <p v-if="item.chain && formatBlockHeight(item.chain.blockNumber)" class="mt-1.5 text-xs font-medium text-cyan-200/95">
                区块高度 | {{ formatBlockHeight(item.chain.blockNumber) }}
              </p>
            </div>
          </li>
        </ol>
      </div>
    </div>
  </section>
</template>
