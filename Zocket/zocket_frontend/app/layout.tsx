import './globals.css';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';

// Configure font with display: swap to prevent FOIT
const inter = Inter({
  subsets: ['latin'],
  display: 'swap',
  preload: true,
  fallback: ['system-ui', 'arial']
});

export const metadata: Metadata = {
  title: 'Task Management System',
  description: 'Manage your tasks efficiently with real-time updates and AI-powered suggestions',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>{children}</body>
    </html>
  );
}