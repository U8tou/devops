<template>
  <ElDrawer v-model="visible" title="绑定部门" size="520px" direction="rtl" @close="handleClose">
    <div class="user-dept-drawer-body">
      <ElScrollbar class="user-dept-drawer-scroll">
        <ElTree
          ref="treeRef"
          :data="deptTree"
          show-checkbox
          node-key="id"
          :default-expand-all="isExpandAll"
          :default-checked-keys="defaultCheckedKeys"
          :props="defaultProps"
          :check-strictly="deptLinkage === 2"
          @check="handleTreeCheck"
        >
          <template #default="{ data }">
            <div style="display: flex; align-items: center">
              <span>{{ data.name }}</span>
            </div>
          </template>
        </ElTree>
      </ElScrollbar>

      <!-- 父子联动配置 -->
      <div class="user-dept-linkage-row">
        <span class="linkage-label">父子节点联动：</span>
        <ElRadioGroup v-model="deptLinkage" size="small">
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
  import { fetchGetDeptAll, fetchGetUserDetail, fetchAssignUserDept } from '@/api/system-manage'
  import { handleTree } from '@/utils/ruoyi'
  import { ElMessage } from 'element-plus'

  interface Props {
    modelValue: boolean
    userData?: { id: string }
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    userData: undefined
  })

  const emit = defineEmits<Emits>()

  const treeRef = ref()
  const isExpandAll = ref(true)
  const isSelectAll = ref(false)
  const loading = ref(false)
  // 数据父子联动配置（1 联动，2 不联动），默认不联动
  const deptLinkage = ref<number>(2)

  // 部门树数据
  const deptTree = ref<Api.SystemManage.DeptTreeNode[]>([])
  // 用户已分配的部门 ID（字符串数组，与接口一致）
  const userDeptIds = ref<string[]>([])

  /**
   * 弹窗显示状态双向绑定
   */
  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  /**
   * 树形组件配置
   */
  const defaultProps = {
    children: 'children',
    label: 'name'
  }

  /**
   * 默认选中的部门 keys
   */
  const defaultCheckedKeys = computed(() => {
    return userDeptIds.value
  })

  /**
   * 根据一组 ID，在当前部门树中只保留叶子节点的 ID
   * 用于在“联动模式”下初始化勾选状态，避免把仅用于标识“半选”的父节点也设为完全勾选
   */
  const getLeafDeptIdsFromTree = (ids: string[]): string[] => {
    if (!ids.length || !deptTree.value.length) return []
    const idSet = new Set(ids.map((id) => String(id)))
    const leafIds: string[] = []

    const traverse = (nodes: Api.SystemManage.DeptTreeNode[]) => {
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

    traverse(deptTree.value)
    return leafIds
  }

  /**
   * 加载部门树数据
   */
  const loadDeptTree = async (): Promise<void> => {
    try {
      loading.value = true
      const res = await fetchGetDeptAll()
      const rows = res?.rows ?? []
      deptTree.value = handleTree(
        rows as unknown as Record<string, unknown>[],
        'id',
        'pid',
        'children'
      ) as unknown as Api.SystemManage.DeptTreeNode[]
    } catch (error) {
      console.error('获取部门列表失败:', error)
      ElMessage.error('获取部门列表失败')
    } finally {
      loading.value = false
    }
  }

  /**
   * 加载用户已分配的部门
   * 根据 swagger，用户详情接口返回 depts 字段（关联部门IDs）
   */
  const loadUserDepts = async (): Promise<void> => {
    if (!props.userData?.id) {
      userDeptIds.value = []
      return
    }

    try {
      loading.value = true
      const detail = await fetchGetUserDetail(props.userData.id)
      // 严格按 swagger：用户详情返回 depts（关联部门IDs）
      const rawIds: string[] = Array.isArray(detail.depts)
        ? detail.depts.map((id: string | number) => String(id))
        : []

      // 根据当前联动模式决定如何恢复勾选：
      // - 联动模式(1)：只勾叶子节点，父级由组件自动计算半选/全选
      // - 不联动模式(2)：按后端返回的原始 ID 全量勾选（包含父级），让父级在“非联动”模式下也能单独表现为选中
      userDeptIds.value =
        deptLinkage.value === 2 ? Array.from(new Set(rawIds)) : getLeafDeptIdsFromTree(rawIds)

      // 设置默认选中
      nextTick(() => {
        if (treeRef.value && userDeptIds.value.length > 0) {
          treeRef.value.setCheckedKeys(userDeptIds.value)
        }
      })
    } catch (error) {
      console.error('获取用户部门失败:', error)
      ElMessage.error('获取用户部门失败')
      userDeptIds.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * 监听弹窗打开，初始化部门树和用户部门数据
   */
  watch(
    () => props.modelValue,
    async (newVal) => {
      if (newVal) {
        await loadDeptTree()
        await loadUserDepts()
      } else {
        userDeptIds.value = []
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
   * 保存部门权限配置
   * - 不联动模式（deptLinkage === 2）：只提交当前勾选节点
   * - 联动模式（deptLinkage === 1）：提交勾选节点 + 半选中的上级节点
   */
  const savePermission = async () => {
    const tree = treeRef.value
    if (!tree) return

    if (!props.userData?.id) {
      ElMessage.error('用户信息不存在，无法保存部门权限')
      return
    }

    const checkedKeys = tree.getCheckedKeys() as Array<string | number>
    const halfCheckedKeys =
      deptLinkage.value === 1 ? (tree.getHalfCheckedKeys?.() as Array<string | number>) || [] : []

    const allIds = [...checkedKeys, ...halfCheckedKeys]
      .map((id) => String(id))
      .filter((id, index, arr) => id && id !== '' && arr.indexOf(id) === index)

    const deptIds = allIds

    try {
      loading.value = true
      // 根据 swagger，deptIds 和 userId 都是必需的
      await fetchAssignUserDept({
        userId: props.userData.id,
        deptIds
      })

      ElMessage.success('部门权限保存成功')
      emit('success')
      handleClose()
    } catch (error) {
      console.error('保存部门权限失败:', error)
      ElMessage.error(error instanceof Error ? error.message : '保存部门权限失败')
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
    Object.values(nodes).forEach((node: any) => {
      node.expanded = !isExpandAll.value
    })

    isExpandAll.value = !isExpandAll.value
  }

  /**
   * 递归获取所有节点的 key（这里使用部门 id）
   * @param nodes 节点列表
   * @returns 所有节点的 key 数组
   */
  const getAllNodeKeys = (nodes: Api.SystemManage.DeptTreeNode[]): string[] => {
    const keys: string[] = []
    const traverse = (nodeList: Api.SystemManage.DeptTreeNode[]): void => {
      nodeList.forEach((node) => {
        if (node.id) keys.push(String(node.id))
        if (node.children?.length) traverse(node.children)
      })
    }
    traverse(nodes)
    return keys
  }

  /**
   * 切换全选/取消全选状态
   */
  const toggleSelectAll = () => {
    const tree = treeRef.value
    if (!tree) return

    if (!isSelectAll.value) {
      const allKeys = getAllNodeKeys(deptTree.value)
      tree.setCheckedKeys(allKeys)
    } else {
      tree.setCheckedKeys([])
    }

    isSelectAll.value = !isSelectAll.value
  }
  /**
   * 处理树节点选中状态变化
   * 同步更新全选按钮状态
   */
  const handleTreeCheck = () => {
    const tree = treeRef.value
    if (!tree) return

    const checkedKeys = tree.getCheckedKeys()
    const allKeys = getAllNodeKeys(deptTree.value)

    isSelectAll.value = checkedKeys.length === allKeys.length && allKeys.length > 0
  }
</script>

<style scoped lang="scss">
  /* 抽屉内容整体占满高度，树区域自适应剩余空间 */
  .user-dept-drawer-body {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .user-dept-drawer-scroll {
    flex: 1;
    min-height: 0;
  }

  .user-dept-linkage-row {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    padding-top: 8px;
  }

  .linkage-label {
    margin-right: 8px;
  }
</style>
