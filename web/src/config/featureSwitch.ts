/**
 * 基于 .env 的功能开关（构建期确定）
 */

/** 解析 VITE 布尔字符串：未配置用 defaultVal；true/1/yes/on 为开；false/0/no/off 为关 */
export function viteEnvBool(raw: string | undefined, defaultVal = true): boolean {
  if (raw === undefined || raw === '') return defaultVal
  const s = String(raw).trim().toLowerCase()
  if (s === 'false' || s === '0' || s === 'no' || s === 'off') return false
  if (s === 'true' || s === '1' || s === 'yes' || s === 'on') return true
  return defaultVal
}

/** 站内信 / 顶部通知中心（SSE 通知抽屉） */
export const siteMessageEnabled = viteEnvBool(import.meta.env.VITE_ENABLE_SITE_MESSAGE, true)

/** IM 聊天室（顶部入口 + 全局聊天窗 WebSocket） */
export const chatRoomEnabled = viteEnvBool(import.meta.env.VITE_ENABLE_CHAT_ROOM, true)
