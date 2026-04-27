<template>
  <ElDrawer v-model="visible" title="操作权限" size="520px" direction="rtl" @close="handleClose">
    <div class="role-menu-drawer-body">
      <ElScrollbar class="role-menu-drawer-scroll">
        <ElEmpty
          v-if="backendMenus.length === 0"
          :description="loading ? '加载中...' : '暂无菜单数据'"
        />
        <ElTree
          v-else
          ref="treeRef"
          class="role-menu-tree"
          :data="backendMenus"
          show-checkbox
          node-key="id"
          :indent="24"
          :default-expand-all="isExpandAll"
          :default-checked-keys="roleMenuIds"
          :props="defaultProps"
          :check-strictly="menuLinkage === 2"
          @check="handleTreeCheck"
        >
          <template #default="{ data }">
            <!-- 普通菜单节点 -->
            <span v-if="data.types !== 2">{{ menuNodeLabel(data) }}</span>
            <!-- 操作按钮节点：单独包一层，方便控制 hover 背景范围 -->
            <span v-else class="role-menu-permission-label">
              {{ menuNodeLabel(data) }}
            </span>
          </template>
        </ElTree>
      </ElScrollbar>

      <!-- 父子联动配置 -->
      <div class="role-menu-linkage-row">
        <span class="linkage-label">父子节点联动：</span>
        <ElRadioGroup v-model="menuLinkage" size="small">
          <ElRadio :label="1">联动</ElRadio>
          <ElRadio :label="2">不联动</ElRadio>
        </ElRadioGroup>
      </div>
    </div>

    <template #footer>
      <ElButton @click="toggleExpandAll">{{ isExpandAll ? '全部收起' : '全部展开' }}</ElButton>
      <ElButton @click="toggleSelectAll" style="margin-left: 8px">{{
        isSelectAll ? '取消全选' : '全部选择'
      }}</ElButton>
      <ElButton type="primary" @click="savePermission">保存</ElButton>
    </template>
  </ElDrawer>
</template>

<script setup lang="ts">
  import { handleTree } from '@/utils/ruoyi'
  import { fetchGetRoleInfo, fetchGetMenuAll, fetchAssignRoleMenu } from '@/api/system-manage'
  import { ElMessage } from 'element-plus'

  interface Props {
    modelValue: boolean
    roleData?: { id: string }
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    roleData: undefined
  })

  const emit = defineEmits<Emits>()

  const treeRef = ref()
  const isExpandAll = ref(true)
  const isSelectAll = ref(false)
  const loading = ref(false)
  // 角色已分配的菜单 ID（来自后端）
  const roleMenuIds = ref<string[]>([])
  // 操作父子联动配置（1 联动，2 不联动）
  const menuLinkage = ref<number>(1)

  // 后端菜单树（树数据来源）
  const backendMenus = ref<Api.SystemManage.MenuTreeNode[]>([])

  /**
   * 弹窗显示状态双向绑定
   */
  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  /**
   * 树节点显示文案：备注优先，否则权限代码
   */
  const menuNodeLabel = (data: Api.SystemManage.MenuTreeNode) =>
    data.remark || data.permis || data.id || ''

  /**
   * 节点自定义 class：权限按钮(types=2)横向排；其父节点标记以便子容器 flex 横向
   */
  const getNodeClass = (data: Api.SystemManage.MenuTreeNode) => {
    if (data.types === 2) return 'is-permission-item'
    const children = data.children
    if (children?.length && children.every((c: Api.SystemManage.MenuTreeNode) => c.types === 2)) {
      return 'has-permission-children'
    }
    return ''
  }

  /**
   * 树形组件配置
   */
  const defaultProps = {
    children: 'children',
    label: (data: any) => menuNodeLabel(data as Api.SystemManage.MenuTreeNode),
    class: (data: any) => getNodeClass(data as Api.SystemManage.MenuTreeNode)
  }

  /**
   * 递归收集后端树所有节点 id
   */
  const getAllBackendNodeIds = (nodes: Api.SystemManage.MenuTreeNode[]): string[] => {
    const ids: string[] = []
    const traverse = (list: Api.SystemManage.MenuTreeNode[]) => {
      list.forEach((node) => {
        if (node.id) ids.push(node.id)
        if (node.children?.length) traverse(node.children)
      })
    }
    traverse(nodes)
    return ids
  }

  /**
   * 根据一组 ID，在当前菜单树中只保留叶子节点的 ID
   * 用于初始化勾选状态，避免把仅用于标识“半选”的父节点也设为完全勾选
   */
  const getLeafMenuIdsFromTree = (ids: string[]): string[] => {
    if (!ids.length || !backendMenus.value.length) return []
    const idSet = new Set(ids.map((id) => String(id)))
    const leafIds: string[] = []

    const traverse = (nodes: Api.SystemManage.MenuTreeNode[]) => {
      nodes.forEach((node) => {
        const nodeId = node.id != null ? String(node.id) : ''
        const children = node.children || []
        const isLeaf = !children.length

        if (isLeaf && nodeId && idSet.has(nodeId)) {
          leafIds.push(nodeId)
        }

        if (children.length) {
          traverse(children)
        }
      })
    }

    traverse(backendMenus.value)
    return leafIds
  }

  /**
   * 加载后端菜单树（只需加载一次）
   * 严格按 swagger：/sys_menu/all 返回 { rows, total }（internal_sys_menu.ListResp）
   */
  const loadBackendMenus = async (): Promise<void> => {
    if (backendMenus.value.length) return
    try {
      const res = await fetchGetMenuAll()
      const rows = res?.rows ?? []
      backendMenus.value = handleTree(
        rows as unknown as Record<string, unknown>[],
        'id',
        'pid',
        'children'
      ) as unknown as Api.SystemManage.MenuTreeNode[]
    } catch (error) {
      console.error('获取菜单列表失败:', error)
      ElMessage.error('获取菜单列表失败')
    }
  }

  /**
   * 加载角色权限数据
   */
  const loadRolePermission = async (): Promise<void> => {
    if (!props.roleData?.id) {
      roleMenuIds.value = []
      return
    }

    try {
      loading.value = true
      const detail: any = await fetchGetRoleInfo(props.roleData.id)
      const rawIds: string[] = (detail.menuIds || []).map((id: string) => String(id))
      menuLinkage.value = detail.menuLinkage ?? 1

      // 根据联动模式决定如何恢复勾选：
      // - 联动模式(1)：只勾叶子节点，父级由组件自动计算半选/全选
      // - 不联动模式(2)：按后端返回的原始 ID 全量勾选（包含父级），让父级在“非联动”模式下也能单独表现为选中
      roleMenuIds.value =
        menuLinkage.value === 2 ? Array.from(new Set(rawIds)) : getLeafMenuIdsFromTree(rawIds)
      // 依赖 roleMenuIds 和后端菜单数据的 defaultCheckedKeys
      // 会自动根据当前 menuIds 初始化 authCheckedMap 和树的默认勾选状态
    } catch (error) {
      console.error('获取角色权限失败:', error)
      ElMessage.error('获取角色权限失败')
    } finally {
      loading.value = false
    }
  }

  /**
   * 监听弹窗打开，从后端加载菜单树与角色已分配权限
   */
  watch(
    () => props.modelValue,
    async (newVal) => {
      if (newVal) {
        loading.value = true
        try {
          await loadBackendMenus()
          await loadRolePermission()
          await nextTick()
          treeRef.value?.setCheckedKeys(roleMenuIds.value)
        } finally {
          loading.value = false
        }
      } else {
        roleMenuIds.value = []
        treeRef.value?.setCheckedKeys([])
      }
    }
  )

  /**
   * 关闭弹窗并清空选中状态
   */
  const handleClose = () => {
    visible.value = false
    treeRef.value?.setCheckedKeys([])
  }

  /**
   * 保存权限配置：将当前树勾选的节点 id 提交给后端
   */
  const savePermission = async () => {
    const tree = treeRef.value
    if (!tree) return

    if (!props.roleData?.id) {
      ElMessage.error('角色信息不存在，无法保存权限')
      return
    }

    // 同时包含勾选和半选中的父级节点 ID，一并提交给后端
    const checkedKeys = tree.getCheckedKeys() as string[]
    const halfCheckedKeys = tree.getHalfCheckedKeys() as string[]
    const menuIds = Array.from(new Set([...checkedKeys, ...halfCheckedKeys]))

    try {
      loading.value = true
      await fetchAssignRoleMenu({
        roleId: props.roleData.id,
        menuIds,
        menuLinkage: menuLinkage.value || 1
      })

      ElMessage.success('权限保存成功')
      emit('success')
      handleClose()
    } catch (error) {
      console.error('保存权限失败:', error)
      ElMessage.error(error instanceof Error ? error.message : '保存权限失败')
    } finally {
      loading.value = false
    }
  }

  /**
   * 切换全部展开/收起状态
   */
  const toggleExpandAll = () => {
    const tree = treeRef.value
    if (!tree) return

    const nodes = tree.store.nodesMap
    // 这里保留 any，因为 Element Plus 的内部节点类型较复杂
    Object.values(nodes).forEach((node: any) => {
      node.expanded = !isExpandAll.value
    })

    isExpandAll.value = !isExpandAll.value
  }

  /**
   * 切换全选/取消全选状态
   */
  const toggleSelectAll = () => {
    const tree = treeRef.value
    if (tree) {
      if (!isSelectAll.value) {
        const allIds = getAllBackendNodeIds(backendMenus.value)
        tree.setCheckedKeys(allIds)
      } else {
        tree.setCheckedKeys([])
      }
    }

    isSelectAll.value = !isSelectAll.value
  }

  /**
   * 处理树节点选中状态变化，同步更新全选按钮状态
   */
  const handleTreeCheck = () => {
    const tree = treeRef.value
    if (!tree) return

    const checkedKeys = tree.getCheckedKeys()
    const allIds = getAllBackendNodeIds(backendMenus.value)
    isSelectAll.value = checkedKeys.length === allIds.length && allIds.length > 0
  }
</script>

<style scoped lang="scss">
  /* 抽屉内容整体占满高度，树区域自适应剩余空间 */
  .role-menu-drawer-body {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .role-menu-drawer-scroll {
    flex: 1;
    min-height: 0;
  }

  .role-menu-linkage-row {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    padding-top: 8px;
  }

  .linkage-label {
    margin-right: 8px;
  }

  /* 操作权限树：按钮权限横向排列，每个按钮浅色背景 */
  .role-menu-tree {
    :deep(.el-tree-node__content) {
      min-height: 32px;
      padding-top: 4px;
      padding-bottom: 4px;
    }

    /* 菜单节点下的“操作权限”子节点组：整组浅色背景，突出这是一组操作按钮 */
    :deep(.has-permission-children) {
      > .el-tree-node__children {
        display: flex;
        flex-wrap: wrap;
        gap: 8px 16px;
        align-items: center;
        padding: 8px 12px;
        margin-left: 12px;
        background: var(--el-fill-color-lighter);
        border-radius: 4px;
      }
    }

    :deep(.is-permission-item) {
      display: inline-flex;
      width: auto;

      .el-tree-node__content {
        padding: 0;
        margin-left: 0;
        background-color: transparent !important;
      }

      .el-tree-node__expand-icon {
        display: none;
      }
    }

    /* 操作按钮文本的实际 hover 区域，背景只在文字 + 复选框附近，不占用左侧空白 */
    :deep(.role-menu-permission-label) {
      display: inline-block;
      padding: 4px 8px;
      border-radius: 4px;
      transition: background-color 0.2s ease;
    }

    :deep(.is-permission-item .el-tree-node__content:hover .role-menu-permission-label) {
      background-color: var(--el-color-primary-light-9);
    }
  }
</style>
