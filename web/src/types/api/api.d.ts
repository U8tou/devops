/**
 * API 接口类型定义模块
 *
 * 提供所有后端接口的类型定义
 *
 * ## 主要功能
 *
 * - 通用类型（分页参数、响应结构等）
 * - 认证类型（登录、用户信息等）
 * - 系统管理类型（用户、角色等）
 * - 全局命名空间声明
 *
 * ## 使用场景
 *
 * - API 请求参数类型约束
 * - API 响应数据类型定义
 * - 接口文档类型同步
 *
 * ## 注意事项
 *
 * - 在 .vue 文件使用需要在 eslint.config.mjs 中配置 globals: { Api: 'readonly' }
 * - 使用全局命名空间，无需导入即可使用
 *
 * ## 使用方式
 *
 * ```typescript
 * const params: Api.Auth.LoginParams = { userName: 'admin', password: '123456' }
 * const response: Api.Auth.UserInfo = await fetchUserInfo()
 * ```
 *
 * @module types/api/api
 * @author Art Design Pro Team
 */

declare namespace Api {
  /** 通用类型 */
  namespace Common {
    /** 获取文件参数（对应 /common/get_file/{dir}/{obj}） */
    interface GetFileParams {
      /** 文件夹 */
      dir: string
      /** 文件名 */
      obj: string
    }

    /** SSE 连接参数（对应 GET /sys_sse） */
    interface SseConnectionParams {
      /** 鉴权 token（可选），有则登录用户 sseId=loginId_connId 定向推送，无则游客 sseId=guest_connId 仅广播 */
      token?: string
    }

    /** WebSocket 连接参数（对应 GET /sys_websocket，query uuid） */
    interface WebSocketConnectionParams {
      /** 连接标识，用于后续 SendMsg 定向推送 */
      uuid: string
    }

    /** IM 聊天室 WebSocket 连接参数（对应 GET /sys_im/chat） */
    interface ImChatConnectionParams {
      /** 鉴权 token（可选），有则显示用户昵称和头像，无则为游客。nickName 为纯昵称，会话由 clientId(loginId_connId) 区分 */
      token?: string
    }

    /** IM 服务端消息基础（type 字段区分类型） */
    interface ImChatServerMessageBase {
      type: 'message' | 'joined' | 'left' | 'error' | 'online' | 'offline'
    }

    /** IM 聊天消息（type=message）。判断 isMe 用 loginId 与本人 loginId 比较 */
    interface ImChatServerMessage extends ImChatServerMessageBase {
      type: 'message'
      from?: string
      nickName?: string
      avatar?: string
      content?: string
      time?: string
      /** 会话唯一标识 loginId_connId */
      clientId?: string
      /** 登录标识，用于判断是否本人消息 */
      loginId?: string
      userId?: string
    }

    /** IM 成员加入/离开等通知（type=joined|left|online|offline），会推送 memberCount */
    interface ImChatServerNotify extends ImChatServerMessageBase {
      type: 'joined' | 'left' | 'online' | 'offline'
      roomId?: string
      nickName?: string
      from?: string
      clientId?: string
      loginId?: string
      userId?: string
      /** 当前房间在线人数。joined=含自己；left=不含自己；online=含新上线者；offline=不含下线者 */
      memberCount?: number
    }

    /** IM 错误消息（type=error） */
    interface ImChatServerError extends ImChatServerMessageBase {
      type: 'error'
      msg?: string
    }

    /** 分页参数 */
    interface PaginationParams {
      /** 当前页码 */
      current: number
      /** 每页条数 */
      size: number
      /** 总条数 */
      total: number
    }

    /** 通用搜索参数 */
    type CommonSearchParams = Pick<PaginationParams, 'current' | 'size'>

    /** 分页响应基础结构 */
    interface PaginatedResponse<T = any> {
      records?: T[]
      rows?: T[] // 部分接口使用 rows
      current?: number
      size?: number
      total: number | string // 部分接口返回 string 类型的 total
    }

    /** 启用状态 */
    type EnableStatus = '1' | '2'
  }

  /** 认证类型 */
  namespace Auth {
    /** 登录参数（对应 internal_sys_auth.LoginReq） */
    interface LoginParams {
      /** 用户名 */
      userName: string
      /** 密码 */
      password: string
      /** 验证码（可选） */
      code?: string
      /** 验证码ID（可选） */
      codeId?: string
    }

    /** 登录响应（对应 internal_sys_auth.LoginResp） */
    interface LoginResponse {
      token: string
      refreshToken: string
    }

    /** 注册请求（对应 internal_sys_auth.RegisterReq） */
    interface RegisterParams {
      userName: string
      password: string
      nickName?: string
      email?: string
    }

    /** 注册响应（对应 internal_sys_auth.RegisterResp） */
    interface RegisterResponse {
      token: string
      refreshToken: string
      needBindRole: boolean
    }

    /** Token 刷新请求（对应 sysauth.RefreshTokenReq） */
    interface RefreshTokenReq {
      /** 刷新令牌 */
      refreshToken: string
    }

    /** 用户信息（对应 internal_sys_auth.InfoResp，仅含文档中字段） */
    interface UserInfo {
      /** 联系地址 */
      address: string
      /** 头像 */
      avatar: string
      /** 创建者 */
      createBy: string
      /** 创建时间 */
      createTime: string
      /** 用户邮箱 */
      email: string
      /** 用户昵称 */
      nickName: string
      /** 电话号码 */
      phone: string
      /** 电话区号 */
      phoneArea: string
      /** 备注 */
      remark: string
      /** 用户性别:1_男,2_女,3_未知 */
      sex: number
      /** 状态:1_正常,2_停用 */
      status: number
      /** 更新者 */
      updateBy: string
      /** 更新时间 */
      updateTime: string
      /** 用户ID */
      userId: string
      /** 用户账号 */
      userName: string
      /** 用户类型:00_系统用户 */
      userType: string

      /** 按钮权限（可选） */
      buttons?: string[]
      /** 菜单权限（可选） */
      menus?: string[]
      /** 关联角色IDs（可选） */
      roles?: string[]
      /** 关联部门IDs（可选） */
      depts?: string[]
      /** 关联岗位IDs（可选） */
      posts?: string[]
    }

    /** 修改密码请求（对应 internal_sys_auth.ChangePasswordReq） */
    interface ChangePasswordReq {
      oldPassword: string
      newPassword: string
    }

    /** 更新个人信息请求（对应 internal_sys_auth.UpdateProfileReq） */
    interface UpdateProfileReq {
      /** 联系地址 */
      address?: string
      /** 用户邮箱 */
      email?: string
      /** 用户昵称 */
      nickName?: string
      /** 电话号码 */
      phone?: string
      /** 电话区号 */
      phoneArea?: string
      /** 用户性别:1_男,2_女,3_未知 */
      sex?: number
    }

    /** 上传头像响应（对应 internal_sys_auth.UploadAvatarResp） */
    interface UploadAvatarResp {
      /** 头像地址（相对路径或完整 URL） */
      avatar: string
    }
  }

  /** 系统管理类型 */
  namespace SystemManage {
    /** 用户列表 */
    type UserList = Api.Common.PaginatedResponse<UserListItem>

    /** 用户列表项（与后端结构体一致，对应 internal_sys_user.ListBody） */
    interface UserListItem {
      id: string
      userName: string
      nickName: string
      avatar: string
      sex: number // 用户性别:1_男,2_女,3_未知
      phone: string
      phoneArea: string
      email: string
      address: string
      status: number // 状态:1_正常,2_停用
      userType: string // 用户类型:00_系统用户
      remark: string
      createBy: string
      createTime: string
      updateBy: string
      updateTime: string
      depts?: string[] // 关联部门IDs（列表返回中包含，不在 required 中）
      roles?: string[] // 关联角色IDs（列表返回中包含，不在 required 中）
      posts?: string[] // 关联岗位IDs
    }

    /** 用户详情（包含关联的角色、部门和岗位，对应 internal_sys_user.GetResp） */
    interface UserDetail {
      id: string
      userName: string
      nickName: string
      avatar: string
      sex: number // 用户性别:1_男,2_女,3_未知
      phone: string
      phoneArea: string
      email: string
      address: string
      status: number // 状态:1_正常,2_停用
      userType: string // 用户类型:00_系统用户
      remark: string
      createBy: string
      createTime: string
      updateBy: string
      updateTime: string
      /** 关联部门IDs */
      depts?: string[]
      /** 关联角色IDs */
      roles?: string[]
      /** 关联岗位IDs */
      posts?: string[]
    }

    /** 用户新增请求（对应 internal_sys_user.AddReq，仅含文档中字段） */
    interface UserAddReq {
      userName: string
      password: string
      nickName: string
      sex: number
      phone: string
      phoneArea: string
      email: string
      address: string
      status: number
      userType: string
      remark: string
      avatar?: string
      deleteTime?: string
    }

    /** 用户编辑请求（对应 internal_sys_user.EditReq，仅含文档中字段） */
    interface UserEditReq {
      id: string
      userName: string
      nickName: string
      sex: number
      phone: string
      phoneArea: string
      email: string
      address: string
      status: number
      userType: string
      remark: string
      password?: string // 密码，空表示不修改
      avatar?: string
      deleteTime?: string
    }

    /** 用户重置密码请求（对应 internal_sys_user.ResetPasswordReq） */
    interface UserResetPasswordReq {
      /** 用户ID */
      userId: string
      /** 新密码 */
      password: string
    }

    /** 用户分配角色请求（对应 internal_sys_user.AssignRoleReq） */
    interface UserAssignRoleReq {
      userId: string
      roleIds: string[]
    }

    /** 用户分配部门请求（对应 internal_sys_user.AssignDeptReq） */
    interface UserAssignDeptReq {
      userId: string
      deptIds: string[]
    }

    /** 用户分配岗位请求（对应 internal_sys_user.AssignPostReq） */
    interface UserAssignPostReq {
      userId: string
      postIds: string[]
    }

    /** 用户导入失败项（对应 internal_sys_user.ImportFail） */
    interface UserImportFail {
      /** Excel 行号（从 2 起为数据行） */
      row: number
      /** 失败原因 */
      reason: string
    }

    /** 用户导入响应（对应 internal_sys_user.ImportResp） */
    interface UserImportResp {
      /** 成功条数 */
      successCount: number
      /** 失败行（行号+原因） */
      failList?: UserImportFail[]
    }

    /** 用户搜索参数（支持基础字段及分页查询，与 /sys_user/list、/sys_user/export 保持一致） */
    type UserSearchParams = Partial<
      Pick<UserListItem, 'id' | 'userName' | 'nickName' | 'sex' | 'phone' | 'email' | 'status'> &
        Api.Common.CommonSearchParams & {
          /** 部门ID */
          deptId?: string
          /** 时间范围（创建时间，Unix 秒时间戳区间：[开始,结束]，导出和列表均可用） */
          timeRange?: string[]
        }
    >

    /** 角色列表 */
    type RoleList = Api.Common.PaginatedResponse<RoleListItem>

    /** 角色列表项（与后端结构体一致，对应 internal_sys_role.ListBody） */
    interface RoleListItem {
      id: string // 角色ID
      name: string // 名称
      role: string // 角色标识
      status: number // 状态:1_正常,2_停用
      menuLinkage?: number // 操作父子联动:1_联动,2_不联动
      deptLinkage?: number // 数据父子联动:1_联动,2_不联动
      sort: number // 显示顺序
      remark: string // 备注
      createTime?: string // 创建时间（字符串格式）
      updateTime?: string // 更新时间（字符串格式）
    }

    /** 角色详情（对应 internal_sys_role.GetResp） */
    interface RoleDetail {
      id: string
      name: string
      role: string
      status: number
      sort: number
      remark: string
      createBy?: string
      createTime?: string
      updateBy?: string
      updateTime?: string
      menuIds?: string[] // 菜单ID列表
      menuLinkage?: number // 操作父子联动
      deptIds?: string[] // 部门ID列表
      deptLinkage?: number // 数据父子联动
    }

    /** 角色新增请求（对应 internal_sys_role.AddReq） */
    interface RoleAddReq {
      name: string
      role: string
      status?: number // 状态:1_正常,2_停用，默认1
      sort?: number // 显示顺序，默认999
      remark?: string
    }

    /** 角色编辑请求（对应 internal_sys_role.EditReq） */
    interface RoleEditReq {
      id: string
      name: string
      role: string
      status?: number // 状态:1_正常,2_停用，默认1
      sort?: number // 显示顺序，默认999
      remark?: string
    }

    /** 角色分配菜单请求（对应 internal_sys_role.AssignMenuReq） */
    interface RoleAssignMenuReq {
      roleId: string
      menuIds: string[]
      menuLinkage: number // 操作父子联动:1_联动,2_不联动
    }

    /** 角色分配部门请求（对应 internal_sys_role.AssignDeptReq） */
    interface RoleAssignDeptReq {
      roleId: string
      deptIds: string[]
      deptLinkage: number // 数据父子联动:1_联动,2_不联动
    }

    /** 角色搜索参数 */
    type RoleSearchParams = Partial<
      Pick<RoleListItem, 'id' | 'name' | 'role' | 'status'> & Api.Common.CommonSearchParams
    >

    /** 部门列表 */
    type DeptList = Api.Common.PaginatedResponse<DeptListItem>

    /** 部门全部列表接口返回（/all 返回 total + rows） */
    interface DeptAllResponse {
      total: string
      rows: DeptListItem[]
    }

    /** 部门列表项（与后端结构体一致，对应 internal_sys_dept.ListBody） */
    interface DeptListItem {
      id: string
      pid?: string // 上级部门ID，不在 required 中
      name: string
      profile?: string // 部门简介，不在 required 中
      leader: string
      phone: string // 负责人电话，在 required 中
      email: string // 负责人邮箱，在 required 中
      sort: number
      status?: number // 状态:1_正常,2_停用，不在 required 中
      createTime?: string // 创建时间，不在 required 中
      updateTime?: string // 更新时间，不在 required 中
    }

    /** 部门详情（对应 internal_sys_dept.GetResp） */
    interface DeptDetail {
      id: string
      pid?: string
      name: string
      profile?: string
      leader: string
      phone: string
      email: string
      sort: number
      status?: number
      createBy?: string
      createTime?: string
      updateBy?: string
      updateTime?: string
    }

    /** 部门新增请求（对应 internal_sys_dept.AddReq） */
    interface DeptAddReq {
      name: string
      leader: string
      sort: number
      pid?: string
      profile?: string
      phone?: string
      email?: string
      status?: number // 状态:1_正常,2_停用
    }

    /** 部门编辑请求（对应 internal_sys_dept.EditReq） */
    interface DeptEditReq {
      id: string
      name: string
      leader: string
      sort: number
      pid?: string
      profile?: string
      phone?: string
      email?: string
      status?: number // 状态:1_正常,2_停用
    }

    /** 部门树节点（handleTree 后的结构，带 children） */
    interface DeptTreeNode extends DeptListItem {
      children?: DeptTreeNode[]
    }

    /** 部门搜索参数 */
    type DeptSearchParams = Partial<
      Pick<
        DeptListItem,
        'id' | 'name' | 'leader' | 'phone' | 'email' | 'status' | 'pid' | 'profile' | 'sort'
      > &
        Api.Common.CommonSearchParams
    >

    /** 岗位列表（对应 internal_sys_post.ListResp） */
    type PostList = Api.Common.PaginatedResponse<PostListItem>

    /** 岗位列表项（对应 internal_sys_post.ListBody） */
    interface PostListItem {
      id: string
      name: string
      sort: number
      status?: number // 状态:1_正常,2_停用
      createTime?: string
      updateTime?: string
    }

    /** 岗位详情（对应 internal_sys_post.GetResp） */
    interface PostDetail {
      id: string
      name: string
      sort: number
      status?: number
      createBy?: string
      createTime?: string
      updateBy?: string
      updateTime?: string
    }

    /** 岗位新增请求（对应 internal_sys_post.AddReq） */
    interface PostAddReq {
      name: string
      sort: number
      status?: number // 状态:1_正常,2_停用
    }

    /** 岗位编辑请求（对应 internal_sys_post.EditReq） */
    interface PostEditReq {
      id: string
      name: string
      sort: number
      status?: number // 状态:1_正常,2_停用
    }

    /** 岗位列表响应 */
    interface PostListResponse {
      rows: PostListItem[]
      total: string
    }

    /** 岗位搜索参数 */
    type PostSearchParams = Partial<
      Pick<PostListItem, 'id' | 'name' | 'status' | 'sort'> & Api.Common.CommonSearchParams
    >

    /** 菜单列表项（对应 internal_sys_menu.ListBody） */
    interface MenuListItem {
      id: string
      pid?: string
      permis: string // 权限代码
      types: number // 类型:1_菜单,2_权限
      remark?: string
      createTime?: string
    }

    /** 菜单详情（对应 internal_sys_menu.GetResp） */
    interface MenuDetail {
      id: string
      pid?: string
      permis: string
      types: number
      remark?: string
      createBy?: string
      createTime?: string
      updateBy?: string
      updateTime?: string
    }

    /** 菜单树节点（对应 internal_sys_menu.TreeResp） */
    interface MenuTreeNode {
      id: string
      pid?: string
      permis: string
      types: number
      remark?: string
      children?: MenuTreeNode[]
    }

    /** 菜单列表响应 */
    interface MenuListResponse {
      rows: MenuListItem[]
      total: string
    }

    /** 菜单新增请求（对应 internal_sys_menu.AddReq，当前 Swagger 未提供 add 接口） */
    interface MenuAddReq {
      permis: string
      types: number // 类型:1_菜单,2_权限
      pid?: string
      remark?: string
    }

    /** 菜单编辑请求（对应 internal_sys_menu.EditReq，当前 Swagger 未提供 edit 接口） */
    interface MenuEditReq {
      id: string
      permis: string
      types: number // 类型:1_菜单,2_权限
      pid?: string
      remark?: string
    }

    /** 菜单搜索参数 */
    type MenuSearchParams = Partial<
      Pick<MenuListItem, 'id' | 'pid' | 'permis' | 'types'> & Api.Common.CommonSearchParams
    >
  }

  /** 自动化流程 dev_process */
  namespace DevProcess {
    interface ProcessListItem {
      id: string
      code: string
      remark: string
      cronExpr?: string
      /** 1 开启定时执行 0 关闭 */
      cronEnabled?: string
      lastExecTime?: string
      /** 最近一次执行总耗时（毫秒），字符串数字；无执行或未记录为 0 */
      lastExecDurationMs?: string
      /** success | failed | cancelled */
      lastExecResult?: string
      createTime?: string
      updateTime?: string
      tags?: ProcessTagItem[]
    }

    interface ProcessTagItem {
      id: string
      /** 空表示字典已删除，仅孤儿关联 */
      name: string
    }

    interface ProcessListResponse {
      rows: ProcessListItem[]
      total: string
    }

    interface ProcessDetail {
      id: string
      code: string
      remark: string
      flow: string
      /** 流程级环境变量（JSON 解析为 KV） */
      env?: Record<string, string>
      /** 1 开启定时执行 0 关闭 */
      cronEnabled?: string
      cronExpr?: string
      lastExecTime?: string
      lastExecDurationMs?: string
      /** success | failed | cancelled */
      lastExecResult?: string
      /** 仅详情接口返回完整执行日志 */
      lastExecLog?: string
      createTime?: string
      createBy?: string
      updateTime?: string
      updateBy?: string
      /** 流程标签；name 为空表示字典已删除 */
      tags?: ProcessTagItem[]
    }

    interface ProcessAddReq {
      code: string
      remark?: string
      flow?: string
      env?: Record<string, string>
      /** 0 | 1 是否开启定时执行 */
      cronEnabled?: number
      cronExpr?: string
      tagIds?: number[]
    }

    interface ProcessAddRes {
      affect: string
      id: string
    }

    /** 仅更新编号、备注 */
    interface ProcessEditReq {
      id: string
      code: string
      remark?: string
      /** 0 | 1 是否开启定时执行 */
      cronEnabled?: number
      cronExpr?: string
      tagIds?: number[]
    }

    /** 仅更新流程 JSON */
    interface ProcessEditFlowReq {
      id: string
      flow: string
    }

    /** 仅更新流程环境变量 JSON */
    interface ProcessEditEnvReq {
      id: string
      env: Record<string, string>
    }

    /** 设置是否开启定时执行（列表开关） */
    interface ProcessSetCronEnabledReq {
      id: string
      enabled: boolean
    }

    /** 删除流程（软删） */
    interface ProcessDelReq {
      ids: string[]
    }

    /** 节点参数快速校验（录入时「验证连接」） */
    interface ValidateNodeReq {
      kind: string
      params?: Record<string, unknown>
    }

    interface ValidateNodeRes {
      ok: boolean
      message: string
      detail?: string
    }

    type ProcessSearchParams = Partial<
      Pick<ProcessListItem, 'id' | 'code' | 'remark'> & Api.Common.CommonSearchParams
    > & {
      /** 逗号分隔的标签 id */
      tagIds?: string
      /** 1/true 筛选含孤儿标签关联的流程 */
      tagOther?: string
      /** include（默认）| exclude：标签条件取反 */
      tagMode?: string
    }

    interface TagListResponse {
      rows: { id: string; name: string }[]
    }

    interface TagAddReq {
      name: string
    }

    interface TagAddRes {
      id: string
    }

    interface TagEditReq {
      id: string
      name: string
    }
  }

  /** 项目管理 dev_project */
  namespace DevProject {
    interface ProjectTagItem {
      id: string
      name: string
    }

    interface ProjectListItem {
      id: string
      name: string
      /** 0 草稿 1 进行中 2 暂停 3 已完成 */
      status: string
      /** 0-100 */
      progress: string
      versionChangelog: string
      createTime?: string
      updateTime?: string
      tags?: ProjectTagItem[]
    }

    interface ProjectListResponse {
      rows: ProjectListItem[]
      total: string
    }

    interface ProjectDetail {
      id: string
      name: string
      status: string
      progress: string
      versionChangelog: string
      mindJson: string
      createTime?: string
      createBy?: string
      updateTime?: string
      updateBy?: string
      tags?: ProjectTagItem[]
    }

    interface ProjectAddReq {
      name: string
      status?: number
      progress?: number
      versionChangelog?: string
      mindJson?: string
      tagIds?: number[]
    }

    interface ProjectAddRes {
      affect: string
      id: string
    }

    interface ProjectEditReq {
      id: string
      name: string
      status?: number
      progress?: number
      versionChangelog?: string
      tagIds?: number[]
    }

    interface ProjectEditMindReq {
      id: string
      mindJson: string
    }

    interface ProjectDelReq {
      ids: string[]
    }

    type ProjectSearchParams = Partial<
      Pick<ProjectListItem, 'name'> & Api.Common.CommonSearchParams
    > & {
      status?: string
      /** 逗号分隔的标签 id */
      tagIds?: string
      tagOther?: string
      tagMode?: string
    }

    interface TagListResponse {
      rows: { id: string; name: string }[]
    }

    interface TagAddReq {
      name: string
    }

    interface TagAddRes {
      id: string
    }

    interface TagEditReq {
      id: string
      name: string
    }
  }
}
