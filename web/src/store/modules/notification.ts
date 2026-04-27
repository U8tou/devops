/**
 * 通知中心状态管理
 *
 * 管理通知、消息两 tab 的数据，与 SSE 消息联动，
 * 持久化到 localStorage
 *
 * @module store/modules/notification
 */

import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { useUserStore } from './user'
import { mittBus } from '@/utils/sys'
import avatar1 from '@/assets/images/avatar/avatar1.webp'
import avatar2 from '@/assets/images/avatar/avatar2.webp'
import avatar3 from '@/assets/images/avatar/avatar3.webp'
import avatar4 from '@/assets/images/avatar/avatar4.webp'
import avatar5 from '@/assets/images/avatar/avatar5.webp'
import avatar6 from '@/assets/images/avatar/avatar6.webp'

import {
  loadNotifications,
  saveNotifications,
  genNotificationId,
  formatNoticeTime,
  trimList,
  type NoticeItem,
  type MessageItem
} from '@/utils/storage/notification-storage'

const MAX_ITEMS = 100

// SSE tag -> 归类：'notice' | 'message'
// 仅 greeting（问候语/温馨提示）归消息，其余皆归通知
const SSE_TAG_TO_TAB: Record<string, 'notice' | 'message'> = {
  greeting: 'message',
  message: 'notice',
  notify: 'notice',
  warning: 'notice'
}

const DEFAULT_AVATARS = [avatar1, avatar2, avatar3, avatar4, avatar5, avatar6]

function pickAvatar(index: number): string {
  return DEFAULT_AVATARS[index % DEFAULT_AVATARS.length]
}

export const useNotificationStore = defineStore('notificationStore', () => {
  const userStore = useUserStore()

  const noticeList = ref<NoticeItem[]>([])
  const msgList = ref<MessageItem[]>([])

  const userId = () => userStore.info?.userId ?? userStore.info?.userName ?? 'guest'

  /** 按时间倒序（最新在前） */
  function sortByTimeDesc<T extends { time: string }>(list: T[]): T[] {
    return [...list].sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime())
  }

  const sortedNoticeList = computed(() => sortByTimeDesc(noticeList.value))
  const sortedMsgList = computed(() => sortByTimeDesc(msgList.value))

  function persist() {
    saveNotifications(userId(), {
      notice: noticeList.value,
      message: msgList.value
    })
  }

  function load() {
    const data = loadNotifications(userId())
    noticeList.value = data.notice
    msgList.value = data.message
  }

  /** 根据 SSE tag 推断展示类型（图标） */
  function sseTagToNoticeType(tag: string): string {
    const map: Record<string, string> = {
      greeting: 'message',
      message: 'notice',
      notify: 'notice',
      warning: 'warning'
    }
    return map[tag] ?? 'notice'
  }

  /** 处理 SSE 新消息，持久化并加入列表 */
  function addSseMessage(tag: string, msg: string) {
    const tab = SSE_TAG_TO_TAB[tag] ?? 'notice'
    const id = genNotificationId()
    const time = formatNoticeTime()

    if (tab === 'notice') {
      const item: NoticeItem = {
        id,
        title: msg,
        time,
        type: sseTagToNoticeType(tag),
        read: false,
        tag
      }
      noticeList.value = trimList([...noticeList.value, item], MAX_ITEMS)
    } else {
      const item: MessageItem = {
        id,
        title: msg,
        time,
        avatar: pickAvatar(noticeList.value.length + msgList.value.length),
        read: false
      }
      msgList.value = trimList([...msgList.value, item], MAX_ITEMS)
    }
    persist()
  }

  /** 当前 tab 全部标为已读 */
  function markCurrentTabRead(tabIndex: number) {
    if (tabIndex === 0) {
      noticeList.value = noticeList.value.map((n) => ({ ...n, read: true }))
    } else if (tabIndex === 1) {
      msgList.value = msgList.value.map((m) => ({ ...m, read: true }))
    }
    persist()
  }

  /** 全部标为已读 */
  function markAllRead() {
    noticeList.value = noticeList.value.map((n) => ({ ...n, read: true }))
    msgList.value = msgList.value.map((m) => ({ ...m, read: true }))
    persist()
  }

  /** 标记单条通知为已读 */
  function markNoticeRead(id: string) {
    const idx = noticeList.value.findIndex((n) => n.id === id)
    if (idx >= 0 && !noticeList.value[idx].read) {
      noticeList.value = noticeList.value.map((n) => (n.id === id ? { ...n, read: true } : n))
      persist()
    }
  }

  /** 标记单条消息为已读 */
  function markMessageRead(id: string) {
    const idx = msgList.value.findIndex((m) => m.id === id)
    if (idx >= 0 && !msgList.value[idx].read) {
      msgList.value = msgList.value.map((m) => (m.id === id ? { ...m, read: true } : m))
      persist()
    }
  }

  /** 删除单条通知 */
  function deleteNotice(id: string) {
    noticeList.value = noticeList.value.filter((n) => n.id !== id)
    persist()
  }

  /** 删除单条消息 */
  function deleteMessage(id: string) {
    msgList.value = msgList.value.filter((m) => m.id !== id)
    persist()
  }

  /** 清空当前选中 tab 的消息（0=通知，1=消息） */
  function clearCurrentTab(tabIndex: number) {
    if (tabIndex === 0) {
      noticeList.value = []
    } else if (tabIndex === 1) {
      msgList.value = []
    }
    persist()
  }

  let inited = false
  function onSseMessage(payload: { tag: string; msg: string }) {
    if (payload) addSseMessage(payload.tag, payload.msg)
  }

  /** 初始化：加载持久化数据并监听 SSE */
  function init() {
    load()
    if (inited) return
    inited = true
    mittBus.on('sseMessage', onSseMessage)
  }

  /** 用户切换时重新加载 */
  watch(
    () => userId(),
    () => load(),
    { immediate: false }
  )

  return {
    noticeList,
    msgList,
    sortedNoticeList,
    sortedMsgList,
    load,
    persist,
    addSseMessage,
    markCurrentTabRead,
    markAllRead,
    markNoticeRead,
    markMessageRead,
    deleteNotice,
    deleteMessage,
    clearCurrentTab,
    init
  }
})
