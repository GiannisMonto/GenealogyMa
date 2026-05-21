import { ReactElement, useEffect, useState } from 'react';
import {
  Table,
  Button,
  Space,
  Input,
  Tag,
  Typography,
  Modal,
  Form,
  Select,
  message,
  Popconfirm,
  Drawer,
  Descriptions,
  Progress,
  Card,
  Row,
  Col,
  Statistic,
} from 'antd';
import {
  PlusOutlined,
  DeleteOutlined,
  SearchOutlined,
  EyeOutlined,
  ReloadOutlined,
  CloudUploadOutlined,
  HistoryOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { backupApi, type BackupDTO, type CreateBackupRequest } from '@/api/backup';
import styles from './index.module.css';

const { Title } = Typography;
const { Search } = Input;

const BACKUP_TYPE_OPTIONS = [
  { label: '全量备份', value: 'full' },
  { label: '增量备份', value: 'incremental' },
  { label: '差异备份', value: 'differential' },
];

const STATUS_COLOR_MAP: Record<string, string> = {
  pending: 'default',
  running: 'processing',
  completed: 'success',
  failed: 'error',
};

const STATUS_TEXT_MAP: Record<string, string> = {
  pending: '等待中',
  running: '进行中',
  completed: '已完成',
  failed: '失败',
};

export function Backup(): ReactElement {
  const [loading, setLoading] = useState(false);
  const [backups, setBackups] = useState<BackupDTO[]>([]);
  const [total, setTotal] = useState(0);
  const [searchText, setSearchText] = useState('');
  const [typeFilter, setTypeFilter] = useState<string | undefined>();
  const [statusFilter, setStatusFilter] = useState<string | undefined>();
  const [modalVisible, setModalVisible] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [selectedBackup, setSelectedBackup] = useState<BackupDTO | null>(null);
  const [cleanModalVisible, setCleanModalVisible] = useState(false);
  const [cleanDays, setCleanDays] = useState(30);
  const [form] = Form.useForm<CreateBackupRequest>();
  const [messageApi, contextHolder] = message.useMessage();

  const fetchBackups = async () => {
    setLoading(true);
    try {
      const result = await backupApi.getBackups({
        type: typeFilter,
        status: statusFilter,
        page_size: 100,
      });
      setBackups(result.data);
      setTotal(result.total);
    } catch {
      messageApi.error('获取备份列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchBackups();
  }, [typeFilter, statusFilter]);

  const filteredBackups = backups.filter(
    (b) =>
      b.name.toLowerCase().includes(searchText.toLowerCase()) ||
      b.database.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleCreate = () => {
    form.resetFields();
    setModalVisible(true);
  };

  const handleView = (backup: BackupDTO) => {
    setSelectedBackup(backup);
    setDetailVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await backupApi.deleteBackup(id);
      messageApi.success('删除成功');
      fetchBackups();
    } catch {
      messageApi.error('删除失败');
    }
  };

  const handleRestore = async (backup: BackupDTO) => {
    try {
      await backupApi.restoreBackup(backup.id);
      messageApi.success('恢复任务已启动');
      setDetailVisible(false);
    } catch {
      messageApi.error('恢复失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      await backupApi.createBackup(values);
      messageApi.success('备份任务已创建');
      setModalVisible(false);
      fetchBackups();
    } catch {
      messageApi.error('创建失败');
    }
  };

  const handleCleanOld = async () => {
    try {
      const result = await backupApi.cleanOldBackups(cleanDays);
      messageApi.success(`已清理 ${result.deleted_count} 个旧备份`);
      setCleanModalVisible(false);
      fetchBackups();
    } catch {
      messageApi.error('清理失败');
    }
  };

  const getStatusPercent = (backup: BackupDTO) => {
    if (backup.status === 'completed') return 100;
    if (backup.status === 'failed') return 100;
    if (backup.status === 'running') return 50;
    return 0;
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CloudUploadOutlined />;
      case 'running':
        return <ReloadOutlined spin />;
      case 'failed':
        return <DeleteOutlined />;
      default:
        return <HistoryOutlined />;
    }
  };

  const columns: ColumnsType<BackupDTO> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '备份名称',
      dataIndex: 'name',
      width: 200,
      render: (name: string, record) => (
        <Button type="link" onClick={() => handleView(record)}>
          {name}
        </Button>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      width: 120,
      render: (type: string) => {
        const option = BACKUP_TYPE_OPTIONS.find((o) => o.value === type);
        return <Tag color="blue">{option?.label || type}</Tag>;
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 120,
      render: (status: string) => (
        <Space>
          {getStatusIcon(status)}
          <Tag color={STATUS_COLOR_MAP[status]}>{STATUS_TEXT_MAP[status] || status}</Tag>
        </Space>
      ),
    },
    {
      title: '数据库',
      dataIndex: 'database',
      width: 120,
    },
    {
      title: '文件大小',
      dataIndex: 'file_size_str',
      width: 100,
      render: (size: string) => size || '-',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: unknown, record: BackupDTO) => (
        <Space size="small">
          <Button size="small" type="link" icon={<EyeOutlined />} onClick={() => handleView(record)}>
            详情
          </Button>
          <Popconfirm
            title="确定删除此备份？"
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

  const completedCount = backups.filter((b) => b.status === 'completed').length;
  const failedCount = backups.filter((b) => b.status === 'failed').length;
  const runningCount = backups.filter((b) => b.status === 'running').length;

  return (
    <div className={styles.container}>
      {contextHolder}
      <div className={styles.header}>
        <Title level={3}>备份管理</Title>
        <Space>
          <Button icon={<HistoryOutlined />} onClick={() => setCleanModalVisible(true)}>
            清理旧备份
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            创建备份
          </Button>
        </Space>
      </div>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic title="总备份数" value={total} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="已完成" value={completedCount} valueStyle={{ color: '#52c41a' }} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="进行中" value={runningCount} valueStyle={{ color: '#1890ff' }} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="失败" value={failedCount} valueStyle={{ color: '#ff4d4f' }} />
          </Card>
        </Col>
      </Row>

      <Card className={styles.filterCard}>
        <Space wrap>
          <Search
            placeholder="搜索备份名称或数据库..."
            allowClear
            enterButton={<SearchOutlined />}
            onSearch={setSearchText}
            style={{ width: 280 }}
          />
          <Select
            placeholder="备份类型"
            allowClear
            style={{ width: 140 }}
            onChange={setTypeFilter}
            options={BACKUP_TYPE_OPTIONS}
          />
          <Select
            placeholder="状态"
            allowClear
            style={{ width: 120 }}
            onChange={setStatusFilter}
            options={[
              { label: '等待中', value: 'pending' },
              { label: '进行中', value: 'running' },
              { label: '已完成', value: 'completed' },
              { label: '失败', value: 'failed' },
            ]}
          />
          <Button icon={<ReloadOutlined />} onClick={fetchBackups}>
            刷新
          </Button>
        </Space>
      </Card>

      <Table
        columns={columns}
        dataSource={filteredBackups}
        rowKey="id"
        loading={loading}
        pagination={{
          pageSize: 10,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (t) => `共 ${t} 条`,
        }}
      />

      <Modal
        title="创建备份"
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item
            name="name"
            label="备份名称"
            rules={[{ required: true, message: '请输入备份名称' }]}
          >
            <Input placeholder="请输入备份名称" />
          </Form.Item>
          <Form.Item
            name="type"
            label="备份类型"
            rules={[{ required: true, message: '请选择备份类型' }]}
          >
            <Select placeholder="请选择备份类型" options={BACKUP_TYPE_OPTIONS} />
          </Form.Item>
          <Form.Item name="database" label="数据库">
            <Input placeholder="留空表示全部数据库" />
          </Form.Item>
        </Form>
      </Modal>

      <Drawer
        title="备份详情"
        open={detailVisible}
        onClose={() => setDetailVisible(false)}
        width={500}
        extra={
          selectedBackup?.status === 'completed' && (
            <Popconfirm
              title="确定恢复此备份？"
              description="当前数据将被备份并覆盖"
              onConfirm={() => selectedBackup && handleRestore(selectedBackup)}
            >
              <Button type="primary">恢复</Button>
            </Popconfirm>
          )
        }
      >
        {selectedBackup && (
          <>
            <Descriptions column={1} bordered>
              <Descriptions.Item label="备份名称">{selectedBackup.name}</Descriptions.Item>
              <Descriptions.Item label="备份类型">
                {BACKUP_TYPE_OPTIONS.find((o) => o.value === selectedBackup.type)?.label ||
                  selectedBackup.type}
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                <Space>
                  {getStatusIcon(selectedBackup.status)}
                  <Tag color={STATUS_COLOR_MAP[selectedBackup.status]}>
                    {STATUS_TEXT_MAP[selectedBackup.status] || selectedBackup.status}
                  </Tag>
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="数据库">{selectedBackup.database || '-'}</Descriptions.Item>
              <Descriptions.Item label="文件路径">{selectedBackup.file_path || '-'}</Descriptions.Item>
              <Descriptions.Item label="文件大小">
                {selectedBackup.file_size_str || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">{selectedBackup.created_at}</Descriptions.Item>
              {selectedBackup.started_at && (
                <Descriptions.Item label="开始时间">{selectedBackup.started_at}</Descriptions.Item>
              )}
              {selectedBackup.completed_at && (
                <Descriptions.Item label="完成时间">{selectedBackup.completed_at}</Descriptions.Item>
              )}
              {selectedBackup.error_message && (
                <Descriptions.Item label="错误信息">
                  <span style={{ color: '#ff4d4f' }}>{selectedBackup.error_message}</span>
                </Descriptions.Item>
              )}
            </Descriptions>

            {selectedBackup.status !== 'completed' && selectedBackup.status !== 'failed' && (
              <div className={styles.backupProgress}>
                <p>备份进度</p>
                <Progress percent={getStatusPercent(selectedBackup)} status="active" />
              </div>
            )}
          </>
        )}
      </Drawer>

      <Modal
        title="清理旧备份"
        open={cleanModalVisible}
        onOk={handleCleanOld}
        onCancel={() => setCleanModalVisible(false)}
      >
        <p>确定清理 {cleanDays} 天之前的备份吗？</p>
        <Input
          type="number"
          value={cleanDays}
          onChange={(e) => setCleanDays(Number(e.target.value))}
          placeholder="天数"
          style={{ marginTop: 16 }}
        />
      </Modal>
    </div>
  );
}