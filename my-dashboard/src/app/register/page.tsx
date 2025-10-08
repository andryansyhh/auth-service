// src/app/register/page.tsx
'use client';

import { Form, Input, Button, Card, message, Typography } from 'antd';
import axios from 'axios';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

const { Title } = Typography;

const RegisterPage = () => {
    const router = useRouter();

    const onFinish = async (values: any) => {
        if (values.password !== values.confirmPassword) {
            message.error("Password dan konfirmasi password tidak cocok!");
            return;
        }

        try {
            await axios.post('http://localhost:8089/register', {
                username: values.username,
                password: values.password,
            });

            message.success('Registrasi berhasil! Silakan login.');
            
            router.push('/login');
        } catch (error: any) {
            const errorMessage = error.response?.data?.error || 'Registrasi gagal! Username mungkin sudah ada.';
            message.error(errorMessage);
        }
    };

    return (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh', background: '#f0f2f5' }}>
            <Card style={{ width: 350 }}>
                <div style={{ textAlign: 'center', marginBottom: '24px' }}>
                    <Title level={3}>Buat Akun Baru</Title>
                </div>
                <Form name="register" onFinish={onFinish} layout="vertical">
                    <Form.Item label="Username" name="username" rules={[{ required: true, message: 'Username tidak boleh kosong!' }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item label="Password" name="password" rules={[{ required: true, message: 'Password tidak boleh kosong!' }]}>
                        <Input.Password />
                    </Form.Item>
                    <Form.Item label="Konfirmasi Password" name="confirmPassword" rules={[{ required: true, message: 'Konfirmasi password tidak boleh kosong!' }]}>
                        <Input.Password />
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" block>
                            Daftar
                        </Button>
                    </Form.Item>
                    <div style={{ textAlign: 'center' }}>
                        Sudah punya akun? <Link href="/login">Login di sini</Link>
                    </div>
                </Form>
            </Card>
        </div>
    );
};

export default RegisterPage;