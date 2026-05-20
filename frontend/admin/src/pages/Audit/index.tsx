import { ReactElement, useEffect, useState } from 'react';
import {
  Table,
  Button,
  Space,
  Input,
  Tag,
  Typography,
  Card,
  DatePicker,
  Select,
  Modal,
  message,
  Tooltip,
} from 'antd';
import {
  SearchOutlined,
  EyeOutlined,
  ExportOutlined,
  FilterOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { Dayjs } from 'dayjs';
import dayjs from 'dayjs';
import styles from './index.module.css';
import { getAuditLogs, getAuditLogDetail, exportAuditLogs, type AuditLogDTO, type AuditFilters } from '@/api/audit';

const { Title } = Typography;
const { RangePicker } = DatePicker;

const MODULE_OPTIONS = [
  { label: '用户模块', value: 'user' },
  { label: '角色模块', value: 'role' },
  { label: '权限模块', value: 'permission' },
  { label: '人物模块', value: 'person' },
  { label: '系统模块', value: 'system' },
];

const ACTION_OPTIONS = [
  { label: '创建', value: 'create' },
  { label: '更新', value: 'update' },
  { label: '删除', value: 'delete' },
  { label: '登录', value: 'login' },
  { label: '登出', value: 'logout' },
];

const getActionTag = (action: string) => {
  const colorMap: Record<string, string> = {
    create: 'green',
    update: 'blue',
    delete: 'red',
    login: 'cyan',
    logout: 'default',
  };
  const labelMap: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除',
    login: '登录',
    logout: '登出',
  };
  return <Tag color={colorMap[action] || 'default'}>{labelMap[action] || action}</Tag>;
};

// Mock data for demonstration
const mockAuditLogs: AuditLogDTO[] = [
  {
    id: 1,
    user_id: 1,
    username: 'admin',
    module: 'user',
    action: 'create',
    resource_type: 'user',
    resource_id: '10',
    description: '创建新用户 testuser',
    ip_address: '192.168.1.100',
    created_at: '2026-05-21T10:30:00Z',
  },
  {
    id: 2,
    user_id: 1,
    username: 'admin',
    module: 'role',
    action: 'update',
    resource_type: 'role',
    resource_id: '5',
    old_value: '{"name":"普通用户"}',
    new_value: '{"name":"高级用户"}',
    ip_address: '192.168.1.100',
    created_at: '2026-05-21T10:25:00Z',
  },
  {
    id: 3,
    user_id: 2,
    username: 'editor',
    module: 'person',
    action: 'update',
    resource_type: 'person',
    resource_id: '1024',
    description: '更新人物信息',
    ip_address: '192.168.1.101',
    created_at: '2026-05-21T10:20:00Z',
  },
  {
    id: 4,
    user_id: 3,
    username: 'viewer',
    module: 'system',
    action: 'login',
    resource_type: 'session',
    resource_id: 'sess_abc123',
    ip_address: '192.168.1.102',
    created_at: '2026-05-21T10:15:00Z',
  },
  {
    id: 5,
    user_id: 1,
    username: 'admin',
    module: 'permission',
    action: 'delete',
    resource_type: 'permission',
    resource_id: '25',
    description: '删除权限 read_users',
    ip_address: '192.168.1.100',
    created_at: '2026-05-21T10:10:00Z',
  },
];

export function Audit(): ReactElement {
  const [loading, setLoading] = useState(false);
  const [auditLogs, setAuditLogs] = useState<AuditLogDTO[]>([]);
  const [filters, setFilters] = useState<AuditFilters>({});
  const [detailVisible, setDetailVisible] = useState(false);
  const [selectedLog, setSelectedLog] = useState<AuditLogDTO | null>(null);
  const [messageApi, contextHolder] = message.useMessage();

  useEffect(() => {
    fetchAuditLogs();
  }, [filters]);

  const fetchAuditLogs = async () => {
    setLoading(true);
    try {
      const data = await getAuditLogs(filters);
      setAuditLogs(data.data);
    } catch {
      // Fallback to mock data when backend is not available
      await new Promise((resolve) => setTimeout(resolve, 500));
      let filtered = [...mockAuditLogs];

      if (filters.keyword) {
        const kw = filters.keyword.toLowerCase();
        filtered = filtered.filter(
          (log) =>
            log.username.toLowerCase().includes(kw) ||
            log.description?.toLowerCase().includes(kw) ||
            log.resource_id.includes(kw)
        );
      }

      if (filters.module) {
        filtered = filtered.filter((log) => log.module === filters.module);
      }

      if (filters.action) {
        filtered = filtered.filter((log) => log.action === filters.action);
      }

      if (filters.start_date) {
        filtered = filtered.filter(
          (log) => new Date(log.created_at) >= new Date(filters.start_date!)
        );
      }

      if (filters.end_date) {
        filtered = filtered.filter(
          (log) => new Date(log.created_at) <= new Date(filters.end_date!)
        );
      }

      setAuditLogs(filtered);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    setFilters((prev) => ({ ...prev, keyword: value }));
  };

  const handleFilterChange = (key: keyof AuditFilters, value: string | undefined) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
  };

  const handleDateChange = (
    dates: [Dayjs | null, Dayjs | null] | null
  ) => {
    if (dates && dates[0] && dates[1]) {
      setFilters((prev) => ({
        ...prev,
        start_date: dates[0]!.toISOString(),
        end_date: dates[1]!.toISOString(),
      }));
    } else {
      setFilters((prev) => ({
        ...prev,
        start_date: undefined,
        end_date: undefined,
      }));
    }
  };

  const handleViewDetail = (log: AuditLogDTO) => {
    setSelectedLog(log);
    setDetailVisible(true);
  };

  const handleExport = async () => {
    try {
      const blob = await exportAuditLogs(filters);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `audit-logs-${dayjs().format('YYYY-MM-DD')}.xlsx`;
      a.click();
      window.URL.revokeObjectURL(url);
      messageApi.success('导出成功');
    } catch {
      messageApi.error('导出失败');
    }
  };

  const columns: ColumnsType<AuditLogDTO> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '时间',
      dataIndex: 'created_at',
      width: 180,
      render: (time: string) => new Date(time).toLocaleString('zh-CN'),
    },
    {
      title: '用户',
      dataIndex: 'username',
      width: 120,
    },
    {
      title: '模块',
      dataIndex: 'module',
      width: 100,
      render: (module: string) => {
        const mod = MODULE_OPTIONS.find((m) => m.value === module);
        return mod?.label || module;
      },
    },
    {
      title: '操作',
      dataIndex: 'action',
      width: 100,
      render: (action: string) => getActionTag(action),
    },
    {
      title: '资源类型',
      dataIndex: 'resource_type',
      width: 120,
    },
    {
      title: '资源ID',
      dataIndex: 'resource_id',
      width: 100,
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
    },
    {
      title: 'IP地址',
      dataIndex: 'ip_address',
      width: 140,
    },
    {
      title: '操作',
      key: 'action',
      width: 80,
      render: (_, record) => (
        <Tooltip title="查看详情">
          <Button
            size="small"
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          />
        </Tooltip>
      ),
    },
  ];

  return (
    <div className={styles.container}>
      {contextHolder}
      <div className={styles.header}>
        <Title level={3}>审计日志</Title>
        <Button icon={<ExportOutlined />} onClick={handleExport}>
          导出
        </Button>
      </div>

      <Card className={styles.filterCard}>
        <Space wrap align="end">
          <Input.Search
            placeholder="搜索用户名、描述..."
            allowClear
            enterButton={<SearchOutlined />}
            onSearch={handleSearch}
            style={{ width: 250 }}
          />
          <Select
            placeholder="选择模块"
            allowClear
            style={{ width: 150 }}
            options={MODULE_OPTIONS}
            onChange={(value) => handleFilterChange('module', value)}
          />
          <Select
            placeholder="选择操作"
            allowClear
            style={{ width: 120 }}
            options={ACTION_OPTIONS}
            onChange={(value) => handleFilterChange('action', value)}
          />
          <RangePicker
            showTime
            onChange={handleDateChange}
          />
          <Button
            icon={<FilterOutlined />}
            onClick={() => setFilters({})}
          >
            重置
          </Button>
        </Space>
      </Card>

      <Table
        columns={columns}
        dataSource={auditLogs}
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
        title="审计日志详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailVisible(false)}>
            关闭
          </Button>,
        ]}
      >
        {selectedLog && (
          <div className={styles.detailContent}>
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>日志ID：</span>
              <span>{selectedLog.id}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>用户：</span>
              <span>{selectedLog.username}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>模块：</span>
              <span>
                {MODULE_OPTIONS.find((m) => m.value === selectedLog.module)?.label ||
                  selectedLog.module}
              </span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>操作：</span>
              {getActionTag(selectedLog.action)}
            </div>
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>资源类型：</span>
              <span>{selectedLog.resource_type}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>资源ID：</span>
              <span>{selectedLog.resource_id}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>IP地址：</span>
              <span>{selectedLog.ip_address || '-'}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>时间：</span>
              <span>{new Date(selectedLog.created_at).toLocaleString('zh-CN')}</span>
            </div>
            <div className={styles.detailRow}>
              <span className={styles.detailLabel}>描述：</span>
              <span>{selectedLog.description || '-'}</span>
            </div>
            {selectedLog.old_value && (
              <div className={styles.detailRow}>
                <span className={styles.detailLabel}>旧值：</span>
                <pre className={styles.codeBlock}>{selectedLog.old_value}</pre>
              </div>
            )}
            {selectedLog.new_value && (
              <div className={styles.detailRow}>
                <span className={styles.detailLabel}>新值：</span>
                <pre className={styles.codeBlock}>{selectedLog.new_value}</pre>
              </div>
            )}
          </div>
        )}
      </Modal>
    </div>
  );
}