<template>
  <ElDialog
    v-model="visible"
    :title="
      dialogType === 'add' ? $t('sys.post.addPost') : $t('sys.post.editPost', { id: form.id })
    "
    width="640px"
    align-center
    @close="handleClose"
  >
    <ElForm ref="formRef" :model="form" :rules="rules" label-position="top">
      <ElFormItem :label="$t('sys.post.postName')" prop="name">
        <ElInput v-model="form.name" :placeholder="$t('sys.post.postNamePlaceholder')" />
      </ElFormItem>
      <ElFormItem :label="$t('sys.user.status')" prop="status">
        <ElRadioGroup v-model="form.status">
          <ElRadio :label="1">{{ $t('sys.common.statusNormal') }}</ElRadio>
          <ElRadio :label="2">{{ $t('sys.common.statusDisabled') }}</ElRadio>
        </ElRadioGroup>
      </ElFormItem>
      <ElFormItem :label="$t('sys.dept.sort')" prop="sort">
        <ElInputNumber
          v-model="form.sort"
          :min="0"
          controls-position="right"
          style="width: 160px"
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
  import { fetchGetPostInfo, fetchAddPost, fetchEditPost } from '@/api/system-manage'

  interface Props {
    modelValue: boolean
    dialogType: 'add' | 'edit'
    postData?: { id: string }
  }

  interface Emits {
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    dialogType: 'add',
    postData: undefined
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

  /**
   * 表单验证规则
   */
  const rules = computed<FormRules>(() => ({
    name: [
      { required: true, message: t('sys.post.postNamePlaceholder'), trigger: 'blur' },
      { min: 2, max: 50, message: t('sys.post.postNameRule'), trigger: 'blur' }
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
  const form = reactive<Api.SystemManage.PostAddReq & { id?: string }>({
    id: '',
    name: '',
    status: 1,
    sort: 999
  })

  /**
   * 重置表单
   */
  const resetForm = (): void => {
    formRef.value?.resetFields()
    Object.assign(form, {
      id: '',
      name: '',
      status: 1,
      sort: 999
    })
  }

  /**
   * 加载岗位详情
   */
  const loadPostData = async (): Promise<void> => {
    if (!props.postData?.id) return

    try {
      submitLoading.value = true
      const detail = await fetchGetPostInfo(props.postData.id)
      Object.assign(form, {
        id: detail.id || '',
        name: detail.name || '',
        status: detail.status ?? 1,
        sort: detail.sort ?? 0
      })
    } catch (error) {
      console.error('获取岗位详情失败:', error)
      ElMessage.error(t('sys.post.getDetailFailed'))
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
        if (props.dialogType === 'edit' && props.postData?.id) {
          await loadPostData()
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

      if (props.dialogType === 'add') {
        const payload: Api.SystemManage.PostAddReq = {
          name: form.name,
          status: form.status,
          sort: form.sort
        }
        await fetchAddPost(payload)
      } else {
        const payload: Api.SystemManage.PostEditReq = {
          id: form.id!,
          name: form.name,
          status: form.status,
          sort: form.sort
        }
        await fetchEditPost(payload)
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
