<template>
  <ArtSearchBar
    ref="searchBarRef"
    v-model="formData"
    :items="formItems"
    :rules="rules"
    @reset="handleReset"
    @search="handleSearch"
  >
  </ArtSearchBar>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  interface Props {
    modelValue: Record<string, any>
  }
  interface Emits {
    (e: 'update:modelValue', value: Record<string, any>): void
    (e: 'search', params: Record<string, any>): void
    (e: 'reset'): void
  }
  const props = defineProps<Props>()
  const emit = defineEmits<Emits>()

  // 表单数据双向绑定
  const searchBarRef = ref()
  const formData = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
  })

  // 校验规则
  const rules = {
    // userName: [{ required: true, message: '请输入用户名', trigger: 'blur' }]
  }

  const { t } = useI18n()
  // 状态选项（根据接口：1 正常，2 停用）
  const statusOptions = computed<{ label: string; value: number; disabled?: boolean }[]>(() => [
    { label: t('sys.common.statusNormal'), value: 1 },
    { label: t('sys.common.statusDisabled'), value: 2 }
  ])

  // 表单配置
  const formItems = computed(() => [
    {
      label: t('sys.user.username'),
      key: 'userName',
      type: 'input',
      placeholder: t('sys.user.usernamePlaceholder'),
      clearable: true
    },
    {
      label: t('sys.user.phone'),
      key: 'phone',
      type: 'input',
      props: { placeholder: t('sys.user.phonePlaceholder'), maxlength: '11' }
    },
    {
      label: t('sys.user.email'),
      key: 'email',
      type: 'input',
      props: { placeholder: t('sys.user.emailPlaceholder') }
    },
    {
      label: t('sys.user.status'),
      key: 'status',
      type: 'select',
      props: {
        placeholder: t('sys.user.statusPlaceholder'),
        options: statusOptions.value
      }
    },
    {
      label: t('sys.user.sex'),
      key: 'sex',
      type: 'radiogroup',
      props: {
        options: [
          { label: t('sys.common.sexMale'), value: 1 },
          { label: t('sys.common.sexFemale'), value: 2 }
        ]
      }
    },
    {
      label: t('sys.user.createDate'),
      key: 'daterange',
      type: 'datetime',
      props: {
        style: { width: '100%' },
        placeholder: t('sys.user.dateRangePlaceholder'),
        type: 'daterange',
        rangeSeparator: t('sys.user.rangeSeparator'),
        startPlaceholder: t('sys.user.startDate'),
        endPlaceholder: t('sys.user.endDate'),
        valueFormat: 'YYYY-MM-DD',
        shortcuts: [
          { text: t('sys.user.today'), value: [new Date(), new Date()] },
          { text: t('sys.user.lastWeek'), value: [new Date(Date.now() - 604800000), new Date()] },
          { text: t('sys.user.lastMonth'), value: [new Date(Date.now() - 2592000000), new Date()] }
        ]
      }
    }
  ])

  // 事件
  function handleReset() {
    console.log('重置表单')
    emit('reset')
  }

  async function handleSearch() {
    await searchBarRef.value.validate()
    emit('search', formData.value)
    console.log('表单数据', formData.value)
  }
</script>
