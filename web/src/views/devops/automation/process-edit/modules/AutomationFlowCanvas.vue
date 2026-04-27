<!-- 画布：自定义节点 + 拖放添加 -->
<template>
  <div class="automation-flow-canvas relative h-full min-h-0 w-full min-w-0 flex-1">
    <div
      class="pointer-events-none absolute left-3 top-3 z-10 flex flex-wrap items-center gap-1.5"
      aria-label="canvas-toolbar"
    >
      <div
        class="pointer-events-auto flex flex-wrap items-center gap-1.5 rounded-md border border-[var(--el-border-color)] bg-[var(--el-bg-color)]/95 px-1.5 py-1 shadow-sm backdrop-blur-sm"
      >
        <ElButton size="small" @click="onResetView">
          {{ $t('dev.workflow.toolbar.resetView') }}
        </ElButton>
        <ElButton size="small" @click="emit('resetFlow')">
          {{ $t('dev.workflow.toolbar.resetFlow') }}
        </ElButton>
        <ElButton size="small" :disabled="!props.canUndo" @click="emit('undo')">
          {{ $t('dev.workflow.toolbar.undo') }}
        </ElButton>
        <ElButton size="small" type="danger" plain @click="emit('clear')">
          {{ $t('dev.workflow.toolbar.clear') }}
        </ElButton>
      </div>
    </div>
    <VueFlow
      v-model:nodes="nodesModel"
      v-model:edges="edgesModel"
      :node-types="nodeTypes"
      :default-viewport="{ x: 0, y: 0, zoom: 1 }"
      :min-zoom="0.2"
      :max-zoom="2"
      :fit-view-on-init="false"
      :edges-updatable="true"
      :delete-key-code="deleteKeys"
      class="h-full w-full"
      @connect="onConnect"
      @edge-update="onEdgeUpdate"
      @node-click="onNodeClick"
      @pane-click="onPaneClick"
      @viewport-change-end="emit('flowSnapshotRefresh')"
    >
      <FlowViewportSync ref="viewportSyncRef" />
      <FlowDndBridge @drop-node="onDropNode" />
      <Background pattern-color="#aaa" :gap="16" />
      <Controls />
    </VueFlow>
  </div>
</template>

<script setup lang="ts">
  import { markRaw, ref } from 'vue'
  import { VueFlow } from '@vue-flow/core'
  import { Background } from '@vue-flow/background'
  import { Controls } from '@vue-flow/controls'
  import type { Connection, Edge, EdgeUpdateEvent, Node, ViewportTransform } from '@vue-flow/core'
  import AutomationNode from './AutomationNode.vue'
  import FlowDndBridge from './FlowDndBridge.vue'
  import FlowViewportSync from './FlowViewportSync.vue'
  import { defaultParams } from './automation-node-params'
  import type { AutomationNodePersistData } from './automation-node-params'
  import type { AutomationKind } from './automation-kinds'

  const props = withDefaults(
    defineProps<{
      /** 是否可撤回（由父级维护历史栈） */
      canUndo?: boolean
    }>(),
    { canUndo: false }
  )

  const emit = defineEmits<{
    select: [nodeId: string | null]
    /** 视口平移/缩放结束，供父级重算与上次保存的 flow 是否一致 */
    flowSnapshotRefresh: []
    resetFlow: []
    clear: []
    undo: []
  }>()

  const viewportSyncRef = ref<{
    snapshotForSave: () => ViewportTransform | null
    restoreViewport: (v: ViewportTransform) => void
    resetView: () => void
  } | null>(null)

  const nodesModel = defineModel<Node[]>('nodes', { required: true })
  const edgesModel = defineModel<Edge[]>('edges', { required: true })

  /** 默认仅 Backspace；补充 Delete 以符合 Windows 习惯 */
  const deleteKeys: string[] = ['Backspace', 'Delete']

  const nodeTypes = {
    automation: markRaw(AutomationNode)
  }

  function onConnect(c: Connection) {
    const edge: Edge = {
      id: `e${c.source}-${c.target}-${Date.now()}`,
      source: c.source,
      target: c.target,
      sourceHandle: c.sourceHandle ?? undefined,
      targetHandle: c.targetHandle ?? undefined,
      updatable: true
    }
    edgesModel.value = [...edgesModel.value, edge]
  }

  function onEdgeUpdate({ edge, connection }: EdgeUpdateEvent) {
    edgesModel.value = edgesModel.value.map((e) =>
      e.id === edge.id
        ? {
            ...e,
            source: connection.source,
            target: connection.target,
            sourceHandle: connection.sourceHandle ?? undefined,
            targetHandle: connection.targetHandle ?? undefined
          }
        : e
    )
  }

  function onDropNode(payload: {
    kind: string
    label: string
    position: { x: number; y: number }
  }) {
    const kind = payload.kind as AutomationKind
    const id = `n-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`
    const data: AutomationNodePersistData = {
      kind,
      label: payload.label,
      flowEnabled: true,
      params: defaultParams(kind)
    }
    const node: Node = {
      id,
      type: 'automation',
      position: payload.position,
      data
    }
    nodesModel.value = [...nodesModel.value, node]
    emit('select', id)
  }

  function onNodeClick({ node }: { node: Node }) {
    emit('select', String(node.id))
  }

  function onPaneClick() {
    emit('select', null)
  }

  function getViewportForSave() {
    return viewportSyncRef.value?.snapshotForSave() ?? null
  }

  function applyViewport(v: ViewportTransform) {
    viewportSyncRef.value?.restoreViewport(v)
  }

  function onResetView() {
    viewportSyncRef.value?.resetView?.()
  }

  defineExpose({
    getViewportForSave,
    applyViewport,
    resetView: onResetView
  })
</script>

<style scoped>
  .automation-flow-canvas :deep(.vue-flow) {
    width: 100%;
    height: 100%;
  }
</style>
