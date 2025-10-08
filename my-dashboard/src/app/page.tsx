// src/app/page.tsx
'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Spin } from 'antd';

const RootPage = () => {
    const router = useRouter();

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            router.replace('/home');
        } else {
            router.replace('/register');
        }
    }, [router]); 

    return (
        <Spin 
            size="large" 
            style={{ 
                display: 'flex', 
                justifyContent: 'center', 
                alignItems: 'center', 
                height: '100vh' 
            }} 
        />
    );
};

export default RootPage;