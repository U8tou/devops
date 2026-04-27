<template>
  <div
    class="node-config-panel node-config-panel--root flex min-h-0 w-full flex-col overflow-hidden border-l border-[var(--el-border-color)] bg-[var(--el-fill-color-blank)]"
  >
    <div class="shrink-0 border-b border-[var(--el-border-color)] px-3 py-2">
      <div v-if="node" class="flex w-full min-w-0 items-center justify-between gap-2">
        <span class="min-w-0 text-sm font-medium text-[var(--el-text-color-primary)]">{{
          $t('dev.workflow.nodeSettings')
        }}</span>
        <div class="node-config-panel__flow-switch flex shrink-0 items-center gap-1.5" @click.stop>
          <span class="whitespace-nowrap text-xs text-[var(--el-text-color-secondary)]">{{
            $t('dev.workflow.flowEnableLabel')
          }}</span>
          <ElSwitch
            v-model="localFlowEnabled"
            size="small"
            :aria-label="$t('dev.workflow.flowEnableLabel')"
          />
        </div>
      </div>
      <template v-else>
        <span class="text-sm font-medium text-[var(--el-text-color-primary)]">{{
          $t('dev.workflow.nodeSettings')
        }}</span>
        <div class="mt-1 space-y-1.5 text-xs leading-relaxed text-[var(--el-text-color-secondary)]">
          <p class="m-0">{{ $t('dev.workflow.nodeSettingsEmpty') }}</p>
          <p class="m-0">{{ $t('dev.workflow.nodeSettingsEnvUsage') }}</p>
        </div>
      </template>
    </div>

    <div v-if="node" class="node-config-panel__main">
      <div
        ref="scrollPanelRef"
        class="node-config-panel__scroll px-3 py-3"
        data-testid="node-config-scroll"
      >
        <div class="node-config-panel__form-wrap">
          <ElForm
            class="node-config-panel__inner-form"
            label-position="top"
            size="small"
            @submit.prevent
          >
            <ElFormItem :label="$t('dev.workflow.fieldLabel')">
              <ElInput v-model="localLabel" maxlength="120" show-word-limit />
            </ElFormItem>

            <div v-if="nodeKind === 'remote_ssh_script'" class="node-config-panel__kind-block">
              <ul
                v-if="remoteSshBranches.length > 1"
                class="mb-2 list-inside list-disc space-y-0.5 rounded border border-[var(--el-border-color-lighter)] bg-[var(--el-fill-color-light)] px-2 py-1.5 text-xs leading-snug text-[var(--el-text-color-secondary)]"
              >
                <li v-for="b in remoteSshBranches" :key="b.sourceId">
                  <span class="font-medium text-[var(--el-text-color-primary)]">{{ b.label }}</span>
                  -> {{ sshLine(b.ssh) }}
                </li>
              </ul>
              <p
                v-if="remoteSshUpstreamHint"
                class="mb-2 rounded border border-[var(--el-color-primary-light-5)] bg-[var(--el-color-primary-light-9)] px-2 py-1.5 text-xs leading-snug text-[var(--el-text-color-secondary)]"
              >
                {{ remoteSshUpstreamHint }}
              </p>
              <p class="mb-2 text-xs leading-snug text-[var(--el-text-color-secondary)]">
                {{ $t('dev.workflow.remoteScriptSshHint') }}
              </p>
              <ElFormItem :label="$t('dev.workflow.params.remoteSshCwd')">
                <ElInput v-model="remoteSshScriptParams.cwd" />
              </ElFormItem>
              <ElFormItem>
                <template #label>
                  <div
                    class="node-config-panel__field-label-row flex w-full min-w-0 items-center justify-between gap-2"
                  >
                    <span class="min-w-0">{{ $t('dev.workflow.params.remoteSshScript') }}</span>
                    <ElButton
                      type="primary"
                      link
                      size="small"
                      class="shrink-0"
                      @click.stop="openScriptEditor('remote_ssh_script')"
                    >
                      <ElIcon class="node-config-panel__label-btn-icon mr-0.5 text-sm">
                        <FullScreen />
                      </ElIcon>
                      {{ $t('dev.workflow.openScriptEditor') }}
                    </ElButton>
                  </div>
                </template>
                <ElInput v-model="remoteSshScriptParams.script" type="textarea" :rows="8" />
              </ElFormItem>
            </div>

            <div v-else-if="nodeKind === 'ssh_connection'" class="node-config-panel__kind-block">
              <ElFormItem :label="$t('dev.workflow.params.sshHostPort')">
                <div class="node-config-panel__host-port-row flex gap-2">
                  <ElInput
                    v-model="sshParams.host"
                    class="min-w-0 flex-1"
                    :placeholder="$t('dev.workflow.params.sshHost')"
                  />
                  <ElInput
                    v-model="sshParams.port"
                    class="node-config-panel__port-input shrink-0"
                    placeholder="22"
                  />
                </div>
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.sshUsername')">
                <ElInput v-model="sshParams.username" />
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.sshAuthType')">
                <ElRadioGroup
                  v-model="sshParams.authType"
                  class="node-config-panel__radio-group w-full"
                >
                  <ElRadio value="password">{{
                    $t('dev.workflow.params.sshAuthPassword')
                  }}</ElRadio>
                  <ElRadio value="key">{{ $t('dev.workflow.params.sshAuthKey') }}</ElRadio>
                </ElRadioGroup>
              </ElFormItem>
              <ElFormItem
                v-if="sshParams.authType === 'password'"
                :label="$t('dev.workflow.params.sshPassword')"
              >
                <ElInput
                  v-model="sshParams.password"
                  type="password"
                  show-password
                  autocomplete="off"
                />
              </ElFormItem>
              <ElFormItem v-else :label="$t('dev.workflow.params.sshPrivateKey')">
                <ElInput v-model="sshParams.privateKey" type="textarea" :rows="6" />
              </ElFormItem>
            </div>

            <div v-else-if="nodeKind === 'git_repo'" class="node-config-panel__kind-block">
              <ElFormItem :label="$t('dev.workflow.params.repoUrl')">
                <ElInput v-model="gitRepo.repositoryUrl" type="textarea" :rows="2" />
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.branch')">
                <ElInput v-model="gitRepo.branch" />
              </ElFormItem>
              <ElFormItem>
                <template #label>
                  <span class="inline-flex items-center gap-1">
                    {{ $t('dev.workflow.params.checkoutSubdir') }}
                    <ElTooltip
                      :content="$t('dev.workflow.params.checkoutSubdirHint')"
                      placement="top"
                      :show-after="300"
                    >
                      <span
                        class="inline-flex cursor-help text-[var(--el-text-color-secondary)] hover:text-[var(--el-color-primary)]"
                        tabindex="0"
                        role="button"
                        :aria-label="$t('dev.workflow.params.checkoutSubdirHint')"
                      >
                        <ElIcon class="text-base leading-none"><QuestionFilled /></ElIcon>
                      </span>
                    </ElTooltip>
                  </span>
                </template>
                <ElInput v-model="gitRepo.checkoutSubdir" />
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.shallowClone')">
                <ElSwitch v-model="gitRepo.shallowClone" />
              </ElFormItem>
              <p class="mb-2 text-xs leading-snug text-[var(--el-text-color-secondary)]">
                {{ $t('dev.workflow.params.gitCredsHint') }}
              </p>
              <ElFormItem :label="$t('dev.workflow.params.gitAuthType')">
                <ElRadioGroup
                  v-model="gitRepo.gitAuthType"
                  class="node-config-panel__radio-group w-full"
                >
                  <ElRadio value="none">{{ $t('dev.workflow.params.gitAuthNone') }}</ElRadio>
                  <ElRadio value="http">{{ $t('dev.workflow.params.gitAuthHttp') }}</ElRadio>
                  <ElRadio value="ssh_key">{{ $t('dev.workflow.params.gitAuthSshKey') }}</ElRadio>
                </ElRadioGroup>
              </ElFormItem>
              <div v-if="gitRepo.gitAuthType === 'http'" class="node-config-panel__kind-block">
                <ElFormItem :label="$t('dev.workflow.params.gitHttpUsername')">
                  <ElInput v-model="gitRepo.httpUsername" autocomplete="off" />
                </ElFormItem>
                <ElFormItem :label="$t('dev.workflow.params.gitHttpPassword')">
                  <ElInput
                    v-model="gitRepo.httpPassword"
                    type="password"
                    show-password
                    autocomplete="off"
                  />
                </ElFormItem>
              </div>
              <div v-if="gitRepo.gitAuthType === 'ssh_key'" class="node-config-panel__kind-block">
                <ElFormItem :label="$t('dev.workflow.params.gitSshPrivateKey')">
                  <ElInput v-model="gitRepo.sshPrivateKey" type="textarea" :rows="6" />
                </ElFormItem>
              </div>
            </div>

            <div v-else-if="nodeKind === 'execute_script'" class="node-config-panel__kind-block">
              <ElFormItem :label="$t('dev.workflow.params.executeScriptCwd')">
                <ElInput v-model="executeScriptParams.cwd" />
              </ElFormItem>
              <ElFormItem>
                <template #label>
                  <div
                    class="node-config-panel__field-label-row flex w-full min-w-0 items-center justify-between gap-2"
                  >
                    <span class="min-w-0">{{ $t('dev.workflow.params.executeScriptBody') }}</span>
                    <ElButton
                      type="primary"
                      link
                      size="small"
                      class="shrink-0"
                      @click.stop="openScriptEditor('execute_script')"
                    >
                      <ElIcon class="node-config-panel__label-btn-icon mr-0.5 text-sm">
                        <FullScreen />
                      </ElIcon>
                      {{ $t('dev.workflow.openScriptEditor') }}
                    </ElButton>
                  </div>
                </template>
                <ElInput v-model="executeScriptParams.script" type="textarea" :rows="10" />
              </ElFormItem>
            </div>

            <div v-else-if="nodeKind === 'upload_servers'" class="node-config-panel__kind-block">
              <p class="mb-2 text-xs leading-snug text-[var(--el-text-color-secondary)]">
                {{ $t('dev.workflow.uploadDeploySshHint') }}
              </p>
              <ElFormItem :label="$t('dev.workflow.params.sshHostPort')">
                <div class="node-config-panel__host-port-row flex gap-2">
                  <ElInput
                    v-model="uploadParams.host"
                    class="min-w-0 flex-1"
                    :placeholder="$t('dev.workflow.params.sshHost')"
                  />
                  <ElInput
                    v-model="uploadParams.port"
                    class="node-config-panel__port-input shrink-0"
                    placeholder="22"
                  />
                </div>
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.sshUsername')">
                <ElInput v-model="uploadParams.username" />
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.sshAuthType')">
                <ElRadioGroup
                  v-model="uploadParams.authType"
                  class="node-config-panel__radio-group w-full"
                >
                  <ElRadio value="password">{{
                    $t('dev.workflow.params.sshAuthPassword')
                  }}</ElRadio>
                  <ElRadio value="key">{{ $t('dev.workflow.params.sshAuthKey') }}</ElRadio>
                </ElRadioGroup>
              </ElFormItem>
              <ElFormItem
                v-if="uploadParams.authType === 'password'"
                :label="$t('dev.workflow.params.sshPassword')"
              >
                <ElInput
                  v-model="uploadParams.password"
                  type="password"
                  show-password
                  autocomplete="off"
                />
              </ElFormItem>
              <ElFormItem v-else :label="$t('dev.workflow.params.sshPrivateKey')">
                <ElInput v-model="uploadParams.privateKey" type="textarea" :rows="6" />
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.remotePath')">
                <ElInput v-model="uploadParams.remotePath" />
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.artifactGlob')">
                <ElInput v-model="uploadParams.artifactGlob" />
              </ElFormItem>
            </div>

            <div v-else-if="nodeKind === 'remote_download'" class="node-config-panel__kind-block">
              <p
                v-if="remoteDownloadUpstreamHint"
                class="mb-2 rounded border border-[var(--el-color-primary-light-5)] bg-[var(--el-color-primary-light-9)] px-2 py-1.5 text-xs leading-snug text-[var(--el-text-color-secondary)]"
              >
                {{ remoteDownloadUpstreamHint }}
              </p>
              <p class="mb-2 text-xs leading-snug text-[var(--el-text-color-secondary)]">
                {{ $t('dev.workflow.remoteDownloadSshHint') }}
              </p>
              <ElFormItem :label="$t('dev.workflow.params.remoteSourcePath')">
                <ElInput v-model="remoteDownload.remoteSourcePath" type="textarea" :rows="2" />
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.outputPath')">
                <ElInput v-model="remoteDownload.outputPath" />
              </ElFormItem>
              <ElFormItem :label="$t('dev.workflow.params.extractArchive')">
                <ElSwitch v-model="remoteDownload.extract" />
              </ElFormItem>
            </div>
          </ElForm>
        </div>
      </div>
      <div
        class="node-config-panel__footer shrink-0 border-t border-[var(--el-border-color)] bg-[var(--el-bg-color)] px-3 py-2 shadow-[0_-6px_12px_-4px_rgba(0,0,0,0.08)]"
      >
        <div class="flex gap-2">
          <ElButton
            v-if="canQuickValidate"
            :loading="validating"
            class="min-w-0 flex-1"
            @click="validateConnection"
          >
            {{ $t('dev.workflow.validateConnection') }}
          </ElButton>
          <ElButton
            type="primary"
            :class="canQuickValidate ? 'min-w-0 flex-1' : 'w-full'"
            @click="saveNode"
          >
            {{ $t('dev.workflow.saveNode') }}
          </ElButton>
        </div>
      </div>
    </div>

    <ElDialog
      v-model="scriptEditorVisible"
      :title="scriptEditorDialogTitle"
      width="min(900px, 96vw)"
      append-to-body
      align-center
      class="node-config-panel__script-dialog"
      :close-on-click-modal="false"
      @closed="onScriptEditorClosed"
    >
      <ElInput
        v-model="scriptEditorDraft"
        type="textarea"
        :rows="24"
        class="node-config-panel__script-dialog-input w-full"
        spellcheck="false"
      />
      <template #footer>
        <ElButton @click="scriptEditorVisible = false">{{ $t('common.cancel') }}</ElButton>
        <ElButton type="primary" @click="applyScriptEditor">{{ $t('common.confirm') }}</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { computed, nextTick, reactive, ref, watch } from 'vue'
  import { FullScreen, QuestionFilled } from '@element-plus/icons-vue'
  import { ElMessage } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import type { Edge, Node } from '@vue-flow/core'
  import {
    normalizeNodeKind,
    type AutomationKind
  } from '@/views/devops/automation/process-edit/modules/automation-kinds'
  import {
    defaultParams,
    findDuplicateGitCheckoutSubdir,
    mergeNodeParams,
    normalizeSshPort,
    type AutomationNodePersistData,
    type ExecuteScriptParams,
    type GitRepoParams,
    type RemoteDownloadParams,
    type RemoteSshScriptParams,
    type SshConnectionParams,
    type UploadServersParams
  } from '@/views/devops/automation/process-edit/modules/automation-node-params'
  import {
    listUpstreamSshBranches,
    resolveUpstreamSsh
  } from '@/views/devops/automation/process-edit/modules/flow-graph-ssh'
  import { fetchDevProcessValidateNode } from '@/api/dev-process'

  const props = withDefaults(
    defineProps<{
      node: Node | null
      nodes?: Node[]
      edges?: Edge[]
    }>(),
    {
      nodes: () => [],
      edges: () => []
    }
  )

  const emit = defineEmits<{
    'update:data': [payload: { id: string; data: AutomationNodePersistData }]
    'session-draft': [payload: { id: string; data: AutomationNodePersistData }]
  }>()

  const { t } = useI18n()

  const scrollPanelRef = ref<HTMLElement | null>(null)

  type ScriptEditorTarget = 'execute_script' | 'remote_ssh_script'
  const scriptEditorVisible = ref(false)
  const scriptEditorTarget = ref<ScriptEditorTarget | null>(null)
  const scriptEditorDraft = ref('')

  const scriptEditorDialogTitle = computed(() => {
    const k = scriptEditorTarget.value
    if (k === 'remote_ssh_script') return t('dev.workflow.params.remoteSshScript')
    if (k === 'execute_script') return t('dev.workflow.params.executeScriptBody')
    return t('dev.workflow.openScriptEditor')
  })

  function openScriptEditor(target: ScriptEditorTarget) {
    scriptEditorTarget.value = target
    if (target === 'execute_script') {
      scriptEditorDraft.value = executeScriptParams.script
    } else {
      scriptEditorDraft.value = remoteSshScriptParams.script
    }
    scriptEditorVisible.value = true
  }

  function applyScriptEditor() {
    const k = scriptEditorTarget.value
    if (k === 'execute_script') {
      executeScriptParams.script = scriptEditorDraft.value
    } else if (k === 'remote_ssh_script') {
      remoteSshScriptParams.script = scriptEditorDraft.value
    }
    scriptEditorVisible.value = false
  }

  function onScriptEditorClosed() {
    scriptEditorTarget.value = null
    scriptEditorDraft.value = ''
  }

  const nodeKind = computed(
    () =>
      normalizeNodeKind(
        String((props.node?.data as AutomationNodePersistData | undefined)?.kind ?? 'git_repo')
      ) as AutomationKind
  )

  function sshLine(eff: SshConnectionParams) {
    const h = eff.host?.trim() || '-'
    const portStr = normalizeSshPort(eff.port)
    const port = portStr === '22' ? '' : `:${portStr}`
    const user = eff.username?.trim() || ''
    return user ? `${user}@${h}${port}` : `${h}${port}`
  }

  const localLabel = ref('')
  const localFlowEnabled = ref(true)

  function buildPayloadFromDraft(): AutomationNodePersistData {
    const k = nodeKind.value
    const paramsMap: Record<AutomationKind, () => AutomationNodePersistData['params']> = {
      git_repo: () => mergeNodeParams('git_repo', { ...gitRepo }),
      ssh_connection: () => mergeNodeParams('ssh_connection', { ...sshParams }),
      remote_ssh_script: () => mergeNodeParams('remote_ssh_script', { ...remoteSshScriptParams }),
      execute_script: () => mergeNodeParams('execute_script', { ...executeScriptParams }),
      upload_servers: () => mergeNodeParams('upload_servers', { ...uploadParams }),
      remote_download: () => mergeNodeParams('remote_download', { ...remoteDownload })
    }
    return {
      kind: k,
      label: localLabel.value,
      flowEnabled: localFlowEnabled.value,
      params: paramsMap[k]()
    }
  }

  function nodesWithDraftApplied(): Node[] {
    if (!props.node) return props.nodes
    const payload = buildPayloadFromDraft()
    const nid = String(props.node.id)
    return props.nodes.map((n) => (String(n.id) === nid ? { ...n, data: payload } : n))
  }

  function saveNode() {
    if (!props.node) return
    const dup = findDuplicateGitCheckoutSubdir(nodesWithDraftApplied())
    if (dup !== null) {
      ElMessage.error(
        t('dev.workflow.gitCheckoutDuplicate', {
          path: dup === '' ? t('dev.workflow.gitCheckoutPathWorkspaceRoot') : dup
        })
      )
      return
    }
    emit('update:data', {
      id: String(props.node.id),
      data: buildPayloadFromDraft()
    })
    ElMessage.success(t('dev.workflow.nodeSaveSuccess'))
  }

  const validating = ref(false)

  const canQuickValidate = computed(
    () =>
      nodeKind.value === 'git_repo' ||
      nodeKind.value === 'ssh_connection' ||
      nodeKind.value === 'upload_servers'
  )

  async function validateConnection() {
    if (!props.node || !canQuickValidate.value) return
    const payload = buildPayloadFromDraft()
    validating.value = true
    try {
      const res = await fetchDevProcessValidateNode({
        kind: payload.kind,
        params: payload.params as unknown as Record<string, unknown>
      })
      if (res.ok) {
        ElMessage.success(t('dev.workflow.validateConnectionOk'))
      } else {
        ElMessage.error(res.message || t('dev.workflow.validateConnectionFail'))
      }
    } finally {
      validating.value = false
    }
  }

  const gitRepo = reactive<GitRepoParams>({ ...defaultParams('git_repo') })
  const sshParams = reactive<SshConnectionParams>({ ...defaultParams('ssh_connection') })
  const remoteSshScriptParams = reactive<RemoteSshScriptParams>({
    ...defaultParams('remote_ssh_script')
  })
  const remoteSshBranches = computed(() => {
    if (!props.node || !props.nodes.length) return []
    return listUpstreamSshBranches(props.node.id, props.nodes, props.edges)
  })

  const remoteSshUpstreamHint = computed(() => {
    if (!props.node || !props.nodes.length) return ''
    const branches = remoteSshBranches.value
    if (branches.length === 0) {
      return t('dev.workflow.remoteScriptNoUpstream')
    }
    if (branches.length === 1) {
      return t('dev.workflow.remoteScriptResolvedSsh', { target: sshLine(branches[0].ssh) })
    }
    return t('dev.workflow.remoteScriptMultiRunHint', { count: branches.length })
  })
  const executeScriptParams = reactive<ExecuteScriptParams>({ ...defaultParams('execute_script') })
  const uploadParams = reactive<UploadServersParams>({ ...defaultParams('upload_servers') })
  const remoteDownload = reactive<RemoteDownloadParams>({ ...defaultParams('remote_download') })
  const remoteDownloadUpstreamHint = computed(() => {
    if (!props.node || !props.nodes.length) return ''
    const up = resolveUpstreamSsh(props.node.id, props.nodes, props.edges)
    if (up?.host?.trim()) {
      return t('dev.workflow.remoteDownloadResolvedSsh', { target: sshLine(up) })
    }
    return t('dev.workflow.remoteDownloadNoUpstream')
  })

  const draftJsonSig = computed(() => {
    if (!props.node) return null
    try {
      return JSON.stringify(buildPayloadFromDraft())
    } catch {
      return null
    }
  })

  // 须先于 session-draft 的 watch：用同步方式灌入 local，避免 async 时 draft 比表单回填更早、误标未保存
  watch(
    () => props.node?.id,
    () => {
      if (!props.node) return
      const d = props.node.data as AutomationNodePersistData | undefined
      const rawKind = String(d?.kind ?? 'git_repo')
      const nk = normalizeNodeKind(rawKind)
      localLabel.value = d?.label ?? ''
      localFlowEnabled.value = d?.flowEnabled !== false
      if (nk === 'remote_ssh_script') {
        Object.assign(
          remoteSshScriptParams,
          mergeNodeParams(rawKind, d?.params) as RemoteSshScriptParams
        )
      } else {
        Object.assign(remoteSshScriptParams, defaultParams('remote_ssh_script'))
      }
      if (nk === 'ssh_connection') {
        Object.assign(sshParams, mergeNodeParams(rawKind, d?.params) as SshConnectionParams)
      } else {
        Object.assign(sshParams, defaultParams('ssh_connection'))
      }
      if (nk === 'git_repo') {
        Object.assign(gitRepo, mergeNodeParams(rawKind, d?.params) as GitRepoParams)
      } else {
        Object.assign(gitRepo, defaultParams('git_repo'))
      }
      if (nk === 'execute_script') {
        Object.assign(
          executeScriptParams,
          mergeNodeParams(rawKind, d?.params) as ExecuteScriptParams
        )
      } else {
        Object.assign(executeScriptParams, defaultParams('execute_script'))
      }
      if (nk === 'upload_servers') {
        Object.assign(uploadParams, mergeNodeParams(rawKind, d?.params) as UploadServersParams)
      } else {
        Object.assign(uploadParams, defaultParams('upload_servers'))
      }
      if (nk === 'remote_download') {
        Object.assign(remoteDownload, mergeNodeParams(rawKind, d?.params) as RemoteDownloadParams)
      } else {
        Object.assign(remoteDownload, defaultParams('remote_download'))
      }
      nextTick(() => {
        const el = scrollPanelRef.value
        if (el) el.scrollTop = 0
      })
    },
    { immediate: true, flush: 'sync' }
  )

  watch(
    draftJsonSig,
    () => {
      if (!props.node || draftJsonSig.value == null) return
      emit('session-draft', {
        id: String(props.node.id),
        data: buildPayloadFromDraft()
      })
    },
    { flush: 'post' }
  )
</script>

<style scoped>
  /** 占满父级 process-edit__node-config（交叉轴 stretch，勿依赖子项 height:100%） */
  .node-config-panel--root {
    box-sizing: border-box;
    display: flex;
    flex: 1 1 0%;
    flex-direction: column;
    min-height: 0;
  }

  /**
   * 中间「表单滚动区 + 底栏」用 grid：比纯 flex 更稳定，避免生产环境下中间行被算成 0 高。
   */
  .node-config-panel__main {
    display: grid;
    flex: 1 1 0%;
    grid-template-rows: minmax(0, 1fr) auto;
    min-height: 0;
    overflow: hidden;
  }

  .node-config-panel__scroll {
    min-height: 0;
    overflow: hidden auto;
    -webkit-overflow-scrolling: touch;
  }

  .node-config-panel__footer {
    flex-shrink: 0;
  }

  .node-config-panel__form-wrap {
    width: 100%;
    min-height: min-content;
    font-size: var(--el-font-size-base);
    line-height: 1.5;
    color: var(--el-text-color-primary);
  }

  /** 须使用 ElForm 注入 label-position；纯 div 套类名无法让表单项顶对齐 */
  .node-config-panel__form-wrap :deep(.node-config-panel__inner-form.el-form--label-top) {
    display: block;
    width: 100%;
  }

  /** 全局 el-ui 对 label 固定高度会破坏顶对齐时的多行标签 */
  .node-config-panel__form-wrap :deep(.el-form--label-top .el-form-item__label) {
    height: auto !important;
    min-height: 0;
    padding-bottom: 4px;
    line-height: 1.45 !important;
  }

  /** 各节点类型字段区（真实 DOM，避免生产构建下 template 片段与 ElForm 插槽异常） */
  .node-config-panel__kind-block {
    width: 100%;
  }

  .node-config-panel__form-wrap :deep(.el-form-item) {
    margin-bottom: 14px;
  }

  .node-config-panel__radio-group.el-radio-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
    align-items: stretch;
    width: 100%;
    line-height: 1.45;
  }

  .node-config-panel__radio-group :deep(.el-radio) {
    box-sizing: border-box;
    display: flex;
    align-items: flex-start;
    width: 100%;
    height: auto;
    padding: 0;
    margin: 0;
  }

  .node-config-panel__radio-group :deep(.el-radio__input) {
    flex-shrink: 0;
    margin-top: 2px;
  }

  .node-config-panel__radio-group :deep(.el-radio__label) {
    padding-left: 8px;
    line-height: inherit;
    word-break: break-word;
    white-space: normal;
  }

  .node-config-panel__port-input {
    width: 88px;
  }

  .node-config-panel__label-btn-icon {
    vertical-align: -0.125em;
  }

  /** 弹窗大编辑区：等宽、足够高度，便于长脚本 */
  :deep(.node-config-panel__script-dialog-input .el-textarea__inner) {
    box-sizing: border-box;
    min-height: min(50vh, 640px);
    font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, 'Liberation Mono', monospace;
    line-height: 1.45;
  }
</style>
