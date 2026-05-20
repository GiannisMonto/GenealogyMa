import { ReactElement } from 'react';
import { Typography } from 'antd';

const { Title } = Typography;

export function MembersDetail(): ReactElement {
  return (
    <div>
      <Title level={3}>成员详情</Title>
      <p>成员详情页面开发中...</p>
    </div>
  );
}