<!-- 部门管理页面 -->
<template>
  <div class="Dept-page art-full-height">
    <!-- 搜索栏 -->
    <ArtSearchBar
      v-model="formFilters"
      :items="formItems"
      :showExpand="false"
      @reset="handleReset"
      @search="handleSearch"
    />

    <ElCard class="art-table-card" shadow="never">
      <!-- 表格头部 -->
      <ArtTableHeader
        :showZebra="false"
        :loading="loading"
        v-model:columns="columnChecks"
        @refresh="handleRefresh"
      >
        <template #left>
          <ElButton v-auth="'sys:dept:add'" @click="handleAddDept" v-ripple>{{
            $t('sys.dept.addDept')
          }}</ElButton>
          <ElButton @click="toggleExpand" v-ripple>
            {{ isExpanded ? $t('table.searchBar.collapse') : $t('table.searchBar.expand') }}
          </ElButton>
        </template>
      </ArtTableHeader>

      <ArtTable
        ref="tableRef"
        rowKey="id"
        :loading="loading"
        :columns="columns"
        :data="filteredTableData"
        :stripe="false"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        :default-expand-all="false"
      />

      <!-- 部门弹窗 -->
      <DeptDialog v-model:visible="dialogVisible" :editData="editData" @submit="handleSubmit" />
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import { useTableColumns } from '@/hooks/core/useTableColumns'
  import DeptDialog from './modules/dept-dialog.vue'
  import type { DeptFormData } from './modules/dept-dialog.vue'
  import { fetchGetDeptAll, fetchDeleteDept } from '@/api/system-manage'
  import { ElTag, ElMessageBox, ElMessage } from 'element-plus'
  import { handleTree, parseTime } from '@/utils/ruoyi'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'

  defineOptions({ name: 'Depts' })

  // 状态管理
  const loading = ref(false)
  const isExpanded = ref(false)
  const tableRef = ref()

  // 弹窗相关
  const dialogVisible = ref(false)
  const editData = ref<{ id: string } | null>(null)

  // 搜索相关
  const initialSearchState = {
    name: '',
    route: ''
  }

  const formFilters = reactive({ ...initialSearchState })
  const appliedFilters = reactive({ ...initialSearchState })

  const { t } = useI18n()
  const formItems = computed(() => [
    {
      label: t('sys.dept.deptName'),
      key: 'name',
      type: 'input',
      props: { clearable: true }
    }
  ])

  onMounted(() => {
    getDeptList()
  })

  /**
   * 获取部门列表数据
   */
  const getDeptList = async (): Promise<void> => {
    loading.value = true

    try {
      const res = await fetchGetDeptAll()
      // 后端返回 { total, rows }，列表在 rows 中
      const rows = res?.rows ?? []
      tableData.value = handleTree(
        rows as unknown as Record<string, unknown>[],
        'id',
        'pid',
        'children'
      ) as unknown as Api.SystemManage.DeptTreeNode[]
    } catch (error) {
      throw error instanceof Error ? error : new Error(t('sys.common.operateFailed'))
    } finally {
      loading.value = false
    }
  }

  // 表格列配置
  const { columnChecks, columns } = useTableColumns(() => [
    {
      prop: 'name',
      label: t('sys.dept.deptName'),
      minWidth: 120,
      formatter: (row: Api.SystemManage.DeptTreeNode) => row.name ?? '-'
    },
    {
      prop: 'leader',
      label: t('sys.dept.leader'),
      formatter: (row: Api.SystemManage.DeptTreeNode) => row.leader ?? '-'
    },
    {
      prop: 'contact',
      label: t('sys.user.contact'),
      minWidth: 180,
      align: 'left',
      formatter: (row: Api.SystemManage.DeptTreeNode) => {
        return h('div', { style: 'display: flex; flex-direction: column; gap: 2px' }, [
          h('div', { style: 'display: flex; align-items: center; gap: 4px; color: #303133' }, [
            h(ArtSvgIcon, {
              icon: 'ri:phone-line',
              style: 'font-size: 13px; color: #909399; flex-shrink: 0'
            }),
            h('span', { style: 'font-size: 13px' }, row.phone || '-')
          ]),
          h(
            'div',
            {
              style: 'display: flex; align-items: center; gap: 4px; font-size: 12px; color: #909399'
            },
            [
              h(ArtSvgIcon, {
                icon: 'ri:mail-line',
                style: 'font-size: 13px; color: #909399; flex-shrink: 0'
              }),
              h('span', {}, row.email || '-')
            ]
          )
        ])
      }
    },
    {
      prop: 'status',
      label: t('sys.user.status'),
      width: 100,
      align: 'center',
      formatter: (row: Api.SystemManage.DeptTreeNode) => {
        const status = Number(row.status)
        const type = status === 1 ? 'success' : 'info'
        const text = status === 1 ? t('sys.common.statusNormal') : t('sys.common.statusDisabled')
        return h(ElTag, { type }, () => text)
      }
    },
    {
      prop: 'timeInfo',
      label: t('sys.user.timeInfo'),
      minWidth: 180,
      align: 'left',
      formatter: (row: Api.SystemManage.DeptTreeNode) => {
        return h('div', { style: 'display: flex; flex-direction: column; gap: 2px' }, [
          h('div', { style: 'display: flex; align-items: center; gap: 4px; color: #303133' }, [
            h(ArtSvgIcon, {
              icon: 'ri:calendar-line',
              style: 'font-size: 14px; color: #909399; flex-shrink: 0'
            }),
            h('span', {}, parseTime(row.createTime) || '-')
          ]),
          h(
            'div',
            {
              style: 'display: flex; align-items: center; gap: 4px; font-size: 12px; color: #909399'
            },
            [
              h(ArtSvgIcon, {
                icon: 'ri:refresh-line',
                style: 'font-size: 14px; color: #909399; flex-shrink: 0'
              }),
              h('span', {}, parseTime(row.updateTime) || '-')
            ]
          )
        ])
      }
    },
    {
      prop: 'operation',
      label: t('sys.user.action'),
      width: 180,
      align: 'right',
      headerAlign: 'center',
      formatter: (row: Api.SystemManage.DeptTreeNode) => {
        return h('div', { style: 'text-align: right' }, [
          h(ArtButtonTable, {
            type: 'edit',
            'v-auth': 'sys:dept:edit',
            onClick: () => handleEditDept(row)
          }),
          h(ArtButtonTable, {
            type: 'delete',
            'v-auth': 'sys:dept:del',
            onClick: () => handleDeleteDept(row)
          })
        ])
      }
    }
  ])

  // 数据相关（部门树，handleTree 返回带 children 的节点）
  const tableData = ref<Api.SystemManage.DeptTreeNode[]>([])

  /**
   * 重置搜索条件
   */
  const handleReset = (): void => {
    Object.assign(formFilters, { ...initialSearchState })
    Object.assign(appliedFilters, { ...initialSearchState })
    getDeptList()
  }

  /**
   * 执行搜索
   */
  const handleSearch = (): void => {
    Object.assign(appliedFilters, { ...formFilters })
    getDeptList()
  }

  /**
   * 刷新部门列表
   */
  const handleRefresh = (): void => {
    getDeptList()
  }

  /**
   * 深度克隆对象
   * @param obj 要克隆的对象
   * @returns 克隆后的对象
   */
  const deepClone = <T,>(obj: T): T => {
    if (obj === null || typeof obj !== 'object') return obj
    if (obj instanceof Date) return new Date(obj) as T
    if (Array.isArray(obj)) return obj.map((item) => deepClone(item)) as T

    const cloned = {} as T
    for (const key in obj) {
      if (Object.prototype.hasOwnProperty.call(obj, key)) {
        cloned[key] = deepClone(obj[key])
      }
    }
    return cloned
  }

  /**
   * 搜索部门（按部门名称、负责人过滤）
   * @param items 部门树节点数组
   * @returns 过滤后的树
   */
  const searchDept = (items: Api.SystemManage.DeptTreeNode[]): Api.SystemManage.DeptTreeNode[] => {
    const results: Api.SystemManage.DeptTreeNode[] = []
    const searchName = appliedFilters.name?.toLowerCase().trim() || ''
    const searchRoute = appliedFilters.route?.toLowerCase().trim() || ''

    for (const item of items) {
      const nameStr = (item.name ?? '').toString().toLowerCase()
      const routeStr = (item.leader ?? '').toString().toLowerCase()
      const nameMatch = !searchName || nameStr.includes(searchName)
      const routeMatch = !searchRoute || routeStr.includes(searchRoute)

      if (item.children?.length) {
        const matchedChildren = searchDept(item.children)
        if (matchedChildren.length > 0) {
          const clonedItem = deepClone(item)
          clonedItem.children = matchedChildren
          results.push(clonedItem)
          continue
        }
      }

      if (nameMatch && routeMatch) {
        results.push(deepClone(item))
      }
    }

    return results
  }

  // 过滤后的表格数据（部门列表不做 auth 子节点转换，直接使用搜索结果）
  const filteredTableData = computed(() => searchDept(tableData.value))

  /**
   * 新增部门
   */
  const handleAddDept = (): void => {
    editData.value = null
    dialogVisible.value = true
  }

  /**
   * 编辑部门
   * @param row 部门行数据
   */
  const handleEditDept = (row: Api.SystemManage.DeptTreeNode): void => {
    if (!row.id) {
      ElMessage.error(t('sys.dept.deptIdNotExist'))
      return
    }
    editData.value = { id: row.id }
    dialogVisible.value = true
  }

  /**
   * 提交部门表单
   * @param formData 部门表单数据
   */
  const handleSubmit = (formData: DeptFormData): void => {
    console.log('提交数据:', formData)
    // TODO: 调用新增/编辑部门 API
    getDeptList()
  }

  /**
   * 删除部门
   */
  const handleDeleteDept = async (
    row?: Api.SystemManage.DeptTreeNode | Record<string, unknown>
  ): Promise<void> => {
    if (!row) return

    try {
      const deptRow = row as Api.SystemManage.DeptTreeNode
      const name = deptRow.name || '该部门'
      const id = deptRow.id

      if (!id) {
        ElMessage.error(t('sys.dept.deptIdNotExist'))
        return
      }

      await ElMessageBox.confirm(t('sys.dept.deleteDeptConfirm', { name }), t('common.tips'), {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      })

      loading.value = true
      await fetchDeleteDept([String(id)])
      ElMessage.success(t('sys.common.deleteSuccess'))
      await getDeptList()
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error instanceof Error ? error.message : t('sys.common.deleteFailed'))
      }
    } finally {
      loading.value = false
    }
  }

  /**
   * 切换展开/收起所有部门
   */
  const toggleExpand = (): void => {
    isExpanded.value = !isExpanded.value
    nextTick(() => {
      if (tableRef.value?.elTableRef && filteredTableData.value) {
        const processRows = (rows: Api.SystemManage.DeptTreeNode[]) => {
          rows.forEach((row) => {
            if (row.children?.length) {
              tableRef.value.elTableRef.toggleRowExpansion(row, isExpanded.value)
              processRows(row.children)
            }
          })
        }
        processRows(filteredTableData.value)
      }
    })
  }
</script>
