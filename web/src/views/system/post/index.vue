<!-- 岗位管理页面 -->
<template>
  <div class="art-full-height">
    <!-- 搜索栏 -->
    <PostSearch
      v-show="showSearchBar"
      v-model="searchForm"
      @search="handleSearch"
      @reset="resetSearchParams"
    ></PostSearch>

    <ElCard
      class="art-table-card"
      shadow="never"
      :style="{ 'margin-top': showSearchBar ? '12px' : '0' }"
    >
      <ArtTableHeader
        v-model:columns="columnChecks"
        v-model:showSearchBar="showSearchBar"
        :loading="loading"
        @refresh="refreshData"
      >
        <template #left>
          <ElSpace wrap>
            <ElButton v-auth="'sys:post:add'" @click="showDialog('add')" v-ripple>
              {{ $t('sys.post.addPost') }}
            </ElButton>
          </ElSpace>
        </template>
      </ArtTableHeader>

      <!-- 表格 -->
      <ArtTable
        :loading="loading"
        :data="data"
        :columns="columns"
        :pagination="pagination"
        @pagination:size-change="handleSizeChange"
        @pagination:current-change="handleCurrentChange"
      >
        <template #timeInfo="{ row }">
          <div style="display: flex; flex-direction: column; gap: 2px">
            <div style="display: flex; gap: 4px; align-items: center; color: #303133">
              <ArtSvgIcon
                icon="ri:calendar-line"
                style="flex-shrink: 0; font-size: 14px; color: #909399"
              />
              <span>{{ row.createTime ? formatPostTime(row.createTime) : '-' }}</span>
            </div>
            <div
              style="display: flex; gap: 4px; align-items: center; font-size: 12px; color: #909399"
            >
              <ArtSvgIcon
                icon="ri:refresh-line"
                style="flex-shrink: 0; font-size: 14px; color: #909399"
              />
              <span>{{ row.updateTime ? formatPostTime(row.updateTime) : '-' }}</span>
            </div>
          </div>
        </template>
        <template #operation="{ row }">
          <div style="display: flex; align-items: center; justify-content: flex-end">
            <ArtButtonTable v-auth="'sys:post:edit'" type="edit" @click="showDialog('edit', row)" />
            <ArtButtonTable v-auth="'sys:post:del'" type="delete" @click="deletePost(row)" />
          </div>
        </template>
      </ArtTable>
    </ElCard>

    <!-- 岗位编辑弹窗 -->
    <PostEditDialog
      v-model="dialogVisible"
      :dialog-type="dialogType"
      :post-data="currentPostData"
      @success="refreshData"
    />
  </div>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import { useTable } from '@/hooks/core/useTable'
  import { fetchGetPostList, fetchDeletePost } from '@/api/system-manage'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import PostSearch from './modules/post-search.vue'
  import PostEditDialog from './modules/post-edit-dialog.vue'
  import { ElTag, ElMessageBox, ElMessage } from 'element-plus'
  import { parseTime } from '@/utils/ruoyi'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'
  import { h } from 'vue'

  defineOptions({ name: 'Post' })

  type PostListItem = Api.SystemManage.PostListItem

  // 搜索表单
  const searchForm = ref({
    name: undefined,
    status: undefined,
    daterange: undefined
  })

  const { t } = useI18n()
  const showSearchBar = ref(true)

  const dialogVisible = ref(false)
  const currentPostData = ref<{ id: string } | undefined>(undefined)

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
    // 核心配置
    core: {
      apiFn: fetchGetPostList,
      apiParams: {
        current: 1,
        size: 20
      },
      // 排除 apiParams 中的属性
      excludeParams: ['daterange'],
      columnsFactory: () => [
        {
          type: 'globalIndex',
          label: t('sys.user.index'),
          width: 80,
          align: 'center'
        },
        {
          prop: 'name',
          label: t('sys.post.postName'),
          minWidth: 180
        },
        {
          prop: 'status',
          label: t('sys.user.status'),
          width: 120,
          align: 'center',
          formatter: (row: PostListItem) => {
            const statusConfig =
              row.status === 1
                ? { type: 'success', text: t('sys.common.statusNormal') }
                : { type: 'info', text: t('sys.common.statusDisabled') }
            return h(
              ElTag,
              { type: statusConfig.type as 'success' | 'info' },
              () => statusConfig.text
            )
          }
        },
        {
          prop: 'sort',
          label: t('sys.dept.sort'),
          width: 100,
          align: 'center'
        },
        {
          prop: 'timeInfo',
          label: t('sys.user.timeInfo'),
          minWidth: 220,
          align: 'left',
          useSlot: true
        },
        {
          prop: 'operation',
          label: t('sys.user.action'),
          width: 150,
          fixed: 'right',
          align: 'right',
          headerAlign: 'center',
          useSlot: true
        }
      ]
    }
  })

  const dialogType = ref<'add' | 'edit'>('add')

  const showDialog = (type: 'add' | 'edit', row?: PostListItem) => {
    dialogVisible.value = true
    dialogType.value = type
    currentPostData.value = row ? { id: row.id } : undefined
  }

  /**
   * 搜索处理
   * @param params 搜索参数
   */
  const handleSearch = (params: Record<string, any>) => {
    // 处理日期区间参数，把 daterange 转换为 startTime 和 endTime
    const { daterange, ...filtersParams } = params
    const [startTime, endTime] = Array.isArray(daterange) ? daterange : [null, null]

    // 搜索参数赋值
    Object.assign(searchParams, { ...filtersParams, startTime, endTime })
    getData()
  }

  /**
   * 格式化岗位时间（统一使用 parseTime 函数）
   */
  const formatPostTime = (time: string | number | null | undefined): string => {
    return parseTime(time) || '-'
  }

  const deletePost = async (row: PostListItem) => {
    try {
      await ElMessageBox.confirm(t('sys.common.deleteConfirm'), t('common.tips'), {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      })
      await fetchDeletePost([row.id])
      ElMessage.success(t('sys.common.deleteSuccess'))
      refreshData()
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error instanceof Error ? error.message : t('sys.common.deleteFailed'))
      }
    }
  }
</script>
