<!-- 流程基础信息：新建 / 编辑编号、备注、定时执行与 Cron -->
<template>
  <ElDialog
    v-model="visible"
    :title="title"
    width="min(520px, 92vw)"
    destroy-on-close
    append-to-body
    @closed="submitting = false"
  >
    <ElForm label-width="100px" class="max-w-xl">
      <ElFormItem :label="$t('dev.process.code')" required>
        <ElInput v-model="form.code" maxlength="64" show-word-limit :disabled="submitting" />
      </ElFormItem>
      <ElFormItem :label="$t('dev.process.remark')">
        <ElInput
          v-model="form.remark"
          type="textarea"
          :rows="3"
          maxlength="500"
          show-word-limit
          :disabled="submitting"
        />
      </ElFormItem>
      <ElFormItem :label="$t('dev.process.cronEnabled')">
        <ElSwitch
          v-model="form.cronEnabled"
          :disabled="submitting"
          :active-value="true"
          :inactive-value="false"
        />
      </ElFormItem>
      <ElFormItem :label="$t('dev.process.tags')">
        <div class="flex w-full max-w-full flex-col gap-2">
          <div v-if="orphanTags.length" class="flex flex-wrap gap-1">
            <ElTag v-for="tg in orphanTags" :key="tg.id" type="info" size="small">
              {{ $t('dev.process.tagOrphan', { id: tg.id }) }}
            </ElTag>
            <span class="text-xs text-[var(--el-text-color-secondary)]">{{
              $t('dev.process.tagOrphanHint')
            }}</span>
          </div>
          <ElSelect
            v-model="form.tagIds"
            multiple
            collapse-tags
            collapse-tags-tooltip
            filterable
            clearable
            class="w-full"
            :disabled="submitting"
            :placeholder="$t('dev.process.tagsPlaceholder')"
          >
            <ElOption
              v-for="opt in tagOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </ElSelect>
        </div>
      </ElFormItem>
      <ElFormItem
        v-if="form.cronEnabled"
        :label="$t('dev.process.cronExpr')"
        required
        :error="cronExprError"
      >
        <div class="flex max-w-full flex-col gap-2">
          <ElSelect
            :model-value="cronPresetValue"
            clearable
            filterable
            :placeholder="$t('dev.process.cronPresetPlaceholder')"
            :disabled="submitting"
            class="w-full"
            @update:model-value="onCronPresetPick"
          >
            <ElOption
              v-for="p in cronPresetOptions"
              :key="p.value"
              :label="p.label"
              :value="p.value"
            />
          </ElSelect>
          <ElInput
            v-model="form.cronExpr"
            :placeholder="$t('dev.process.cronPlaceholder')"
            :disabled="submitting"
            clearable
          />
          <p
            v-if="cronExprHint"
            class="m-0 text-xs leading-snug text-[var(--el-text-color-secondary)]"
          >
            {{ cronExprHint }}
          </p>
        </div>
      </ElFormItem>
    </ElForm>
    <template #footer>
      <ElButton @click="visible = false">{{ $t('common.cancel') }}</ElButton>
      <ElButton type="primary" :loading="submitting" @click="submit">{{
        $t('common.confirm')
      }}</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { computed, reactive, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { ElMessage } from 'element-plus'
  import {
    fetchDevProcessAdd,
    fetchDevProcessEdit,
    fetchDevProcessTagList
  } from '@/api/dev-process'
  import { describeCronExpr, isValidCronExpr } from '@/utils/cron-expr'

  /** 与后端 DefaultCronExprEveryMinute 一致：每分钟 */
  const CRON_DEFAULT_EVERY_MINUTE = '* * * * *'

  const props = defineProps<{
    mode: 'create' | 'edit'
    processId?: string
    /** 打开「编辑」时由列表传入 */
    initialCode?: string
    initialRemark?: string
    /** '1' | '0' */
    initialCronEnabled?: string
    initialCronExpr?: string
    /** 编辑时详情中的标签（name 为空为孤儿） */
    initialTags?: { id: string; name: string }[]
  }>()

  const emit = defineEmits<{
    saved: [payload: { id?: string }]
  }>()

  const visible = defineModel<boolean>({ required: true })

  const { t, locale } = useI18n()

  const CRON_PRESETS = [
    { value: '*/15 * * * *', i18nKey: 'cronPreset15m' as const },
    { value: '0 * * * *', i18nKey: 'cronPresetHourly' as const },
    { value: '0 4 * * *', i18nKey: 'cronPresetDaily' as const },
    { value: '0 9 * * 1', i18nKey: 'cronPresetWeeklyMon' as const }
  ]

  const cronLocale = computed(() => (String(locale.value).startsWith('zh') ? 'zh' : 'en'))

  const cronPresetOptions = computed(() =>
    CRON_PRESETS.map((p) => ({
      value: p.value,
      label: t(`dev.process.${p.i18nKey}`)
    }))
  )

  const cronPresetValue = computed(() => {
    if (!form.cronEnabled) return ''
    const e = form.cronExpr.trim()
    return CRON_PRESETS.some((p) => p.value === e) ? e : ''
  })

  function onCronPresetPick(v: string | null | undefined) {
    if (v) form.cronExpr = v
  }

  const cronExprError = computed(() => {
    if (!form.cronEnabled) return ''
    const raw = form.cronExpr.trim()
    if (!raw) return ''
    return isValidCronExpr(form.cronExpr) ? '' : t('dev.process.cronInvalid')
  })

  const cronExprHint = computed(() => {
    if (!form.cronEnabled) return ''
    const raw = form.cronExpr.trim()
    if (!raw || !isValidCronExpr(form.cronExpr)) return ''
    return describeCronExpr(form.cronExpr, cronLocale.value) ?? ''
  })

  const submitting = ref(false)
  const tagOptions = ref<{ label: string; value: string }[]>([])
  const form = reactive({
    code: '',
    remark: '',
    cronEnabled: false,
    cronExpr: '',
    tagIds: [] as string[]
  })

  const orphanTags = computed(() =>
    (props.initialTags ?? []).filter((t) => !t.name || !String(t.name).trim())
  )

  async function loadTagOptions() {
    try {
      const res = await fetchDevProcessTagList()
      tagOptions.value = (res.rows ?? []).map((r) => ({ label: r.name, value: r.id }))
    } catch {
      tagOptions.value = []
    }
  }

  const title = computed(() =>
    props.mode === 'create' ? t('dev.process.metaCreateTitle') : t('dev.process.metaEditTitle')
  )

  function syncFormFromProps() {
    if (props.mode === 'create') {
      form.code = ''
      form.remark = ''
      form.cronEnabled = false
      form.cronExpr = ''
      form.tagIds = []
      return
    }
    form.code = props.initialCode ?? ''
    form.remark = props.initialRemark ?? ''
    form.cronEnabled = props.initialCronEnabled === '1'
    form.cronExpr = props.initialCronExpr ?? ''
    form.tagIds = (props.initialTags ?? [])
      .filter((t) => t.name && String(t.name).trim())
      .map((t) => t.id)
  }

  watch(
    () => visible.value,
    (v) => {
      if (v) {
        syncFormFromProps()
        loadTagOptions()
      }
    }
  )

  watch(
    () => form.cronEnabled,
    (on) => {
      if (on && !form.cronExpr.trim()) {
        form.cronExpr = CRON_DEFAULT_EVERY_MINUTE
      }
    }
  )

  watch(
    () => [
      props.initialCode,
      props.initialRemark,
      props.initialCronEnabled,
      props.initialCronExpr,
      props.initialTags
    ],
    () => {
      if (visible.value && props.mode === 'edit') syncFormFromProps()
    }
  )

  async function submit() {
    const code = form.code.trim()
    if (!code) {
      ElMessage.warning(t('dev.process.codeRequired'))
      return
    }
    if (form.cronEnabled) {
      const ce = form.cronExpr.trim()
      if (!ce) {
        ElMessage.warning(t('dev.process.cronRequired'))
        return
      }
      if (!isValidCronExpr(ce)) {
        ElMessage.warning(t('dev.process.cronInvalid'))
        return
      }
    }
    submitting.value = true
    try {
      if (props.mode === 'create') {
        const res = await fetchDevProcessAdd({
          code,
          remark: form.remark,
          flow: '{"nodes":[],"edges":[]}',
          cronEnabled: form.cronEnabled ? 1 : 0,
          cronExpr: form.cronEnabled ? form.cronExpr.trim() : ''
        })
        ElMessage.success(t('dev.process.saveSuccess'))
        visible.value = false
        emit('saved', { id: res.id })
        return
      }
      const id = props.processId
      if (!id) {
        ElMessage.error(t('sys.common.operateFailed'))
        return
      }
      await fetchDevProcessEdit({
        id,
        code,
        remark: form.remark,
        cronEnabled: form.cronEnabled ? 1 : 0,
        cronExpr: form.cronEnabled ? form.cronExpr.trim() : '',
        tagIds: form.tagIds.map((id) => Number(id))
      })
      ElMessage.success(t('dev.process.saveSuccess'))
      visible.value = false
      emit('saved', {})
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : t('sys.common.operateFailed'))
    } finally {
      submitting.value = false
    }
  }
</script>
