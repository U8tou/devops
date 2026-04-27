<!-- 置于 VueFlow 子树内，用于读写视口（与 flow JSON 一并持久化） -->
<script setup lang="ts">
  import { useVueFlow } from '@vue-flow/core'
  import type { ViewportTransform } from '@vue-flow/core'

  const { toObject, setViewport, fitView } = useVueFlow()

  /** 将视口适配到全部节点 */
  function resetView() {
    fitView({ padding: 0.15, duration: 250 })
  }

  function snapshotForSave(): ViewportTransform | null {
    try {
      const o = toObject()
      return o.viewport ?? null
    } catch {
      return null
    }
  }

  function restoreViewport(v: ViewportTransform) {
    setViewport(v)
  }

  defineExpose({
    snapshotForSave,
    restoreViewport,
    resetView
  })
</script>

<template>
  <span class="hidden" aria-hidden="true" />
</template>
