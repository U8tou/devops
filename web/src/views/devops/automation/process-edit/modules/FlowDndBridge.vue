<!-- 必须在 VueFlow 子树内使用 useVueFlow，用于将拖放坐标转换为画布坐标 -->
<script setup lang="ts">
  import { useVueFlow } from '@vue-flow/core'
  import { nextTick, onUnmounted, watch } from 'vue'

  const emit = defineEmits<{
    dropNode: [payload: { kind: string; label: string; position: { x: number; y: number } }]
  }>()

  const { screenToFlowCoordinate, vueFlowRef } = useVueFlow()

  /** 捕获阶段：避免被子元素（连线层、节点层）吞掉，保证能 preventDefault 允许 drop */
  const CAPTURE = true

  function handleDragEnter(e: DragEvent) {
    e.preventDefault()
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault()
    if (e.dataTransfer) e.dataTransfer.dropEffect = 'copy'
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault()
    e.stopPropagation()
    const raw = e.dataTransfer?.getData('application/json') || e.dataTransfer?.getData('text/plain')
    if (!raw) return
    let parsed: { kind: string; label: string }
    try {
      parsed = JSON.parse(raw) as { kind: string; label: string }
    } catch {
      return
    }
    if (!parsed.kind) return
    const position = screenToFlowCoordinate({ x: e.clientX, y: e.clientY })
    emit('dropNode', {
      kind: parsed.kind,
      label: parsed.label || '',
      position
    })
  }

  let cleanup: (() => void) | undefined

  function bindToFlowRoot(el: HTMLElement) {
    cleanup?.()
    cleanup = undefined
    el.addEventListener('dragenter', handleDragEnter, CAPTURE)
    el.addEventListener('dragover', handleDragOver, CAPTURE)
    el.addEventListener('drop', handleDrop, CAPTURE)
    cleanup = () => {
      el.removeEventListener('dragenter', handleDragEnter, CAPTURE)
      el.removeEventListener('dragover', handleDragOver, CAPTURE)
      el.removeEventListener('drop', handleDrop, CAPTURE)
    }
  }

  watch(
    vueFlowRef,
    (el, _prev, onCleanup) => {
      cleanup?.()
      cleanup = undefined
      if (!el) return
      nextTick(() => bindToFlowRoot(el))
      onCleanup(() => cleanup?.())
    },
    { immediate: true, flush: 'post' }
  )

  onUnmounted(() => {
    cleanup?.()
  })
</script>

<template>
  <span class="hidden" aria-hidden="true" />
</template>
