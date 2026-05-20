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
  Tree,
  type DataNode,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
  FolderOutlined,
  SafetyOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { PermissionDTO } from '@shared/api/types';
import { permissionApi } from '@/api/rbac';
import styles from './index.module.css';

const { Title } = Typography;

interface PermissionFormData {
  code: string;
  name: string;
  module: string;
  description?: string;
}

export function Permissions(): ReactElement {
  const [loading, setLoading] = useState(false);
  const [permissions, setPermissions] = useState<PermissionDTO[]>([]);
  const [searchText, setSearchText] = useState('');
  const [modalVisible, setModalVisible] = useState(false);
  const [editingPermission, setEditingPermission] = useState<PermissionDTO | null>(null);
  const [form] = Form.useForm<PermissionFormData>();
  const [messageApi, contextHolder] = message.useMessage();

  const fetchPermissions = async () => {
    setLoading(true);
    try {
      const data = await permissionApi.getPermissions();
      setPermissions(data);
    } catch {
      messageApi.error('获取权限列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPermissions();
  }, []);

  const filteredPermissions = permissions.filter(
    (perm) =>
      perm.name.toLowerCase().includes(searchText.toLowerCase()) ||
      perm.code.toLowerCase().includes(searchText.toLowerCase()) ||
      perm.module.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleCreate = () => {
    setEditingPermission(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (perm: PermissionDTO) => {
    setEditingPermission(perm);
    form.setFieldsValue({
      code: perm.code,
      name: perm.name,
      module: perm.module,
      description: perm.description,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await permissionApi.deletePermission(id);
      messageApi.success('删除成功');
      fetchPermissions();
    } catch {
      messageApi.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingPermission) {
        await permissionApi.updatePermission(editingPermission.id, {
          name: values.name,
          description: values.description,
        });
        messageApi.success('更新成功');
      } else {
        await permissionApi.createPermission({
          code: values.code,
          name: values.name,
          module: values.module,
          description: values.description,
        });
        messageApi.success('创建成功');
      }
      setModalVisible(false);
      fetchPermissions();
    } catch {
      messageApi.error('操作失败');
    }
  };

  const getModuleColor = (module: string): string => {
    const colors: Record<string, string> = {
      admin: 'red',
      person: 'blue',
      genealogy: 'green',
      culture: 'purple',
      memorial: 'orange',
      cemetery: 'cyan',
      community: 'magenta',
    };
    return colors[module] || 'default';
  };

  const getModuleTreeData = (): DataNode[] => {
    const modules = [...new Set(permissions.map((p) => p.module))];
    return modules.map((module) => ({
      title: module.toUpperCase(),
      key: module,
      icon: <FolderOutlined />,
      children: permissions
        .filter((p) => p.module === module)
        .map((perm) => ({
          title: perm.name,
          key: perm.id.toString(),
        })),
    }));
  };

  const columns: ColumnsType<PermissionDTO> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '权限名称',
      dataIndex: 'name',
      width: 200,
      render: (name: string) => <span style={{ fontWeight: 500 }}>{name}</span>,
    },
    {
      title: '权限代码',
      dataIndex: 'code',
      width: 180,
      render: (code: string) => (
        <Tag color="geekblue" style={{ fontFamily: 'monospace' }}>
          {code}
        </Tag>
      ),
    },
    {
      title: '模块',
      dataIndex: 'module',
      width: 120,
      render: (module: string) => (
        <Tag color={getModuleColor(module)}>{module}</Tag>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      render: (_, record) => (
        <Space size="small">
          <Button size="small" type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确认删除"
            description="删除后无法恢复，确定要删除吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确认"
            cancelText="取消"
          >
            <Button size="small" type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className={styles.container}>
      {contextHolder}
      <div className={styles.header}>
        <Title level={3}>权限管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增权限
        </Button>
      </div>

      <div className={styles.content}>
        <Card className={styles.treeCard} title="权限分组">
          <Tree
            showIcon
            showLine={{ showLeafIcon: false }}
            treeData={getModuleTreeData()}
            style={{ padding: '8px 0' }}
          />
        </Card>

        <Card className={styles.tableCard}>
          <div className={styles.searchBar}>
            <Input.Search
              placeholder="搜索权限名称、代码或模块..."
              allowClear
              enterButton={<SearchOutlined />}
              onSearch={setSearchText}
              style={{ width: 300 }}
            />
          </div>

          <Table
            columns={columns}
            dataSource={filteredPermissions}
            rowKey="id"
            loading={loading}
            pagination={{
              pageSize: 10,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total) => `共 ${total} 条`,
            }}
          />
        </Card>
      </div>

      <Modal
        title={editingPermission ? '编辑权限' : '新增权限'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          {!editingPermission && (
            <Form.Item
              name="code"
              label="权限代码"
              rules={[
                { required: true, message: '请输入权限代码' },
                { pattern: /^[a-z_]+$/, message: '权限代码只能包含小写字母和下划线' },
              ]}
            >
              <Input placeholder="如: person_read" disabled={!!editingPermission} />
            </Form.Item>
          )}
          <Form.Item
            name="name"
            label="权限名称"
            rules={[{ required: true, message: '请输入权限名称' }]}
          >
            <Input placeholder="请输入权限名称" />
          </Form.Item>
          {!editingPermission && (
            <Form.Item
              name="module"
              label="模块"
              rules={[{ required: true, message: '请选择模块' }]}
            >
              <Select placeholder="选择所属模块">
                <Select.Option value="admin">系统管理</Select.Option>
                <Select.Option value="person">人物管理</Select.Option>
                <Select.Option value="genealogy">族谱管理</Select.Option>
                <Select.Option value="culture">文化管理</Select.Option>
                <Select.Option value="memorial">纪念馆管理</Select.Option>
                <Select.Option value="cemetery">墓园管理</Select.Option>
                <Select.Option value="community">社区管理</Select.Option>
              </Select>
            </Form.Item>
          )}
          <Form.Item name="description" label="描述">
            <Input.TextArea placeholder="请输入权限描述" rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}