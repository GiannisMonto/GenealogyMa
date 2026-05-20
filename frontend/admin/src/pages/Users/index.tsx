import { ReactElement, useEffect, useState } from 'react';
import {
  Table,
  Button,
  Space,
  Input,
  Tag,
  Typography,
  Card,
  Modal,
  Form,
  Select,
  message,
  Popconfirm,
  Avatar,
  Dropdown,
  type MenuProps,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
  UserOutlined,
  LockOutlined,
  CheckCircleOutlined,
  StopOutlined,
  MoreOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { UserDTO } from '@shared/api/types';
import { userApi } from '@/api/user';
import styles from './index.module.css';

const { Title } = Typography;

interface UserFormData {
  username: string;
  email: string;
  display_name?: string;
  password?: string;
  role_ids?: number[];
  status?: 'active' | 'inactive' | 'banned';
}

export function Users(): ReactElement {
  const [loading, setLoading] = useState(false);
  const [users, setUsers] = useState<UserDTO[]>([]);
  const [searchText, setSearchText] = useState('');
  const [modalVisible, setModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState<UserDTO | null>(null);
  const [resetPasswordVisible, setResetPasswordVisible] = useState(false);
  const [form] = Form.useForm<UserFormData>();
  const [resetForm] = Form.useForm();
  const [messageApi, contextHolder] = message.useMessage();

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const data = await userApi.getUsers();
      setUsers(data);
    } catch {
      messageApi.error('获取用户列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const filteredUsers = users.filter(
    (user) =>
      user.username.toLowerCase().includes(searchText.toLowerCase()) ||
      user.email.toLowerCase().includes(searchText.toLowerCase()) ||
      user.display_name?.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleCreate = () => {
    setEditingUser(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (user: UserDTO) => {
    setEditingUser(user);
    form.setFieldsValue({
      username: user.username,
      email: user.email,
      display_name: user.display_name,
      status: user.status,
      role_ids: user.roles?.map((r) => r.id),
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await userApi.deleteUser(id);
      messageApi.success('删除成功');
      fetchUsers();
    } catch {
      messageApi.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingUser) {
        await userApi.updateUser(editingUser.id, {
          email: values.email,
          display_name: values.display_name,
          status: values.status,
          role_ids: values.role_ids,
        });
        messageApi.success('更新成功');
      } else {
        await userApi.createUser({
          username: values.username,
          email: values.email,
          password: values.password!,
          display_name: values.display_name,
          role_ids: values.role_ids,
        });
        messageApi.success('创建成功');
      }
      setModalVisible(false);
      fetchUsers();
    } catch {
      messageApi.error('操作失败');
    }
  };

  const handleResetPassword = async () => {
    if (!editingUser) return;
    try {
      const values = await resetForm.validateFields();
      await userApi.resetPassword(editingUser.id, { new_password: values.new_password });
      messageApi.success('密码重置成功');
      setResetPasswordVisible(false);
    } catch {
      messageApi.error('密码重置失败');
    }
  };

  const handleStatusChange = async (user: UserDTO, status: 'active' | 'inactive' | 'banned') => {
    try {
      await userApi.updateStatus(user.id, status);
      messageApi.success('状态更新成功');
      fetchUsers();
    } catch {
      messageApi.error('状态更新失败');
    }
  };

  const getStatusTag = (status: string) => {
    switch (status) {
      case 'active':
        return <Tag color="success">正常</Tag>;
      case 'inactive':
        return <Tag color="warning">未激活</Tag>;
      case 'banned':
        return <Tag color="error">已禁用</Tag>;
      default:
        return <Tag>{status}</Tag>;
    }
  };

  const getMenuItems = (user: UserDTO): MenuProps['items'] => [
    {
      key: 'edit',
      icon: <EditOutlined />,
      label: '编辑',
      onClick: () => handleEdit(user),
    },
    {
      key: 'reset',
      icon: <LockOutlined />,
      label: '重置密码',
      onClick: () => {
        setEditingUser(user);
        resetForm.setFieldsValue({ new_password: '' });
        setResetPasswordVisible(true);
      },
    },
    { type: 'divider' },
    user.status === 'active'
      ? {
          key: 'deactivate',
          icon: <StopOutlined />,
          label: '禁用',
          onClick: () => handleStatusChange(user, 'inactive'),
        }
      : {
          key: 'activate',
          icon: <CheckCircleOutlined />,
          label: '激活',
          onClick: () => handleStatusChange(user, 'active'),
        },
  ];

  const columns: ColumnsType<UserDTO> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '用户',
      key: 'user',
      width: 250,
      render: (_, record) => (
        <Space>
          <Avatar icon={<UserOutlined />} src={record.avatar_url} />
          <div>
            <div style={{ fontWeight: 500 }}>{record.display_name || record.username}</div>
            <div style={{ fontSize: 12, color: '#999' }}>@{record.username}</div>
          </div>
        </Space>
      ),
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      width: 200,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '角色',
      key: 'roles',
      width: 200,
      render: (_, record) => (
        <Space wrap>
          {record.roles?.map((role) => (
            <Tag key={role.id} color="blue">{role.name}</Tag>
          ))}
        </Space>
      ),
    },
    {
      title: '最后登录',
      dataIndex: 'last_login_at',
      width: 180,
      render: (time: string | null) => time ? new Date(time).toLocaleString('zh-CN') : '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      render: (_, record) => (
        <Space size="small">
          <Button size="small" type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Dropdown menu={{ items: getMenuItems(record) }} trigger={['click']}>
            <Button size="small" type="link" icon={<MoreOutlined />} />
          </Dropdown>
        </Space>
      ),
    },
  ];

  return (
    <div className={styles.container}>
      {contextHolder}
      <div className={styles.header}>
        <Title level={3}>用户管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增用户
        </Button>
      </div>

      <Card className={styles.filterCard}>
        <Input.Search
          placeholder="搜索用户名、邮箱或昵称..."
          allowClear
          enterButton={<SearchOutlined />}
          onSearch={setSearchText}
          style={{ width: 300 }}
        />
      </Card>

      <Table
        columns={columns}
        dataSource={filteredUsers}
        rowKey="id"
        loading={loading}
        pagination={{
          pageSize: 10,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`,
        }}
      />

      <Modal
        title={editingUser ? '编辑用户' : '新增用户'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="请输入用户名" disabled={!!editingUser} />
          </Form.Item>
          {!editingUser && (
            <Form.Item
              name="password"
              label="密码"
              rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '密码至少6位' }]}
            >
              <Input.Password placeholder="请输入密码" />
            </Form.Item>
          )}
          <Form.Item
            name="email"
            label="邮箱"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input placeholder="请输入邮箱" />
          </Form.Item>
          <Form.Item name="display_name" label="显示名称">
            <Input placeholder="请输入显示名称" />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select placeholder="选择状态">
              <Select.Option value="active">正常</Select.Option>
              <Select.Option value="inactive">未激活</Select.Option>
              <Select.Option value="banned">禁用</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="重置密码"
        open={resetPasswordVisible}
        onOk={handleResetPassword}
        onCancel={() => setResetPasswordVisible(false)}
      >
        <Form form={resetForm} layout="vertical">
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[{ required: true, message: '请输入新密码' }, { min: 6, message: '密码至少6位' }]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}