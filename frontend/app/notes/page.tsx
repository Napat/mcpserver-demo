'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import apiClient from '@/services/api';
import axios, { AxiosError } from 'axios';

interface Note {
    id: number;
    title: string;
    content: string;
    user_id: number;
    created_at: string;
    updated_at: string;
}

const NotesPage = () => {
    const { isAuthenticated, logout } = useAuth();
    const router = useRouter();
    const [notes, setNotes] = useState<Note[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const handleLogout = useCallback(() => {
        logout();
        localStorage.removeItem('token');
        router.push('/login');
    }, [logout, router]);

    const fetchNotes = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            console.log('กำลังดึงข้อมูลบันทึก...');
            const response = await apiClient.get<Note[]>('/api/notes');
            console.log('ได้รับข้อมูลบันทึกแล้ว:', response.data);
            setNotes(response.data);
        } catch (err: unknown) {
            console.error("ไม่สามารถดึงข้อมูลบันทึกได้:", err);
            let errorMessage = 'ไม่สามารถโหลดบันทึกได้ กรุณาลองอีกครั้ง';
            if (axios.isAxiosError(err)) {
                const serverError = err as AxiosError<{ message?: string }>;
                if (serverError.response?.status === 401) {
                    handleLogout();
                    return;
                } else if (serverError.response?.data?.message) {
                    errorMessage = serverError.response.data.message;
                }
            } else if (err instanceof Error) {
                errorMessage = err.message;
            }
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    }, [handleLogout]);

    useEffect(() => {
        // Initial auth check
        if (typeof window !== 'undefined' && !isAuthenticated) {
            const token = localStorage.getItem('token');
            if (!token) {
                router.push('/login');
                return;
            }
        }
        fetchNotes();
    }, [isAuthenticated, router, fetchNotes]);

    const handleDelete = async (id: number) => {
        if (!confirm('คุณแน่ใจหรือไม่ว่าต้องการลบบันทึกนี้?')) {
            return;
        }
        try {
            await apiClient.delete(`/api/notes/${id}`);
            setNotes(notes.filter(note => note.id !== id)); // ลบบันทึกออกจาก state
        } catch (err: unknown) {
            console.error("ไม่สามารถลบบันทึกได้:", err);
            let errorMessage = 'ไม่สามารถลบบันทึกได้ กรุณาลองอีกครั้ง';
            if (axios.isAxiosError(err)) {
                const serverError = err as AxiosError<{ message?: string }>;
                if (serverError.response?.status === 401) {
                    handleLogout();
                } else if (serverError.response?.status === 403) {
                    errorMessage = 'คุณไม่มีสิทธิ์ในการลบบันทึกนี้';
                } else if (serverError.response?.data?.message) {
                    errorMessage = serverError.response.data.message;
                }
            } else if (err instanceof Error) {
                errorMessage = err.message;
            }
            alert(errorMessage); // แสดงข้อผิดพลาดในการลบด้วย alert เพื่อความง่าย
        }
    };

    if (loading) {
        return <div className="text-center mt-10 text-gray-300">กำลังโหลดบันทึก...</div>;
    }

    if (error) {
        return <div className="text-center mt-10 text-red-400">ข้อผิดพลาด: {error}</div>;
    }

    return (
        <div className="container mx-auto p-4 bg-gray-800 min-h-screen">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold text-gray-100">บันทึกของฉัน</h1>
                <div>
                    <Link
                        href="/notes/new"
                        className="bg-blue-700 hover:bg-blue-800 text-white px-4 py-2 rounded mr-2"
                    >
                        เพิ่มบันทึกใหม่
                    </Link>
                    <button
                        onClick={handleLogout}
                        className="bg-gray-700 hover:bg-gray-600 text-gray-300 px-4 py-2 rounded"
                    >
                        ออกจากระบบ
                    </button>
                </div>
            </div>

            {notes.length === 0 ? (
                <div className="text-center py-10 text-gray-400">ไม่มีบันทึก เริ่มสร้างบันทึกแรกเลย!</div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {notes.map((note) => (
                        <div key={note.id} className="bg-gray-700 p-4 rounded shadow border border-gray-600">
                            <h2 className="text-xl font-bold mb-2 text-gray-100">{note.title}</h2>
                            <p className="text-gray-300 mb-4">{note.content}</p>
                            <div className="flex justify-end">
                                <Link
                                    href={`/notes/${note.id}`}
                                    className="text-blue-400 hover:text-blue-300 mr-4"
                                >
                                    แก้ไข
                                </Link>
                                <button
                                    onClick={() => handleDelete(note.id)}
                                    className="text-red-400 hover:text-red-300"
                                >
                                    ลบ
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default NotesPage; 