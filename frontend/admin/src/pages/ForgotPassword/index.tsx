import { ReactElement, useState } from 'react';
import { Form, Input, Button, Card, Typography, message, Checkbox } from 'antd';
import { MailOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import styles from './index.module.css';

const { Title, Text } = Typography;

interface ForgotPasswordForm {
  email: string;
}

export function ForgotPassword(): ReactElement {
  const navigate = useNavigate();
  const [messageApi, contextHolder] = message.useMessage();
  const [isLoading, setIsLoading] = useState(false);
  const [emailSent, setEmailSent] = useState(false);
  const [form] = Form.useForm();

  const onFinish = async (values: ForgotPasswordForm) => {
    setIsLoading(true);
    try {
      // Simulate API call - in production this would call the backend
      await new Promise((resolve) => setTimeout(resolve, 1500));
      console.log('Password reset requested for:', values.email);
      setEmailSent(true);
      messageApi.success('重置链接已发送到您的邮箱');
    } catch {
      messageApi.error('发送失败，请稍后重试');
    } finally {
      setIsLoading(false);
    }
  };

  if (emailSent) {
    return (
      <div className={styles.container}>
        {contextHolder}
        <Card className={styles.card}>
          <div className={styles.successContent}>
            <Title level={4}>邮件已发送</Title>
            <Text type="secondary">
              我们已向您的邮箱发送了密码重置链接，请查收。
            </Text>
            <Button type="primary" onClick={() => navigate('/login')} style={{ marginTop: 24 }}>
              返回登录
            </Button>
          </div>
        </Card>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {contextHolder}
      <Card className={styles.card}>
        <div className={styles.header}>
          <Title level={3}>忘记密码</Title>
          <Text type="secondary">输入您的注册邮箱，我们会发送重置链接</Text>
        </div>
        <Form form={form} name="forgot-password" onFinish={onFinish} layout="vertical" size="large">
          <Form.Item
            name="email"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input prefix={<MailOutlined />} placeholder="注册邮箱" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={isLoading}>
              发送重置链接
            </Button>
          </Form.Item>
          <Form.Item>
            <div className={styles.backLogin}>
              想起密码了？<Link to="/login">返回登录</Link>
            </div>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}