import { Geist, Geist_Mono } from "next/font/google";
import { useRouter } from "next/router";
import { useEffect } from "react";
import React from "react";
import { useAuth } from "@/contexts/AuthContext";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function Home() {
  const router = useRouter();
  const { loading } = useAuth();

  useEffect(() => {
    if (!loading) {
      router.push("/playground");
    }
  }, [loading, router]);

  return (
    <div className={`${geistSans.variable} ${geistMono.variable} w-full h-screen bg-white dark:bg-black flex items-center justify-center`}>
      <p className="text-black dark:text-white">Loading...</p>
    </div>
  );
}
