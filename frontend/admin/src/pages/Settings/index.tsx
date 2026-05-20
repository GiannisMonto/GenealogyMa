import { ReactElement, useState } from 'react';
import {
  Card,
  Form,
  Input,
  Switch,
  Button,
  Space,
  Typography,
  message,
  Select,
  InputNumber,
  Divider,
} from 'antd';
import {
  SaveOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import styles from './index.module.css';

const { Title, Text } = Typography;
const { TextArea } = Input;

interface SystemConfigFormData {
  siteName: string;
  siteDescription: string;
  maintenanceMode: boolean;
  allowRegister: boolean;
  defaultRole: string;
  uploadMaxSize: number;
  uploadAllowedTypes: string;
  logLevel: string;
  sessionExpireHours: number;
}

export function Settings(): ReactElement {
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm<SystemConfigFormData>();
  const [messageApi, contextHolder] = message.useMessage();

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      // In real implementation, this would call an API to save settings
      console.log('Saving settings:', values);
      messageApi.success('设置保存成功');
    } catch {
      messageApi.error('请检查表单填写');
    }
  };

  const handleReset = () => {
    form.resetFields();
    messageApi.info('已重置为默认值');
  };

  return (
    <div className={styles.container}>
      {contextHolder}
      <div className={styles.header}>
        <Title level={3}>系统设置</Title>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={handleReset}>
            重置
          </Button>
          <Button type="primary" icon={<SaveOutlined />} onClick={handleSubmit}>
            保存设置
          </Button>
        </Space>
      </div>

      <div className={styles.grid}>
        <Card title="基本信息" className={styles.card}>
          <Form
            form={form}
            layout="vertical"
            initialValues={{
              siteName: '族谱数字化管理平台',
              siteDescription: '为姓氏宗族提供完整的数字化族谱管理解决方案',
              maintenanceMode: false,
              allowRegister: true,
              defaultRole: 'verified_user',
            }}
          >
            <Form.Item
              name="siteName"
              label="网站名称"
              rules={[{ required: true, message: '请输入网站名称' }]}
            >
              <Input placeholder="请输入网站名称" />
            </Form.Item>

            <Form.Item
              name="siteDescription"
              label="网站描述"
              rules={[{ required: true, message: '请输入网站描述' }]}
            >
              <TextArea placeholder="请输入网站描述" rows={3} />
            </Form.Item>

            <Form.Item
              name="defaultRole"
              label="新用户默认角色"
              rules={[{ required: true, message: '请选择默认角色' }]}
            >
              <Select
                options={[
                  { value: 'guest', label: '游客' },
                  { value: 'verified_user', label: '认证用户' },
                ]}
              />
            </Form.Item>
          </Form>
        </Card>

        <Card title="功能开关" className={styles.card}>
          <div className={styles.switchList}>
            <div className={styles.switchItem}>
              <div>
                <Text strong>维护模式</Text>
                <br />
                <Text type="secondary">开启后仅管理员可访问</Text>
              </div>
              <Form.Item name="maintenanceMode" valuePropName="checked" noStyle>
                <Switch />
              </Form.Item>
            </div>

            <Divider />

            <div className={styles.switchItem}>
              <div>
                <Text strong>开放注册</Text>
                <br />
                <Text type="secondary">允许新用户注册账号</Text>
              </div>
              <Form.Item name="allowRegister" valuePropName="checked" noStyle>
                <Switch />
              </Form.Item>
            </div>
          </div>
        </Card>

        <Card title="上传设置" className={styles.card}>
          <Form
            form={form}
            layout="vertical"
            initialValues={{
              uploadMaxSize: 10,
              uploadAllowedTypes: 'jpg,jpeg,png,gif,pdf',
            }}
          >
            <Form.Item
              name="uploadMaxSize"
              label="最大上传大小 (MB)"
              rules={[{ required: true, message: '请输入最大上传大小' }]}
            >
              <InputNumber min={1} max={100} />
            </Form.Item>

            <Form.Item
              name="uploadAllowedTypes"
              label="允许的文件类型"
              rules={[{ required: true, message: '请输入允许的文件类型' }]}
            >
              <Input placeholder="jpg,jpeg,png,gif,pdf" />
            </Form.Item>
          </Form>
        </Card>

        <Card title="日志设置" className={styles.card}>
          <Form
            form={form}
            layout="vertical"
            initialValues={{
              logLevel: 'info',
              sessionExpireHours: 24,
            }}
          >
            <Form.Item
              name="logLevel"
              label="日志级别"
              rules={[{ required: true, message: '请选择日志级别' }]}
            >
              <Select
                options={[
                  { value: 'debug', label: 'Debug' },
                  { value: 'info', label: 'Info' },
                  { value: 'warn', label: 'Warn' },
                  { value: 'error', label: 'Error' },
                ]}
              />
            </Form.Item>

            <Form.Item
              name="sessionExpireHours"
              label="会话过期时间 (小时)"
              rules={[{ required: true, message: '请输入会话过期时间' }]}
            >
              <InputNumber min={1} max={168} />
            </Form.Item>
          </Form>
        </Card>
      </div>
    </div>
  );
}