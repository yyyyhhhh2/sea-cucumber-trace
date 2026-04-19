<script setup lang="ts">
import { onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import axios from "axios";
import api from "../api/client";

type Batch = {
  id: number;
  batchNo: string;
  breedArea: string;
  spec: string;
  quantity: string;
  productName?: string;
  farmBase?: string;
  org?: { name: string };
};

const batches = ref<Batch[]>([]);
const loading = ref(true);
const err = ref("");
const ok = ref("");

const form = ref({
  batchNo: "",
  orgId: 1,
  productName: "",
  farmBase: "",
  quality: "",
  catchDate: "",
  breedArea: "",
  spec: "",
  quantity: "",
});

async function load() {
  loading.value = true;
  err.value = "";
  try {
    const { data } = await api.get("/batches");
    batches.value = data.items || [];
  } catch {
    err.value = "加载失败，请确认已登录且后端已启动。";
  } finally {
    loading.value = false;
  }
}

async function createBatch() {
  ok.value = "";
  err.value = "";
  if (!form.value.batchNo.trim()) {
    err.value = "请填写批次号";
    return;
  }
  try {
    const body: Record<string, unknown> = {
      batchNo: form.value.batchNo.trim(),
      orgId: form.value.orgId,
      breedArea: form.value.breedArea.trim(),
      spec: form.value.spec.trim(),
      quantity: form.value.quantity.trim(),
    };
    if (form.value.productName.trim()) body.productName = form.value.productName.trim();
    if (form.value.farmBase.trim()) body.farmBase = form.value.farmBase.trim();
    if (form.value.quality.trim()) body.quality = form.value.quality.trim();
    if (form.value.catchDate) {
      body.catchDate = new Date(`${form.value.catchDate}T12:00:00`).toISOString();
    }
    await api.post("/batches", body);
    form.value = {
      batchNo: "",
      orgId: form.value.orgId,
      productName: "",
      farmBase: "",
      quality: "",
      catchDate: "",
      breedArea: "",
      spec: "",
      quantity: "",
    };
    ok.value = "批次已创建";
    await load();
  } catch (e: unknown) {
    if (axios.isAxiosError(e) && typeof e.message === "string") {
      err.value = e.message;
    } else {
      err.value = "创建失败（权限或批次号重复）";
    }
  }
}

onMounted(load);
</script>

<template>
  <section class="mx-auto max-w-6xl px-6 py-10 md:py-12">
    <div class="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
      <div>
        <h1 class="font-display text-3xl font-semibold tracking-tight text-white">企业工作台</h1>
        <p class="mt-1.5 max-w-xl text-sm text-slate-400">
          管理批次并在时间线上登记溯源事件；管理员可为任意企业创建批次。
        </p>
      </div>
      <button
        type="button"
        class="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-200 transition hover:border-cyan-400/35 hover:bg-white/10"
        @click="load"
      >
        刷新列表
      </button>
    </div>

    <div class="mt-10 grid gap-8 lg:grid-cols-[1.15fr_0.85fr]">
      <div class="glass noise rounded-[1.75rem] p-6 shadow-card ring-1 ring-white/[0.05] md:p-8">
        <h2 class="text-base font-semibold text-white">新建批次</h2>
        <p class="mt-1 text-xs text-slate-500">带 <span class="text-slate-400">*</span> 为必填，其余选填。</p>

        <div class="mt-6 grid gap-4">
          <div class="grid gap-4 md:grid-cols-2">
            <div class="md:col-span-2">
              <label class="mb-1.5 block text-xs font-medium text-slate-400"
                >批次号 <span class="text-rose-300/90">*</span></label
              >
              <input
                v-model="form.batchNo"
                class="input-field font-mono"
                placeholder="如 HSC-2026-0001"
              />
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-400">企业 orgId</label>
              <input
                v-model.number="form.orgId"
                type="number"
                min="1"
                class="input-field"
              />
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-400">采捕日期</label>
              <input v-model="form.catchDate" type="date" class="input-field" />
            </div>
            <div class="md:col-span-2">
              <label class="mb-1.5 block text-xs font-medium text-slate-400">产品名称</label>
              <input
                v-model="form.productName"
                class="input-field"
                placeholder="如 大连刺参 (精品级)"
              />
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-400">养殖基地</label>
              <input
                v-model="form.farmBase"
                class="input-field"
                placeholder="养殖场名称"
              />
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-400">质量等级</label>
              <input
                v-model="form.quality"
                class="input-field"
                placeholder="如 合格 — 国家A级标准"
              />
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-400">养殖海域</label>
              <input v-model="form.breedArea" class="input-field" placeholder="海域或分区" />
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-400">规格</label>
              <input v-model="form.spec" class="input-field" placeholder="规格 / 等级" />
            </div>
            <div class="md:col-span-2">
              <label class="mb-1.5 block text-xs font-medium text-slate-400">数量</label>
              <input v-model="form.quantity" class="input-field" placeholder="如 500 kg" />
            </div>
          </div>
        </div>

        <button type="button" class="btn-primary mt-6 w-full" @click="createBatch">
          创建批次
        </button>
        <p
          v-if="ok"
          class="mt-4 rounded-xl border border-emerald-400/25 bg-emerald-400/10 px-3 py-2 text-sm text-emerald-100"
        >
          {{ ok }}
        </p>
        <p
          v-if="err"
          class="mt-4 rounded-xl border border-rose-500/25 bg-rose-500/10 px-3 py-2 text-sm text-rose-100"
        >
          {{ err }}
        </p>
      </div>

      <div class="glass noise h-fit rounded-[1.75rem] p-6 shadow-card ring-1 ring-white/[0.05] md:p-7">
        <h2 class="text-base font-semibold text-white">使用提示</h2>
        <ul class="mt-4 space-y-3 text-sm leading-relaxed text-slate-400">
          <li class="flex gap-2">
            <span class="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-cyan-400/80" />
            <span>企业用户仅能使用本企业 <code class="rounded bg-white/10 px-1 font-mono text-xs">orgId</code>。</span>
          </li>
          <li class="flex gap-2">
            <span class="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-cyan-400/80" />
            <span>登记溯源事件可通过 API <code class="rounded bg-white/10 px-1 font-mono text-[11px]">POST /batches/:id/events</code>；哈希将自动上链（演示环境为 mock）。</span>
          </li>
          <li class="flex gap-2">
            <span class="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-cyan-400/80" />
            <span>后端健康检查：<code class="rounded bg-white/10 px-1 font-mono text-[11px]">GET /api/health</code></span>
          </li>
        </ul>
      </div>
    </div>

    <div class="mt-14">
      <div class="flex items-baseline justify-between gap-4">
        <h2 class="text-lg font-semibold text-white">批次列表</h2>
        <span class="text-xs text-slate-500">{{ batches.length }} 条</span>
      </div>

      <div v-if="loading" class="mt-6 grid gap-4 md:grid-cols-2">
        <div v-for="i in 4" :key="i" class="glass rounded-2xl p-5">
          <div class="skeleton mb-3 h-4 w-2/3" />
          <div class="skeleton h-3 w-1/3" />
          <div class="mt-4 grid grid-cols-2 gap-2">
            <div class="skeleton h-8" />
            <div class="skeleton h-8" />
          </div>
        </div>
      </div>

      <div
        v-else-if="!batches.length"
        class="mt-6 rounded-2xl border border-dashed border-white/10 bg-white/[0.02] px-6 py-14 text-center text-sm text-slate-500"
      >
        暂无批次。可在上方创建，或使用公开溯源页查看演示批次。
      </div>

      <div v-else class="mt-6 grid gap-4 md:grid-cols-2">
        <article
          v-for="b in batches"
          :key="b.id"
          class="group glass rounded-2xl p-5 transition hover:border-cyan-400/25 hover:shadow-card"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <p class="truncate font-mono text-sm text-cyan-200">{{ b.batchNo }}</p>
              <p v-if="b.productName" class="mt-1 truncate text-sm text-white/90">
                {{ b.productName }}
              </p>
              <p class="mt-1 text-xs text-slate-500">{{ b.org?.name }}</p>
            </div>
            <RouterLink
              class="shrink-0 rounded-full bg-white/5 px-3 py-1.5 text-xs font-medium text-white ring-1 ring-white/10 transition hover:bg-cyan-400/15 hover:ring-cyan-400/30"
              :to="{ name: 'trace', params: { batchNo: b.batchNo } }"
            >
              溯源
            </RouterLink>
          </div>
          <dl class="mt-4 grid grid-cols-2 gap-3 text-xs">
            <div v-if="b.farmBase">
              <dt class="text-slate-500">基地</dt>
              <dd class="text-slate-200">{{ b.farmBase }}</dd>
            </div>
            <div>
              <dt class="text-slate-500">海域</dt>
              <dd class="text-slate-200">{{ b.breedArea || "—" }}</dd>
            </div>
            <div>
              <dt class="text-slate-500">规格</dt>
              <dd class="text-slate-200">{{ b.spec || "—" }}</dd>
            </div>
          </dl>
        </article>
      </div>
    </div>
  </section>
</template>
