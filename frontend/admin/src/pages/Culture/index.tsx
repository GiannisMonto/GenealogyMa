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
  InputNumber,
  message,
  Popconfirm,
  Drawer,
  Descriptions,
  Tabs,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SearchOutlined,
  EyeOutlined,
  BookOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { DocumentDTO, StoryDTO, FamilyTeachingsDTO } from '@shared/api/types';
import { cultureApi } from '@/api/culture';
import styles from './index.module.css';

const { Title } = Typography;

const categoryMap: Record<string, { label: string; color: string }> = {
  classic: { label: '经典', color: 'gold' },
  genealogy: { label: '族谱', color: 'blue' },
  memorial: { label: '纪念', color: 'purple' },
  history: { label: '历史', color: 'green' },
};

export function Culture(): ReactElement {
  const [activeTab, setActiveTab] = useState('documents');
  const [loading, setLoading] = useState(false);
  const [documents, setDocuments] = useState<DocumentDTO[]>([]);
  const [stories, setStories] = useState<StoryDTO[]>([]);
  const [teachings, setTeachings] = useState<FamilyTeachingsDTO[]>([]);
  const [searchText, setSearchText] = useState('');
  const [modalVisible, setModalVisible] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [editingItem, setEditingItem] = useState<DocumentDTO | StoryDTO | FamilyTeachingsDTO | null>(null);
  const [selectedItem, setSelectedItem] = useState<DocumentDTO | StoryDTO | FamilyTeachingsDTO | null>(null);
  const [form] = Form.useForm();
  const [messageApi, contextHolder] = message.useMessage();

  const fetchDocuments = async () => {
    setLoading(true);
    try {
      const data = await cultureApi.getDocuments();
      setDocuments(data);
    } catch {
      messageApi.error('获取文献列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchStories = async () => {
    setLoading(true);
    try {
      const data = await cultureApi.getStories();
      setStories(data);
    } catch {
      messageApi.error('获取故事列表失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchTeachings = async () => {
    setLoading(true);
    try {
      const data = await cultureApi.getFamilyTeachings();
      setTeachings(data);
    } catch {
      messageApi.error('获取家训列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (activeTab === 'documents') fetchDocuments();
    else if (activeTab === 'stories') fetchStories();
    else fetchTeachings();
  }, [activeTab]);

  const filteredDocuments = documents.filter(
    (d) =>
      d.title.toLowerCase().includes(searchText.toLowerCase()) ||
      d.author.toLowerCase().includes(searchText.toLowerCase()) ||
      d.category.toLowerCase().includes(searchText.toLowerCase())
  );

  const filteredStories = stories.filter(
    (s) =>
      s.title.toLowerCase().includes(searchText.toLowerCase()) ||
      s.era.toLowerCase().includes(searchText.toLowerCase())
  );

  const filteredTeachings = teachings.filter(
    (t) => t.title.toLowerCase().includes(searchText.toLowerCase())
  );

  const handleCreate = () => {
    setEditingItem(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (item: DocumentDTO | StoryDTO | FamilyTeachingsDTO) => {
    setEditingItem(item);
    if (activeTab === 'documents') {
      const doc = item as DocumentDTO;
      form.setFieldsValue({
        title: doc.title,
        content: doc.content,
        category: doc.category,
        author: doc.author,
        dynasty: doc.dynasty,
        source: doc.source,
      });
    } else if (activeTab === 'stories') {
      const story = item as StoryDTO;
      form.setFieldsValue({
        title: story.title,
        content: story.content,
        era: story.era,
        category: story.category,
      });
    } else {
      const teaching = item as FamilyTeachingsDTO;
      form.setFieldsValue({
        title: teaching.title,
        content: teaching.content,
        generation: teaching.generation,
        origin_text: teaching.origin_text,
        meaning: teaching.meaning,
      });
    }
    setModalVisible(true);
  };

  const handleView = (item: DocumentDTO | StoryDTO | FamilyTeachingsDTO) => {
    setSelectedItem(item);
    setDetailVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      if (activeTab === 'documents') await cultureApi.deleteDocument(id);
      else if (activeTab === 'stories') await cultureApi.deleteStory(id);
      else await cultureApi.deleteFamilyTeaching(id);
      messageApi.success('删除成功');
      if (activeTab === 'documents') fetchDocuments();
      else if (activeTab === 'stories') fetchStories();
      else fetchTeachings();
    } catch {
      messageApi.error('删除失败');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (activeTab === 'documents') {
        if (editingItem) {
          await cultureApi.updateDocument((editingItem as DocumentDTO).id, values);
        } else {
          await cultureApi.createDocument(values);
        }
      } else if (activeTab === 'stories') {
        if (editingItem) {
          await cultureApi.updateStory((editingItem as StoryDTO).id, values);
        } else {
          await cultureApi.createStory(values);
        }
      } else {
        if (editingItem) {
          await cultureApi.updateFamilyTeaching((editingItem as FamilyTeachingsDTO).id, values);
        } else {
          await cultureApi.createFamilyTeaching(values);
        }
      }
      messageApi.success(editingItem ? '更新成功' : '创建成功');
      setModalVisible(false);
      if (activeTab === 'documents') fetchDocuments();
      else if (activeTab === 'stories') fetchStories();
      else fetchTeachings();
    } catch {
      messageApi.error('操作失败');
    }
  };

  const documentColumns: ColumnsType<DocumentDTO> = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    {
      title: '标题',
      dataIndex: 'title',
      width: 200,
      render: (title: string, record) => (
        <Button type="link" onClick={() => handleView(record)}>
          {title}
        </Button>
      ),
    },
    {
      title: '分类',
      dataIndex: 'category',
      width: 100,
      render: (cat: string) => {
        const info = categoryMap[cat] || { label: cat, color: 'default' };
        return <Tag color={info.color}>{info.label}</Tag>;
      },
    },
    { title: '作者', dataIndex: 'author', width: 120 },
    { title: '朝代', dataIndex: 'dynasty', width: 80 },
    { title: '浏览', dataIndex: 'view_count', width: 70 },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: unknown, record: DocumentDTO) => (
        <Space size="small">
          <Button size="small" type="link" icon={<EyeOutlined />} onClick={() => handleView(record)}>
            详情
          </Button>
          <Button size="small" type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除此文献？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const storyColumns: ColumnsType<StoryDTO> = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    {
      title: '标题',
      dataIndex: 'title',
      width: 200,
      render: (title: string, record) => (
        <Button type="link" onClick={() => handleView(record)}>
          {title}
        </Button>
      ),
    },
    { title: '时代', dataIndex: 'era', width: 100 },
    { title: '分类', dataIndex: 'category', width: 100 },
    { title: '浏览', dataIndex: 'view_count', width: 70 },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: unknown, record: StoryDTO) => (
        <Space size="small">
          <Button size="small" type="link" icon={<EyeOutlined />} onClick={() => handleView(record)}>
            详情
          </Button>
          <Button size="small" type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除此故事？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const teachingColumns: ColumnsType<FamilyTeachingsDTO> = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    {
      title: '标题',
      dataIndex: 'title',
      width: 200,
      render: (title: string, record) => (
        <Button type="link" onClick={() => handleView(record)}>
          {title}
        </Button>
      ),
    },
    { title: '世代', dataIndex: 'generation', width: 80 },
    { title: '引用', dataIndex: 'usage_count', width: 70 },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_: unknown, record: FamilyTeachingsDTO) => (
        <Space size="small">
          <Button size="small" type="link" icon={<EyeOutlined />} onClick={() => handleView(record)}>
            详情
          </Button>
          <Button size="small" type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除此家训？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const renderModalForm = () => {
    if (activeTab === 'documents') {
      return (
        <>
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
            <Input placeholder="文献标题" />
          </Form.Item>
          <Form.Item name="content" label="内容" rules={[{ required: true, message: '请输入内容' }]}>
            <Input.TextArea rows={4} placeholder="文献内容" />
          </Form.Item>
          <Form.Item name="category" label="分类" rules={[{ required: true, message: '请选择分类' }]}>
            <Select placeholder="选择分类">
              <Select.Option value="classic">经典</Select.Option>
              <Select.Option value="genealogy">族谱</Select.Option>
              <Select.Option value="memorial">纪念</Select.Option>
              <Select.Option value="history">历史</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="author" label="作者">
            <Input placeholder="作者" />
          </Form.Item>
          <Form.Item name="dynasty" label="朝代">
            <Input placeholder="如：清朝" />
          </Form.Item>
          <Form.Item name="source" label="来源">
            <Input placeholder="来源" />
          </Form.Item>
        </>
      );
    } else if (activeTab === 'stories') {
      return (
        <>
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
            <Input placeholder="故事标题" />
          </Form.Item>
          <Form.Item name="content" label="内容" rules={[{ required: true, message: '请输入内容' }]}>
            <Input.TextArea rows={4} placeholder="故事内容" />
          </Form.Item>
          <Form.Item name="era" label="时代">
            <Input placeholder="如：明朝" />
          </Form.Item>
          <Form.Item name="category" label="分类">
            <Input placeholder="故事分类" />
          </Form.Item>
        </>
      );
    } else {
      return (
        <>
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
            <Input placeholder="家训标题" />
          </Form.Item>
          <Form.Item name="content" label="内容" rules={[{ required: true, message: '请输入内容' }]}>
            <Input.TextArea rows={4} placeholder="家训内容" />
          </Form.Item>
          <Form.Item name="generation" label="世代">
            <InputNumber placeholder="世代" min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="origin_text" label="原文">
            <Input.TextArea rows={2} placeholder="原文" />
          </Form.Item>
          <Form.Item name="meaning" label="释义">
            <Input.TextArea rows={2} placeholder="释义" />
          </Form.Item>
        </>
      );
    }
  };

  const renderDetail = () => {
    if (!selectedItem) return null;
    if (activeTab === 'documents') {
      const doc = selectedItem as DocumentDTO;
      const info = categoryMap[doc.category] || { label: doc.category, color: 'default' };
      return (
        <Descriptions column={1} bordered>
          <Descriptions.Item label="标题">{doc.title}</Descriptions.Item>
          <Descriptions.Item label="分类"><Tag color={info.color}>{info.label}</Tag></Descriptions.Item>
          <Descriptions.Item label="作者">{doc.author || '-'}</Descriptions.Item>
          <Descriptions.Item label="朝代">{doc.dynasty || '-'}</Descriptions.Item>
          <Descriptions.Item label="来源">{doc.source || '-'}</Descriptions.Item>
          <Descriptions.Item label="内容">{doc.content}</Descriptions.Item>
          <Descriptions.Item label="浏览次数">{doc.view_count}</Descriptions.Item>
          <Descriptions.Item label="收藏次数">{doc.collect_count}</Descriptions.Item>
          <Descriptions.Item label="创建时间">{doc.created_at}</Descriptions.Item>
        </Descriptions>
      );
    } else if (activeTab === 'stories') {
      const story = selectedItem as StoryDTO;
      return (
        <Descriptions column={1} bordered>
          <Descriptions.Item label="标题">{story.title}</Descriptions.Item>
          <Descriptions.Item label="时代">{story.era || '-'}</Descriptions.Item>
          <Descriptions.Item label="分类">{story.category || '-'}</Descriptions.Item>
          <Descriptions.Item label="内容">{story.content}</Descriptions.Item>
          <Descriptions.Item label="浏览次数">{story.view_count}</Descriptions.Item>
          <Descriptions.Item label="创建时间">{story.created_at}</Descriptions.Item>
        </Descriptions>
      );
    } else {
      const teaching = selectedItem as FamilyTeachingsDTO;
      return (
        <Descriptions column={1} bordered>
          <Descriptions.Item label="标题">{teaching.title}</Descriptions.Item>
          <Descriptions.Item label="世代">{teaching.generation}</Descriptions.Item>
          <Descriptions.Item label="内容">{teaching.content}</Descriptions.Item>
          <Descriptions.Item label="原文">{teaching.origin_text || '-'}</Descriptions.Item>
          <Descriptions.Item label="释义">{teaching.meaning || '-'}</Descriptions.Item>
          <Descriptions.Item label="引用次数">{teaching.usage_count}</Descriptions.Item>
          <Descriptions.Item label="创建时间">{teaching.created_at}</Descriptions.Item>
        </Descriptions>
      );
    }
  };

  const getTabLabel = () => {
    if (activeTab === 'documents') return '文献';
    if (activeTab === 'stories') return '故事';
    return '家训';
  };

  return (
    <div className={styles.container}>
      {contextHolder}
      <div className={styles.header}>
        <Title level={3}>文化管理</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          新增{getTabLabel()}
        </Button>
      </div>

      <Card className={styles.filterCard}>
        <Input.Search
          placeholder={`搜索${getTabLabel()}...`}
          allowClear
          enterButton={<SearchOutlined />}
          onSearch={setSearchText}
          style={{ width: 300 }}
        />
      </Card>

      <Tabs
        activeKey={activeTab}
        onChange={(key) => { setActiveTab(key); setSearchText(''); }}
        items={[
          { key: 'documents', label: '文献', icon: <BookOutlined /> },
          { key: 'stories', label: '故事' },
          { key: 'teachings', label: '家训' },
        ]}
      />

      {activeTab === 'documents' && (
        <Table columns={documentColumns} dataSource={filteredDocuments} rowKey="id" loading={loading} pagination={{ pageSize: 10, showTotal: (total) => `共 ${total} 条` }} />
      )}
      {activeTab === 'stories' && (
        <Table columns={storyColumns} dataSource={filteredStories} rowKey="id" loading={loading} pagination={{ pageSize: 10, showTotal: (total) => `共 ${total} 条` }} />
      )}
      {activeTab === 'teachings' && (
        <Table columns={teachingColumns} dataSource={filteredTeachings} rowKey="id" loading={loading} pagination={{ pageSize: 10, showTotal: (total) => `共 ${total} 条` }} />
      )}

      <Modal
        title={editingItem ? `编辑${getTabLabel()}` : `新增${getTabLabel()}`}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        destroyOnClose
        width={600}
      >
        <Form form={form} layout="vertical" preserve={false}>
          {renderModalForm()}
        </Form>
      </Modal>

      <Drawer
        title={`${getTabLabel()}详情`}
        open={detailVisible}
        onClose={() => setDetailVisible(false)}
        width={500}
      >
        {renderDetail()}
      </Drawer>
    </div>
  );
}
