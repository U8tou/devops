/**
 * HTTP 错误处理模块
 *
 * 提供统一的 HTTP 请求错误处理机制
 *
 * ## 主要功能
 *
 * - 自定义 HttpError 错误类，封装错误信息、状态码、时间戳等
 * - 错误拦截和转换，将 Axios 错误转换为标准的 HttpError
 * - 错误消息国际化处理，根据状态码返回对应的多语言错误提示
 * - 错误日志记录，便于问题追踪和调试
 * - 错误和成功消息的统一展示
 * - 类型守卫函数，用于判断错误类型
 *
 * ## 使用场景
 *
 * - HTTP 请求拦截器中统一处理错误
 * - 业务代码中捕获和处理特定错误
 * - 错误日志收集和上报
 *
 * @module utils/http/error
 * @author Art Design Pro Team
 */
import { unref } from 'vue'
import { AxiosError } from 'axios'
import { ApiStatus } from './status'
import i18n, { $t } from '@/locales'
import { LanguageEnum } from '@/enums/appEnum'

// 错误响应接口
export interface ErrorResponse {
  /** 错误状态码 */
  code: number
  /** 错误消息 */
  msg: string
  /** 错误附加数据 */
  data?: unknown
}

// 错误日志数据接口
export interface ErrorLogData {
  /** 错误状态码 */
  code: number
  /** 错误消息 */
  message: string
  /** 错误附加数据 */
  data?: unknown
  /** 错误发生时间戳 */
  timestamp: string
  /** 请求 URL */
  url?: string
  /** 请求方法 */
  method?: string
  /** 错误堆栈信息 */
  stack?: string
}

// 自定义 HttpError 类
export class HttpError extends Error {
  public readonly code: number
  public readonly data?: unknown
  public readonly timestamp: string
  public readonly url?: string
  public readonly method?: string

  constructor(
    message: string,
    code: number,
    options?: {
      data?: unknown
      url?: string
      method?: string
    }
  ) {
    super(message)
    this.name = 'HttpError'
    this.code = code
    this.data = options?.data
    this.timestamp = new Date().toISOString()
    this.url = options?.url
    this.method = options?.method
  }

  public toLogData(): ErrorLogData {
    return {
      code: this.code,
      message: this.message,
      data: this.data,
      timestamp: this.timestamp,
      url: this.url,
      method: this.method,
      stack: this.stack
    }
  }
}

/**
 * 获取错误消息
 * @param status 错误状态码
 * @returns 错误消息
 */
const getErrorMessage = (status: number): string => {
  const errorMap: Record<number, string> = {
    [ApiStatus.error]: 'httpMsg.badRequest',
    [ApiStatus.unauthorized]: 'httpMsg.unauthorized',
    [ApiStatus.forbidden]: 'httpMsg.forbidden',
    [ApiStatus.notFound]: 'httpMsg.notFound',
    [ApiStatus.methodNotAllowed]: 'httpMsg.methodNotAllowed',
    [ApiStatus.requestTimeout]: 'httpMsg.requestTimeout',
    [ApiStatus.internalServerError]: 'httpMsg.internalServerError',
    [ApiStatus.badGateway]: 'httpMsg.badGateway',
    [ApiStatus.serviceUnavailable]: 'httpMsg.serviceUnavailable',
    [ApiStatus.gatewayTimeout]: 'httpMsg.gatewayTimeout'
  }

  return $t(errorMap[status] || 'httpMsg.internalServerError')
}

/** 后端结构体字段名与界面常用中文名的映射（可按业务逐步补全） */
const VALIDATION_FIELD_LABEL_ZH: Record<string, string> = {
  Address: '地址',
  UserName: '用户名',
  NickName: '昵称',
  Phone: '手机号',
  Email: '邮箱',
  Remark: '备注'
}

/** 将 Go/Gin validator 等后端校验文案转为面向用户的提示（尽量保留可读细节） */
function humanizeValidatorDetail(raw: string): string {
  const match = raw.match(/Field validation for '([^']+)' failed on the '(\w+)' tag/)
  if (match) {
    const fieldKey = match[1]
    const field =
      unref(i18n.global.locale) === LanguageEnum.ZH
        ? (VALIDATION_FIELD_LABEL_ZH[fieldKey] ?? fieldKey)
        : fieldKey
    const tag = match[2]
    if (tag === 'required') {
      // 占位符须避免与 vue-i18n 的 {named} 插值冲突，否则 $t 会把字段名吃掉
      return $t('httpMsg.fieldRequired').replace(/__FIELD__/g, field)
    }
    return $t('httpMsg.fieldValidationRule')
      .replace(/__FIELD__/g, field)
      .replace(/__RULE__/g, tag)
  }
  const stripped = raw.replace(/^Key:\s*'[^']*'\s*Error:\s*/i, '').trim()
  return stripped || raw
}

/**
 * 根据接口返回的 msg / data 拼出展示用错误文案（业务 code 非 200 或 HTTP 4xx 体）
 */
export function formatBusinessErrorMessage(msg: string | undefined, data: unknown): string {
  const detail =
    typeof data === 'string' && data.trim() !== '' ? humanizeValidatorDetail(data.trim()) : ''
  const m = (msg ?? '').trim()
  const genericMsg = /^bad request$/i.test(m) || m === ''

  if (detail) {
    if (genericMsg) {
      return `${$t('httpMsg.badRequest')}：${detail}`
    }
    return m.includes(detail) ? m : `${m}：${detail}`
  }
  if (m) return m
  return $t('httpMsg.badRequest')
}

/**
 * 处理错误
 * @param error 错误对象
 * @returns 错误对象
 */
export function handleError(error: AxiosError<ErrorResponse>): never {
  // 处理取消的请求
  if (error.code === 'ERR_CANCELED') {
    console.warn('Request cancelled:', error.message)
    throw new HttpError($t('httpMsg.requestCancelled'), ApiStatus.error)
  }

  const statusCode = error.response?.status
  const errorMessage = error.response?.data?.msg || error.message
  const requestConfig = error.config
  const resData = error.response?.data as ErrorResponse | undefined

  // 处理网络错误
  if (!error.response) {
    throw new HttpError($t('httpMsg.networkError'), ApiStatus.error, {
      url: requestConfig?.url,
      method: requestConfig?.method?.toUpperCase()
    })
  }

  // HTTP 400：展示后端返回的校验/参数说明，避免被误映射为 500 文案
  if (statusCode === ApiStatus.error && resData) {
    const message = formatBusinessErrorMessage(resData.msg, resData.data)
    throw new HttpError(message, statusCode, {
      data: error.response.data,
      url: requestConfig?.url,
      method: requestConfig?.method?.toUpperCase()
    })
  }

  // 处理 HTTP 状态码错误
  const message = statusCode
    ? getErrorMessage(statusCode)
    : errorMessage || $t('httpMsg.requestFailed')
  throw new HttpError(message, statusCode || ApiStatus.error, {
    data: error.response.data,
    url: requestConfig?.url,
    method: requestConfig?.method?.toUpperCase()
  })
}

/**
 * 显示错误消息
 * @param error 错误对象
 * @param showMessage 是否显示错误消息
 */
export function showError(error: HttpError, showMessage: boolean = true): void {
  if (showMessage) {
    ElMessage.error(error.message)
  }
  // 记录错误日志
  console.error('[HTTP Error]', error.toLogData())
}

/**
 * 显示成功消息
 * @param message 成功消息
 * @param showMessage 是否显示消息
 */
export function showSuccess(message: string, showMessage: boolean = true): void {
  if (showMessage) {
    ElMessage.success(message)
  }
}

/**
 * 判断是否为 HttpError 类型
 * @param error 错误对象
 * @returns 是否为 HttpError 类型
 */
export const isHttpError = (error: unknown): error is HttpError => {
  return error instanceof HttpError
}
