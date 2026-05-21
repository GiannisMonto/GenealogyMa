import { ReactElement, useEffect, useState } from 'react';
import {
  Table,
  Button,
  Space,
  Input,
  Typography,
  Card,
  Modal,
  Form,
  message,
  Popconfirm,
  Drawer,
  Descriptions,
  Progress,
  Tag,
  Select,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
  EyeOutlined,
  HomeOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { MemorialHallDTO } from '@shared/api/types';
import { memorialApi } from '@/api/memorial';
import styles from './index.module.css';

const { Title } = Typography;

const TABLET_TYPE_MAP: Record<string, { color: string; label: string }> = {
  ancestor: { color: 'gold', label: '祖先牌位' },
  martyr: { color: 'red', label: '烈士牌位' },
  sage: { color: 'purple', label: '圣贤牌位' },
  founder: { color: 'blue', label: '创始人牌位' },
};

interface HallFormData {
  name: string;
  description?: string;
  province: string;
  city: string;
  district?: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  build_year?: number;
  style?: string;
  total_tablet?: number;
}

export function Memorial(): ReactElement {
  const [loading, setLoading] = useState(false);
  const [halls, setHalls] = useState<MemorialHallDTO[]>([]);
  const [searchText, setSearchText] = useState('');
  const [modalVisible, setModalVisible] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [editingHall, setEditingHall] = useState<MemorialHallDTO | null>(null);
  const [selectedHall, setSelectedHall] = useState<MemorialHallDTO | null>(null);
  const [form] = Form.useForm<HallFormData>();
  const [messageApi, contextHolder] = message.useMessage();

  const fetchHalls = async () => {
    setLoading(true);
    try {
      const data = await memorialApi.getHalls();
      setHalls(data);
    } catch {
      messageApi.error('获取宗祠列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchHalls();
  }, []);

  const filteredHalls = halls.filter(
    (h) =>
      h.name.toLowerCase().includes(searchText.toLowerCase()) ||
      h.province.includes(searchText) ||
      h.city.includes(searchText)
  );

  const handleCreate = () => {
    setEditingHall(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (hall: MemorialHallDTO) => {
    setEditingHall(hall);
    form.setFieldsValue({
      name: hall.name,
      description: hall.description,
      province: hall.province,
      city: hall.city,
      district: hall.district,
      address: hall.address,
      latitude: hall.latitude,
      longitude: hall.longitude,
      build_year: hall.build_year,
      style: hall.style,
      total_tablet: hall.total_tablet,
    });
    setModalVisible(true);
  };

  const handleView = async (hall: MemorialHallDTO) => {
    try {
      const fullHall = await memorialApi.getHallWithTablets(hall.id);
      setSelectedHall(fullHall);
      setDetailVisible(true);
    } catch {
      messageApi.error('获取宗祠详情失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await memorialApi.deleteHall(id);
      messageApi.success('删除成功');
      fetchHalls();
    } catch {
      messageApi.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingHall) {
        await memorialApi.updateHall(editingHall.id, values);
        messageApi.success('更新成功');
      } else {
        await memorialApi.createHall(values);
        messageApi.success('创建成功');
      }
      setModalVisible(false);
      fetchHalls();
    } catch {
      messageApi.error('操作失败');
    }
  };

  const getUsagePercent = (hall: MemorialHallDTO) => {
    if (hall.total_tablet === 0) return 0;
    return Math.round((hall.used_tablet / hall.total_tablet) * 100);
  };

  const columns: ColumnsType<MemorialHallDTO> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '宗祠名称',
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
      title: '牌位使用',
      key: 'usage',
      width: 200,
      render: (_, record) => {
        const percent = getUsagePercent(record);
        const status = percent >= 90 ? 'exception' : percent >= 70 ? 'normal' : 'success';
        return (
          <Space direction="vertical" size="small" style={{ width: '100%' }}>
            <Progress percent={percent} size="small" status={status} />
            <span style={{ fontSize: 12, color: '#999' }}>
              {record.used_tablet} / {record.total_tablet}
            </span>
          </Space>
        );
      },
    },
    {
      title: '建筑风格',
      dataIndex: 'style',
      width: 120,
      render: (style: string) => style || '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: unknown, record: MemorialHallDTO) => (
        <Space size="small">
          <Button size="small" type="link" icon={<EyeOutlined />} onClick={() => handleView(record)}>
            详情
          </Button>
          <Button size="small" type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确定删除此宗祠？"
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
        <Title level={3}>宗祠管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增宗祠
        </Button>
      </div>

      <Card className={styles.filterCard}>
        <Input.Search
          placeholder="搜索宗祠名称或地区..."
          allowClear
          enterButton={<SearchOutlined />}
          onSearch={setSearchText}
          style={{ width: 300 }}
        />
      </Card>

      <Table
        columns={columns}
        dataSource={filteredHalls}
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
        title={editingHall ? '编辑宗祠' : '新增宗祠'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        destroyOnClose
        width={600}
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item
            name="name"
            label="宗祠名称"
            rules={[{ required: true, message: '请输入宗祠名称' }]}
          >
            <Input placeholder="请输入宗祠名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea placeholder="宗祠描述..." rows={3} />
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
          <Form.Item name="style" label="建筑风格">
            <Input placeholder="如：岭南风格" />
          </Form.Item>
          <Form.Item name="build_year" label="建造年份">
            <Input type="number" placeholder="建造年份" />
          </Form.Item>
          <Form.Item name="total_tablet" label="牌位总数">
            <Input type="number" placeholder="牌位总数" />
          </Form.Item>
        </Form>
      </Modal>

      <Drawer
        title="宗祠详情"
        open={detailVisible}
        onClose={() => setDetailVisible(false)}
        width={600}
      >
        {selectedHall && (
          <>
            <Descriptions column={1} bordered>
              <Descriptions.Item label="宗祠名称">{selectedHall.name}</Descriptions.Item>
              <Descriptions.Item label="描述">{selectedHall.description || '-'}</Descriptions.Item>
              <Descriptions.Item label="省份">{selectedHall.province}</Descriptions.Item>
              <Descriptions.Item label="城市">{selectedHall.city}</Descriptions.Item>
              <Descriptions.Item label="区县">{selectedHall.district || '-'}</Descriptions.Item>
              <Descriptions.Item label="详细地址">{selectedHall.address || '-'}</Descriptions.Item>
              <Descriptions.Item label="建筑风格">{selectedHall.style || '-'}</Descriptions.Item>
              <Descriptions.Item label="建造年份">{selectedHall.build_year || '-'}</Descriptions.Item>
              <Descriptions.Item label="经纬度">
                {selectedHall.latitude}, {selectedHall.longitude}
              </Descriptions.Item>
              <Descriptions.Item label="牌位使用情况">
                <Space direction="vertical">
                  <Progress
                    percent={getUsagePercent(selectedHall)}
                    status={
                      getUsagePercent(selectedHall) >= 90
                        ? 'exception'
                        : 'normal'
                    }
                  />
                  <span>
                    已使用 {selectedHall.used_tablet} / 总数 {selectedHall.total_tablet}
                  </span>
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">{selectedHall.created_at}</Descriptions.Item>
              <Descriptions.Item label="更新时间">{selectedHall.updated_at}</Descriptions.Item>
            </Descriptions>

            {selectedHall.tablets && selectedHall.tablets.length > 0 && (
              <>
                <Title level={5} style={{ marginTop: 24 }}>牌位列表</Title>
                <Table
                  dataSource={selectedHall.tablets}
                  rowKey="id"
                  size="small"
                  pagination={{ pageSize: 5 }}
                  columns={[
                    { title: 'ID', dataIndex: 'id', width: 60 },
                    { title: '姓名', dataIndex: 'person_name', width: 100 },
                    {
                      title: '类型',
                      dataIndex: 'tablet_type',
                      width: 100,
                      render: (type: string) => {
                        const info = TABLET_TYPE_MAP[type] || { color: 'default', label: type };
                        return <Tag color={info.color}>{info.label}</Tag>;
                      },
                    },
                    { title: '位置', dataIndex: 'position', width: 150 },
                    { title: '代数', dataIndex: 'generation', width: 80 },
                  ]}
                />
              </>
            )}
          </>
        )}
      </Drawer>
    </div>
  );
}