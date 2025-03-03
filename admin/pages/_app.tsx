import "@/styles/globals.css";
import type { AppProps } from "next/app";
import { Toaster } from 'sonner';
import React from "react";
import { AuthProvider } from '@/contexts/AuthContext';

export default function AdminApp({ Component, pageProps }: AppProps) {
  return (
    <AuthProvider>
      <Component {...pageProps} />
      <Toaster position="bottom-right" />
    </AuthProvider>
  );
}
