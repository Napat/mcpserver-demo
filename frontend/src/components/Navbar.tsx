'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

const Navbar = () => {
  const { isAuthenticated, logout } = useAuth();
  const router = useRouter();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = () => {
    logout();
    localStorage.removeItem('token'); // Clear token on logout
    router.push('/login'); // Redirect to login page
    setMenuOpen(false); // ปิดเมนูเมื่อออกจากระบบ
  };

  const toggleMenu = () => {
    setMenuOpen(!menuOpen);
  };

  return (
    <nav className="bg-gray-800 text-white p-4 shadow-md">
      <div className="container mx-auto">
        <div className="flex justify-between items-center">
          <Link href="/" className="font-bold text-xl hover:text-blue-300">
            MCPServer
          </Link>
          
          {/* Mobile menu button */}
          <button
            onClick={toggleMenu}
            className="md:hidden focus:outline-none"
            aria-label="Toggle menu"
          >
            <svg
              className="h-6 w-6 fill-current"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              {menuOpen ? (
                <path
                  fillRule="evenodd"
                  clipRule="evenodd"
                  d="M18.278 16.864a1 1 0 0 1-1.414 1.414l-4.829-4.828-4.828 4.828a1 1 0 0 1-1.414-1.414l4.828-4.829-4.828-4.828a1 1 0 0 1 1.414-1.414l4.829 4.828 4.828-4.828a1 1 0 1 1 1.414 1.414l-4.828 4.829 4.828 4.828z"
                />
              ) : (
                <path
                  fillRule="evenodd"
                  d="M4 5h16a1 1 0 0 1 0 2H4a1 1 0 1 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2zm0 6h16a1 1 0 0 1 0 2H4a1 1 0 0 1 0-2z"
                />
              )}
            </svg>
          </button>

          {/* Desktop navigation */}
          <div className="hidden md:block space-x-4">
            {isAuthenticated ? (
              <>
                <Link href="/profile" className="hover:text-blue-300">
                  Profile
                </Link>
                <Link href="/notes" className="hover:text-blue-300">
                  Notes
                </Link>
                <button
                  onClick={handleLogout}
                  className="bg-red-600 hover:bg-red-700 text-white font-bold py-1 px-3 rounded text-sm"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link href="/login" className="hover:text-blue-300">
                  Login
                </Link>
                <Link href="/register" className="hover:text-blue-300">
                  Register
                </Link>
              </>
            )}
          </div>
        </div>

        {/* Mobile navigation */}
        {menuOpen && (
          <div className="md:hidden mt-4 bg-gray-700 p-4 rounded-lg">
            {isAuthenticated ? (
              <div className="flex flex-col space-y-3">
                <Link 
                  href="/profile" 
                  className="hover:text-blue-300 py-2 px-3"
                  onClick={() => setMenuOpen(false)}
                >
                  Profile
                </Link>
                <Link 
                  href="/notes" 
                  className="hover:text-blue-300 py-2 px-3"
                  onClick={() => setMenuOpen(false)}
                >
                  Notes
                </Link>
                <button
                  onClick={handleLogout}
                  className="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-3 rounded text-sm text-left"
                >
                  Logout
                </button>
              </div>
            ) : (
              <div className="flex flex-col space-y-3">
                <Link 
                  href="/login" 
                  className="hover:text-blue-300 py-2 px-3"
                  onClick={() => setMenuOpen(false)}
                >
                  Login
                </Link>
                <Link 
                  href="/register" 
                  className="hover:text-blue-300 py-2 px-3"
                  onClick={() => setMenuOpen(false)}
                >
                  Register
                </Link>
              </div>
            )}
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar; 