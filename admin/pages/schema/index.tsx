import { Geist, Geist_Mono } from "next/font/google";
import Header from "@/components/layout/Header";
import Sidebar from "@/components/layout/Sidebar";
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

export default function Schema() {
  const { user, loading, logout } = useAuth();

  if (loading) {
    return (
      <div className={`${geistSans.variable} ${geistMono.variable} w-full h-screen bg-white dark:bg-black flex items-center justify-center`}>
        <p className="text-black dark:text-white">Loading...</p>
      </div>
    );
  }

  return (
    <div
      className={`${geistSans.variable} ${geistMono.variable} w-full h-screen bg-white dark:bg-black flex flex-col`}
    >
      <Header username={user?.username} onLogout={logout} />
      <aside className="w-full h-full grid grid-cols-6">
        <Sidebar />
        <div className="col-span-5 flex flex-col justify-start items-center">
          <div className="w-full h-max flex justify-between items-center border-b-[1px] border-gray-200 border-dotted px-5 py-3">
            <h2 className="text-base font-semibold text-black dark:text-white">
              Schema Builder
            </h2>
          </div>
          
          <div className="poo w-full h-full grid grid-cols-5">
            <div className="h-full border-r-[1px] border-gray-200 border-dotted">

            </div>
            <div className="col-span-4 h-full">

            </div>
          </div>
        </div>
      </aside>
    </div>
  );
}