import { ReactElement, useState, useRef, useEffect } from 'react';
import { Form, Input, Button, Card, Typography, message, Checkbox } from 'antd';
import { UserOutlined, LockOutlined, CloseCircleOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';
import styles from './index.module.css';

const { Title, Text } = Typography;

interface LoginForm {
  username: string;
  password: string;
  captcha: string;
}

interface CaptchaState {
  text: string;
  dataUrl: string;
}

// Simple SVG captcha generator (no backend needed)
const generateCaptcha = (): CaptchaState => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789';
  const length = 4;
  let text = '';
  for (let i = 0; i < length; i++) {
    text += chars.charAt(Math.floor(Math.random() * chars.length));
  }

  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="120" height="40" style="background:#f0f2f5">
      <text x="10" y="28" font-family="Courier New" font-size="24" fill="#1890ff" font-weight="bold" transform="rotate(${-5 + Math.random() * 10} 60 20)">${text.split('').join(' ')}</text>
      ${Array.from({ length: 6 }, (_, i) => `<line x1="${Math.random() * 120}" y1="${Math.random() * 40}" x2="${Math.random() * 120}" y2="${Math.random() * 40}" stroke="#ccc" stroke-width="1" opacity="${0.3 + Math.random() * 0.3}"/>`).join('')}
    </svg>
  `;
  return { text, dataUrl: `data:image/svg+xml;base64,${btoa(svg)}` };
};

export function Login(): ReactElement {
  const navigate = useNavigate();
  const { login, isLoading } = useAuthStore();
  const [messageApi, contextHolder] = message.useMessage();
  const [captcha, setCaptcha] = useState<CaptchaState>({ text: '', dataUrl: '' });
  const [captchaInput, setCaptchaInput] = useState('');
  const [localError, setLocalError] = useState('');
  const [form] = Form.useForm();
  const captchaRef = useRef<HTMLImageElement>(null);

  useEffect(() => {
    refreshCaptcha();
  }, []);

  const refreshCaptcha = () => {
    setCaptcha(generateCaptcha());
    setCaptchaInput('');
    setLocalError('');
  };

  const onFinish = async (values: LoginForm) => {
    setLocalError('');

    // Client-side captcha validation
    if (captchaInput.toLowerCase() !== captcha.text.toLowerCase()) {
      setLocalError('验证码错误');
      refreshCaptcha();
      return;
    }

    try {
      await login(values.username, values.password);
      messageApi.success('登录成功');
      navigate('/dashboard');
    } catch {
      messageApi.error('用户名或密码错误');
      refreshCaptcha();
    }
  };

  const onCaptchaChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCaptchaInput(e.target.value);
    if (localError) setLocalError('');
  };

  return (
    <div className={styles.container}>
      {contextHolder}
      <Card className={styles.card}>
        <div className={styles.header}>
          <Title level={3}>族谱管理平台</Title>
          <Text type="secondary">请登录您的账号</Text>
        </div>
        <Form form={form} name="login" onFinish={onFinish} layout="vertical" size="large">
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名" autoComplete="username" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" autoComplete="current-password" />
          </Form.Item>
          <Form.Item
            name="captcha"
            validateStatus={localError ? 'error' : ''}
            help={localError}
          >
            <div className={styles.captchaRow}>
              <Input
                prefix={<LockOutlined />}
                placeholder="验证码"
                value={captchaInput}
                onChange={onCaptchaChange}
                maxLength={4}
                style={{ flex: 1 }}
              />
              <img
                ref={captchaRef}
                src={captcha.dataUrl}
                alt="验证码"
                className={styles.captchaImage}
                onClick={refreshCaptcha}
                title="点击刷新"
              />
              <Button type="text" onClick={refreshCaptcha} icon={<CloseCircleOutlined />} />
            </div>
          </Form.Item>
          <Form.Item>
            <div className={styles.optionsRow}>
              <Checkbox>记住我</Checkbox>
              <Link to="/forgot-password" className={styles.forgotLink}>忘记密码？</Link>
            </div>
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