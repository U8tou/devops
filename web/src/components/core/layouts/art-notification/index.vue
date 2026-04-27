<!-- 通知组件 -->
<template>
  <!-- 直接使用抽屉作为主展示 -->
  <ElDrawer
    v-model="drawerVisible"
    size="400px"
    direction="rtl"
    class="art-notification-drawer"
    :with-header="false"
  >
    <div class="flex flex-col h-full">
      <div class="flex-cb shrink-0 pb-2">
        <span class="text-base font-medium text-g-800 dark:text-g-200">
          {{ $t('notice.title') }}
        </span>
        <div class="flex items-center gap-2">
          <ElButton type="primary" link size="small" @click="handleClearCurrentTab">
            {{ $t('notice.btnClear') }}
          </ElButton>
          <ElButton type="primary" link size="small" @click="handleMarkAllRead">
            {{ $t('notice.btnRead') }}
          </ElButton>
          <ElIcon class="c-p text-lg" @click="drawerVisible = false">
            <Close />
          </ElIcon>
        </div>
      </div>
      <ul class="flex items-end gap-4 border-b border-g-200 dark:border-g-700 shrink-0">
        <li
          v-for="(item, index) in barList"
          :key="index"
          class="text-[13px] text-g-700 dark:text-g-300 c-p select-none border-b-2 border-transparent"
          :class="{ 'bar-active': drawerTabIndex === index }"
          @click="drawerTabIndex = index"
        >
          {{ item.name }} ({{ item.num }})
        </li>
      </ul>
      <div class="flex-1 overflow-y-auto pt-2 pb-3 min-h-0 [&::-webkit-scrollbar]:w-1">
        <!-- 通知列表 -->
        <ul v-show="drawerTabIndex === 0" class="space-y-0">
          <li
            v-for="item in sortedNoticeList"
            :key="item.id"
            class="flex items-center gap-3 px-2 py-3 c-p rounded-lg hover:bg-g-200/60 dark:hover:bg-g-700/40"
            @click="openNoticeDetail(item)"
          >
            <div
              class="size-9 shrink-0 rounded-lg flex-cc"
              :class="[getNoticeIconStyle(item.type, item.tag, item.read).iconClass]"
            >
              <ArtSvgIcon
                class="text-lg !bg-transparent"
                :icon="getNoticeIconStyle(item.type, item.tag, item.read).icon"
              />
            </div>
            <div class="flex-1 min-w-0">
              <h4
                class="text-sm text-g-900 dark:text-g-100 truncate"
                :class="{ 'opacity-75': item.read }"
              >
                {{ item.title }}
              </h4>
              <p class="mt-0.5 text-xs text-g-500">{{ item.time }}</p>
            </div>
          </li>
        </ul>
        <!-- 消息列表 -->
        <ul v-show="drawerTabIndex === 1" class="space-y-0">
          <li
            v-for="item in sortedMsgList"
            :key="item.id"
            class="flex items-center gap-3 px-2 py-3 c-p rounded-lg hover:bg-g-200/60 dark:hover:bg-g-700/40"
            @click="openMessageDetail(item)"
          >
            <div
              class="size-9 shrink-0 rounded-lg flex-cc"
              :class="[getMessageIconStyle(item.read).iconClass]"
            >
              <ArtSvgIcon
                class="text-lg !bg-transparent"
                :icon="getMessageIconStyle(item.read).icon"
              />
            </div>
            <div class="flex-1 min-w-0">
              <h4
                class="text-sm text-g-900 dark:text-g-100 truncate"
                :class="{ 'opacity-75': item.read }"
              >
                {{ item.title }}
              </h4>
              <p class="mt-0.5 text-xs text-g-500">{{ item.time }}</p>
            </div>
          </li>
        </ul>
        <!-- 空状态 -->
        <div
          v-show="
            (drawerTabIndex === 0 && noticeList.length === 0) ||
            (drawerTabIndex === 1 && msgList.length === 0)
          "
          class="flex flex-col items-center justify-center py-16 text-g-500"
        >
          <ArtSvgIcon icon="system-uicons:inbox" class="text-5xl" />
          <p class="mt-3 text-xs">{{ $t('notice.text[0]') }}{{ barList[drawerTabIndex]?.name }}</p>
        </div>
      </div>
    </div>
  </ElDrawer>

  <!-- 通知详情弹窗 -->
  <ElDialog
    v-model="noticeDetailVisible"
    :title="$t('notice.noticeDetail')"
    width="480px"
    destroy-on-close
    class="art-notice-detail-dialog art-message-detail-dialog"
    @open="onNoticeDetailOpen"
  >
    <template v-if="currentNotice">
      <div
        class="rounded-xl border border-g-200 dark:border-g-700 bg-g-50 dark:bg-g-800/40 px-5 py-4 min-h-[80px] text-[15px] text-g-800 dark:text-g-200 leading-[1.7] whitespace-pre-wrap break-words"
      >
        {{ currentNotice.title }}
      </div>
      <div class="mt-4 flex justify-between items-center">
        <span
          class="inline-flex items-center rounded-md bg-g-100 dark:bg-g-700/60 px-2.5 py-1 text-xs text-g-600 dark:text-g-400"
        >
          {{ $t('notice.messageTime') }}: {{ currentNotice.time }}
        </span>
        <ElButton type="danger" link size="small" @click="handleDeleteNotice">
          {{ $t('notice.delete') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>

  <!-- 消息详情弹窗 -->
  <ElDialog
    v-model="messageDetailVisible"
    :title="$t('notice.messageDetail')"
    width="480px"
    destroy-on-close
    class="art-message-detail-dialog"
    @open="onMessageDetailOpen"
  >
    <template v-if="currentMessage">
      <div
        class="rounded-xl border border-g-200 dark:border-g-700 bg-g-50 dark:bg-g-800/40 px-5 py-4 min-h-[80px] text-[15px] text-g-800 dark:text-g-200 leading-[1.7] whitespace-pre-wrap break-words"
      >
        {{ currentMessage.title }}
      </div>
      <div class="mt-4 flex justify-between items-center">
        <span
          class="inline-flex items-center rounded-md bg-g-100 dark:bg-g-700/60 px-2.5 py-1 text-xs text-g-600 dark:text-g-400"
        >
          {{ $t('notice.messageTime') }}: {{ currentMessage.time }}
        </span>
        <ElButton type="danger" link size="small" @click="handleDeleteMessage">
          {{ $t('notice.delete') }}
        </ElButton>
      </div>
    </template>
  </ElDialog>
</template>

<script setup lang="ts">
  import { computed, ref, watch } from 'vue'
  import { storeToRefs } from 'pinia'
  import { useI18n } from 'vue-i18n'
  import { ElMessageBox } from 'element-plus'
  import { Close } from '@element-plus/icons-vue'
  import { useNotificationStore } from '@/store/modules/notification'
  import type { MessageItem, NoticeItem } from '@/utils/storage/notification-storage'

  defineOptions({ name: 'ArtNotification' })

  interface NoticeStyle {
    icon: string
    iconClass: string
  }

  const { t } = useI18n()
  const notificationStore = useNotificationStore()

  const { noticeList, msgList, sortedNoticeList, sortedMsgList } = storeToRefs(notificationStore)

  const props = defineProps<{
    value: boolean
  }>()

  const emit = defineEmits<{
    'update:value': [value: boolean]
  }>()

  const noticeDetailVisible = ref(false)
  const currentNotice = ref<NoticeItem | null>(null)
  const messageDetailVisible = ref(false)
  const currentMessage = ref<MessageItem | null>(null)
  const drawerVisible = ref(false)
  const drawerTabIndex = ref(0)

  // 标签栏数据（括号内展示未读数量）
  const barList = computed(() => [
    { name: t('notice.bar[0]'), num: noticeList.value.filter((n) => !n.read).length },
    { name: t('notice.bar[1]'), num: msgList.value.filter((m) => !m.read).length }
  ])

  // 样式映射（支持 type 和 tag）
  const noticeStyleMap: Record<string, NoticeStyle> = {
    email: {
      icon: 'ri:mail-line',
      iconClass: 'bg-warning/12 text-warning'
    },
    message: {
      icon: 'ri:volume-down-line',
      iconClass: 'bg-success/12 text-success'
    },
    collection: {
      icon: 'ri:heart-3-line',
      iconClass: 'bg-danger/12 text-danger'
    },
    user: {
      icon: 'ri:volume-down-line',
      iconClass: 'bg-info/12 text-info'
    },
    notice: {
      icon: 'ri:notification-3-line',
      iconClass: 'bg-theme/12 text-theme'
    },
    warning: {
      icon: 'ri:error-warning-line',
      iconClass: 'bg-warning/12 text-warning'
    }
  }

  function getNoticeStyle(type: string, tag?: string): NoticeStyle {
    const key = tag || type
    return noticeStyleMap[key] || noticeStyleMap[type] || noticeStyleMap['notice']
  }

  /** 通知图标：与消息一致，未读=彩色，已读=灰色；email 已读时用拆开信封 */
  function getNoticeIconStyle(type: string, tag: string | undefined, read: boolean): NoticeStyle {
    const base = getNoticeStyle(type, tag)
    const readStyle = 'bg-g-400/12 text-g-500 dark:bg-g-600/20 dark:text-g-500'
    if (read) {
      const isEmail = (tag || type) === 'email'
      return {
        icon: isEmail ? 'ri:mail-open-line' : base.icon,
        iconClass: readStyle
      }
    }
    return base
  }

  /** 邮件图标：未读=未拆信封，已读=拆开信封 */
  function getMessageIconStyle(read: boolean): NoticeStyle {
    return read
      ? {
          icon: 'ri:mail-open-line',
          iconClass: 'bg-g-400/12 text-g-500 dark:bg-g-600/20 dark:text-g-500'
        }
      : { icon: 'ri:mail-line', iconClass: 'bg-info/12 text-info' }
  }

  /** 清空当前选中 tab 的消息 */
  const handleClearCurrentTab = async () => {
    const list = drawerTabIndex.value === 0 ? noticeList.value : msgList.value
    if (list.length === 0) return
    const hasUnread = list.some((item) => !item.read)
    if (hasUnread) {
      try {
        await ElMessageBox.confirm(t('notice.clearConfirm'), t('notice.title'), {
          confirmButtonText: t('notice.btnClear'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        })
        notificationStore.clearCurrentTab(drawerTabIndex.value)
      } catch {
        // 用户取消
      }
    } else {
      notificationStore.clearCurrentTab(drawerTabIndex.value)
    }
  }

  /** 当前选中 tab 全部标为已读 */
  const handleMarkAllRead = () => {
    notificationStore.markCurrentTabRead(drawerTabIndex.value)
  }

  /** 打开通知详情（点开即标为已读） */
  const openNoticeDetail = (item: NoticeItem) => {
    currentNotice.value = item
    noticeDetailVisible.value = true
  }

  /** 打开消息详情（仿邮件，点开即标为已读） */
  const openMessageDetail = (item: MessageItem) => {
    currentMessage.value = item
    messageDetailVisible.value = true
  }

  /** 通知详情弹窗打开时标记为已读 */
  const onNoticeDetailOpen = () => {
    if (currentNotice.value?.id) {
      notificationStore.markNoticeRead(currentNotice.value.id)
    }
  }

  /** 弹窗打开时标记该消息为已读 */
  const onMessageDetailOpen = () => {
    if (currentMessage.value?.id) {
      notificationStore.markMessageRead(currentMessage.value.id)
    }
  }

  /** 删除当前通知 */
  const handleDeleteNotice = async () => {
    const id = currentNotice.value?.id
    if (!id) return
    try {
      await ElMessageBox.confirm(t('notice.deleteNoticeConfirm'), t('notice.noticeDetail'), {
        confirmButtonText: t('notice.delete'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      })
      notificationStore.deleteNotice(id)
      noticeDetailVisible.value = false
      currentNotice.value = null
    } catch {
      // 用户取消
    }
  }

  /** 删除当前消息 */
  const handleDeleteMessage = async () => {
    const id = currentMessage.value?.id
    if (!id) return
    try {
      await ElMessageBox.confirm(t('notice.deleteConfirm'), t('notice.messageDetail'), {
        confirmButtonText: t('notice.delete'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      })
      notificationStore.deleteMessage(id)
      messageDetailVisible.value = false
      currentMessage.value = null
    } catch {
      // 用户取消
    }
  }

  watch(
    () => props.value,
    (newValue) => {
      // 点击铃铛：直接打开抽屉
      if (newValue) {
        const hasUnreadNotice = noticeList.value.some((n) => !n.read)
        const hasUnreadMsg = msgList.value.some((m) => !m.read)
        drawerTabIndex.value = hasUnreadNotice ? 0 : hasUnreadMsg ? 1 : 0
        drawerVisible.value = true
        // 重置外部 v-model，避免下拉面板逻辑
        emit('update:value', false)
      }
    }
  )
</script>

<style scoped>
  @reference '@styles/core/tailwind.css';

  .bar-active {
    color: var(--theme-color) !important;
    border-bottom: 2px solid var(--theme-color);
  }
</style>

<style lang="scss">
  /* 消息详情弹窗：标题与正文紧凑间距 */
  .art-message-detail-dialog.el-dialog {
    .el-dialog__header {
      padding-bottom: 0;
      margin-bottom: 4px;
    }

    .el-dialog__body {
      padding-top: 0 !important;
      padding-bottom: 20px;
    }
  }

  /* 抽屉：紧凑布局，压缩空白 */
  .art-notification-drawer.el-drawer {
    .el-drawer__body {
      display: flex;
      flex-direction: column;
      height: 100%;
      padding: 0 20px;
      overflow: hidden;
    }
  }
</style>
