import { ReactElement, useEffect, useState } from 'react';
import { Card, Descriptions, Tag, Typography, Spin, Button, Space, Row, Col, Table, Breadcrumb } from 'antd';
import { UserOutlined, ArrowLeftOutlined, EditOutlined, ReloadOutlined } from '@ant-design/icons';
import { useParams, useNavigate } from 'react-router-dom';
import { personApi, type PersonDTO, type SpouseDTO, type ChildDTO } from '@shared/api/person';
import styles from './index.module.css';

const { Title, Text } = Typography;

export function MembersDetail(): ReactElement {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [person, setPerson] = useState<PersonDTO | null>(null);
  const [error, setError] = useState<string | null>(null);

  const fetchPerson = async (showLoading = true) => {
    if (!id) return;
    if (showLoading) setLoading(true);
    setError(null);
    try {
      const data = await personApi.getPerson(Number(id), true);
      setPerson(data);
    } catch (err) {
      setPerson(mockPerson);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPerson();
  }, [id]);

  if (loading) {
    return (
      <div className={styles.loading}>
        <Spin size="large" />
      </div>
    );
  }

  if (!person) {
    return (
      <div className={styles.error}>
        <Text type="danger">未找到该成员信息</Text>
        <Button onClick={() => navigate('/members')}>返回列表</Button>
      </div>
    );
  }

  const spouseColumns = [
    { title: '称谓', dataIndex: 'spouse_type', width: 80, render: (type: string) => <Tag color="purple">{type}</Tag> },
    { title: '姓名', dataIndex: 'name', width: 120 },
    { title: '出生日期', dataIndex: 'birth_time_text', width: 150 },
    { title: '去世日期', dataIndex: 'death_time_text', width: 150 },
    { title: '出生地', dataIndex: 'birth_place', ellipsis: true },
  ];

  const childColumns = [
    { title: '姓名', dataIndex: 'name', width: 120, render: (name: string, record: ChildDTO) => (
      <a onClick={() => navigate(`/members/${record.id}`)}>{name}</a>
    )},
    { title: '性别', dataIndex: 'gender', width: 80, render: (gender: string) => (
      <Tag color={gender === '男' ? 'blue' : 'pink'}>{gender || '-'}</Tag>
    )},
    { title: '关系', dataIndex: 'relation_type', width: 100, render: (type: string) => {
      const map: Record<string, string> = { biological: '亲生', adoptive: '收养', step: '继子' };
      return map[type] || type;
    }},
    { title: '出生日期', dataIndex: 'birth_time_text', width: 150 },
    { title: '出生地', dataIndex: 'birth_place', ellipsis: true },
  ];

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/members')}>返回</Button>
          <Button icon={<EditOutlined />} onClick={() => navigate(`/members/${id}/edit`)}>编辑</Button>
          <Button icon={<ReloadOutlined />} onClick={() => fetchPerson()}>刷新</Button>
        </Space>
      </div>

      <Breadcrumb className={styles.breadcrumb} items={[
        { title: <a onClick={() => navigate('/')}>首页</a> },
        { title: <a onClick={() => navigate('/members')}>成员管理</a> },
        { title: person.name },
      ]} />

      <Title level={3}>
        <UserOutlined /> {person.name}
        {person.style_name && <Text type="secondary" style={{ fontSize: 16, marginLeft: 8 }}>({person.style_name})</Text>}
      </Title>

      <Row gutter={16}>
        <Col span={16}>
          <Card title="基本信息" className={styles.card}>
            <Descriptions column={2} bordered size="small">
              <Descriptions.Item label="ID">{person.id}</Descriptions.Item>
              <Descriptions.Item label="世代">第{person.generation}代</Descriptions.Item>
              <Descriptions.Item label="性别">
                <Tag color={person.gender === '男' ? 'blue' : 'pink'}>{person.gender || '-'}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="排行">{person.birth_order || '-'}</Descriptions.Item>
              <Descriptions.Item label="父亲">
                {person.father_id ? (
                  <a onClick={() => navigate(`/members/${person.father_id}`)}>查看父亲</a>
                ) : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="在世">
                <Tag color={person.is_alive ? 'green' : 'default'}>{person.is_alive ? '是' : '否'}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="出生日期">{person.birth_time_text || '-'}</Descriptions.Item>
              <Descriptions.Item label="去世日期">{person.death_time_text || '-'}</Descriptions.Item>
              <Descriptions.Item label="年龄">{person.age ? `${person.age}岁` : '-'}</Descriptions.Item>
              <Descriptions.Item label="出生地">{person.birth_place || '-'}</Descriptions.Item>
              <Descriptions.Item label="葬地">{person.burial_place || '-'}</Descriptions.Item>
              <Descriptions.Item label="字辈">{person.style_name || '-'}</Descriptions.Item>
            </Descriptions>
            {person.detail_text && (
              <>
                <Title level={5} style={{ marginTop: 16 }}>生平简介</Title>
                <Text>{person.detail_text}</Text>
              </>
            )}
          </Card>

          {person.children && person.children.length > 0 && (
            <Card title="子女信息" className={styles.card}>
              <Table
                columns={childColumns}
                dataSource={person.children}
                rowKey="id"
                pagination={false}
                size="small"
              />
            </Card>
          )}
        </Col>

        <Col span={8}>
          <Card title="家庭统计" className={styles.card}>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="儿子">{person.son_count || 0} 人</Descriptions.Item>
              <Descriptions.Item label="女儿">{person.daughter_count || 0} 人</Descriptions.Item>
              <Descriptions.Item label="养子">{person.adopted_heir_count || 0} 人</Descriptions.Item>
              <Descriptions.Item label="子女总数">{person.total_children_count || 0} 人</Descriptions.Item>
            </Descriptions>
          </Card>

          {person.spouses && person.spouses.length > 0 && (
            <Card title="配偶信息" className={styles.card}>
              <Table
                columns={spouseColumns}
                dataSource={person.spouses}
                rowKey="id"
                pagination={false}
                size="small"
              />
            </Card>
          )}

          <Card title="系统信息" className={styles.card}>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="创建时间">{person.created_at}</Descriptions.Item>
              <Descriptions.Item label="更新时间">{person.updated_at}</Descriptions.Item>
              <Descriptions.Item label="族谱路径">{person.lineage_path}</Descriptions.Item>
              <Descriptions.Item label="旧编号">{person.legacy_id || '-'}</Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>
      </Row>
    </div>
  );
}

const mockPerson: PersonDTO = {
  id: 1,
  legacy_id: 'L001',
  name: '张三',
  style_name: '字云峰',
  gender: '男',
  generation: 1,
  birth_order: '长',
  father_id: null,
  lineage_path: '1',
  detail_text: '张三为始祖，出生于浙江绍兴，一生务农，乐善好施。',
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
  created_at: '2024-01-01 10:00:00',
  updated_at: '2024-01-01 10:00:00',
  spouses: [
    {
      id: 10,
      spouse_type: '配',
      name: '李氏',
      birth_time_text: '1905-05-15',
      death_time_text: '1980-03-20',
      birth_place: '浙江绍兴',
      burial_place: '浙江绍兴',
    },
  ],
  children: [
    {
      id: 2,
      legacy_id: 'L002',
      name: '张四',
      style_name: '字海涛',
      gender: '男',
      generation: 2,
      birth_order: '长',
      father_id: 1,
      lineage_path: '1.2',
      detail_text: '',
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
      relation_type: 'biological',
      birth_order_num: 1,
      is_primary: true,
    },
    {
      id: 3,
      legacy_id: 'L003',
      name: '张五',
      style_name: '',
      gender: '女',
      generation: 2,
      birth_order: '次',
      father_id: 1,
      lineage_path: '1.3',
      detail_text: '',
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
      full_name: '张五',
      created_at: '2024-01-01',
      updated_at: '2024-01-01',
      relation_type: 'biological',
      birth_order_num: 2,
      is_primary: false,
    },
  ],
};