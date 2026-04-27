/**
 * v-auth 权限指令
 *
 * 基于当前用户的权限标识控制 DOM 的显示与隐藏。
 * 无对应权限时元素会从 DOM 中移除。
 *
 * ## 主要功能
 *
 * - 权限验证 - 根据用户信息中的按钮权限（info.buttons）判断
 * - DOM 控制 - 无权限时移除元素，而非仅隐藏
 * - 响应式更新 - 权限变化时在 updated 阶段重新校验
 *
 * ## 使用示例
 *
 * ```vue
 * <el-button v-auth="'sys:user:add'">新增用户</el-button>
 * <el-button v-auth="'sys:user:edit'">编辑</el-button>
 * ```
 *
 * ## 注意事项
 *
 * - 直接移除 DOM，不使用 v-if
 * - 权限来源：用户信息中的 buttons 数组（接口 /api/sys_auth/info 返回）
 *
 * @module directives/auth
 * @author Art Design Pro Team
 */

import { useUserStore } from '@/store/modules/user'
import { App, Directive, DirectiveBinding } from 'vue'

interface AuthBinding extends DirectiveBinding {
  value: string
}

function checkAuthPermission(el: HTMLElement, binding: AuthBinding): void {
  const userStore = useUserStore()
  // 全部由后端返回的权限标识（buttons）来判断，不判断超级管理员
  const userPermissions: string[] = userStore.info?.buttons ?? []
  const hasPermission = userPermissions.includes(binding.value)

  if (!hasPermission) {
    removeElement(el)
  }
}

function removeElement(el: HTMLElement): void {
  if (el.parentNode) {
    el.parentNode.removeChild(el)
  }
}

const authDirective: Directive = {
  mounted: checkAuthPermission,
  updated: checkAuthPermission
}

export function setupAuthDirective(app: App): void {
  app.directive('auth', authDirective)
}
