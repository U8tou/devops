/// <reference types="vite/client" />

declare module 'nprogress'

declare module 'crypto-js'

declare module 'vue-img-cutter'

declare module 'file-saver'

declare module 'simple-mind-map/src/plugins/Export.js'

declare module 'simple-mind-map/src/plugins/ExportPDF.js'

declare module 'simple-mind-map/src/svg/icons.js' {
  const icons: {
    nodeIconList: Array<{ type: string; list: Array<{ name: string; icon: string }> }>
  }
  export default icons
}

declare module 'simple-mind-map/src/utils/index.js' {
  export function mergerIconList(
    list: Array<{ type: string; name?: string; list: Array<{ name: string; icon: string }> }>
  ): Array<{ type: string; name?: string; list: Array<{ name: string; icon: string }> }>
  export function walk(
    root: unknown,
    parent: unknown,
    beforeCallback?: (node: unknown) => void | boolean
  ): void
}

declare module 'qrcode.vue' {
  export type Level = 'L' | 'M' | 'Q' | 'H'
  export type RenderAs = 'canvas' | 'svg'
  export type GradientType = 'linear' | 'radial'
  export interface ImageSettings {
    src: string
    height: number
    width: number
    excavate: boolean
  }
  export interface QRCodeProps {
    value: string
    size?: number
    level?: Level
    background?: string
    foreground?: string
    renderAs?: RenderAs
  }
  const QrcodeVue: any
  export default QrcodeVue
}

// 全局变量声明
declare const __APP_VERSION__: string // 版本号

interface ImportMetaEnv {
  /** 是否开启站内信（通知中心），默认 true */
  readonly VITE_ENABLE_SITE_MESSAGE?: string
  /** 是否开启聊天室，默认 true */
  readonly VITE_ENABLE_CHAT_ROOM?: string
  /** 与后端 app.encryptKey 一致时启用 API 加解密；留空则明文 */
  readonly VITE_API_ENCRYPT_KEY?: string
}
