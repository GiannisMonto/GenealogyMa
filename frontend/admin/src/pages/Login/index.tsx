import { ReactElement, useState } from 'react';
import { Form, Input, Button, Card, Typography, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';
import styles from './index.module.css';

const { Title, Text } = Typography;

interface LoginForm {
  username: string;
  password: string;
}

export function Login(): ReactElement {
  const navigate = useNavigate();
  const { login, isLoading } = useAuthStore();
  const [messageApi, contextHolder] = message.useMessage();

  const onFinish = async (values: LoginForm) => {
    try {
      await login(values.username, values.password);
      messageApi.success('登录成功');
      navigate('/dashboard');
    } catch (error) {
      messageApi.error('用户名或密码错误');
    }
  };

  return (
    <div className={styles.container}>
      {contextHolder}
      <Card className={styles.card}>
        <div className={styles.header}>
          <Title level={3}>族谱管理平台</Title>
          <Text type="secondary">请登录您的账号</Text>
        </div>
        <Form name="login" onFinish={onFinish} layout="vertical" size="large">
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={isLoading}>
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}