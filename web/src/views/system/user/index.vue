<!-- 左树右表示例页面 -->
<template>
  <div class="art-full-height">
    <div class="box-border flex gap-4 h-full max-md:block max-md:gap-0 max-md:h-auto">
      <div class="flex-shrink-0 w-58 h-full max-md:w-full max-md:h-auto max-md:mb-5">
        <ElCard class="tree-card art-card-xs flex flex-col h-full mt-0" shadow="never">
          <template #header>
            <b>{{ $t('sys.common.dept') }}</b>
          </template>
          <ElScrollbar>
            <ElTree
              :data="treeData"
              :props="treeProps"
              node-key="id"
              default-expand-all
              highlight-current
              @node-click="handleNodeClick"
            />
          </ElScrollbar>
        </ElCard>
      </div>

      <div class="flex flex-col flex-grow min-w-0">
        <UserSearch v-model="searchForm" @search="handleSearch" @reset="handleReset" />

        <ElCard class="flex flex-col flex-1 min-h-0 art-table-card" shadow="never">
          <ArtTableHeader v-model:columns="columnChecks" :loading="loading" @refresh="refreshData">
            <template #left>
              <ElSpace wrap>
                <ElButton v-auth="'sys:user:add'" @click="showDialog('add')" v-ripple>{{
                  $t('sys.user.addUser')
                }}</ElButton>
                <ElButton
                  v-auth="'sys:user:export'"
                  :loading="exportLoading"
                  @click="handleExport"
                  v-ripple
                >
                  {{ $t('sys.common.export') }}
                </ElButton>
                <ElButton v-auth="'sys:user:import'" @click="importDialogVisible = true" v-ripple>
                  {{ $t('sys.common.import') }}
                </ElButton>
              </ElSpace>
            </template>
          </ArtTableHeader>

          <ArtTable
            rowKey="id"
            :loading="loading"
            :data="data"
            :columns="columns"
            :pagination="pagination"
            @pagination:size-change="handleSizeChange"
            @pagination:current-change="handleCurrentChange"
          >
            <template #userInfo="{ row }">
              <div style="display: flex; gap: 12px; align-items: center">
                <!-- 头像 -->
                <div style="flex-shrink: 0">
                  <img
                    v-if="row.avatar && !avatarErrorMap[row.id]"
                    :src="resolveImageUrl(row.avatar)"
                    :alt="row.nickName || row.userName"
                    style="
                      display: block;
                      width: 40px;
                      height: 40px;
                      object-fit: cover;
                      border-radius: 4px;
                    "
                    @error="(e) => handleAvatarError(e, row)"
                  />
                  <div
                    v-else
                    style="
                      display: flex;
                      align-items: center;
                      justify-content: center;
                      width: 40px;
                      height: 40px;
                      font-size: 16px;
                      font-weight: 600;
                      color: #fff;
                      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
                      border-radius: 4px;
                    "
                  >
                    {{ row.nickName?.charAt(0) || row.userName?.charAt(0) || 'U' }}
                  </div>
                </div>
                <!-- 用户信息 -->
                <div style="display: flex; flex: 1; flex-direction: column; gap: 2px; min-width: 0">
                  <div style="font-weight: 500; color: #303133; word-break: break-all">
                    {{ row.nickName || '-' }}
                  </div>
                  <div style="font-size: 12px; color: #909399; word-break: break-all">
                    {{ row.userName || '-' }}
                  </div>
                </div>
              </div>
            </template>
            <template #contact="{ row }">
              <div style="display: flex; flex-direction: column; gap: 2px">
                <div style="display: flex; gap: 6px; align-items: center; color: #303133">
                  <ArtSvgIcon
                    icon="ri:phone-line"
                    style="flex-shrink: 0; font-size: 14px; color: #909399"
                  />
                  <span>{{ row.phone || '-' }}</span>
                </div>
                <div
                  style="
                    display: flex;
                    gap: 6px;
                    align-items: center;
                    font-size: 12px;
                    color: #909399;
                  "
                >
                  <ArtSvgIcon
                    icon="ri:mail-line"
                    style="flex-shrink: 0; font-size: 14px; color: #909399"
                  />
                  <span>{{ row.email || '-' }}</span>
                </div>
              </div>
            </template>
            <template #timeInfo="{ row }">
              <div style="display: flex; flex-direction: column; gap: 2px">
                <div style="display: flex; gap: 4px; align-items: center; color: #303133">
                  <ArtSvgIcon
                    icon="ri:calendar-line"
                    style="flex-shrink: 0; font-size: 14px; color: #909399"
                  />
                  <span>{{ row.createTime ? formatTime(row.createTime) : '-' }}</span>
                </div>
                <div
                  style="
                    display: flex;
                    gap: 4px;
                    align-items: center;
                    font-size: 12px;
                    color: #909399;
                  "
                >
                  <ArtSvgIcon
                    icon="ri:refresh-line"
                    style="flex-shrink: 0; font-size: 14px; color: #909399"
                  />
                  <span>{{ row.updateTime ? formatTime(row.updateTime) : '-' }}</span>
                </div>
              </div>
            </template>
            <template #deptId="{ row }">
              <div
                v-if="getDeptNames(row.depts).length > 0"
                style="display: flex; flex-wrap: wrap; gap: 4px; align-items: center"
              >
                <ElTag
                  v-for="(deptName, index) in getDeptNames(row.depts)"
                  :key="index"
                  type="success"
                  size="small"
                >
                  {{ deptName }}
                </ElTag>
              </div>
              <span v-else style="font-size: 12px; color: #909399">-</span>
            </template>
            <template #roles="{ row }">
              <div
                v-if="getRoleNames(row.roles).length > 0"
                style="display: flex; flex-wrap: wrap; gap: 4px; align-items: center"
              >
                <ElTag
                  v-for="(roleName, index) in getRoleNames(row.roles)"
                  :key="index"
                  type="info"
                  size="small"
                >
                  {{ roleName }}
                </ElTag>
              </div>
              <span v-else style="font-size: 12px; color: #909399">-</span>
            </template>
            <template #posts="{ row }">
              <div
                v-if="getPostNames(row.posts).length > 0"
                style="display: flex; flex-wrap: wrap; gap: 4px; align-items: center"
              >
                <ElTag
                  v-for="(postName, index) in getPostNames(row.posts)"
                  :key="index"
                  type="warning"
                  size="small"
                >
                  {{ postName }}
                </ElTag>
              </div>
              <span v-else style="font-size: 12px; color: #909399">-</span>
            </template>
            <template #operation="{ row }">
              <div style="display: flex; align-items: center; justify-content: flex-end">
                <ArtButtonTable
                  v-auth="'sys:user:edit'"
                  type="edit"
                  @click="showDetailDialog(row)"
                />
                <ArtButtonMore
                  :list="moreButtonList"
                  @click="(item) => handleMoreClick(item, row)"
                />
              </div>
            </template>
          </ArtTable>
        </ElCard>
      </div>
    </div>

    <!-- 用户新增弹窗 -->
    <UserDialog
      v-model:visible="dialogVisible"
      :type="dialogType"
      :user-data="currentUserData"
      @submit="handleDialogSubmit"
    />

    <!-- 用户详情弹窗 -->
    <UserDetailDialog
      v-model="detailDialogVisible"
      :user-data="currentUserData"
      @success="refreshData"
    />

    <!-- 用户绑定角色弹窗 -->
    <UserRoleDialog
      v-model="roleDialogVisible"
      :user-data="currentUserData"
      @success="refreshData"
    />

    <!-- 用户绑定部门弹窗 -->
    <UserDeptDialog
      v-model="deptDialogVisible"
      :user-data="currentUserData"
      @success="refreshData"
    />

    <!-- 用户绑定岗位弹窗 -->
    <UserPostDialog
      v-model="postDialogVisible"
      :user-data="currentUserData"
      @success="refreshData"
    />

    <!-- 用户重置密码弹窗 -->
    <ElDialog
      v-model="resetPwdDialogVisible"
      :title="$t('sys.user.resetPwd')"
      width="400px"
      align-center
    >
      <div>
        <p class="mb-3 text-sm text-g-600">
          {{
            $t('sys.user.resetPwdTitle', {
              name: resetPwdTargetUser?.nickName || resetPwdTargetUser?.userName || '-'
            })
          }}
        </p>
        <ElInput
          v-model="resetPwdForm.password"
          type="password"
          show-password
          :placeholder="$t('sys.user.resetPwdPlaceholder')"
        />
      </div>
      <template #footer>
        <ElButton @click="resetPwdDialogVisible = false" :disabled="resetPwdLoading">{{
          $t('common.cancel')
        }}</ElButton>
        <ElButton type="primary" @click="handleResetPwdConfirm" :loading="resetPwdLoading">
          {{ $t('common.confirm') }}
        </ElButton>
      </template>
    </ElDialog>

    <!-- 用户导入弹窗 -->
    <ElDialog
      v-model="importDialogVisible"
      :title="$t('sys.user.importTitle')"
      width="480px"
      destroy-on-close
      @closed="clearImportFile"
    >
      <ElUpload
        drag
        :auto-upload="false"
        accept=".xlsx,.xls"
        :limit="1"
        :file-list="importFileList"
        :on-change="onImportFileChange"
        :on-remove="onImportFileRemove"
      >
        <div class="import-upload-inner">
          <ElIcon class="import-upload-icon"><UploadFilled /></ElIcon>
          <div class="import-upload-text" v-html="$t('sys.user.importDragTip')"></div>
        </div>
      </ElUpload>
      <div class="import-options">
        <ElCheckbox v-model="importUpdateExisting">{{
          $t('sys.user.importUpdateExist')
        }}</ElCheckbox>
        <p class="import-format-hint">{{ $t('sys.user.importFormatTip') }}</p>
        <ElLink type="primary" :underline="false" @click="handleDownloadTemplate">
          {{ $t('sys.common.downloadTemplate') }}
        </ElLink>
      </div>
      <template #footer>
        <ElButton @click="importDialogVisible = false">{{ $t('common.cancel') }}</ElButton>
        <ElButton type="primary" :loading="importLoading" @click="handleImportConfirm">
          {{ $t('common.confirm') }}
        </ElButton>
      </template>
    </ElDialog>

    <!-- 导入失败明细弹窗 -->
    <ElDialog
      v-model="importFailDialogVisible"
      :title="$t('sys.user.importFailDetail')"
      width="500px"
      destroy-on-close
    >
      <ElTable :data="importFailList" border max-height="300">
        <ElTableColumn prop="row" :label="$t('sys.common.rowNo')" width="80" align="center" />
        <ElTableColumn
          prop="reason"
          :label="$t('sys.common.failReason')"
          min-width="200"
          show-overflow-tooltip
        />
      </ElTable>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { reactive, h } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useTable } from '@/hooks/core/useTable'
  import {
    fetchGetUserList,
    fetchGetDeptAll,
    fetchDeleteUser,
    fetchResetUserPassword,
    fetchGetRoleAll,
    fetchGetPostAll,
    fetchExportUser,
    fetchImportUser,
    fetchGetUserTemplate
  } from '@/api/system-manage'
  import UserSearch from '@/views/system/user/modules/user-search.vue'
  import { DialogType } from '@/types'
  import UserDialog from './modules/user-dialog.vue'
  import UserDetailDialog from './modules/user-detail-dialog.vue'
  import UserRoleDialog from './modules/user-role-dialog.vue'
  import UserDeptDialog from './modules/user-dept-dialog.vue'
  import UserPostDialog from './modules/user-post-dialog.vue'
  import { resolveImageUrl } from '@/utils/resolveImageUrl'
  import { handleTree, parseTime, blobValidate } from '@/utils/ruoyi'
  import FileSaver from 'file-saver'
  import type { UploadFile } from 'element-plus'
  import { ButtonMoreItem } from '@/components/core/forms/art-button-more/index.vue'
  import ArtButtonMore from '@/components/core/forms/art-button-more/index.vue'
  import ArtButtonTable from '@/components/core/forms/art-button-table/index.vue'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'
  import { ElMessageBox, ElMessage, ElTag, ElLink, ElCheckbox } from 'element-plus'
  import { UploadFilled } from '@element-plus/icons-vue'

  type UserListItem = Api.SystemManage.UserListItem

  defineOptions({ name: 'TreeTable' })

  // 弹窗相关
  const dialogType = ref<DialogType>('add')
  const dialogVisible = ref(false)
  const detailDialogVisible = ref(false)
  const roleDialogVisible = ref(false)
  const deptDialogVisible = ref(false)
  const postDialogVisible = ref(false)
  const currentUserData = ref<{ id: string } | undefined>(undefined)

  // 重置密码弹窗
  const resetPwdDialogVisible = ref(false)
  const resetPwdLoading = ref(false)
  const resetPwdTargetUser = ref<UserListItem | null>(null)
  const resetPwdForm = reactive({
    password: ''
  })

  // 导入导出
  const exportLoading = ref(false)
  const templateLoading = ref(false)
  const importDialogVisible = ref(false)
  const importFileList = ref<UploadFile[]>([])
  const importFileRef = ref<File | null>(null)
  const importUpdateExisting = ref(false)
  const importLoading = ref(false)
  const importFailDialogVisible = ref(false)
  const importFailList = ref<Api.SystemManage.UserImportFail[]>([])

  const { t } = useI18n()
  // 更多操作按钮列表
  const moreButtonList = computed(() => [
    { key: 'role', label: t('sys.user.bindRole'), icon: 'ri:user-star-line' },
    { key: 'dept', label: t('sys.user.bindDept'), icon: 'ri:organization-chart' },
    { key: 'post', label: t('sys.user.bindPost'), icon: 'ri:briefcase-line' },
    {
      key: 'resetPwd',
      label: t('sys.user.resetPwd'),
      icon: 'ri:key-2-line',
      auth: 'sys:user:reset_pwd'
    },
    {
      key: 'delete',
      label: t('sys.user.deleteUser'),
      icon: 'ri:delete-bin-4-line',
      color: '#f56c6c',
      auth: 'sys:user:del'
    }
  ])

  // 部门树数据（来自接口）
  const treeData = ref<Api.SystemManage.DeptTreeNode[]>([])
  const selectedDeptId = ref<string | undefined>(undefined)

  // 角色列表（用于角色名称映射）
  const roleList = ref<Api.SystemManage.RoleListItem[]>([])
  // 岗位列表（用于岗位名称映射）
  const postList = ref<Api.SystemManage.PostListItem[]>([])

  const treeProps = {
    children: 'children',
    label: 'name'
  }

  /**
   * 部门树点击：根据部门筛选用户
   */
  const handleNodeClick = (data: Api.SystemManage.DeptTreeNode) => {
    if (!data || !data.id) {
      console.warn('部门节点数据无效:', data)
      return
    }
    selectedDeptId.value = data.id
    // 设置部门ID到搜索参数并查询用户数据
    const paramsRecord = searchParams as Record<string, any>
    paramsRecord.deptId = data.id
    // 重置到第一页
    paramsRecord.current = 1
    getData()
  }

  /**
   * 加载部门树
   */
  const loadDeptTree = async () => {
    try {
      const res = await fetchGetDeptAll()
      const rows = res?.rows ?? []
      treeData.value = handleTree(
        rows as unknown as Record<string, unknown>[],
        'id',
        'pid',
        'children'
      ) as unknown as Api.SystemManage.DeptTreeNode[]
    } catch (error) {
      console.error('获取部门树失败:', error)
    }
  }

  /**
   * 根据部门ID获取部门名称
   */
  const getDeptName = (deptId?: string): string => {
    if (!deptId) return ''

    const findDept = (
      nodes: Api.SystemManage.DeptTreeNode[]
    ): Api.SystemManage.DeptTreeNode | null => {
      for (const node of nodes) {
        if (node.id === deptId) {
          return node
        }
        if (node.children && node.children.length > 0) {
          const found = findDept(node.children)
          if (found) return found
        }
      }
      return null
    }

    const dept = findDept(treeData.value)
    return dept?.name || ''
  }

  /**
   * 根据部门ID数组获取部门名称数组（用于标签展示）
   */
  const getDeptNames = (deptIds?: string[]): string[] => {
    if (!deptIds || deptIds.length === 0) return []
    return deptIds.map((deptId) => getDeptName(deptId)).filter((name) => name)
  }

  /**
   * 根据角色ID数组获取角色名称
   */
  const getRoleNames = (roleIds?: string[]): string[] => {
    if (!roleIds || roleIds.length === 0 || roleList.value.length === 0) return []
    return roleIds
      .map((roleId) => {
        const role = roleList.value.find((r) => String(r.id) === String(roleId))
        return role?.name || ''
      })
      .filter((name) => name)
  }

  /**
   * 根据岗位ID数组获取岗位名称
   */
  const getPostNames = (postIds?: string[]): string[] => {
    if (!postIds || postIds.length === 0 || postList.value.length === 0) return []
    return postIds
      .map((postId) => {
        const post = postList.value.find((p) => String(p.id) === String(postId))
        return post?.name || ''
      })
      .filter((name) => name)
  }

  /**
   * 格式化时间（统一使用 parseTime 函数）
   */
  const formatTime = (time: string | number | null | undefined): string => {
    return parseTime(time) || '-'
  }

  // 记录每个用户头像是否加载失败
  const avatarErrorMap = reactive<Record<string, boolean>>({})

  /**
   * 头像加载失败处理
   */
  const handleAvatarError = (event: Event, row: UserListItem): void => {
    // 标记该用户的头像加载失败，将显示默认占位符
    if (row.avatar && !avatarErrorMap[row.id]) {
      avatarErrorMap[row.id] = true
    }
  }

  /**
   * 加载角色列表
   */
  const loadRoleList = async () => {
    try {
      const roles = await fetchGetRoleAll()
      roleList.value = roles || []
    } catch (error) {
      console.error('获取角色列表失败:', error)
    }
  }

  /**
   * 加载岗位列表
   */
  const loadPostList = async () => {
    try {
      const res = await fetchGetPostAll()
      postList.value = res?.rows ?? []
    } catch (error) {
      console.error('获取岗位列表失败:', error)
    }
  }

  /**
   * 显示新增用户弹窗
   */
  const showDialog = (type: DialogType): void => {
    dialogType.value = type
    currentUserData.value = undefined
    nextTick(() => {
      dialogVisible.value = true
    })
  }

  /**
   * 显示用户详情弹窗
   */
  const showDetailDialog = (row: UserListItem): void => {
    currentUserData.value = { id: row.id }
    detailDialogVisible.value = true
  }

  /**
   * 处理更多操作点击
   */
  const handleMoreClick = (item: ButtonMoreItem, row: UserListItem) => {
    currentUserData.value = { id: row.id }
    switch (item.key) {
      case 'role':
        roleDialogVisible.value = true
        break
      case 'dept':
        deptDialogVisible.value = true
        break
      case 'post':
        postDialogVisible.value = true
        break
      case 'resetPwd':
        openResetPwdDialog(row)
        break
      case 'delete':
        deleteUser(row)
        break
    }
  }

  /**
   * 打开重置密码弹窗
   */
  const openResetPwdDialog = (row: UserListItem) => {
    resetPwdTargetUser.value = row
    resetPwdForm.password = ''
    resetPwdDialogVisible.value = true
  }

  /**
   * 确认重置用户密码
   */
  const handleResetPwdConfirm = async () => {
    if (!resetPwdTargetUser.value) {
      ElMessage.error(t('sys.user.userNotExist'))
      return
    }
    const pwd = resetPwdForm.password.trim()
    if (!pwd) {
      ElMessage.error(t('sys.user.pleaseInputNewPwd'))
      return
    }
    if (pwd.length < 6 || pwd.length > 100) {
      ElMessage.error(t('sys.user.pwdLengthRule'))
      return
    }

    try {
      resetPwdLoading.value = true
      await fetchResetUserPassword({
        userId: resetPwdTargetUser.value.id,
        password: pwd
      })
      ElMessage.success(t('sys.user.resetPwdSuccess'))
      resetPwdDialogVisible.value = false
    } catch (error) {
      ElMessage.error(error instanceof Error ? error.message : t('sys.user.resetPwdFailed'))
    } finally {
      resetPwdLoading.value = false
    }
  }

  /**
   * 删除用户
   */
  const deleteUser = async (row: UserListItem) => {
    try {
      await ElMessageBox.confirm(
        t('sys.user.deleteUserConfirm', { name: row.nickName || row.userName || '-' }),
        t('common.tips'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
      )
      await fetchDeleteUser([row.id])
      ElMessage.success(t('sys.common.deleteSuccess'))
      refreshData()
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error instanceof Error ? error.message : t('sys.common.deleteFailed'))
      }
    }
  }

  /**
   * 处理弹窗提交事件
   */
  const handleDialogSubmit = async () => {
    refreshData()
  }

  /**
   * 下载 Blob 为文件（校验非 JSON 错误响应）
   */
  const downloadBlob = async (blob: Blob, defaultFilename: string): Promise<void> => {
    if (!blobValidate(blob)) {
      const text = await blob.text()
      try {
        const json = JSON.parse(text) as { msg?: string; message?: string }
        ElMessage.error(json.msg || json.message || t('sys.user.downloadFailed'))
      } catch {
        ElMessage.error(t('sys.user.downloadFailed'))
      }
      return
    }
    FileSaver.saveAs(blob, defaultFilename)
  }

  /**
   * 导出用户列表
   */
  const handleExport = async () => {
    if (exportLoading.value) return
    exportLoading.value = true
    try {
      const params = searchParams as Api.SystemManage.UserSearchParams
      const blob = await fetchExportUser(params)
      await downloadBlob(blob, t('sys.user.userListExcel'))
      ElMessage.success(t('sys.user.exportSuccess'))
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error instanceof Error ? error.message : t('sys.user.exportFailed'))
      }
    } finally {
      exportLoading.value = false
    }
  }

  /**
   * 下载用户导入模板
   */
  const handleDownloadTemplate = async () => {
    if (templateLoading.value) return
    templateLoading.value = true
    try {
      const blob = await fetchGetUserTemplate()
      await downloadBlob(blob, t('sys.user.userTemplateExcel'))
      ElMessage.success(t('sys.user.templateSuccess'))
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error instanceof Error ? error.message : t('sys.user.templateFailed'))
      }
    } finally {
      templateLoading.value = false
    }
  }

  /** 清空导入弹窗内已选文件 */
  const clearImportFile = () => {
    importFileList.value = []
    importFileRef.value = null
  }

  const onImportFileChange = (uploadFile: UploadFile) => {
    const file = uploadFile?.raw
    if (!file) return
    const name = file.name?.toLowerCase() || ''
    if (!name.endsWith('.xlsx') && !name.endsWith('.xls')) {
      ElMessage.warning(t('sys.user.pleaseSelectExcel'))
      return
    }
    importFileRef.value = file
    importFileList.value = [uploadFile]
  }

  const onImportFileRemove = () => {
    importFileRef.value = null
    importFileList.value = []
  }

  /** 导入弹窗内点击确定：上传文件并处理结果 */
  const handleImportConfirm = async () => {
    const file = importFileRef.value
    if (!file) {
      ElMessage.warning(t('sys.user.pleaseSelectFile'))
      return
    }
    importLoading.value = true
    try {
      const res = await fetchImportUser(file, { updateExisting: importUpdateExisting.value })
      const successCount = res?.successCount ?? 0
      const failList = res?.failList ?? []
      if (successCount > 0) {
        ElMessage.success(t('sys.user.importSuccessCount', { count: successCount }))
        refreshData()
        importDialogVisible.value = false
      }
      if (failList.length > 0) {
        importFailList.value = failList
        importFailDialogVisible.value = true
        if (successCount === 0) {
          ElMessage.warning(t('sys.user.importNoSuccess'))
        }
      }
    } catch (error) {
      ElMessage.error(error instanceof Error ? error.message : t('sys.user.importFailed'))
    } finally {
      importLoading.value = false
    }
  }

  // 搜索表单
  const searchForm = ref({
    userName: undefined as string | undefined,
    phone: undefined as string | undefined,
    email: undefined as string | undefined,
    sex: undefined as number | undefined,
    status: undefined as number | undefined,
    daterange: undefined as [string, string] | undefined
  })

  const {
    data,
    columns,
    columnChecks,
    loading,
    pagination,
    refreshData,
    handleSizeChange,
    handleCurrentChange,
    searchParams,
    getData
  } = useTable({
    core: {
      apiFn: fetchGetUserList,
      apiParams: {
        current: 1,
        size: 20,
        userName: '',
        phone: '',
        email: ''
      },
      // 排除搜索表单中的日期区间字段
      excludeParams: ['daterange'],
      columnsFactory: () => [
        {
          type: 'globalIndex',
          label: t('sys.user.index'),
          width: 80,
          align: 'center'
        },
        {
          prop: 'userInfo',
          label: t('sys.user.userInfo'),
          align: 'left',
          minWidth: 160,
          useSlot: true
        },
        {
          prop: 'contact',
          label: t('sys.user.contact'),
          align: 'left',
          minWidth: 140,
          useSlot: true
        },
        {
          prop: 'deptId',
          label: t('sys.user.dept'),
          align: 'left',
          minWidth: 120,
          useSlot: true
        },
        {
          prop: 'roles',
          label: t('sys.user.role'),
          align: 'left',
          minWidth: 120,
          useSlot: true
        },
        {
          prop: 'posts',
          label: t('sys.user.post'),
          align: 'left',
          minWidth: 120,
          useSlot: true
        },
        {
          prop: 'status',
          label: t('sys.user.status'),
          width: 80,
          align: 'center',
          formatter: (row: UserListItem) => {
            const status = Number(row.status)
            const statusConfig =
              status === 1
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
          prop: 'timeInfo',
          label: t('sys.user.timeInfo'),
          minWidth: 160,
          align: 'left',
          useSlot: true
        },
        {
          prop: 'operation',
          label: t('sys.user.action'),
          width: 150,
          align: 'right',
          headerAlign: 'center',
          fixed: 'right',
          useSlot: true
        }
      ]
    }
  })

  /**
   * 搜索
   */
  const handleSearch = (params: Record<string, any>) => {
    const paramsRecord = searchParams as Record<string, any>
    // 处理日期区间参数，把 daterange 转换为 startTime 和 endTime
    const { daterange, ...filtersParams } = params
    const [startTime, endTime] = Array.isArray(daterange) ? daterange : [null, null]

    Object.assign(paramsRecord, { ...filtersParams, startTime, endTime })
    // 如果搜索时没有指定部门，保持当前选中的部门
    if (!filtersParams.deptId && selectedDeptId.value) {
      paramsRecord.deptId = selectedDeptId.value
    }
    getData()
  }

  /**
   * 重置搜索
   */
  const handleReset = () => {
    Object.assign(searchForm.value, {
      userName: undefined,
      phone: undefined,
      email: undefined,
      sex: undefined,
      status: undefined,
      daterange: undefined
    })
    // 清除选中的部门
    selectedDeptId.value = undefined
    const paramsRecord = searchParams as Record<string, any>
    Object.keys(paramsRecord).forEach((key) => {
      // 保留分页参数，清空过滤条件
      if (!['current', 'size'].includes(key)) {
        paramsRecord[key] = undefined
      }
    })
    selectedDeptId.value = undefined
    getData()
  }

  onMounted(async () => {
    await loadDeptTree()
    await loadRoleList()
    await loadPostList()
  })
</script>

<style scoped>
  .tree-card :deep(.el-card__body) {
    flex: 1;
    min-height: 0;
    padding: 10px 2px 10px 10px;
  }

  .import-upload-inner {
    padding: 24px 0;
  }

  .import-upload-icon {
    font-size: 48px;
    color: var(--el-text-color-placeholder);
  }

  .import-upload-text {
    margin-top: 8px;
    font-size: 14px;
    color: var(--el-text-color-secondary);
  }

  .import-upload-text em {
    font-style: normal;
    color: var(--el-color-primary);
  }

  .import-options {
    margin-top: 16px;
  }

  .import-options .el-checkbox {
    display: block;
    margin-bottom: 8px;
  }

  .import-format-hint {
    margin: 0 0 8px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }
</style>

<style>
  /* 确保头像列的表头和内容都居中对齐 */
  .art-table :deep(.el-table__header-wrapper .el-table__header th[data-column-key='avatar']) {
    text-align: center !important;
  }

  .art-table :deep(.el-table__body-wrapper .el-table__body td[data-column-key='avatar'] .cell) {
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
    padding: 0 !important;
  }
</style>
