'use client';

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useEffect, useState } from "react";

export default function HomePage() {
  const { isAuthenticated, logout } = useAuth();
  const router = useRouter();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // เช็คสถานะ authentication ใน client-side เท่านั้น
    setLoading(false);
  }, [isAuthenticated]);

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  if (loading) {
    return (
      <div className="text-center">
        <p className="text-lg mb-6 text-gray-300 dark:text-gray-300">กำลังโหลด...</p>
      </div>
    );
  }

  return (
    <div className="text-center p-4 max-w-screen-md mx-auto dark:bg-gray-900">
      <h1 className="text-3xl sm:text-4xl font-bold mb-4 text-gray-100 dark:text-gray-100">ยินดีต้อนรับสู่ MCPServer</h1>
      <p className="text-base sm:text-lg mb-6 text-gray-300 dark:text-gray-300">
        นี่คือโปรเจคตัวอย่างที่สาธิตการใช้งาน Go backend ร่วมกับ Next.js frontend
      </p>
      
      {isAuthenticated ? (
        // แสดงเมื่อล็อกอินแล้ว
        <div className="space-y-4 sm:space-y-6">
          <p className="text-lg text-gray-200 dark:text-gray-200">คุณได้เข้าสู่ระบบแล้ว!</p>
          <div className="flex flex-col sm:flex-row justify-center gap-3 sm:space-x-4">
            <Link href="/notes" className="bg-gray-700 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-white">
              บันทึกของฉัน
            </Link>
            <Link href="/profile" className="bg-gray-700 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-white">
              โปรไฟล์
            </Link>
            <button
              onClick={handleLogout}
              className="bg-red-700 hover:bg-red-600 text-white font-bold py-2 px-4 rounded dark:bg-red-700 dark:hover:bg-red-600 dark:text-white"
            >
              ออกจากระบบ
            </button>
          </div>
        </div>
      ) : (
        // แสดงเมื่อยังไม่ได้ล็อกอิน
        <div className="flex flex-col sm:flex-row justify-center gap-3 sm:space-x-4">
          <Link href="/login" className="bg-gray-700 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-white">
            เข้าสู่ระบบ
          </Link>
          <Link href="/register" className="bg-blue-700 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded dark:bg-blue-700 dark:hover:bg-blue-600 dark:text-white">
            ลงทะเบียน
          </Link>
        </div>
      )}
    </div>
  );
}
