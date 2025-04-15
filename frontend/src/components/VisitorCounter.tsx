'use client';

import { useEffect, useState } from 'react';

export default function VisitorCounter() {
  const [visitorCount, setVisitorCount] = useState<number>(0);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const incrementVisitorCount = async () => {
      setIsLoading(true);
      setError(null);
      
      try {
        console.log('Trying to increment visitor count...');
        // ใช้ URL ที่กำหนดใน .env.local หรือ default ถ้าไม่มี
        const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
        console.log('Using API URL:', API_URL);
        
        // เรียกใช้ API เพื่อเพิ่มจำนวนผู้เข้าชม
        const response = await fetch(`${API_URL}/api/visitors`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          // มั่นใจว่า credentials ถูกส่งไปด้วย
          credentials: 'include',
          mode: 'cors',
        });
        
        if (!response.ok) {
          throw new Error(`Server responded with status: ${response.status}`);
        }
        
        const data = await response.json();
        console.log('Response data:', data);
        setVisitorCount(data.visitor_count);
      } catch (error) {
        console.error('Error incrementing visitor count:', error);
        setError(error instanceof Error ? error.message : 'Unknown error');
        
        // หากมีข้อผิดพลาด ลองดึงข้อมูลจำนวนโดยไม่เพิ่ม
        try {
          console.log('Trying to get visitor count instead...');
          const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
          
          const response = await fetch(`${API_URL}/api/visitors`, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            },
            credentials: 'include',
            mode: 'cors',
          });
          
          if (!response.ok) {
            throw new Error(`Server responded with status: ${response.status}`);
          }
          
          const data = await response.json();
          console.log('GET Response data:', data);
          setVisitorCount(data.visitor_count);
        } catch (getError) {
          console.error('Error getting visitor count:', getError);
        }
      } finally {
        setIsLoading(false);
      }
    };

    incrementVisitorCount();
  }, []);

  return (
    <div className="text-center p-2 bg-gray-800 rounded-lg">
      {isLoading ? (
        <p className="text-sm text-gray-400">กำลังโหลดข้อมูลผู้เข้าชม...</p>
      ) : error ? (
        <div>
          <p className="text-sm text-red-400">ไม่สามารถโหลดข้อมูลผู้เข้าชม</p>
          <p className="text-xs text-red-500">{error}</p>
        </div>
      ) : (
        <p className="text-sm text-gray-300">
          <span className="font-semibold">จำนวนผู้เข้าชม:</span> {visitorCount.toLocaleString()} คน
        </p>
      )}
    </div>
  );
} 