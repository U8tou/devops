<!-- 个人中心页面 -->
<template>
  <div class="w-full h-full p-0 bg-transparent border-none shadow-none">
    <div class="relative flex-b mt-2.5 max-md:block max-md:mt-1">
      <div class="w-112 mr-5 max-md:w-full max-md:mr-0">
        <div class="art-card-sm relative p-9 pb-6 overflow-hidden text-center">
          <img class="absolute top-0 left-0 w-full h-50 object-cover" src="@imgs/user/bg.webp" />
          <img
            class="relative z-10 w-20 h-20 mt-30 mx-auto object-cover border-2 border-white rounded-full c-p"
            :src="avatarUrl"
            @click="openAvatarDialog"
          />
          <h2 class="mt-5 text-xl font-normal">
            {{ userInfo.nickName || '未设置昵称' }}
          </h2>
          <p class="mt-2 mb-1 text-sm text-g-500">
            {{ userInfo.userName || `${date}，欢迎回来` }}
          </p>

          <div class="w-75 mx-auto mt-7.5 text-left">
            <div class="mt-2.5">
              <ArtSvgIcon icon="ri:mail-line" class="text-g-700" />
              <span class="ml-2 text-sm">{{ userInfo.email || '未设置邮箱' }}</span>
            </div>

            <div class="mt-2.5">
              <ArtSvgIcon icon="ri:user-3-line" class="text-g-700" />
              <span class="ml-2 text-sm">
                {{ userInfo.userType === '00' ? '系统用户' : '普通用户' }}
              </span>
            </div>
            <div class="mt-2.5">
              <ArtSvgIcon icon="ri:map-pin-line" class="text-g-700" />
              <span class="ml-2 text-sm">{{ userInfo.address || '未设置地址' }}</span>
            </div>
            <div class="mt-2.5">
              <ArtSvgIcon icon="ri:user-star-line" class="text-g-700" />
              <span class="ml-2 text-sm">
                {{ roleTagList.length ? roleTagList.join('，') : '暂无角色信息' }}
              </span>
            </div>
            <div class="mt-2.5">
              <ArtSvgIcon icon="ri:organization-chart" class="text-g-700" />
              <span class="ml-2 text-sm">
                {{ deptTagList.length ? deptTagList.join('，') : '暂无部门信息' }}
              </span>
            </div>
            <div class="mt-2.5">
              <ArtSvgIcon icon="ri:briefcase-line" class="text-g-700" />
              <span class="ml-2 text-sm">
                {{ postTagList.length ? postTagList.join('，') : '暂无岗位信息' }}
              </span>
            </div>
            <div class="mt-2.5">
              <ArtSvgIcon icon="ri:calendar-line" class="text-g-700" />
              <span class="ml-2 text-sm">
                {{
                  userInfo.createTime
                    ? parseTime(userInfo.createTime, '{y}-{m}-{d} {h}:{i}')
                    : '未知'
                }}
              </span>
            </div>
          </div>

          <div class="mt-10">
            <h3 class="text-sm font-medium">标签</h3>
            <div class="flex flex-wrap justify-center mt-3.5">
              <div
                v-for="item in lableList"
                :key="item"
                class="py-1 px-1.5 mr-2.5 mb-2.5 text-xs border border-g-300 rounded"
              >
                {{ item }}
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="flex-1 overflow-hidden max-md:w-full max-md:mt-3.5">
        <div class="art-card-sm">
          <h1 class="p-4 text-xl font-normal border-b border-g-300">基本设置</h1>

          <ElForm
            :model="form"
            class="box-border p-5 [&>.el-row_.el-form-item]:w-[calc(50%-10px)] [&>.el-row_.el-input]:w-full [&>.el-row_.el-select]:w-full"
            ref="ruleFormRef"
            :rules="rules"
            label-width="86px"
            label-position="top"
          >
            <ElRow>
              <ElFormItem label="用户名">
                <ElInput :model-value="userInfo.userName" disabled />
              </ElFormItem>
              <ElFormItem label="性别" prop="sex" class="ml-5">
                <ElSelect v-model="form.sex" placeholder="Select" :disabled="!isEdit">
                  <ElOption
                    v-for="item in options"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                  />
                </ElSelect>
              </ElFormItem>
            </ElRow>

            <ElRow>
              <ElFormItem label="昵称" prop="nikeName">
                <ElInput v-model="form.nikeName" :disabled="!isEdit" />
              </ElFormItem>
              <ElFormItem label="邮箱" prop="email" class="ml-5">
                <ElInput v-model="form.email" :disabled="!isEdit" />
              </ElFormItem>
            </ElRow>

            <ElRow>
              <ElFormItem label="手机" prop="mobile">
                <ElInput v-model="form.mobile" :disabled="!isEdit" />
              </ElFormItem>
              <ElFormItem label="地址" prop="address" class="ml-5">
                <ElInput v-model="form.address" :disabled="!isEdit" />
              </ElFormItem>
            </ElRow>

            <!-- 去掉个人介绍字段 -->

            <div class="flex-c justify-end [&_.el-button]:!w-27.5">
              <ElButton type="primary" class="w-22.5" v-ripple @click="edit">
                {{ isEdit ? '保存' : '编辑' }}
              </ElButton>
            </div>
          </ElForm>
        </div>

        <div class="art-card-sm my-5">
          <h1 class="p-4 text-xl font-normal border-b border-g-300">更改密码</h1>

          <ElForm
            :model="pwdForm"
            class="box-border p-5"
            label-width="86px"
            label-position="top"
            ref="pwdFormRef"
            :rules="pwdRules"
          >
            <ElFormItem label="当前密码" prop="password">
              <ElInput
                v-model="pwdForm.password"
                type="password"
                :disabled="!isEditPwd"
                show-password
              />
            </ElFormItem>

            <ElFormItem label="新密码" prop="newPassword">
              <ElInput
                v-model="pwdForm.newPassword"
                type="password"
                :disabled="!isEditPwd"
                show-password
              />
            </ElFormItem>

            <ElFormItem label="确认新密码" prop="confirmPassword">
              <ElInput
                v-model="pwdForm.confirmPassword"
                type="password"
                :disabled="!isEditPwd"
                show-password
              />
            </ElFormItem>

            <div class="flex-c justify-end [&_.el-button]:!w-27.5">
              <ElButton type="primary" class="w-22.5" v-ripple @click="editPwd">
                {{ isEditPwd ? '保存' : '编辑' }}
              </ElButton>
            </div>
          </ElForm>
        </div>
      </div>
    </div>

    <!-- 更换头像弹窗 -->
    <ElDialog
      v-model="avatarDialogVisible"
      title="更换头像"
      width="560px"
      align-center
      destroy-on-close
    >
      <ArtCutterImg
        :is-modal="false"
        :show-preview="false"
        :box-width="520"
        :box-height="380"
        :cut-width="320"
        :cut-height="320"
        :file-type="'png'"
        :quality="0.9"
        :img-url="avatarPreviewUrl"
        @update:imgUrl="handleAvatarPreviewChange"
      />
      <template #footer>
        <ElButton @click="avatarDialogVisible = false" :disabled="avatarUploading">取消</ElButton>
        <ElButton
          type="primary"
          @click="submitAvatar"
          :loading="avatarUploading"
          :disabled="!avatarCroppedUrl || avatarUploading"
        >
          上传头像
        </ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { useUserStore } from '@/store/modules/user'
  import {
    fetchGetUserInfo,
    fetchChangePassword,
    fetchUpdateProfile,
    fetchUploadAvatar
  } from '@/api/auth'
  import { fetchGetRoleAll, fetchGetDeptAll, fetchGetPostAll } from '@/api/system-manage'
  import type { FormInstance, FormRules } from 'element-plus'
  import { ElMessage } from 'element-plus'
  import defaultAvatar from '@imgs/user/avatar.webp'
  import { resolveImageUrl } from '@/utils/resolveImageUrl'
  import { parseTime } from '@/utils/ruoyi'
  import ArtCutterImg from '@/components/core/media/art-cutter-img/index.vue'

  defineOptions({ name: 'UserCenter' })

  const userStore = useUserStore()
  const userInfo = computed(() => userStore.getUserInfo)

  const isEdit = ref(false)
  const isEditPwd = ref(false)
  const date = ref('')
  const ruleFormRef = ref<FormInstance>()
  const pwdFormRef = ref<FormInstance>()

  const avatarDialogVisible = ref(false)
  const avatarPreviewUrl = ref('')
  const avatarCroppedUrl = ref('')
  const avatarUploading = ref(false)

  const roleAllList = ref<Api.SystemManage.RoleListItem[]>([])
  const deptAllList = ref<Api.SystemManage.DeptListItem[]>([])
  const postAllList = ref<Api.SystemManage.PostListItem[]>([])

  /**
   * 头像地址（处理相对路径 & 默认图）
   */
  const avatarUrl = computed(() => {
    const avatar = userInfo.value.avatar
    if (!avatar) return defaultAvatar
    return resolveImageUrl(avatar) || defaultAvatar
  })

  /**
   * 用户信息表单
   */
  const form = reactive({
    realName: '',
    nikeName: '',
    email: '',
    mobile: '',
    address: '',
    sex: '3',
    des: '',
    phoneArea: ''
  })

  /**
   * 密码修改表单
   */
  const pwdForm = reactive({
    password: '',
    newPassword: '',
    confirmPassword: ''
  })

  /**
   * 表单验证规则
   */
  const rules = reactive<FormRules>({
    realName: [
      { required: true, message: '请输入姓名', trigger: 'blur' },
      { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
    ],
    nikeName: [
      { required: true, message: '请输入昵称', trigger: 'blur' },
      { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
    ],
    email: [{ required: true, message: '请输入邮箱', trigger: 'blur' }],
    mobile: [{ required: true, message: '请输入手机号码', trigger: 'blur' }],
    address: [{ required: true, message: '请输入地址', trigger: 'blur' }],
    sex: [{ required: true, message: '请选择性别', trigger: 'blur' }]
  })

  /**
   * 密码表单验证规则
   */
  const pwdRules = reactive<FormRules>({
    password: [
      { required: true, message: '请输入当前密码', trigger: 'blur' },
      { min: 6, max: 50, message: '长度在 6 到 50 个字符', trigger: 'blur' }
    ],
    newPassword: [
      { required: true, message: '请输入新密码', trigger: 'blur' },
      { min: 6, max: 50, message: '长度在 6 到 50 个字符', trigger: 'blur' }
    ],
    confirmPassword: [
      { required: true, message: '请再次输入新密码', trigger: 'blur' },
      {
        validator: (_rule, value, callback) => {
          if (!value) {
            callback(new Error('请再次输入新密码'))
          } else if (value !== pwdForm.newPassword) {
            callback(new Error('两次输入的新密码不一致'))
          } else {
            callback()
          }
        },
        trigger: 'blur'
      }
    ]
  })

  /**
   * 性别选项
   */
  const options = [
    { value: '1', label: '男' },
    { value: '2', label: '女' }
  ]

  /**
   * 用户标签列表
   */
  const lableList = computed<Array<string>>(() => {
    const list: string[] = []
    const info = userInfo.value

    if (info.userType === '00') {
      list.push('系统用户')
    } else {
      list.push('普通用户')
    }

    if (info.status === 1) {
      list.push('账号正常')
    } else if (info.status === 2) {
      list.push('账号停用')
    }

    if (info.roles && info.roles.length && roleAllList.value.length) {
      const roleMap = new Map(roleAllList.value.map((item) => [item.id, item.name]))
      list.push(...info.roles.map((id) => roleMap.get(id) || id).map((name) => `${name}`))
    }

    if (info.posts && info.posts.length && postAllList.value.length) {
      const postMap = new Map(postAllList.value.map((item) => [item.id, item.name]))
      list.push(...info.posts.map((id) => postMap.get(id) || id).map((name) => `${name}`))
    }

    return list.length ? list : ['普通用户']
  })

  /**
   * 角色标签列表（根据当前用户角色ID映射成名称）
   */
  const roleTagList = computed<Array<string>>(() => {
    if (!userInfo.value.roles || !userInfo.value.roles.length || !roleAllList.value.length) {
      return []
    }
    const map = new Map(roleAllList.value.map((item) => [item.id, item.name]))
    return userInfo.value.roles
      .map((id) => map.get(id))
      .filter((name): name is string => Boolean(name))
  })

  /**
   * 部门标签列表（根据当前用户部门ID映射成名称）
   */
  const deptTagList = computed<Array<string>>(() => {
    if (!userInfo.value.depts || !userInfo.value.depts.length || !deptAllList.value.length) {
      return []
    }
    const map = new Map(deptAllList.value.map((item) => [item.id, item.name]))
    return userInfo.value.depts
      .map((id) => map.get(id))
      .filter((name): name is string => Boolean(name))
  })

  /**
   * 岗位标签列表（根据当前用户岗位ID映射成名称）
   */
  const postTagList = computed<Array<string>>(() => {
    if (!userInfo.value.posts || !userInfo.value.posts.length || !postAllList.value.length) {
      return []
    }
    const map = new Map(postAllList.value.map((item) => [item.id, item.name]))
    return userInfo.value.posts
      .map((id) => map.get(id))
      .filter((name): name is string => Boolean(name))
  })

  /**
   * 打开头像裁剪弹窗
   */
  const openAvatarDialog = () => {
    avatarPreviewUrl.value = avatarUrl.value
    avatarCroppedUrl.value = ''
    avatarDialogVisible.value = true
  }

  /**
   * dataURL 转 File
   */
  const dataURLToFile = (dataUrl: string, fileName: string): File => {
    const arr = dataUrl.split(',')
    const mimeMatch = arr[0].match(/:(.*?);/)
    const mime = mimeMatch ? mimeMatch[1] : 'image/png'
    const bstr = atob(arr[1])
    let n = bstr.length
    const u8arr = new Uint8Array(n)

    while (n--) {
      u8arr[n] = bstr.charCodeAt(n)
    }

    return new File([u8arr], fileName, { type: mime })
  }

  /**
   * 头像裁剪完成回调：上传并更新用户信息
   */
  const handleAvatarPreviewChange = (dataUrl: string) => {
    // 裁剪结果用于后续上传
    avatarCroppedUrl.value = dataUrl
    // 同步给预览地址，保持裁剪框中仍然显示最新图片
    if (dataUrl) {
      avatarPreviewUrl.value = dataUrl
    }
  }

  /**
   * 提交头像上传
   */
  const submitAvatar = async () => {
    if (!avatarCroppedUrl.value) return

    try {
      avatarUploading.value = true
      const file = dataURLToFile(avatarCroppedUrl.value, 'avatar.png')
      const res = await fetchUploadAvatar(file)

      const newAvatar = res.avatar

      // 更新本地用户信息
      userStore.setUserInfo({
        ...(userInfo.value as Api.Auth.UserInfo),
        avatar: newAvatar
      })

      ElMessage.success('头像更新成功')
      avatarDialogVisible.value = false
    } catch (error) {
      console.error(error)
      ElMessage.error('头像上传失败，请稍后重试')
    } finally {
      avatarUploading.value = false
    }
  }

  /**
   * 使用用户信息初始化个人资料表单
   */
  const initFormFromUserInfo = () => {
    const info = userInfo.value
    if (!info) return

    form.realName = info.nickName || info.userName || ''
    form.nikeName = info.nickName || ''
    form.email = info.email || ''
    form.mobile = info.phone || ''
    form.address = info.address || ''
    form.sex = info.sex ? String(info.sex) : '3'
    form.des = info.remark || ''
    form.phoneArea = info.phoneArea || ''
  }

  /**
   * 加载最新的用户信息
   */
  const loadUserInfo = async () => {
    const data = await fetchGetUserInfo()
    userStore.setUserInfo(data)
    initFormFromUserInfo()
  }

  /**
   * 加载全部角色、部门、岗位数据，用于将ID映射为名称
   */
  const loadRoleDeptAndPost = async () => {
    try {
      const [roleRes, deptRes, postRes] = await Promise.all([
        fetchGetRoleAll(),
        fetchGetDeptAll(),
        fetchGetPostAll()
      ])
      roleAllList.value = roleRes || []
      deptAllList.value = deptRes.rows || []
      postAllList.value = postRes?.rows ?? []
    } catch (error) {
      console.error('加载角色/部门/岗位信息失败', error)
    }
  }

  onMounted(async () => {
    getDate()
    // 初始化表单
    initFormFromUserInfo()
    // 再从接口拉取最新信息
    try {
      await loadUserInfo()
    } catch {
      // 获取失败时保持本地信息
    }
    // 加载角色、部门、岗位映射数据
    await loadRoleDeptAndPost()
  })

  /**
   * 根据当前时间获取问候语
   */
  const getDate = () => {
    const h = new Date().getHours()

    if (h >= 6 && h < 9) date.value = '早上好'
    else if (h >= 9 && h < 11) date.value = '上午好'
    else if (h >= 11 && h < 13) date.value = '中午好'
    else if (h >= 13 && h < 18) date.value = '下午好'
    else if (h >= 18 && h < 24) date.value = '晚上好'
    else date.value = '很晚了，早点睡'
  }

  /**
   * 切换用户信息编辑状态
   */
  const edit = async () => {
    // 当前为非编辑状态，点击后进入编辑模式
    if (!isEdit.value) {
      initFormFromUserInfo()
      isEdit.value = true
      return
    }

    // 保存
    if (!ruleFormRef.value) return
    await ruleFormRef.value.validate(async (valid) => {
      if (!valid) return
      const payload: Api.Auth.UpdateProfileReq = {
        nickName: form.nikeName,
        email: form.email,
        phone: form.mobile,
        phoneArea: form.phoneArea || userInfo.value.phoneArea,
        address: form.address,
        sex: Number(form.sex || 3)
      }

      try {
        await fetchUpdateProfile(payload)
        await loadUserInfo()
        isEdit.value = false
        ElMessage.success('个人信息已保存')
      } catch (error) {
        console.error(error)
        ElMessage.error('保存个人信息失败，请稍后重试')
      }
    })
  }

  /**
   * 切换密码编辑状态
   */
  const editPwd = async () => {
    // 当前为非编辑状态，点击后进入编辑模式
    if (!isEditPwd.value) {
      pwdForm.password = ''
      pwdForm.newPassword = ''
      pwdForm.confirmPassword = ''
      isEditPwd.value = true
      return
    }

    // 保存密码
    if (!pwdFormRef.value) return
    await pwdFormRef.value.validate(async (valid) => {
      if (!valid) return

      try {
        await fetchChangePassword({
          oldPassword: pwdForm.password,
          newPassword: pwdForm.newPassword
        })
        ElMessage.success('密码修改成功')
        isEditPwd.value = false
        pwdForm.password = ''
        pwdForm.newPassword = ''
        pwdForm.confirmPassword = ''
      } catch (error) {
        console.error(error)
        ElMessage.error('密码修改失败，请稍后重试')
      }
    })
  }
</script>
