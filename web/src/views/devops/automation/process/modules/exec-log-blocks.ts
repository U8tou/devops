/**
 * 按节点步骤划分日志块（与后端 `[id] 节点显示名: start (YYYY-MM-DD...)` 对齐；名称可含空格、冒号）
 */
const NODE_START_LINE = /^\[[^\]]+\]\s+.+:\s+start\s+\(\d{4}-\d{2}-\d{2}/

export function splitExecLogIntoBlocks(text: string): string[][] {
  const lines = text.split(/\r?\n/)
  const blocks: string[][] = []
  let cur: string[] = []
  for (const line of lines) {
    if (NODE_START_LINE.test(line)) {
      if (cur.length) blocks.push(cur)
      cur = [line]
    } else if (cur.length === 0) {
      cur = [line]
    } else {
      cur.push(line)
    }
  }
  if (cur.length) blocks.push(cur)
  return blocks
}
