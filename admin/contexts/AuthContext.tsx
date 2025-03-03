import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/router';
import { AuthUserType } from '@/types/user';

interface AuthContextType {
  user: AuthUserType | null;
  token: string | null;
  loading: boolean;
  logout: () => void;
  setUserAndToken: (user: AuthUserType, token: string) => void;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  token: null,
  loading: true,
  logout: () => {},
  setUserAndToken: () => {},
});

export const useAuth = () => useContext(AuthContext);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<AuthUserType | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  const validateToken = async (token: string) => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/auth/validate`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Token validation failed');
      }

      const data = await response.json();
      return data.valid;
    } catch (error) {
      console.error('Token validation error:', error);
      return false;
    }
  };

  useEffect(() => {
    const initAuth = async () => {
      try {
        if (router.pathname.startsWith('/auth/')) {
          setLoading(false);
          return;
        }

        const storedToken = localStorage.getItem('token');
        const storedUser = localStorage.getItem('user');

        if (!storedToken || !storedUser) {
          router.push('/auth/login');
          return;
        }

        if (!router.pathname.startsWith('/auth/')) {
          const isValid = await validateToken(storedToken);
          if (!isValid) {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            router.push('/auth/login');
            return;
          }
        }

        setToken(storedToken);
        setUser(JSON.parse(storedUser));
      } catch (error) {
        console.error('Failed to initialize auth context:', error);
        router.push('/auth/login');
      } finally {
        setLoading(false);
      }
    };

    initAuth();
  }, [router.pathname]);

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
    setToken(null);
    router.push('/auth/login');
  };

  const setUserAndToken = (newUser: AuthUserType, newToken: string) => {
    localStorage.setItem('user', JSON.stringify(newUser));
    localStorage.setItem('token', newToken);
    setUser(newUser);
    setToken(newToken);
  };

  return (
    <AuthContext.Provider value={{ user, token, loading, logout, setUserAndToken }}>
      {children}
    </AuthContext.Provider>
  );
};