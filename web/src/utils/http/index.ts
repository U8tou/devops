/**
 * HTTP 请求封装模块
 * 基于 Axios 封装的 HTTP 请求工具，提供统一的请求/响应处理
 *
 * ## 主要功能
 *
 * - 请求/响应拦截器（自动添加 Token、统一错误处理）
 * - 401 时若有 refreshToken 则自动刷新短 token 并重试原请求，用户无感知
 * - 刷新失败则退出登录（带防抖机制）
 * - 请求失败自动重试（可配置）
 * - 统一的成功/错误消息提示
 * - 支持 GET/POST/PUT/DELETE 等常用方法
 *
 * @module utils/http
 * @author Art Design Pro Team
 */

import axios, { AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { useUserStore } from '@/store/modules/user'
import { ApiStatus } from './status'
import { HttpError, handleError, showError, showSuccess, formatBusinessErrorMessage } from './error'
import { $t } from '@/locales'
import { BaseResponse } from '@/types'
import {
  API_CRYPTO_HEADER,
  API_CRYPTO_HEADER_VALUE,
  encryptUtf8ToEnvelopeB64,
  isApiCryptoEnabled,
  maybeDecryptResponseBody
} from './apiCrypto'

/** 公开认证接口：请求体需为明文 JSON，供后端 BodyParser 绑定（与加密信封 { c } 互斥） */
const SKIP_ENCRYPT_PATH_MARKERS = [
  'sys_auth/login',
  'sys_auth/register',
  'sys_auth/refresh_token'
] as const

function shouldSkipRequestBodyEncrypt(url: string | undefined): boolean {
  if (!url) return false
  return SKIP_ENCRYPT_PATH_MARKERS.some((m) => url.includes(m))
}

/** 请求配置常量 */
const REQUEST_TIMEOUT = 15000
const LOGOUT_DELAY = 500
const MAX_RETRIES = 0
const RETRY_DELAY = 1000
const UNAUTHORIZED_DEBOUNCE_TIME = 3000

/** 401防抖状态 */
let isUnauthorizedErrorShown = false
let unauthorizedTimer: NodeJS.Timeout | null = null

/** Token 刷新状态（并发 401 时仅发起一次刷新） */
let isRefreshing = false
const failedQueue: Array<{ resolve: (token: string) => void; reject: (err: unknown) => void }> = []

/** 扩展 AxiosRequestConfig */
interface ExtendedAxiosRequestConfig extends AxiosRequestConfig {
  showErrorMessage?: boolean
  showSuccessMessage?: boolean
}

const { VITE_API_URL, VITE_WITH_CREDENTIALS } = import.meta.env

/** Axios实例 */
const axiosInstance = axios.create({
  timeout: REQUEST_TIMEOUT,
  baseURL: VITE_API_URL,
  withCredentials: VITE_WITH_CREDENTIALS === 'true',
  validateStatus: (status) => status >= 200 && status < 300,
  transformResponse: [
    (data, headers) => {
      const decrypted = maybeDecryptResponseBody(data, headers)
      const payload = typeof decrypted === 'string' ? decrypted : data
      const contentType =
        headers['content-type'] ||
        (typeof (headers as any).get === 'function' && (headers as any).get('content-type'))
      if (typeof contentType === 'string' && contentType.includes('application/json')) {
        try {
          return JSON.parse(payload as string)
        } catch {
          return payload
        }
      }
      return payload
    }
  ]
})

/** 请求拦截器 */
axiosInstance.interceptors.request.use(
  (request: InternalAxiosRequestConfig) => {
    const { accessToken } = useUserStore()
    if (accessToken) request.headers.set('Authorization', accessToken)

    const skipCrypto =
      !isApiCryptoEnabled() ||
      request.responseType === 'blob' ||
      request.data instanceof FormData ||
      shouldSkipRequestBodyEncrypt(request.url)

    if (!skipCrypto) {
      const method = (request.method || 'GET').toUpperCase()
      const hasEntityBody =
        request.data !== undefined && request.data !== null && request.data !== ''
      if (hasEntityBody) {
        if (!request.headers['Content-Type']) {
          request.headers.set('Content-Type', 'application/json')
        }
        const raw = typeof request.data === 'string' ? request.data : JSON.stringify(request.data)
        request.headers.set(API_CRYPTO_HEADER, API_CRYPTO_HEADER_VALUE)
        request.data = JSON.stringify({ c: encryptUtf8ToEnvelopeB64(raw) })
      } else if (method === 'GET' || method === 'HEAD') {
        request.headers.set(API_CRYPTO_HEADER, API_CRYPTO_HEADER_VALUE)
      }
    } else if (
      request.data &&
      !(request.data instanceof FormData) &&
      !request.headers['Content-Type']
    ) {
      request.headers.set('Content-Type', 'application/json')
      request.data = JSON.stringify(request.data)
    }

    return request
  },
  (error) => {
    showError(createHttpError($t('httpMsg.requestConfigError'), ApiStatus.error))
    return Promise.reject(error)
  }
)

/** 处理刷新队列（刷新成功时让等待中的请求重试） */
function processQueue(token: string | null, err: unknown = null) {
  failedQueue.forEach(({ resolve, reject }) => (err ? reject(err) : resolve(token!)))
  failedQueue.length = 0
}

/** 尝试用 refreshToken 刷新短 token，返回新 token */
async function tryRefreshToken(): Promise<string> {
  const userStore = useUserStore()
  const refreshTokenValue = userStore.refreshToken
  if (!refreshTokenValue) throw new Error('No refresh token')

  const res = await axiosInstance.request<BaseResponse<Api.Auth.LoginResponse>>({
    method: 'POST',
    url: '/api/sys_auth/refresh_token',
    data: { refreshToken: refreshTokenValue }
  })
  const payload = (res.data as BaseResponse<Api.Auth.LoginResponse>).data
  if (!payload?.token) throw new Error('Invalid refresh response')
  userStore.setToken(payload.token, payload.refreshToken ?? refreshTokenValue)
  return payload.token
}

/** 处理 401：有 refreshToken 则尝试刷新并重试，否则退出 */
function handle401AndMaybeRetry(
  originalConfig: InternalAxiosRequestConfig
): Promise<AxiosResponse> {
  const userStore = useUserStore()

  if (!originalConfig.url?.includes('refresh_token') && userStore.refreshToken) {
    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        failedQueue.push({
          resolve: (token) => {
            originalConfig.headers.set('Authorization', token)
            resolve(axiosInstance.request(originalConfig))
          },
          reject
        })
      })
    }

    isRefreshing = true
    return tryRefreshToken()
      .then((newToken) => {
        processQueue(newToken, null)
        originalConfig.headers.set('Authorization', newToken)
        return axiosInstance.request(originalConfig)
      })
      .catch((refreshErr) => {
        processQueue(null, refreshErr)
        doLogout()
        return Promise.reject(refreshErr)
      })
      .finally(() => {
        isRefreshing = false
      })
  }

  doLogout()
  return Promise.reject(createHttpError($t('httpMsg.unauthorized'), ApiStatus.unauthorized))
}

/** 退出登录（延迟执行，供 401 场景调用） */
function doLogout() {
  if (!isUnauthorizedErrorShown) {
    isUnauthorizedErrorShown = true
    logOut()
    unauthorizedTimer = setTimeout(resetUnauthorizedError, UNAUTHORIZED_DEBOUNCE_TIME)
    showError(createHttpError($t('httpMsg.unauthorized'), ApiStatus.unauthorized), true)
  }
}

/** 响应拦截器 */
axiosInstance.interceptors.response.use(
  (response: AxiosResponse<BaseResponse>) => {
    if (response.config.responseType === 'blob') return response
    const { code, msg, data } = response.data
    if (code === ApiStatus.success) return response
    if (code === ApiStatus.unauthorized) {
      return handle401AndMaybeRetry(response.config as InternalAxiosRequestConfig).catch((err) =>
        Promise.reject(err instanceof HttpError ? err : handleError(err as any))
      )
    }
    throw createHttpError(
      formatBusinessErrorMessage(msg, data) || $t('httpMsg.requestFailed'),
      code
    )
  },
  (error) => {
    const originalConfig = error.config as InternalAxiosRequestConfig | undefined
    if (error.response?.status === ApiStatus.unauthorized && originalConfig) {
      return handle401AndMaybeRetry(originalConfig).catch((err) =>
        Promise.reject(err instanceof HttpError ? err : handleError(error))
      )
    }
    return Promise.reject(handleError(error))
  }
)

/** 统一创建HttpError */
function createHttpError(message: string, code: number) {
  return new HttpError(message, code)
}

/** 重置401防抖状态 */
function resetUnauthorizedError() {
  isUnauthorizedErrorShown = false
  if (unauthorizedTimer) clearTimeout(unauthorizedTimer)
  unauthorizedTimer = null
}

/** 退出登录函数 */
function logOut() {
  setTimeout(() => {
    useUserStore().logOut()
  }, LOGOUT_DELAY)
}

/** 是否需要重试 */
function shouldRetry(statusCode: number) {
  return [
    ApiStatus.requestTimeout,
    ApiStatus.internalServerError,
    ApiStatus.badGateway,
    ApiStatus.serviceUnavailable,
    ApiStatus.gatewayTimeout
  ].includes(statusCode)
}

/** 请求重试逻辑 */
async function retryRequest<T>(
  config: ExtendedAxiosRequestConfig,
  retries: number = MAX_RETRIES
): Promise<T> {
  try {
    return await request<T>(config)
  } catch (error) {
    if (retries > 0 && error instanceof HttpError && shouldRetry(error.code)) {
      await delay(RETRY_DELAY)
      return retryRequest<T>(config, retries - 1)
    }
    throw error
  }
}

/** 延迟函数 */
function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

/** 请求函数 */
async function request<T = any>(config: ExtendedAxiosRequestConfig): Promise<T> {
  // POST | PUT 参数自动填充
  if (
    ['POST', 'PUT'].includes(config.method?.toUpperCase() || '') &&
    config.params &&
    !config.data
  ) {
    config.data = config.params
    config.params = undefined
  }

  try {
    const res = await axiosInstance.request(config)

    // 二进制文件响应直接返回 res.data（Blob）
    if (config.responseType === 'blob') {
      return res.data as T
    }

    // 显示成功消息
    if (config.showSuccessMessage && (res.data as BaseResponse)?.msg) {
      showSuccess((res.data as BaseResponse).msg)
    }

    return (res.data as BaseResponse<T>).data as T
  } catch (error) {
    if (error instanceof HttpError && error.code !== ApiStatus.unauthorized) {
      const showMsg = config.showErrorMessage !== false
      showError(error, showMsg)
    }
    return Promise.reject(error)
  }
}

/** API方法集合 */
const api = {
  get<T>(config: ExtendedAxiosRequestConfig) {
    return retryRequest<T>({ ...config, method: 'GET' })
  },
  post<T>(config: ExtendedAxiosRequestConfig) {
    return retryRequest<T>({ ...config, method: 'POST' })
  },
  put<T>(config: ExtendedAxiosRequestConfig) {
    return retryRequest<T>({ ...config, method: 'PUT' })
  },
  del<T>(config: ExtendedAxiosRequestConfig) {
    return retryRequest<T>({ ...config, method: 'DELETE' })
  },
  request<T>(config: ExtendedAxiosRequestConfig) {
    return retryRequest<T>(config)
  }
}

export default api
