import { AppRouteRecord } from '@/types/router'

export const dashboardRoutes: AppRouteRecord = {
  name: 'Dashboard',
  path: '/dashboard',
  component: '/index/index',
  meta: {
    title: 'menus.dashboard.title',
    icon: 'ri:dashboard-2-line',
    menus: ['R_SUPER', 'R_ADMIN', 'dashboard']
  },
  children: [
    {
      path: 'console',
      name: 'Console',
      component: '/dashboard/console',
      meta: {
        title: 'menus.dashboard.console',
        icon: 'ri:layout-grid-line',
        keepAlive: false,
        fixedTab: true,
        menus: ['dashboard:workbench']
      }
    }
  ]
}
