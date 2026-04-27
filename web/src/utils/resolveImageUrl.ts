/**
 * 解析后端返回的图片/头像等资源 URL
 * - 若为网络图片（http/https）则原样返回
 * - 若为相对路径则拼接当前页面的域名或 ip:端口（window.location.origin）
 *
 * @module utils/resolveImageUrl
 */

/**
 * 获取当前页面的 origin（协议 + 域名或 ip + 端口）
 * 非浏览器环境（如 SSR）返回空字符串
 */
function getCurrentOrigin(): string {
  if (typeof window === 'undefined' || !window.location) return ''
  return window.location.origin || ''
}

/**
 * 将后端返回的图片/头像地址解析为可访问的完整 URL
 * @param url 后端返回的地址（相对路径如 /api/uploads/avatar/xx.png 或完整 URL）
 * @returns 可用的完整 URL，若 url 为空则返回空字符串
 */
export function resolveImageUrl(url: string | undefined | null): string {
  if (url == null || url === '') return ''

  const trimmed = String(url).trim()
  if (!trimmed) return ''

  // 已是网络图片，不处理
  if (/^https?:\/\//i.test(trimmed)) return trimmed

  // 相对路径：拼接当前域名或 ip:端口
  const base = getCurrentOrigin().replace(/\/$/, '')
  const path = trimmed.startsWith('/') ? trimmed : `/${trimmed}`
  console.log(`🚀 resolveImageUrl: `, `${base}/api${path}`)
  return base ? `${base}/api/uploads${path}` : trimmed
}
