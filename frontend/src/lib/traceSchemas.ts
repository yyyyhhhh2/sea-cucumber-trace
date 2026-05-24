export type TraceStageValue =
  | "breeding"
  | "harvest"
  | "processing"
  | "packaging"
  | "logistics"
  | "retail";

export type SelectOption = {
  value: string;
  label: string;
};

export type SchemaField = {
  key: string;
  label: string;
  placeholder: string;
  type?: "text" | "number" | "select" | "textarea";
  options?: SelectOption[];
  required?: boolean;
};

export type StageDefinition = {
  value: TraceStageValue;
  label: string;
  codePrefix: string;
  defaultTitle: string;
  description: string;
  fields: SchemaField[];
};

export type BatchExtraForm = {
  productCode: string;
  unit: string;
  packageType: string;
  shelfLifeDays: string;
  storageRequirement: string;
  inspectionReportNo: string;
  inspectorName: string;
  certificateNo: string;
  traceTagType: string;
  standardCode: string;
  originCity: string;
  remark: string;
};

export type EventMetaForm = {
  eventCode: string;
  status: string;
  sopCode: string;
  note: string;
};

export type DetailItem = {
  label: string;
  value: string;
};

export type DetailSection = {
  title: string;
  items: DetailItem[];
};

export type DetailView = {
  summary: DetailItem[];
  sections: DetailSection[];
};

export const qualityLevelOptions: SelectOption[] = [
  { value: "A", label: "A 级" },
  { value: "B", label: "B 级" },
  { value: "Premium", label: "Premium" },
  { value: "Qualified", label: "合格" },
];

export const unitOptions: SelectOption[] = [
  { value: "kg", label: "kg" },
  { value: "box", label: "箱" },
  { value: "bag", label: "袋" },
];

export const traceTagOptions: SelectOption[] = [
  { value: "qr-code", label: "二维码" },
  { value: "nfc", label: "NFC 标签" },
  { value: "barcode", label: "条形码" },
];

export const booleanStatusOptions: SelectOption[] = [
  { value: "yes", label: "是" },
  { value: "no", label: "否" },
];

export const eventStatusOptions: SelectOption[] = [
  { value: "completed", label: "已完成" },
  { value: "pending", label: "待补录" },
  { value: "exception", label: "异常待处理" },
];

export const stageDefinitions: StageDefinition[] = [
  {
    value: "breeding",
    label: "养殖",
    codePrefix: "BRD",
    defaultTitle: "养殖巡检记录",
    description: "记录养殖环境、投喂批次和日常巡检结果。",
    fields: [
      { key: "pondId", label: "养殖池/网箱编号", placeholder: "例如 C-12", required: true },
      { key: "waterTempC", label: "水温（°C）", placeholder: "例如 12.5", type: "number", required: true },
      { key: "salinityPpt", label: "盐度（‰）", placeholder: "例如 30", type: "number" },
      { key: "feedBatchNo", label: "投喂批次号", placeholder: "例如 FEED-20260509-01" },
      { key: "inspectionResult", label: "巡检结论", placeholder: "例如 水质正常、活性良好", required: true },
    ],
  },
  {
    value: "harvest",
    label: "采捕",
    codePrefix: "HVT",
    defaultTitle: "采捕交接记录",
    description: "记录采捕方式、重量和现场验收情况。",
    fields: [
      { key: "harvestMethod", label: "采捕方式", placeholder: "例如 人工采捕", required: true },
      { key: "vesselName", label: "作业船只/班组", placeholder: "例如 辽渔 08 号" },
      { key: "grossWeightKg", label: "毛重（kg）", placeholder: "例如 520", type: "number", required: true },
      { key: "onsiteTempC", label: "现场温度（°C）", placeholder: "例如 4", type: "number" },
      { key: "acceptanceResult", label: "验收结论", placeholder: "例如 抽检合格，允许入库", required: true },
    ],
  },
  {
    value: "processing",
    label: "加工",
    codePrefix: "PRC",
    defaultTitle: "加工质检记录",
    description: "记录车间、工艺和质量复核结果。",
    fields: [
      { key: "workshopNo", label: "车间/产线编号", placeholder: "例如 Workshop-03", required: true },
      { key: "processType", label: "加工工艺", placeholder: "例如 清洗、分级、预煮", required: true },
      { key: "centerTempC", label: "中心温度（°C）", placeholder: "例如 3", type: "number" },
      { key: "qaResult", label: "质检结果", placeholder: "例如 微生物抽检合格", required: true },
      { key: "packageLotNo", label: "关联包装批号", placeholder: "例如 PKG-20260509-02" },
    ],
  },
  {
    value: "packaging",
    label: "包装",
    codePrefix: "PKG",
    defaultTitle: "包装赋码记录",
    description: "记录包装规格、数量和追溯码绑定状态。",
    fields: [
      { key: "packageType", label: "包装形式", placeholder: "例如 真空袋 + 外箱", required: true },
      { key: "packageSpec", label: "包装规格", placeholder: "例如 500g/袋", required: true },
      { key: "packageCount", label: "包装数量", placeholder: "例如 100", type: "number", required: true },
      { key: "labelInspector", label: "标签复核人", placeholder: "例如 李四" },
      { key: "qrBindingStatus", label: "追溯码绑定", placeholder: "请选择", type: "select", options: booleanStatusOptions, required: true },
    ],
  },
  {
    value: "logistics",
    label: "物流",
    codePrefix: "LOG",
    defaultTitle: "冷链发运记录",
    description: "记录承运信息、运输温度和交接状态。",
    fields: [
      { key: "logisticsProvider", label: "承运商", placeholder: "例如 大连冷链物流", required: true },
      { key: "vehicleNo", label: "车辆/运单号", placeholder: "例如 辽B-12345", required: true },
      { key: "departureTempC", label: "发运温度（°C）", placeholder: "例如 2", type: "number", required: true },
      { key: "destination", label: "目的地", placeholder: "例如 上海嘉定分拨中心", required: true },
      { key: "handoverStatus", label: "交接状态", placeholder: "例如 已签收，无破损", required: true },
    ],
  },
  {
    value: "retail",
    label: "零售",
    codePrefix: "RTL",
    defaultTitle: "门店上架记录",
    description: "记录上架门店、陈列环境和销售端可见状态。",
    fields: [
      { key: "channelName", label: "销售渠道", placeholder: "例如 精品商超", required: true },
      { key: "storeName", label: "门店名称", placeholder: "例如 SeaTrace 大连体验店", required: true },
      { key: "shelfTempC", label: "陈列温度（°C）", placeholder: "例如 4", type: "number" },
      { key: "displayArea", label: "陈列区域", placeholder: "例如 冷柜 A 区" },
      { key: "consumerStatus", label: "终端状态", placeholder: "例如 已上架，可扫码查询", required: true },
    ],
  },
];

export const stageDefinitionMap = Object.fromEntries(
  stageDefinitions.map((definition) => [definition.value, definition]),
) as Record<TraceStageValue, StageDefinition>;

export function stageLabel(stage: string): string {
  return stageDefinitionMap[stage as TraceStageValue]?.label || stage;
}

export function eventStatusLabel(status: string): string {
  return eventStatusOptions.find((option) => option.value === status)?.label || status;
}

export function createEmptyBatchExtraForm(): BatchExtraForm {
  return {
    productCode: "",
    unit: "kg",
    packageType: "",
    shelfLifeDays: "",
    storageRequirement: "0-4°C 冷藏",
    inspectionReportNo: "",
    inspectorName: "",
    certificateNo: "",
    traceTagType: "qr-code",
    standardCode: "",
    originCity: "大连",
    remark: "",
  };
}

export function buildBatchExtraPayload(form: BatchExtraForm) {
  return compactValue({
    schemaVersion: "batch-profile.v1",
    product: {
      productCode: form.productCode.trim(),
      unit: form.unit.trim(),
      packageType: form.packageType.trim(),
      shelfLifeDays: form.shelfLifeDays.trim(),
      storageRequirement: form.storageRequirement.trim(),
    },
    quality: {
      inspectionReportNo: form.inspectionReportNo.trim(),
      inspectorName: form.inspectorName.trim(),
      certificateNo: form.certificateNo.trim(),
      standardCode: form.standardCode.trim(),
    },
    traceability: {
      traceTagType: form.traceTagType.trim(),
      originCity: form.originCity.trim(),
    },
    note: form.remark.trim(),
  });
}

export function buildBatchExtraJson(form: BatchExtraForm): string {
  return JSON.stringify(buildBatchExtraPayload(form), null, 2);
}

export function buildEventCode(stage: string, at = new Date()): string {
  const definition = stageDefinitionMap[stage as TraceStageValue];
  const prefix = definition?.codePrefix || "EVT";
  const yyyy = at.getFullYear();
  const mm = pad(at.getMonth() + 1);
  const dd = pad(at.getDate());
  const hh = pad(at.getHours());
  const min = pad(at.getMinutes());
  const sec = pad(at.getSeconds());
  return `${prefix}-${yyyy}${mm}${dd}-${hh}${min}${sec}`;
}

export function defaultSopCode(stage: string): string {
  const definition = stageDefinitionMap[stage as TraceStageValue];
  const prefix = definition?.codePrefix || "EVT";
  return `SOP-${prefix}-001`;
}

export function createEventMetaForm(stage: string, at = new Date()): EventMetaForm {
  return {
    eventCode: buildEventCode(stage, at),
    status: "completed",
    sopCode: defaultSopCode(stage),
    note: "",
  };
}

export function createStageDetailForm(stage: string): Record<string, string> {
  const definition = stageDefinitionMap[stage as TraceStageValue];
  const initial: Record<string, string> = {};
  for (const field of definition?.fields || []) {
    initial[field.key] = "";
  }
  return initial;
}

export function parseEvidenceUrls(raw: string): string[] {
  return raw
    .split(/[\n,]+/)
    .map((value) => value.trim())
    .filter(Boolean);
}

export function buildEventDetailPayload(
  stage: string,
  meta: EventMetaForm,
  details: Record<string, string>,
  evidenceList: string[],
) {
  const definition = stageDefinitionMap[stage as TraceStageValue];
  const record: Record<string, string> = {};
  for (const field of definition?.fields || []) {
    const value = (details[field.key] || "").trim();
    if (value) record[field.key] = value;
  }

  return compactValue({
    schemaVersion: "event-detail.v1",
    stage,
    stageLabel: definition?.label || stage,
    eventCode: meta.eventCode.trim(),
    status: meta.status.trim() || "completed",
    sopCode: meta.sopCode.trim(),
    record,
    attachments: evidenceList,
    note: meta.note.trim(),
  });
}

export function buildEventDetailJson(
  stage: string,
  meta: EventMetaForm,
  details: Record<string, string>,
  evidenceList: string[],
): string {
  return JSON.stringify(buildEventDetailPayload(stage, meta, details, evidenceList), null, 2);
}

export function buildStructuredDetailView(detailJson?: string, evidenceUrls?: string): DetailView | null {
  const raw = (detailJson || "").trim();
  const fallbackAttachments = parseEvidenceUrls(evidenceUrls || "");

  if (!raw) {
    if (!fallbackAttachments.length) return null;
    return {
      summary: [],
      sections: [
        {
          title: "佐证材料",
          items: fallbackAttachments.map((value, index) => ({ label: `附件 ${index + 1}`, value })),
        },
      ],
    };
  }

  try {
    const parsed = JSON.parse(raw) as Record<string, unknown>;
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
      return fallbackTextView(raw, fallbackAttachments);
    }

    if (parsed.schemaVersion === "event-detail.v1") {
      return structuredEventView(parsed, fallbackAttachments);
    }

    const genericSections = genericObjectSections(parsed);
    if (fallbackAttachments.length) {
      genericSections.push({
        title: "佐证材料",
        items: fallbackAttachments.map((value, index) => ({ label: `附件 ${index + 1}`, value })),
      });
    }
    return { summary: [], sections: genericSections };
  } catch {
    return fallbackTextView(raw, fallbackAttachments);
  }
}

function structuredEventView(parsed: Record<string, unknown>, fallbackAttachments: string[]): DetailView {
  const stage = typeof parsed.stage === "string" ? parsed.stage : "";
  const definition = stageDefinitionMap[stage as TraceStageValue];
  const summary: DetailItem[] = [];
  const sections: DetailSection[] = [];

  if (typeof parsed.eventCode === "string" && parsed.eventCode.trim()) {
    summary.push({ label: "事件编码", value: parsed.eventCode.trim() });
  }
  if (typeof parsed.status === "string" && parsed.status.trim()) {
    summary.push({ label: "执行状态", value: eventStatusLabel(parsed.status.trim()) });
  }
  if (typeof parsed.sopCode === "string" && parsed.sopCode.trim()) {
    summary.push({ label: "SOP 编号", value: parsed.sopCode.trim() });
  }

  const record = isPlainObject(parsed.record) ? parsed.record : null;
  if (definition && record) {
    const items = definition.fields
      .map((field) => ({
        label: field.label,
        value: normalizeValue(record[field.key]),
      }))
      .filter((item) => item.value);
    if (items.length) {
      sections.push({ title: `${definition.label}详情`, items });
    }
  } else if (record) {
    const items = objectEntries(record);
    if (items.length) {
      sections.push({ title: "事件详情", items });
    }
  }

  if (typeof parsed.note === "string" && parsed.note.trim()) {
    sections.push({
      title: "补充说明",
      items: [{ label: "备注", value: parsed.note.trim() }],
    });
  }

  const attachments = Array.isArray(parsed.attachments)
    ? parsed.attachments.map((value) => String(value).trim()).filter(Boolean)
    : fallbackAttachments;
  if (attachments.length) {
    sections.push({
      title: "佐证材料",
      items: attachments.map((value, index) => ({ label: `附件 ${index + 1}`, value })),
    });
  }

  return { summary, sections };
}

function genericObjectSections(parsed: Record<string, unknown>): DetailSection[] {
  const baseItems: DetailItem[] = [];
  const sections: DetailSection[] = [];

  for (const [key, value] of Object.entries(parsed)) {
    if (isPlainObject(value)) {
      const items = objectEntries(value);
      if (items.length) sections.push({ title: humanizeKey(key), items });
      continue;
    }

    if (Array.isArray(value)) {
      const items = value
        .map((item) => normalizeValue(item))
        .filter(Boolean)
        .map((item, index) => ({ label: `${humanizeKey(key)} ${index + 1}`, value: item }));
      if (items.length) sections.push({ title: humanizeKey(key), items });
      continue;
    }

    const normalized = normalizeValue(value);
    if (normalized) {
      baseItems.push({ label: humanizeKey(key), value: normalized });
    }
  }

  if (baseItems.length) {
    sections.unshift({ title: "事件详情", items: baseItems });
  }
  return sections;
}

function objectEntries(value: Record<string, unknown>): DetailItem[] {
  return Object.entries(value)
    .map(([key, itemValue]) => ({
      label: humanizeKey(key),
      value: normalizeValue(itemValue),
    }))
    .filter((item) => item.value);
}

function fallbackTextView(raw: string, attachments: string[]): DetailView {
  const sections: DetailSection[] = [
    {
      title: "事件详情",
      items: [{ label: "原始内容", value: raw }],
    },
  ];

  if (attachments.length) {
    sections.push({
      title: "佐证材料",
      items: attachments.map((value, index) => ({ label: `附件 ${index + 1}`, value })),
    });
  }

  return { summary: [], sections };
}

function normalizeValue(value: unknown): string {
  if (value == null) return "";
  if (typeof value === "string") return value.trim();
  if (typeof value === "number" || typeof value === "bigint") return String(value);
  if (typeof value === "boolean") return value ? "是" : "否";
  if (Array.isArray(value)) {
    return value.map((item) => normalizeValue(item)).filter(Boolean).join("、");
  }
  if (isPlainObject(value)) {
    return JSON.stringify(value);
  }
  return String(value).trim();
}

function humanizeKey(key: string): string {
  const dictionary: Record<string, string> = {
    schemaVersion: "Schema 版本",
    stage: "阶段编码",
    stageLabel: "阶段名称",
    eventCode: "事件编码",
    status: "执行状态",
    sopCode: "SOP 编号",
    note: "备注",
    productCode: "产品编码",
    unit: "计量单位",
    packageType: "包装形式",
    shelfLifeDays: "保质期（天）",
    storageRequirement: "储运要求",
    inspectionReportNo: "检测报告编号",
    inspectorName: "质检负责人",
    certificateNo: "合格证编号",
    standardCode: "执行标准",
    traceTagType: "追溯标签类型",
    originCity: "产地城市",
  };
  if (dictionary[key]) return dictionary[key];
  return key.replace(/([a-z0-9])([A-Z])/g, "$1 $2");
}

function compactValue<T>(value: T): T {
  if (Array.isArray(value)) {
    return value
      .map((item) => compactValue(item))
      .filter((item) => !isEmpty(item)) as unknown as T;
  }

  if (isPlainObject(value)) {
    const next: Record<string, unknown> = {};
    for (const [key, item] of Object.entries(value)) {
      const compacted = compactValue(item);
      if (!isEmpty(compacted)) next[key] = compacted;
    }
    return next as T;
  }

  return value;
}

function isEmpty(value: unknown): boolean {
  if (value == null) return true;
  if (typeof value === "string") return value.trim() === "";
  if (Array.isArray(value)) return value.length === 0;
  if (isPlainObject(value)) return Object.keys(value).length === 0;
  return false;
}

function isPlainObject(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function pad(value: number): string {
  return String(value).padStart(2, "0");
}
