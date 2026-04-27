import { AppRouteRecord } from '@/types/router'

export const devopsRoutes: AppRouteRecord = {
  name: 'Devops',
  path: '/devops',
  component: '/index/index',
  meta: {
    title: 'menus.devops.title',
    icon: 'ri:kanban-view',
    menus: ['R_SUPER', 'R_ADMIN', 'devops']
  },
  children: [
    {
      path: 'process',
      name: 'DevopsProcessList',
      component: '/devops/automation/process',
      meta: {
        title: 'menus.devops.processList',
        icon: 'ri:list-check-2',
        keepAlive: true,
        menus: ['R_SUPER', 'R_ADMIN', 'dev:process']
      }
    },
    {
      path: 'process/edit/:id',
      name: 'DevopsProcessEdit',
      component: '/devops/automation/process-edit',
      meta: {
        title: 'menus.devops.processEdit',
        isHide: true,
        keepAlive: false,
        activePath: '/devops/process'
      }
    },
    {
      path: 'project',
      name: 'DevopsProjectList',
      component: '/devops/project/index',
      meta: {
        title: 'menus.devops.projectList',
        icon: 'ri:folder-chart-line',
        keepAlive: true,
        menus: ['R_SUPER', 'R_ADMIN', 'dev:project']
      }
    },
    {
      path: 'project/mind/:id',
      name: 'DevopsProjectMind',
      component: '/devops/project/mind',
      meta: {
        title: 'menus.devops.projectMind',
        isHide: true,
        keepAlive: false,
        activePath: '/devops/project'
      }
    }
  ]
}
