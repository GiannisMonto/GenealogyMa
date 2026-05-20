import { ReactElement, useEffect, useState } from 'react';
import { Table, Button, Space, Input, Select, Tag, Typography, Card } from 'antd';
import { SearchOutlined, ReloadOutlined, PlusOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { useNavigate } from 'react-router-dom';
import { personApi } from '@shared/api/person';
import type { PersonDTO } from '@shared/api/types';
import styles from './index.module.css';

const { Title } = Typography;

interface SearchParams {
  keyword: string;
  gender: string;
  generation: number | null;
  page: number;
  pageSize: number;
}

const defaultSearchParams: SearchParams = {
  keyword: '',
  gender: '',
  generation: null,
  page: 1,
  pageSize: 20,
};

export function Members(): ReactElement {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [searchParams, setSearchParams] = useState<SearchParams>(defaultSearchParams);
  const [data, setData] = useState<PersonDTO[]>([]);
  const [pagination, setPagination] = useState({ total: 0, page: 1, pageSize: 20 });

  const fetchData = async () => {
    setLoading(true);
    try {
      const result = await personApi.searchPersons({
        keyword: searchParams.keyword || undefined,
        gender: (searchParams.gender as '男' | '女') || undefined,
        generation: searchParams.generation || undefined,
        page: searchParams.page,
        pageSize: searchParams.pageSize,
      });
      setData(result.data);
      setPagination({
        total: result.total,
        page: result.page,
        pageSize: result.pageSize,
      });
    } catch {
      // 使用模拟数据
      setData(mockData);
      setPagination({ total: mockData.length, page: 1, pageSize: 20 });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [searchParams]);

  const handleSearch = (value: string) => {
    setSearchParams((prev) => ({ ...prev, keyword: value, page: 1 }));
  };

  const handleFilterChange = (key: keyof SearchParams, value: string | number | null) => {
    setSearchParams((prev) => ({ ...prev, [key]: value, page: 1 }));
  };

  const handleTableChange = (page: number, pageSize: number) => {
    setSearchParams((prev) => ({ ...prev, page, pageSize }));
  };

  const columns: ColumnsType<PersonDTO> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
      fixed: 'left',
    },
    {
      title: '姓名',
      dataIndex: 'name',
      width: 120,
      fixed: 'left',
      render: (name: string, record) => (
        <a onClick={() => navigate(`/members/${record.id}`)}>{name}</a>
      ),
    },
    {
      title: '字号',
      dataIndex: 'style_name',
      width: 100,
    },
    {
      title: '性别',
      dataIndex: 'gender',
      width: 80,
      render: (gender: string) => (
        <Tag color={gender === '男' ? 'blue' : 'pink'}>{gender || '-'}</Tag>
      ),
    },
    {
      title: '世代',
      dataIndex: 'generation',
      width: 80,
      sorter: (a, b) => (a.generation || 0) - (b.generation || 0),
    },
    {
      title: '父亲',
      dataIndex: 'father_id',
      width: 100,
      render: (fatherId: number | null) => fatherId || '-',
    },
    {
      title: '出生日期',
      dataIndex: 'birth_time_text',
      width: 120,
    },
    {
      title: '去世日期',
      dataIndex: 'death_time_text',
      width: 120,
    },
    {
      title: '出生地',
      dataIndex: 'birth_place',
      width: 150,
      ellipsis: true,
    },
    {
      title: '葬地',
      dataIndex: 'burial_place',
      width: 150,
      ellipsis: true,
    },
    {
      title: '子女数',
      dataIndex: 'total_children_count',
      width: 100,
      render: (count: number) => (count ? `${count} 人` : '-'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right',
      render: (_: unknown, record: PersonDTO) => (
        <Space size="small">
          <Button size="small" type="link" onClick={() => navigate(`/members/${record.id}`)}>
            详情
          </Button>
          <Button size="small" type="link" onClick={() => navigate(`/members/${record.id}/edit`)}>
            编辑
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Title level={3}>成员管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/members/create')}>
          新增成员
        </Button>
      </div>

      <Card className={styles.filterCard}>
        <Space wrap size="middle">
          <Input.Search
            placeholder="搜索姓名、字号..."
            allowClear
            enterButton={<SearchOutlined />}
            onSearch={handleSearch}
            style={{ width: 260 }}
          />
          <Select
            placeholder="性别"
            allowClear
            style={{ width: 120 }}
            value={searchParams.gender || undefined}
            onChange={(value) => handleFilterChange('gender', value || '')}
            options={[
              { label: '男', value: '男' },
              { label: '女', value: '女' },
            ]}
          />
          <Select
            placeholder="世代"
            allowClear
            style={{ width: 120 }}
            value={searchParams.generation}
            onChange={(value) => handleFilterChange('generation', value)}
            options={Array.from({ length: 22 }, (_, i) => ({
              label: `第${i + 1}代`,
              value: i + 1,
            }))}
          />
          <Button icon={<ReloadOutlined />} onClick={() => setSearchParams(defaultSearchParams)}>
            重置
          </Button>
        </Space>
      </Card>

      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        scroll={{ x: 1500 }}
        pagination={{
          current: pagination.page,
          pageSize: pagination.pageSize,
          total: pagination.total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`,
          onChange: handleTableChange,
        }}
      />
    </div>
  );
}

// 模拟数据
const mockData: PersonDTO[] = [
  {
    id: 1,
    legacy_id: 'L001',
    name: '张三',
    style_name: '字云峰',
    gender: '男',
    generation: 1,
    birth_order: '长',
    father_id: null,
    lineage_path: '1',
    detail_text: '张三为始祖',
    birth_time_text: '1900-01-01',
    death_time_text: '1970-01-01',
    birth_place: '浙江绍兴',
    burial_place: '浙江绍兴',
    son_count: 3,
    daughter_count: 2,
    adopted_heir_count: 0,
    total_children_count: 5,
    age: 70,
    is_alive: false,
    full_name: '张三',
    created_at: '2024-01-01',
    updated_at: '2024-01-01',
  },
  {
    id: 2,
    legacy_id: 'L002',
    name: '张四',
    style_name: '字海涛',
    gender: '男',
    generation: 2,
    birth_order: '次',
    father_id: 1,
    lineage_path: '1.2',
    detail_text: '张三次子',
    birth_time_text: '1930-05-15',
    death_time_text: '',
    birth_place: '浙江绍兴',
    burial_place: '',
    son_count: 2,
    daughter_count: 1,
    adopted_heir_count: 0,
    total_children_count: 3,
    age: 94,
    is_alive: true,
    full_name: '张四',
    created_at: '2024-01-01',
    updated_at: '2024-01-01',
  },
  {
    id: 3,
    legacy_id: 'L003',
    name: '张小红',
    style_name: '',
    gender: '女',
    generation: 2,
    birth_order: '三',
    father_id: 1,
    lineage_path: '1.3',
    detail_text: '张三三女',
    birth_time_text: '1935-08-20',
    death_time_text: '',
    birth_place: '浙江绍兴',
    burial_place: '',
    son_count: 0,
    daughter_count: 2,
    adopted_heir_count: 0,
    total_children_count: 2,
    age: 89,
    is_alive: true,
    full_name: '张小红',
    created_at: '2024-01-01',
    updated_at: '2024-01-01',
  },
];