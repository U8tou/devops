<!-- 置于 VueFlow 子树内：弹窗/容器尺寸稳定后再 fit + 刷新句柄（避免首帧 pane 为 0 导致边错位） -->
<script setup lang="ts">
  import { nextTick, watch } from 'vue'
  import { useResizeObserver, useDebounceFn } from '@vueuse/core'
  import { useVueFlow, useNodesInitialized } from '@vue-flow/core'

  const props = withDefaults(
    defineProps<{
      /** 由外层 ElDialog @opened 或数据就绪后递增，强制重算 */
      layoutNonce?: number
    }>(),
    {
      layoutNonce: 0
    }
  )

  const { fitView, updateNodeInternals, nodes, dimensions, vueFlowRef } = useVueFlow()
  const nodesInitialized = useNodesInitialized()

  function allNodeIds(): string[] {
    return nodes.value.map((n) => n.id)
  }

  async function doFullRefit() {
    const dim = dimensions.value
    if (!dim?.width || !dim?.height) return

    const ids = allNodeIds()
    if (!ids.length) return

    try {
      await fitView({
        padding: 0.18,
        duration: 0,
        maxZoom: 1,
        minZoom: 0.06,
        includeHiddenNodes: false
      })
    } catch {
      /* ignore */
    }

    await nextTick()
    updateNodeInternals(ids)
    await nextTick()
    updateNodeInternals(ids)
    requestAnimationFrame(() => {
      updateNodeInternals(ids)
      requestAnimationFrame(() => updateNodeInternals(ids))
    })
  }

  const debouncedRefit = useDebounceFn(() => {
    void doFullRefit()
  }, 80)

  /** 容器从隐藏变为可见、或弹窗动画结束后尺寸变化时触发 */
  useResizeObserver(vueFlowRef, (entries) => {
    const { width, height } = entries[0]?.contentRect ?? { width: 0, height: 0 }
    if (width > 32 && height > 32) {
      debouncedRefit()
    }
  })

  watch(
    () => [dimensions.value.width, dimensions.value.height] as const,
    ([w, h], prev) => {
      const [pw, ph] = prev ?? [0, 0]
      if (w > 0 && h > 0 && (pw === 0 || ph === 0)) {
        debouncedRefit()
      }
    }
  )

  watch(
    () => props.layoutNonce,
    () => {
      if (props.layoutNonce > 0) debouncedRefit()
    }
  )

  watch(
    [nodesInitialized, () => dimensions.value.width, () => dimensions.value.height],
    ([ready, w, h]) => {
      if (ready && w > 0 && h > 0) debouncedRefit()
    },
    { immediate: true }
  )
</script>

<template>
  <span class="hidden" aria-hidden="true" />
</template>
