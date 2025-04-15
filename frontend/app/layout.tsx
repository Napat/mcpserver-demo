import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Navbar from "@/components/Navbar";
import { AuthProvider } from "@/contexts/AuthContext";
import VisitorCounter from "@/components/VisitorCounter";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "MCPServer Frontend",
  description: "Frontend for the MCPServer sample project",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className={`${inter.className} bg-gray-900 text-gray-100 dark:bg-gray-900 dark:text-gray-100`}>
        <AuthProvider>
          <Navbar />
          <main className="container mx-auto p-4 mt-4">
            {children}
          </main>
          <footer className="text-center text-gray-400 text-sm mt-8 p-4 border-t border-gray-800">
            <VisitorCounter />
            <p className="mt-2">
              © {new Date().getFullYear()} MCPServer Sample Project
            </p>
          </footer>
        </AuthProvider>
      </body>
    </html>
  );
}
