<!-- 自动化：占位页，与控制台、流程列表同级；后续可在此挂载子模块 -->
<template>
  <div class="art-full-height flex flex-col items-center justify-center gap-2 p-6">
    <ArtSvgIcon icon="ri:robot-2-line" class="text-4xl text-[var(--el-text-color-secondary)]" />
    <p class="text-sm text-[var(--el-text-color-secondary)]">
      {{ $t('menus.devops.automationPlaceholder') }}
    </p>
  </div>
</template>

<script setup lang="ts">
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'

  defineOptions({ name: 'DevopsAutomationHome' })
</script>
