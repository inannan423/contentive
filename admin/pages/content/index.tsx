import { Geist, Geist_Mono } from "next/font/google";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import Header from "@/components/layout/Header";
import Sidebar from "@/components/layout/Sidebar";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function Content() {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check if user is logged in
    const storedUser = localStorage.getItem("user");
    const token = localStorage.getItem("token");

    if (!storedUser || !token) {
      router.push("/auth/login");
      return;
    }

    try {
      setUser(JSON.parse(storedUser));
    } catch (error) {
      console.error("Failed to parse user data", error);
      router.push("/auth/login");
      return;
    }

    setLoading(false);
  }, [router]);

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    router.push("/auth/login");
  };

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
      <Header username={user?.username} onLogout={handleLogout} />
      <aside className="w-full h-full grid grid-cols-6">
        <Sidebar />
        <div className="col-span-5 flex flex-col justify-start items-center">
          <div className="w-full h-max flex justify-between items-center border-b-[1px] border-gray-200 border-dotted px-5 py-3">
            <h2 className="text-base font-semibold text-black dark:text-white">
              Content Management
            </h2>
          </div>
          
          <div className="max-w-4xl w-full h-full pt-5 px-5 border-l-[1px] border-r-[1px] border-gray-200 border-dotted">
            <div className="flex w-full justify-start items-center gap-2">
              <p className="font-semibold text-lg text-black dark:text-white">
                Content Types and Entries
              </p>
            </div>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-4">
              Manage your content types and entries here. Create, edit, and delete content as needed.
            </p>
          </div>
        </div>
      </aside>
    </div>
  );
}