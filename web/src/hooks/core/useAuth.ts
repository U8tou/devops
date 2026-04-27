/**
 * useAuth - 权限验证管理
 *
 * 提供统一的权限验证功能，支持前端和后端两种权限模式。
 * 用于控制页面按钮、操作等功能的显示和访问权限。
 *
 * ## 主要功能
 *
 * 1. 权限检查 - 检查用户是否拥有指定的权限标识
 * 2. 双模式支持 - 自动适配前端模式和后端模式的权限验证
 * 3. 前后端模式均以用户信息中的 buttons 作为按钮权限列表
 *
 * ## 使用示例
 *
 * ```typescript
 * const { hasAuth } = useAuth()
 *
 * // 检查是否有新增权限
 * if (hasAuth('add')) {
 *   // 显示新增按钮
 * }
 *
 * // 在模板中使用
 * <el-button v-if="hasAuth('edit')">编辑</el-button>
 * <el-button v-if="hasAuth('delete')">删除</el-button>
 * ```
 *
 * @module useAuth
 * @author Art Design Pro Team
 */

import { storeToRefs } from 'pinia'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()

export const useAuth = () => {
  const { info } = storeToRefs(userStore)

  // 用户按钮权限（接口返回的 buttons，如 ['sys:user:add', 'sys:user:edit']）
  const userAuthList = info.value?.buttons ?? []

  /**
   * 检查是否拥有某权限标识（全部由后端返回的 buttons 来判断，不判断超级管理员）
   * @param auth 权限标识
   * @returns 是否有权限
   */
  const hasAuth = (auth: string): boolean => {
    return userAuthList.includes(auth)
  }

  return {
    hasAuth
  }
}
