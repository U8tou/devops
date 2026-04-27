<!-- 流程画布编辑（与基础信息编号、备注分离） -->
<template>
  <div class="process-edit art-full-height flex flex-col gap-3">
    <ElCard shadow="never" class="shrink-0">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="min-w-0 flex-1">
          <div class="truncate text-base font-medium text-[var(--el-text-color-primary)]">
            {{ headerCode || '—' }}
          </div>
          <div
            v-if="headerRemark"
            class="mt-0.5 truncate text-sm text-[var(--el-text-color-secondary)]"
          >
            {{ headerRemark }}
          </div>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <ElButton
            v-auth="'dev:process:edit'"
            :disabled="loading || !idParam"
            @click="openFlowEnvDialog"
          >
            {{ $t('dev.process.flowEnv.dialogOpen') }}
          </ElButton>
          <ElButton v-auth="'dev:process:edit'" @click="openMetaDialog">
            {{ $t('dev.process.editMeta') }}
          </ElButton>
          <ElTooltip
            :disabled="!hasUnsavedNodeEdits"
            :content="$t('dev.process.saveFlowNeedNodeSave')"
            placement="top"
          >
            <span class="inline-flex">
              <ElButton
                type="primary"
                :loading="saving"
                :disabled="loading || hasUnsavedNodeEdits || !flowIsDirty"
                @click="saveFlow"
              >
                {{ $t('dev.process.saveFlow') }}
              </ElButton>
            </span>
          </ElTooltip>
          <ElButton @click="back">{{ $t('dev.process.back') }}</ElButton>
          <ElDivider direction="vertical" class="!mx-1 h-6" />
          <ElButton :disabled="loading" @click="insertStandardPipeline">
            {{ $t('dev.workflow.insertStandardPipeline') }}
          </ElButton>
        </div>
      </div>
    </ElCard>

    <ElCard
      shadow="never"
      class="process-edit__flow-card min-h-0 flex flex-1 flex-col overflow-hidden"
    >
      <div
        class="process-edit__workspace flex min-h-0 flex-1 overflow-hidden rounded-lg border border-[var(--el-border-color)]"
      >
        <WorkflowPalette />
        <AutomationFlowCanvas
          ref="flowCanvasRef"
          v-model:nodes="nodes"
          v-model:edges="edges"
          :can-undo="canUndoFlow"
          @select="onFlowSelect"
          @flow-snapshot-refresh="onFlowSnapshotRefresh"
          @reset-flow="onToolbarResetFlow"
          @clear="onToolbarClear"
          @undo="onToolbarUndo"
        />
        <div class="process-edit__node-config">
          <NodeConfigPanel
            :node="selectedNodeForPanel"
            :nodes="nodes"
            :edges="edges"
            @update:data="onPatchNodeData"
            @session-draft="onNodeSessionDraft"
          />
        </div>
      </div>
    </ElCard>

    <ProcessFlowEnvDialog v-model="flowEnvVisible" :process-id="idParam" />
    <ProcessMetaDialog
      :key="`flow-${idParam}`"
      v-model="metaVisible"
      mode="edit"
      :process-id="idParam"
      :initial-code="headerCode"
      :initial-remark="headerRemark"
      :initial-cron-enabled="headerCronEnabled"
      :initial-cron-expr="headerCronExpr"
      :initial-tags="headerTags"
      @saved="onMetaSavedFromFlow"
    />
  </div>
</template>

<script setup lang="ts">
  import { computed, nextTick, provide } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useRoute, useRouter } from 'vue-router'
  import type { Edge, Node, ViewportTransform } from '@vue-flow/core'
  import { fetchDevProcessDetail, fetchDevProcessEditFlow } from '@/api/dev-process'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import '@vue-flow/core/dist/style.css'
  import '@vue-flow/core/dist/theme-default.css'
  import WorkflowPalette from './modules/WorkflowPalette.vue'
  import AutomationFlowCanvas from './modules/AutomationFlowCanvas.vue'
  import NodeConfigPanel from '@/components/devops/automation/NodeConfigPanel.vue'
  import ProcessFlowEnvDialog from './modules/process-flow-env.vue'
  import ProcessMetaDialog from '../process/modules/process-meta-dialog.vue'
  import {
    KIND_TO_TITLE_KEY,
    normalizeNodeKind,
    type AutomationKind
  } from './modules/automation-kinds'
  import type { AutomationNodePersistData } from './modules/automation-node-params'
  import { findDuplicateGitCheckoutSubdir, mergeNodeParams } from './modules/automation-node-params'
  import { buildStandardPipeline } from './modules/standard-pipeline'
  import { computeFlowSkippedNodeIds } from './modules/flow-graph-skip'

  defineOptions({ name: 'DevopsProcessEdit' })

  const { t } = useI18n()
  const route = useRoute()
  const router = useRouter()

  const loading = ref(false)
  const saving = ref(false)
  const headerCode = ref('')
  const headerRemark = ref('')
  const headerCronEnabled = ref('0')
  const headerCronExpr = ref('')
  const headerTags = ref<{ id: string; name: string }[]>([])
  const metaVisible = ref(false)
  const flowEnvVisible = ref(false)

  const nodes = ref<Node[]>([])
  const edges = ref<Edge[]>([])

  /** 与最近一次「加载 / 成功保存」时的 flow JSON 一致则为 false，用于禁用「保存流程」 */
  const lastSavedFlowPayload = ref<string | null>(null)
  /** 视口平移/缩放不修改 nodes/edges，依赖此项触发脏检查重算 */
  const flowSnapshotRev = ref(0)

  function onFlowSnapshotRefresh() {
    flowSnapshotRev.value += 1
  }

  const flowIsDirty = computed(() => {
    void flowSnapshotRev.value
    const last = lastSavedFlowPayload.value
    if (last === null) return false
    return flowPayload() !== last
  })

  /** 撤回栈：仅记录画布 nodes/edges */
  const flowUndoStack = ref<{ nodes: Node[]; edges: Edge[] }[]>([])
  const flowHistorySkip = ref(false)
  const flowHistoryApplying = ref(false)
  const canUndoFlow = computed(() => flowUndoStack.value.length > 1)

  const flowCanvasRef = ref<{
    getViewportForSave: () => ViewportTransform | null
    applyViewport: (v: ViewportTransform) => void
    resetView?: () => void
  } | null>(null)

  function cloneFlowState(): { nodes: Node[]; edges: Edge[] } {
    return {
      nodes: JSON.parse(JSON.stringify(nodes.value)) as Node[],
      edges: JSON.parse(JSON.stringify(edges.value)) as Edge[]
    }
  }

  function resetFlowHistory() {
    flowUndoStack.value = [cloneFlowState()]
  }

  function pushFlowHistoryIfChanged() {
    if (flowHistorySkip.value || flowHistoryApplying.value) return
    const snap = cloneFlowState()
    const last = flowUndoStack.value[flowUndoStack.value.length - 1]
    if (last && JSON.stringify(last) === JSON.stringify(snap)) return
    flowUndoStack.value.push(snap)
    if (flowUndoStack.value.length > 40) flowUndoStack.value.shift()
  }

  let flowHistoryTimer: ReturnType<typeof setTimeout> | null = null
  watch(
    [nodes, edges],
    () => {
      if (flowHistoryTimer) clearTimeout(flowHistoryTimer)
      flowHistoryTimer = setTimeout(() => {
        flowHistoryTimer = null
        pushFlowHistoryIfChanged()
      }, 450)
    },
    { deep: true }
  )

  async function onToolbarResetFlow() {
    try {
      await ElMessageBox.confirm(t('dev.workflow.toolbar.confirmResetFlow'), t('common.tips'), {
        type: 'warning',
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel')
      })
    } catch {
      return
    }
    const id = idParam.value
    if (!id) return
    try {
      const d = await fetchDevProcessDetail(id)
      await applyFlow(d.flow || '')
      selectedNodeId.value = null
      ElMessage.success(t('dev.workflow.toolbar.resetFlowDone'))
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : t('dev.process.loadFailed'))
    }
  }

  async function onToolbarClear() {
    try {
      await ElMessageBox.confirm(t('dev.workflow.toolbar.confirmClear'), t('common.tips'), {
        type: 'warning',
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel')
      })
    } catch {
      return
    }
    flowHistorySkip.value = true
    clearNodeSessionDrafts()
    nodes.value = []
    edges.value = []
    selectedNodeId.value = null
    await nextTick()
    resetFlowHistory()
    flowHistorySkip.value = false
    ElMessage.success(t('dev.workflow.toolbar.clearDone'))
  }

  function onToolbarUndo() {
    if (flowUndoStack.value.length < 2) return
    flowHistoryApplying.value = true
    flowUndoStack.value.pop()
    const prev = flowUndoStack.value[flowUndoStack.value.length - 1]
    if (prev) {
      clearNodeSessionDrafts()
      nodes.value = JSON.parse(JSON.stringify(prev.nodes)) as Node[]
      edges.value = JSON.parse(JSON.stringify(prev.edges)) as Edge[]
    }
    selectedNodeId.value = null
    nextTick(() => {
      flowHistoryApplying.value = false
      flowCanvasRef.value?.resetView?.()
    })
  }

  const selectedNodeId = ref<string | null>(null)
  provide('flowSelectedNodeId', selectedNodeId)
  provide('flowNodes', nodes)
  provide('flowEdges', edges)

  /** 侧栏编辑中的未落盘参数（与 nodes 中已保存的 data 对比） */
  const nodeSessionDraft = ref<Record<string, AutomationNodePersistData>>({})

  const selectedNodeForPanel = computed(() => {
    if (!selectedNodeId.value) return null
    const sid = String(selectedNodeId.value)
    const n = nodes.value.find((x) => String(x.id) === sid) ?? null
    if (!n) return null
    const d = nodeSessionDraft.value[sid]
    if (d) {
      return { ...n, data: d } as Node
    }
    return n
  })

  function nodePersistEqual(a: AutomationNodePersistData, b: unknown): boolean {
    return JSON.stringify(a) === JSON.stringify(b)
  }

  const hasUnsavedNodeEdits = computed(() => {
    for (const n of nodes.value) {
      const id = String(n.id)
      const d = nodeSessionDraft.value[id]
      if (d == null) continue
      if (nodePersistEqual(d, n.data)) continue
      return true
    }
    return false
  })

  const flowNodeDirtyIds = computed(() => {
    const r: string[] = []
    for (const n of nodes.value) {
      const id = String(n.id)
      const d = nodeSessionDraft.value[id]
      if (d == null) continue
      if (nodePersistEqual(d, n.data)) continue
      r.push(id)
    }
    return r
  })

  provide('flowNodeDirtyIds', flowNodeDirtyIds)

  const flowSkippedIdSet = computed(() =>
    computeFlowSkippedNodeIds(nodes.value, edges.value, nodeSessionDraft.value)
  )
  provide('flowSkippedIdSet', flowSkippedIdSet)

  function onNodeSessionDraft(payload: { id: string; data: AutomationNodePersistData }) {
    const id = String(payload.id)
    const n = nodes.value.find((x) => String(x.id) === id)
    if (!n) return
    if (nodePersistEqual(payload.data, n.data)) {
      const next = { ...nodeSessionDraft.value }
      delete next[id]
      nodeSessionDraft.value = next
    } else {
      nodeSessionDraft.value = { ...nodeSessionDraft.value, [id]: payload.data }
    }
  }

  function clearNodeSessionDrafts() {
    nodeSessionDraft.value = {}
  }

  const idParam = computed(() => {
    const raw = route.params.id
    return typeof raw === 'string' && raw.length > 0 ? raw : ''
  })

  function onFlowSelect(id: string | null) {
    selectedNodeId.value = id == null ? null : String(id)
  }

  function onPatchNodeData(payload: { id: string; data: AutomationNodePersistData }) {
    const pid = String(payload.id)
    nodes.value = nodes.value.map((n) => (String(n.id) === pid ? { ...n, data: payload.data } : n))
    const ne = { ...nodeSessionDraft.value }
    delete ne[pid]
    nodeSessionDraft.value = ne
  }

  watch(
    nodes,
    (list) => {
      if (
        selectedNodeId.value &&
        !list.some((n) => String(n.id) === String(selectedNodeId.value))
      ) {
        selectedNodeId.value = null
      }
    },
    { deep: true }
  )

  async function insertStandardPipeline() {
    if (nodes.value.length > 0) {
      try {
        await ElMessageBox.confirm(
          t('dev.workflow.insertStandardPipelineConfirm'),
          t('common.tips'),
          {
            type: 'warning',
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel')
          }
        )
      } catch {
        return
      }
    }
    const built = buildStandardPipeline(t)
    flowHistorySkip.value = true
    clearNodeSessionDrafts()
    nodes.value = built.nodes
    edges.value = built.edges
    selectedNodeId.value = null
    nextTick(() => {
      resetFlowHistory()
      flowHistorySkip.value = false
    })
  }

  /** 将边端点与节点 id 统一为 string，避免 number / string 混用导致选中态与图算法失效 */
  function normalizeFlowEdges(list: Edge[]): Edge[] {
    return list.map((e) => ({
      ...e,
      source: String(e.source),
      target: String(e.target)
    }))
  }

  /** 兼容旧版保存的节点（无 type / 无 kind / 无 params） */
  function normalizeLoadedNodes(list: Node[]): Node[] {
    return list.map((n, i) => {
      const d = n.data as Partial<AutomationNodePersistData> | undefined
      const rawKind = String(d?.kind ?? 'git_repo')
      const kind = normalizeNodeKind(rawKind) as AutomationKind
      const titleKey = KIND_TO_TITLE_KEY[kind]
      const label =
        d?.label?.trim() || (titleKey ? t(`dev.workflow.${titleKey}`) : '') || `Node ${i + 1}`
      return {
        ...n,
        id: String(n.id ?? `automation-${i}`),
        type: 'automation' as const,
        data: {
          kind,
          label,
          flowEnabled: d?.flowEnabled !== false,
          params: mergeNodeParams(rawKind, d?.params)
        }
      }
    })
  }

  function flowPayload(): string {
    const viewport = flowCanvasRef.value?.getViewportForSave?.() ?? undefined
    return JSON.stringify({
      nodes: nodes.value,
      edges: edges.value,
      ...(viewport ? { viewport } : {})
    })
  }

  async function applyFlow(raw: string) {
    flowHistorySkip.value = true
    nodes.value = []
    edges.value = []
    if (!raw || !raw.trim()) {
      await nextTick()
      clearNodeSessionDrafts()
      resetFlowHistory()
      flowHistorySkip.value = false
      return
    }
    try {
      const o = JSON.parse(raw) as {
        nodes?: Node[]
        edges?: Edge[]
        viewport?: ViewportTransform
      }
      if (Array.isArray(o.nodes) && o.nodes.length) {
        nodes.value = normalizeLoadedNodes(o.nodes)
      }
      if (Array.isArray(o.edges)) edges.value = normalizeFlowEdges(o.edges)
      await nextTick()
      await nextTick()
      if (o.viewport && flowCanvasRef.value) {
        flowCanvasRef.value.applyViewport(o.viewport)
      } else {
        await nextTick()
        flowCanvasRef.value?.resetView?.()
      }
    } catch {
      /* keep cleared */
    } finally {
      await nextTick()
      clearNodeSessionDrafts()
      resetFlowHistory()
      flowHistorySkip.value = false
      await nextTick()
      const syncLastSaved = () => {
        lastSavedFlowPayload.value = flowPayload()
      }
      syncLastSaved()
      // resetView / 视口动画可能晚于本帧，多拍一次，避免基线与当前 payload 暂时不一致
      setTimeout(syncLastSaved, 0)
      setTimeout(syncLastSaved, 320)
    }
  }

  async function load() {
    const id = idParam.value
    if (!id) {
      router.replace({ name: 'DevopsProcessList' })
      return
    }
    loading.value = true
    try {
      const d = await fetchDevProcessDetail(id)
      headerCode.value = d.code
      headerRemark.value = d.remark || ''
      headerCronEnabled.value = d.cronEnabled === '1' ? '1' : '0'
      headerCronExpr.value = d.cronExpr || ''
      headerTags.value = d.tags ?? []
      await applyFlow(d.flow || '')
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : t('dev.process.loadFailed'))
    } finally {
      loading.value = false
    }
  }

  async function saveFlow() {
    const id = idParam.value
    if (!id) return
    if (hasUnsavedNodeEdits.value) {
      ElMessage.error(t('dev.process.saveFlowNeedNodeSave'))
      return
    }
    const dup = findDuplicateGitCheckoutSubdir(nodes.value)
    if (dup !== null) {
      ElMessage.error(
        t('dev.workflow.gitCheckoutDuplicate', {
          path: dup === '' ? t('dev.workflow.gitCheckoutPathWorkspaceRoot') : dup
        })
      )
      return
    }
    if (!flowIsDirty.value) return

    const flow = flowPayload()
    saving.value = true
    try {
      await fetchDevProcessEditFlow({ id, flow })
      lastSavedFlowPayload.value = flow
      ElMessage.success(t('dev.process.saveSuccess'))
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : t('sys.common.operateFailed'))
    } finally {
      saving.value = false
    }
  }

  function openFlowEnvDialog() {
    if (!idParam.value) return
    flowEnvVisible.value = true
  }

  function openMetaDialog() {
    metaVisible.value = true
  }

  function onMetaSavedFromFlow() {
    const id = idParam.value
    if (!id) return
    fetchDevProcessDetail(id)
      .then((d) => {
        headerCode.value = d.code
        headerRemark.value = d.remark || ''
        headerCronEnabled.value = d.cronEnabled === '1' ? '1' : '0'
        headerCronExpr.value = d.cronExpr || ''
        headerTags.value = d.tags ?? []
      })
      .catch(() => {
        /* ignore */
      })
  }

  async function back() {
    if (hasUnsavedNodeEdits.value) {
      try {
        await ElMessageBox.confirm(t('dev.process.backUnsavedConfirm'), t('common.tips'), {
          type: 'warning',
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel')
        })
      } catch {
        return
      }
    }
    router.push({ name: 'DevopsProcessList' })
  }

  watch(idParam, () => load(), { immediate: true })
</script>

<style scoped>
  .process-edit__flow-card :deep(.el-card__body) {
    display: flex;
    flex: 1;
    flex-direction: column;
    min-height: 0;
    overflow: hidden;
  }

  .process-edit__workspace {
    align-items: stretch;
    height: min(640px, calc(100vh - 280px));
    min-height: 420px;
  }

  /**
   * 右侧参数栏：在 flex 行内占满交叉轴高度；勿用 height:100%（部分环境下百分比未解析，侧栏子 flex 区成 0 高）。
   */
  .process-edit__node-config {
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    align-self: stretch;
    width: 300px;
    min-width: 300px;
    max-width: 300px;
    min-height: 0;
    overflow: hidden;
  }
</style>
