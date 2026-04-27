<!-- 项目列表 dev_project -->
<template>
  <div class="art-full-height">
    <ArtSearchBar
      v-model="formFilters"
      :items="formItems"
      :show-expand="false"
      @reset="handleReset"
      @search="handleSearch"
    >
      <template #tagFilter="{ modelValue: m }">
        <div class="tag-filter-combo flex min-w-0 max-w-full items-stretch gap-1.5">
          <ElSelect
            class="tag-filter-mode shrink-0"
            style="width: 88px"
            :model-value="m.tagMatchMode ?? 'include'"
            @update:model-value="(v) => (m.tagMatchMode = v)"
          >
            <ElOption :label="$t('dev.process.tagMatchInclude')" value="include" />
            <ElOption :label="$t('dev.process.tagMatchExclude')" value="exclude" />
          </ElSelect>
          <ElSelect
            class="tag-filter-tags min-w-0 flex-1"
            multiple
            collapse-tags
            :max-collapse-tags="1"
            collapse-tags-tooltip
            clearable
            filterable
            :placeholder="$t('dev.process.tagsFilterPlaceholder')"
            :model-value="m.tagFilter ?? []"
            @update:model-value="(v) => (m.tagFilter = v)"
          >
            <ElOption
              v-for="opt in tagFilterOptions"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </ElSelect>
        </div>
      </template>
    </ArtSearchBar>

    <ElCard class="art-table-card" shadow="never" style="margin-top: 12px">
      <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
        <template #left>
          <ElButton v-auth="'dev:project:add'" type="primary" v-ripple @click="openDialog('add')">
            {{ $t('dev.project.createProject') }}
          </ElButton>
          <ElButton v-if="canOpenProjectTagManage" @click="tagManageVisible = true">
            {{ $t('dev.process.tagManage') }}
          </ElButton>
        </template>
      </ArtTableHeader>

      <ArtTable
        class="w-full"
        width="100%"
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
        <template #tags="{ row }">
          <div class="flex max-w-[220px] flex-wrap gap-1">
            <template v-if="row.tags?.length">
              <ElTag v-for="tg in row.tags" :key="tg.id" :type="tg.name ? undefined : 'info'">
                {{ tg.name || $t('dev.process.tagOrphan', { id: tg.id }) }}
              </ElTag>
            </template>
            <span v-else class="text-xs text-[var(--el-text-color-secondary)]">-</span>
          </div>
        </template>
        <template #status="{ row }">
          <ElTag :type="statusTagType(row.status)" size="small">{{
            statusLabel(row.status)
          }}</ElTag>
        </template>
        <template #progress="{ row }">
          <span class="text-sm">{{ formatProgress(row.progress) }}</span>
        </template>
        <template #versionChangelog="{ row }">
          <div class="flex justify-center">
            <ElButton
              v-if="(row.versionChangelog ?? '').trim()"
              type="primary"
              link
              size="small"
              @click="openChangelogDialog(row)"
            >
              {{ $t('dev.project.viewChangelog') }}
            </ElButton>
            <span v-else class="text-xs text-[var(--el-text-color-secondary)]">-</span>
          </div>
        </template>
        <template #timeInfo="{ row }">
          <div class="flex w-full min-w-0 flex-col gap-0.5">
            <div class="flex items-center gap-1 text-[var(--el-text-color-primary)]">
              <ArtSvgIcon
                icon="ri:calendar-line"
                class="shrink-0 text-sm text-[var(--el-text-color-secondary)]"
              />
              <span class="whitespace-nowrap text-sm">{{
                row.createTime ? formatTime(row.createTime) : '-'
              }}</span>
            </div>
            <div class="flex items-center gap-1 text-xs text-[var(--el-text-color-secondary)]">
              <ArtSvgIcon icon="ri:refresh-line" class="shrink-0 text-sm" />
              <span class="whitespace-nowrap">{{
                row.updateTime ? formatTime(row.updateTime) : '-'
              }}</span>
            </div>
          </div>
        </template>
        <template #operation="{ row }">
          <div class="flex shrink-0 flex-nowrap items-center justify-end gap-0.5">
            <ElTooltip :content="$t('dev.process.edit')" placement="top" :show-after="300">
              <span class="inline-flex items-center">
                <ArtButtonTable
                  v-auth="'dev:project:edit'"
                  icon="ri:pencil-line"
                  icon-class="bg-theme/12 text-theme"
                  @click="openDialog('edit', row)"
                />
              </span>
            </ElTooltip>
            <ElTooltip :content="$t('dev.project.requirements')" placement="top" :show-after="300">
              <span class="inline-flex items-center">
                <ArtButtonTable
                  v-auth="'dev:project:get'"
                  icon="ri:mind-map"
                  icon-class="bg-success/12 text-success"
                  @click="goMind(row.id)"
                />
              </span>
            </ElTooltip>
            <ElTooltip :content="$t('dev.project.delete')" placement="top" :show-after="300">
              <span class="inline-flex items-center">
                <ArtButtonTable
                  v-auth="'dev:project:del'"
                  icon="ri:delete-bin-line"
                  icon-class="bg-danger/12 text-danger"
                  @click="deleteRow(row)"
                />
              </span>
            </ElTooltip>
          </div>
        </template>
      </ArtTable>
    </ElCard>

    <ProjectEditDialog
      v-model="dialogVisible"
      :dialog-type="dialogType"
      :project-id="currentId"
      @success="refreshData"
    />

    <ProjectTagManageDialog v-model="tagManageVisible" />

    <ElDialog
      v-model="changelogDialogVisible"
      width="min(720px, 92vw)"
      destroy-on-close
      append-to-body
      :title="changelogDialogTitle"
    >
      <div
        class="max-h-[min(480px,60vh)] overflow-y-auto whitespace-pre-wrap break-words text-sm leading-relaxed text-[var(--el-text-color-primary)]"
      >
        {{ changelogDialogContent }}
      </div>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useRouter } from 'vue-router'
  import { useTable } from '@/hooks/core/useTable'
  import {
    fetchDevProjectDel,
    fetchDevProjectList,
    fetchDevProjectTagList
  } from '@/api/dev-project'
  import {
    ElDialog,
    ElMessage,
    ElMessageBox,
    ElOption,
    ElSelect,
    ElTag,
    ElTooltip
  } from 'element-plus'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'
  import { parseTime } from '@/utils/ruoyi'
  import ProjectEditDialog from './modules/project-edit-dialog.vue'
  import ProjectTagManageDialog from './modules/project-tag-manage-dialog.vue'
  import { useUserStore } from '@/store/modules/user'

  defineOptions({ name: 'DevopsProjectList' })

  type Row = Api.DevProject.ProjectListItem

  const { t } = useI18n()
  const router = useRouter()
  const userStore = useUserStore()

  const canOpenProjectTagManage = computed(() => {
    const b = userStore.info?.buttons ?? []
    return (
      b.includes('dev:project:edit') ||
      b.includes('dev:project:tag:list') ||
      b.includes('dev:project:tag:add') ||
      b.includes('dev:project:tag:edit') ||
      b.includes('dev:project:tag:del')
    )
  })

  /** 筛选：标签多选值含魔法项 __other__ 表示「其他」 */
  const TAG_FILTER_OTHER = '__other__'
  const tagManageVisible = ref(false)
  const tagDictRows = ref<{ id: string; name: string }[]>([])

  const formFilters = reactive<{
    name: string
    status: string
    tagMatchMode: 'include' | 'exclude'
    tagFilter: string[]
  }>({
    name: '',
    status: '',
    tagMatchMode: 'include',
    tagFilter: []
  })

  async function loadTagDictForFilter() {
    try {
      const res = await fetchDevProjectTagList()
      tagDictRows.value = res.rows ?? []
    } catch {
      tagDictRows.value = []
    }
  }

  const tagFilterOptions = computed(() => [
    ...tagDictRows.value.map((r) => ({ label: r.name, value: r.id })),
    { label: t('dev.process.tagOther'), value: TAG_FILTER_OTHER }
  ])

  const formItems = computed(() => [
    {
      label: t('dev.project.searchName'),
      key: 'name',
      type: 'input',
      span: 6,
      props: { clearable: true }
    },
    {
      label: t('dev.project.status'),
      key: 'status',
      type: 'select',
      span: 6,
      props: {
        clearable: true,
        options: [
          { label: t('dev.project.statusDraft'), value: '0' },
          { label: t('dev.project.statusActive'), value: '1' },
          { label: t('dev.project.statusPaused'), value: '2' },
          { label: t('dev.project.statusDone'), value: '3' }
        ]
      }
    },
    {
      label: t('dev.process.tagSearch'),
      key: 'tagFilter',
      type: 'select',
      span: 6,
      props: {
        multiple: true,
        collapseTags: true,
        collapseTagsTooltip: true,
        clearable: true,
        filterable: true,
        options: tagFilterOptions.value
      }
    }
  ])

  watch(tagManageVisible, (v) => {
    if (!v) {
      loadTagDictForFilter()
      refreshData()
    }
  })

  const dialogVisible = ref(false)
  const dialogType = ref<'add' | 'edit'>('add')
  const currentId = ref<string | undefined>(undefined)

  const changelogDialogVisible = ref(false)
  const changelogDialogTitle = ref('')
  const changelogDialogContent = ref('')

  function openChangelogDialog(row: Row) {
    changelogDialogTitle.value = `${row.name} — ${t('dev.project.versionChangelog')}`
    changelogDialogContent.value = (row.versionChangelog ?? '').trim()
    changelogDialogVisible.value = true
  }

  const {
    columns,
    columnChecks,
    data,
    loading,
    pagination,
    getData,
    searchParams,
    resetSearchParams,
    handleSizeChange,
    handleCurrentChange,
    refreshData
  } = useTable({
    core: {
      apiFn: fetchDevProjectList,
      apiParams: {
        current: 1,
        size: 20,
        name: '',
        status: '',
        tagIds: '',
        tagOther: '',
        tagMode: ''
      },
      columnsFactory: () => [
        {
          type: 'globalIndex',
          label: '#',
          width: 56,
          align: 'center'
        },
        {
          prop: 'name',
          label: t('dev.project.name'),
          minWidth: 200,
          showOverflowTooltip: true
        },
        {
          prop: 'tags',
          label: t('dev.process.tags'),
          minWidth: 160,
          useSlot: true
        },
        {
          prop: 'status',
          label: t('dev.project.status'),
          width: 104,
          align: 'center',
          headerAlign: 'center',
          useSlot: true
        },
        {
          prop: 'progress',
          label: t('dev.project.progress'),
          width: 88,
          align: 'center',
          headerAlign: 'center',
          useSlot: true
        },
        {
          prop: 'versionChangelog',
          label: t('dev.project.versionChangelog'),
          width: 100,
          align: 'center',
          headerAlign: 'center',
          useSlot: true
        },
        {
          prop: 'timeInfo',
          label: t('sys.user.timeInfo'),
          minWidth: 196,
          align: 'left',
          useSlot: true
        },
        {
          prop: 'operation',
          label: t('sys.user.action'),
          width: 168,
          minWidth: 168,
          fixed: 'right',
          align: 'right',
          headerAlign: 'center',
          useSlot: true
        }
      ]
    }
  })

  const formatTime = (time: string | number | null | undefined) => parseTime(time) || '-'

  function formatProgress(p: string | undefined) {
    const n = Number(p)
    if (!Number.isFinite(n)) return '-'
    return `${Math.round(n)}%`
  }

  function statusLabel(s: string | undefined) {
    switch (String(s ?? '')) {
      case '0':
        return t('dev.project.statusDraft')
      case '1':
        return t('dev.project.statusActive')
      case '2':
        return t('dev.project.statusPaused')
      case '3':
        return t('dev.project.statusDone')
      default:
        return '-'
    }
  }

  function statusTagType(s: string | undefined): 'info' | 'success' | 'warning' | 'danger' {
    switch (String(s ?? '')) {
      case '1':
        return 'success'
      case '2':
        return 'warning'
      case '3':
        return 'info'
      default:
        return 'info'
    }
  }

  function buildTagSearchParams() {
    const ids: string[] = []
    let other = false
    for (const v of formFilters.tagFilter ?? []) {
      if (v === TAG_FILTER_OTHER) other = true
      else ids.push(v)
    }
    return { tagIds: ids.join(','), tagOther: other ? '1' : '' }
  }

  const handleSearch = () => {
    const tagQ = buildTagSearchParams()
    Object.assign(searchParams, {
      name: formFilters.name?.trim() ?? '',
      status: formFilters.status ?? '',
      tagIds: tagQ.tagIds,
      tagOther: tagQ.tagOther,
      tagMode: formFilters.tagMatchMode === 'exclude' ? 'exclude' : 'include'
    })
    getData()
  }

  const handleReset = () => {
    formFilters.name = ''
    formFilters.status = ''
    formFilters.tagMatchMode = 'include'
    formFilters.tagFilter = []
    resetSearchParams()
    Object.assign(searchParams, {
      name: '',
      status: '',
      tagIds: '',
      tagOther: '',
      tagMode: ''
    })
    getData()
  }

  onMounted(() => {
    loadTagDictForFilter()
  })

  function openDialog(type: 'add' | 'edit', row?: Row) {
    dialogType.value = type
    currentId.value = row?.id
    dialogVisible.value = true
  }

  function goMind(id: string) {
    router.push({ name: 'DevopsProjectMind', params: { id } })
  }

  async function deleteRow(row: Row) {
    try {
      await ElMessageBox.confirm(
        t('dev.project.deleteConfirm', { name: row.name }),
        t('common.tips'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
      )
      await fetchDevProjectDel([row.id])
      ElMessage.success(t('sys.common.deleteSuccess'))
      refreshData()
    } catch (e) {
      if (e !== 'cancel') {
        ElMessage.error(e instanceof Error ? e.message : t('sys.common.deleteFailed'))
      }
    }
  }
</script>
