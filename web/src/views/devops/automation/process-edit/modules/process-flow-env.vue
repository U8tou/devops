<!-- 流程环境变量弹窗：与 process 模块共用实现，便于本目录内引用 -->
<script setup lang="ts">
  import ProcessFlowEnvDialog from '../../process/modules/process-flow-env-dialog.vue'

  const visible = defineModel<boolean>({ required: true })

  defineProps<{
    processId: string
  }>()
</script>

<template>
  <ProcessFlowEnvDialog v-model="visible" :process-id="processId" />
</template>
