import { useRouter } from 'next/router';
import { useAuth } from '@/contexts/AuthContext';
import { ReactNode } from 'react';
import React from 'react';

interface PermissionGuardProps {
  children: ReactNode;
  requiredRole?: 'super_admin' | 'content_admin' | 'editor' | 'viewer';
}

export default function PermissionGuard({ children, requiredRole }: PermissionGuardProps) {
  const router = useRouter();
  const { user, loading } = useAuth();

  // If still loading, show loading state with full height
  if (loading) {
    return (
      <div className="w-full h-full flex items-center justify-center">
        <p className="text-black dark:text-white">Loading...</p>
      </div>
    );
  }

  // If no user is logged in, redirect to login
  if (!user) {
    router.push('/auth/login');
    return (
      <div className="w-full h-full flex items-center justify-center">
        <p className="text-black dark:text-white">Redirecting...</p>
      </div>
    );
  }

  // If a specific role is required and user doesn't have it, redirect to home
  if (requiredRole && user.role !== requiredRole) {
    router.push('/');
    return (
      <div className="w-full h-full flex items-center justify-center">
        <p className="text-black dark:text-white">Access denied. Redirecting...</p>
      </div>
    );
  }

  // If all checks pass, render the children with full height wrapper
  return <div className="w-full h-full overflow-auto">{children}</div>;
}