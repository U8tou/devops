/**
 * IM 聊天状态管理
 *
 * 管理 IM 聊天室 WebSocket 连接、消息列表、持久化，
 * 监听退出登录事件进行清理
 *
 * @module store/modules/imChat
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getImChatUrl } from '@/api/common'
import { mittBus } from '@/utils/sys'
import { useUserStore } from './user'
import {
  loadImChatMessages,
  saveImChatMessages,
  clearAllImChatMessages,
  genImChatMessageId,
  type ImChatMessage
} from '@/utils/storage/im-chat-storage'

const DEFAULT_ROOM_ID = '1'
let persistTimer: ReturnType<typeof setTimeout> | null = null
const PERSIST_DEBOUNCE_MS = 300

const RECONNECT_DELAY_INITIAL = 1000
const RECONNECT_DELAY_MAX = 30000
const RECONNECT_BACKOFF_FACTOR = 1.5

export const useImChatStore = defineStore('imChatStore', () => {
  const messages = ref<ImChatMessage[]>([])
  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const currentRoomId = ref<string>(DEFAULT_ROOM_ID)
  const hasJoinedOnce = ref(false)
  const isDrawerVisible = ref(false)
  const notificationText = ref('')
  const isReconnecting = ref(false)
  /** 当前房间在线人数（joined/left/online/offline 推送） */
  const memberCount = ref(0)
  /** 是否有未读消息（抽屉关闭时收到他人消息则 true，点击打开抽屉后置为已读） */
  const hasUnreadMessages = ref(false)
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let reconnectDelay = RECONNECT_DELAY_INITIAL
  let intentionalClose = false
  let notificationClearTimer: ReturnType<typeof setTimeout> | null = null
  const NOTIFICATION_DISPLAY_MS = 3000

  /** 获取当前用户 ID */
  function getUserId(): string {
    const userStore = useUserStore()
    return userStore.info?.userId ?? userStore.info?.userName ?? 'guest'
  }

  /** 获取当前用户的 loginId（用于判断 isMe，登录用户为 userId，游客为 guest/0） */
  function getMyLoginId(): string {
    const userStore = useUserStore()
    return userStore.info?.userId ?? userStore.info?.userName ?? 'guest'
  }

  /** 获取 accessToken */
  function getAccessToken(): string {
    const userStore = useUserStore()
    return userStore.accessToken ?? ''
  }

  /** 持久化消息（防抖） */
  function persistMessages() {
    if (persistTimer) clearTimeout(persistTimer)
    persistTimer = setTimeout(() => {
      saveImChatMessages(getUserId(), messages.value)
      persistTimer = null
    }, PERSIST_DEBOUNCE_MS)
  }

  /**
   * 判断消息是否为本人的（通过 loginId 判断）
   * 优先用 loginId/userId 比较；否则从 clientId(loginId_connId) 解析 loginId 比较
   */
  function computeIsMe(msg: { loginId?: string; clientId?: string; userId?: string }): boolean {
    const myLoginId = getMyLoginId()
    const msgLoginId = msg.loginId ?? msg.userId
    if (msgLoginId != null && String(msgLoginId) !== '') {
      if (String(msgLoginId) === String(myLoginId)) return true
      if (myLoginId === 'guest' && String(msgLoginId) === '0') return true
    }
    const cid = msg.clientId ?? ''
    if (cid && cid.includes('_')) {
      const extracted = cid.split('_')[0]
      if (String(extracted) === String(myLoginId)) return true
      if (myLoginId === 'guest' && (extracted === '0' || extracted === 'guest')) return true
    }
    return false
  }

  /** 从 storage 加载消息，按 loginId 重新计算 isMe */
  function loadFromStorage() {
    const loaded = loadImChatMessages(getUserId())
    messages.value = loaded.map((m) => ({
      ...m,
      isMe: computeIsMe({ loginId: m.loginId, clientId: m.clientId, userId: m.userId })
    }))
  }

  /** 展示通知消息，3 秒后恢复为默认提示 */
  function scheduleNotificationClear(text: string) {
    if (notificationClearTimer) {
      clearTimeout(notificationClearTimer)
      notificationClearTimer = null
    }
    notificationText.value = text
    notificationClearTimer = setTimeout(() => {
      notificationText.value = ''
      notificationClearTimer = null
    }, NOTIFICATION_DISPLAY_MS)
  }

  /** 建立 WebSocket 连接 */
  function connect() {
    const token = getAccessToken()
    const url = getImChatUrl(token || undefined)

    if (ws.value?.readyState === WebSocket.OPEN) return
    if (ws.value?.readyState === WebSocket.CONNECTING) return

    try {
      const socket = new WebSocket(url)
      ws.value = socket

      socket.onopen = () => {
        isConnected.value = true
        reconnectDelay = RECONNECT_DELAY_INITIAL
      }

      socket.onerror = (e) => {
        isConnected.value = false
        console.warn('[ImChat] WebSocket 连接错误，请检查后端服务与代理配置:', url, e)
      }

      socket.onmessage = (e: MessageEvent) => {
        try {
          const data = JSON.parse(e.data as string)
          const type = data.type
          if (type === 'message') {
            const formatTime = (t: string | number | undefined): string => {
              if (t == null || t === '') {
                return new Date().toLocaleTimeString([], {
                  hour: '2-digit',
                  minute: '2-digit'
                })
              }
              if (typeof t === 'number') {
                return new Date(t).toLocaleTimeString([], {
                  hour: '2-digit',
                  minute: '2-digit'
                })
              }
              return String(t)
            }
            const isMe = computeIsMe({
              loginId: data.loginId,
              clientId: data.clientId ?? data.from,
              userId: data.userId
            })
            const msg: ImChatMessage = {
              id: genImChatMessageId(),
              from: data.from ?? '',
              nickName: data.nickName ?? '',
              avatar: data.avatar ?? '',
              content: data.content ?? '',
              time: formatTime(data.time),
              isMe,
              clientId: data.clientId ?? data.from ?? '',
              loginId: data.loginId ?? data.userId,
              userId: data.userId
            }
            messages.value.push(msg)
            persistMessages()
            if (!isMe && !isDrawerVisible.value) hasUnreadMessages.value = true
          } else if (type === 'joined') {
            isConnected.value = true
            if (typeof data.memberCount === 'number') memberCount.value = data.memberCount
            const nickName = data.nickName ?? data.from ?? data.clientId ?? ''
            scheduleNotificationClear(nickName ? `${nickName} 加入聊天室` : '有成员加入聊天室')
          } else if (type === 'left') {
            if (typeof data.memberCount === 'number') memberCount.value = data.memberCount
            const nickName = data.nickName ?? data.from ?? data.clientId ?? ''
            scheduleNotificationClear(nickName ? `${nickName} 离开聊天室` : '有成员离开聊天室')
          } else if (type === 'online') {
            if (typeof data.memberCount === 'number') memberCount.value = data.memberCount
            const nickName = data.nickName ?? data.from ?? data.clientId ?? ''
            scheduleNotificationClear(nickName ? `${nickName} 上线了` : '有成员上线')
          } else if (type === 'offline') {
            if (typeof data.memberCount === 'number') memberCount.value = data.memberCount
            const nickName = data.nickName ?? data.from ?? data.clientId ?? ''
            scheduleNotificationClear(nickName ? `${nickName} 下线了` : '有成员下线')
          } else if (type === 'error') {
            console.warn('[ImChat] 服务端错误:', data.msg ?? data)
            scheduleNotificationClear(data.msg ?? '发生错误')
            isConnected.value = false
          }
        } catch (err) {
          console.warn('[ImChat] 解析消息失败:', e.data, err)
        }
      }

      socket.onclose = (e) => {
        isConnected.value = false
        ws.value = null
        if (e.code !== 1000 && !intentionalClose) {
          console.warn('[ImChat] WebSocket 连接关闭，将尝试重连:', e.code, e.reason)
          scheduleReconnect()
        }
      }
    } catch (err) {
      console.error('[ImChat] WebSocket 连接失败:', err)
      isConnected.value = false
    }
  }

  /** 加入聊天室 */
  function joinRoom(roomId: string = DEFAULT_ROOM_ID) {
    currentRoomId.value = roomId
    loadFromStorage()

    const socket = ws.value
    if (!socket || socket.readyState !== WebSocket.OPEN) {
      if (!socket) connect()
      // 等待连接建立后再发送 join（最多重试 30 次，约 3 秒）
      let retries = 0
      const maxRetries = 30
      const tryJoin = () => {
        if (ws.value?.readyState === WebSocket.OPEN) {
          ws.value!.send(JSON.stringify({ type: 'join', roomId }))
        } else if (retries < maxRetries) {
          retries++
          setTimeout(tryJoin, 100)
        }
      }
      setTimeout(tryJoin, 200)
      return
    }

    socket.send(JSON.stringify({ type: 'join', roomId }))
  }

  /** 尝试重连（断线时调用） */
  function scheduleReconnect() {
    if (isReconnecting.value || !hasJoinedOnce.value || !currentRoomId.value) return
    if (reconnectTimer) return

    isReconnecting.value = true
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      intentionalClose = false
      connect()
      joinRoom(currentRoomId.value)
      isReconnecting.value = false
      reconnectDelay = Math.min(reconnectDelay * RECONNECT_BACKOFF_FACTOR, RECONNECT_DELAY_MAX)
    }, reconnectDelay)
  }

  /** 取消重连 */
  function cancelReconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    reconnectDelay = RECONNECT_DELAY_INITIAL
    isReconnecting.value = false
  }

  /** 离开聊天室 */
  function leaveRoom() {
    intentionalClose = true
    cancelReconnect()
    if (notificationClearTimer) {
      clearTimeout(notificationClearTimer)
      notificationClearTimer = null
    }
    notificationText.value = ''
    hasUnreadMessages.value = false
    const socket = ws.value
    if (socket?.readyState === WebSocket.OPEN && currentRoomId.value) {
      try {
        socket.send(JSON.stringify({ type: 'leave', roomId: currentRoomId.value }))
      } catch (e) {
        console.warn('[ImChat] 发送 leave 失败:', e)
      }
      socket.close()
    }
    ws.value = null
    isConnected.value = false
    memberCount.value = 0
    currentRoomId.value = ''
    hasJoinedOnce.value = false
  }

  /** 发送消息（仅发送到服务端，由服务端回显后展示，避免重复） */
  function sendMessage(content: string) {
    const trimmed = content?.trim()
    if (!trimmed) return

    const socket = ws.value
    if (socket?.readyState === WebSocket.OPEN && currentRoomId.value) {
      socket.send(
        JSON.stringify({ type: 'message', roomId: currentRoomId.value, content: trimmed })
      )
    }
  }

  /** 清空消息并清除持久化 */
  function clearMessages() {
    messages.value = []
    clearAllImChatMessages()
  }

  /** 初始化连接（布局加载时调用，确保关闭抽屉时也能收到消息并亮起未读绿点） */
  function initConnection() {
    initLogoutListener()
    if (!hasJoinedOnce.value) {
      connect()
      joinRoom('1')
      hasJoinedOnce.value = true
    }
  }

  /** 打开聊天抽屉，点击即标记已读 */
  function openDrawer() {
    isDrawerVisible.value = true
    hasUnreadMessages.value = false
    initConnection()
  }

  /** 关闭聊天抽屉 */
  function closeDrawer() {
    isDrawerVisible.value = false
  }

  let logoutListenerSetup = false
  /** 初始化退出登录监听 */
  function initLogoutListener() {
    if (logoutListenerSetup) return
    logoutListenerSetup = true
    mittBus.on('imChatLogout', () => {
      leaveRoom()
      clearMessages()
    })
  }

  return {
    messages,
    ws,
    isConnected,
    currentRoomId,
    hasJoinedOnce,
    isDrawerVisible,
    notificationText,
    isReconnecting,
    memberCount,
    hasUnreadMessages,
    initConnection,
    connect,
    joinRoom,
    leaveRoom,
    sendMessage,
    clearMessages,
    loadFromStorage,
    openDrawer,
    closeDrawer,
    initLogoutListener
  }
})
