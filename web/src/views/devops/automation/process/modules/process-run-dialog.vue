<!-- 仿 n8n：上可视化流转，下控制台（5:3）；执行日志由 SSE 实时推送 -->
<template>
  <ElDialog
    v-model="dialogVisible"
    :title="dialogTitle"
    width="min(960px, 94vw)"
    top="4vh"
    destroy-on-close
    append-to-body
    class="process-run-dialog"
    @closed="onDialogClosed"
    @opened="onRunDialogOpened"
  >
    <div
      v-if="!flowJson?.trim()"
      class="py-8 text-center text-sm text-[var(--el-text-color-secondary)]"
    >
      {{ $t('dev.process.runEmptyFlow') }}
    </div>
    <div
      v-else
      class="process-run-dialog__body grid min-h-0 grid-rows-[minmax(0,5fr)_minmax(0,3fr)] overflow-hidden rounded-lg border border-[var(--el-border-color)] bg-[var(--el-fill-color-blank)]"
      :style="{ height: bodyHeight }"
    >
      <!-- 上 5：可视化 -->
      <div
        class="process-run-dialog__flow row-start-1 flex min-h-0 flex-col border-b border-[var(--el-border-color)]"
      >
        <div class="shrink-0 border-b border-[var(--el-border-color-lighter)] px-2 py-1.5">
          <span class="text-xs text-[var(--el-text-color-secondary)]">{{
            $t('dev.process.runFlowCaption')
          }}</span>
        </div>
        <div class="min-h-0 flex-1">
          <ProcessRunFlow
            :key="runFlowMountKey"
            v-model:nodes="runNodes"
            v-model:edges="runEdges"
            :layout-nonce="flowLayoutNonce"
          />
        </div>
      </div>
      <!-- 下 3：控制台 -->
      <div class="process-run-dialog__console row-start-2 flex min-h-0 flex-col">
        <div class="shrink-0 border-b border-[var(--el-border-color-lighter)] px-2 py-1.5">
          <span class="text-xs text-[var(--el-text-color-secondary)]">{{
            $t('dev.process.runConsoleCaption')
          }}</span>
        </div>
        <div
          ref="consoleRef"
          class="process-run-dialog__console-out min-h-0 flex-1 overflow-y-auto scroll-smooth py-2 font-mono text-[12px] leading-relaxed"
        >
          <div
            v-for="(block, idx) in consoleLogBlocks"
            :key="idx"
            class="whitespace-pre-wrap break-all px-3 py-1.5 text-[#a3e635]"
            :class="idx % 2 === 0 ? 'bg-[rgba(255,255,255,0.04)]' : 'bg-[rgba(255,255,255,0.09)]'"
          >
            {{ block.join('\n') }}
          </div>
          <div v-if="runPhase === 'running'" class="mt-1 px-3 animate-pulse text-[#86efac]">_</div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex w-full items-center justify-between gap-2">
        <span v-if="runPhase === 'done'" class="text-xs text-[var(--el-color-success)]">{{
          $t('dev.process.runDone')
        }}</span>
        <span v-else-if="runPhase === 'error'" class="text-xs text-[var(--el-color-danger)]">{{
          runError || $t('dev.process.runFailed')
        }}</span>
        <span v-else class="text-xs text-[var(--el-text-color-secondary)]">{{
          $t('dev.process.runSseHint')
        }}</span>
        <div class="ml-auto flex flex-wrap items-center justify-end gap-2">
          <ElButton v-if="runPhase === 'running'" type="danger" plain @click="cancelRun">
            {{ $t('dev.process.runStop') }}
          </ElButton>
          <ElButton v-if="canRestartRun" type="warning" plain @click="restartRun">
            {{ $t('dev.process.runRestart') }}
          </ElButton>
          <ElButton type="primary" @click="dialogVisible = false">{{
            $t('dev.process.runClose')
          }}</ElButton>
        </div>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { computed, nextTick, ref, watch, withDefaults } from 'vue'
  import { useI18n } from 'vue-i18n'
  import '@vue-flow/core/dist/style.css'
  import '@vue-flow/core/dist/theme-default.css'
  import type { Edge, Node, ViewportTransform } from '@vue-flow/core'
  import ProcessRunFlow from './process-run-flow.vue'
  import { splitExecLogIntoBlocks } from './exec-log-blocks'
  import { normalizeNodeKind } from '../../process-edit/modules/automation-kinds'
  import type { AutomationNodePersistData } from '../../process-edit/modules/automation-node-params'
  import { mergeNodeParams } from '../../process-edit/modules/automation-node-params'
  import type { RunStatus } from './run-workflow-node.vue'
  import { fetchDevProcessRunCancel } from '@/api/dev-process'
  import { ElMessageBox } from 'element-plus'
  import { useUserStore } from '@/store/modules/user'

  const props = withDefaults(
    defineProps<{
      flowJson: string
      /** 流程编号，用于标题 */
      processCode: string
      /** 流程主键，用于 SSE run_stream */
      processId?: string
    }>(),
    { processId: '' }
  )

  const dialogVisible = defineModel<boolean>({ required: true })

  const emit = defineEmits<{
    /** 弹窗完全关闭后（含点关闭、点 X），便于父级刷新列表等 */
    closed: []
  }>()

  const { t } = useI18n()
  const userStore = useUserStore()

  const runNodes = ref<Node[]>([])
  const runEdges = ref<Edge[]>([])
  const runFlowMountKey = ref(0)
  /** 弹窗动画结束 / 流程数据挂载后递增，迫使 Vue Flow 在真实容器尺寸下重算连线 */
  const flowLayoutNonce = ref(0)
  const logLines = ref<string[]>([])
  const runPhase = ref<'idle' | 'running' | 'done' | 'error'>('idle')
  const runError = ref('')
  const consoleRef = ref<HTMLElement | null>(null)
  /** 用户终止或服务端「已取消」后的 idle，才允许与「执行完成」一样点重启；非首次打开前的空闲 */
  const allowRestartAfterStop = ref(false)

  let eventSource: EventSource | null = null
  /** 正常结束或已处理错误，避免 EventSource onerror 误判 */
  let sseFinished = false

  const bodyHeight = 'min(72vh, 680px)'

  const canRestartRun = computed(
    () =>
      !!props.flowJson?.trim() &&
      !!props.processId?.trim() &&
      (runPhase.value === 'done' || (runPhase.value === 'idle' && allowRestartAfterStop.value))
  )

  const dialogTitle = computed(() => t('dev.process.runTitle', { code: props.processCode || '—' }))

  /** 控制台按节点分段栅栏色（与「查看日志」一致） */
  const consoleLogBlocks = computed(() => splitExecLogIntoBlocks(logLines.value.join('\n')))

  function closeEventSource() {
    if (eventSource) {
      eventSource.close()
      eventSource = null
    }
  }

  function buildRunStreamUrl(): string {
    const token = userStore.accessToken || ''
    const base = (import.meta.env.VITE_API_URL as string | undefined)?.replace(/\/$/, '') ?? ''
    const q = `id=${encodeURIComponent(props.processId)}&token=${encodeURIComponent(token)}`
    if (base) {
      return `${base}/api/dev_process/run_stream?${q}`
    }
    return `/api/dev_process/run_stream?${q}`
  }

  function onDialogClosed() {
    cancelRun()
    logLines.value = []
    runNodes.value = []
    runEdges.value = []
    runPhase.value = 'idle'
    runError.value = ''
    allowRestartAfterStop.value = false
    flowLayoutNonce.value = 0
    emit('closed')
  }

  function onRunDialogOpened() {
    nextTick(() => {
      requestAnimationFrame(() => {
        flowLayoutNonce.value += 1
      })
    })
  }

  function cancelRun() {
    closeEventSource()
    if (props.processId) {
      fetchDevProcessRunCancel(props.processId).catch(() => {})
    }
    if (runPhase.value === 'running') {
      appendLog(t('dev.process.runCancelled'))
      allowRestartAfterStop.value = true
      runPhase.value = 'idle'
      sseFinished = true
    }
  }

  async function restartRun() {
    try {
      await ElMessageBox.confirm(t('dev.process.runRestartConfirm'), t('common.tips'), {
        type: 'warning',
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel')
      })
    } catch {
      return
    }
    sseFinished = false
    cancelRun()
    runPhase.value = 'idle'
    runError.value = ''
    await startRunFromSse()
  }

  function parseFlow(raw: string): {
    nodes: Node[]
    edges: Edge[]
    viewport: ViewportTransform | null
  } {
    try {
      const o = JSON.parse(raw) as {
        nodes?: Node[]
        edges?: Edge[]
        viewport?: ViewportTransform
      }
      const rawNodes = Array.isArray(o.nodes) ? o.nodes : []
      const nodes = rawNodes.map((n) => {
        const d = n.data as Partial<AutomationNodePersistData> | undefined
        const rawKind = String(d?.kind ?? 'git_repo')
        const kind = normalizeNodeKind(rawKind)
        return {
          ...n,
          position: n.position ?? { x: 0, y: 0 },
          type: 'automation' as const,
          data: {
            kind,
            label: d?.label ?? '',
            flowEnabled: d?.flowEnabled !== false,
            params: mergeNodeParams(rawKind, d?.params),
            runStatus: 'pending' as RunStatus
          }
        }
      })
      const edges = Array.isArray(o.edges) ? o.edges : []
      const vp = o.viewport
      const viewport =
        vp && typeof vp === 'object' && typeof vp.zoom === 'number'
          ? ({ x: vp.x ?? 0, y: vp.y ?? 0, zoom: vp.zoom } as ViewportTransform)
          : null
      return { nodes, edges, viewport }
    } catch {
      return { nodes: [], edges: [], viewport: null }
    }
  }

  /** 合并节点执行态（状态 + 启动时间 + 耗时，与 SSE 一致） */
  function applyNodeRun(
    id: string,
    patch: Partial<{
      runStatus: RunStatus
      runStartedAtMs: number
      runDurationMs: number | undefined
      /** 上传/下载字节数，仅 upload_servers / remote_download 等节点可能有值 */
      runTransferBytes: number | undefined
    }>
  ) {
    runNodes.value = runNodes.value.map((n) => {
      if (n.id !== id) return n
      const prev = (n.data as Record<string, unknown>) ?? {}
      return {
        ...n,
        data: {
          ...prev,
          ...patch
        }
      }
    })
  }

  function appendLog(line: string) {
    logLines.value = [...logLines.value, line]
    nextTick(() => {
      const el = consoleRef.value
      if (el) el.scrollTop = el.scrollHeight
    })
  }

  async function startRunFromSse() {
    sseFinished = false
    allowRestartAfterStop.value = false
    runPhase.value = 'running'
    runError.value = ''
    logLines.value = []

    const { nodes, edges } = parseFlow(props.flowJson)
    runFlowMountKey.value += 1
    runEdges.value = edges
    runNodes.value = nodes
    await nextTick()
    flowLayoutNonce.value += 1

    if (!nodes.length) {
      runPhase.value = 'done'
      appendLog(t('dev.process.runEmptyFlow'))
      return
    }

    if (!props.processId?.trim()) {
      runPhase.value = 'error'
      runError.value = t('dev.process.runMissingProcessId')
      return
    }

    appendLog(t('dev.process.runStarted', { code: props.processCode || '—' }))
    closeEventSource()

    const es = new EventSource(buildRunStreamUrl())
    eventSource = es

    es.addEventListener('progress', (ev: Event) => {
      const me = ev as MessageEvent
      if (!me.data) return
      try {
        const p = JSON.parse(me.data as string) as {
          type: string
          nodeId?: string
          kind?: string
          line?: string
          ok?: boolean
          startedAtMs?: number
          durationMs?: number
          transferBytes?: number
          skipped?: boolean
        }
        if (p.type === 'log') {
          appendLog(p.line ?? '')
        } else if (p.type === 'node_start' && p.nodeId) {
          applyNodeRun(p.nodeId, {
            runStatus: 'running',
            runStartedAtMs: p.startedAtMs ?? Date.now(),
            runDurationMs: undefined,
            runTransferBytes: undefined
          })
        } else if (p.type === 'node_end' && p.nodeId) {
          const tb = p.transferBytes
          if (p.skipped) {
            applyNodeRun(p.nodeId, {
              runStatus: 'skipped',
              runDurationMs: undefined,
              runTransferBytes: undefined
            })
          } else {
            applyNodeRun(p.nodeId, {
              runStatus: p.ok ? 'success' : 'error',
              runDurationMs: p.durationMs,
              runTransferBytes:
                typeof tb === 'number' && Number.isFinite(tb) && tb > 0 ? tb : undefined
            })
          }
        }
      } catch {
        /* ignore */
      }
    })

    es.addEventListener('done', (ev: Event) => {
      const me = ev as MessageEvent
      sseFinished = true
      try {
        const d = JSON.parse((me.data as string) || '{}') as { ok?: boolean; error?: string }
        if (d.ok) {
          runPhase.value = 'done'
        } else {
          runPhase.value = 'error'
          runError.value = d.error || ''
        }
      } catch {
        runPhase.value = 'done'
      }
      closeEventSource()
    })

    es.addEventListener('cancelled', () => {
      sseFinished = true
      appendLog(t('dev.process.runCancelled'))
      allowRestartAfterStop.value = true
      runPhase.value = 'idle'
      closeEventSource()
    })

    es.addEventListener('run_error', (ev: Event) => {
      const me = ev as MessageEvent
      sseFinished = true
      try {
        const d = JSON.parse((me.data as string) || '{}') as { ok?: boolean; error?: string }
        runError.value = d.error || ''
        runPhase.value = 'error'
      } catch {
        runPhase.value = 'error'
      }
      closeEventSource()
    })

    es.onerror = () => {
      if (sseFinished) return
      if (eventSource?.readyState === EventSource.CLOSED) {
        if (!sseFinished) {
          runPhase.value = 'error'
          runError.value = t('dev.process.runSseFailed')
        }
        return
      }
      runPhase.value = 'error'
      runError.value = t('dev.process.runSseFailed')
      closeEventSource()
    }
  }

  watch(dialogVisible, (v) => {
    if (v && props.flowJson?.trim()) {
      void startRunFromSse()
    }
  })
</script>

<style scoped>
  .process-run-dialog__console-out {
    color: #bef264;
    background: linear-gradient(180deg, #0f172a 0%, #020617 100%);
  }
</style>

<style>
  .process-run-dialog .el-dialog__body {
    padding-top: 8px;
  }
</style>
