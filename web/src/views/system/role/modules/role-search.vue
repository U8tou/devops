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
  const { t } = useI18n()
  const searchBarRef = ref()

  const formData = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
  })

  const rules = {}

  const statusOptions = computed(() => [
    { label: t('sys.common.statusNormal'), value: 1 },
    { label: t('sys.common.statusDisabled'), value: 2 }
  ])

  const formItems = computed(() => [
    {
      label: t('sys.role.roleName'),
      key: 'name',
      type: 'input',
      placeholder: t('sys.role.roleNamePlaceholder'),
      clearable: true
    },
    {
      label: t('sys.role.roleKey'),
      key: 'role',
      type: 'input',
      placeholder: t('sys.role.roleKeyPlaceholder'),
      clearable: true
    },
    {
      label: t('sys.role.roleStatus'),
      key: 'status',
      type: 'select',
      props: {
        placeholder: t('sys.user.statusPlaceholder'),
        options: statusOptions.value,
        clearable: true
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

  /**
   * 处理重置事件
   */
  const handleReset = () => {
    emit('reset')
  }

  /**
   * 处理搜索事件
   * 验证表单后触发搜索
   */
  const handleSearch = async () => {
    await searchBarRef.value.validate()
    emit('search', formData.value)
  }
</script>
