import { AppRouteRecord } from '@/types/router'
import { dashboardRoutes } from './dashboard'
import { systemRoutes } from './system'
import { safeguardRoutes } from '../backend/safeguard'

export const routeModules: AppRouteRecord[] = [
  dashboardRoutes,
  systemRoutes,
  safeguardRoutes
]
