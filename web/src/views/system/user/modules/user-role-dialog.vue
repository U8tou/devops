<template>
  <ElDrawer v-model="visible" title="绑定角色" size="520px" direction="rtl" @close="handleClose">
    <div class="user-role-drawer-body">
      <ElScrollbar class="user-role-drawer-scroll">
        <ElTree
          ref="treeRef"
          :data="roleTree"
          show-checkbox
          node-key="id"
          :default-expand-all="isExpandAll"
          :default-checked-keys="defaultCheckedKeys"
          :props="defaultProps"
          check-strictly
        >
          <template #default="{ data }">
            <div style="display: flex; align-items: center">
              <span>{{ data.name }}</span>
            </div>
          </template>
        </ElTree>
      </ElScrollbar>
    </div>

    <template #footer>
      <ElButton @click="toggleSelectAll">{{ isSelectAll ? '取消全选' : '全部选择' }}</ElButton>
      <ElButton type="primary" @click="savePermission">保存</ElButton>
    </template>
  </ElDrawer>
</template>

<script setup lang="ts">
  import { fetchGetRoleAll, fetchGetUserDetail, fetchAssignUserRole } from '@/api/system-manage'
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

  // 角色列表数据（扁平结构）
  const roleList = ref<Api.SystemManage.RoleListItem[]>([])
  // 用户已分配的角色 ID（字符串数组，与接口一致）
  const userRoleIds = ref<string[]>([])

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
   * 角色树数据（扁平列表转为树形，这里角色是扁平的，不需要树形结构）
   */
  const roleTree = computed(() => {
    return roleList.value.map((role) => ({
      id: String(role.id),
      name: role.name,
      role: role.role
    }))
  })

  /**
   * 默认选中的角色 keys
   */
  const defaultCheckedKeys = computed(() => {
    return userRoleIds.value
  })

  /**
   * 加载角色列表
   */
  const loadRoleList = async (): Promise<void> => {
    try {
      loading.value = true
      const roles = await fetchGetRoleAll()
      roleList.value = roles || []
    } catch (error) {
      console.error('获取角色列表失败:', error)
      ElMessage.error('获取角色列表失败')
    } finally {
      loading.value = false
    }
  }

  /**
   * 加载用户已分配的角色
   * 根据 swagger，用户详情接口返回 roles 字段（关联角色IDs）
   */
  const loadUserRoles = async (): Promise<void> => {
    if (!props.userData?.id) {
      userRoleIds.value = []
      return
    }

    try {
      loading.value = true
      const detail = await fetchGetUserDetail(props.userData.id)
      // 严格按 swagger：用户详情返回 roles（关联角色IDs）
      userRoleIds.value = Array.isArray(detail.roles) ? detail.roles.map((id) => String(id)) : []

      // 设置默认选中
      nextTick(() => {
        if (treeRef.value && userRoleIds.value.length > 0) {
          treeRef.value.setCheckedKeys(userRoleIds.value)
        }
      })
    } catch (error) {
      console.error('获取用户角色失败:', error)
      ElMessage.error('获取用户角色失败')
      userRoleIds.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * 监听弹窗打开，初始化角色列表和用户角色数据
   */
  watch(
    () => props.modelValue,
    async (newVal) => {
      if (newVal) {
        await loadRoleList()
        await loadUserRoles()
      } else {
        userRoleIds.value = []
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
   * 保存角色权限配置
   */
  const savePermission = async () => {
    const tree = treeRef.value
    if (!tree) return

    if (!props.userData?.id) {
      ElMessage.error('用户信息不存在，无法保存角色权限')
      return
    }

    const checkedKeys = tree.getCheckedKeys() as string[]
    const roleIds = checkedKeys.filter((id) => id && id !== '').map((id) => String(id))

    try {
      loading.value = true
      // 根据 swagger，roleIds 和 userId 都是必需的
      await fetchAssignUserRole({
        userId: props.userData.id,
        roleIds
      })

      ElMessage.success('角色权限保存成功')
      emit('success')
      handleClose()
    } catch (error) {
      console.error('保存角色权限失败:', error)
      ElMessage.error(error instanceof Error ? error.message : '保存角色权限失败')
    } finally {
      loading.value = false
    }
  }

  /**
   * 切换全选/取消全选状态
   */
  const toggleSelectAll = () => {
    const tree = treeRef.value
    if (!tree) return

    if (!isSelectAll.value) {
      const allKeys = roleTree.value.map((role) => String(role.id))
      tree.setCheckedKeys(allKeys)
    } else {
      tree.setCheckedKeys([])
    }

    isSelectAll.value = !isSelectAll.value
  }
</script>

<style scoped lang="scss">
  /* 抽屉内容整体占满高度，树区域自适应剩余空间 */
  .user-role-drawer-body {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .user-role-drawer-scroll {
    flex: 1;
    min-height: 0;
  }

  /* 目前不需要额外的父子联动配置样式 */
</style>
