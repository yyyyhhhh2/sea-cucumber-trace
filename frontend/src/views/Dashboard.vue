<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { RouterLink } from "vue-router";
import axios from "axios";
import api from "../api/client";
import { useAuthStore } from "../stores/auth";
import {
  type BatchExtraForm,
  type DetailView,
  type TraceStageValue,
  buildBatchExtraJson,
  buildEventCode,
  buildEventDetailJson,
  buildStructuredDetailView,
  createEmptyBatchExtraForm,
  createEventMetaForm,
  createStageDetailForm,
  defaultSopCode,
  eventStatusOptions,
  parseEvidenceUrls,
  qualityLevelOptions,
  stageDefinitionMap,
  stageDefinitions,
  stageLabel,
  traceTagOptions,
  unitOptions,
} from "../lib/traceSchemas";

type Batch = {
  id: number;
  batchNo: string;
  orgId: number;
  productName?: string;
  farmBase?: string;
  quality?: string;
  catchDate?: string;
  breedArea?: string;
  breedStartDate?: string;
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
  chain?: {
    txId?: string;
    status?: string;
    blockNumber?: number;
  };
};

type TimelineItem = TimelineApiItem & {
  detailView: DetailView | null;
};

type ChecklistItem = {
  label: string;
  filled: boolean;
};

const auth = useAuthStore();
const batches = ref<Batch[]>([]);
const selectedBatchId = ref<number | null>(null);
const timeline = ref<TimelineItem[]>([]);
const loading = ref(true);
const timelineLoading = ref(false);
const savingBatch = ref(false);
const savingEvent = ref(false);
const error = ref("");
const notice = ref("");
const showAdvancedMeta = ref(false);

const currentBatch = computed(() => batches.value.find((batch) => batch.id === selectedBatchId.value) || null);
const currentStageDefinition = computed(
  () => stageDefinitionMap[eventForm.value.stage as TraceStageValue] || stageDefinitions[0],
);

const batchForm = ref({
  batchNo: "",
  orgId: auth.user?.orgId ?? 1,
  productName: "",
  farmBase: "",
  quality: qualityLevelOptions[0].value,
  catchDate: "",
  breedArea: "",
  breedStartDate: "",
  spec: "",
  quantity: "",
});

const batchExtraForm = ref<BatchExtraForm>(createEmptyBatchExtraForm());

const eventForm = ref({
  stage: "breeding",
  title: stageDefinitionMap.breeding.defaultTitle,
  occurredAt: toDatetimeLocal(new Date()),
  location: "",
  operatorName: auth.user?.displayName || "",
});

const eventMetaForm = ref(createEventMetaForm("breeding"));
const eventDetailForm = ref<Record<string, string>>(createStageDetailForm("breeding"));
const evidenceInput = ref("");

const timelineStageCounts = computed(() => {
  const counts: Partial<Record<TraceStageValue, number>> = {};
  for (const item of timeline.value) {
    const stage = item.event.stage as TraceStageValue;
    counts[stage] = (counts[stage] || 0) + 1;
  }
  return counts;
});

const recommendedStage = computed(
  () => stageDefinitions.find((stage) => !timelineStageCounts.value[stage.value]) || null,
);

const stageProgress = computed(() =>
  stageDefinitions.map((stage, index) => ({
    value: stage.value,
    label: stage.label,
    step: index + 1,
    count: timelineStageCounts.value[stage.value] || 0,
    completed: Boolean(timelineStageCounts.value[stage.value]),
    isCurrent: eventForm.value.stage === stage.value,
    isRecommended: recommendedStage.value?.value === stage.value,
  })),
);

const completedStageCount = computed(() => stageProgress.value.filter((item) => item.completed).length);

const batchChecklist = computed<ChecklistItem[]>(() => [
  { label: "批次号", filled: Boolean(batchForm.value.batchNo.trim()) },
  { label: "产品名称", filled: Boolean(batchForm.value.productName.trim()) },
  { label: "养殖起始日期", filled: Boolean(batchForm.value.breedStartDate) },
  { label: "采捕日期", filled: Boolean(batchForm.value.catchDate) },
  { label: "养殖基地", filled: Boolean(batchForm.value.farmBase.trim()) },
  { label: "养殖海域", filled: Boolean(batchForm.value.breedArea.trim()) },
  { label: "产品规格", filled: Boolean(batchForm.value.spec.trim()) },
  { label: "数量", filled: Boolean(batchForm.value.quantity.trim()) },
  { label: "产品编码", filled: Boolean(batchExtraForm.value.productCode.trim()) },
  { label: "包装形式", filled: Boolean(batchExtraForm.value.packageType.trim()) },
]);

const batchProgress = computed(() => ({
  completed: batchChecklist.value.filter((item) => item.filled).length,
  total: batchChecklist.value.length,
  missingLabels: batchChecklist.value.filter((item) => !item.filled).map((item) => item.label),
}));

const eventChecklist = computed<ChecklistItem[]>(() => [
  { label: "事件标题", filled: Boolean(eventForm.value.title.trim()) },
  { label: "事件地点", filled: Boolean(eventForm.value.location.trim()) },
  { label: "操作人", filled: Boolean(eventForm.value.operatorName.trim()) },
  { label: "事件编码", filled: Boolean(eventMetaForm.value.eventCode.trim()) },
  { label: "SOP 编号", filled: Boolean(eventMetaForm.value.sopCode.trim()) },
  ...currentStageDefinition.value.fields
    .filter((field) => field.required)
    .map((field) => ({
      label: field.label,
      filled: Boolean((eventDetailForm.value[field.key] || "").trim()),
    })),
]);

const eventProgress = computed(() => ({
  completed: eventChecklist.value.filter((item) => item.filled).length,
  total: eventChecklist.value.length,
  missingLabels: eventChecklist.value.filter((item) => !item.filled).map((item) => item.label),
}));

const workflowHint = computed(() => {
  if (!currentBatch.value) return "请先选择批次，再开始事件登记。";
  if (recommendedStage.value) return `建议优先登记${recommendedStage.value.label}阶段，保持标准链路顺序。`;
  return "六个标准阶段均已有记录，可继续补录或修正细节。";
});

const currentBatchSummary = computed(() => {
  if (!currentBatch.value) return "未选择批次";
  return [currentBatch.value.productName, currentBatch.value.spec, currentBatch.value.quality]
    .filter(Boolean)
    .join(" / ") || "已建立批次上下文";
});

const batchExtraPreview = computed(() => buildBatchExtraJson(batchExtraForm.value));

function toDatetimeLocal(date: Date) {
  const offset = date.getTimezoneOffset();
  return new Date(date.getTime() - offset * 60 * 1000).toISOString().slice(0, 16);
}

function parseJsonObject(raw?: string): Record<string, unknown> | null {
  const content = (raw || "").trim();
  if (!content) return null;
  try {
    const parsed = JSON.parse(content);
    return parsed && typeof parsed === "object" && !Array.isArray(parsed)
      ? (parsed as Record<string, unknown>)
      : null;
  } catch {
    return null;
  }
}

function asObject(value: unknown): Record<string, unknown> | null {
  return value && typeof value === "object" && !Array.isArray(value)
    ? (value as Record<string, unknown>)
    : null;
}

function normalizeCode(raw: string) {
  return raw.trim().toUpperCase().replace(/[^A-Z0-9]+/g, "-").replace(/^-+|-+$/g, "");
}

function pad(value: number) {
  return String(value).padStart(2, "0");
}

function buildSuggestedBatchNo(date = new Date()) {
  return `HSC-${date.getFullYear()}${pad(date.getMonth() + 1)}${pad(date.getDate())}-${String(auth.user?.orgId ?? batchForm.value.orgId ?? 1).padStart(2, "0")}`;
}

function apiError(err: unknown, fallback: string) {
  if (axios.isAxiosError(err) && err.message) return err.message;
  return fallback;
}

function enrichTimelineItem(item: TimelineApiItem): TimelineItem {
  return {
    ...item,
    detailView: buildStructuredDetailView(item.event.detailJson, item.event.evidenceUrls),
  };
}

function latestTimelineItem(stage?: TraceStageValue): TimelineItem | null {
  const source = stage ? timeline.value.filter((item) => item.event.stage === stage) : timeline.value;
  return source.length ? source[source.length - 1] : null;
}

function inferStageLocation(stage: TraceStageValue): string {
  const latestSameStage = latestTimelineItem(stage);
  if (latestSameStage?.event.location) return latestSameStage.event.location;

  const latestAnyStage = latestTimelineItem();
  const batch = currentBatch.value;

  switch (stage) {
    case "breeding":
      return batch?.farmBase || batch?.breedArea || latestAnyStage?.event.location || "";
    case "harvest":
      return batch?.breedArea || batch?.farmBase || latestAnyStage?.event.location || "";
    case "processing":
    case "packaging":
      return latestAnyStage?.event.location || batch?.farmBase || "";
    default:
      return latestAnyStage?.event.location || "";
  }
}

function resetBatchForm() {
  batchForm.value = {
    batchNo: "",
    orgId: auth.user?.orgId ?? batchForm.value.orgId,
    productName: "",
    farmBase: "",
    quality: qualityLevelOptions[0].value,
    catchDate: "",
    breedArea: "",
    breedStartDate: "",
    spec: "",
    quantity: "",
  };
  batchExtraForm.value = createEmptyBatchExtraForm();
}

function resetEventForm(nextStage = eventForm.value.stage) {
  eventForm.value = {
    stage: nextStage,
    title: stageDefinitionMap[nextStage as TraceStageValue]?.defaultTitle || "",
    occurredAt: toDatetimeLocal(new Date()),
    location: "",
    operatorName: auth.user?.displayName || "",
  };
  eventMetaForm.value = createEventMetaForm(nextStage);
  eventDetailForm.value = createStageDetailForm(nextStage);
  evidenceInput.value = "";
}

function applyBatchTemplate() {
  const now = new Date();
  if (!batchForm.value.batchNo.trim()) {
    batchForm.value.batchNo = buildSuggestedBatchNo(now);
  }

  const normalizedBatchNo = normalizeCode(batchForm.value.batchNo);
  if (!batchExtraForm.value.productCode.trim() && normalizedBatchNo) {
    batchExtraForm.value.productCode = `SC-${normalizedBatchNo}`;
  }
  if (!batchExtraForm.value.packageType.trim() && batchForm.value.spec.trim()) {
    batchExtraForm.value.packageType = `标准包装 / ${batchForm.value.spec.trim()}`;
  }
  if (!batchExtraForm.value.standardCode.trim()) {
    batchExtraForm.value.standardCode = "SC/T 3208-2024";
  }
  if (!batchExtraForm.value.inspectorName.trim() && auth.user?.displayName) {
    batchExtraForm.value.inspectorName = auth.user.displayName;
  }
  if (!batchExtraForm.value.inspectionReportNo.trim()) {
    batchExtraForm.value.inspectionReportNo = `QC-${now.toISOString().slice(0, 10).replace(/-/g, "")}-${pad(now.getHours())}${pad(now.getMinutes())}`;
  }
  if (!batchExtraForm.value.certificateNo.trim() && normalizedBatchNo) {
    batchExtraForm.value.certificateNo = `CERT-${normalizedBatchNo}`;
  }
  if (!batchExtraForm.value.storageRequirement.trim()) {
    batchExtraForm.value.storageRequirement = "0-4°C 冷藏";
  }
  if (!batchExtraForm.value.originCity.trim()) {
    batchExtraForm.value.originCity = "大连";
  }

  error.value = "";
  notice.value = "已带入企业标准模板，可按实际情况微调。";
}

function applyEventTemplate(force = false) {
  const stage = eventForm.value.stage as TraceStageValue;
  const latestAnyStage = latestTimelineItem();
  const extraPayload = parseJsonObject(currentBatch.value?.extraJson);
  const productMeta = asObject(extraPayload?.product);

  if (force || !eventForm.value.title.trim() || eventForm.value.title === stageDefinitionMap[stage].defaultTitle) {
    eventForm.value.title = stageDefinitionMap[stage].defaultTitle;
  }
  if (force || !eventForm.value.location.trim()) {
    eventForm.value.location = inferStageLocation(stage);
  }
  if (force || !eventForm.value.operatorName.trim()) {
    eventForm.value.operatorName = auth.user?.displayName || latestAnyStage?.event.operatorName || "";
  }
  if (force || !eventMetaForm.value.eventCode.trim()) {
    eventMetaForm.value.eventCode = buildEventCode(stage);
  }
  if (force || !eventMetaForm.value.sopCode.trim()) {
    eventMetaForm.value.sopCode = defaultSopCode(stage);
  }

  const nextDetailForm = { ...eventDetailForm.value };
  switch (stage) {
    case "breeding":
      if (force || !nextDetailForm.pondId) {
        nextDetailForm.pondId = currentBatch.value?.farmBase || "";
      }
      if (force || !nextDetailForm.inspectionResult) {
        nextDetailForm.inspectionResult = "水质正常、活性良好";
      }
      break;
    case "harvest":
      if (force || !nextDetailForm.harvestMethod) {
        nextDetailForm.harvestMethod = "人工采捕";
      }
      if (force || !nextDetailForm.acceptanceResult) {
        nextDetailForm.acceptanceResult = "现场抽检合格，允许入库";
      }
      break;
    case "processing":
      if (force || !nextDetailForm.qaResult) {
        nextDetailForm.qaResult = "复检合格";
      }
      break;
    case "packaging":
      if (force || !nextDetailForm.packageType) {
        nextDetailForm.packageType = String(productMeta?.packageType || batchExtraForm.value.packageType || "");
      }
      if (force || !nextDetailForm.packageSpec) {
        nextDetailForm.packageSpec = currentBatch.value?.spec || "";
      }
      if (force || !nextDetailForm.labelInspector) {
        nextDetailForm.labelInspector = auth.user?.displayName || "";
      }
      if (force || !nextDetailForm.qrBindingStatus) {
        nextDetailForm.qrBindingStatus = "yes";
      }
      break;
    case "logistics":
      if (force || !nextDetailForm.handoverStatus) {
        nextDetailForm.handoverStatus = "待交接";
      }
      break;
    case "retail":
      if (force || !nextDetailForm.consumerStatus) {
        nextDetailForm.consumerStatus = "已上架，可扫码查询";
      }
      break;
  }

  eventDetailForm.value = nextDetailForm;

  if (force) {
    error.value = "";
    notice.value = `已带入${stageDefinitionMap[stage].label}阶段标准模板。`;
  }
}

function syncWorkbenchToBatch(pickRecommendedStage = false) {
  if (pickRecommendedStage && recommendedStage.value) {
    eventForm.value.stage = recommendedStage.value.value;
  }
  applyEventTemplate(false);
}

function selectStage(stage: TraceStageValue) {
  if (eventForm.value.stage === stage) return;
  eventForm.value.stage = stage;
}

function jumpToRecommendedStage() {
  if (!currentBatch.value) {
    error.value = "请先选择一个批次。";
    return;
  }
  if (!recommendedStage.value) {
    error.value = "";
    notice.value = "当前批次的标准阶段已覆盖完成，可按需补录。";
    return;
  }
  eventForm.value.stage = recommendedStage.value.value;
  applyEventTemplate(false);
}

function validateBatchForm(): string {
  if (!batchForm.value.batchNo.trim()) return "请填写批次号。";
  if (!batchForm.value.productName.trim()) return "请填写产品名称。";
  if (!batchForm.value.breedStartDate) return "请填写养殖起始日期。";
  if (!batchForm.value.catchDate) return "请填写采捕日期。";
  if (!batchForm.value.farmBase.trim()) return "请填写养殖基地。";
  if (!batchForm.value.breedArea.trim()) return "请填写养殖海域。";
  if (!batchForm.value.spec.trim()) return "请填写产品规格。";
  if (!batchForm.value.quantity.trim()) return "请填写数量。";
  if (!batchExtraForm.value.productCode.trim()) return "请填写产品编码。";
  if (!batchExtraForm.value.packageType.trim()) return "请填写包装形式。";
  return "";
}

function validateEventForm(): string {
  if (!currentBatch.value) return "请先选择一个批次。";
  if (!eventForm.value.title.trim()) return "请填写事件标题。";
  if (!eventForm.value.location.trim()) return "请填写事件地点。";
  if (!eventForm.value.operatorName.trim()) return "请填写操作人。";
  if (!eventMetaForm.value.eventCode.trim()) return "请填写事件编码。";
  if (!eventMetaForm.value.sopCode.trim()) return "请填写 SOP 编号。";

  for (const field of currentStageDefinition.value.fields) {
    if (field.required && !(eventDetailForm.value[field.key] || "").trim()) {
      return `请填写${field.label}。`;
    }
  }
  return "";
}

async function loadBatches(selectFirst = true) {
  loading.value = true;
  error.value = "";
  try {
    const { data } = await api.get("/batches");
    batches.value = data.items || [];
    if (selectFirst && !selectedBatchId.value && batches.value.length) {
      selectedBatchId.value = batches.value[0].id;
      await loadTimeline(batches.value[0].id);
      syncWorkbenchToBatch(true);
    } else if (selectedBatchId.value) {
      await loadTimeline(selectedBatchId.value);
      applyEventTemplate(false);
    }
  } catch (err) {
    error.value = apiError(err, "批次列表加载失败，请确认后端服务已启动并已登录。");
  } finally {
    loading.value = false;
  }
}

async function selectBatch(batch: Batch) {
  selectedBatchId.value = batch.id;
  await loadTimeline(batch.id);
  syncWorkbenchToBatch(true);
}

async function loadTimeline(batchId: number) {
  timelineLoading.value = true;
  error.value = "";
  try {
    const { data } = await api.get(`/batches/${batchId}/timeline`);
    timeline.value = (data.items || []).map(enrichTimelineItem);
  } catch (err) {
    timeline.value = [];
    error.value = apiError(err, "时间线加载失败。");
  } finally {
    timelineLoading.value = false;
  }
}

async function createBatch() {
  notice.value = "";
  error.value = "";

  const validationError = validateBatchForm();
  if (validationError) {
    error.value = validationError;
    return;
  }

  savingBatch.value = true;
  try {
    const body: Record<string, unknown> = {
      batchNo: batchForm.value.batchNo.trim(),
      orgId: Number(batchForm.value.orgId),
      productName: batchForm.value.productName.trim(),
      farmBase: batchForm.value.farmBase.trim(),
      quality: batchForm.value.quality.trim(),
      breedArea: batchForm.value.breedArea.trim(),
      spec: batchForm.value.spec.trim(),
      quantity: batchForm.value.quantity.trim(),
      extraJson: batchExtraPreview.value,
    };

    if (batchForm.value.catchDate) {
      body.catchDate = new Date(`${batchForm.value.catchDate}T12:00:00`).toISOString();
    }
    if (batchForm.value.breedStartDate) {
      body.breedStartDate = new Date(`${batchForm.value.breedStartDate}T12:00:00`).toISOString();
    }

    const { data } = await api.post("/batches", body);
    selectedBatchId.value = data.id;
    notice.value = "批次创建成功，可以继续登记流转记录。";
    resetBatchForm();
    await loadBatches(false);
    syncWorkbenchToBatch(true);
  } catch (err) {
    error.value = apiError(err, "批次创建失败，可能是批次号重复或权限不足。");
  } finally {
    savingBatch.value = false;
  }
}

async function addEvent() {
  notice.value = "";
  error.value = "";

  const validationError = validateEventForm();
  if (validationError) {
    error.value = validationError;
    return;
  }

  savingEvent.value = true;
  try {
    const evidenceList = parseEvidenceUrls(evidenceInput.value);
    const detailJson = buildEventDetailJson(
      eventForm.value.stage,
      eventMetaForm.value,
      eventDetailForm.value,
      evidenceList,
    );

    const body = {
      stage: eventForm.value.stage,
      title: eventForm.value.title.trim(),
      occurredAt: eventForm.value.occurredAt
        ? new Date(eventForm.value.occurredAt).toISOString()
        : new Date().toISOString(),
      location: eventForm.value.location.trim(),
      operatorName: eventForm.value.operatorName.trim(),
      detailJson,
      evidenceUrls: evidenceList.join(", "),
    };

    const { data } = await api.post(`/batches/${currentBatch.value!.id}/events`, body);
    const status = data.chain?.status === "success" ? "并已完成链上锚定" : "但链上锚定未完成";
    notice.value = `事件登记成功，${status}。`;
    resetEventForm(eventForm.value.stage);
    await loadTimeline(currentBatch.value!.id);
    applyEventTemplate(false);
  } catch (err) {
    error.value = apiError(err, "事件登记失败。");
  } finally {
    savingEvent.value = false;
  }
}

watch(
  () => eventForm.value.stage,
  (next, prev) => {
    const nextDefinition = stageDefinitionMap[next as TraceStageValue];
    const prevDefinition = prev ? stageDefinitionMap[prev as TraceStageValue] : null;

    if (!eventForm.value.title.trim() || eventForm.value.title === prevDefinition?.defaultTitle) {
      eventForm.value.title = nextDefinition.defaultTitle;
    }

    if (!eventMetaForm.value.eventCode.trim() || eventMetaForm.value.eventCode.startsWith(`${prevDefinition?.codePrefix || "EVT"}-`)) {
      eventMetaForm.value.eventCode = buildEventCode(next);
    }

    if (!eventMetaForm.value.sopCode.trim() || eventMetaForm.value.sopCode === defaultSopCode(prev || next)) {
      eventMetaForm.value.sopCode = defaultSopCode(next);
    }

    eventDetailForm.value = createStageDetailForm(next);
    applyEventTemplate(false);
  },
);

onMounted(() => {
  loadBatches();
});
</script>

<template>
  <section class="mx-auto max-w-6xl px-6 py-10 md:py-12">
    <div class="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
      <div>
        <p class="text-xs uppercase tracking-[0.24em] text-cyan-200/70">Trace Workspace</p>
        <h1 class="mt-2 font-display text-3xl font-semibold tracking-tight text-white">企业追溯工作台</h1>
        <p class="mt-2 max-w-3xl text-sm leading-relaxed text-slate-400">
          围绕批次建档、环节登记和进度查看，集中处理企业追溯业务。
        </p>
      </div>
      <button type="button" class="btn-ghost px-4 py-2" @click="loadBatches(false)">刷新数据</button>
    </div>

    <div v-if="notice" class="mt-6 rounded-2xl border border-emerald-400/25 bg-emerald-400/10 px-4 py-3 text-sm text-emerald-100">
      {{ notice }}
    </div>
    <div v-if="error" class="mt-6 rounded-2xl border border-rose-500/25 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
      {{ error }}
    </div>

    <div class="mt-8 grid gap-4 lg:grid-cols-3">
      <section class="glass rounded-2xl p-4">
        <p class="text-xs uppercase tracking-[0.18em] text-slate-500">当前批次</p>
        <p class="mt-2 truncate font-mono text-sm text-cyan-100">{{ currentBatch?.batchNo || "未选择" }}</p>
        <p class="mt-2 text-sm text-white">{{ currentBatchSummary }}</p>
      </section>
      <section class="glass rounded-2xl p-4">
        <p class="text-xs uppercase tracking-[0.18em] text-slate-500">环节进度</p>
        <p class="mt-2 text-2xl font-semibold text-white">{{ completedStageCount }}/{{ stageDefinitions.length }}</p>
        <p class="mt-2 text-sm text-slate-400">已完成环节越多，批次流转信息越完整。</p>
      </section>
      <section class="glass rounded-2xl p-4">
        <p class="text-xs uppercase tracking-[0.18em] text-slate-500">当前待办</p>
        <p class="mt-2 text-sm font-semibold text-white">
          {{ recommendedStage ? `${recommendedStage.label}环节` : "当前批次已登记完成" }}
        </p>
        <p class="mt-2 text-sm text-slate-400">{{ workflowHint }}</p>
      </section>
    </div>

    <div class="mt-6 grid gap-6 xl:grid-cols-[1.02fr_0.98fr]">
      <section class="glass noise rounded-[1.75rem] p-6 shadow-card ring-1 ring-white/[0.05]">
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div>
            <h2 class="text-base font-semibold text-white">1. 创建批次档案</h2>
            <p class="mt-1 text-sm text-slate-400">录入批次基础信息，用于后续生产、流转和追溯登记。</p>
          </div>
          <span class="rounded-full bg-cyan-400/10 px-3 py-1 text-[11px] font-medium text-cyan-100 ring-1 ring-cyan-400/20">
            批次档案
          </span>
        </div>

        <div class="mt-5 grid gap-3 md:grid-cols-[1fr_auto]">
          <div class="rounded-[1.25rem] border border-white/[0.08] bg-white/[0.03] px-4 py-3">
            <p class="text-xs uppercase tracking-[0.18em] text-slate-500">录入完成度</p>
            <p class="mt-2 text-2xl font-semibold text-white">{{ batchProgress.completed }}/{{ batchProgress.total }}</p>
            <p class="mt-1 text-xs text-slate-400">
              {{ batchProgress.missingLabels.length ? `待补：${batchProgress.missingLabels.slice(0, 4).join("、")}` : "核心字段已齐备，可直接建档。" }}
            </p>
          </div>
          <button type="button" class="btn-ghost px-4 py-3 text-sm" @click="applyBatchTemplate">快捷带入</button>
        </div>

        <div class="mt-5 space-y-6">
          <div>
            <p class="text-xs font-medium uppercase tracking-[0.2em] text-slate-500">基础档案</p>
            <div class="mt-4 grid gap-4 md:grid-cols-2">
              <label class="md:col-span-2">
                <span class="mb-1.5 block text-xs font-medium text-slate-400">批次号 *</span>
                <input v-model="batchForm.batchNo" class="input-field font-mono" placeholder="HSC-2026-0001" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">产品名称 *</span>
                <input v-model="batchForm.productName" class="input-field" placeholder="大连刺参（精品级）" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">养殖起始日期 *</span>
                <input v-model="batchForm.breedStartDate" type="date" class="input-field" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">采捕日期 *</span>
                <input v-model="batchForm.catchDate" type="date" class="input-field" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">养殖基地 *</span>
                <input v-model="batchForm.farmBase" class="input-field" placeholder="基地名称" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">养殖海域 *</span>
                <input v-model="batchForm.breedArea" class="input-field" placeholder="大连长海海域 A 区" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">质量等级</span>
                <select v-model="batchForm.quality" class="input-field">
                  <option v-for="option in qualityLevelOptions" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">产品规格 *</span>
                <input v-model="batchForm.spec" class="input-field" placeholder="500g/袋" />
              </label>
              <label class="md:col-span-2">
                <span class="mb-1.5 block text-xs font-medium text-slate-400">数量 *</span>
                <input v-model="batchForm.quantity" class="input-field" placeholder="500" />
              </label>
            </div>
          </div>

          <div class="rounded-[1.35rem] border border-white/[0.08] bg-white/[0.03] p-5">
            <div class="flex flex-wrap items-start justify-between gap-4">
              <div>
                <p class="text-sm font-semibold text-white">补充信息</p>
                <p class="mt-1 text-xs leading-relaxed text-slate-400">完善产品识别、包装和展示信息，方便内部管理与对外展示。</p>
              </div>
              <span class="rounded-full bg-white/5 px-3 py-1 text-[11px] text-slate-300 ring-1 ring-white/10">产品信息</span>
            </div>
            <div class="mt-4 grid gap-4 md:grid-cols-2">
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">产品编码 *</span>
                <input v-model="batchExtraForm.productCode" class="input-field font-mono" placeholder="SC-SKU-2026-001" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">计量单位</span>
                <select v-model="batchExtraForm.unit" class="input-field">
                  <option v-for="option in unitOptions" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">包装形式 *</span>
                <input v-model="batchExtraForm.packageType" class="input-field" placeholder="真空袋 + 外箱" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">保质期（天）</span>
                <input v-model="batchExtraForm.shelfLifeDays" type="number" min="0" class="input-field" placeholder="365" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">储运要求</span>
                <input v-model="batchExtraForm.storageRequirement" class="input-field" placeholder="0-4°C 冷藏" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">产地城市</span>
                <input v-model="batchExtraForm.originCity" class="input-field" placeholder="大连" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">检测报告编号</span>
                <input v-model="batchExtraForm.inspectionReportNo" class="input-field" placeholder="QC-20260509-008" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">质检负责人</span>
                <input v-model="batchExtraForm.inspectorName" class="input-field" placeholder="王敏" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">合格证编号</span>
                <input v-model="batchExtraForm.certificateNo" class="input-field" placeholder="CERT-20260509-02" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">执行标准</span>
                <input v-model="batchExtraForm.standardCode" class="input-field" placeholder="SC/T 3208-2024" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">追溯标签类型</span>
                <select v-model="batchExtraForm.traceTagType" class="input-field">
                  <option v-for="option in traceTagOptions" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </label>
              <label class="md:col-span-2">
                <span class="mb-1.5 block text-xs font-medium text-slate-400">补充说明</span>
                <textarea v-model="batchExtraForm.remark" rows="3" class="input-field resize-y" placeholder="例如：该批次用于精品渠道首发，已完成出厂抽检。" />
              </label>
            </div>
          </div>

        </div>

        <button type="button" class="btn-primary mt-6 w-full" :disabled="savingBatch" @click="createBatch">
          {{ savingBatch ? "创建中..." : "创建批次档案" }}
        </button>
      </section>

      <section class="glass noise rounded-[1.75rem] p-6 shadow-card ring-1 ring-white/[0.05]">
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div>
            <h2 class="text-base font-semibold text-white">2. 登记追溯事件</h2>
            <p class="mt-1 text-sm text-slate-400">按业务环节登记生产与流转信息，形成完整批次记录。</p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <span class="rounded-full bg-white/5 px-3 py-1 text-[11px] text-slate-300 ring-1 ring-white/10">
              当前阶段：{{ currentStageDefinition.label }}
            </span>
            <RouterLink
              v-if="currentBatch"
              class="rounded-full bg-cyan-400/10 px-3 py-1 text-[11px] text-cyan-100 ring-1 ring-cyan-400/20 hover:bg-cyan-400/15"
              :to="{ name: 'trace', params: { batchNo: currentBatch.batchNo } }"
            >
              查看消费者页面
            </RouterLink>
          </div>
        </div>

        <div class="mt-5 space-y-6">
          <div class="rounded-[1.35rem] border border-white/[0.08] bg-white/[0.03] p-5">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div>
                <p class="text-sm font-semibold text-white">业务环节</p>
                <p class="mt-1 text-xs leading-relaxed text-slate-400">{{ workflowHint }}</p>
              </div>
              <button type="button" class="btn-ghost px-3 py-2 text-xs" :disabled="!currentBatch" @click="jumpToRecommendedStage">
                定位待办环节
              </button>
            </div>
            <div class="mt-4 grid gap-2 md:grid-cols-3 xl:grid-cols-6">
              <button
                v-for="stage in stageProgress"
                :key="stage.value"
                type="button"
                class="rounded-2xl border px-3 py-3 text-left transition"
                :class="stage.isCurrent
                  ? 'border-cyan-400/40 bg-cyan-400/12'
                  : stage.completed
                    ? 'border-emerald-400/25 bg-emerald-400/8 hover:border-emerald-300/30'
                    : stage.isRecommended
                      ? 'border-amber-400/35 bg-amber-400/10 hover:border-amber-300/45'
                      : 'border-white/10 bg-white/[0.02] hover:border-white/20'"
                @click="selectStage(stage.value)"
              >
                <span class="text-[11px] text-slate-500">步骤 {{ stage.step }}</span>
                <span class="mt-2 block text-sm font-semibold" :class="stage.isCurrent ? 'text-cyan-100' : 'text-white'">
                  {{ stage.label }}
                </span>
                <span class="mt-1 block text-[11px]" :class="stage.completed ? 'text-emerald-200' : stage.isRecommended ? 'text-amber-200' : 'text-slate-500'">
                  {{ stage.count ? `${stage.count} 条记录` : stage.isRecommended ? "待处理" : "待登记" }}
                </span>
              </button>
            </div>
          </div>

          <div class="grid gap-3 md:grid-cols-[1fr_auto_auto]">
            <div class="rounded-[1.25rem] border border-white/[0.08] bg-white/[0.03] px-4 py-3">
              <div class="flex flex-wrap gap-2">
                <span class="rounded-full bg-cyan-400/10 px-3 py-1 text-[11px] text-cyan-100 ring-1 ring-cyan-400/20">
                  当前阶段：{{ currentStageDefinition.label }}
                </span>
                <span class="rounded-full bg-white/5 px-3 py-1 text-[11px] text-slate-300 ring-1 ring-white/10">
                  已完成 {{ eventProgress.completed }}/{{ eventProgress.total }}
                </span>
                <span class="rounded-full bg-white/5 px-3 py-1 text-[11px] text-slate-300 ring-1 ring-white/10">
                  标准字段自动生成
                </span>
              </div>
              <p class="mt-3 text-xs text-slate-400">
                {{ eventProgress.missingLabels.length ? `待补：${eventProgress.missingLabels.slice(0, 5).join("、")}` : "当前表单必填项已齐备，可直接登记。" }}
              </p>
            </div>
            <button type="button" class="btn-ghost px-4 py-3 text-sm" :disabled="!currentBatch" @click="applyEventTemplate(true)">
              带入当前阶段模板
            </button>
            <button type="button" class="btn-ghost px-4 py-3 text-sm" @click="showAdvancedMeta = !showAdvancedMeta">
              {{ showAdvancedMeta ? "收起标准字段" : "展开标准字段" }}
            </button>
          </div>

          <div class="grid gap-4 md:grid-cols-2">
            <label class="md:col-span-2">
              <span class="mb-1.5 block text-xs font-medium text-slate-400">批次上下文</span>
              <div class="input-field flex min-h-[78px] items-center bg-white/[0.02] text-sm text-slate-200">
                <div v-if="currentBatch" class="min-w-0">
                  <p class="truncate font-mono text-cyan-100">{{ currentBatch.batchNo }}</p>
                  <p class="mt-1 truncate text-xs text-slate-400">{{ currentBatchSummary }}</p>
                </div>
                <span v-else class="text-slate-500">请先从左侧选择批次</span>
              </div>
            </label>
            <label class="md:col-span-2">
              <span class="mb-1.5 block text-xs font-medium text-slate-400">事件标题 *</span>
              <input v-model="eventForm.title" class="input-field" placeholder="请输入事件标题" />
            </label>
            <label>
              <span class="mb-1.5 block text-xs font-medium text-slate-400">发生时间</span>
              <input v-model="eventForm.occurredAt" type="datetime-local" class="input-field" />
            </label>
            <label>
              <span class="mb-1.5 block text-xs font-medium text-slate-400">事件地点 *</span>
              <input v-model="eventForm.location" class="input-field" placeholder="例如 大连加工厂一号车间" />
            </label>
            <label>
              <span class="mb-1.5 block text-xs font-medium text-slate-400">操作人 *</span>
              <input v-model="eventForm.operatorName" class="input-field" placeholder="请输入责任人姓名" />
            </label>
          </div>

          <div v-if="showAdvancedMeta" class="rounded-[1.35rem] border border-white/[0.08] bg-white/[0.03] p-5">
            <div class="flex flex-wrap items-start justify-between gap-4">
              <div>
                <p class="text-sm font-semibold text-white">系统信息</p>
                <p class="mt-1 text-xs leading-relaxed text-slate-400">通常无需调整，特殊情况下可手动修正。</p>
              </div>
              <span class="rounded-full bg-white/5 px-3 py-1 text-[11px] text-slate-300 ring-1 ring-white/10">
                系统生成
              </span>
            </div>
            <div class="mt-4 grid gap-4 md:grid-cols-3">
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">事件编码 *</span>
                <input v-model="eventMetaForm.eventCode" class="input-field font-mono" placeholder="系统自动生成，可手动调整" />
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">执行状态</span>
                <select v-model="eventMetaForm.status" class="input-field">
                  <option v-for="option in eventStatusOptions" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
              </label>
              <label>
                <span class="mb-1.5 block text-xs font-medium text-slate-400">SOP 编号 *</span>
                <input v-model="eventMetaForm.sopCode" class="input-field font-mono" placeholder="SOP-BRD-001" />
              </label>
            </div>
          </div>

          <div class="rounded-[1.35rem] border border-white/[0.08] bg-white/[0.03] p-5">
            <div class="flex flex-wrap items-start justify-between gap-4">
              <div>
                <p class="text-sm font-semibold text-white">{{ currentStageDefinition.label }}阶段明细</p>
                <p class="mt-1 text-xs leading-relaxed text-slate-400">{{ currentStageDefinition.description }}</p>
              </div>
              <span class="rounded-full bg-white/5 px-3 py-1 text-[11px] text-slate-300 ring-1 ring-white/10">
                环节信息
              </span>
            </div>

            <div class="mt-4 grid gap-4 md:grid-cols-2">
              <label
                v-for="field in currentStageDefinition.fields"
                :key="field.key"
                :class="field.type === 'textarea' ? 'md:col-span-2' : ''"
              >
                <span class="mb-1.5 block text-xs font-medium text-slate-400">
                  {{ field.label }}<span v-if="field.required"> *</span>
                </span>
                <select
                  v-if="field.type === 'select'"
                  v-model="eventDetailForm[field.key]"
                  class="input-field"
                >
                  <option value="">请选择</option>
                  <option v-for="option in field.options || []" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </option>
                </select>
                <textarea
                  v-else-if="field.type === 'textarea'"
                  v-model="eventDetailForm[field.key]"
                  rows="3"
                  class="input-field resize-y"
                  :placeholder="field.placeholder"
                />
                <input
                  v-else
                  v-model="eventDetailForm[field.key]"
                  class="input-field"
                  :type="field.type === 'number' ? 'number' : 'text'"
                  :step="field.type === 'number' ? '0.1' : undefined"
                  :placeholder="field.placeholder"
                />
              </label>

              <label class="md:col-span-2">
                <span class="mb-1.5 block text-xs font-medium text-slate-400">佐证材料链接</span>
                <textarea
                  v-model="evidenceInput"
                  rows="3"
                  class="input-field resize-y"
                  placeholder="每行一个链接，或使用逗号分隔多个链接"
                />
              </label>

              <label class="md:col-span-2">
                <span class="mb-1.5 block text-xs font-medium text-slate-400">补充说明</span>
                <textarea
                  v-model="eventMetaForm.note"
                  rows="3"
                  class="input-field resize-y"
                  placeholder="记录异常说明、补录原因或备注信息"
                />
              </label>
            </div>
          </div>

        </div>

        <button type="button" class="btn-primary mt-6 w-full" :disabled="savingEvent || !currentBatch" @click="addEvent">
          {{ savingEvent ? "登记中..." : "提交事件记录" }}
        </button>
      </section>
    </div>

    <div class="mt-10 grid gap-6 lg:grid-cols-[0.78fr_1.22fr]">
      <section>
        <div class="flex items-baseline justify-between">
          <h2 class="text-lg font-semibold text-white">批次列表</h2>
          <span class="text-xs text-slate-500">{{ batches.length }} 条</span>
        </div>

        <div v-if="loading" class="mt-4 space-y-3">
          <div v-for="i in 3" :key="i" class="skeleton h-24 rounded-2xl" />
        </div>
        <div v-else-if="!batches.length" class="mt-4 rounded-2xl border border-dashed border-white/10 px-6 py-10 text-center text-sm text-slate-500">
          暂无批次，请先新建批次。
        </div>
        <div v-else class="mt-4 space-y-3">
          <button
            v-for="batch in batches"
            :key="batch.id"
            type="button"
            class="w-full rounded-2xl border p-4 text-left transition"
            :class="batch.id === selectedBatchId ? 'border-cyan-400/35 bg-cyan-400/10' : 'border-white/10 bg-white/5 hover:border-white/20'"
            @click="selectBatch(batch)"
          >
            <span class="block truncate font-mono text-sm text-cyan-100">{{ batch.batchNo }}</span>
            <span class="mt-1 block truncate text-sm text-white/90">{{ batch.productName || "未填写产品名称" }}</span>
            <div class="mt-2 flex flex-wrap gap-2 text-[11px] text-slate-500">
              <span>{{ batch.org?.name || `orgId ${batch.orgId}` }}</span>
              <span v-if="batch.quality">| {{ batch.quality }}</span>
              <span v-if="batch.spec">| {{ batch.spec }}</span>
            </div>
          </button>
        </div>
      </section>

      <section>
        <div class="flex items-baseline justify-between gap-4">
          <h2 class="text-lg font-semibold text-white">当前批次时间线</h2>
          <span class="text-xs text-slate-500">{{ timeline.length }} 个事件</span>
        </div>

        <div v-if="timelineLoading" class="mt-4 space-y-3">
          <div v-for="i in 3" :key="i" class="skeleton h-28 rounded-2xl" />
        </div>
        <div v-else-if="!currentBatch" class="mt-4 rounded-2xl border border-dashed border-white/10 px-6 py-10 text-center text-sm text-slate-500">
          请选择一个批次。
        </div>
        <div v-else-if="!timeline.length" class="mt-4 rounded-2xl border border-dashed border-white/10 px-6 py-10 text-center text-sm text-slate-500">
          当前批次暂无事件记录。
        </div>
        <ol v-else class="mt-4 border-l border-white/10 pl-5">
          <li v-for="item in timeline" :key="item.event.id" class="relative pb-6 last:pb-0">
            <span class="absolute -left-[25px] top-2 h-2.5 w-2.5 rounded-full bg-cyan-300 shadow-[0_0_14px_rgba(103,232,249,0.9)]" />
            <article class="glass rounded-2xl p-4">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <p class="text-sm font-semibold text-white">{{ item.event.title }}</p>
                  <p class="mt-1 text-xs text-slate-500">
                    {{ new Date(item.event.occurredAt).toLocaleString() }}
                    <span v-if="item.event.location"> | {{ item.event.location }}</span>
                    <span v-if="item.event.operatorName"> | {{ item.event.operatorName }}</span>
                  </p>
                </div>
                <span class="rounded-full bg-white/5 px-2.5 py-1 text-xs text-cyan-100 ring-1 ring-white/10">
                  {{ stageLabel(item.event.stage) }}
                </span>
              </div>

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

              <div class="mt-4 space-y-1 text-xs">
                <p class="break-all text-slate-400">
                  数据哈希：<span class="font-mono text-cyan-100">{{ item.event.dataHash }}</span>
                </p>
                <p v-if="item.chain?.txId" class="break-all text-emerald-200">
                  链上交易：<span class="font-mono">{{ item.chain.txId }}</span>
                  <span class="text-slate-500">（{{ item.chain.status }}）</span>
                </p>
                <p v-else class="text-amber-200">链上状态：{{ item.chain?.status || "pending" }}</p>
              </div>
            </article>
          </li>
        </ol>
      </section>
    </div>
  </section>
</template>
