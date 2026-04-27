import iconsSvg from 'simple-mind-map/src/svg/icons.js'

/**
 * 任务「已启动」（原 task_none）：绿圈 + 三角，与库内置 progress 色系 #12BB37 一致。
 * 节点 data.icon 使用 task_none，经 simple-mind-map iconList 解析为 SVG。
 */
export const TASK_NONE_SVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024"><circle cx="512" cy="512" r="400" fill="#fff" stroke="#12BB37" stroke-width="72"/><path fill="#12BB37" d="M430 380L430 644L690 512z"/></svg>`

/** 标签行七色实心圆点（与库 priority 色值对齐） */
const TAG_DOT_COLORS = ['#E93B30', '#FA8D2E', '#F7E983', '#12BB37', '#2E66FA', '#9B59B6', '#6D768D']

export const colorDotIconList = {
  name: '标签',
  type: 'color',
  list: TAG_DOT_COLORS.map((fill, i) => ({
    name: String(i + 1),
    icon: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024"><circle cx="512" cy="512" r="380" fill="${fill}"/></svg>`
  }))
}

/** 与「标签」七色协调的标记底色（库默认 sign 为 #6D768D 单色，此处按序着色） */
const SIGN_RING_PALETTE = [
  '#E93B30',
  '#FA8D2E',
  '#E6D84A',
  '#12BB37',
  '#2E66FA',
  '#9B59B6',
  '#1ABC9C'
]

/**
 * 覆盖 `type: sign`：将 SVG 中灰色 #6D768D 换为调色板色（含圆环与局部强调路径），供侧栏与画布一致。
 */
export function buildColoredSignIconGroup() {
  const sign = iconsSvg.nodeIconList.find((g: { type: string }) => g.type === 'sign')
  if (!sign?.list?.length) return null
  return {
    name: (sign as { name?: string }).name ?? '标记图标',
    type: 'sign' as const,
    list: sign.list.map((item: { name: string; icon: string }, index: number) => {
      const c = SIGN_RING_PALETTE[index % SIGN_RING_PALETTE.length]
      return {
        name: item.name,
        icon: item.icon.replace(/#6D768D/gi, c)
      }
    })
  }
}

/** 合并进 MindMap 的扩展图标分组（type_name → task_none） */
export const mindMapTaskIconList = [
  {
    name: '任务',
    type: 'task',
    list: [{ name: 'none', icon: TASK_NONE_SVG }]
  }
]
