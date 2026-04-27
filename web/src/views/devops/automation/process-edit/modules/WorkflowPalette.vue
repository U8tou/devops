<!-- 左侧可拖拽节点面板（仿 n8n） -->
<template>
  <div
    class="workflow-palette flex h-full min-h-0 w-[220px] shrink-0 flex-col overflow-hidden border-r border-[var(--el-border-color)] bg-[var(--el-fill-color-blank)]"
  >
    <div class="border-b border-[var(--el-border-color)] px-3 py-2">
      <span class="text-sm font-medium text-[var(--el-text-color-primary)]">{{
        $t('dev.workflow.paletteTitle')
      }}</span>
      <p class="mt-0.5 text-xs text-[var(--el-text-color-secondary)]">
        {{ $t('dev.workflow.dragHint') }}
      </p>
    </div>
    <div class="flex flex-1 flex-col gap-2 overflow-y-auto p-2">
      <div
        v-for="item in PALETTE_ITEMS"
        :key="item.kind"
        draggable="true"
        class="workflow-palette__item cursor-grab rounded-md border border-[var(--el-border-color-lighter)] bg-[var(--el-bg-color)] px-2.5 py-2 transition-shadow hover:shadow-sm active:cursor-grabbing"
        :style="{ borderLeft: `3px solid ${item.accent}` }"
        @dragstart="(e) => onDragStart(e, item)"
      >
        <div class="flex items-center gap-2">
          <ArtSvgIcon :icon="item.icon" class="text-lg shrink-0" :style="{ color: item.accent }" />
          <span class="text-sm text-[var(--el-text-color-primary)]">{{
            $t(`dev.workflow.${item.titleKey}`)
          }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'
  import { PALETTE_ITEMS, type PaletteItemDef } from './automation-kinds'

  const { t } = useI18n()

  function onDragStart(e: DragEvent, item: PaletteItemDef) {
    if (!e.dataTransfer) return
    const payload = JSON.stringify({
      kind: item.kind,
      label: t(`dev.workflow.${item.titleKey}`)
    })
    // 与画布侧 dropEffect: copy 一致；move+copy 不匹配时部分浏览器会拒绝 drop
    e.dataTransfer.setData('application/json', payload)
    e.dataTransfer.setData('text/plain', payload)
    e.dataTransfer.effectAllowed = 'copy'
  }
</script>

<style scoped>
  .workflow-palette__item:focus-visible {
    outline: 2px solid var(--el-color-primary);
    outline-offset: 1px;
  }
</style>
