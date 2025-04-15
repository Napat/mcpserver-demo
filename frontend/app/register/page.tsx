'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import apiClient from '@/services/api';
import axios, { AxiosError } from 'axios';

// Define a type for the expected register response (similar to login)
interface RegisterResponse {
  token: string;
  // user: any; // Add user type if needed
}

const RegisterPage = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [gender, setGender] = useState('other'); // Default gender
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const response = await apiClient.post<RegisterResponse>('/api/auth/register', {
        email,
        password,
        first_name: firstName,
        last_name: lastName,
        gender,
      });

      if (response.data && response.data.token) {
        login(response.data.token); // Login user immediately after registration
        localStorage.setItem('token', response.data.token);
        router.push('/profile'); // Redirect to profile
      } else {
        setError('Registration failed: Invalid response from server.');
      }
    } catch (err: unknown) {
      console.error("Registration error:", err);
      let errorMessage = 'An unexpected error occurred. Please try again.';

      if (axios.isAxiosError(err)) {
        const serverError = err as AxiosError<{ message?: string }>;
        if (serverError.response?.data?.message) {
          errorMessage = serverError.response.data.message;
        } else if (serverError.response?.status === 400) {
             // You might get more specific errors from the backend for 400
             errorMessage = 'Registration failed. Please check your input.';
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
      <h1 className="text-2xl font-bold mb-6 text-center text-gray-100">Register</h1>
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
            placeholder="Minimum 6 characters"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={6}
            disabled={loading}
          />
        </div>

        <div className="flex flex-col sm:flex-row sm:space-x-4 space-y-4 sm:space-y-0">
          <div className="w-full">
            <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="firstName">
              First Name
            </label>
            <input
              className="shadow appearance-none border bg-gray-700 border-gray-600 rounded w-full py-2 px-3 text-gray-100 focus:outline-none focus:border-blue-500 focus:shadow-outline"
              id="firstName"
              type="text"
              placeholder="Your First Name"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              required
              disabled={loading}
            />
          </div>

          <div className="w-full">
            <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="lastName">
              Last Name
            </label>
            <input
              className="shadow appearance-none border bg-gray-700 border-gray-600 rounded w-full py-2 px-3 text-gray-100 focus:outline-none focus:border-blue-500 focus:shadow-outline"
              id="lastName"
              type="text"
              placeholder="Your Last Name"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
              required
              disabled={loading}
            />
          </div>
        </div>

        <div>
          <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="gender">
            Gender
          </label>
          <select
            id="gender"
            className="shadow border bg-gray-700 border-gray-600 rounded w-full py-2 px-3 text-gray-100 focus:outline-none focus:border-blue-500 focus:shadow-outline"
            value={gender}
            onChange={(e) => setGender(e.target.value)}
            disabled={loading}
          >
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>
        </div>

        <div className="pt-2">
          <button
            className={`bg-blue-700 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
            type="submit"
            disabled={loading}
          >
            {loading ? 'Registering...' : 'Register'}
          </button>
        </div>
        <p className="text-center text-gray-400 text-sm mt-4">
          Already have an account?{' '}
          <Link href="/login" className="text-blue-400 hover:text-blue-300">
            Login here
          </Link>
        </p>
      </form>
    </div>
  );
};

export default RegisterPage; 