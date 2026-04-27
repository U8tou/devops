/**
 * 通用方法封装（若依风格）
 * Copyright (c) 2019 ruoyi
 */

/** 日期格式化占位符 */
type DateFormatKey = 'y' | 'm' | 'd' | 'h' | 'i' | 's' | 'a'

/** 数据字典项 */
export interface DictItem {
  value: string
  label: string
}

/** 带 params 的请求参数 */
export interface ParamsWithRange {
  params?: Record<string, unknown>
  [key: string]: unknown
}

/** 树节点（含可选的 children） */
export interface TreeNode extends Record<string, unknown> {
  children?: TreeNode[]
}

/**
 * 日期格式化
 * @param time 日期：Date、时间戳(10/13位)或日期字符串
 * @param pattern 格式，默认 '{y}-{m}-{d} {h}:{i}:{s}'
 */
export function parseTime(time?: Date | string | number | null, pattern?: string): string | null {
  if (time === undefined || time === null) {
    return null
  }
  const format = pattern ?? '{y}-{m}-{d} {h}:{i}:{s}'
  let date: Date
  if (time instanceof Date) {
    date = time
  } else {
    let normalized: string | number = time
    if (typeof time === 'string' && /^[0-9]+$/.test(time)) {
      normalized = parseInt(time, 10)
    } else if (typeof time === 'string') {
      normalized = time
        .replace(/-/gm, '/')
        .replace('T', ' ')
        .replace(/\.[\d]{3}/gm, '')
    }
    if (typeof normalized === 'number' && normalized.toString().length === 10) {
      normalized = normalized * 1000
    }
    date = new Date(normalized)
  }
  const formatObj: Record<DateFormatKey, number> = {
    y: date.getFullYear(),
    m: date.getMonth() + 1,
    d: date.getDate(),
    h: date.getHours(),
    i: date.getMinutes(),
    s: date.getSeconds(),
    a: date.getDay()
  }
  const timeStr = format.replace(/{(y|m|d|h|i|s|a)+}/g, (_match, key: DateFormatKey) => {
    const value = formatObj[key]
    if (key === 'a') {
      return ['日', '一', '二', '三', '四', '五', '六'][value]
    }
    if (value === undefined || value === null) {
      return '00'
    }
    // 年四位；月日时分秒两位补零（修复原先用 Number('0'+n) 再转字符串会丢失前导零的问题）
    if (key === 'y') {
      return String(value).padStart(4, '0')
    }
    return String(value).padStart(2, '0')
  })
  return timeStr
}

/**
 * 表单重置（传入 form 实例）
 * @param formRef 表单 ref，如 ElForm 的 ref 值
 */
export function resetForm(formRef: { resetFields?: () => void } | null | undefined): void {
  formRef?.resetFields?.()
}

/**
 * 添加日期范围到请求参数
 */
export function addDateRange(
  params: ParamsWithRange,
  dateRange?: unknown[] | null,
  propName?: string
): ParamsWithRange {
  const search = { ...params }
  search.params =
    typeof search.params === 'object' && search.params !== null && !Array.isArray(search.params)
      ? { ...search.params }
      : {}
  const range = Array.isArray(dateRange) ? dateRange : []
  if (propName === undefined) {
    ;(search.params as Record<string, unknown>)['beginTime'] = range[0]
    ;(search.params as Record<string, unknown>)['endTime'] = range[1]
  } else {
    ;(search.params as Record<string, unknown>)['begin' + propName] = range[0]
    ;(search.params as Record<string, unknown>)['end' + propName] = range[1]
  }
  return search
}

/**
 * 回显数据字典（单个值）
 */
export function selectDictLabel(
  datas: Record<string, DictItem> | DictItem[],
  value: unknown
): string {
  if (value === undefined) {
    return ''
  }
  const actions: string[] = []
  const list = Array.isArray(datas) ? datas : Object.values(datas)
  const found = list.some((item) => {
    if (String(item.value) === String(value)) {
      actions.push(item.label)
      return true
    }
    return false
  })
  if (!found) {
    actions.push(String(value))
  }
  return actions.join('')
}

/**
 * 回显数据字典（多个值，逗号或指定分隔符）
 */
export function selectDictLabels(
  datas: Record<string, DictItem> | DictItem[],
  value: string | string[] | undefined,
  separator?: string
): string {
  if (value === undefined || (Array.isArray(value) && value.length === 0)) {
    return ''
  }
  const str = Array.isArray(value) ? value.join(',') : value
  const currentSeparator = separator ?? ','
  const temp = str.split(currentSeparator)
  const actions: string[] = []
  for (const item of temp) {
    const trimmed = item.trim()
    if (!trimmed) continue
    let match = false
    const list = Array.isArray(datas) ? datas : Object.values(datas)
    for (const dict of list) {
      if (String(dict.value) === trimmed) {
        actions.push(dict.label + currentSeparator)
        match = true
        break
      }
    }
    if (!match) {
      actions.push(trimmed + currentSeparator)
    }
  }
  const joined = actions.join('')
  return joined.length > 0 ? joined.slice(0, joined.length - 1) : ''
}

/**
 * 字符串格式化（%s 占位）
 */
export function sprintf(str: string, ...args: unknown[]): string {
  let i = 0
  let flag = true
  const result = str.replace(/%s/g, () => {
    const arg = args[i++]
    if (arg === undefined) {
      flag = false
      return ''
    }
    return String(arg)
  })
  return flag ? result : ''
}

/**
 * 将 undefined/null/'undefined'/'null' 转为空字符串
 */
export function parseStrEmpty(str: unknown): string {
  if (str === undefined || str === null || str === 'undefined' || str === 'null') {
    return ''
  }
  return String(str)
}

/**
 * 深度合并对象（递归）
 */
export function mergeRecursive<T extends Record<string, unknown>>(
  source: T,
  target: Record<string, unknown>
): T {
  for (const p of Object.keys(target)) {
    try {
      const val = target[p]
      if (val && typeof val === 'object' && !Array.isArray(val) && val.constructor === Object) {
        ;(source as Record<string, unknown>)[p] = mergeRecursive(
          ((source as Record<string, unknown>)[p] as Record<string, unknown>) ?? {},
          val as Record<string, unknown>
        )
      } else {
        ;(source as Record<string, unknown>)[p] = val
      }
    } catch {
      ;(source as Record<string, unknown>)[p] = target[p]
    }
  }
  return source
}

export interface HandleTreeOptions {
  id?: string
  parentId?: string
  children?: string
}

/**
 * 构造树型结构数据
 * @param data 扁平数据源
 * @param id id 字段名，默认 'id'
 * @param parentId 父节点字段名，默认 'parentId'
 * @param children 子节点字段名，默认 'children'
 */
export function handleTree<T extends Record<string, unknown>>(
  data: T[],
  id?: string,
  parentId?: string,
  children?: string
): (T & { [K: string]: unknown })[] {
  const idKey = id ?? 'id'
  const parentIdKey = parentId ?? 'parentId'
  const childrenKey = children ?? 'children'

  const childrenListMap: Record<string, T[]> = {}
  const nodeIds: Record<string, T> = {}
  const tree: (T & { [K: string]: unknown })[] = []

  for (const d of data) {
    const pid = d[parentIdKey] as string | number | undefined
    const key = String(pid ?? '')
    if (childrenListMap[key] == null) {
      childrenListMap[key] = []
    }
    const idVal = d[idKey]
    if (idVal !== undefined && idVal !== null) {
      nodeIds[String(idVal)] = d
    }
    childrenListMap[key].push(d)
  }

  for (const d of data) {
    const pid = d[parentIdKey] as string | number | undefined
    if (nodeIds[String(pid ?? '')] == null) {
      tree.push(d as T & { [K: string]: unknown })
    }
  }

  function adaptToChildrenList(o: Record<string, unknown>): void {
    const oid = o[idKey]
    const list = oid !== undefined && oid !== null ? childrenListMap[String(oid)] : null
    if (list != null) {
      o[childrenKey] = list as unknown
    }
    const childList = o[childrenKey] as Record<string, unknown>[] | undefined
    if (childList?.length) {
      for (const c of childList) {
        adaptToChildrenList(c)
      }
    }
  }

  for (const t of tree) {
    adaptToChildrenList(t as Record<string, unknown>)
  }
  return tree
}

/**
 * 参数序列化为 query 字符串（不包含前导 ?）
 */
export function tansParams(params: Record<string, unknown>): string {
  let result = ''
  for (const propName of Object.keys(params)) {
    const value = params[propName]
    const part = encodeURIComponent(propName) + '='
    if (value !== null && value !== '' && value !== undefined) {
      if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
        const obj = value as Record<string, unknown>
        for (const key of Object.keys(obj)) {
          if (obj[key] !== null && obj[key] !== '' && obj[key] !== undefined) {
            const paramKey = propName + '[' + key + ']'
            result +=
              encodeURIComponent(paramKey) + '=' + encodeURIComponent(String(obj[key])) + '&'
          }
        }
      } else {
        result += part + encodeURIComponent(String(value)) + '&'
      }
    }
  }
  return result
}

/**
 * 规范化路径：去双斜杠、去末尾斜杠
 */
export function getNormalPath(p: string): string {
  if (!p || p.length === 0 || p === 'undefined') {
    return p
  }
  let res = p.replace(/\/\/+/g, '/')
  if (res.endsWith('/')) {
    res = res.slice(0, res.length - 1)
  }
  return res
}

/**
 * 验证是否为 Blob 且非 JSON（常用于判断是否为文件流）
 */
export function blobValidate(data: Blob): boolean {
  return data.type !== 'application/json'
}
