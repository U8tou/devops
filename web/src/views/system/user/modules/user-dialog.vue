<template>
  <ElDrawer
    v-model="dialogVisible"
    :title="$t('sys.user.addUser')"
    size="640px"
    direction="rtl"
    :destroy-on-close="false"
    @close="initFormData"
  >
    <div class="flex flex-col h-full">
      <div class="flex-1 overflow-auto pr-1">
        <ElForm ref="formRef" :model="formData" :rules="rules" label-position="top">
          <ElFormItem :label="$t('sys.user.username')" prop="userName">
            <ElInput
              v-model="formData.userName"
              :placeholder="$t('sys.user.usernamePlaceholder')"
            />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.nickname')" prop="nickName">
            <ElInput
              v-model="formData.nickName"
              :placeholder="$t('sys.user.nicknamePlaceholder')"
            />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.phone')" prop="phone">
            <ElInput v-model="formData.phone" :placeholder="$t('sys.user.phonePlaceholder')" />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.email')" prop="email">
            <ElInput v-model="formData.email" :placeholder="$t('sys.user.emailPlaceholder')" />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.password')" prop="password">
            <ElInput
              v-model="formData.password"
              type="password"
              :placeholder="$t('sys.user.passwordPlaceholder')"
              show-password
            />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.sex')" prop="sex">
            <ElRadioGroup v-model="formData.sex">
              <ElRadio :label="1">{{ $t('sys.common.sexMale') }}</ElRadio>
              <ElRadio :label="2">{{ $t('sys.common.sexFemale') }}</ElRadio>
              <ElRadio :label="3">{{ $t('sys.common.sexUnknown') }}</ElRadio>
            </ElRadioGroup>
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.address')" prop="address">
            <ElInput
              v-model="formData.address"
              type="textarea"
              :rows="2"
              :placeholder="$t('sys.user.addressPlaceholder')"
            />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.remark')" prop="remark">
            <ElInput
              v-model="formData.remark"
              type="textarea"
              :rows="3"
              :placeholder="$t('sys.user.remarkPlaceholder')"
            />
          </ElFormItem>
        </ElForm>
      </div>
      <div class="mt-4 flex justify-end gap-3 pt-3">
        <ElButton @click="dialogVisible = false" :disabled="submitLoading">{{
          $t('common.cancel')
        }}</ElButton>
        <ElButton type="primary" @click="handleSubmit" :loading="submitLoading">{{
          $t('common.confirm')
        }}</ElButton>
      </div>
    </div>
  </ElDrawer>
</template>

<script setup lang="ts">
  import type { FormInstance, FormRules } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import { ElMessage } from 'element-plus'
  import { fetchAddUser, fetchGetRoleAll } from '@/api/system-manage'

  interface Props {
    visible: boolean
    type: string
    userData?: Partial<Api.SystemManage.UserListItem>
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'submit'): void
  }

  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 角色列表数据
  const roleList = ref<Api.SystemManage.RoleListItem[]>([])
  const roleListLoading = ref(false)

  // 对话框显示控制
  const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
  })

  // 表单实例
  const formRef = ref<FormInstance>()

  // 表单数据
  const formData = reactive({
    userName: '',
    nickName: '',
    phone: '',
    phoneArea: '+86',
    email: '',
    sex: 1,
    password: '',
    address: '',
    remark: '',
    avatar: '',
    userType: '00'
  })

  const submitLoading = ref(false)

  const { t } = useI18n()
  // 表单验证规则
  const rules = computed<FormRules>(() => ({
    userName: [
      { required: true, message: t('sys.user.usernamePlaceholder'), trigger: 'blur' },
      { min: 2, max: 30, message: t('sys.user.usernameRule'), trigger: 'blur' }
    ],
    nickName: [
      { required: true, message: t('sys.user.nicknamePlaceholder'), trigger: 'blur' },
      { min: 2, max: 30, message: t('sys.user.usernameRule'), trigger: 'blur' }
    ],
    phone: [
      { required: true, message: t('sys.user.phonePlaceholder'), trigger: 'blur' },
      { pattern: /^1[3-9]\d{9}$/, message: t('sys.user.phoneFormatRule'), trigger: 'blur' }
    ],
    email: [{ type: 'email', message: t('sys.user.emailPlaceholder'), trigger: 'blur' }],
    password: [
      { required: true, message: t('sys.user.passwordPlaceholder'), trigger: 'blur' },
      { min: 6, max: 100, message: t('sys.user.pwdLengthRule'), trigger: 'blur' }
    ],
    sex: [{ required: true, message: t('sys.user.sex'), trigger: 'change' }]
  }))

  /**
   * 加载角色列表
   */
  const loadRoleList = async () => {
    if (roleList.value.length > 0) return
    try {
      roleListLoading.value = true
      const roles = await fetchGetRoleAll()
      roleList.value = roles || []
    } catch (error) {
      console.error('获取角色列表失败:', error)
    } finally {
      roleListLoading.value = false
    }
  }

  /**
   * 初始化表单数据
   */
  const initFormData = () => {
    Object.assign(formData, {
      userName: '',
      nickName: '',
      phone: '',
      phoneArea: '+86',
      email: '',
      sex: 1,
      password: '',
      address: '',
      remark: '',
      avatar: '',
      userType: '00'
    })
  }

  /**
   * 监听对话框状态变化
   */
  watch(
    () => props.visible,
    async (visible) => {
      if (visible && props.type === 'add') {
        initFormData()
        await loadRoleList()
        nextTick(() => {
          formRef.value?.clearValidate()
        })
      }
    }
  )

  /**
   * 提交表单
   */
  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
      submitLoading.value = true

      // 构建新增请求参数（严格对应 internal_sys_user.AddReq）
      const payload: Api.SystemManage.UserAddReq = {
        userName: formData.userName,
        password: formData.password,
        nickName: formData.nickName,
        phone: formData.phone,
        phoneArea: formData.phoneArea,
        email: formData.email || '',
        sex: formData.sex,
        address: formData.address || '',
        remark: formData.remark || '',
        status: 1,
        userType: formData.userType,
        avatar: formData.avatar || ''
      }

      await fetchAddUser(payload)
      ElMessage.success(t('sys.common.addSuccess'))
      emit('submit')
      dialogVisible.value = false
    } catch (error) {
      if (error !== false) {
        ElMessage.error(error instanceof Error ? error.message : t('sys.common.operateFailed'))
      }
    } finally {
      submitLoading.value = false
    }
  }
</script>
