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
  Descriptions,
  Progress,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
  EyeOutlined,
  MapPinOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { CemeteryDTO, GraveDTO } from '@shared/api/types';
import { cemeteryApi } from '@/api/cemetery';
import styles from './index.module.css';

const { Title } = Typography;

interface CemeteryFormData {
  name: string;
  description?: string;
  province: string;
  city: string;
  district?: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  total_grave?: number;
  image_url?: string;
}

export function Cemetery(): ReactElement {
  const [loading, setLoading] = useState(false);
  const [cemeteries, setCemeteries] = useState<CemeteryDTO[]>([]);
  const [searchText, setSearchText] = useState('');
  const [modalVisible, setModalVisible] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [editingCemetery, setEditingCemetery] = useState<CemeteryDTO | null>(null);
  const [selectedCemetery, setSelectedCemetery] = useState<CemeteryDTO | null>(null);
  const [form] = Form.useForm<CemeteryFormData>();
  const [messageApi, contextHolder] = message.useMessage();

  const fetchCemeteries = async () => {
    setLoading(true);
    try {
      const data = await cemeteryApi.getCemeteries();
      setCemeteries(data);
    } catch {
      messageApi.error('获取墓园列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCemeteries();
  }, []);

  const filteredCemeteries = cemeteries.filter(
    (c) =>
      c.name.toLowerCase().includes(searchText.toLowerCase()) ||
      c.province.toLowerCase().includes(searchText.toLowerCase()) ||
      c.city.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleCreate = () => {
    setEditingCemetery(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (cemetery: CemeteryDTO) => {
    setEditingCemetery(cemetery);
    form.setFieldsValue({
      name: cemetery.name,
      description: cemetery.description,
      province: cemetery.province,
      city: cemetery.city,
      district: cemetery.district,
      address: cemetery.address,
      latitude: cemetery.latitude,
      longitude: cemetery.longitude,
      total_grave: cemetery.total_grave,
      image_url: cemetery.image_url,
    });
    setModalVisible(true);
  };

  const handleView = (cemetery: CemeteryDTO) => {
    setSelectedCemetery(cemetery);
    setDetailVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await cemeteryApi.deleteCemetery(id);
      messageApi.success('删除成功');
      fetchCemeteries();
    } catch {
      messageApi.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingCemetery) {
        await cemeteryApi.updateCemetery(editingCemetery.id, values);
        messageApi.success('更新成功');
      } else {
        await cemeteryApi.createCemetery(values);
        messageApi.success('创建成功');
      }
      setModalVisible(false);
      fetchCemeteries();
    } catch {
      messageApi.error('操作失败');
    }
  };

  const getUsagePercent = (cemetery: CemeteryDTO) => {
    if (cemetery.total_grave === 0) return 0;
    return Math.round((cemetery.used_grave / cemetery.total_grave) * 100);
  };

  const columns: ColumnsType<CemeteryDTO> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '墓园名称',
      dataIndex: 'name',
      width: 200,
      render: (name: string, record) => (
        <Button type="link" onClick={() => handleView(record)}>
          {name}
        </Button>
      ),
    },
    {
      title: '位置',
      key: 'location',
      width: 200,
      render: (_, record) => (
        <span>
          {record.province} {record.city} {record.district}
        </span>
      ),
    },
    {
      title: '墓位使用',
      key: 'usage',
      width: 200,
      render: (_, record) => {
        const percent = getUsagePercent(record);
        const status = percent >= 90 ? 'exception' : percent >= 70 ? 'normal' : 'success';
        return (
          <Space direction="vertical" size="small" style={{ width: '100%' }}>
            <Progress percent={percent} size="small" status={status} />
            <span style={{ fontSize: 12, color: '#999' }}>
              {record.used_grave} / {record.total_grave}
            </span>
          </Space>
        );
      },
    },
    {
      title: '地址',
      dataIndex: 'address',
      ellipsis: true,
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: unknown, record: CemeteryDTO) => (
        <Space size="small">
          <Button size="small" type="link" icon={<EyeOutlined />} onClick={() => handleView(record)}>
            详情
          </Button>
          <Button size="small" type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确定删除此墓园？"
            description="删除后无法恢复"
            onConfirm={() => handleDelete(record.id)}
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
        <Title level={3}>墓园管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增墓园
        </Button>
      </div>

      <Card className={styles.filterCard}>
        <Input.Search
          placeholder="搜索墓园名称或地区..."
          allowClear
          enterButton={<SearchOutlined />}
          onSearch={setSearchText}
          style={{ width: 300 }}
        />
      </Card>

      <Table
        columns={columns}
        dataSource={filteredCemeteries}
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
        title={editingCemetery ? '编辑墓园' : '新增墓园'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        destroyOnClose
        width={600}
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item
            name="name"
            label="墓园名称"
            rules={[{ required: true, message: '请输入墓园名称' }]}
          >
            <Input placeholder="请输入墓园名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea placeholder="墓园描述..." rows={3} />
          </Form.Item>
          <Form.Item
            name="province"
            label="省份"
            rules={[{ required: true, message: '请输入省份' }]}
          >
            <Input placeholder="如：广东省" />
          </Form.Item>
          <Form.Item
            name="city"
            label="城市"
            rules={[{ required: true, message: '请输入城市' }]}
          >
            <Input placeholder="如：深圳市" />
          </Form.Item>
          <Form.Item name="district" label="区县">
            <Input placeholder="如：南山区" />
          </Form.Item>
          <Form.Item name="address" label="详细地址">
            <Input placeholder="详细地址" />
          </Form.Item>
          <Form.Item name="total_grave" label="墓位总数">
            <Input type="number" placeholder="墓位总数" />
          </Form.Item>
        </Form>
      </Modal>

      <Drawer
        title="墓园详情"
        open={detailVisible}
        onClose={() => setDetailVisible(false)}
        width={500}
      >
        {selectedCemetery && (
          <Descriptions column={1} bordered>
            <Descriptions.Item label="墓园名称">{selectedCemetery.name}</Descriptions.Item>
            <Descriptions.Item label="描述">{selectedCemetery.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="省份">{selectedCemetery.province}</Descriptions.Item>
            <Descriptions.Item label="城市">{selectedCemetery.city}</Descriptions.Item>
            <Descriptions.Item label="区县">{selectedCemetery.district || '-'}</Descriptions.Item>
            <Descriptions.Item label="详细地址">{selectedCemetery.address || '-'}</Descriptions.Item>
            <Descriptions.Item label="经纬度">
              {selectedCemetery.latitude}, {selectedCemetery.longitude}
            </Descriptions.Item>
            <Descriptions.Item label="墓位使用情况">
              <Space direction="vertical">
                <Progress
                  percent={getUsagePercent(selectedCemetery)}
                  status={
                    getUsagePercent(selectedCemetery) >= 90
                      ? 'exception'
                      : 'normal'
                  }
                />
                <span>
                  已使用 {selectedCemetery.used_grave} / 总数 {selectedCemetery.total_grave}
                </span>
              </Space>
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">{selectedCemetery.created_at}</Descriptions.Item>
            <Descriptions.Item label="更新时间">{selectedCemetery.updated_at}</Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>
    </div>
  );
}