import { ReactElement } from 'react';
import { Typography } from 'antd';

const { Title } = Typography;

export function Audit(): ReactElement {
  return (
    <div>
      <Title level={3}>审计日志</Title>
      <p>审计日志页面开发中...</p>
    </div>
  );
}