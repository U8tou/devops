import type { Edge, Node } from '@vue-flow/core'
import type { AutomationNodePersistData } from './automation-node-params'

/**
 * 节点 data.flowEnabled === false 时，本节点与「沿出边能到达」的所有节点均不执行、画布上显示为灰色。
 * sessionDraft 为侧栏未落盘数据，与 process-edit 中 selectedNode 一致。
 */
export function computeFlowSkippedNodeIds(
  nodes: Node[],
  edges: Edge[],
  sessionDraft: Record<string, AutomationNodePersistData>
): Set<string> {
  const isEnabled = (id: string): boolean => {
    const n = nodes.find((x) => String(x.id) === id)
    if (!n?.data) return true
    const d = sessionDraft[id] ?? (n.data as AutomationNodePersistData)
    return d.flowEnabled !== false
  }

  const rootOff: string[] = []
  for (const n of nodes) {
    const id = String(n.id)
    if (!isEnabled(id)) rootOff.push(id)
  }

  const next = new Map<string, string[]>()
  for (const e of edges) {
    const s = String(e.source)
    const t = String(e.target)
    if (!next.has(s)) next.set(s, [])
    next.get(s)!.push(t)
  }
  return bfsDownstreamFrom(rootOff, next)
}

function bfsDownstreamFrom(roots: string[], adj: Map<string, string[]>): Set<string> {
  const out = new Set<string>()
  const q = [...roots]
  while (q.length) {
    const id = q.shift()!
    if (out.has(id)) continue
    out.add(id)
    for (const t of adj.get(id) ?? []) {
      if (!out.has(t)) q.push(t)
    }
  }
  return out
}
