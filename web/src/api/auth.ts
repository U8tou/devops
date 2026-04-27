import request from '@/utils/http'

/**
 * 登录（对应 /sys_auth/login）
 * @param params 登录参数（对应 internal_sys_auth.LoginReq）
 * @returns 登录响应（对应 internal_sys_auth.LoginResp）
 */
export function fetchLogin(params: Api.Auth.LoginParams) {
  return request.post<Api.Auth.LoginResponse>({
    url: '/api/sys_auth/login',
    data: params
  })
}

/**
 * 用户注册（对应 /sys_auth/register，公开接口）
 */
export function fetchRegister(data: Api.Auth.RegisterParams) {
  return request.post<Api.Auth.RegisterResponse>({
    url: '/api/sys_auth/register',
    data
  })
}

/**
 * 获取用户信息（对应 /sys_auth/info）
 * @returns 用户信息（对应 internal_sys_auth.InfoResp）
 */
export function fetchGetUserInfo() {
  return request.get<Api.Auth.UserInfo>({
    url: '/api/sys_auth/info'
  })
}

/**
 * 使用 RefreshToken 换取新的 AccessToken（对应 /sys_auth/refresh_token）
 * @param data 刷新令牌（对应 sysauth.RefreshTokenReq）
 * @returns 新的 token 和 refreshToken（对应 sysauth.LoginResp）
 */
export function fetchRefreshToken(data: Api.Auth.RefreshTokenReq) {
  return request.post<Api.Auth.LoginResponse>({
    url: '/api/sys_auth/refresh_token',
    data
  })
}

/**
 * 登出（对应 /sys_auth/logout）
 * @returns 登出响应
 */
export function fetchLogout() {
  return request.post<void>({
    url: '/api/sys_auth/logout'
  })
}

/**
 * 当前用户修改密码（对应 /sys_auth/change_password）
 * @param data 原密码与新密码（对应 internal_sys_auth.ChangePasswordReq）
 * @returns 统一响应（对应 resp.MyResp）
 */
export function fetchChangePassword(data: Api.Auth.ChangePasswordReq) {
  return request.post<unknown>({
    url: '/api/sys_auth/change_password',
    data
  })
}

/**
 * 当前登录用户更新个人信息（对应 /sys_auth/profile）
 * @param data 个人信息（对应 internal_sys_auth.UpdateProfileReq）
 * @returns 统一响应（对应 resp.MyResp）
 */
export function fetchUpdateProfile(data: Api.Auth.UpdateProfileReq) {
  return request.post<unknown>({
    url: '/api/sys_auth/profile',
    data
  })
}

/**
 * 当前登录用户上传头像（对应 /sys_auth/avatar）
 * @param file 头像图片文件
 * @returns 头像地址（对应 internal_sys_auth.UploadAvatarResp）
 */
export function fetchUploadAvatar(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  return request.post<Api.Auth.UploadAvatarResp>({
    url: '/api/sys_auth/avatar',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
