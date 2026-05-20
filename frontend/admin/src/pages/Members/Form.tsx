import { ReactElement, useState, useEffect } from 'react';
import { Form, Input, Select, Button, Card, Typography, message, Space, Row, Col, Radio } from 'antd';
import { ArrowLeftOutlined, SaveOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import { personApi, type PersonDTO, type CreatePersonRequest, type UpdatePersonRequest } from '@shared/api/person';
import { useAuthStore } from '@/store/auth';
import styles from './index.module.css';

const { Title, Text } = Typography;

interface PersonFormData {
  name: string;
  style_name: string;
  gender: '男' | '女';
  generation: number;
  birth_order: string;
  father_id: number | null;
  birth_time_text: string;
  death_time_text: string;
  birth_place: string;
  burial_place: string;
  detail_text: string;
}

const initialFormData: PersonFormData = {
  name: '',
  style_name: '',
  gender: '男',
  generation: 1,
  birth_order: '',
  father_id: null,
  birth_time_text: '',
  death_time_text: '',
  birth_place: '',
  burial_place: '',
  detail_text: '',
};

export function PersonForm(): ReactElement {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { token } = useAuthStore();
  const [messageApi, contextHolder] = message.useMessage();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [fetchLoading, setFetchLoading] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const [formData, setFormData] = useState<PersonFormData>(initialFormData);
  const [fatherSearchValue, setFatherSearchValue] = useState('');
  const [fatherOptions, setFatherOptions] = useState<PersonDTO[]>([]);
  const [selectedFather, setSelectedFather] = useState<PersonDTO | null>(null);

  useEffect(() => {
    if (id && id !== 'new') {
      setIsEdit(true);
      fetchPersonData(Number(id));
    }
  }, [id]);

  const fetchPersonData = async (personId: number) => {
    setFetchLoading(true);
    try {
      const data = await personApi.getPerson(personId, true);
      const formValues: PersonFormData = {
        name: data.name || '',
        style_name: data.style_name || '',
        gender: data.gender || '男',
        generation: data.generation || 1,
        birth_order: data.birth_order || '',
        father_id: data.father_id || null,
        birth_time_text: data.birth_time_text || '',
        death_time_text: data.death_time_text || '',
        birth_place: data.birth_place || '',
        burial_place: data.burial_place || '',
        detail_text: data.detail_text || '',
      };
      form.setFieldsValue(formValues);
      setFormData(formValues);
      if (data.father_id) {
        setSelectedFather(data.father);
      }
    } catch {
      messageApi.error('获取人物信息失败');
    } finally {
      setFetchLoading(false);
    }
  };

  const searchFathers = async (keyword: string) => {
    if (!keyword || keyword.length < 1) {
      setFatherOptions([]);
      return;
    }
    try {
      const result = await personApi.searchPersons({ keyword, page: 1, page_size: 20 });
      setFatherOptions(result.data.filter(p => p.id !== (id && id !== 'new' ? Number(id) : 0)));
    } catch {
      setFatherOptions([]);
    }
  };

  useEffect(() => {
    const debounce = setTimeout(() => {
      if (fatherSearchValue) {
        searchFathers(fatherSearchValue);
      } else {
        setFatherOptions([]);
      }
    }, 300);
    return () => clearTimeout(debounce);
  }, [fatherSearchValue]);

  const onFinish = async (values: PersonFormData) => {
    if (!token) {
      messageApi.error('未登录或登录已过期');
      return;
    }
    setLoading(true);
    try {
      if (isEdit && id) {
        const request: UpdatePersonRequest = {
          name: values.name,
          style_name: values.style_name,
          gender: values.gender,
          generation: values.generation,
          birth_order: values.birth_order,
          father_id: values.father_id ?? undefined,
          birth_time_text: values.birth_time_text,
          death_time_text: values.death_time_text,
          birth_place: values.birth_place,
          burial_place: values.burial_place,
          detail_text: values.detail_text,
        };
        await personApi.updatePerson(Number(id), request, token);
        messageApi.success('更新成功');
      } else {
        const request: CreatePersonRequest = {
          name: values.name,
          style_name: values.style_name,
          gender: values.gender,
          generation: values.generation,
          birth_order: values.birth_order,
          father_id: values.father_id ?? undefined,
          birth_time_text: values.birth_time_text,
          death_time_text: values.death_time_text,
          birth_place: values.birth_place,
          burial_place: values.burial_place,
          detail_text: values.detail_text,
        };
        await personApi.createPerson(request, token);
        messageApi.success('创建成功');
      }
      navigate('/members');
    } catch {
      messageApi.error(isEdit ? '更新失败' : '创建失败');
    } finally {
      setLoading(false);
    }
  };

  const handleFatherSelect = (personId: number) => {
    const person = fatherOptions.find(p => p.id === personId);
    if (person) {
      setSelectedFather(person);
      form.setFieldsValue({ father_id: personId });
    }
  };

  return (
    <div className={styles.container}>
      {contextHolder}
      <div className={styles.header}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(isEdit ? `/members/${id}` : '/members')}>
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>{isEdit ? '编辑人物' : '新增人物'}</Title>
        </Space>
      </div>

      <Card title="基本信息" className={styles.card}>
        <Form
          form={form}
          layout="vertical"
          onFinish={onFinish}
          initialValues={formData}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="name"
                label="姓名"
                rules={[{ required: true, message: '请输入姓名' }]}
              >
                <Input placeholder="请输入姓名" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="style_name" label="字辈">
                <Input placeholder="请输入字辈" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="gender"
                label="性别"
                rules={[{ required: true, message: '请选择性别' }]}
              >
                <Radio.Group>
                  <Radio value="男">男</Radio>
                  <Radio value="女">女</Radio>
                </Radio.Group>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="generation" label="世代">
                <Input type="number" min={1} placeholder="请输入世代" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="birth_order" label="排行">
                <Select allowClear placeholder="请选择排行">
                  {['长', '次', '三', '四', '五', '六', '七', '八', '幼'].map(o => (
                    <Select.Option key={o} value={o}>{o}</Select.Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="father_id" label="父亲">
                <Select
                  showSearch
                  allowClear
                  placeholder="搜索选择父亲"
                  value={selectedFather ? selectedFather.id : null}
                  onSearch={(value) => setFatherSearchValue(value)}
                  onClear={() => {
                    setSelectedFather(null);
                    setFatherOptions([]);
                    setFatherSearchValue('');
                  }}
                  onSelect={handleFatherSelect}
                  filterOption={false}
                  notFoundContent={null}
                >
                  {fatherOptions.map(person => (
                    <Select.Option key={person.id} value={person.id}>
                      {person.name} (第{person.generation}代 {person.birth_order || ''})
                    </Select.Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="birth_time_text" label="出生日期">
                <Input placeholder="如：1900-01-01" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="death_time_text" label="去世日期">
                <Input placeholder="如：1970-01-01" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="birth_place" label="出生地">
                <Input placeholder="请输入出生地" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="burial_place" label="葬地">
                <Input placeholder="请输入葬地" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item name="detail_text" label="生平简介">
            <Input.TextArea rows={4} placeholder="请输入生平简介" />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" icon={<SaveOutlined />} loading={loading}>
                {isEdit ? '保存' : '创建'}
              </Button>
              <Button onClick={() => navigate(isEdit ? `/members/${id}` : '/members')}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}