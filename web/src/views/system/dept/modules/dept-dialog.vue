<template>
  <ElDialog
    :title="dialogTitle"
    :model-value="visible"
    @update:model-value="handleCancel"
    width="640px"
    align-center
    class="dept-dialog"
    @closed="handleClosed"
  >
    <ArtForm
      ref="formRef"
      v-model="form"
      :items="formItems"
      :rules="rules"
      :span="width > 640 ? 12 : 24"
      :gutter="20"
      label-position="top"
      :show-reset="false"
      :show-submit="false"
    />

    <template #footer>
      <span class="dialog-footer">
        <ElButton @click="handleCancel" :disabled="submitLoading">
          {{ $t('common.cancel') }}
        </ElButton>
        <ElButton type="primary" :loading="submitLoading" @click="handleSubmit">
          {{ $t('common.confirm') }}
        </ElButton>
      </span>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import type { FormRules } from 'element-plus'
  import { useI18n } from 'vue-i18n'
  import { ElMessage } from 'element-plus'
  import type { FormItem } from '@/components/core/forms/art-form/index.vue'
  import ArtForm from '@/components/core/forms/art-form/index.vue'
  import { useWindowSize } from '@vueuse/core'
  import {
    fetchAddDept,
    fetchEditDept,
    fetchGetDeptAll,
    fetchGetDeptInfo
  } from '@/api/system-manage'
  import { handleTree } from '@/utils/ruoyi'

  const { width } = useWindowSize()

  /** 部门表单/行数据（与后端结构一致：id, pid, name, profile, leader, phone, email, sort, status） */
  export interface DeptFormData {
    id: string
    pid: string // 上级部门ID，0 为顶级
    name: string
    profile: string
    leader: string
    phone: string
    email: string
    sort: number
    status: number // 1 正常, 0 停用
  }

  /** 部门树节点（含 children） */
  type DeptTreeNode = Api.SystemManage.DeptListItem & { children?: DeptTreeNode[] }

  interface Props {
    visible: boolean
    editData?: { id: string } | null
  }

  interface Emits {
    (e: 'update:visible', value: boolean): void
    (e: 'submit', data: DeptFormData): void
  }

  const props = withDefaults(defineProps<Props>(), {
    visible: false,
    editData: null
  })

  const emit = defineEmits<Emits>()

  const formRef = ref()
  const isEdit = ref(false)

  const form = reactive<DeptFormData>({
    id: '0',
    pid: '0',
    name: '',
    profile: '',
    leader: '',
    phone: '',
    email: '',
    sort: 1,
    status: 1
  })

  const { t } = useI18n()
  const deptTreeLoading = ref(false)
  const deptTree = ref<DeptTreeNode[]>([])

  const rules = computed<FormRules>(() => ({
    name: [
      { required: true, message: t('sys.dept.deptNameRule'), trigger: 'blur' },
      { min: 2, max: 50, message: t('sys.user.usernameRule'), trigger: 'blur' }
    ],
    leader: [{ required: true, message: t('sys.dept.leaderPlaceholder'), trigger: 'blur' }],
    phone: [
      {
        pattern: /^1[3-9]\d{9}$|^$/,
        message: t('sys.dept.leaderPhoneRule'),
        trigger: 'blur'
      }
    ],
    email: [{ type: 'email', message: t('sys.dept.leaderEmailPlaceholder'), trigger: 'blur' }],
    sort: [
      { required: true, message: t('sys.dept.sortPlaceholder'), trigger: 'blur' },
      { type: 'number', min: 0, message: t('sys.dept.sortNegativeRule'), trigger: 'blur' }
    ],
    status: [{ required: true, message: t('sys.user.statusPlaceholder'), trigger: 'change' }]
  }))

  const statusOptions = computed(() => [
    { label: t('sys.common.statusNormal'), value: 1 },
    { label: t('sys.common.statusDisabled'), value: 0 }
  ])

  /** 上级部门树（编辑时排除当前节点及其子节点，避免循环） */
  const parentTreeOptions = computed<DeptTreeNode[]>(() => {
    const excludeId = form.id ? Number(form.id) : null
    if (!excludeId) return deptTree.value
    const filterTree = (nodes: DeptTreeNode[]): DeptTreeNode[] =>
      nodes
        .filter((n) => Number(n.id) !== excludeId)
        .map((n) => ({
          ...n,
          children: n.children?.length ? filterTree(n.children) : undefined
        }))
    return filterTree(deptTree.value)
  })

  /** 带“顶级”选项的树（用于下拉，顶级 id 为 0） */
  const parentTreeOptionsWithRoot = computed(() => [
    { id: '0', name: t('sys.dept.topDept'), children: parentTreeOptions.value } as DeptTreeNode
  ])

  const formItems = computed<FormItem[]>(() => {
    const switchSpan = width.value < 640 ? 24 : 12
    return [
      {
        label: t('sys.dept.parentDept'),
        key: 'pid',
        type: 'treeselect',
        span: 24,
        props: {
          data: parentTreeOptionsWithRoot.value,
          props: { label: 'name', value: 'id' },
          nodeKey: 'id',
          placeholder: t('sys.dept.parentPlaceholderTop'),
          clearable: true,
          style: { width: '100%' },
          checkStrictly: true,
          renderAfterExpand: false
        }
      },
      {
        label: t('sys.dept.deptName'),
        key: 'name',
        type: 'input',
        props: { placeholder: t('sys.dept.deptNamePlaceholder') }
      },
      {
        label: t('sys.dept.profile'),
        key: 'profile',
        type: 'input',
        props: {
          placeholder: t('sys.dept.profilePlaceholder'),
          type: 'textarea',
          autosize: { minRows: 2 }
        },
        span: 24
      },
      {
        label: t('sys.dept.leader'),
        key: 'leader',
        type: 'input',
        props: { placeholder: t('sys.dept.leaderPlaceholder') },
        span: switchSpan
      },
      {
        label: t('sys.dept.leaderPhone'),
        key: 'phone',
        type: 'input',
        props: { placeholder: t('sys.dept.leaderPhonePlaceholder') },
        span: switchSpan
      },
      {
        label: t('sys.dept.leaderEmail'),
        key: 'email',
        type: 'input',
        props: { placeholder: t('sys.dept.leaderEmailPlaceholder') },
        span: switchSpan
      },
      {
        label: t('sys.dept.sort'),
        key: 'sort',
        type: 'number',
        props: { min: 0, controlsPosition: 'right', style: { width: '100%' } },
        span: switchSpan
      },
      {
        label: t('sys.user.status'),
        key: 'status',
        type: 'select',
        props: {
          placeholder: t('sys.user.statusPlaceholder'),
          options: statusOptions.value,
          style: { width: '100%' }
        },
        span: switchSpan
      }
    ]
  })

  const dialogTitle = computed(() =>
    isEdit.value && form.id ? t('sys.dept.editDept', { id: form.id }) : t('sys.dept.addDeptTitle')
  )

  const resetForm = (): void => {
    formRef.value?.reset?.()
    Object.assign(form, {
      id: '0',
      pid: '0',
      name: '',
      profile: '',
      leader: '',
      phone: '',
      email: '',
      sort: 1,
      status: 1
    })
  }

  /** 打开弹窗时拉取部门树 */
  const loadDeptTree = async (): Promise<void> => {
    deptTreeLoading.value = true
    try {
      const res = await fetchGetDeptAll()
      const rows = res?.rows ?? []
      deptTree.value = handleTree(
        rows as unknown as Record<string, unknown>[],
        'id',
        'pid',
        'children'
      ) as unknown as DeptTreeNode[]
    } finally {
      deptTreeLoading.value = false
    }
  }

  const loadFormData = async (): Promise<void> => {
    if (!props.editData || !props.editData.id) return

    try {
      submitLoading.value = true
      const detail = await fetchGetDeptInfo(props.editData.id)
      isEdit.value = true
      form.id = detail.id || '0'
      form.pid = detail.pid !== undefined && detail.pid !== null ? String(detail.pid) : '0'
      form.name = detail.name ?? ''
      form.profile = detail.profile ?? ''
      form.leader = detail.leader ?? ''
      form.phone = detail.phone ?? ''
      form.email = detail.email ?? ''
      form.sort = detail.sort !== undefined && detail.sort !== null ? Number(detail.sort) : 1
      form.status =
        detail.status !== undefined && detail.status !== null ? Number(detail.status) : 1
    } catch (error) {
      console.error('获取部门详情失败:', error)
      ElMessage.error(t('sys.dept.getDetailFailed'))
    } finally {
      submitLoading.value = false
    }
  }

  const submitLoading = ref(false)

  const handleSubmit = async (): Promise<void> => {
    if (!formRef.value) return
    try {
      await formRef.value.validate()
      submitLoading.value = true
      const payload: Api.SystemManage.DeptListItem = {
        id: form.id,
        pid: form.pid ?? '0',
        name: form.name,
        profile: form.profile,
        leader: form.leader,
        phone: form.phone,
        email: form.email,
        sort: form.sort,
        status: form.status
      }
      if (isEdit.value && form.id) {
        await fetchEditDept(payload)
      } else {
        await fetchAddDept({ ...payload })
      }
      ElMessage.success(isEdit.value ? t('sys.common.editSuccess') : t('sys.common.addSuccess'))
      handleCancel()
      emit('submit', { ...form })
    } catch (err) {
      if (err !== false && err !== undefined) {
        ElMessage.error(err instanceof Error ? err.message : t('sys.common.formValidateFailed'))
      }
    } finally {
      submitLoading.value = false
    }
  }

  const handleCancel = (): void => {
    emit('update:visible', false)
  }

  const handleClosed = (): void => {
    resetForm()
    isEdit.value = false
  }

  watch(
    () => props.visible,
    async (newVal) => {
      if (newVal) {
        await loadDeptTree()
        nextTick(async () => {
          if (props.editData && props.editData.id !== undefined && props.editData.id !== null) {
            await loadFormData()
          } else {
            resetForm()
            isEdit.value = false
          }
        })
      }
    }
  )
</script>
