<!-- 流程标签字典 CRUD -->
<template>
  <ElDialog
    v-model="visible"
    :title="$t('dev.process.tagManageTitle')"
    width="min(480px, 92vw)"
    destroy-on-close
    append-to-body
    @open="load"
  >
    <div v-loading="loading" class="min-h-[120px]">
      <ElTable :data="rows" border size="small" class="w-full">
        <ElTableColumn :label="$t('dev.process.tagNameCol')" min-width="200">
          <template #default="{ row }">
            <ElInput
              v-if="editingId === row.id"
              v-model="editingName"
              size="small"
              maxlength="128"
              show-word-limit
              @keyup.enter="saveEdit(row)"
            />
            <span v-else>{{ row.name }}</span>
          </template>
        </ElTableColumn>
        <ElTableColumn :label="$t('sys.user.action')" width="160" align="center" fixed="right">
          <template #default="{ row }">
            <template v-if="editingId === row.id">
              <ElButton type="primary" link size="small" @click="saveEdit(row)">
                {{ $t('common.confirm') }}
              </ElButton>
              <ElButton link size="small" @click="cancelEdit">{{ $t('common.cancel') }}</ElButton>
            </template>
            <template v-else>
              <ElButton
                type="primary"
                link
                size="small"
                :disabled="!canEditTag"
                @click="startEdit(row)"
              >
                {{ $t('dev.process.edit') }}
              </ElButton>
              <ElButton type="danger" link size="small" :disabled="!canDelTag" @click="remove(row)">
                {{ $t('common.delete') }}
              </ElButton>
            </template>
          </template>
        </ElTableColumn>
      </ElTable>
      <div class="mt-3">
        <ElButton
          v-if="canAddTag"
          type="primary"
          plain
          size="small"
          :loading="adding"
          @click="addTag"
        >
          {{ $t('dev.process.tagAdd') }}
        </ElButton>
      </div>
    </div>
    <template #footer>
      <ElButton @click="visible = false">{{ $t('common.cancel') }}</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import {
    fetchDevProcessTagAdd,
    fetchDevProcessTagDel,
    fetchDevProcessTagEdit,
    fetchDevProcessTagList
  } from '@/api/dev-process'
  import { useUserStore } from '@/store/modules/user'

  defineOptions({ name: 'ProcessTagManageDialog' })

  const visible = defineModel<boolean>({ required: true })

  const { t } = useI18n()
  const userStore = useUserStore()

  const loading = ref(false)
  const adding = ref(false)
  const rows = ref<{ id: string; name: string }[]>([])
  const editingId = ref<string | null>(null)
  const editingName = ref('')

  const buttons = computed(() => userStore.info?.buttons ?? [])
  const canAddTag = computed(
    () =>
      buttons.value.includes('dev:process:edit') || buttons.value.includes('dev:process:tag:add')
  )
  const canEditTag = computed(
    () =>
      buttons.value.includes('dev:process:edit') || buttons.value.includes('dev:process:tag:edit')
  )
  const canDelTag = computed(
    () =>
      buttons.value.includes('dev:process:edit') || buttons.value.includes('dev:process:tag:del')
  )

  async function load() {
    loading.value = true
    editingId.value = null
    try {
      const res = await fetchDevProcessTagList()
      rows.value = res.rows ?? []
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : t('dev.process.tagLoadFailed'))
    } finally {
      loading.value = false
    }
  }

  function startEdit(row: { id: string; name: string }) {
    editingId.value = row.id
    editingName.value = row.name
  }

  function cancelEdit() {
    editingId.value = null
    editingName.value = ''
  }

  async function saveEdit(row: { id: string; name: string }) {
    const name = editingName.value.trim()
    if (!name) {
      ElMessage.warning(t('dev.process.tagNameRequired'))
      return
    }
    if (name === row.name && editingId.value === row.id) {
      cancelEdit()
      return
    }
    try {
      await fetchDevProcessTagEdit({ id: row.id, name })
      ElMessage.success(t('dev.process.tagSaveOk'))
      cancelEdit()
      await load()
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : t('sys.common.operateFailed'))
    }
  }

  async function addTag() {
    try {
      const { value } = await ElMessageBox.prompt(
        t('dev.process.tagNamePlaceholder'),
        t('dev.process.tagAdd'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          inputPattern: /\S+/,
          inputErrorMessage: t('dev.process.tagNameRequired')
        }
      )
      const name = String(value ?? '').trim()
      if (!name) return
      adding.value = true
      try {
        await fetchDevProcessTagAdd({ name })
        ElMessage.success(t('dev.process.tagSaveOk'))
        await load()
      } finally {
        adding.value = false
      }
    } catch {
      /* cancel */
    }
  }

  async function remove(row: { id: string; name: string }) {
    try {
      await ElMessageBox.confirm(
        t('dev.process.tagDeleteConfirm', { name: row.name }),
        t('common.tips'),
        { type: 'warning' }
      )
    } catch {
      return
    }
    try {
      await fetchDevProcessTagDel(row.id)
      ElMessage.success(t('dev.process.tagDeleteOk'))
      await load()
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : t('sys.common.operateFailed'))
    }
  }
</script>
