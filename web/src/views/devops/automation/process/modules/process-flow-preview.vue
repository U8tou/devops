<!-- 只读流程预览（Vue Flow） -->
<template>
  <div class="process-flow-preview">
    <VueFlow
      v-model:nodes="nodes"
      v-model:edges="edges"
      :node-types="nodeTypes"
      :nodes-draggable="false"
      :nodes-connectable="false"
      :elements-selectable="false"
      :zoom-on-scroll="true"
      fit-view-on-init
      class="process-flow-preview__canvas"
    >
      <Background pattern-color="#aaa" :gap="16" />
    </VueFlow>
  </div>
</template>

<script setup lang="ts">
  import { markRaw, provide, ref, watch } from 'vue'
  import { VueFlow } from '@vue-flow/core'
  import { Background } from '@vue-flow/background'
  import type { Edge, Node } from '@vue-flow/core'
  import AutomationNode from '../../process-edit/modules/AutomationNode.vue'
  import type { AutomationNodePersistData } from '../../process-edit/modules/automation-node-params'
  import { mergeNodeParams } from '../../process-edit/modules/automation-node-params'
  import { normalizeNodeKind } from '../../process-edit/modules/automation-kinds'
  import '@vue-flow/core/dist/style.css'
  import '@vue-flow/core/dist/theme-default.css'

  const props = defineProps<{
    flowJson: string
  }>()

  const nodes = ref<Node[]>([])
  const edges = ref<Edge[]>([])
  provide('flowNodes', nodes)
  provide('flowEdges', edges)

  const nodeTypes = {
    automation: markRaw(AutomationNode)
  }

  watch(
    () => props.flowJson,
    (raw) => {
      try {
        const o = JSON.parse(raw) as { nodes?: Node[]; edges?: Edge[] }
        const rawNodes = Array.isArray(o.nodes) ? o.nodes : []
        nodes.value = rawNodes.map((n) => {
          const d = n.data as Partial<AutomationNodePersistData> | undefined
          const rawKind = String(d?.kind ?? 'git_repo')
          const kind = normalizeNodeKind(rawKind)
          return {
            ...n,
            type: 'automation' as const,
            data: {
              kind,
              label: d?.label ?? '',
              params: mergeNodeParams(rawKind, d?.params)
            }
          }
        })
        edges.value = Array.isArray(o.edges) ? o.edges : []
      } catch {
        nodes.value = []
        edges.value = []
      }
    },
    { immediate: true }
  )
</script>

<style scoped>
  .process-flow-preview {
    height: 280px;
    overflow: hidden;
    background: var(--el-fill-color-blank);
    border: 1px solid var(--el-border-color);
    border-radius: 8px;
  }

  .process-flow-preview__canvas {
    width: 100%;
    height: 280px;
  }
</style>
