'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import apiClient from '@/services/api';
import axios, { AxiosError } from 'axios';

// Define User type based on your models.User (adjust fields as needed)
interface UserProfile {
    id: number;
    email: string;
    first_name: string;
    last_name: string;
    gender: string;
    profile_image_url?: string;
    role: number; // Or a more specific role type if defined
    active: boolean;
    last_login_time?: string;
    created_at: string;
    updated_at: string;
}

// Define Login History type
interface LoginHistoryEntry {
    id: number;
    user_id: number;
    login_time: string;
    ip_address: string;
    user_agent: string;
    created_at: string;
}


const ProfilePage = () => {
    const { isAuthenticated, logout } = useAuth();
    const router = useRouter();
    const [user, setUser] = useState<UserProfile | null>(null);
    const [loginHistory, setLoginHistory] = useState<LoginHistoryEntry[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [historyLimit, setHistoryLimit] = useState(10);
    const [imageFile, setImageFile] = useState<File | null>(null);
    const [uploading, setUploading] = useState(false);
    const [uploadError, setUploadError] = useState<string | null>(null);

    const handleLogout = useCallback(() => {
        logout();
        localStorage.removeItem('token');
        router.push('/login');
    }, [logout, router]);

    useEffect(() => {
        // Initial check on client-side if context reports not authenticated
        if (typeof window !== 'undefined' && !isAuthenticated) {
            const token = localStorage.getItem('token');
            console.log('Profile page - Token in localStorage:', token ? 'exists' : 'not found');
            if (!token) {
                 router.push('/login');
                 return; // Stop execution if redirecting
            }
            // If token exists, proceed. Context might update or API call will handle expiry.
        }

        const fetchData = async () => {
            // Only fetch if authenticated or token exists (handled above)
            if (typeof window !== 'undefined' && (!isAuthenticated && !localStorage.getItem('token'))) {
                console.log('Not fetching profile data - not authenticated');
                return; // Don't fetch if definitely not logged in
            }

            console.log('Starting to fetch profile data...');
            setLoading(true);
            setError(null);
            try {
                // Fetch profile data - แก้ไขเส้นทาง API กลับไปเป็น /api/me
                console.log('Fetching profile from:', '/api/me');
                const profileRes = await apiClient.get<UserProfile>('/api/me');
                console.log('Profile data received:', profileRes.data);
                setUser(profileRes.data);

                // Fetch login history - แก้ไขเส้นทาง API กลับไปเป็น /api/me/login-history
                console.log('Fetching login history with limit:', historyLimit);
                const historyRes = await apiClient.get<LoginHistoryEntry[]>(`/api/me/login-history?limit=${historyLimit}`);
                console.log('Login history received:', historyRes.data);
                setLoginHistory(Array.isArray(historyRes.data) ? historyRes.data : []);

            } catch (err: unknown) {
                console.error("Failed to fetch profile data:", err);
                 let errorMessage = 'Failed to load profile data. Please try again.';
                 if (axios.isAxiosError(err)) {
                     const serverError = err as AxiosError<{ message?: string }>;
                     console.log('Error status:', serverError.response?.status);
                     console.log('Error data:', serverError.response?.data);
                     if (serverError.response?.status === 401) {
                         handleLogout(); // Logout user on 401
                         return; // Stop execution
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
        };

        console.log('Profile useEffect running, authenticated:', isAuthenticated);
        fetchData();
    }, [isAuthenticated, router, historyLimit, handleLogout]);


    const handleHistoryLimitChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const limit = parseInt(e.target.value, 10);
        if (!isNaN(limit) && limit > 0) {
            setHistoryLimit(limit);
        }
    };

     const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            setImageFile(e.target.files[0]);
            setUploadError(null); // Clear previous upload errors
        }
    };

    const handleImageUpload = async () => {
        if (!imageFile) {
            setUploadError("Please select an image file first.");
            return;
        }
        setUploading(true);
        setUploadError(null);

        const formData = new FormData();
        formData.append('image', imageFile);

        try {
            // แก้ไขเส้นทาง API กลับไปเป็น /api/me/profile-image
            console.log('Uploading profile image...');
            const response = await apiClient.post<{ image_url: string }>('/api/me/profile-image', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                },
            });
            console.log('Upload response:', response.data);
            // Update user profile state with new image URL
            if (user && response.data.image_url) {
                setUser({ ...user, profile_image_url: response.data.image_url });
            }
            setImageFile(null); // Clear file input after successful upload
            const fileInput = document.getElementById('profileImageInput') as HTMLInputElement;
            if (fileInput) fileInput.value = "";

        } catch (err: unknown) {
            console.error("Failed to upload profile image:", err);
            let errorMessage = 'Failed to upload image. Please try again.';
             if (axios.isAxiosError(err)) {
                 const serverError = err as AxiosError<{ message?: string }>;
                 if (serverError.response?.data?.message) {
                     errorMessage = serverError.response.data.message;
                 } else if (serverError.response?.status === 401) {
                    handleLogout();
                    return;
                 }
             } else if (err instanceof Error) {
                 errorMessage = err.message;
             }
            setUploadError(errorMessage);
        } finally {
            setUploading(false);
        }
    };

    // Render loading state or error before checking user
    if (loading) {
        return <div className="text-center mt-10">Loading profile...</div>;
    }

    if (error) {
        return <div className="text-center mt-10 text-red-500">Error: {error}</div>;
    }

    // If not loading and no error, but user is still null, it means redirect should have happened
    // or there's an issue. Show a generic message or redirect again.
    if (!user) {
        // This might indicate an issue if loading is false and error is null
        // For safety, redirect again or show a message
        if (typeof window !== 'undefined') router.push('/login');
        return <div className="text-center mt-10">Redirecting to login...</div>;
    }

    // User is loaded and available here
    return (
        <div className="bg-gray-800 p-6 rounded-lg shadow-md border border-gray-700">
            <h1 className="text-3xl font-bold mb-6 text-gray-100">Your Profile</h1>

            {/* Profile Details Section */}
            <div className="mb-8 p-4 border border-gray-700 rounded bg-gray-850">
                 <h2 className="text-xl font-semibold mb-4 text-gray-200">Details</h2>
                 <div className="flex flex-col sm:flex-row items-center mb-4">
                     <img
                        src={user.profile_image_url || '/default-avatar.png'} // Provide a default avatar path in /public
                        alt="Profile"
                        className="w-24 h-24 rounded-full mr-0 sm:mr-4 mb-4 sm:mb-0 object-cover border border-gray-600"
                    />
                    <div className="text-center sm:text-left text-gray-300">
                        <p><strong>Email:</strong> {user.email}</p>
                        <p><strong>Name:</strong> {user.first_name} {user.last_name}</p>
                        <p><strong>Gender:</strong> {user.gender}</p>
                        {/* Add more fields as needed */}
                    </div>
                 </div>

                {/* Profile Image Upload */}
                 <div className="mt-4 pt-4 border-t border-gray-700">
                    <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="profileImageInput">
                        Update Profile Image
                    </label>
                    <div className="flex flex-col sm:flex-row items-start sm:items-center space-y-2 sm:space-y-0 sm:space-x-2">
                        <input
                            id="profileImageInput"
                            type="file"
                            accept="image/*"
                            onChange={handleFileChange}
                            className="text-sm text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-gray-700 file:text-blue-400 hover:file:bg-gray-600 w-full sm:w-auto"
                            disabled={uploading}
                        />
                        <button
                            onClick={handleImageUpload}
                            disabled={!imageFile || uploading}
                            className={`bg-blue-700 hover:bg-blue-600 text-white font-bold py-1 px-3 rounded text-sm disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto`}
                        >
                            {uploading ? 'Uploading...' : 'Upload'}
                        </button>
                    </div>
                    {uploadError && <p className="mt-2 text-red-400 text-xs">{uploadError}</p>}
                 </div>
            </div>

            {/* Login History Section */}
            <div className="mb-8 p-4 border border-gray-700 rounded bg-gray-850">
                 <h2 className="text-xl font-semibold mb-4 text-gray-200">Login History</h2>
                 <div className="mb-4">
                    <label htmlFor="historyLimit" className="mr-2 text-sm font-medium text-gray-300">Show last:</label>
                    <input
                        type="number"
                        id="historyLimit"
                        value={historyLimit}
                        onChange={handleHistoryLimitChange}
                        min="1"
                        max="100" // Set a reasonable max limit
                        className="w-20 border rounded px-2 py-1 text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 bg-gray-700 border-gray-600 text-gray-200"
                    />
                    <span className="ml-2 text-sm text-gray-300">entries</span>
                 </div>
                {loading && loginHistory.length === 0 && <p className="text-gray-400">Loading history...</p>}
                {!loading && loginHistory.length === 0 && <p className="text-gray-400">No login history found.</p>}
                {loginHistory.length > 0 && (
                    <ul className="space-y-2 max-h-60 overflow-y-auto border border-gray-700 rounded p-2 bg-gray-800">
                        {loginHistory.map((entry) => (
                            <li key={entry.id} className="text-sm p-2 border-b border-gray-700 last:border-b-0 bg-gray-750 rounded-md">
                                <div className="text-gray-200"><strong>Time:</strong> {entry.login_time && new Date(entry.login_time).getFullYear() > 1970 ? new Date(entry.login_time).toLocaleString() : 'Invalid date'}</div>
                                <div className="text-gray-400"><strong>IP:</strong> {entry.ip_address}</div>
                                <div className="text-gray-400 text-xs break-all"><strong>Device:</strong> {entry.user_agent}</div>
                            </li>
                        ))}
                    </ul>
                )}
            </div>

            <button
                onClick={handleLogout}
                className="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full sm:w-auto"
            >
                Logout
            </button>
        </div>
    );
};

export default ProfilePage; 