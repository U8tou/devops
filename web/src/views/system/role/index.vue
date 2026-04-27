<!-- 角色管理页面 -->
<template>
  <div class="art-full-height">
    <RoleSearch
      v-show="showSearchBar"
      v-model="searchForm"
      @search="handleSearch"
      @reset="resetSearchParams"
    ></RoleSearch>

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
            <ElButton v-auth="'sys:role:add'" @click="showDialog('add')" v-ripple>{{
              $t('sys.role.addRole')
            }}</ElButton>
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
              <span>{{ row.createTime ? formatRoleTime(row.createTime) : '-' }}</span>
            </div>
            <div
              style="display: flex; gap: 4px; align-items: center; font-size: 12px; color: #909399"
            >
              <ArtSvgIcon
                icon="ri:refresh-line"
                style="flex-shrink: 0; font-size: 14px; color: #909399"
              />
              <span>{{ row.updateTime ? formatRoleTime(row.updateTime) : '-' }}</span>
            </div>
          </div>
        </template>
        <template #operation="{ row }">
          <div style="display: flex; align-items: center; justify-content: flex-end">
            <ArtButtonTable v-auth="'sys:role:edit'" type="edit" @click="showDialog('edit', row)" />
            <ArtButtonMore :list="moreButtonList" @click="(item) => buttonMoreClick(item, row)" />
          </div>
        </template>
      </ArtTable>
    </ElCard>

    <!-- 角色编辑弹窗 -->
    <RoleEditDialog
      v-model="dialogVisible"
      :dialog-type="dialogType"
      :role-data="currentRoleData"
      @success="refreshData"
    />

    <!-- 菜单权限弹窗 -->
    <RoleMenuDialog
      v-model="menuPermissionDialog"
      :role-data="currentRoleData"
      @success="refreshData"
    />

    <!-- 数据权限弹窗 -->
    <RoleDeptDialog
      v-model="deptPermissionDialog"
      :role-data="currentRoleData"
      @success="refreshData"
    />
  </div>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import { ButtonMoreItem } from '@/components/core/forms/art-button-more/index.vue'
  import { useTable } from '@/hooks/core/useTable'
  import { fetchGetRoleList, fetchDeleteRole } from '@/api/system-manage'
  import ArtButtonMore from '@/components/core/forms/art-button-more/index.vue'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import RoleSearch from './modules/role-search.vue'
  import RoleEditDialog from './modules/role-edit-dialog.vue'
  import RoleMenuDialog from './modules/role-menu-dialog.vue'
  import RoleDeptDialog from './modules/role-dept-dialog.vue'
  import { ElTag, ElMessageBox, ElMessage } from 'element-plus'
  import { parseTime } from '@/utils/ruoyi'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'

  defineOptions({ name: 'Role' })

  type RoleListItem = Api.SystemManage.RoleListItem

  // 搜索表单
  const searchForm = ref({
    name: undefined,
    role: undefined,
    status: undefined,
    daterange: undefined
  })

  const { t } = useI18n()
  const showSearchBar = ref(false)

  const dialogVisible = ref(false)
  const menuPermissionDialog = ref(false)
  const deptPermissionDialog = ref(false)
  const currentRoleData = ref<{ id: string } | undefined>(undefined)

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
      apiFn: fetchGetRoleList,
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
          label: t('sys.role.roleName'),
          minWidth: 120
        },
        {
          prop: 'role',
          label: t('sys.role.roleKey'),
          minWidth: 120
        },
        {
          prop: 'status',
          label: t('sys.user.status'),
          width: 100,
          formatter: (row: RoleListItem) => {
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
          width: 80
        },
        {
          prop: 'remark',
          label: t('sys.user.remark'),
          minWidth: 150,
          showOverflowTooltip: true
        },
        {
          prop: 'timeInfo',
          label: t('sys.user.timeInfo'),
          minWidth: 160,
          align: 'left',
          useSlot: true
        },
        {
          prop: 'operation',
          label: t('sys.user.action'),
          width: 120,
          fixed: 'right',
          align: 'right',
          headerAlign: 'center',
          useSlot: true
        }
      ]
    }
  })

  const dialogType = ref<'add' | 'edit'>('add')

  // 更多操作按钮列表
  const moreButtonList = computed(() => [
    {
      key: 'menu',
      label: t('sys.role.menuPerm'),
      icon: 'ri:shield-keyhole-line',
      auth: 'sys:role:assign_menu'
    },
    {
      key: 'dept',
      label: t('sys.role.dataPerm'),
      icon: 'ri:organization-chart',
      auth: 'sys:role:assign_dept'
    },
    {
      key: 'delete',
      label: t('sys.role.deleteRole'),
      icon: 'ri:delete-bin-4-line',
      color: '#f56c6c',
      auth: 'sys:role:del'
    }
  ])

  const showDialog = (type: 'add' | 'edit', row?: RoleListItem) => {
    dialogVisible.value = true
    dialogType.value = type
    currentRoleData.value = row ? { id: row.id } : undefined
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

  const buttonMoreClick = (item: ButtonMoreItem, row: RoleListItem) => {
    switch (item.key) {
      case 'menu':
        showMenuPermissionDialog(row)
        break
      case 'dept':
        showDeptPermissionDialog(row)
        break
      case 'delete':
        deleteRole(row)
        break
    }
  }

  const showMenuPermissionDialog = (row?: RoleListItem) => {
    menuPermissionDialog.value = true
    currentRoleData.value = row ? { id: row.id } : undefined
  }

  const showDeptPermissionDialog = (row?: RoleListItem) => {
    deptPermissionDialog.value = true
    currentRoleData.value = row ? { id: row.id } : undefined
  }

  /**
   * 格式化角色时间（统一使用 parseTime 函数）
   */
  const formatRoleTime = (time: string | number | null | undefined): string => {
    return parseTime(time) || '-'
  }

  const deleteRole = async (row: RoleListItem) => {
    try {
      await ElMessageBox.confirm(t('sys.role.deleteRoleConfirm'), t('common.tips'), {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      })
      await fetchDeleteRole([row.id])
      ElMessage.success(t('sys.common.deleteSuccess'))
      refreshData()
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error instanceof Error ? error.message : t('sys.common.deleteFailed'))
      }
    }
  }
</script>
