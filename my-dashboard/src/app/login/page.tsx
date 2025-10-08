'use client';

import { Form, Input, Button, Card, message, Typography } from 'antd';
import axios from 'axios';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

const { Title } = Typography;

const LoginPage = () => {
    const router = useRouter();

    const onFinish = async (values: any) => {
        try {
            const response = await axios.post('http://localhost:8089/login', {
                username: values.username,
                password: values.password,
            });

            localStorage.setItem('token', response.data.token);
            message.success('Login berhasil!');
            
            router.push('/dashboard');
        } catch (error) {
            message.error('Login gagal! Cek kembali username dan password.');
        }
    };

    return (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh', background: '#f0f2f5' }}>
            <Card style={{ width: 350 }}>
                 <div style={{ textAlign: 'center', marginBottom: '24px' }}>
                    <Title level={3}>Selamat Datang</Title>
                </div>
                <Form name="login" onFinish={onFinish} layout="vertical">
                    <Form.Item label="Username" name="username" rules={[{ required: true, message: 'Username tidak boleh kosong!' }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item label="Password" name="password" rules={[{ required: true, message: 'Password tidak boleh kosong!' }]}>
                        <Input.Password />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" block>
                            Log in
                        </Button>
                    </Form.Item>
                    
                    <div style={{ textAlign: 'center' }}>
                        Belum punya akun? <Link href="/register">Daftar di sini</Link>
                    </div>

                </Form>
            </Card>
        </div>
    );
};

export default LoginPage;
