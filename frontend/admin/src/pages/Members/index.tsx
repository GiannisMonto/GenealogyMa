import { ReactElement } from 'react';
import { Typography } from 'antd';

const { Title } = Typography;

export function Members(): ReactElement {
  return (
    <div>
      <Title level={3}>成员管理</Title>
      <p>成员列表页面开发中...</p>
    </div>
  );
}