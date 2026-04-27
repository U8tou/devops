<template>
  <ElDrawer v-model="visible" title="绑定岗位" size="520px" direction="rtl" @close="handleClose">
    <div class="user-post-drawer-body">
      <ElScrollbar class="user-post-drawer-scroll">
        <ElTree
          ref="treeRef"
          :data="postTree"
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
  import { fetchGetPostAll, fetchGetUserDetail, fetchAssignUserPost } from '@/api/system-manage'
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

  // 岗位列表数据（扁平结构）
  const postList = ref<Api.SystemManage.PostListItem[]>([])
  // 用户已分配的岗位 ID（字符串数组，与接口一致）
  const userPostIds = ref<string[]>([])

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
   * 岗位树数据（扁平列表，岗位是扁平的）
   */
  const postTree = computed(() => {
    return postList.value.map((post) => ({
      id: String(post.id),
      name: post.name
    }))
  })

  /**
   * 默认选中的岗位 keys
   */
  const defaultCheckedKeys = computed(() => {
    return userPostIds.value
  })

  /**
   * 加载岗位列表
   */
  const loadPostList = async (): Promise<void> => {
    try {
      loading.value = true
      const res = await fetchGetPostAll()
      postList.value = res?.rows ?? []
    } catch (error) {
      console.error('获取岗位列表失败:', error)
      ElMessage.error('获取岗位列表失败')
    } finally {
      loading.value = false
    }
  }

  /**
   * 加载用户已分配的岗位
   * 用户详情接口返回 posts 字段（关联岗位IDs）
   */
  const loadUserPosts = async (): Promise<void> => {
    if (!props.userData?.id) {
      userPostIds.value = []
      return
    }

    try {
      loading.value = true
      const detail = await fetchGetUserDetail(props.userData.id)
      userPostIds.value = Array.isArray(detail.posts) ? detail.posts.map((id) => String(id)) : []

      nextTick(() => {
        if (treeRef.value && userPostIds.value.length > 0) {
          treeRef.value.setCheckedKeys(userPostIds.value)
        }
      })
    } catch (error) {
      console.error('获取用户岗位失败:', error)
      ElMessage.error('获取用户岗位失败')
      userPostIds.value = []
    } finally {
      loading.value = false
    }
  }

  /**
   * 监听弹窗打开，初始化岗位列表和用户岗位数据
   */
  watch(
    () => props.modelValue,
    async (newVal) => {
      if (newVal) {
        await loadPostList()
        await loadUserPosts()
      } else {
        userPostIds.value = []
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
   * 保存岗位配置
   */
  const savePermission = async () => {
    const tree = treeRef.value
    if (!tree) return

    if (!props.userData?.id) {
      ElMessage.error('用户信息不存在，无法保存岗位')
      return
    }

    const checkedKeys = tree.getCheckedKeys() as string[]
    const postIds = checkedKeys.filter((id) => id && id !== '').map((id) => String(id))

    try {
      loading.value = true
      await fetchAssignUserPost({
        userId: props.userData.id,
        postIds
      })

      ElMessage.success('岗位保存成功')
      emit('success')
      handleClose()
    } catch (error) {
      console.error('保存岗位失败:', error)
      ElMessage.error(error instanceof Error ? error.message : '保存岗位失败')
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
      const allKeys = postTree.value.map((post) => String(post.id))
      tree.setCheckedKeys(allKeys)
    } else {
      tree.setCheckedKeys([])
    }

    isSelectAll.value = !isSelectAll.value
  }
</script>

<style scoped lang="scss">
  .user-post-drawer-body {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .user-post-drawer-scroll {
    flex: 1;
    min-height: 0;
  }
</style>
