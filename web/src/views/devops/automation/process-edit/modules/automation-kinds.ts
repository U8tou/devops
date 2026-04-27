/** 自动化节点类型（与持久化 JSON 中 data.kind 一致） */
export type AutomationKind =
  | 'git_repo'
  | 'ssh_connection'
  | 'remote_ssh_script'
  | 'execute_script'
  | 'upload_servers'
  | 'remote_download'

/** 旧版已合并为 execute_script，仅用于兼容加载 */
export type LegacyAutomationKind = 'build' | 'remote_script' | 'start_service'

const KNOWN_KINDS: readonly AutomationKind[] = [
  'git_repo',
  'ssh_connection',
  'remote_ssh_script',
  'execute_script',
  'upload_servers',
  'remote_download'
]

export function normalizeNodeKind(kind: string): AutomationKind {
  if (kind === 'build' || kind === 'remote_script' || kind === 'start_service') {
    return 'execute_script'
  }
  if (KNOWN_KINDS.includes(kind as AutomationKind)) {
    return kind as AutomationKind
  }
  return 'git_repo'
}

export interface PaletteItemDef {
  kind: AutomationKind
  /** i18n key: dev.workflow.* */
  titleKey: string
  icon: string
  /** 用于节点左侧色条 */
  accent: string
}

export const PALETTE_ITEMS: PaletteItemDef[] = [
  { kind: 'git_repo', titleKey: 'gitRepo', icon: 'ri:git-branch-line', accent: '#6366f1' },
  {
    kind: 'ssh_connection',
    titleKey: 'sshConnection',
    icon: 'ri:key-2-line',
    accent: '#9333ea'
  },
  {
    kind: 'remote_ssh_script',
    titleKey: 'remoteSshScript',
    icon: 'ri:earth-line',
    accent: '#ea580c'
  },
  {
    kind: 'execute_script',
    titleKey: 'executeScript',
    icon: 'ri:terminal-box-line',
    accent: '#f59e0b'
  },
  {
    kind: 'upload_servers',
    titleKey: 'uploadServers',
    icon: 'ri:upload-cloud-2-line',
    accent: '#0ea5e9'
  },
  {
    kind: 'remote_download',
    titleKey: 'remoteDownload',
    icon: 'ri:download-cloud-2-line',
    accent: '#14b8a6'
  }
]

export function findPaletteDef(kind: string): PaletteItemDef | undefined {
  return PALETTE_ITEMS.find((p) => p.kind === normalizeNodeKind(kind))
}

/** 与 dev.workflow.* 文案键对应 */
export const KIND_TO_TITLE_KEY: Record<AutomationKind, string> = {
  git_repo: 'gitRepo',
  ssh_connection: 'sshConnection',
  remote_ssh_script: 'remoteSshScript',
  execute_script: 'executeScript',
  upload_servers: 'uploadServers',
  remote_download: 'remoteDownload'
}
