import React, { ReactElement, lazy } from 'react';

export interface RouteConfig {
  path: string;
  name: string;
  element: ReactElement;
  children?: RouteConfig[];
  meta?: {
    requiresAuth?: boolean;
    permission?: string;
    icon?: string;
  };
}

export const routes: RouteConfig[] = [
  {
    path: '/login',
    name: '登录',
    element: React.createElement(lazy(() => import('@/pages/Login'))),
    meta: { requiresAuth: false },
  },
  {
    path: '/forgot-password',
    name: '忘记密码',
    element: React.createElement(lazy(() => import('@/pages/ForgotPassword'))),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    name: '控制台',
    element: React.createElement(lazy(() => import('@/pages/Dashboard'))),
    meta: { requiresAuth: true },
  },
  {
    path: '/dashboard',
    name: '仪表盘',
    element: React.createElement(lazy(() => import('@/pages/Dashboard'))),
    meta: { requiresAuth: true },
  },
  {
    path: '/members',
    name: '成员管理',
    element: React.createElement(lazy(() => import('@/pages/Members'))),
    meta: { requiresAuth: true, permission: 'person:read' },
  },
  {
    path: '/members/:id',
    name: '成员详情',
    element: React.createElement(lazy(() => import('@/pages/Members/Detail'))),
    meta: { requiresAuth: true, permission: 'person:read' },
  },
  {
    path: '/members/:id/edit',
    name: '编辑成员',
    element: React.createElement(lazy(() => import('@/pages/Members/Form'))),
    meta: { requiresAuth: true, permission: 'person:write' },
  },
  {
    path: '/members/new',
    name: '新增成员',
    element: React.createElement(lazy(() => import('@/pages/Members/Form'))),
    meta: { requiresAuth: true, permission: 'person:write' },
  },
  {
    path: '/users',
    name: '用户管理',
    element: React.createElement(lazy(() => import('@/pages/Users'))),
    meta: { requiresAuth: true, permission: 'admin:user' },
  },
  {
    path: '/roles',
    name: '角色管理',
    element: React.createElement(lazy(() => import('@/pages/Roles'))),
    meta: { requiresAuth: true, permission: 'admin:role' },
  },
  {
    path: '/audit',
    name: '审计日志',
    element: React.createElement(lazy(() => import('@/pages/Audit'))),
    meta: { requiresAuth: true, permission: 'admin:audit' },
  },
  {
    path: '/settings',
    name: '系统设置',
    element: React.createElement(lazy(() => import('@/pages/Settings'))),
    meta: { requiresAuth: true, permission: 'admin:settings' },
  },
  {
    path: '/permissions',
    name: '权限管理',
    element: React.createElement(lazy(() => import('@/pages/Permissions'))),
    meta: { requiresAuth: true, permission: 'admin:permission' },
  },
  {
    path: '/cemetery',
    name: '墓园管理',
    element: React.createElement(lazy(() => import('@/pages/Cemetery'))),
    meta: { requiresAuth: true, permission: 'cemetery:read' },
  },
  {
    path: '/culture',
    name: '文化管理',
    element: React.createElement(lazy(() => import('@/pages/Culture'))),
    meta: { requiresAuth: true, permission: 'culture:read' },
  },
  {
    path: '*',
    name: '404',
    element: React.createElement(lazy(() => import('@/pages/Error/404'))),
  },
];

export const getFlattenRoutes = (routeList: RouteConfig[]): RouteConfig[] => {
  const result: RouteConfig[] = [];

  const flatten = (routes: RouteConfig[]) => {
    routes.forEach((route) => {
      result.push(route);
      if (route.children) {
        flatten(route.children);
      }
    });
  };

  flatten(routeList);
  return result;
};