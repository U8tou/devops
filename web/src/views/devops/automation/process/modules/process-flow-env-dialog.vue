<!-- 流程环境变量：弹窗编辑（与基础信息弹窗交互一致） -->
<template>
  <ElDialog
    v-model="visible"
    :title="$t('dev.process.flowEnv.title')"
    width="min(640px, 94vw)"
    destroy-on-close
    append-to-body
    @closed="onClosed"
  >
    <p class="mb-4 text-xs leading-relaxed text-[var(--el-text-color-secondary)]">
      {{ $t('dev.process.flowEnv.hint') }}
    </p>
    <div v-loading="loading" class="min-h-[120px]">
      <div class="flex flex-col gap-2">
        <div v-for="(row, idx) in rows" :key="idx" class="flex flex-wrap items-center gap-2">
          <ElInput
            v-model="row.key"
            class="!w-40 sm:!w-44"
            :placeholder="$t('dev.env.key')"
            :disabled="!canEdit || saving"
          />
          <ElInput
            v-model="row.value"
            class="min-w-0 flex-1"
            :placeholder="$t('dev.env.value')"
            :disabled="!canEdit || saving"
          />
          <ElButton v-if="canEdit" type="danger" plain :disabled="saving" @click="removeRow(idx)">
            {{ $t('dev.env.removeRow') }}
          </ElButton>
        </div>
      </div>
      <ElButton v-if="canEdit" class="mt-3" :disabled="saving" @click="addRow">
        {{ $t('dev.env.addRow') }}
      </ElButton>
    </div>
    <template #footer>
      <ElButton @click="visible = false">{{ $t('common.cancel') }}</ElButton>
      <ElButton type="primary" :loading="saving" :disabled="!canEdit" @click="submit">
        {{ $t('common.confirm') }}
      </ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { fetchDevProcessDetail, fetchDevProcessEditEnv } from '@/api/dev-process'
  import { useUserStore } from '@/store/modules/user'
  import { ElMessage } from 'element-plus'
  import { useI18n } from 'vue-i18n'

  defineOptions({ name: 'ProcessFlowEnvDialog' })

  const props = defineProps<{
    processId: string
  }>()

  const visible = defineModel<boolean>({ required: true })

  const { t } = useI18n()
  const userStore = useUserStore()

  const rows = ref<{ key: string; value: string }[]>([{ key: '', value: '' }])
  const loading = ref(false)
  const saving = ref(false)

  const canEdit = computed<boolean>(() =>
    (userStore.info?.buttons ?? []).includes('dev:process:edit')
  )

  function envToRows(env: Record<string, string>) {
    const entries = Object.entries(env || {})
    if (entries.length === 0) {
      return [{ key: '', value: '' }]
    }
    return entries.map(([key, value]) => ({ key, value }))
  }

  function rowsToEnv(): Record<string, string> {
    const out: Record<string, string> = {}
    for (const r of rows.value) {
      const k = (r.key || '').trim()
      if (k !== '') {
        out[k] = r.value ?? ''
      }
    }
    return out
  }

  async function load() {
    const id = props.processId?.trim()
    if (!id) return
    loading.value = true
    try {
      const d = await fetchDevProcessDetail(id)
      rows.value = envToRows(d.env || {})
    } catch {
      ElMessage.error(t('dev.process.flowEnv.loadFailed'))
    } finally {
      loading.value = false
    }
  }

  function addRow() {
    rows.value.push({ key: '', value: '' })
  }

  function removeRow(idx: number) {
    rows.value.splice(idx, 1)
    if (rows.value.length === 0) {
      rows.value.push({ key: '', value: '' })
    }
  }

  function onClosed() {
    saving.value = false
  }

  async function submit() {
    const id = props.processId?.trim()
    if (!id) return
    for (const r of rows.value) {
      const k = (r.key || '').trim()
      if (k === '' && (r.value || '').trim() !== '') {
        ElMessage.warning(t('dev.env.keyRequired'))
        return
      }
    }
    saving.value = true
    try {
      await fetchDevProcessEditEnv({ id, env: rowsToEnv() })
      ElMessage.success(t('dev.process.flowEnv.saveSuccess'))
      visible.value = false
    } finally {
      saving.value = false
    }
  }

  watch(visible, (v) => {
    if (v && props.processId?.trim()) {
      load()
    }
  })
</script>
