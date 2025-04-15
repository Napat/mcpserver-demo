'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useRouter, useParams } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import apiClient from '@/services/api';
import axios, { AxiosError, AxiosResponse } from 'axios';

interface Note {
    id?: number; // Optional for new notes
    title: string;
    content: string;
}

const NoteFormPage = () => {
    const { isAuthenticated, logout } = useAuth();
    const router = useRouter();
    const params = useParams<{ id?: string }>(); // Get ID from URL
    const noteId = params?.id && params.id !== 'new' ? parseInt(params.id, 10) : null;
    const isNewNote = noteId === null;

    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [loading, setLoading] = useState(!isNewNote); // Load only if editing
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleLogout = useCallback(() => {
        logout();
        localStorage.removeItem('token');
        router.push('/login');
    }, [logout, router]);

    useEffect(() => {
        // Auth check
        if (typeof window !== 'undefined' && !isAuthenticated) {
            const token = localStorage.getItem('token');
            if (!token) {
                router.push('/login');
                return;
            }
        }

        if (!isNewNote && noteId) {
            setLoading(true);
            setError(null);
            apiClient.get<Note>(`/api/notes/${noteId}`)
                .then((response: AxiosResponse<Note>) => {
                    setTitle(response.data.title);
                    setContent(response.data.content);
                })
                .catch((err: unknown) => {
                    console.error("Failed to fetch note:", err);
                    let errorMessage = 'ไม่สามารถโหลดรายละเอียดบันทึกได้';
                     if (axios.isAxiosError(err)) {
                         const serverError = err as AxiosError<{ message?: string }>;
                         if (serverError.response?.status === 401) {
                             handleLogout();
                             return;
                         } else if (serverError.response?.status === 403) {
                            errorMessage = 'คุณไม่มีสิทธิ์ในการดูบันทึกนี้';
                         } else if (serverError.response?.status === 404) {
                            errorMessage = 'ไม่พบบันทึกที่ต้องการ';
                         } else if (serverError.response?.data?.message) {
                             errorMessage = serverError.response.data.message;
                         }
                     } else if (err instanceof Error) {
                         errorMessage = err.message;
                     }
                    setError(errorMessage);
                })
                .finally(() => {
                    setLoading(false);
                });
        }
    }, [isAuthenticated, router, noteId, isNewNote, handleLogout]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        setError(null);

        const noteData: Note = { title, content };
        const request = isNewNote
            ? apiClient.post('/api/notes', noteData)
            : apiClient.put(`/api/notes/${noteId}`, noteData);

        try {
            await request;
            router.push('/notes'); // Redirect to notes list after saving
        } catch (err: unknown) {
            console.error("Failed to save note:", err);
            let errorMessage = `ไม่สามารถ${isNewNote ? 'สร้าง' : 'อัปเดต'}บันทึกได้ กรุณาลองอีกครั้ง`;
            if (axios.isAxiosError(err)) {
                const serverError = err as AxiosError<{ message?: string }>;
                if (serverError.response?.status === 401) {
                    handleLogout();
                    return;
                } else if (serverError.response?.status === 403) {
                     errorMessage = `คุณไม่มีสิทธิ์ในการ${isNewNote ? 'สร้าง' : 'อัปเดต'}บันทึกนี้`;
                 } else if (serverError.response?.data?.message) {
                    errorMessage = serverError.response.data.message;
                }
            } else if (err instanceof Error) {
                errorMessage = err.message;
            }
            setError(errorMessage);
        } finally {
            setSaving(false);
        }
    };

    if (loading) {
        return <div className="text-center mt-10 text-gray-300">กำลังโหลดบันทึก...</div>;
    }

    // If editing and error occurred during load
    if (!isNewNote && error && !title) {
        return <div className="text-center mt-10 text-red-400">ข้อผิดพลาด: {error}</div>;
    }

    return (
        <div className="max-w-2xl mx-auto mt-10 bg-gray-800 p-8 rounded-lg shadow-md border border-gray-700">
            <h1 className="text-2xl font-bold mb-6 text-gray-100">{isNewNote ? 'สร้างบันทึกใหม่' : 'แก้ไขบันทึก'}</h1>
            <form onSubmit={handleSubmit}>
                {error && <p className="mb-4 text-red-400 text-sm">{error}</p>}
                <div className="mb-4">
                    <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="title">
                        หัวข้อ
                    </label>
                    <input
                        className="shadow appearance-none border bg-gray-700 border-gray-600 rounded w-full py-2 px-3 text-gray-100 leading-tight focus:outline-none focus:border-blue-500 focus:shadow-outline"
                        id="title"
                        type="text"
                        placeholder="หัวข้อบันทึก"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        required
                        disabled={saving}
                    />
                </div>
                <div className="mb-6">
                    <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="content">
                        เนื้อหา
                    </label>
                    <textarea
                        className="shadow appearance-none border bg-gray-700 border-gray-600 rounded w-full py-2 px-3 text-gray-100 leading-tight focus:outline-none focus:border-blue-500 focus:shadow-outline h-40"
                        id="content"
                        placeholder="เขียนบันทึกของคุณที่นี่..."
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                        required
                        disabled={saving}
                    />
                </div>
                <div className="flex items-center justify-end space-x-4">
                     <Link href="/notes" className="text-gray-400 hover:text-gray-200 text-sm">
                        ยกเลิก
                    </Link>
                    <button
                        className={`bg-blue-700 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline ${saving ? 'opacity-50 cursor-not-allowed' : ''}`}
                        type="submit"
                        disabled={saving}
                    >
                        {saving ? 'กำลังบันทึก...' : 'บันทึก'}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default NoteFormPage; 