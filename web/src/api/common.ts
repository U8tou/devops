import request from '@/utils/http'

const { VITE_API_URL } = import.meta.env

/**
 * 获取文件（对应 /common/get_file/{dir}/{obj}）
 * @param dir 文件夹
 * @param obj 文件名
 * @returns 文件内容（Blob，如 application/octet-stream）
 */
export function fetchGetFile(dir: string, obj: string) {
  return request.get<Blob>({
    url: `/api/common/get_file/${encodeURIComponent(dir)}/${encodeURIComponent(obj)}`,
    responseType: 'blob'
  })
}

/**
 * 获取 SSE 连接地址（对应 GET /sys_sse）
 * connId 由后端用分布式 ID 生成，前端无需传。传 token 鉴权：有 token 且有效则为登录用户（sseId=loginId_connId），
 * 可被 SendToUser 定向通知；无 token 或无效为游客（sseId=guest_connId），仅能接收广播。
 * EventSource 不支持自定义 header，鉴权需将 token 传入 URL。
 * 使用方式：new EventSource(getSseUrl(token))
 * @param token 鉴权 token（可选），有则登录用户定向推送，无则游客广播
 * @returns 完整 SSE URL（text/event-stream）
 */
export function getSseUrl(token?: string): string {
  const base = (VITE_API_URL || '').replace(/\/$/, '')
  const url = token
    ? `${base}/api/sys_sse?token=${encodeURIComponent(token)}`
    : `${base}/api/sys_sse`
  return url
}

/**
 * 获取 WebSocket 连接地址（对应 GET /sys_websocket）
 * 建立 WebSocket 全双工连接，需传 uuid 标识；服务端通过 syswebsocket.SendMsg(uuid, tag, msg) 推送；客户端可发文本/二进制，服务端会 echo 文本消息。
 * 使用方式：new WebSocket(getWebSocketUrl(uuid))
 * @param uuid 连接标识，用于后续定向推送
 * @returns 完整 WebSocket URL（ws 或 wss）
 */
export function getWebSocketUrl(uuid: string): string {
  const base = (VITE_API_URL || '').replace(/\/$/, '')
  const wsBase = base.replace(/^http/, 'ws')
  return `${wsBase}/api/sys_websocket?uuid=${encodeURIComponent(uuid)}`
}

/**
 * 获取 IM 聊天室 WebSocket 连接地址（对应 GET /sys_im/chat）
 * 可选传 token 鉴权：有 token 且有效则显示登录用户昵称和头像，无 token 为游客。
 * nickName 为纯昵称（如 小王、游客）；会话唯一性由 loginId、connId（clientId）区分，判断本人消息需用 loginId。
 * 消息格式：客户端→服务端 join(roomId)/leave(roomId)/message(roomId,content)；
 * 服务端→客户端 message、joined、left、error、online(成员上线广播)、offline(成员下线广播)。
 * joined/left/online/offline 推送 memberCount（当前房间在线人数）。
 * 使用方式：new WebSocket(getImChatUrl(token))
 * @param token 鉴权 token（可选），有则显示用户昵称和头像，无则为游客
 * @returns 完整 IM WebSocket URL（ws 或 wss）
 */
export function getImChatUrl(token?: string): string {
  let base = (VITE_API_URL || '').replace(/\/$/, '')
  // 开发环境 VITE_API_URL 常为 "/"，需用当前 origin 构造 ws URL 以正确走代理
  if (!base || base === '/') {
    base = typeof window !== 'undefined' ? window.location.origin : ''
  }
  const wsBase = base.replace(/^https?/, (m: string) => (m === 'https' ? 'wss' : 'ws'))
  return token
    ? `${wsBase}/api/sys_im/chat?token=${encodeURIComponent(token)}`
    : `${wsBase}/api/sys_im/chat`
}
