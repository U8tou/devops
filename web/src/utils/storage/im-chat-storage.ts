/**
 * IM 聊天消息持久化存储
 *
 * 按用户维度将聊天消息持久化到 localStorage，
 * 退出登录时需清空
 *
 * @module utils/storage/im-chat-storage
 */

const IM_CHAT_PREFIX = 'sys-im-chat-'

/** IM 聊天消息项 */
export interface ImChatMessage {
  id: string
  from: string
  nickName: string
  avatar: string
  content: string
  time: string
  isMe: boolean
  /** 会话唯一标识 loginId_connId */
  clientId?: string
  /** 登录标识，用于判断 isMe（message.loginId === 本人 loginId） */
  loginId?: string
  /** 用户 ID（与 loginId 可能相同） */
  userId?: string
}

export interface ImChatPersistData {
  messages: ImChatMessage[]
}

function getStorageKey(userId: string): string {
  return `${IM_CHAT_PREFIX}${userId || 'guest'}`
}

/**
 * 从 localStorage 加载 IM 聊天消息
 */
export function loadImChatMessages(userId: string): ImChatMessage[] {
  try {
    const key = getStorageKey(userId)
    const raw = localStorage.getItem(key)
    if (!raw) {
      return []
    }
    const data = JSON.parse(raw) as ImChatPersistData
    return Array.isArray(data.messages) ? data.messages : []
  } catch {
    return []
  }
}

/**
 * 持久化 IM 聊天消息到 localStorage
 */
export function saveImChatMessages(userId: string, messages: ImChatMessage[]): void {
  try {
    const key = getStorageKey(userId)
    localStorage.setItem(key, JSON.stringify({ messages }))
  } catch (e) {
    console.warn('[ImChatStorage] 持久化失败:', e)
  }
}

/**
 * 清空 IM 聊天消息（指定用户）
 */
export function clearImChatMessages(userId: string): void {
  try {
    const key = getStorageKey(userId)
    localStorage.removeItem(key)
  } catch (e) {
    console.warn('[ImChatStorage] 清空失败:', e)
  }
}

/**
 * 清空所有用户的 IM 聊天消息（用于登出时）
 */
export function clearAllImChatMessages(): void {
  try {
    const keysToRemove: string[] = []
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (key && key.startsWith(IM_CHAT_PREFIX)) {
        keysToRemove.push(key)
      }
    }
    keysToRemove.forEach((key) => localStorage.removeItem(key))
  } catch (e) {
    console.warn('[ImChatStorage] 清空全部失败:', e)
  }
}

/**
 * 生成消息唯一 ID
 */
export function genImChatMessageId(): string {
  return `im_${Date.now()}_${Math.random().toString(36).slice(2, 9)}`
}
