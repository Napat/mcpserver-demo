'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import apiClient from '@/services/api';
import axios, { AxiosError } from 'axios';

// Define a type for the expected login response
interface LoginResponse {
  token: string;
  // user: any; // Add user type if needed
}

const LoginPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      // แก้ไขเส้นทาง API กลับไปเป็น /api/auth/login
      console.log('Attempting login with:', { email });
      const response = await apiClient.post<LoginResponse>('/api/auth/login', { email, password });
      console.log('Login response received:', response.data ? 'has data' : 'no data');

      if (response.data && response.data.token) {
        login(response.data.token); // Update auth state
        localStorage.setItem('token', response.data.token); // Store token
        router.push('/profile'); // Redirect to profile page after login
      } else {
        setError('Login failed: Invalid response from server.');
      }
    } catch (err: unknown) {
      console.error("Login error:", err);
      let errorMessage = 'An unexpected error occurred. Please try again.';

      if (axios.isAxiosError(err)) {
        const serverError = err as AxiosError<{ message?: string }>;
        if (serverError.response) {
          if (serverError.response.data && typeof serverError.response.data.message === 'string') {
            errorMessage = serverError.response.data.message;
          } else if (serverError.response.status === 401) {
            errorMessage = 'Invalid credentials. Please try again.';
          }
        }
      } else if (err instanceof Error) {
        errorMessage = err.message;
      }

      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full sm:max-w-md mx-auto my-4 sm:mt-10 bg-gray-800 p-6 sm:p-8 rounded-lg shadow-md border border-gray-700">
      <h1 className="text-2xl font-bold mb-6 text-center text-gray-100">Login</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && <p className="mb-4 text-red-400 text-sm text-center">{error}</p>}
        <div>
          <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="email">
            Email
          </label>
          <input
            className="shadow appearance-none border bg-gray-700 border-gray-600 rounded w-full py-2 px-3 text-gray-100 focus:outline-none focus:border-blue-500 focus:shadow-outline"
            id="email"
            type="email"
            placeholder="you@example.com"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            disabled={loading}
          />
        </div>
        <div>
          <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="password">
            Password
          </label>
          <input
            className="shadow appearance-none border bg-gray-700 border-gray-600 rounded w-full py-2 px-3 text-gray-100 focus:outline-none focus:border-blue-500 focus:shadow-outline"
            id="password"
            type="password"
            placeholder="******************"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            disabled={loading}
          />
        </div>
        <div className="pt-2">
          <button
            className={`bg-blue-700 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
            type="submit"
            disabled={loading}
          >
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </div>
        <p className="text-center text-gray-400 text-sm mt-4">
          Don&apos;t have an account?{' '}
          <Link href="/register" className="text-blue-400 hover:text-blue-300">
            Register here
          </Link>
        </p>
      </form>
    </div>
  );
};

export default LoginPage; 