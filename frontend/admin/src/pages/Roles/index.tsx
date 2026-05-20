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
  Drawer,
  Transfer,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SettingOutlined,
  SearchOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { RoleDTO, PermissionDTO } from '@shared/api/types';
import { roleApi, permissionApi } from '@/api/rbac';
import styles from './index.module.css';

const { Title } = Typography;

interface RoleFormData {
  code: string;
  name: string;
  description?: string;
}

export function Roles(): ReactElement {
  const [loading, setLoading] = useState(false);
  const [roles, setRoles] = useState<RoleDTO[]>([]);
  const [permissions, setPermissions] = useState<PermissionDTO[]>([]);
  const [searchText, setSearchText] = useState('');
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRole, setEditingRole] = useState<RoleDTO | null>(null);
  const [permissionDrawerVisible, setPermissionDrawerVisible] = useState(false);
  const [currentRolePermissions, setCurrentRolePermissions] = useState<PermissionDTO[]>([]);
  const [form] = Form.useForm<RoleFormData>();
  const [messageApi, contextHolder] = message.useMessage();

  const fetchRoles = async () => {
    setLoading(true);
    try {
      const data = await roleApi.getRoles();
      setRoles(data);
    } catch {
      messageApi.error('获取角色列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchPermissions = async () => {
    try {
      const data = await permissionApi.getPermissions();
      setPermissions(data);
    } catch {
      messageApi.error('获取权限列表失败');
    }
  };

  useEffect(() => {
    fetchRoles();
    fetchPermissions();
  }, []);

  const filteredRoles = roles.filter(
    (role) =>
      role.name.toLowerCase().includes(searchText.toLowerCase()) ||
      role.code.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleCreate = () => {
    setEditingRole(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (role: RoleDTO) => {
    setEditingRole(role);
    form.setFieldsValue({
      code: role.code,
      name: role.name,
      description: role.description,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await roleApi.deleteRole(id);
      messageApi.success('删除成功');
      fetchRoles();
    } catch {
      messageApi.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingRole) {
        await roleApi.updateRole(editingRole.id, values);
        messageApi.success('更新成功');
      } else {
        await roleApi.createRole(values);
        messageApi.success('创建成功');
      }
      setModalVisible(false);
      fetchRoles();
    } catch {
      messageApi.error('操作失败');
    }
  };

  const handlePermissionSettings = async (role: RoleDTO) => {
    try {
      const perms = await roleApi.getRolePermissions(role.id);
      setCurrentRolePermissions(perms);
      setEditingRole(role);
      setPermissionDrawerVisible(true);
    } catch {
      messageApi.error('获取权限列表失败');
    }
  };

  const handlePermissionChange = async (targetKeys: string[]) => {
    if (!editingRole) return;
    try {
      await roleApi.assignPermissions(editingRole.id, targetKeys.map(Number));
      messageApi.success('权限分配成功');
      setPermissionDrawerVisible(false);
    } catch {
      messageApi.error('权限分配失败');
    }
  };

  const columns: ColumnsType<RoleDTO> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '角色代码',
      dataIndex: 'code',
      width: 150,
      render: (code: string) => <Tag color="blue">{code}</Tag>,
    },
    {
      title: '角色名称',
      dataIndex: 'name',
      width: 150,
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
    },
    {
      title: '系统角色',
      dataIndex: 'is_system',
      width: 100,
      render: (isSystem: boolean) =>
        isSystem ? <Tag color="green">是</Tag> : <Tag>否</Tag>,
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: unknown, record: RoleDTO) => (
        <Space size="small">
          <Button
            size="small"
            type="link"
            icon={<SettingOutlined />}
            onClick={() => handlePermissionSettings(record)}
          >
            权限
          </Button>
          <Button
            size="small"
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          {!record.is_system && (
            <Popconfirm
              title="确定删除此角色？"
              onConfirm={() => handleDelete(record.id)}
            >
              <Button size="small" type="link" danger icon={<DeleteOutlined />}>
                删除
              </Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  const permissionDataSource = permissions.map((p) => ({
    key: String(p.id),
    title: `${p.name} (${p.code})`,
    description: p.description,
  }));

  const assignedKeys = currentRolePermissions.map((p) => String(p.id));

  return (
    <div className={styles.container}>
      {contextHolder}
      <div className={styles.header}>
        <Title level={3}>角色管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增角色
        </Button>
      </div>

      <Card className={styles.filterCard}>
        <Input.Search
          placeholder="搜索角色名称或代码..."
          allowClear
          enterButton={<SearchOutlined />}
          onSearch={setSearchText}
          style={{ width: 300 }}
        />
      </Card>

      <Table
        columns={columns}
        dataSource={filteredRoles}
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
        title={editingRole ? '编辑角色' : '新增角色'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item
            name="code"
            label="角色代码"
            rules={[{ required: true, message: '请输入角色代码' }]}
          >
            <Input placeholder="如: role_admin" disabled={!!editingRole} />
          </Form.Item>
          <Form.Item
            name="name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          >
            <Input placeholder="如: 管理员" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea placeholder="角色描述..." rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      <Drawer
        title="分配权限"
        open={permissionDrawerVisible}
        onClose={() => setPermissionDrawerVisible(false)}
        width={500}
      >
        <Transfer
          dataSource={permissionDataSource}
          titles={['可分配权限', '已拥有权限']}
          targetKeys={assignedKeys}
          render={(item) => item.title || ''}
          onChange={handlePermissionChange}
          listStyle={{ width: 400, height: 400 }}
        />
      </Drawer>
    </div>
  );
}