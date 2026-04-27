<template>
  <ElDialog
    v-model="visible"
    :title="
      dialogType === 'add' ? $t('sys.role.addRole') : $t('sys.role.editRole', { id: form.id })
    "
    width="640px"
    align-center
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-position="top">
      <ElFormItem :label="$t('sys.role.roleName')" prop="name">
        <ElInput v-model="form.name" :placeholder="$t('sys.role.roleNamePlaceholder')" />
      </ElFormItem>
      <ElFormItem :label="$t('sys.role.roleKey')" prop="role">
        <ElInput v-model="form.role" :placeholder="$t('sys.role.roleKeyPlaceholder')" />
      </ElFormItem>
      <ElFormItem :label="$t('sys.user.status')" prop="status">
        <ElRadioGroup v-model="form.status">
          <ElRadio :label="1">{{ $t('sys.common.statusNormal') }}</ElRadio>
          <ElRadio :label="2">{{ $t('sys.common.statusDisabled') }}</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem :label="$t('sys.role.sortLabel')" prop="sort">
        <ElInputNumber v-model="form.sort" :min="0" controls-position="right" style="width: 100%" />
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
    <template #footer>
      <ElButton @click="handleClose" :disabled="submitLoading">{{ $t('common.cancel') }}</ElButton>
      <ElButton type="primary" @click="handleSubmit" :loading="submitLoading">{{
        $t('common.confirm')
      }}</ElButton>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import type { FormInstance, FormRules } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import { ElMessage } from 'element-plus'
  import { fetchGetRoleInfo, fetchAddRole, fetchEditRole } from '@/api/system-manage'

  interface Props {
    modelValue: boolean
    dialogType: 'add' | 'edit'
    roleData?: { id: string }
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    dialogType: 'add',
    roleData: undefined
  })

  const emit = defineEmits<Emits>()
  const { t } = useI18n()

  const formRef = ref<FormInstance>()
  const submitLoading = ref(false)

  /**
   * 弹窗显示状态双向绑定
   */
  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  })

  const rules = computed<FormRules>(() => ({
    name: [
      { required: true, message: t('sys.role.roleNamePlaceholder'), trigger: 'blur' },
      { min: 2, max: 50, message: t('sys.user.usernameRule'), trigger: 'blur' }
    ],
    role: [
      { required: true, message: t('sys.role.roleKeyPlaceholder'), trigger: 'blur' },
      { min: 2, max: 50, message: t('sys.user.usernameRule'), trigger: 'blur' }
    ],
    status: [{ required: true, message: t('sys.user.statusPlaceholder'), trigger: 'change' }],
    sort: [
      { required: true, message: t('sys.post.sortPlaceholder'), trigger: 'blur' },
      { type: 'number', min: 0, message: t('sys.post.sortNegativeRule'), trigger: 'blur' }
    ]
  }))

  /**
   * 表单数据
   */
  const form = reactive<Api.SystemManage.RoleListItem>({
    id: '',
    name: '',
    role: '',
    status: 1,
    sort: 0,
    remark: ''
  })

  /**
   * 重置表单
   */
  const resetForm = (): void => {
    formRef.value?.resetFields()
    Object.assign(form, {
      id: '',
      name: '',
      role: '',
      status: 1,
      menuLinkage: 0,
      deptLinkage: 0,
      sort: 0,
      remark: '',
      menuIds: [],
      createTime: 0,
      updateTime: 0
    })
  }

  /**
   * 加载角色详情
   */
  const loadRoleData = async (): Promise<void> => {
    if (!props.roleData?.id) return

    try {
      submitLoading.value = true
      const detail = await fetchGetRoleInfo(props.roleData.id)
      Object.assign(form, {
        id: detail.id || '',
        name: detail.name || '',
        role: detail.role || '',
        status: detail.status ?? 1,
        sort: detail.sort ?? 0,
        remark: detail.remark || '',
        createTime: detail.createTime || 0
      })
    } catch (error) {
      console.error('获取角色详情失败:', error)
      ElMessage.error(t('sys.role.getDetailFailed'))
    } finally {
      submitLoading.value = false
    }
  }

  /**
   * 监听弹窗打开，初始化表单数据
   */
  watch(
    () => props.modelValue,
    async (newVal) => {
      if (newVal) {
        if (props.dialogType === 'edit' && props.roleData?.id) {
          await loadRoleData()
        } else {
          resetForm()
        }
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
   * 验证通过后调用接口保存数据
   */
  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()
      submitLoading.value = true

      const payload: Api.SystemManage.RoleListItem = {
        id: form.id,
        name: form.name,
        role: form.role,
        status: form.status,
        sort: form.sort,
        remark: form.remark,
        createTime: form.createTime
      }

      if (props.dialogType === 'add') {
        await fetchAddRole(payload)
      } else {
        await fetchEditRole(payload)
      }

      const message =
        props.dialogType === 'add' ? t('sys.common.addSuccess') : t('sys.common.editSuccess')
      ElMessage.success(message)
      emit('success')
      handleClose()
    } catch (error) {
      ElMessage.error(error instanceof Error ? error.message : t('sys.common.operateFailed'))
    } finally {
      submitLoading.value = false
    }
  }
</script>
