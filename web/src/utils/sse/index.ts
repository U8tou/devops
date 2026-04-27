/**
 * SSE (Server-Sent Events) 服务模块
 *
 * 建立与后端的 SSE 长连接，接收实时消息推送
 * 后端消息格式：event: message\ndata:tag-msg\n\n
 *
 * @module utils/sse
 */

import { h } from 'vue'
import { Icon } from '@iconify/vue'
import { getSseUrl } from '@/api/common'
import { ElNotification } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { mittBus } from '@/utils/sys'

/** SSE 消息解析结果 */
export interface SseMessage {
  tag: string
  msg: string
}

/** 消息类型标签（与 Go 后端 Tag 常量对应） */
const SseTag = {
  greeting: 'greeting', // 问候语
  message: 'message', // 消息
  notify: 'notify', // 通知
  warning: 'warning' // 警告
} as const

/** tag -> { type, icon, iconClass } 映射，仅用图标表示类型 */
const TAG_NOTIFICATION_MAP: Record<
  string,
  { type: 'success' | 'info' | 'warning' | 'error'; icon: string; iconClass: string }
> = {
  [SseTag.greeting]: {
    type: 'success',
    icon: 'ri:emotion-happy-line',
    iconClass: 'bg-success/12 text-success'
  },
  [SseTag.message]: {
    type: 'info',
    icon: 'ri:message-2-line',
    iconClass: 'bg-info/12 text-info'
  },
  [SseTag.notify]: {
    type: 'info',
    icon: 'ri:notification-3-line',
    iconClass: 'bg-theme/12 text-theme'
  },
  [SseTag.warning]: {
    type: 'warning',
    icon: 'ri:error-warning-line',
    iconClass: 'bg-warning/12 text-warning'
  }
}

let eventSource: EventSource | null = null
let currentToken: string | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let reconnectDelay = 3000
const RECONNECT_BASE_DELAY = 3000
const RECONNECT_MAX_DELAY = 60000

/**
 * 解析 SSE data 格式，后端格式为 "tag-msg"
 */
function parseSseData(data: string): SseMessage {
  const idx = data.indexOf('-')
  if (idx >= 0) {
    return {
      tag: data.slice(0, idx).trim(),
      msg: data.slice(idx + 1).trim()
    }
  }
  return { tag: '', msg: data }
}

/**
 * 展示 SSE 消息（右侧边弹出）
 * tag 用图标表示，不展示文字；消息内容合理布局
 */
function showSseNotification(message: SseMessage): void {
  const config = TAG_NOTIFICATION_MAP[message.tag] ?? {
    type: 'info' as const,
    icon: 'ri:information-line',
    iconClass: 'bg-info/12 text-info'
  }

  const iconWrap = h(
    'div',
    {
      class: `flex-shrink-0 size-9 rounded-full flex items-center justify-center ${config.iconClass}`
    },
    h(Icon, { icon: config.icon, class: 'text-lg' })
  )
  const msgWrap = h(
    'div',
    {
      class:
        'flex-1 min-w-0 pt-0.5 text-sm text-g-800 dark:text-g-200 leading-relaxed whitespace-pre-wrap break-words'
    },
    message.msg
  )
  const content = h('div', { class: 'flex items-start gap-3 min-w-0' }, [iconWrap, msgWrap])

  const instance = ElNotification({
    message: content,
    type: config.type,
    position: 'top-right',
    offset: 35,
    duration: config.type === 'warning' ? 6000 : 4500,
    zIndex: 10000,
    customClass: 'sse-notification-no-title sse-notification-clickable',
    showClose: false,
    onClick: () => {
      instance.close()
    }
  })
}

/**
 * 建立 SSE 连接（内部实现，支持延迟重连）
 */
function doConnect(token: string | null): void {
  try {
    const url = getSseUrl(token || undefined)
    eventSource = new EventSource(url)
    currentToken = token

    eventSource.addEventListener('message', (e: MessageEvent) => {
      const parsed = parseSseData(e.data || '')
      showSseNotification(parsed)
      // 推送到通知中心并持久化
      mittBus.emit('sseMessage', { tag: parsed.tag, msg: parsed.msg })
    })

    eventSource.addEventListener('error', () => {
      if (!eventSource) return
      eventSource.close()
      eventSource = null
      console.warn('[SSE] 连接异常，将在', reconnectDelay / 1000, '秒后重连')
      scheduleReconnect()
    })

    eventSource.addEventListener('open', () => {
      reconnectDelay = RECONNECT_BASE_DELAY
      console.info('[SSE] 已建立连接')
    })
  } catch (err) {
    console.error('[SSE] 连接失败:', err)
    scheduleReconnect()
  }
}

/** 调度延迟重连（指数退避，避免重连风暴） */
function scheduleReconnect(): void {
  if (reconnectTimer) return
  if (currentToken === null) return

  reconnectTimer = setTimeout(() => {
    reconnectTimer = null
    if (currentToken !== null) {
      doConnect(currentToken)
      reconnectDelay = Math.min(reconnectDelay * 2, RECONNECT_MAX_DELAY)
    }
  }, reconnectDelay)
}

/**
 * 建立 SSE 连接
 * @param token 鉴权 token（可选），有则登录用户定向推送，无则游客仅接收广播
 */
export function connectSse(token?: string): void {
  const tokenVal = token || null
  if (eventSource && currentToken === tokenVal) return

  clearReconnectTimer()
  disconnectSse()
  currentToken = tokenVal
  reconnectDelay = RECONNECT_BASE_DELAY

  if (tokenVal !== null) {
    doConnect(tokenVal)
  }
}

/** 清除重连定时器 */
function clearReconnectTimer(): void {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
}

/**
 * 断开 SSE 连接
 */
export function disconnectSse(): void {
  clearReconnectTimer()
  if (eventSource) {
    eventSource.close()
    eventSource = null
    currentToken = null
    reconnectDelay = RECONNECT_BASE_DELAY
    console.info('[SSE] 已断开连接')
  }
}

/**
 * 获取当前 SSE 连接状态
 */
export function isSseConnected(): boolean {
  return eventSource?.readyState === EventSource.OPEN
}

/**
 * 在登录状态下注册 SSE（供布局/应用入口调用）
 * 有 accessToken 时带 token 连接（登录用户定向推送），登出时需手动调用 disconnectSse
 */
export function registerSseWhenLoggedIn(): () => void {
  const token = useUserStore().accessToken
  if (token) {
    connectSse(token)
  }
  return disconnectSse
}
