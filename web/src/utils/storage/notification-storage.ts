/**
 * 通知消息持久化存储
 *
 * 按用户维度将通知数据持久化到 localStorage，
 * 支持通知、消息两种类型
 *
 * @module utils/storage/notification-storage
 */

const NOTIFICATION_PREFIX = 'sys-notification-'

/** 通知类型 */
export type NoticeTabType = 'notice' | 'message'

/** 通知项 */
export interface NoticeItem {
  id: string
  title: string
  time: string
  type: string
  read: boolean
  /** 原始 tag，用于图标映射 */
  tag?: string
}

/** 消息项（带头像） */
export interface MessageItem {
  id: string
  title: string
  time: string
  avatar: string
  read: boolean
}

export interface NotificationPersistData {
  notice: NoticeItem[]
  message: MessageItem[]
}

function getStorageKey(userId: string): string {
  return `${NOTIFICATION_PREFIX}${userId || 'guest'}`
}

/**
 * 从 localStorage 加载通知数据
 */
export function loadNotifications(userId: string): NotificationPersistData {
  try {
    const key = getStorageKey(userId)
    const raw = localStorage.getItem(key)
    if (!raw) {
      return { notice: [], message: [] }
    }
    const data = JSON.parse(raw) as NotificationPersistData
    return {
      notice: Array.isArray(data.notice) ? data.notice : [],
      message: Array.isArray(data.message) ? data.message : []
    }
  } catch {
    return { notice: [], message: [] }
  }
}

/**
 * 持久化通知数据到 localStorage
 */
export function saveNotifications(userId: string, data: NotificationPersistData): void {
  try {
    const key = getStorageKey(userId)
    localStorage.setItem(key, JSON.stringify(data))
  } catch (e) {
    console.warn('[NotificationStorage] 持久化失败:', e)
  }
}

/**
 * 限制数组长度，保留最新 N 条
 */
export function trimList<T>(arr: T[], max: number): T[] {
  if (arr.length <= max) return arr
  return arr.slice(-max)
}

/**
 * 生成唯一 ID
 */
export function genNotificationId(): string {
  return `n_${Date.now()}_${Math.random().toString(36).slice(2, 9)}`
}

/**
 * 格式化时间为 YYYY-M-D H:mm
 */
export function formatNoticeTime(date: Date = new Date()): string {
  const y = date.getFullYear()
  const m = date.getMonth() + 1
  const d = date.getDate()
  const h = date.getHours()
  const min = date.getMinutes()
  return `${y}-${m}-${d} ${h}:${min.toString().padStart(2, '0')}`
}
