'use client'; // Required for Context API

import React, { createContext, useState, useContext, ReactNode, useEffect } from 'react';

interface AuthContextType {
  isAuthenticated: boolean;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // ตรวจสอบ token เมื่อเริ่มแอพ
  useEffect(() => {
    // ตรวจสอบเฉพาะฝั่ง client เท่านั้น
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token');
      if (token) {
        console.log('Found existing token, setting authenticated state');
        setIsAuthenticated(true);
      }
    }
  }, []);

  const login = (token: string) => {
    // เก็บ token ใน localStorage
    localStorage.setItem('token', token);
    console.log("Login with token:", token);
    setIsAuthenticated(true);
  };

  const logout = () => {
    // ลบ token ออกจาก localStorage
    localStorage.removeItem('token');
    console.log("Logout");
    setIsAuthenticated(false);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}; 