import { AppRouteRecord } from '@/types/router'

export const systemRoutes: AppRouteRecord = {
  path: '/system',
  name: 'System',
  component: '/index/index',
  meta: {
    title: 'menus.system.title',
    icon: 'ri:settings-3-line',
    menus: ['R_SUPER', 'R_ADMIN', 'sys']
  },
  children: [
    {
      path: 'user',
      name: 'User',
      component: '/system/user',
      meta: {
        title: 'menus.system.user',
        icon: 'ri:user-3-line',
        keepAlive: true,
        menus: ['R_SUPER', 'R_ADMIN', 'sys:user'],
        authList: [
          { title: '用户列表', authMark: 'sys:user:list' },
          { title: '用户详情', authMark: 'sys:user:get' },
          { title: '新增用户', authMark: 'sys:user:add' },
          { title: '编辑用户', authMark: 'sys:user:edit' },
          { title: '删除用户', authMark: 'sys:user:del' },
          { title: '导入用户', authMark: 'sys:user:import' },
          { title: '导出用户', authMark: 'sys:user:export' },
          { title: '绑定角色', authMark: 'sys:user:assign_role' },
          { title: '绑定部门', authMark: 'sys:user:assign_dept' },
          { title: '绑定岗位', authMark: 'sys:user:assign_post' }
        ]
      }
    },
    {
      path: 'role',
      name: 'Role',
      component: '/system/role',
      meta: {
        title: 'menus.system.role',
        icon: 'ri:shield-user-line',
        keepAlive: true,
        menus: ['R_SUPER', 'sys:role'],
        authList: [
          { title: '角色列表', authMark: 'sys:role:list' },
          { title: '角色详情', authMark: 'sys:role:get' },
          { title: '新增角色', authMark: 'sys:role:add' },
          { title: '编辑角色', authMark: 'sys:role:edit' },
          { title: '删除角色', authMark: 'sys:role:del' },
          { title: '操作权限', authMark: 'sys:role:assign_menu' },
          { title: '数据权限', authMark: 'sys:role:assign_dept' }
        ]
      }
    },
    {
      path: 'user-center',
      name: 'UserCenter',
      component: '/system/user-center',
      meta: {
        title: 'menus.system.userCenter',
        isHide: true,
        keepAlive: true,
        isHideTab: true
      }
    },
    {
      path: 'dept',
      name: 'Dept',
      component: '/system/dept',
      meta: {
        title: 'menus.system.dept',
        icon: 'ri:organization-chart',
        keepAlive: true,
        menus: ['R_SUPER', 'sys:dept'],
        authList: [
          { title: '部门列表', authMark: 'sys:dept:list' },
          { title: '部门详情', authMark: 'sys:dept:get' },
          { title: '新增部门', authMark: 'sys:dept:add' },
          { title: '编辑部门', authMark: 'sys:dept:edit' },
          { title: '删除部门', authMark: 'sys:dept:del' }
        ]
      }
    },
    {
      path: 'post',
      name: 'Post',
      component: '/system/post',
      meta: {
        title: 'menus.system.post',
        icon: 'ri:briefcase-2-line',
        keepAlive: true,
        menus: ['R_SUPER', 'sys:post'],
        authList: [
          { title: '岗位列表', authMark: 'sys:post:list' },
          { title: '岗位详情', authMark: 'sys:post:get' },
          { title: '新增岗位', authMark: 'sys:post:add' },
          { title: '编辑岗位', authMark: 'sys:post:edit' },
          { title: '删除岗位', authMark: 'sys:post:del' }
        ]
      }
    }
    // {
    //   path: 'menu',
    //   name: 'Menus',
    //   component: '/system/menu',
    //   meta: {
    //     title: 'menus.system.menu',
    //     keepAlive: true,
    //     menus: ['R_SUPER'],
    //     authList: [
    //       { title: '新增', authMark: 'add' },
    //       { title: '编辑', authMark: 'edit' },
    //       { title: '删除', authMark: 'delete' }
    //     ]
    //   }
    // }
  ]
}
