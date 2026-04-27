import type { Edge, Node } from '@vue-flow/core'
import { normalizeNodeKind } from './automation-kinds'
import type {
  AutomationNodePersistData,
  SshConnectionParams,
  UploadServersParams
} from './automation-node-params'
import { mergeNodeParams, normalizeSshPort } from './automation-node-params'

/** 持久化 JSON 可能把 id 存成 number，边两端必须与节点 id 比较方式一致 */
function graphId(v: unknown): string {
  return v == null ? '' : String(v)
}

/** 从单个节点直接提取可用的 SSH（不含「沿上游继承」） */
export function sshParamsFromNodeData(
  data: AutomationNodePersistData | undefined
): SshConnectionParams | null {
  if (!data?.params) return null
  const k = normalizeNodeKind(String(data.kind))
  if (k === 'ssh_connection') {
    const p = mergeNodeParams('ssh_connection', data.params) as SshConnectionParams
    return {
      ...p,
      port: normalizeSshPort(p.port)
    }
  }
  if (k === 'upload_servers') {
    const p = mergeNodeParams('upload_servers', data.params) as UploadServersParams
    return {
      host: p.host,
      port: normalizeSshPort(p.port),
      username: p.username,
      authType: p.authType,
      password: p.password,
      privateKey: p.privateKey
    }
  }
  return null
}

/**
 * 沿入边回溯，找到第一个能提供 SSH 的节点（SSH连接 或 上传部署）。
 */
export function resolveUpstreamSsh(
  nodeId: string | number,
  nodes: Node[],
  edges: Edge[],
  visited = new Set<string>()
): SshConnectionParams | null {
  const nid = graphId(nodeId)
  if (!nid || visited.has(nid)) return null
  visited.add(nid)
  const preds = edges
    .filter((e) => graphId(e.target) === nid)
    .map((e) => graphId(e.source))
    .sort()
  for (const sid of preds) {
    const n = nodes.find((x) => graphId(x.id) === sid)
    const data = n?.data as AutomationNodePersistData | undefined
    const direct = sshParamsFromNodeData(data)
    if (direct?.host?.trim()) return direct
    const up = resolveUpstreamSsh(sid, nodes, edges, visited)
    if (up?.host?.trim()) return up
  }
  return null
}

const EMPTY_SSH: SshConnectionParams = {
  host: '',
  port: '22',
  username: '',
  authType: 'key',
  password: '',
  privateKey: ''
}

/** 列出指向该节点的每条入边对应的 SSH（多路合并时用于展示；执行时各机依次跑同一脚本） */
export type UpstreamSshBranch = {
  sourceId: string
  ssh: SshConnectionParams
  /** 上游节点显示名 */
  label: string
}

/** 与 Go CollectSshTargetsForRemoteScript 一致：按入边源 id 排序；同一 user@host:port 去重 */
export function listUpstreamSshBranches(
  nodeId: string | number,
  nodes: Node[],
  edges: Edge[]
): UpstreamSshBranch[] {
  const nid = graphId(nodeId)
  const preds = edges
    .filter((e) => graphId(e.target) === nid)
    .map((e) => graphId(e.source))
    .sort()
  const out: UpstreamSshBranch[] = []
  const seen = new Set<string>()
  for (const sid of preds) {
    const n = nodes.find((x) => graphId(x.id) === sid)
    const data = n?.data as AutomationNodePersistData | undefined
    const direct = sshParamsFromNodeData(data)
    let sp: SshConnectionParams | null = direct?.host?.trim() ? direct : null
    if (!sp) {
      sp = resolveUpstreamSsh(sid, nodes, edges, new Set())
    }
    if (sp?.host?.trim()) {
      const key = `${(sp.username ?? '').trim().toLowerCase()}@${(sp.host ?? '').trim().toLowerCase()}:${normalizeSshPort(sp.port)}`
      if (seen.has(key)) continue
      seen.add(key)
      const label = (data?.label ?? '').trim() || sid
      out.push({ sourceId: sid, ssh: sp, label })
    }
  }
  return out
}

/** 画布摘要：取第一个上游 SSH（多机时由 listUpstreamSshBranches 另行展示数量） */
export function resolveEffectiveSshForRemoteScript(
  nodeId: string,
  nodes: Node[],
  edges: Edge[]
): SshConnectionParams {
  const br = listUpstreamSshBranches(nodeId, nodes, edges)
  if (br.length === 0) return EMPTY_SSH
  return br[0].ssh
}

/** 远程下载：SSH 仅来自上游 */
export function resolveEffectiveSshForRemoteDownload(
  nodeId: string,
  nodes: Node[],
  edges: Edge[]
): SshConnectionParams {
  return resolveUpstreamSsh(nodeId, nodes, edges) ?? EMPTY_SSH
}
