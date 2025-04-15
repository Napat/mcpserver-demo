import axios from 'axios';

// Determine the base URL for the API
// In development, it might be the Go backend directly (e.g., http://localhost:8080)
// In production (with static export), it might be the same origin or a configured URL
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'; // Point to Go backend without /api

console.log('API Client initialized with base URL:', API_BASE_URL);

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to include the JWT token if available
apiClient.interceptors.request.use(
  (config) => {
    // Check if running on the client side
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token');
      if (token) {
        // Ensure headers object exists
        config.headers = config.headers || {};
        config.headers.Authorization = `Bearer ${token}`;
        console.log(`API Request to ${config.url} with token`);
      } else {
        console.log(`API Request to ${config.url} without token`);
      }
    }
    return config;
  },
  (error) => {
    console.error('API Request interceptor error:', error);
    return Promise.reject(error);
  }
);

// Add a response interceptor to log responses and errors
apiClient.interceptors.response.use(
  (response) => {
    console.log(`API Response from ${response.config.url}: Status ${response.status}`);
    return response;
  },
  (error) => {
    if (error.response) {
      console.error(`API Error from ${error.config.url}: Status ${error.response.status}`, error.response.data);
    } else if (error.request) {
      console.error('API Error: No response received', error.request);
    } else {
      console.error('API Error:', error.message);
    }
    return Promise.reject(error);
  }
);

export default apiClient;

// Define types for API responses (optional but recommended)
// Example:
// export interface UserProfile {
//   id: number;
//   email: string;
//   first_name: string;
//   last_name: string;
//   // ... other fields
// }

// export interface LoginResponse {
//   token: string;
//   user: UserProfile;
// } 