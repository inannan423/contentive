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

  // If still loading, show nothing
  if (loading) return null;

  // If no user is logged in, redirect to login
  if (!user) {
    router.push('/auth/login');
    return null;
  }

  // If a specific role is required and user doesn't have it, redirect to home
  if (requiredRole && user.role !== requiredRole) {
    router.push('/');
    return null;
  }

  // If all checks pass, render the children
  return <React.Fragment>{children}</React.Fragment>;
}