<!-- 执行弹窗内：只读流程图（节点带 runStatus）；视口自动缩放以包含全部节点 -->
<template>
  <div class="process-run-flow h-full w-full min-h-[200px]">
    <VueFlow
      v-model:nodes="nodesModel"
      v-model:edges="edgesModel"
      :node-types="nodeTypes"
      :default-viewport="{ x: 0, y: 0, zoom: 1 }"
      :default-edge-options="defaultEdgeOptions"
      :nodes-draggable="false"
      :nodes-connectable="false"
      :elements-selectable="false"
      :zoom-on-scroll="true"
      :min-zoom="0.06"
      :max-zoom="1.5"
      :fit-view-on-init="false"
      class="process-run-flow__canvas"
    >
      <ProcessRunFlowFit :layout-nonce="layoutNonce" />
      <Background pattern-color="#94a3b8" :gap="14" :size="1" />
    </VueFlow>
  </div>
</template>

<script setup lang="ts">
  import { markRaw } from 'vue'
  import { VueFlow } from '@vue-flow/core'
  import { Background } from '@vue-flow/background'
  import type { DefaultEdgeOptions, Edge, Node } from '@vue-flow/core'
  import RunWorkflowNode from './run-workflow-node.vue'
  import ProcessRunFlowFit from './process-run-flow-fit.vue'

  defineProps<{
    /** 弹窗 opened / 数据就绪时递增，驱动在真实尺寸下重算连线 */
    layoutNonce?: number
  }>()

  const nodesModel = defineModel<Node[]>('nodes', { required: true })
  const edgesModel = defineModel<Edge[]>('edges', { required: true })

  /** 正交折线，多入边时锚点更稳定 */
  const defaultEdgeOptions: DefaultEdgeOptions = {
    type: 'smoothstep'
  }

  const nodeTypes = {
    automation: markRaw(RunWorkflowNode)
  }
</script>

<style scoped>
  .process-run-flow__canvas {
    width: 100%;
    height: 100%;
  }

  .process-run-flow :deep(.vue-flow) {
    width: 100%;
    height: 100%;
  }
</style>
