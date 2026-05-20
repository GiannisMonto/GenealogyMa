import { ReactElement } from 'react';
import { Typography } from 'antd';

const { Title } = Typography;

export function Roles(): ReactElement {
  return (
    <div>
      <Title level={3}>角色管理</Title>
      <p>角色管理页面开发中...</p>
    </div>
  );
}