<!-- 执行态流程节点（与编辑态区分：显示运行状态） -->
<template>
  <div
    class="run-workflow-node rounded-lg border bg-[var(--el-bg-color)] shadow-sm min-w-[180px] max-w-[240px] transition-colors"
    :class="statusClass"
    :style="{ borderLeftWidth: '3px', borderLeftColor: accent }"
  >
    <Handle class="run-workflow-node__handle" type="target" :position="Position.Left" />
    <div class="relative px-3 py-2 pr-8">
      <div
        class="absolute right-2 top-2 flex h-5 w-5 items-center justify-center rounded-full text-[10px] font-bold"
        :class="badgeClass"
      >
        <span v-if="status === 'running'" class="run-wf-spin">⟳</span>
        <span v-else-if="status === 'success'">✓</span>
        <span v-else-if="status === 'error'">✕</span>
        <span v-else-if="status === 'skipped'" title="skipped">⊘</span>
        <span v-else class="opacity-40">○</span>
      </div>
      <div class="mb-0.5 flex items-center gap-1.5">
        <ArtSvgIcon :icon="icon" class="text-base shrink-0" :style="{ color: accent }" />
        <span class="truncate text-sm font-medium text-[var(--el-text-color-primary)]">{{
          labelText
        }}</span>
      </div>
      <p class="text-[11px] text-[var(--el-text-color-secondary)]">{{ statusLabel }}</p>
      <p
        v-if="timingStartLine"
        class="mt-0.5 text-[10px] leading-tight text-[var(--el-text-color-secondary)]"
      >
        {{ timingStartLine }}
      </p>
      <p
        v-if="timingDurationLine"
        class="mt-0.5 text-[10px] leading-tight text-[var(--el-text-color-secondary)]"
      >
        {{ timingDurationLine }}
      </p>
      <p
        v-if="transferSizeLine"
        class="mt-0.5 text-[10px] leading-tight text-[var(--el-text-color-secondary)]"
      >
        {{ transferSizeLine }}
      </p>
    </div>
    <Handle class="run-workflow-node__handle" type="source" :position="Position.Right" />
  </div>
</template>

<script setup lang="ts">
  import { computed, nextTick, onMounted, watch } from 'vue'
  import { Handle, Position, useVueFlow } from '@vue-flow/core'
  import type { NodeProps } from '@vue-flow/core'
  import { useI18n } from 'vue-i18n'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'
  import {
    findPaletteDef,
    KIND_TO_TITLE_KEY,
    normalizeNodeKind,
    type AutomationKind
  } from '../../process-edit/modules/automation-kinds'
  import type { AutomationNodePersistData } from '../../process-edit/modules/automation-node-params'

  export type RunStatus = 'pending' | 'running' | 'success' | 'error' | 'skipped'

  export type RunWorkflowNodeData = AutomationNodePersistData & {
    runStatus?: RunStatus
    /** 节点开始执行时间（Unix 毫秒），来自 SSE node_start */
    runStartedAtMs?: number
    /** 节点执行耗时（毫秒），来自 SSE node_end */
    runDurationMs?: number
    /** 上传/下载字节数，来自 SSE node_end transferBytes */
    runTransferBytes?: number
  }

  /** 与后端 devflow.formatByteSize 一致：二进制 KiB/MiB… */
  function formatByteSize(n: number): string {
    if (!Number.isFinite(n) || n < 0) return ''
    if (n < 1024) return `${Math.floor(n)} B`
    let x = n / 1024
    let u = 0
    const units = ['KiB', 'MiB', 'GiB', 'TiB']
    while (x >= 1024 && u < units.length - 1) {
      x /= 1024
      u++
    }
    return `${x.toFixed(2)} ${units[u]}`
  }

  function formatRunStartMs(ms: number): string {
    const d = new Date(ms)
    return d.toLocaleString(undefined, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    })
  }

  function formatRunDurationMs(ms: number): string {
    if (!Number.isFinite(ms) || ms < 0) return ''
    if (ms < 1000) return `${Math.round(ms)}ms`
    return `${(ms / 1000).toFixed(2)}s`
  }

  const props = defineProps<NodeProps<RunWorkflowNodeData>>()

  const { t } = useI18n()
  const { updateNodeInternals } = useVueFlow()

  const status = computed(() => props.data?.runStatus ?? 'pending')

  /** 状态切换后重算连线锚点（避免边框/内容变化后句柄漂移） */
  watch(status, () => {
    nextTick(() => updateNodeInternals([props.id]))
  })

  onMounted(() => {
    nextTick(() => updateNodeInternals([props.id]))
  })

  const def = computed(() => findPaletteDef(props.data?.kind ?? ''))
  const accent = computed(() => def.value?.accent ?? '#64748b')
  const icon = computed(() => def.value?.icon ?? 'ri:flow-chart')

  const labelText = computed(() => {
    if (props.data?.label?.trim()) return props.data.label
    const nk = normalizeNodeKind(String(props.data?.kind ?? '')) as AutomationKind
    const tk = KIND_TO_TITLE_KEY[nk]
    return tk ? t(`dev.workflow.${tk}`) : nk
  })

  const statusLabel = computed(() => {
    switch (status.value) {
      case 'running':
        return t('dev.process.runStatus.running')
      case 'success':
        return t('dev.process.runStatus.success')
      case 'error':
        return t('dev.process.runStatus.error')
      case 'skipped':
        return t('dev.process.runStatus.skipped')
      default:
        return t('dev.process.runStatus.pending')
    }
  })

  const timingStartLine = computed(() => {
    const ms = props.data?.runStartedAtMs
    if (ms == null || !Number.isFinite(ms)) return ''
    return t('dev.process.runNodeStartAt', { time: formatRunStartMs(ms) })
  })

  const timingDurationLine = computed(() => {
    if (status.value === 'skipped') return ''
    if (status.value !== 'success' && status.value !== 'error') return ''
    const ms = props.data?.runDurationMs
    if (ms == null || !Number.isFinite(ms)) return ''
    return t('dev.process.runNodeDuration', { duration: formatRunDurationMs(ms) })
  })

  const transferSizeLine = computed(() => {
    if (status.value === 'skipped') return ''
    if (status.value !== 'success' && status.value !== 'error') return ''
    const b = props.data?.runTransferBytes
    if (b == null || !Number.isFinite(b) || b <= 0) return ''
    return t('dev.process.runNodeTransfer', { size: formatByteSize(b) })
  })

  const statusClass = computed(() => {
    switch (status.value) {
      case 'running':
        return 'run-workflow-node--running border-[var(--el-color-primary-light-5)]'
      case 'success':
        return 'run-workflow-node--ok border-[var(--el-color-success-light-5)]'
      case 'error':
        return 'run-workflow-node--err border-[var(--el-color-danger-light-5)]'
      case 'skipped':
        return 'run-workflow-node--skipped border-[var(--el-border-color-darker)] opacity-80'
      default:
        return 'border-[var(--el-border-color)] opacity-90'
    }
  })

  const badgeClass = computed(() => {
    switch (status.value) {
      case 'running':
        return 'bg-[var(--el-color-primary-light-7)] text-[var(--el-color-primary)]'
      case 'success':
        return 'bg-[var(--el-color-success-light-7)] text-[var(--el-color-success)]'
      case 'error':
        return 'bg-[var(--el-color-danger-light-7)] text-[var(--el-color-danger)]'
      case 'skipped':
        return 'bg-[var(--el-fill-color-dark)] text-[var(--el-text-color-secondary)]'
      default:
        return 'bg-[var(--el-fill-color)] text-[var(--el-text-color-secondary)]'
    }
  })
</script>

<style scoped>
  .run-workflow-node__handle {
    width: 8px;
    height: 8px;
    background: var(--el-bg-color);
    border-width: 2px;
  }

  .run-workflow-node--running {
    box-shadow: 0 0 0 1px rgb(64 158 255 / 35%);
    animation: runWfPulse 1.2s ease-in-out infinite;
  }

  @keyframes runWfPulse {
    0%,
    100% {
      box-shadow: 0 0 0 1px rgb(64 158 255 / 25%);
    }

    50% {
      box-shadow: 0 0 0 4px rgb(64 158 255 / 15%);
    }
  }

  .run-wf-spin {
    display: inline-block;
    animation: runWfSpin 0.9s linear infinite;
  }

  @keyframes runWfSpin {
    to {
      transform: rotate(-360deg);
    }
  }
</style>
