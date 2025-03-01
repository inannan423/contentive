import { Geist, Geist_Mono } from "next/font/google";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

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
  useEffect(() => {
    // Check if user is logged in
    const storedUser = localStorage.getItem("user");
    const token = localStorage.getItem("token");

    if (!storedUser || !token) {
      router.push("/auth/login");
      return;
    }

    // Redirect to playground page
    router.push("/playground");
  }, [router]);

  return (
    <div className={`${geistSans.variable} ${geistMono.variable} w-full h-screen bg-white dark:bg-black flex items-center justify-center`}>
      <p className="text-black dark:text-white">Loading...</p>
    </div>
  );
}
