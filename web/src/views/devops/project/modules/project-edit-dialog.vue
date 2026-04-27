<template>
  <ElDialog
    v-model="visible"
    :title="dialogType === 'add' ? $t('dev.project.addTitle') : $t('dev.project.editTitle')"
    width="520px"
    destroy-on-close
    append-to-body
    @closed="onClosed"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-width="100px">
      <ElFormItem :label="$t('dev.project.name')" prop="name">
        <ElInput v-model="form.name" maxlength="200" show-word-limit clearable />
      </ElFormItem>
      <ElFormItem :label="$t('dev.project.status')" prop="status">
        <ElSelect v-model="form.status" style="width: 100%">
          <ElOption :label="$t('dev.project.statusDraft')" :value="0" />
          <ElOption :label="$t('dev.project.statusActive')" :value="1" />
          <ElOption :label="$t('dev.project.statusPaused')" :value="2" />
          <ElOption :label="$t('dev.project.statusDone')" :value="3" />
        </ElSelect>
      </ElFormItem>
      <ElFormItem :label="$t('dev.project.progress')" prop="progress">
        <ElSlider v-model="form.progress" :min="0" :max="100" show-input />
      </ElFormItem>
      <ElFormItem :label="$t('dev.project.versionChangelog')" prop="versionChangelog">
        <ElInput
          v-model="form.versionChangelog"
          type="textarea"
          :rows="5"
          :placeholder="$t('dev.project.changelogPlaceholder')"
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
  import { useI18n } from 'vue-i18n'
  import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
  import {
    fetchDevProjectAdd,
    fetchDevProjectDetail,
    fetchDevProjectEdit,
    fetchDevProjectTagList
  } from '@/api/dev-project'

  const props = defineProps<{
    modelValue: boolean
    dialogType: 'add' | 'edit'
    projectId?: string
  }>()

  const emit = defineEmits<{
    (e: 'update:modelValue', v: boolean): void
    (e: 'success'): void
  }>()

  const { t } = useI18n()
  const visible = computed({
    get: () => props.modelValue,
    set: (v: boolean) => emit('update:modelValue', v)
  })

  const formRef = ref<FormInstance>()
  const submitting = ref(false)
  const tagOptions = ref<{ label: string; value: string }[]>([])

  const form = reactive({
    name: '',
    status: 0 as number,
    progress: 0,
    versionChangelog: '',
    tagIds: [] as string[]
  })

  const orphanTags = computed(() =>
    (detailTags.value ?? []).filter((t) => !t.name || !String(t.name).trim())
  )

  /** 详情中的标签（编辑时用于孤儿展示与字典多选） */
  const detailTags = ref<Api.DevProject.ProjectTagItem[] | undefined>(undefined)

  async function loadTagOptions() {
    try {
      const res = await fetchDevProjectTagList()
      tagOptions.value = (res.rows ?? []).map((r: { id: string; name: string }) => ({
        label: r.name,
        value: r.id
      }))
    } catch {
      tagOptions.value = []
    }
  }

  const rules: FormRules = {
    name: [{ required: true, message: () => t('dev.project.nameRequired'), trigger: 'blur' }]
  }

  watch(
    () => visible.value,
    async (v) => {
      if (!v) return
      if (props.dialogType === 'edit' && props.projectId) {
        try {
          const d = await fetchDevProjectDetail(props.projectId)
          form.name = d.name ?? ''
          form.status = Number(d.status ?? 0)
          form.progress = Number(d.progress ?? 0)
          form.versionChangelog = d.versionChangelog ?? ''
          detailTags.value = d.tags
          form.tagIds = (d.tags ?? [])
            .filter((tg) => tg.name && String(tg.name).trim())
            .map((tg) => tg.id)
        } catch {
          ElMessage.error(t('dev.project.loadFailed'))
          visible.value = false
        }
      } else {
        form.name = ''
        form.status = 0
        form.progress = 0
        form.versionChangelog = ''
        detailTags.value = undefined
        form.tagIds = []
      }
      loadTagOptions()
      nextTick(() => formRef.value?.clearValidate())
    }
  )

  function onClosed() {
    formRef.value?.resetFields()
  }

  async function submit() {
    await formRef.value?.validate().catch(() => Promise.reject())
    submitting.value = true
    try {
      if (props.dialogType === 'add') {
        await fetchDevProjectAdd({
          name: form.name,
          status: form.status,
          progress: form.progress,
          versionChangelog: form.versionChangelog,
          mindJson: '{}',
          tagIds: form.tagIds.map((id) => Number(id))
        })
      } else if (props.projectId) {
        await fetchDevProjectEdit({
          id: props.projectId,
          name: form.name,
          status: form.status,
          progress: form.progress,
          versionChangelog: form.versionChangelog,
          tagIds: form.tagIds.map((id) => Number(id))
        })
      }
      ElMessage.success(t('dev.project.saveSuccess'))
      visible.value = false
      emit('success')
    } finally {
      submitting.value = false
    }
  }
</script>
