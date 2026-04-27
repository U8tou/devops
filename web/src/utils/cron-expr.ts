/**
 * 与后端 cronvalidate 一致：仅标准 crontab（可选秒），含 CRON_TZ= 前缀；不接受 @daily、@every 等描述符。
 */
import { CronExpressionParser } from 'cron-parser'
import { toString as cronToString } from 'cronstrue'
import 'cronstrue/locales/zh_CN'

const CRON_TZ_PREFIX = /^CRON_TZ=\S+\s+/

function stripCronTz(expr: string): string {
  return expr.replace(CRON_TZ_PREFIX, '')
}

/** 是否与后端 cronvalidate.ValidateExpr 一致（标准五/六段，无 @ 描述符）。 */
export function isValidCronExpr(expr: string): boolean {
  const trimmed = expr.trim()
  if (!trimmed) return false
  const body = stripCronTz(trimmed)
  if (/^@/.test(body.trim())) return false

  try {
    CronExpressionParser.parse(body)
    return true
  } catch {
    return false
  }
}

export function describeCronExpr(expr: string, locale: 'zh' | 'en'): string | undefined {
  if (!isValidCronExpr(expr)) return undefined
  const trimmed = expr.trim()
  const body = stripCronTz(trimmed)
  const cronstrueLocale = locale === 'zh' ? 'zh_CN' : 'en'

  try {
    return cronToString(body, { locale: cronstrueLocale })
  } catch {
    return undefined
  }
}
