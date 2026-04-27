import type { Node } from '@vue-flow/core'
import type { AutomationKind } from './automation-kinds'
import { normalizeNodeKind } from './automation-kinds'

/** port 在旧数据或接口中可能是 number，直接 .trim() 会在运行时报错（侧栏/画布摘要等） */
export function normalizeSshPort(p: unknown): string {
  if (p == null || p === '') return '22'
  if (typeof p === 'number' && Number.isFinite(p)) return String(Math.trunc(p))
  if (typeof p === 'string') {
    const t = p.trim()
    return t === '' ? '22' : t
  }
  const s = String(p).trim()
  return s === '' ? '22' : s
}

/** Git 凭证：无 / HTTPS（用户名+密码或 PAT）/ SSH 私钥（与仓库 URL 协议一致） */
export type GitRepoAuthType = 'none' | 'http' | 'ssh_key'

/** 各节点持久化参数（写入 flow JSON） */
export interface GitRepoParams {
  repositoryUrl: string
  branch: string
  /** 相对工作区根目录的检出子目录；留空则直接检出到工作区根目录 */
  checkoutSubdir: string
  shallowClone: boolean
  /** 凭证类型 */
  gitAuthType: GitRepoAuthType
  /** HTTPS：用户名（GitHub PAT 时可填任意占位用户名，或 oauth2） */
  httpUsername: string
  /** HTTPS：密码或个人访问令牌 */
  httpPassword: string
  /** SSH：PEM 私钥（用于 git@host:path 或 ssh://） */
  sshPrivateKey: string
}

/** 执行脚本（合并原：打包构建 / 执行远程脚本 / 启动服务） */
export interface ExecuteScriptParams {
  /** 工作目录（相对仓库根或工作区） */
  cwd: string
  /** Shell 脚本内容（可含构建、远端命令、启动服务等） */
  script: string
}

/** 上传部署：SSH 目标与【SSH 连接】一致（结构化），另含远端目录与产物 glob */
export interface UploadServersParams {
  host: string
  port: string
  username: string
  authType: 'password' | 'key'
  password: string
  privateKey: string
  /** 远端部署目录 */
  remotePath: string
  /** 要同步的本地产物 glob，相对构建目录 */
  artifactGlob: string
}

export interface RemoteDownloadParams {
  /** 远端服务器上的文件或目录绝对路径（日志、备份、产物等，经 SCP/SFTP 拉取） */
  remoteSourcePath: string
  /** 保存到工作区根目录下的相对路径（可含文件名） */
  outputPath: string
  /** 下载完成后是否尝试解压 zip / tar.gz 等 */
  extract: boolean
}

/** SSH 连接（凭证仅存于流程 JSON，请妥善保管） */
export interface SshConnectionParams {
  host: string
  port: string
  username: string
  authType: 'password' | 'key'
  password: string
  /** PEM 私钥内容（key 认证时） */
  privateKey: string
}

/** 在远端通过 SSH 执行脚本；SSH 完全由上游「SSH连接」或「上传部署」提供（见 flow-graph-ssh） */
export interface RemoteSshScriptParams {
  /** 远端工作目录 */
  cwd: string
  /** Shell 脚本 */
  script: string
}

export type NodeParamsByKind = {
  git_repo: GitRepoParams
  ssh_connection: SshConnectionParams
  remote_ssh_script: RemoteSshScriptParams
  execute_script: ExecuteScriptParams
  upload_servers: UploadServersParams
  remote_download: RemoteDownloadParams
}

export type AutomationNodePersistData = {
  kind: AutomationKind | string
  label: string
  /**
   * 为 false 时本节点不执行，且从本节点出边可达的下游均不执行；缺省为 true（保持旧数据兼容）
   */
  flowEnabled?: boolean
  params: NodeParamsByKind[AutomationKind]
}

export function defaultParams<K extends AutomationKind>(kind: K): NodeParamsByKind[K] {
  switch (kind) {
    case 'git_repo':
      return {
        repositoryUrl: 'https://github.com/example/app.git',
        branch: 'main',
        checkoutSubdir: '',
        shallowClone: true,
        gitAuthType: 'none',
        httpUsername: '',
        httpPassword: '',
        sshPrivateKey: ''
      } as NodeParamsByKind[K]
    case 'ssh_connection':
      return {
        host: '192.168.1.10',
        port: '22',
        username: 'deploy',
        authType: 'key',
        password: '',
        privateKey: ''
      } as NodeParamsByKind[K]
    case 'remote_ssh_script':
      return {
        cwd: '/opt/app',
        script: 'set -e\necho "remote script"'
      } as NodeParamsByKind[K]
    case 'execute_script':
      return {
        cwd: '.',
        script: 'set -e\nnpm ci\nnpm run build'
      } as NodeParamsByKind[K]
    case 'upload_servers':
      return {
        host: '192.168.1.10',
        port: '22',
        username: 'deploy',
        authType: 'key',
        password: '',
        privateKey: '',
        remotePath: '/opt/releases',
        artifactGlob: 'dist/**'
      } as NodeParamsByKind[K]
    case 'remote_download':
      return {
        remoteSourcePath: '/var/log/myapp/app.log',
        outputPath: 'artifacts/remote/app.log',
        extract: false
      } as NodeParamsByKind[K]
    default: {
      const _k: never = kind
      return _k
    }
  }
}

/** 将旧版三种节点参数合并为 ExecuteScriptParams */
function migrateLegacyToExecuteScript(
  originalKind: string,
  raw: Record<string, unknown>
): ExecuteScriptParams {
  const base = defaultParams('execute_script')
  if (originalKind === 'build') {
    return {
      cwd: typeof raw.cwd === 'string' ? raw.cwd : '.',
      script: typeof raw.command === 'string' ? raw.command : base.script
    }
  }
  if (originalKind === 'remote_script') {
    return {
      cwd: typeof raw.cwd === 'string' ? raw.cwd : '.',
      script: typeof raw.shellScript === 'string' ? raw.shellScript : base.script
    }
  }
  if (originalKind === 'start_service') {
    const name = typeof raw.serviceName === 'string' ? raw.serviceName.trim() : ''
    const cmd = typeof raw.startCommand === 'string' ? raw.startCommand.trim() : ''
    let script = ''
    if (name && cmd) script = `# ${name}\n${cmd}`
    else if (cmd) script = cmd
    else if (name) script = `# ${name}`
    else script = base.script
    return { cwd: '.', script }
  }
  return base
}

function mergeUploadServersParams(
  base: UploadServersParams,
  raw: Record<string, unknown>
): UploadServersParams {
  const remotePath = typeof raw.remotePath === 'string' ? raw.remotePath : base.remotePath
  const artifactGlob = typeof raw.artifactGlob === 'string' ? raw.artifactGlob : base.artifactGlob
  const authType =
    raw.authType === 'key' || raw.authType === 'password' ? raw.authType : base.authType
  const password = typeof raw.password === 'string' ? raw.password : base.password
  const privateKey = typeof raw.privateKey === 'string' ? raw.privateKey : base.privateKey

  let host = typeof raw.host === 'string' ? raw.host.trim() : ''
  let port =
    typeof raw.port === 'string'
      ? raw.port.trim()
      : typeof raw.port === 'number' && Number.isFinite(raw.port)
        ? String(Math.trunc(raw.port))
        : ''
  let username = typeof raw.username === 'string' ? raw.username.trim() : ''

  const legacyHosts = raw.hosts
  if (!host && typeof legacyHosts === 'string' && legacyHosts.trim()) {
    host =
      legacyHosts
        .split(/\r?\n/)
        .map((l) => l.trim())
        .find(Boolean) ?? ''
  }

  // 旧版：单字段 host 为 user@host:port
  const legacyCombined =
    host.includes('@') &&
    !(typeof raw.port === 'string' && raw.port.trim()) &&
    !(typeof raw.username === 'string' && raw.username.trim())
  if (legacyCombined || (host.includes('@') && !username)) {
    const parsed = sshFromUploadHostString(host)
    if (parsed) {
      host = parsed.host
      port = parsed.port || port || '22'
      username = parsed.username || username
    }
  }

  if (!port.trim()) port = base.port || '22'
  if (!host) host = base.host
  if (!username) username = base.username

  return {
    host,
    port,
    username,
    authType,
    password,
    privateKey,
    remotePath,
    artifactGlob
  }
}

/** 与 devflow ParseUploadHost 行为一致，供 merge 旧数据使用 */
function sshFromUploadHostString(hostRaw: string): SshConnectionParams | null {
  const s = hostRaw.trim()
  if (!s) return null
  const at = s.indexOf('@')
  if (at === -1) {
    const lastColon = s.lastIndexOf(':')
    if (lastColon > 0 && /^\d{1,5}$/.test(s.slice(lastColon + 1))) {
      return {
        username: '',
        host: s.slice(0, lastColon),
        port: s.slice(lastColon + 1),
        authType: 'key',
        password: '',
        privateKey: ''
      }
    }
    return {
      username: '',
      host: s,
      port: '22',
      authType: 'key',
      password: '',
      privateKey: ''
    }
  }
  const user = s.slice(0, at)
  const rest = s.slice(at + 1)
  const lastColon = rest.lastIndexOf(':')
  if (lastColon > 0 && /^\d{1,5}$/.test(rest.slice(lastColon + 1))) {
    return {
      username: user,
      host: rest.slice(0, lastColon),
      port: rest.slice(lastColon + 1),
      authType: 'key',
      password: '',
      privateKey: ''
    }
  }
  return {
    username: user,
    host: rest,
    port: '22',
    authType: 'key',
    password: '',
    privateKey: ''
  }
}

/** 合并加载的旧数据与当前默认值；originalKind 可为旧版 build 等 */
export function mergeNodeParams(
  originalKind: string,
  raw: unknown
): NodeParamsByKind[AutomationKind] {
  const nk = normalizeNodeKind(originalKind)
  if (nk === 'execute_script') {
    if (originalKind === 'execute_script') {
      const base = defaultParams('execute_script')
      if (!raw || typeof raw !== 'object') return base
      const o = raw as Record<string, unknown>
      return {
        cwd: typeof o.cwd === 'string' ? o.cwd : base.cwd,
        script: typeof o.script === 'string' ? o.script : base.script
      }
    }
    if (!raw || typeof raw !== 'object') {
      return migrateLegacyToExecuteScript(originalKind, {})
    }
    return migrateLegacyToExecuteScript(originalKind, raw as Record<string, unknown>)
  }

  const base = defaultParams(nk)
  if (!raw || typeof raw !== 'object') return base
  const o = raw as Record<string, unknown>
  if (nk === 'upload_servers') {
    return mergeUploadServersParams(
      base as UploadServersParams,
      o
    ) as NodeParamsByKind[AutomationKind]
  }
  if (nk === 'ssh_connection') {
    const b = base as SshConnectionParams
    const authType = o.authType === 'key' || o.authType === 'password' ? o.authType : b.authType
    return {
      host: typeof o.host === 'string' ? o.host : b.host,
      port:
        typeof o.port === 'string'
          ? o.port
          : typeof o.port === 'number' && Number.isFinite(o.port)
            ? String(Math.trunc(o.port))
            : b.port,
      username: typeof o.username === 'string' ? o.username : b.username,
      authType,
      password: typeof o.password === 'string' ? o.password : b.password,
      privateKey: typeof o.privateKey === 'string' ? o.privateKey : b.privateKey
    } as NodeParamsByKind[AutomationKind]
  }
  if (nk === 'remote_ssh_script') {
    const b = base as RemoteSshScriptParams
    return {
      cwd: typeof o.cwd === 'string' ? o.cwd : b.cwd,
      script: typeof o.script === 'string' ? o.script : b.script
    } as NodeParamsByKind[AutomationKind]
  }
  if (nk === 'remote_download') {
    const b = base as RemoteDownloadParams
    let remoteSourcePath =
      typeof o.remoteSourcePath === 'string' ? o.remoteSourcePath : b.remoteSourcePath
    if (!remoteSourcePath.trim() && typeof o.downloadUrl === 'string' && o.downloadUrl.trim()) {
      remoteSourcePath = o.downloadUrl.trim()
    }
    return {
      remoteSourcePath,
      outputPath: typeof o.outputPath === 'string' ? o.outputPath : b.outputPath,
      extract: typeof o.extract === 'boolean' ? o.extract : b.extract
    } as NodeParamsByKind[AutomationKind]
  }
  if (nk === 'git_repo') {
    const b = base as GitRepoParams
    const authType =
      o.gitAuthType === 'http' || o.gitAuthType === 'ssh_key' || o.gitAuthType === 'none'
        ? o.gitAuthType
        : b.gitAuthType
    return {
      repositoryUrl: typeof o.repositoryUrl === 'string' ? o.repositoryUrl : b.repositoryUrl,
      branch: typeof o.branch === 'string' ? o.branch : b.branch,
      checkoutSubdir: typeof o.checkoutSubdir === 'string' ? o.checkoutSubdir : b.checkoutSubdir,
      shallowClone: typeof o.shallowClone === 'boolean' ? o.shallowClone : b.shallowClone,
      gitAuthType: authType,
      httpUsername: typeof o.httpUsername === 'string' ? o.httpUsername : b.httpUsername,
      httpPassword: typeof o.httpPassword === 'string' ? o.httpPassword : b.httpPassword,
      sshPrivateKey: typeof o.sshPrivateKey === 'string' ? o.sshPrivateKey : b.sshPrivateKey
    } as NodeParamsByKind[AutomationKind]
  }
  return { ...base, ...o } as NodeParamsByKind[AutomationKind]
}

/** 检出目录归一化，用于比较多个 git_repo 节点是否指向同一工作区路径 */
export function normalizeCheckoutSubdirKey(raw: string): string {
  let s = String(raw ?? '')
    .trim()
    .replace(/\\/g, '/')
  while (s.endsWith('/')) s = s.slice(0, -1)
  if (s === '' || s === '.') return ''
  if (s.startsWith('./')) s = s.slice(2)
  return s
}

/**
 * 若画布上存在多个「拉取 GIT 仓库」节点使用相同检出目录，返回该归一化路径（空字符串表示工作区根）；否则返回 null。
 */
export function findDuplicateGitCheckoutSubdir(nodes: Node[]): string | null {
  const counts = new Map<string, number>()
  for (const n of nodes) {
    const d = n.data as AutomationNodePersistData | undefined
    if (!d || normalizeNodeKind(String(d.kind)) !== 'git_repo') continue
    const p = mergeNodeParams('git_repo', d.params) as GitRepoParams
    const key = normalizeCheckoutSubdirKey(p.checkoutSubdir)
    counts.set(key, (counts.get(key) ?? 0) + 1)
  }
  for (const [key, c] of counts) {
    if (c > 1) return key
  }
  return null
}
