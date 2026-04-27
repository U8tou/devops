<!-- 自动化流程节点（n8n 风格：左进右出 + 参数摘要） -->
<template>
  <div
    class="automation-node relative rounded-lg border border-[var(--el-border-color)] bg-[var(--el-bg-color)] shadow-sm min-w-[180px] max-w-[240px] transition-shadow"
    :class="{
      'automation-node--selected': isHighlighted,
      'automation-node--flow-off': isFlowSkipped
    }"
    :style="{ borderLeftWidth: '3px', borderLeftColor: accent }"
  >
    <div
      v-if="hasUnsavedInPanel"
      class="automation-node__unsaved-badge"
      :title="t('dev.process.nodeUnsavedHint')"
      role="img"
      :aria-label="t('dev.process.nodeUnsavedHint')"
    >
      <span class="automation-node__unsaved-ico" aria-hidden="true">!</span>
    </div>
    <Handle class="automation-node__handle" type="target" :position="Position.Left" />
    <div class="px-3 py-2">
      <div class="mb-1 flex items-center gap-1.5">
        <ArtSvgIcon :icon="icon" class="text-base shrink-0" :style="{ color: accent }" />
        <span class="truncate text-sm font-medium text-[var(--el-text-color-primary)]">{{
          labelText
        }}</span>
      </div>
      <p
        v-if="summaryText"
        class="text-xs leading-snug text-[var(--el-text-color-secondary)] line-clamp-2"
      >
        {{ summaryText }}
      </p>
      <p
        v-else-if="hintText"
        class="text-xs leading-snug text-[var(--el-text-color-secondary)] line-clamp-2"
      >
        {{ hintText }}
      </p>
    </div>
    <Handle class="automation-node__handle" type="source" :position="Position.Right" />
  </div>
</template>

<script setup lang="ts">
  import { computed, inject, type ComputedRef, type Ref } from 'vue'
  import type { Edge, Node } from '@vue-flow/core'
  import { Handle, Position } from '@vue-flow/core'
  import type { NodeProps } from '@vue-flow/core'
  import { useI18n } from 'vue-i18n'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'
  import {
    findPaletteDef,
    KIND_TO_TITLE_KEY,
    normalizeNodeKind,
    type AutomationKind
  } from './automation-kinds'
  import type {
    AutomationNodePersistData,
    ExecuteScriptParams,
    GitRepoParams,
    RemoteDownloadParams,
    RemoteSshScriptParams,
    SshConnectionParams,
    UploadServersParams
  } from './automation-node-params'
  import { mergeNodeParams, normalizeSshPort } from './automation-node-params'
  import { listUpstreamSshBranches, resolveEffectiveSshForRemoteDownload } from './flow-graph-ssh'

  const props = defineProps<NodeProps<AutomationNodePersistData>>()

  const flowSelectedId = inject<Ref<string | null> | undefined>('flowSelectedNodeId', undefined)
  const flowNodes = inject<Ref<Node[]> | undefined>('flowNodes', undefined)
  const flowEdges = inject<Ref<Edge[]> | undefined>('flowEdges', undefined)
  const flowNodeDirtyIds = inject<ComputedRef<string[]>>('flowNodeDirtyIds')
  const flowSkippedIdSet = inject<ComputedRef<Set<string>> | undefined>(
    'flowSkippedIdSet',
    undefined
  )

  const isFlowSkipped = computed(() => flowSkippedIdSet?.value.has(String(props.id)) ?? false)

  const hasUnsavedInPanel = computed(() => {
    const ids = flowNodeDirtyIds?.value
    if (!ids?.length) return false
    return ids.includes(String(props.id))
  })
  const isHighlighted = computed(
    () =>
      props.selected ||
      (flowSelectedId?.value != null && String(flowSelectedId.value) === String(props.id))
  )

  const { t } = useI18n()

  const normalizedKind = computed(() => normalizeNodeKind(String(props.data?.kind ?? '')))

  const def = computed(() => findPaletteDef(normalizedKind.value))
  const accent = computed(() => def.value?.accent ?? '#64748b')
  const icon = computed(() => def.value?.icon ?? 'ri:flow-chart')

  const labelText = computed(() => {
    if (props.data?.label?.trim()) return props.data.label
    const nk = normalizedKind.value as AutomationKind
    const tk = KIND_TO_TITLE_KEY[nk]
    return tk ? t(`dev.workflow.${tk}`) : nk
  })

  const hintText = computed(() => {
    const nk = normalizedKind.value
    if (!nk) return ''
    const key = `dev.workflow.hints.${nk}`
    const s = t(key)
    return s === key ? '' : s
  })

  function shortUrl(url: string) {
    const u = url
      .replace(/\.git$/i, '')
      .split(/[/:]/)
      .filter(Boolean)
    return u.slice(-2).join('/') || url
  }

  function sshBrief(eff: SshConnectionParams) {
    const h = eff.host?.trim() || '—'
    const portStr = normalizeSshPort(eff.port)
    const port = portStr === '22' ? '' : `:${portStr}`
    const user = eff.username?.trim() || ''
    return user ? `${user}@${h}${port}` : `${h}${port}`
  }

  const summaryText = computed(() => {
    const data = props.data
    if (!data?.params || !data.kind) return ''
    const nk = normalizedKind.value
    switch (nk) {
      case 'git_repo': {
        const p = data.params as GitRepoParams
        return `${shortUrl(p.repositoryUrl)} · ${p.branch}`
      }
      case 'ssh_connection': {
        const p = data.params as SshConnectionParams
        const host = p.host?.trim() || '—'
        const portStr = normalizeSshPort(p.port)
        const port = portStr === '22' ? '' : `:${portStr}`
        const user = p.username?.trim() || ''
        return user ? `${user}@${host}${port}` : `${host}${port}`
      }
      case 'remote_ssh_script': {
        const p = mergeNodeParams('remote_ssh_script', data.params) as RemoteSshScriptParams
        const scriptHead = (p.script.trim().split('\n')[0] ?? '').slice(0, 48)
        const cwdPart = p.cwd?.trim() && p.cwd !== '.' ? `${p.cwd}: ` : ''
        const ns = flowNodes?.value
        const es = flowEdges?.value
        if (ns && es) {
          const branches = listUpstreamSshBranches(props.id, ns, es)
          if (branches.length === 0) {
            return `${t('dev.workflow.remoteScriptNoUpstream')} · ${cwdPart}${scriptHead}`
          }
          if (branches.length === 1) {
            return `${t('dev.workflow.remoteScriptUpstream')} ${sshBrief(branches[0].ssh)} · ${cwdPart}${scriptHead}`
          }
          return `${t('dev.workflow.remoteScriptMultiTargets', { n: branches.length })} · ${cwdPart}${scriptHead}`
        }
        return `${t('dev.workflow.remoteScriptUpstreamShort')} · ${cwdPart}${scriptHead}`
      }
      case 'execute_script': {
        const p = data.params as ExecuteScriptParams
        const line = p.script.trim().split('\n')[0] ?? ''
        const head = line.slice(0, 88)
        return p.cwd && p.cwd !== '.' ? `${p.cwd}: ${head}` : head
      }
      case 'upload_servers': {
        const p = mergeNodeParams('upload_servers', data.params) as UploadServersParams
        const host = p.host?.trim() || '—'
        const portStr = normalizeSshPort(p.port)
        const port = portStr === '22' ? '' : `:${portStr}`
        const user = p.username?.trim() || ''
        const target = user ? `${user}@${host}${port}` : `${host}${port}`
        return `${target} → ${p.remotePath}`
      }
      case 'remote_download': {
        const p = mergeNodeParams('remote_download', data.params) as RemoteDownloadParams
        const src = (p.remoteSourcePath || '').trim() || '—'
        const out = (p.outputPath || '').trim() || '—'
        const ns = flowNodes?.value
        const es = flowEdges?.value
        if (ns && es) {
          const eff = resolveEffectiveSshForRemoteDownload(props.id, ns, es)
          if (!eff.host?.trim()) {
            return `${t('dev.workflow.remoteScriptNoUpstream')} · ${src} → ${out}`
          }
          return `${t('dev.workflow.remoteScriptUpstream')} ${sshBrief(eff)} · ${src} → ${out}`
        }
        return `${t('dev.workflow.remoteScriptUpstreamShort')} · ${src} → ${out}`
      }
      default:
        return ''
    }
  })
</script>

<style scoped>
  .automation-node__unsaved-badge {
    position: absolute;
    top: 5px;
    right: 5px;
    left: auto;
    z-index: 2;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 18px;
    height: 18px;
    color: #fff;

    /* 黄色警告，与主题 warning 色一致、在深浅背景下可读性较好 */
    background: var(--el-color-warning);
    border-radius: 9999px;
    box-shadow: 0 0 0 2px var(--el-bg-color);
  }

  .automation-node__unsaved-ico {
    font-size: 12px;
    font-weight: 800;
    line-height: 1;
  }

  .automation-node__handle {
    width: 8px;
    height: 8px;
    background: var(--el-bg-color);
    border-width: 2px;
  }

  .automation-node--selected {
    box-shadow:
      0 0 0 2px var(--el-color-primary-light-5),
      0 4px 12px rgb(0 0 0 / 8%);
  }

  /** 流程开关关闭：本节点及下游不执行 */
  .automation-node--flow-off {
    filter: grayscale(0.92);
    border-color: var(--el-border-color-darker) !important;
    box-shadow: none;
    opacity: 0.55;
  }

  .automation-node--flow-off .automation-node__handle {
    opacity: 0.6;
  }
</style>
