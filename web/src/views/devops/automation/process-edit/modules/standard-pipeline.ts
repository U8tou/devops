import type { Edge, Node } from '@vue-flow/core'
import type { AutomationKind } from './automation-kinds'
import { mergeNodeParams } from './automation-node-params'
import type { AutomationNodePersistData } from './automation-node-params'

/**
 * 标准流水线：拉取代码 → 打包构建 → 上传部署 → 启动服务（最后一环为远程脚本，SSH 沿用上游上传节点）
 */
const STANDARD_PIPELINE_STEPS: Array<{
  kind: AutomationKind
  /** i18n: dev.workflow.standardPipeline.* */
  labelKey: string
  /** 与 kind 对应的初始参数 */
  params: Record<string, unknown>
}> = [
  {
    kind: 'git_repo',
    labelKey: 'standardPipeline.pullCode',
    params: {
      repositoryUrl: 'https://github.com/example/app.git',
      branch: 'main',
      checkoutSubdir: '',
      shallowClone: true
    }
  },
  {
    kind: 'execute_script',
    labelKey: 'standardPipeline.packageBuild',
    params: {
      cwd: '.',
      script: 'npm ci && npm run build'
    }
  },
  {
    kind: 'upload_servers',
    labelKey: 'standardPipeline.uploadDeploy',
    params: {
      host: '192.168.1.10',
      port: '22',
      username: 'deploy',
      authType: 'key',
      password: '',
      privateKey: '',
      remotePath: '/opt/releases',
      artifactGlob: 'dist/**'
    }
  },
  {
    kind: 'remote_ssh_script',
    labelKey: 'standardPipeline.startService',
    params: {
      cwd: '/opt/app',
      script: 'set -e\n# 在此填写服务启动命令（如 systemctl restart myapp）'
    }
  }
]

export function buildStandardPipeline(t: (key: string) => string): {
  nodes: Node[]
  edges: Edge[]
} {
  const base = Date.now()
  const nodes: Node[] = STANDARD_PIPELINE_STEPS.map((step, i) => {
    const data: AutomationNodePersistData = {
      kind: step.kind,
      label: t(`dev.workflow.${step.labelKey}`),
      flowEnabled: true,
      params: mergeNodeParams(step.kind, step.params)
    }
    return {
      id: `n-${base}-${i}`,
      type: 'automation',
      position: { x: 48 + i * 260, y: 96 },
      data
    }
  })

  const edges: Edge[] = []
  for (let i = 0; i < nodes.length - 1; i++) {
    edges.push({
      id: `e-${nodes[i].id}-${nodes[i + 1].id}`,
      source: nodes[i].id,
      target: nodes[i + 1].id,
      updatable: true
    })
  }

  return { nodes, edges }
}
