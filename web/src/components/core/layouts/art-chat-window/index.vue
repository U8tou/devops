<!-- 系统聊天窗口 -->
<template>
  <div>
    <ElDrawer
      :model-value="imChatStore.isDrawerVisible"
      @update:model-value="
        (v) => {
          if (!v) imChatStore.closeDrawer()
        }
      "
      :size="isMobile ? '100%' : '480px'"
      :with-header="false"
    >
      <div class="mb-5 flex-cb">
        <div>
          <span class="text-base font-medium">{{ $t('chat.title', 'IM 聊天室') }}</span>
          <div class="mt-1.5 flex items-center gap-2">
            <div class="flex-c gap-1">
              <div
                class="h-2 w-2 rounded-full"
                :class="isOnline ? 'bg-success/100' : 'bg-danger/100'"
              ></div>
              <span class="text-xs text-g-600">{{ isOnline ? '在线' : '离线' }}</span>
            </div>
            <span v-if="isOnline && memberCount > 0" class="text-xs text-g-500">
              {{ $t('chat.memberCount', { count: memberCount }) }}
            </span>
          </div>
        </div>
        <div class="flex-c gap-2">
          <ElButton type="default" size="small" @click="handleClearMessages">
            {{ $t('chat.clearMessages', '清空消息') }}
          </ElButton>
          <ElIcon class="c-p" :size="20" @click="closeChat">
            <Close />
          </ElIcon>
        </div>
      </div>

      <!-- 常驻通知横幅：通知为空时显示退出提示，否则显示上线/下线等通知，展示不全时才滚动 -->
      <div
        class="mb-3 flex items-center gap-2 overflow-hidden rounded-md bg-g-200/60 px-3 py-2 text-xs text-g-600 dark:bg-g-600/30"
      >
        <ArtSvgIcon icon="ri:megaphone-line" class="shrink-0 text-base text-theme" />
        <div ref="notificationContainerRef" class="relative min-w-0 flex-1 overflow-hidden">
          <!-- 用于测量的隐藏元素，检测文本是否溢出 -->
          <div
            ref="notificationMeasureRef"
            class="absolute left-0 top-0 whitespace-nowrap opacity-0 pointer-events-none"
            style="max-width: none; visibility: hidden"
            aria-hidden="true"
          >
            {{ displayNotification }}
          </div>
          <div v-if="needMarquee" class="chat-notification-marquee inline-flex whitespace-nowrap">
            <span class="inline-block pr-8">{{ displayNotification }}</span>
            <span class="inline-block pr-8">{{ displayNotification }}</span>
          </div>
          <div v-else class="overflow-hidden truncate whitespace-nowrap">
            {{ displayNotification }}
          </div>
        </div>
      </div>

      <div class="flex h-[calc(100%-140px)] flex-col">
        <!-- 聊天消息区域 -->
        <div
          class="flex-1 overflow-y-auto border-t-d px-4 py-7.5 [&::-webkit-scrollbar]:!w-1"
          ref="messageContainer"
        >
          <div
            v-if="chatMessages.length === 0"
            class="flex flex-1 flex-col items-center justify-center py-12 text-g-500 text-sm"
          >
            <span>{{
              isReconnecting
                ? '正在重连...'
                : isOnline
                  ? '暂无消息，发送一条开始聊天吧'
                  : '连接中...'
            }}</span>
          </div>
          <template v-else v-for="(message, index) in chatMessages" :key="message.id || index">
            <div
              :class="[
                'mb-7.5 flex w-full items-start gap-2',
                message.isMe ? 'flex-row-reverse' : 'flex-row'
              ]"
            >
              <ElAvatar :size="32" :src="resolveAvatar(message.avatar)" class="shrink-0" />
              <div
                :class="['flex max-w-[70%] flex-col', message.isMe ? 'items-end' : 'items-start']"
              >
                <div
                  :class="[
                    'mb-1 flex flex-col',
                    message.isMe ? 'items-end text-right' : 'items-start text-left'
                  ]"
                >
                  <span class="text-[13px] font-medium text-g-800 dark:text-g-200">
                    {{ parsedName(message).main }}
                  </span>
                  <span
                    v-if="parsedName(message).suffix"
                    class="mt-0.5 text-[10px] text-g-500 dark:text-g-400"
                  >
                    {{ parsedName(message).suffix }}
                  </span>
                </div>
                <div
                  :class="[
                    'rounded-md px-3.5 py-2.5 text-sm leading-[1.4] text-g-900',
                    message.isMe
                      ? 'message-right bg-theme/15'
                      : 'message-left bg-g-200/80 dark:bg-g-600/50'
                  ]"
                >
                  {{ message.content }}
                </div>
                <span class="mt-0.5 text-[10px] text-g-500">{{ message.time }}</span>
              </div>
            </div>
          </template>
        </div>

        <!-- 聊天输入区域 -->
        <div class="px-4 pt-4">
          <ElInput
            v-model="messageText"
            type="textarea"
            :rows="3"
            :placeholder="$t('chat.inputPlaceholder', '输入消息')"
            resize="none"
            @keyup.enter.prevent="sendMessage"
          >
            <template #append>
              <div class="flex gap-2 py-2">
                <ElButton :icon="Paperclip" circle plain />
                <ElButton :icon="Picture" circle plain />
                <ElButton type="primary" @click="sendMessage" v-ripple>
                  {{ $t('chat.send', '发送') }}
                </ElButton>
              </div>
            </template>
          </ElInput>
          <div class="mt-3 flex-cb">
            <div class="flex-c">
              <ArtSvgIcon icon="ri:image-line" class="mr-5 c-p text-g-600 text-lg" />
              <ArtSvgIcon icon="ri:emotion-happy-line" class="mr-5 c-p text-g-600 text-lg" />
            </div>
            <ElButton type="primary" @click="sendMessage" v-ripple class="min-w-20">
              {{ $t('chat.send', '发送') }}
            </ElButton>
          </div>
        </div>
      </div>
    </ElDrawer>
  </div>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import { Picture, Paperclip, Close } from '@element-plus/icons-vue'
  import { ElMessageBox } from 'element-plus'
  import { useImChatStore } from '@/store/modules/imChat'
  import { resolveImageUrl } from '@/utils/resolveImageUrl'
  import defaultAvatar from '@imgs/user/avatar.webp'

  defineOptions({ name: 'ArtChatWindow' })

  const { t } = useI18n()
  const MOBILE_BREAKPOINT = 640
  const SCROLL_DELAY = 100

  const imChatStore = useImChatStore()
  const {
    messages: chatMessages,
    isConnected,
    notificationText,
    isReconnecting,
    memberCount
  } = storeToRefs(imChatStore)

  const exitWarningText = computed(() => t('chat.exitWarning', '若是退出系统，将清空聊天内容'))
  const displayNotification = computed(() => notificationText.value || exitWarningText.value)

  const { width } = useWindowSize()
  const notificationContainerRef = ref<HTMLElement | null>(null)
  const notificationMeasureRef = ref<HTMLElement | null>(null)
  const needMarquee = ref(false)

  /** 检测文本是否展示不全，仅展示不全时启用滚动 */
  const checkNotificationOverflow = () => {
    nextTick(() => {
      const container = notificationContainerRef.value
      const measureEl = notificationMeasureRef.value
      if (!container || !measureEl) {
        needMarquee.value = false
        return
      }
      needMarquee.value = measureEl.scrollWidth > container.clientWidth
    })
  }

  watch([displayNotification, () => width.value], checkNotificationOverflow, {
    immediate: true,
    flush: 'post'
  })
  const isMobile = computed(() => width.value < MOBILE_BREAKPOINT)
  const messageText = ref('')
  const messageContainer = ref<HTMLElement | null>(null)

  const isOnline = computed(() => isConnected.value)

  const resolveAvatar = (avatar: string | undefined): string => {
    if (!avatar) return defaultAvatar
    return resolveImageUrl(avatar) || defaultAvatar
  }

  /** 昵称展示：nickName 为纯昵称（如 小王、游客），直接使用 */
  const parsedName = (msg: {
    from?: string
    nickName?: string
  }): { main: string; suffix: string } => {
    const main = msg.nickName || msg.from || '未知'
    return { main, suffix: '' }
  }

  const scrollToBottom = (): void => {
    nextTick(() => {
      setTimeout(() => {
        if (messageContainer.value) {
          messageContainer.value.scrollTop = messageContainer.value.scrollHeight
        }
      }, SCROLL_DELAY)
    })
  }

  const sendMessage = (): void => {
    const text = messageText.value.trim()
    if (!text) return
    imChatStore.sendMessage(text)
    messageText.value = ''
    scrollToBottom()
  }

  const closeChat = (): void => {
    imChatStore.closeDrawer()
  }

  onMounted(() => {
    imChatStore.initConnection()
  })

  const handleClearMessages = (): void => {
    ElMessageBox.confirm(t('chat.clearConfirm'), t('chat.clearMessages'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    }).then(() => {
      imChatStore.clearMessages()
    })
  }

  watch(chatMessages, () => scrollToBottom(), { deep: true })

  watch(
    () => imChatStore.isDrawerVisible,
    (visible) => {
      if (visible) scrollToBottom()
    }
  )
</script>

<style scoped>
  .chat-notification-marquee {
    animation: chatMarquee 15s linear infinite;
  }

  .chat-notification-marquee:hover {
    animation-play-state: paused;
  }

  @keyframes chatMarquee {
    0% {
      transform: translateX(0);
    }

    100% {
      transform: translateX(-50%);
    }
  }
</style>
