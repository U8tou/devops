<template>
  <ElDrawer
    v-model="visible"
    :title="form.id ? $t('sys.user.titleWithId', { id: form.id }) : $t('sys.user.title')"
    size="640px"
    direction="rtl"
    :destroy-on-close="false"
    @close="handleClose"
  >
    <div class="flex flex-col h-full" v-loading="loading">
      <div class="flex-1 overflow-auto pr-1">
        <ElForm ref="formRef" :model="form" :rules="rules" label-position="top">
          <ElFormItem :label="$t('sys.user.username')" prop="userName">
            <ElInput v-model="form.userName" :placeholder="$t('sys.user.usernamePlaceholder')" />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.nickname')" prop="nickName">
            <ElInput v-model="form.nickName" :placeholder="$t('sys.user.nicknamePlaceholder')" />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.phone')" prop="phone">
            <ElInput v-model="form.phone" :placeholder="$t('sys.user.phonePlaceholder')" />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.email')" prop="email">
            <ElInput v-model="form.email" :placeholder="$t('sys.user.emailPlaceholder')" />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.sex')" prop="sex">
            <ElRadioGroup v-model="form.sex">
              <ElRadio :label="1">{{ $t('sys.common.sexMale') }}</ElRadio>
              <ElRadio :label="2">{{ $t('sys.common.sexFemale') }}</ElRadio>
              <ElRadio :label="3">{{ $t('sys.common.sexUnknown') }}</ElRadio>
            </ElRadioGroup>
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.status')" prop="status">
            <ElRadioGroup v-model="form.status">
              <ElRadio :label="1">{{ $t('sys.common.statusNormal') }}</ElRadio>
              <ElRadio :label="2">{{ $t('sys.common.statusDisabled') }}</ElRadio>
            </ElRadioGroup>
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.address')" prop="address">
            <ElInput
              v-model="form.address"
              type="textarea"
              :rows="2"
              :placeholder="$t('sys.user.addressPlaceholder')"
            />
          </ElFormItem>
          <ElFormItem :label="$t('sys.user.remark')" prop="remark">
            <ElInput
              v-model="form.remark"
              type="textarea"
              :rows="3"
              :placeholder="$t('sys.user.remarkPlaceholder')"
            />
          </ElFormItem>
        </ElForm>
      </div>
      <div class="mt-4 flex justify-end gap-3 pt-3">
        <ElButton @click="handleClose" :disabled="submitLoading">{{
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
  import { fetchGetUserDetail, fetchEditUser } from '@/api/system-manage'

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

  const formRef = ref<FormInstance>()
  const loading = ref(false)
  const submitLoading = ref(false)

  /**
   * 弹窗显示状态双向绑定
   */
  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  const { t } = useI18n()
  /**
   * 表单验证规则
   */
  const rules = computed<FormRules>(() => ({
    userName: [
      { required: true, message: t('sys.user.usernamePlaceholder'), trigger: 'blur' },
      { min: 2, max: 50, message: t('sys.user.usernameRule'), trigger: 'blur' }
    ],
    nickName: [
      { required: true, message: t('sys.user.nicknamePlaceholder'), trigger: 'blur' },
      { min: 2, max: 50, message: t('sys.user.usernameRule'), trigger: 'blur' }
    ],
    phone: [
      { required: true, message: t('sys.user.phonePlaceholder'), trigger: 'blur' },
      { pattern: /^1[3-9]\d{9}$/, message: t('sys.user.phoneFormatRule'), trigger: 'blur' }
    ],
    email: [{ type: 'email', message: t('sys.user.emailPlaceholder'), trigger: 'blur' }],
    sex: [{ required: true, message: t('sys.user.sex'), trigger: 'change' }],
    status: [{ required: true, message: t('sys.user.statusPlaceholder'), trigger: 'change' }]
  }))

  /**
   * 表单数据
   */
  const form = reactive<Partial<Api.SystemManage.UserListItem>>({
    id: '',
    userName: '',
    nickName: '',
    phone: '',
    phoneArea: '',
    email: '',
    sex: 3,
    status: 1,
    address: '',
    remark: ''
  })

  /**
   * 用户详情数据（用于保存登录相关字段）
   */
  const userDetail = ref<Api.SystemManage.UserDetail | null>(null)

  /**
   * 重置表单
   */
  const resetForm = (): void => {
    formRef.value?.resetFields()
    Object.assign(form, {
      id: '',
      userName: '',
      nickName: '',
      phone: '',
      phoneArea: '',
      email: '',
      sex: 3,
      status: 1,
      address: '',
      remark: ''
    })
  }

  /**
   * 加载用户详情
   */
  const loadUserData = async (): Promise<void> => {
    if (!props.userData?.id) return

    try {
      loading.value = true
      const detail = await fetchGetUserDetail(props.userData.id)
      userDetail.value = detail
      Object.assign(form, {
        id: detail.id || '',
        userName: detail.userName || '',
        nickName: detail.nickName || '',
        phone: detail.phone || '',
        phoneArea: detail.phoneArea || '',
        email: detail.email || '',
        sex: detail.sex ?? 3,
        status: detail.status ?? 1,
        address: detail.address || '',
        remark: detail.remark || ''
      })
    } catch (error) {
      console.error('获取用户详情失败:', error)
      ElMessage.error(t('sys.common.getDetailFailed'))
    } finally {
      loading.value = false
    }
  }

  /**
   * 监听弹窗打开，初始化表单数据
   */
  watch(
    () => props.modelValue,
    async (newVal) => {
      if (newVal && props.userData?.id) {
        await loadUserData()
      } else {
        resetForm()
      }
    }
  )

  /**
   * 关闭弹窗并重置表单
   */
  const handleClose = () => {
    visible.value = false
    resetForm()
  }

  /**
   * 提交表单
   */
  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
      submitLoading.value = true

      if (!form.id) {
        ElMessage.error(t('sys.user.userIdNotExist'))
        return
      }

      // 构建编辑请求参数（严格对应 internal_sys_user.EditReq）
      const payload: Api.SystemManage.UserEditReq = {
        id: form.id!,
        userName: form.userName || '',
        nickName: form.nickName || '',
        phone: form.phone || '',
        phoneArea: form.phoneArea || '+86',
        email: form.email || '',
        sex: form.sex ?? 3,
        status: form.status ?? 1,
        address: form.address || '',
        remark: form.remark || '',
        userType: userDetail.value?.userType || '00'
      }

      await fetchEditUser(payload)
      ElMessage.success(t('sys.common.editSuccess'))
      emit('success')
      handleClose()
    } catch (error) {
      if (error !== false) {
        ElMessage.error(error instanceof Error ? error.message : t('sys.common.operateFailed'))
      }
    } finally {
      submitLoading.value = false
    }
  }
</script>
