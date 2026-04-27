<!-- 需求思维导图 dev_project mind_json — 标记：多标签 + 库内置图标 -->
<template>
  <div class="art-full-height dev-project-mind flex min-h-[480px] flex-col">
    <div
      class="flex flex-wrap items-center gap-2 border-b border-[var(--el-border-color-lighter)] px-4 py-2"
    >
      <ElButton text @click="goBack">
        <ArtSvgIcon icon="ri:arrow-left-line" class="mr-1" />
        {{ $t('dev.process.back') }}
      </ElButton>
      <span class="text-base font-medium text-[var(--el-text-color-primary)]">{{
        projectName || $t('menus.devops.projectMind')
      }}</span>
      <div class="ml-auto flex flex-wrap items-center gap-2">
        <ElButton
          v-auth="'dev:project:edit'"
          type="primary"
          :loading="saving"
          :disabled="!isMindDirty || saving"
          @click="saveMind(false)"
        >
          {{ $t('dev.project.mindSave') }}
        </ElButton>
        <ElButton v-auth="'dev:project:edit'" @click="exportPng">{{
          $t('dev.project.mindExportPng')
        }}</ElButton>
        <ElButton v-auth="'dev:project:edit'" @click="exportPdf">{{
          $t('dev.project.mindExportPdf')
        }}</ElButton>
      </div>
    </div>
    <div class="flex min-h-0 flex-1">
      <div class="mind-stage relative min-h-0 min-w-0 flex-1 bg-[var(--el-fill-color-lighter)]">
        <div
          v-if="activeNode"
          class="mind-node-info-panel absolute left-5 top-5 z-20 flex w-[260px] max-w-[calc(100%-2.5rem)] flex-col overflow-hidden rounded-2xl border border-white/65 bg-[color:color-mix(in_srgb,var(--el-bg-color)_88%,white_12%)] shadow-[0_1px_2px_rgba(15,23,42,0.05)] backdrop-blur-[10px]"
        >
          <div
            class="border-b border-[color:color-mix(in_srgb,var(--el-border-color-lighter)_70%,transparent)] px-4 py-3"
          >
            <p class="text-sm font-semibold text-[var(--el-text-color-primary)]">
              {{ $t('dev.project.mindTaskMark') }}
            </p>
          </div>
          <ElScrollbar max-height="calc(100vh - 220px)" class="min-h-0 flex-1">
            <div class="space-y-4 px-4 py-3 text-xs leading-6">
              <div>
                <p class="mb-1 font-medium text-[var(--el-text-color-secondary)]">
                  {{ $t('dev.project.mindNodeCanvasLabel') }}
                </p>
                <div class="whitespace-pre-wrap break-words text-[var(--el-text-color-primary)]">
                  {{ activeNodeText || '-' }}
                </div>
              </div>
              <div>
                <p class="mb-1 font-medium text-[var(--el-text-color-secondary)]">
                  {{ $t('dev.project.mindNodeDetailDesc') }}
                </p>
                <div
                  v-if="activeNodeDetailDesc"
                  class="whitespace-pre-wrap break-words text-[var(--el-text-color-regular)]"
                >
                  {{ activeNodeDetailDesc }}
                </div>
                <p v-else class="text-[var(--el-text-color-placeholder)]">
                  {{ $t('dev.project.mindNodeDescEmpty') }}
                </p>
              </div>
            </div>
          </ElScrollbar>
        </div>
        <div ref="containerRef" class="h-full min-h-0 min-w-0" />
      </div>
      <div
        class="flex min-h-0 w-[280px] shrink-0 flex-col border-l border-[var(--el-border-color-lighter)] bg-[var(--el-bg-color)] p-2.5"
      >
        <template v-if="!activeNode">
          <p class="mb-1 shrink-0 text-sm font-medium text-[var(--el-text-color-primary)]">
            {{ $t('dev.project.mindSidebarMarkersTab') }}
          </p>
          <p class="mb-2.5 shrink-0 text-xs leading-relaxed text-[var(--el-text-color-secondary)]">
            {{ $t('dev.project.mindNoSelection') }}
          </p>
        </template>
        <ElTabs
          v-else
          v-model="sidebarTab"
          class="mind-sidebar-tabs flex min-h-0 min-w-0 flex-1 flex-col"
        >
          <ElTabPane v-if="false" :label="$t('dev.project.mindTaskMark')" name="info">
            <div
              class="max-h-[min(52vh,380px)] space-y-3 overflow-y-auto pr-0.5 text-xs leading-relaxed"
            >
              <div>
                <p class="mb-1 font-medium text-[var(--el-text-color-secondary)]">
                  {{ $t('dev.project.mindNodeCanvasLabel') }}
                </p>
                <div class="whitespace-pre-wrap break-words text-[var(--el-text-color-secondary)]">
                  {{ activeNodeText || '—' }}
                </div>
              </div>
              <div>
                <p class="mb-1 font-medium text-[var(--el-text-color-secondary)]">
                  {{ $t('dev.project.mindNodeDetailDesc') }}
                </p>
                <div
                  v-if="activeNodeDetailDesc"
                  class="whitespace-pre-wrap break-words text-[var(--el-text-color-regular)]"
                >
                  {{ activeNodeDetailDesc }}
                </div>
                <p v-else class="text-[var(--el-text-color-placeholder)]">
                  {{ $t('dev.project.mindNodeDescEmpty') }}
                </p>
              </div>
            </div>
          </ElTabPane>
          <ElTabPane :label="$t('dev.project.mindSidebarMarkersTab')" name="markers">
            <div class="flex min-h-0 flex-1 flex-col gap-2 overflow-hidden pt-1">
              <ElScrollbar class="min-h-0 flex-1 pr-0.5">
                <div v-for="section in markerSections" :key="section.key" class="mb-3 last:mb-0">
                  <p class="mb-1 text-xs font-medium text-[var(--el-text-color-secondary)]">
                    {{ $t(section.titleKey) }}
                  </p>
                  <div class="flex flex-wrap gap-1">
                    <ElTooltip
                      v-for="item in section.items"
                      :key="item.key"
                      :content="item.tooltip"
                      placement="top"
                    >
                      <button
                        type="button"
                        class="mind-marker-tile flex h-7 w-7 items-center justify-center rounded border border-[var(--el-border-color)] bg-[var(--el-fill-color-blank)] p-0.5 transition hover:border-[var(--el-color-primary)] disabled:cursor-not-allowed disabled:opacity-50"
                        :class="{
                          '!border-[var(--el-color-primary)] ring-1 ring-[var(--el-color-primary)]/30':
                            selectedIconKeys.includes(item.key)
                        }"
                        :disabled="!canEdit"
                        @click="toggleIconKey(item.key)"
                      >
                        <span class="marker-svg" v-html="item.svg" />
                      </button>
                    </ElTooltip>
                  </div>
                </div>
              </ElScrollbar>
              <ElButton
                class="w-full shrink-0"
                size="small"
                :disabled="!canEdit"
                @click="clearAllMarkers"
              >
                {{ $t('dev.project.mindTaskClear') }}
              </ElButton>
            </div>
          </ElTabPane>
          <ElTabPane :label="$t('dev.project.mindSidebarDescTab')" name="desc">
            <div class="flex min-h-0 flex-1 flex-col pt-1">
              <ElInput
                v-model="descDraft"
                type="textarea"
                :autosize="{ minRows: 10, maxRows: 22 }"
                :disabled="!canEdit"
                :placeholder="$t('dev.project.mindNodeDescPlaceholder')"
                class="mind-desc-textarea"
                @input="scheduleDescCommit"
                @blur="commitNodeDesc"
              />
            </div>
          </ElTabPane>
        </ElTabs>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import { useRoute, useRouter } from 'vue-router'
  import MindMap from 'simple-mind-map'
  import Export from 'simple-mind-map/src/plugins/Export.js'
  import ExportPDF from 'simple-mind-map/src/plugins/ExportPDF.js'
  import iconsSvg from 'simple-mind-map/src/svg/icons.js'
  import { mergerIconList, walk } from 'simple-mind-map/src/utils/index.js'
  import 'simple-mind-map/dist/simpleMindMap.esm.css'
  import { fetchDevProjectDetail, fetchDevProjectEditMind } from '@/api/dev-project'
  import { ElMessage } from 'element-plus'
  import ArtSvgIcon from '@/components/core/base/art-svg-icon/index.vue'
  import { useUserStore } from '@/store/modules/user'
  import {
    buildColoredSignIconGroup,
    colorDotIconList,
    mindMapTaskIconList
  } from './mind-task-icons'

  type MindMapRendererApi = {
    setNodeDataRender: (n: unknown, d: Record<string, unknown>, notRender?: boolean) => void
  }

  /** 自定义节点字段：写入 mind_json，不参与画布节点正文渲染 */
  ;(
    MindMap as unknown as { extendNodeDataNoStylePropList: (list: string[]) => void }
  ).extendNodeDataNoStylePropList(['projectNodeDesc'])

  const coloredSignIconGroup = buildColoredSignIconGroup()

  defineOptions({ name: 'DevopsProjectMind' })

  const { t } = useI18n()
  const route = useRoute()
  const router = useRouter()
  const userStore = useUserStore()

  const canEdit = computed(() => (userStore.info?.buttons ?? []).includes('dev:project:edit'))

  const projectId = computed(() => String(route.params.id ?? ''))

  const isMindDirty = computed(() => {
    mindDataRevision.value
    const mm = mindMap.value
    if (!mm || !canEdit.value) return false
    const snap = persistedMindJsonSnapshot.value
    if (!snap) return false
    try {
      return JSON.stringify(mm.getData(true)) !== snap
    } catch {
      return true
    }
  })

  const containerRef = ref<HTMLDivElement | null>(null)
  const mindMap = shallowRef<InstanceType<typeof MindMap> | null>(null)
  const projectName = ref('')
  const saving = ref(false)
  /** 与上次成功持久化（含自动保存）一致的 getData(true) 快照，用于启用/禁用保存按钮 */
  const persistedMindJsonSnapshot = ref('')
  /** simple-mind-map 非响应式，依赖此计数在 data_change 时触发脏检测 */
  const mindDataRevision = ref(0)
  const activeNode = ref<any>(null)
  /** 右侧栏 Tab：节点信息 | 标签（图标）| 描述 */
  const sidebarTab = ref<'info' | 'markers' | 'desc'>('markers')
  /** 侧栏「描述」编辑草稿，与节点 data.projectNodeDesc 同步 */
  const descDraft = ref('')
  let descCommitTimer: ReturnType<typeof setTimeout> | null = null
  let saveTimer: ReturnType<typeof setTimeout> | null = null
  let structureRollupTimer: ReturnType<typeof setTimeout> | null = null
  let taskIconRollupTimer: ReturnType<typeof setTimeout> | null = null
  /** 用于检测新增/删除节点等结构变化，避免仅改文案时重复全树汇总 */
  let lastTreeStructureFingerprint = ''
  /** 各节点 icon 数组快照，用于画布侧移除/修改标记（不经侧栏 applyNodeIcons）时仍能触发父链重算 */
  let lastTaskIconStateFingerprint = ''

  const mergedIconList = mergerIconList([
    ...iconsSvg.nodeIconList,
    colorDotIconList,
    ...mindMapTaskIconList,
    ...(coloredSignIconGroup ? [coloredSignIconGroup] : [])
  ])

  type MarkerItem = { key: string; svg: string; tooltip: string }
  type MarkerSection = { key: string; titleKey: string; items: MarkerItem[] }

  /**
   * 与 simple-mind-map `progress` 分组一致：progress_1～7 为圆饼 1/8～7/8，progress_8 为对勾 100%。
   */
  function progressTooltip(name: string): string {
    const map: Record<string, string> = {
      '1': '12.5%',
      '2': '25%',
      '3': '37.5%',
      '4': '50%',
      '5': '62.5%',
      '6': '75%',
      '7': '87.5%',
      '8': t('dev.project.mindTaskDone')
    }
    return map[name] ?? `progress_${name}`
  }

  const markerSections = computed((): MarkerSection[] => {
    const find = (type: string) => mergedIconList.find((g: { type: string }) => g.type === type)
    const sections: MarkerSection[] = []

    const color = find('color')
    if (color) {
      sections.push({
        key: 'color',
        titleKey: 'dev.project.mindMarkerTag',
        items: color.list.map((item: { name: string; icon: string }) => ({
          key: `color_${item.name}`,
          svg: item.icon,
          tooltip: `color_${item.name}`
        }))
      })
    }

    const priority = find('priority')
    if (priority) {
      sections.push({
        key: 'priority',
        titleKey: 'dev.project.mindMarkerPriority',
        items: priority.list
          .filter((item: { name: string }) => Number(item.name) <= 7)
          .map((item: { name: string; icon: string }) => ({
            key: `priority_${item.name}`,
            svg: item.icon,
            tooltip: `P${item.name}`
          }))
      })
    }

    const taskG = find('task')
    const progressG = find('progress')
    const taskItems: MarkerItem[] = []
    const noneItem = taskG?.list.find((i: { name: string }) => i.name === 'none')
    if (noneItem) {
      taskItems.push({
        key: 'task_none',
        svg: noneItem.icon,
        tooltip: t('dev.project.mindTaskNotStarted')
      })
    }
    if (progressG) {
      progressG.list.forEach((item: { name: string; icon: string }) => {
        taskItems.push({
          key: `progress_${item.name}`,
          svg: item.icon,
          tooltip: progressTooltip(item.name)
        })
      })
    }
    if (taskItems.length) {
      sections.push({
        key: 'task',
        titleKey: 'dev.project.mindMarkerTask',
        items: taskItems
      })
    }

    const expression = find('expression')
    if (expression) {
      sections.push({
        key: 'expression',
        titleKey: 'dev.project.mindMarkerExpression',
        items: expression.list.map((item: { name: string; icon: string }) => ({
          key: `expression_${item.name}`,
          svg: item.icon,
          tooltip: `expression_${item.name}`
        }))
      })
    }

    const sign = find('sign')
    if (sign) {
      const slice = (from: number, to: number, sk: string, titleKey: string): MarkerSection => ({
        key: sk,
        titleKey,
        items: sign.list
          .filter((item: { name: string }) => {
            const n = Number(item.name)
            return n >= from && n <= to
          })
          .map((item: { name: string; icon: string }) => ({
            key: `sign_${item.name}`,
            svg: item.icon,
            tooltip: `sign_${item.name}`
          }))
      })
      sections.push(slice(1, 7, 'signA', 'dev.project.mindMarkerSignA'))
      sections.push(slice(8, 14, 'signB', 'dev.project.mindMarkerSignB'))
      sections.push(slice(15, 23, 'signC', 'dev.project.mindMarkerSignC'))
    }

    return sections
  })

  watch(activeNode, (node, prev) => {
    if (descCommitTimer) {
      clearTimeout(descCommitTimer)
      descCommitTimer = null
    }
    const mm = mindMap.value
    if (prev && mm && canEdit.value) {
      const stored = String(prev.nodeData?.data?.projectNodeDesc ?? '')
      if (descDraft.value !== stored) {
        const renderer = (mm as { renderer: MindMapRendererApi }).renderer
        renderer.setNodeDataRender(prev, { projectNodeDesc: descDraft.value }, true)
      }
    }
    sidebarTab.value = 'markers'
    descDraft.value = String(node?.nodeData?.data?.projectNodeDesc ?? '')
  })

  const selectedIconKeys = computed(() => {
    const icons = activeNode.value?.nodeData?.data?.icon
    return Array.isArray(icons) ? icons : []
  })

  const activeNodeText = computed(() => {
    const n = activeNode.value
    if (!n?.nodeData?.data?.text) return ''
    let raw = String(n.nodeData.data.text)
    raw = raw
      .replace(/<br\s*\/?>/gi, '\n')
      .replace(/<\/p>\s*<p[^>]*>/gi, '\n')
      .replace(/<\/div>\s*<div[^>]*>/gi, '\n')
    raw = raw.replace(/<[^>]+>/g, '')
    return raw.trimEnd().slice(0, 2000)
  })

  /** 节点扩展描述（仅存于数据，画布不展示） */
  const activeNodeDetailDesc = computed(() => {
    const n = activeNode.value
    const v = n?.nodeData?.data?.projectNodeDesc
    if (v == null || String(v).trim() === '') return ''
    return String(v).trimEnd().slice(0, 8000)
  })

  function commitNodeDesc() {
    const mm = mindMap.value
    const node = activeNode.value
    if (!mm || !node || !canEdit.value) return
    const stored = String(node.nodeData?.data?.projectNodeDesc ?? '')
    if (descDraft.value === stored) return
    const renderer = (mm as { renderer: MindMapRendererApi }).renderer
    renderer.setNodeDataRender(node, { projectNodeDesc: descDraft.value }, true)
  }

  function scheduleDescCommit() {
    if (!canEdit.value) return
    if (descCommitTimer) clearTimeout(descCommitTimer)
    descCommitTimer = setTimeout(() => {
      descCommitTimer = null
      commitNodeDesc()
    }, 500)
  }

  /** 与右侧侧栏分组顺序一致：标签 → 优先级 → 任务 → 表情 → 标记（sign 按编号） */
  function iconOrderRank(key: string): number {
    if (key.startsWith('color_')) return 1_000_000 + numSuffix(key, 'color_')
    if (key.startsWith('priority_')) return 2_000_000 + numSuffix(key, 'priority_')
    if (key === 'task_none') return 3_000_000
    if (key.startsWith('progress_')) return 3_000_100 + numSuffix(key, 'progress_')
    if (key.startsWith('expression_')) return 4_000_000 + numSuffix(key, 'expression_')
    if (key.startsWith('sign_')) return 5_000_000 + numSuffix(key, 'sign_')
    return 9_000_000
  }

  function numSuffix(key: string, prefix: string): number {
    if (!key.startsWith(prefix)) return 0
    const n = Number(key.slice(prefix.length))
    return Number.isFinite(n) ? n : 0
  }

  function sortIconsBySidebarOrder(icons: string[]): string[] {
    if (icons.length <= 1) return [...icons]
    return [...icons].sort((a, b) => iconOrderRank(a) - iconOrderRank(b))
  }

  /** 加载历史数据后统一顺序，避免旧 JSON 中 icon 数组顺序与侧栏不一致 */
  function normalizeAllIconsOrder(mm: any) {
    const root = mm.renderer?.root
    if (!root) return
    const updates: { node: unknown; icon: string[] }[] = []
    walk(root, null, (node: any) => {
      const icons = node.nodeData?.data?.icon
      if (!Array.isArray(icons) || icons.length <= 1) return
      const sorted = sortIconsBySidebarOrder(icons)
      if (sorted.join('\0') !== icons.join('\0')) {
        updates.push({ node, icon: sorted })
      }
    })
    updates.forEach(({ node, icon }) => {
      mm.renderer.setNodeDataRender(node, { icon }, true)
    })
    if (updates.length) mm.render()
  }

  function deriveProgressFromIcons(icons: string[]): number | undefined {
    const hit = icons.find((k) => k === 'task_none' || k.startsWith('progress_'))
    if (!hit) return undefined
    if (hit === 'task_none') return 0
    /** 与库 SVG 扇形一致：每档 12.5%（1/8 圆） */
    const map: Record<string, number> = {
      progress_1: 12.5,
      progress_2: 25,
      progress_3: 37.5,
      progress_4: 50,
      progress_5: 62.5,
      progress_6: 75,
      progress_7: 87.5,
      progress_8: 100
    }
    return map[hit] ?? 0
  }

  function hasTaskIconInList(icons: string[] | undefined): boolean {
    if (!Array.isArray(icons)) return false
    return icons.some((k) => k === 'task_none' || k.startsWith('progress_'))
  }

  /**
   * 向上汇总父级时：仅「已完成」progress_8 计 100%；
   * 非已完成（含已启动、中间进度、无任务图标）均按 0% 计，与「已启动」一致。
   * nd 为树数据片段：与 renderer 上 node.nodeData 同形 { data, children }。
   */
  function getRollupPercentFromNodeDataSlice(nd: any): number {
    const icons = nd?.data?.icon
    if (!Array.isArray(icons)) return 0
    const hit = icons.find((k) => k === 'task_none' || k.startsWith('progress_'))
    if (!hit) return 0
    if (hit === 'progress_8') return 100
    return 0
  }

  /** 平均进度映射为任务图标（与 deriveProgressFromIcons 的 1/8 档一致） */
  function percentToTaskIconKey(avg: number): string {
    const a = Math.max(0, Math.min(100, avg))
    if (a <= 0) return 'task_none'
    if (a >= 100) return 'progress_8'
    const centers: Array<{ pct: number; key: string }> = [
      { pct: 0, key: 'task_none' },
      { pct: 12.5, key: 'progress_1' },
      { pct: 25, key: 'progress_2' },
      { pct: 37.5, key: 'progress_3' },
      { pct: 50, key: 'progress_4' },
      { pct: 62.5, key: 'progress_5' },
      { pct: 75, key: 'progress_6' },
      { pct: 87.5, key: 'progress_7' },
      { pct: 100, key: 'progress_8' }
    ]
    let bestKey = 'task_none'
    let bestD = Infinity
    for (const { pct, key } of centers) {
      const d = Math.abs(a - pct)
      if (d < bestD) {
        bestD = d
        bestKey = key
      }
    }
    return bestKey
  }

  function stripTaskKeys(icons: string[]): string[] {
    return icons.filter((k) => k !== 'task_none' && !k.startsWith('progress_'))
  }

  function mergeTaskKeyIntoIcons(icons: string[] | undefined, taskKey: string): string[] {
    const base = stripTaskKeys(Array.isArray(icons) ? [...icons] : [])
    base.push(taskKey)
    return sortIconsBySidebarOrder(base)
  }

  /** 树结构指纹：仅 uid 与前序，用于判断增删节点（与图标/文本变更区分） */
  function getTreeStructureFingerprint(mm: any): string {
    try {
      const d = mm.getData(true) as { root?: unknown }
      const root = d?.root as { data?: { uid?: string }; children?: unknown[] } | undefined
      if (!root) return ''
      const parts: string[] = []
      function walkData(node: typeof root) {
        if (!node?.data) return
        parts.push(String(node.data.uid ?? ''))
        const ch = node.children
        if (Array.isArray(ch)) {
          for (const c of ch) walkData(c as typeof root)
        }
      }
      walkData(root)
      return `${parts.length}:${parts.join('\0')}`
    } catch {
      return ''
    }
  }

  /** 全树各节点 icon 列表指纹（含 uid），用于检测 SET_NODE_ICON 等导致的变更 */
  function getTaskIconStateFingerprint(mm: any): string {
    try {
      const d = mm.getData(true) as { root?: unknown }
      const root = d?.root as
        | { data?: { uid?: string; icon?: unknown }; children?: unknown[] }
        | undefined
      if (!root) return ''
      const parts: string[] = []
      function walkData(node: typeof root) {
        if (!node?.data) return
        const uid = String(node.data.uid ?? '')
        const icons = node.data.icon
        const iconStr = Array.isArray(icons) ? icons.slice().sort().join('\0') : ''
        parts.push(`${uid}:${iconStr}`)
        const ch = node.children
        if (Array.isArray(ch)) {
          for (const c of ch) walkData(c as typeof root)
        }
      }
      walkData(root)
      return parts.join('\n')
    } catch {
      return ''
    }
  }

  /** 根据直接子节点平均进度，重写父节点任务图标（保留其它标签类图标） */
  function rollupParentTaskFromChildren(mm: any, parent: any, notRender = false) {
    /** 分母以 nodeData.children 为准；仅「带有任务类图标」的子节点参与汇总，去掉任务图标的分支不计入（如 1/2→仅 1 支有任务时按 1/1 算） */
    const rawChildren = parent?.nodeData?.children
    if (!Array.isArray(rawChildren) || rawChildren.length === 0) return
    const participating = rawChildren.filter((nd: any) => hasTaskIconInList(nd?.data?.icon))
    const raw = parent.nodeData?.data?.icon
    const cur = Array.isArray(raw) ? [...raw] : []

    if (participating.length === 0) {
      if (!hasTaskIconInList(cur)) return
      const next = stripTaskKeys(cur)
      const data: Record<string, unknown> = { icon: next }
      if (next.length === 0) data.progress = 0
      else {
        const p = deriveProgressFromIcons(next)
        data.progress = p !== undefined ? p : 0
      }
      mm.renderer.setNodeDataRender(parent, data, notRender)
      return
    }

    const parts: number[] = []
    for (const nd of participating) {
      parts.push(getRollupPercentFromNodeDataSlice(nd))
    }
    const sum = parts.reduce((a, b) => a + b, 0)
    const avg = sum / parts.length
    const taskKey = percentToTaskIconKey(avg)
    const next = mergeTaskKeyIntoIcons(cur, taskKey)
    const data: Record<string, unknown> = { icon: next }
    const p = deriveProgressFromIcons(next)
    if (p !== undefined) data.progress = p
    else data.progress = 0
    mm.renderer.setNodeDataRender(parent, data, notRender)
  }

  /** 仅当父级或任一子节点涉及任务进度时，才需要按子节点数重算父级（避免无任务分支被写入 task_none） */
  function parentNeedsTaskRollup(parent: any): boolean {
    const rawChildren = parent?.nodeData?.children
    if (!Array.isArray(rawChildren) || rawChildren.length === 0) return false
    if (hasTaskIconInList(parent.nodeData?.data?.icon)) return true
    for (const nd of rawChildren) {
      if (hasTaskIconInList(nd?.data?.icon)) return true
    }
    return false
  }

  /** 结构变化后：对有任务语义的分支重算父级汇总（分母随子节点数变化） */
  function syncAllParentTaskRollups(mm: any) {
    const root = mm.renderer?.root
    if (!root) return
    const parents: any[] = []
    walk(root, null, (node: any) => {
      if (node.nodeData?.children?.length && parentNeedsTaskRollup(node)) parents.push(node)
    })
    parents.forEach((p) => rollupParentTaskFromChildren(mm, p, true))
    mm.render()
  }

  /** 子节点任务变更后，沿父链联动更新各级祖先的任务进度图标 */
  function rollupAncestorsTask(mm: any, modifiedNode: any) {
    if (!mm || !modifiedNode || !canEdit.value) return
    let p = modifiedNode.parent
    while (p) {
      rollupParentTaskFromChildren(mm, p)
      p = p.parent
    }
  }

  function toggleIconKey(iconKey: string) {
    const mm = mindMap.value
    const node = activeNode.value
    if (!mm || !node || !canEdit.value) return
    const raw = node.nodeData.data?.icon
    const cur = Array.isArray(raw) ? [...raw] : []
    let next: string[]
    if (cur.includes(iconKey)) {
      next = cur.filter((k) => k !== iconKey)
    } else {
      next = [...cur]
      if (iconKey.startsWith('priority_')) {
        next = next.filter((k) => !k.startsWith('priority_'))
      } else if (iconKey === 'task_none' || iconKey.startsWith('progress_')) {
        next = next.filter((k) => k !== 'task_none' && !k.startsWith('progress_'))
      } else if (iconKey.startsWith('color_')) {
        next = next.filter((k) => !k.startsWith('color_'))
      } else if (iconKey.startsWith('expression_')) {
        next = next.filter((k) => !k.startsWith('expression_'))
      }
      next.push(iconKey)
    }
    applyNodeIcons(next)
  }

  function clearAllMarkers() {
    const mm = mindMap.value
    const node = activeNode.value
    if (!mm || !node || !canEdit.value) return
    applyNodeIcons([])
  }

  function applyNodeIcons(next: string[]) {
    const mm = mindMap.value
    const node = activeNode.value
    if (!mm || !node) return
    const prevIcons = node.nodeData?.data?.icon
    const sorted = sortIconsBySidebarOrder(next)
    const data: Record<string, unknown> = { icon: sorted }
    if (sorted.length === 0) {
      data.progress = 0
    } else {
      const p = deriveProgressFromIcons(sorted)
      if (p !== undefined) data.progress = p
    }
    /**
     * 任务类图标变更后均沿祖先链重算：含「已启动」「进行中」与「已完成」；
     * 非「已完成」时与「已启动」一致参与父级汇总，不再仅限已启动/已完成才触发。
     * 清除全部任务类图标时亦需重算（见 taskProgressCleared）。
     */
    const taskProgressCleared = hasTaskIconInList(prevIcons) && !hasTaskIconInList(sorted)
    // 必须用 setNodeDataRender：仅 SET_NODE_DATA 不会触发节点重绘，画布上的图标不会立刻更新
    const renderer = (
      mm as { renderer: { setNodeDataRender: (n: unknown, d: Record<string, unknown>) => void } }
    ).renderer
    renderer.setNodeDataRender(node, data)
    if (taskProgressCleared || hasTaskIconInList(sorted)) {
      rollupAncestorsTask(mm, node)
    }
    lastTaskIconStateFingerprint = getTaskIconStateFingerprint(mm)
  }

  function defaultFullData() {
    return {
      layout: 'logicalStructure',
      root: {
        data: {
          text: t('dev.project.mindDefaultRoot'),
          expand: true,
          isActive: false,
          progress: 0
        },
        children: []
      }
    }
  }

  function parseMindJson(raw: string) {
    const fallback = defaultFullData()
    if (!raw || raw.trim() === '' || raw.trim() === '{}') return fallback
    try {
      const j = JSON.parse(raw) as Record<string, unknown>
      if (j && typeof j === 'object' && 'root' in j && j.root) {
        return j as typeof fallback
      }
      if (j && typeof j === 'object' && 'data' in j) {
        return { ...fallback, root: j as (typeof fallback)['root'] }
      }
    } catch {
      /* ignore */
    }
    return fallback
  }

  function goBack() {
    router.push({ name: 'DevopsProjectList' })
  }

  async function saveMind(silent: boolean) {
    const mm = mindMap.value
    if (!mm || !projectId.value) return
    if (!canEdit.value) return
    const json = JSON.stringify(mm.getData(true))
    saving.value = true
    try {
      await fetchDevProjectEditMind({ id: projectId.value, mindJson: json })
      persistedMindJsonSnapshot.value = json
      if (!silent) ElMessage.success(t('dev.project.mindSaved'))
    } catch (e) {
      if (!silent) ElMessage.error(e instanceof Error ? e.message : t('dev.project.loadFailed'))
    } finally {
      saving.value = false
    }
  }

  function scheduleAutoSave() {
    if (!canEdit.value) return
    if (saveTimer) clearTimeout(saveTimer)
    saveTimer = setTimeout(() => {
      saveTimer = null
      saveMind(true)
    }, 2000)
  }

  function scheduleStructureRollup() {
    if (!canEdit.value) return
    if (structureRollupTimer) clearTimeout(structureRollupTimer)
    structureRollupTimer = setTimeout(() => {
      structureRollupTimer = null
      const mm = mindMap.value
      if (!mm || !canEdit.value) return
      const fp = getTreeStructureFingerprint(mm)
      if (fp === lastTreeStructureFingerprint) return
      lastTreeStructureFingerprint = fp
      syncAllParentTaskRollups(mm)
      lastTaskIconStateFingerprint = getTaskIconStateFingerprint(mm)
    }, 280)
  }

  /** 画布上改图标（含移除任务/进度）会走 history 的 data_change，但不经过 applyNodeIcons，需单独比对 icon 指纹 */
  function scheduleTaskIconRollupFromDataChange() {
    if (!canEdit.value) return
    if (taskIconRollupTimer) clearTimeout(taskIconRollupTimer)
    taskIconRollupTimer = setTimeout(() => {
      taskIconRollupTimer = null
      const mm = mindMap.value
      if (!mm || !canEdit.value) return
      const fp = getTaskIconStateFingerprint(mm)
      if (fp === lastTaskIconStateFingerprint) return
      syncAllParentTaskRollups(mm)
      lastTaskIconStateFingerprint = getTaskIconStateFingerprint(mm)
    }, 300)
  }

  async function exportPng() {
    const mm = mindMap.value
    if (!mm) return
    const name = (projectName.value || 'mind-map').replace(/[\\/:*?"<>|]/g, '_')
    await mm.export('png', true, name)
  }

  async function exportPdf() {
    const mm = mindMap.value
    if (!mm) return
    const name = (projectName.value || 'mind-map').replace(/[\\/:*?"<>|]/g, '_')
    await mm.export('pdf', true, name)
  }

  onMounted(async () => {
    MindMap.usePlugin(ExportPDF)
    MindMap.usePlugin(Export)

    if (!projectId.value) {
      ElMessage.error(t('dev.project.mindLoadFailed'))
      return
    }

    let full: ReturnType<typeof defaultFullData>
    try {
      const detail = await fetchDevProjectDetail(projectId.value)
      projectName.value = detail.name ?? ''
      full = parseMindJson(detail.mindJson ?? '{}')
    } catch {
      ElMessage.error(t('dev.project.mindLoadFailed'))
      return
    }

    await nextTick()
    const el = containerRef.value
    if (!el) return

    const mm = new MindMap({
      el,
      readonly: !canEdit.value,
      iconList: mergedIconList,
      themeConfig: {
        iconSize: 16
      },
      data: {
        data: {
          text: t('dev.project.mindDefaultRoot'),
          expand: true,
          isActive: false
        },
        children: []
      }
    } as any)
    mm.setFullData(full as any)
    mindMap.value = mm
    normalizeAllIconsOrder(mm)
    lastTreeStructureFingerprint = getTreeStructureFingerprint(mm)
    lastTaskIconStateFingerprint = getTaskIconStateFingerprint(mm)
    persistedMindJsonSnapshot.value = JSON.stringify(mm.getData(true))

    mm.on('node_active', (node: any) => {
      activeNode.value = node
    })

    mm.on('data_change', () => {
      mindDataRevision.value++
      scheduleAutoSave()
      scheduleStructureRollup()
      scheduleTaskIconRollupFromDataChange()
    })
  })

  onBeforeUnmount(() => {
    if (descCommitTimer) {
      clearTimeout(descCommitTimer)
      descCommitTimer = null
    }
    const mmUnmount = mindMap.value
    const nodeUnmount = activeNode.value
    if (mmUnmount && nodeUnmount && canEdit.value) {
      const stored = String(nodeUnmount.nodeData?.data?.projectNodeDesc ?? '')
      if (descDraft.value !== stored) {
        const renderer = (mmUnmount as { renderer: MindMapRendererApi }).renderer
        renderer.setNodeDataRender(nodeUnmount, { projectNodeDesc: descDraft.value }, true)
      }
    }
    if (saveTimer) {
      clearTimeout(saveTimer)
      saveTimer = null
    }
    if (structureRollupTimer) {
      clearTimeout(structureRollupTimer)
      structureRollupTimer = null
    }
    if (taskIconRollupTimer) {
      clearTimeout(taskIconRollupTimer)
      taskIconRollupTimer = null
    }
    mindMap.value?.destroy()
    mindMap.value = null
  })
</script>

<style scoped>
  .mind-stage {
    isolation: isolate;
  }

  .dev-project-mind :deep(.mindMapContainer) {
    width: 100%;
    height: 100%;
  }

  .mind-node-info-panel :deep(.el-scrollbar__wrap) {
    overscroll-behavior: contain;
  }

  .marker-svg {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 1rem;
    height: 1rem;
  }

  .marker-svg :deep(svg) {
    display: block;
    width: 100%;
    height: 100%;
  }

  .mind-sidebar-tabs :deep(.el-tabs__header) {
    flex-shrink: 0;
    margin: 0 0 0.5rem;
  }

  .mind-sidebar-tabs :deep(.el-tabs__content) {
    display: flex;
    flex: 1;
    flex-direction: column;
    min-height: 0;
    overflow: hidden;
  }

  .mind-sidebar-tabs :deep(.el-tab-pane) {
    display: flex;
    flex: 1;
    flex-direction: column;
    min-height: 0;
    overflow: hidden;
  }

  .mind-desc-textarea {
    width: 100%;
  }

  .mind-desc-textarea :deep(.el-textarea__inner) {
    min-height: 10rem;
  }
</style>
