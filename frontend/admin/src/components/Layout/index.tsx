import { ReactElement, useState, useEffect } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import { Layout, Menu, Avatar, Dropdown, Button, Breadcrumb } from 'antd';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  TeamOutlined,
  UserOutlined,
  SafetyOutlined,
  FileTextOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { useAuthStore } from '@/store/auth';
import { useAppStore } from '@/store/app';
import { routes, getFlattenRoutes } from '@/routes/routes';
import styles from './index.module.css';

const { Header, Sider, Content } = Layout;

const iconMap: Record<string, ReactElement> = {
  DashboardOutlined: <DashboardOutlined />,
  TeamOutlined: <TeamOutlined />,
  UserOutlined: <UserOutlined />,
  SafetyOutlined: <SafetyOutlined />,
  FileTextOutlined: <FileTextOutlined />,
};

export function AppLayout(): ReactElement {
  const location = useLocation();
  const { user, logout } = useAuthStore();
  const { sidebarCollapsed, toggleSidebar, setBreadcrumbs } = useAppStore();
  const [selectedKeys, setSelectedKeys] = useState<string[]>(['/dashboard']);

  useEffect(() => {
    const path = location.pathname;
    const currentRoute = getFlattenRoutes(routes).find((r) => r.path === path);
    if (currentRoute) {
      setSelectedKeys([path]);
      const crumbs = [{ path: '/', name: '首页' }];
      if (path !== '/' && path !== '/dashboard') {
        crumbs.push({ path, name: currentRoute.name });
      }
      setBreadcrumbs(crumbs);
    }
  }, [location.pathname, setBreadcrumbs]);

  const menuItems: MenuProps['items'] = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: '仪表盘',
    },
    {
      key: '/members',
      icon: <TeamOutlined />,
      label: '成员管理',
    },
    {
      key: '/users',
      icon: <UserOutlined />,
      label: '用户管理',
    },
    {
      key: '/roles',
      icon: <SafetyOutlined />,
      label: '角色管理',
    },
    {
      key: '/audit',
      icon: <FileTextOutlined />,
      label: '审计日志',
    },
  ];

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      danger: true,
    },
  ];

  const handleMenuClick: MenuProps['onClick'] = (e) => {
    if (e.key === 'logout') {
      logout();
    }
  };

  return (
    <Layout className={styles.layout}>
      <Sider trigger={null} collapsible collapsed={sidebarCollapsed} className={styles.sider}>
        <div className={styles.logo}>
          {sidebarCollapsed ? 'G' : '族谱管理'}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={selectedKeys}
          items={menuItems}
          onClick={({ key }) => window.location.href = key}
        />
      </Sider>
      <Layout>
        <Header className={styles.header}>
          <Button
            type="text"
            icon={sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={toggleSidebar}
            className={styles.trigger}
          />
          <div className={styles.headerRight}>
            <Dropdown menu={{ items: userMenuItems, onClick: handleMenuClick }} placement="bottomRight">
              <Avatar icon={<UserOutlined />} className={styles.avatar}>
                {user?.username?.[0]?.toUpperCase()}
              </Avatar>
            </Dropdown>
          </div>
        </Header>
        <Content className={styles.content}>
          <Breadcrumb
            className={styles.breadcrumb}
            items={useAppStore.getState().breadcrumbs.map((b) => ({ title: b.name }))}
          />
          <div className={styles.main}>
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  );
}