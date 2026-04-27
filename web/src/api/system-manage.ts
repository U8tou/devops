import request from '@/utils/http'

// ==================== 用户管理接口 ====================

/**
 * 获取用户列表（分页）（对应 /sys_user/list）
 * @param params 搜索参数
 * @returns 用户列表（对应 internal_sys_user.ListResp）
 */
export function fetchGetUserList(params: Api.SystemManage.UserSearchParams) {
  return request.get<Api.SystemManage.UserList>({
    url: '/api/sys_user/list',
    params
  })
}

/**
 * 获取用户详情（对应 /sys_user/get）
 * @param id 用户ID
 * @returns 用户详情（对应 internal_sys_user.GetResp）
 */
export function fetchGetUserDetail(id: string) {
  return request.get<Api.SystemManage.UserDetail>({
    url: '/api/sys_user/get',
    params: { id }
  })
}

/**
 * 获取所有用户（对应 /sys_user/all）
 * @returns 用户列表（对应 internal_sys_user.ListResp）
 */
export function fetchGetUserAll() {
  return request.get<Api.SystemManage.UserList>({
    url: '/api/sys_user/all'
  })
}

/**
 * 新增用户（对应 /sys_user/add）
 * @param data 新增请求（对应 internal_sys_user.AddReq）
 * @returns 响应（对应 internal_sys_user.AddResp）
 */
export function fetchAddUser(data: Api.SystemManage.UserAddReq) {
  return request.post<{ affect: string }>({
    url: '/api/sys_user/add',
    data
  })
}

/**
 * 编辑用户（对应 /sys_user/edit）
 * @param data 编辑请求（对应 internal_sys_user.EditReq）
 * @returns 响应（对应 internal_sys_user.EditResp）
 */
export function fetchEditUser(data: Api.SystemManage.UserEditReq) {
  return request.put<{ affect: string }>({
    url: '/api/sys_user/edit',
    data
  })
}

/**
 * 删除用户（对应 /sys_user/del）
 * @param ids 用户ID数组
 * @returns 响应（对应 internal_sys_user.DelResp）
 */
export function fetchDeleteUser(ids: string[]) {
  return request.del<{ affect: string }>({
    url: '/api/sys_user/del',
    data: { ids }
  })
}

/**
 * 为指定用户重置密码（对应 /sys_user/reset_pwd）
 * @param data 重置密码请求（对应 internal_sys_user.ResetPasswordReq）
 * @returns 响应（对应 internal_sys_user.ResetPasswordResp）
 */
export function fetchResetUserPassword(data: Api.SystemManage.UserResetPasswordReq) {
  return request.post<{ affect: string }>({
    url: '/api/sys_user/reset_pwd',
    data: {
      userId: String(data.userId),
      password: data.password
    }
  })
}

/**
 * 为用户分配角色（对应 /sys_user/assign_role）
 * @param data 分配请求（对应 internal_sys_user.AssignRoleReq）
 * @returns 响应（对应 internal_sys_user.AssignRoleResp）
 */
export function fetchAssignUserRole(data: Api.SystemManage.UserAssignRoleReq) {
  return request.post<{ affect: string }>({
    url: '/api/sys_user/assign_role',
    data: {
      userId: String(data.userId),
      roleIds: data.roleIds.map((id) => String(id))
    }
  })
}

/**
 * 为用户分配部门（对应 /sys_user/assign_dept）
 * @param data 分配请求（对应 internal_sys_user.AssignDeptReq）
 * @returns 响应（对应 internal_sys_user.AssignDeptResp）
 */
export function fetchAssignUserDept(data: Api.SystemManage.UserAssignDeptReq) {
  return request.post<{ affect: string }>({
    url: '/api/sys_user/assign_dept',
    data: {
      userId: String(data.userId),
      deptIds: data.deptIds.map((id) => String(id))
    }
  })
}

/**
 * 为用户分配岗位（对应 /sys_user/assign_post）
 * @param data 分配请求（对应 internal_sys_user.AssignPostReq）
 * @returns 响应（对应 internal_sys_user.AssignPostResp）
 */
export function fetchAssignUserPost(data: Api.SystemManage.UserAssignPostReq) {
  return request.post<{ affect: string }>({
    url: '/api/sys_user/assign_post',
    data: {
      userId: String(data.userId),
      postIds: data.postIds.map((id) => String(id))
    }
  })
}

/**
 * 按条件导出用户列表为 Excel（对应 /sys_user/export）
 * @param params 与 list 相同的查询参数
 * @returns 文件 Blob（application/octet-stream，如 用户列表.xlsx）
 */
export function fetchExportUser(params?: Api.SystemManage.UserSearchParams) {
  return request.get<Blob>({
    url: '/api/sys_user/export',
    params,
    responseType: 'blob'
  })
}

/**
 * 上传 Excel 导入用户（对应 /sys_user/import）
 * @param file 按模板填写的 xlsx 文件
 * @param options 可选：updateExisting 是否更新已存在的用户数据（true/1 为是）
 * @returns 导入结果（对应 internal_sys_user.ImportResp）
 */
export function fetchImportUser(file: File, options?: { updateExisting?: boolean }) {
  const formData = new FormData()
  formData.append('file', file)
  if (options?.updateExisting !== undefined) {
    formData.append('updateExisting', String(options.updateExisting))
  }
  return request.post<Api.SystemManage.UserImportResp>({
    url: '/api/sys_user/import',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 下载用户导入 Excel 模板（对应 /sys_user/template）
 * @returns 模板文件 Blob（application/octet-stream，如 用户导入模板.xlsx）
 */
export function fetchGetUserTemplate() {
  return request.get<Blob>({
    url: '/api/sys_user/template',
    responseType: 'blob'
  })
}

// ==================== 角色管理接口 ====================

/**
 * 获取角色列表（分页）（对应 /sys_role/list）
 * @param params 搜索参数
 * @returns 角色列表（对应 internal_sys_role.ListResp）
 */
export function fetchGetRoleList(params: Api.SystemManage.RoleSearchParams) {
  return request.get<Api.SystemManage.RoleList>({
    url: '/api/sys_role/list',
    params
  })
}

/**
 * 获取所有角色（对应 /sys_role/all）
 * @returns 角色列表（对应 internal_sys_role.ListBody[]）
 */
export function fetchGetRoleAll() {
  return request.get<Api.SystemManage.RoleListItem[]>({
    url: '/api/sys_role/all'
  })
}

/**
 * 获取角色详情（对应 /sys_role/get）
 * @param id 角色ID
 * @returns 角色详情（对应 internal_sys_role.GetResp）
 */
export function fetchGetRoleInfo(id: string) {
  return request.get<Api.SystemManage.RoleDetail>({
    url: '/api/sys_role/get',
    params: { id }
  })
}

/**
 * 新增角色（对应 /sys_role/add）
 * @param data 新增请求（对应 internal_sys_role.AddReq）
 * @returns 响应（对应 internal_sys_role.AddResp）
 */
export function fetchAddRole(data: Api.SystemManage.RoleAddReq) {
  return request.post<{ id: string }>({
    url: '/api/sys_role/add',
    data
  })
}

/**
 * 编辑角色（对应 /sys_role/edit）
 * @param data 编辑请求（对应 internal_sys_role.EditReq）
 * @returns 响应（对应 internal_sys_role.EditResp）
 */
export function fetchEditRole(data: Api.SystemManage.RoleEditReq) {
  return request.put<{ affect: string }>({
    url: '/api/sys_role/edit',
    data
  })
}

/**
 * 删除角色（对应 /sys_role/del）
 * @param ids 角色ID数组（字符串数组）
 * @returns 响应（对应 internal_sys_role.DelResp）
 */
export function fetchDeleteRole(ids: string[]) {
  return request.del<{ affect: string }>({
    url: '/api/sys_role/del',
    data: { ids }
  })
}

/**
 * 为角色分配菜单权限（对应 /sys_role/assign_menu）
 * @param data 分配请求（对应 internal_sys_role.AssignMenuReq）
 * @returns 响应（对应 internal_sys_role.AssignMenuResp）
 */
export function fetchAssignRoleMenu(data: Api.SystemManage.RoleAssignMenuReq) {
  return request.post<{ affect: string }>({
    url: '/api/sys_role/assign_menu',
    data: {
      roleId: String(data.roleId),
      menuIds: data.menuIds.map((id) => String(id)),
      menuLinkage: data.menuLinkage
    }
  })
}

/**
 * 为角色分配部门权限（对应 /sys_role/assign_dept）
 * @param data 分配请求（对应 internal_sys_role.AssignDeptReq）
 * @returns 响应（对应 internal_sys_role.AssignDeptResp）
 */
export function fetchAssignRoleDept(data: Api.SystemManage.RoleAssignDeptReq) {
  return request.post<{ affect: string }>({
    url: '/api/sys_role/assign_dept',
    data: {
      roleId: String(data.roleId),
      deptIds: data.deptIds.map((id) => String(id)),
      deptLinkage: data.deptLinkage
    }
  })
}

// ==================== 部门管理接口 ====================

/**
 * 获取部门列表（分页）（对应 /sys_dept/list）
 * @param params 搜索参数
 * @returns 部门列表（对应 internal_sys_dept.ListResp）
 */
export function fetchGetDeptList(params?: Api.SystemManage.DeptSearchParams) {
  return request.get<Api.SystemManage.DeptList>({
    url: '/api/sys_dept/list',
    params
  })
}

/**
 * 获取部门列表（全部）（对应 /sys_dept/all）
 * @returns 部门列表（对应 internal_sys_dept.ListResp）
 */
export function fetchGetDeptAll() {
  return request.get<Api.SystemManage.DeptAllResponse>({
    url: '/api/sys_dept/all'
  })
}

/**
 * 获取部门详情（对应 /sys_dept/get）
 * @param id 部门ID
 * @returns 部门详情（对应 internal_sys_dept.GetResp）
 */
export function fetchGetDeptInfo(id: string) {
  return request.get<Api.SystemManage.DeptDetail>({
    url: '/api/sys_dept/get',
    params: { id }
  })
}

/**
 * 新增部门（对应 /sys_dept/add）
 * @param data 新增请求（对应 internal_sys_dept.AddReq）
 * @returns 响应（对应 internal_sys_dept.AddResp）
 */
export function fetchAddDept(data: Api.SystemManage.DeptAddReq) {
  return request.post<{ affect: string }>({
    url: '/api/sys_dept/add',
    data
  })
}

/**
 * 编辑部门（对应 /sys_dept/edit）
 * @param data 编辑请求（对应 internal_sys_dept.EditReq）
 * @returns 响应（对应 internal_sys_dept.EditResp）
 */
export function fetchEditDept(data: Api.SystemManage.DeptEditReq) {
  return request.put<{ affect: string }>({
    url: '/api/sys_dept/edit',
    data
  })
}

/**
 * 删除部门（对应 /sys_dept/del，ids 走 query csv）
 * @param ids 部门ID数组
 * @returns 响应（对应 internal_sys_dept.DelResp）
 */
export function fetchDeleteDept(ids: string[]) {
  return request.del<{ affect: string }>({
    url: '/api/sys_dept/del',
    params: { ids: ids.length ? ids.join(',') : '' }
  })
}

// ==================== 岗位管理接口 ====================

/**
 * 获取岗位列表（分页）（对应 /sys_post/list）
 * @param params 搜索参数
 * @returns 岗位列表（对应 internal_sys_post.ListResp）
 */
export function fetchGetPostList(params?: Api.SystemManage.PostSearchParams) {
  return request.get<Api.SystemManage.PostListResponse>({
    url: '/api/sys_post/list',
    params
  })
}

/**
 * 获取岗位列表（全部）（对应 /sys_post/all）
 * @returns 岗位列表（对应 internal_sys_post.ListResp）
 */
export function fetchGetPostAll() {
  return request.get<Api.SystemManage.PostListResponse>({
    url: '/api/sys_post/all'
  })
}

/**
 * 获取岗位详情（对应 /sys_post/get）
 * @param id 岗位ID
 * @returns 岗位详情（对应 internal_sys_post.GetResp）
 */
export function fetchGetPostInfo(id: string) {
  return request.get<Api.SystemManage.PostDetail>({
    url: '/api/sys_post/get',
    params: { id }
  })
}

/**
 * 新增岗位（对应 /sys_post/add）
 * @param data 新增请求（对应 internal_sys_post.AddReq）
 * @returns 响应（对应 internal_sys_post.AddResp）
 */
export function fetchAddPost(data: Api.SystemManage.PostAddReq) {
  return request.post<{ affect: string }>({
    url: '/api/sys_post/add',
    data
  })
}

/**
 * 编辑岗位（对应 /sys_post/edit）
 * @param data 编辑请求（对应 internal_sys_post.EditReq）
 * @returns 响应（对应 internal_sys_post.EditResp）
 */
export function fetchEditPost(data: Api.SystemManage.PostEditReq) {
  return request.put<{ affect: string }>({
    url: '/api/sys_post/edit',
    data
  })
}

/**
 * 删除岗位（对应 /sys_post/del）
 * @param ids 岗位ID数组
 * @returns 响应（对应 internal_sys_post.DelResp）
 */
export function fetchDeletePost(ids: string[]) {
  return request.del<{ affect: string }>({
    url: '/api/sys_post/del',
    data: { ids }
  })
}

// ==================== 菜单管理接口 ====================

/**
 * 获取菜单列表（分页）（对应 /sys_menu/list）
 * @param params 搜索参数
 * @returns 菜单列表（对应 internal_sys_menu.ListResp）
 */
export function fetchGetMenuList(params?: Api.SystemManage.MenuSearchParams) {
  return request.get<Api.SystemManage.MenuListResponse>({
    url: '/api/sys_menu/list',
    params
  })
}

/**
 * 获取所有菜单（全量列表）（对应 /sys_menu/all）
 * @returns 菜单全量列表（对应 internal_sys_menu.ListResp）
 */
export function fetchGetMenuAll() {
  return request.get<Api.SystemManage.MenuListResponse>({
    url: '/api/sys_menu/all'
  })
}

/**
 * 获取菜单详情（对应 /sys_menu/get）
 * @param id 菜单ID（字符串）
 * @returns 菜单详情（对应 internal_sys_menu.GetResp）
 */
export function fetchGetMenuInfo(id: string) {
  return request.get<Api.SystemManage.MenuDetail>({
    url: '/api/sys_menu/get',
    params: { id }
  })
}
