'use client';

import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, message, Popconfirm, Spin, Typography, Layout, Menu } from 'antd';
import axios from 'axios';
import { useRouter } from 'next/navigation';
import { UserOutlined, LogoutOutlined } from '@ant-design/icons';

const { Header, Content, Footer } = Layout;
const { Title } = Typography;

// Main component for the dashboard page
const DashboardPage = () => {
    // State management
    const [users, setUsers] = useState([]); // Stores the list of users
    const [loading, setLoading] = useState(true); // For the loading spinner
    const [isModalVisible, setIsModalVisible] = useState(false); // To show/hide the edit modal
    const [editingUser, setEditingUser] = useState<any>(null); // To store the user being edited
    const [form] = Form.useForm(); // To control the form inside the modal
    const router = useRouter();

    // Helper function to get the JWT token from localStorage
    const getToken = () => localStorage.getItem('token');

    // Function to fetch all users from the backend
    const fetchUsers = async () => {
        const token = getToken();
        if (!token) {
            router.push('/login'); // Redirect to login if not authenticated
            return;
        }

        setLoading(true);
        try {
            const response = await axios.get('http://localhost:8089/users', {
                headers: { Authorization: `Bearer ${token}` },
            });
            setUsers(response.data || []); // Ensure it's an array
        } catch (error) {
            message.error('Gagal mengambil data user. Sesi Anda mungkin telah berakhir.');
            router.push('/login'); // Redirect on error (e.g., token expired)
        } finally {
            setLoading(false);
        }
    };

    // This effect runs once when the page loads
    useEffect(() => {
        fetchUsers();
    }, []);

    // Function to handle user deletion
    const handleDelete = async (id: number) => {
        try {
            const token = getToken();
            await axios.delete(`http://localhost:8089/users/${id}`, {
                headers: { Authorization: `Bearer ${token}` },
            });
            message.success('User berhasil dihapus');
            fetchUsers(); 
        } catch (error) {
            message.error('Gagal menghapus user.');
        }
    };
    
    const handleModalOk = () => {
        form.validateFields().then(async (values) => {
            try {
                const token = getToken();
                const payload = {
                    username: values.username,
                    ...(values.password && { password: values.password }),
                };

                await axios.put(`http://localhost:8089/users/${editingUser.id}`, payload, {
                     headers: { Authorization: `Bearer ${token}` },
                });
                
                message.success('User berhasil di-update');
                setIsModalVisible(false); 
                fetchUsers(); 
            } catch (error) {
                message.error('Gagal menyimpan user.');
            }
        });
    };
    
    const showEditModal = (user: any) => {
        setEditingUser(user);
        form.setFieldsValue({
            username: user.username,
            password: '', 
        });
        setIsModalVisible(true);
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        message.success('Logout berhasil!');
        router.push('/login');
    };

    const columns = [
        { title: 'ID', dataIndex: 'id', key: 'id', sorter: (a:any, b:any) => a.id - b.id, },
        { title: 'Username', dataIndex: 'username', key: 'username' },
        { 
            title: 'Created At', 
            dataIndex: 'created_at', 
            key: 'created_at', 
            render: (text:string) => new Date(text).toLocaleString('id-ID')
        },
        {
            title: 'Action',
            key: 'action',
            render: (_: any, record: any) => (
                <span>
                    <Button type="link" onClick={() => showEditModal(record)}>Edit</Button>
                    <Popconfirm title="Yakin ingin menghapus user ini?" onConfirm={() => handleDelete(record.id)} okText="Ya" cancelText="Tidak">
                        <Button type="link" danger>Delete</Button>
                    </Popconfirm>
                </span>
            ),
        },
    ];

    if (loading) {
        return <Spin size="large" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }} />;
    }

    return (
        <Layout style={{ minHeight: '100vh' }}>
            <Header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Title level={3} style={{ color: 'white', margin: 0 }}>Dashboard Admin</Title>
                <Button type="primary" danger icon={<LogoutOutlined />} onClick={handleLogout}>
                    Logout
                </Button>
            </Header>
            <Content style={{ padding: '24px 50px' }}>
                <div style={{ background: '#fff', padding: 24, minHeight: 280 }}>
                    <Title level={2}>Manajemen User</Title>
                    <Table 
                        columns={columns} 
                        dataSource={users.map((u: any) => ({...u, key: u.id}))} 
                        bordered
                    />
                </div>
            </Content>
            <Footer style={{ textAlign: 'center' }}>
                Admin Dashboard ©{new Date().getFullYear()}
            </Footer>

            {/* Edit User Modal */}
            <Modal
                title="Edit User"
                visible={isModalVisible}
                onOk={handleModalOk}
                onCancel={() => setIsModalVisible(false)}
                okText="Simpan"
                cancelText="Batal"
                destroyOnClose
            >
                <Form form={form} layout="vertical" name="edit_user_form" initialValues={{ username: editingUser?.username }}>
                    <Form.Item name="username" label="Username" rules={[{ required: true, message: 'Username tidak boleh kosong!' }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="password" label="Password Baru (Opsional)">
                        <Input.Password placeholder="Kosongkan jika tidak ingin ganti" />
                    </Form.Item>
                </Form>
            </Modal>
        </Layout>
    );
};

export default DashboardPage;

