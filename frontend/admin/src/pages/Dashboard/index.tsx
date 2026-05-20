import { ReactElement } from 'react';
import { Card, Typography, Row, Col, Statistic } from 'antd';
import { TeamOutlined, UserOutlined, SafetyOutlined, FileTextOutlined } from '@ant-design/icons';
import styles from './index.module.css';

const { Title } = Typography;

export function Dashboard(): ReactElement {
  return (
    <div className={styles.container}>
      <Title level={3}>控制台</Title>
      <Row gutter={16}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总成员数"
              value={5950}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="用户数"
              value={128}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="角色数"
              value={6}
              prefix={<SafetyOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="审计日志"
              value={1024}
              prefix={<FileTextOutlined />}
              valueStyle={{ color: '#eb2f96' }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
}