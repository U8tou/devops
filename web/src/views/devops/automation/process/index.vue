<!-- 自动化流程列表 -->
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
          <ElButton v-auth="'dev:process:add'" type="primary" v-ripple @click="openCreateMeta">
            {{ $t('dev.process.add') }}
          </ElButton>
          <ElButton v-if="canOpenProcessTagManage" @click="tagManageVisible = true">
            {{ $t('dev.process.tagManage') }}
          </ElButton>
        </template>
      </ArtTableHeader>

      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
        <template #timeInfo="{ row }">
          <div class="flex min-w-[188px] flex-col gap-0.5">
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
        <template #cronEnabled="{ row }">
          <ElTooltip :content="$t('dev.process.cronEnabledHint')" placement="top" :show-after="300">
            <span class="inline-flex items-center justify-center">
              <ElSwitch
                :model-value="row.cronEnabled === '1'"
                :disabled="!canEditProcess || cronSwitchingId === row.id"
                @change="(val: string | number | boolean) => onCronEnabledChange(row, val === true)"
              />
            </span>
          </ElTooltip>
        </template>
        <template #cronExpr="{ row }">
          <span class="text-sm">{{ row.cronExpr ? row.cronExpr : '-' }}</span>
        </template>
        <template #lastExecTime="{ row }">
          <span class="block whitespace-nowrap text-sm">{{
            row.lastExecTime && row.lastExecTime !== '0' ? formatTime(row.lastExecTime) : '-'
          }}</span>
        </template>
        <template #lastExecDurationMs="{ row }">
          <span class="text-sm">{{ formatExecDurationMs(row.lastExecDurationMs) }}</span>
        </template>
        <template #lastExecResult="{ row }">
          <ElTag v-if="execStatus(row) === 'success'" type="success" size="small">{{
            $t('dev.process.execStatusSuccess')
          }}</ElTag>
          <ElTag v-else-if="execStatus(row) === 'failed'" type="danger" size="small">{{
            $t('dev.process.execStatusFailed')
          }}</ElTag>
          <ElTag v-else-if="execStatus(row) === 'cancelled'" type="warning" size="small">{{
            $t('dev.process.execStatusCancelled')
          }}</ElTag>
          <span v-else class="text-sm text-[var(--el-text-color-secondary)]">-</span>
        </template>
        <template #lastExecLog="{ row }">
          <ElButton
            v-if="hasExecLogEntry(row)"
            type="primary"
            link
            size="small"
            @click="openExecLog(row)"
          >
            {{ $t('dev.process.viewExecLog') }}
          </ElButton>
          <span v-else class="text-xs text-[var(--el-text-color-secondary)]">-</span>
        </template>
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
        <template #operation="{ row }">
          <div class="flex shrink-0 flex-nowrap items-center justify-end gap-0.5">
            <ElTooltip :content="$t('dev.process.run')" placement="top" :show-after="300">
              <span class="inline-flex items-center">
                <ArtButtonTable
                  v-auth="'dev:process:get'"
                  icon="ri:play-fill"
                  icon-class="bg-success/12 text-success"
                  @click="openRun(row)"
                />
              </span>
            </ElTooltip>
            <ElTooltip :content="$t('dev.process.editFlow')" placement="top" :show-after="300">
              <span class="inline-flex items-center">
                <ArtButtonTable
                  v-auth="'dev:process:edit'"
                  icon="ri:flow-chart"
                  icon-class="bg-theme/12 text-theme"
                  @click="goFlowEdit(row.id)"
                />
              </span>
            </ElTooltip>
            <ElTooltip :content="$t('dev.process.delete')" placement="top" :show-after="300">
              <span class="inline-flex items-center">
                <ArtButtonTable
                  v-auth="'dev:process:del'"
                  icon="ri:delete-bin-line"
                  icon-class="bg-danger/12 text-danger"
                  @click="deleteProcess(row)"
                />
              </span>
            </ElTooltip>
          </div>
        </template>
      </ArtTable>
    </ElCard>

    <ProcessRunDialog
      v-model="runVisible"
      :flow-json="runFlowJson"
      :process-code="runCode"
      :process-id="runProcessId"
      @closed="refreshData"
    />

    <ProcessMetaDialog v-model="metaVisible" mode="create" @saved="onMetaSaved" />

    <ProcessTagManageDialog v-model="tagManageVisible" />

    <ElDialog v-model="logDialogVisible" width="min(720px, 92vw)" destroy-on-close append-to-body>
      <template #header>
        <div class="flex flex-col items-start gap-1 pr-6">
          <span class="text-base font-medium text-[var(--el-text-color-primary)]">{{
            $t('dev.process.execLogDialogTitle')
          }}</span>
          <span
            v-if="logDialogDurationLabel"
            class="text-sm font-normal text-[var(--el-text-color-secondary)]"
          >
            {{ $t('dev.process.lastExecDuration') }}：{{ logDialogDurationLabel }}
          </span>
        </div>
      </template>
      <div
        v-if="logDialogContent.trim()"
        class="m-0 max-h-[min(60vh,520px)] overflow-auto rounded border border-[var(--el-border-color-lighter)]"
      >
        <div
          v-for="(block, idx) in execLogBlocks"
          :key="idx"
          class="exec-log-block whitespace-pre-wrap break-all px-3 py-2 font-mono text-xs leading-relaxed"
          :class="idx % 2 === 0 ? 'bg-[var(--el-fill-color-lighter)]' : 'bg-[var(--el-fill-color)]'"
        >
          {{ block.join('\n') }}
        </div>
      </div>
      <p v-else class="text-sm text-[var(--el-text-color-secondary)]">
        {{ $t('dev.process.execLogEmpty') }}
      </p>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useRouter } from 'vue-router'
  import { useTable } from '@/hooks/core/useTable'
  import {
    fetchDevProcessList,
    fetchDevProcessDetail,
    fetchDevProcessDel,
    fetchDevProcessSetCronEnabled,
    fetchDevProcessTagList
  } from '@/api/dev-process'
  import { useUserStore } from '@/store/modules/user'
  import {
    ElButton,
    ElDialog,
    ElMessage,
    ElMessageBox,
    ElOption,
    ElSelect,
    ElSwitch,
    ElTag,
    ElTooltip
  } from 'element-plus'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'
  import { parseTime } from '@/utils/ruoyi'
  import ProcessRunDialog from './modules/process-run-dialog.vue'
  import ProcessMetaDialog from './modules/process-meta-dialog.vue'
  import ProcessTagManageDialog from './modules/process-tag-manage-dialog.vue'
  import { splitExecLogIntoBlocks } from './modules/exec-log-blocks'

  defineOptions({ name: 'DevopsProcessList' })

  type Row = Api.DevProcess.ProcessListItem

  const { t } = useI18n()
  const router = useRouter()
  const userStore = useUserStore()

  const canEditProcess = computed(() =>
    (userStore.info?.buttons ?? []).includes('dev:process:edit')
  )

  const canOpenProcessTagManage = computed(() => {
    const b = userStore.info?.buttons ?? []
    return (
      b.includes('dev:process:edit') ||
      b.includes('dev:process:tag:list') ||
      b.includes('dev:process:tag:add') ||
      b.includes('dev:process:tag:edit') ||
      b.includes('dev:process:tag:del')
    )
  })

  const cronSwitchingId = ref<string | null>(null)

  /** 筛选：标签多选值含魔法项 __other__ 表示「其他」 */
  const TAG_FILTER_OTHER = '__other__'
  const tagManageVisible = ref(false)
  const tagDictRows = ref<{ id: string; name: string }[]>([])

  const formFilters = reactive<{
    code: string
    tagMatchMode: 'include' | 'exclude'
    tagFilter: string[]
  }>({ code: '', tagMatchMode: 'include', tagFilter: [] })

  async function loadTagDictForFilter() {
    try {
      const res = await fetchDevProcessTagList()
      tagDictRows.value = res.rows ?? []
    } catch {
      tagDictRows.value = []
    }
  }

  /** 标签筛选：多选项（含「其他」），供 #tagFilter 插槽使用 */
  const tagFilterOptions = computed(() => [
    ...tagDictRows.value.map((r) => ({ label: r.name, value: r.id })),
    { label: t('dev.process.tagOther'), value: TAG_FILTER_OTHER }
  ])

  const formItems = computed(() => [
    {
      label: t('dev.process.searchCode'),
      key: 'code',
      type: 'input',
      span: 6,
      props: { clearable: true }
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

  const runVisible = ref(false)
  const runFlowJson = ref('')
  const runCode = ref('')
  const runProcessId = ref('')

  const metaVisible = ref(false)

  watch(tagManageVisible, (v) => {
    if (!v) {
      loadTagDictForFilter()
      refreshData()
    }
  })

  const logDialogVisible = ref(false)
  const logDialogContent = ref('')
  /** 详情接口返回的最近一次总耗时（毫秒，字符串），用于日志弹窗标题旁展示 */
  const logDialogDurationMs = ref('')

  /** 去掉前导空白，避免 pre 首行不顶格（接口文本或模板在标签后换行会进预格式化内容） */
  const logDialogDisplay = computed(() => (logDialogContent.value ?? '').trimStart())

  /** 按节点分段，与列表表格栅栏色交替 */
  const execLogBlocks = computed(() => splitExecLogIntoBlocks(logDialogDisplay.value))

  /** 列表 / 日志弹窗：展示最近一次执行总耗时（毫秒） */
  function formatExecDurationMs(ms: string | number | null | undefined): string {
    const n = typeof ms === 'number' ? ms : Number(ms)
    if (!Number.isFinite(n) || n <= 0) return '-'
    if (n < 1000) return `${Math.round(n)}ms`
    return `${(n / 1000).toFixed(2)}s`
  }

  const logDialogDurationLabel = computed(() => formatExecDurationMs(logDialogDurationMs.value))

  function hasExecLogEntry(row: Row): boolean {
    const t = row.lastExecTime
    return t != null && t !== '' && t !== '0'
  }

  /** 与后端列表 normalize 一致，并 trim 避免空白导致不匹配 */
  function execStatus(row: Row): 'success' | 'failed' | 'cancelled' | '' {
    const s = (row.lastExecResult ?? '').trim().toLowerCase()
    if (s === 'success' || s === 'failed' || s === 'cancelled') return s
    return ''
  }

  async function openExecLog(row: Row) {
    logDialogContent.value = ''
    logDialogDurationMs.value = ''
    logDialogVisible.value = true
    try {
      const res = await fetchDevProcessDetail(row.id)
      logDialogContent.value = res.lastExecLog ?? ''
      logDialogDurationMs.value = res.lastExecDurationMs ?? ''
    } catch (e) {
      logDialogVisible.value = false
      ElMessage.error(e instanceof Error ? e.message : t('dev.process.loadFailed'))
    }
  }

  const formatTime = (time: string | number | null | undefined) => parseTime(time) || '-'

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
      apiFn: fetchDevProcessList,
      apiParams: {
        current: 1,
        size: 20,
        code: '',
        tagIds: '',
        tagOther: '',
        tagMode: ''
      },
      columnsFactory: () => [
        {
          type: 'globalIndex',
          label: '#',
          width: 64,
          align: 'center'
        },
        {
          prop: 'code',
          label: t('dev.process.code'),
          width: 120,
          align: 'center'
        },
        {
          prop: 'cronEnabled',
          label: t('dev.process.cronEnabled'),
          width: 120,
          align: 'center',
          useSlot: true
        },
        {
          prop: 'cronExpr',
          label: t('dev.process.cronExpr'),
          minWidth: 120,
          showOverflowTooltip: true,
          useSlot: true
        },
        {
          prop: 'lastExecTime',
          label: t('dev.process.lastExecTime'),
          minWidth: 178,
          useSlot: true
        },
        {
          prop: 'lastExecDurationMs',
          label: t('dev.process.lastExecDuration'),
          width: 110,
          align: 'center',
          useSlot: true
        },
        {
          prop: 'lastExecResult',
          label: t('dev.process.lastExecResult'),
          width: 100,
          align: 'center',
          useSlot: true
        },
        {
          prop: 'lastExecLog',
          label: t('dev.process.lastExecLog'),
          width: 100,
          align: 'center',
          useSlot: true
        },
        {
          prop: 'timeInfo',
          label: t('sys.user.timeInfo'),
          minWidth: 200,
          useSlot: true
        },
        {
          prop: 'tags',
          label: t('dev.process.tags'),
          minWidth: 160,
          useSlot: true
        },
        {
          prop: 'remark',
          label: t('dev.process.remark'),
          minWidth: 140,
          showOverflowTooltip: true
        },
        {
          prop: 'operation',
          label: t('sys.user.action'),
          minWidth: 180,
          width: 180,
          fixed: 'right',
          align: 'right',
          headerAlign: 'center',
          useSlot: true
        }
      ]
    }
  })

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
      code: formFilters.code ?? '',
      tagIds: tagQ.tagIds,
      tagOther: tagQ.tagOther,
      tagMode: formFilters.tagMatchMode === 'exclude' ? 'exclude' : 'include'
    })
    getData()
  }

  const handleReset = () => {
    formFilters.code = ''
    formFilters.tagMatchMode = 'include'
    formFilters.tagFilter = []
    resetSearchParams()
    Object.assign(searchParams, { code: '', tagIds: '', tagOther: '', tagMode: '' })
    getData()
  }

  onMounted(() => {
    loadTagDictForFilter()
  })

  function openCreateMeta() {
    metaVisible.value = true
  }

  function onMetaSaved(payload: { id?: string }) {
    if (payload.id) {
      router.push({ name: 'DevopsProcessEdit', params: { id: payload.id } })
    }
  }

  const goFlowEdit = (id: string) => {
    router.push({ name: 'DevopsProcessEdit', params: { id } })
  }

  async function onCronEnabledChange(row: Row, enabled: boolean): Promise<void> {
    if (!canEditProcess.value) return
    if (enabled) {
      try {
        await ElMessageBox.confirm(t('dev.process.cronEnableConfirm'), t('common.tips'), {
          type: 'warning',
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel')
        })
      } catch {
        return
      }
    }
    cronSwitchingId.value = row.id
    try {
      await fetchDevProcessSetCronEnabled({ id: row.id, enabled })
      ElMessage.success(t('dev.process.cronEnabledUpdateOk'))
      refreshData()
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : t('sys.common.operateFailed'))
    } finally {
      cronSwitchingId.value = null
    }
  }

  const openRun = async (row: Row) => {
    try {
      await ElMessageBox.confirm(t('dev.process.runConfirm'), t('common.tips'), {
        type: 'warning',
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel')
      })
    } catch {
      return
    }
    try {
      const res = await fetchDevProcessDetail(row.id)
      runProcessId.value = row.id
      runFlowJson.value = res.flow || ''
      runCode.value = res.code
      runVisible.value = true
    } catch (e) {
      ElMessage.error(e instanceof Error ? e.message : t('dev.process.loadFailed'))
    }
  }

  async function deleteProcess(row: Row) {
    try {
      await ElMessageBox.confirm(
        t('dev.process.deleteConfirm', { code: row.code || row.id }),
        t('common.tips'),
        {
          type: 'warning',
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel')
        }
      )
      await fetchDevProcessDel([row.id])
      ElMessage.success(t('sys.common.deleteSuccess'))
      refreshData()
    } catch (e) {
      if (e !== 'cancel') {
        ElMessage.error(e instanceof Error ? e.message : t('sys.common.deleteFailed'))
      }
    }
  }
</script>
