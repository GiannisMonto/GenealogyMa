import { ReactElement } from 'react';
import { Typography } from 'antd';

const { Title } = Typography;

export function Users(): ReactElement {
  return (
    <div>
      <Title level={3}>用户管理</Title>
      <p>用户管理页面开发中...</p>
    </div>
  );
}