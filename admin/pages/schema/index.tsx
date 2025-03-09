import { Geist, Geist_Mono } from "next/font/google";
import Header from "@/components/layout/Header";
import Sidebar from "@/components/layout/Sidebar";
import React from "react";
import { useAuth } from "@/contexts/AuthContext";
import { ContentTypeType } from "@/types/content_type";
import { toast } from "sonner";
import { useRouter } from "next/router";


const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function Schema() {
  const { user, loading, token, logout } = useAuth();
  const router = useRouter();

  const [contentTypes, setContentTypes] = React.useState<ContentTypeType[]>();

  React.useEffect(() => {
    if (!loading) {

      fetchContentTypes();
    }
  }, [loading, token, router]);

  const fetchContentTypes = async () => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/content-types`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        toast.dismiss();
        toast.error("Failed to fetch content types");
      } else {
        const data = await response.json();
        setContentTypes(data);
        console.log(data);
      }
    } catch (error) {
      console.error("Error fetching users:", error);
      toast.dismiss();
      toast.error("Failed to fetch users");
    }
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
              {
                contentTypes && contentTypes.map((contentType) => {
                  return (
                    <div key={contentType.ID} className="w-full h-full flex flex-col justify-start items-center">
                      <div className="w-full h-max flex justify-between items-center border-b-[1px] border-gray-200 border-dotted px-5 py-3">
                        <h2 className="text-base font-semibold text-black dark:text-white">
                          {contentType.name}
                        </h2>
                      </div>
                    </div>
                  )
                })
              }
            </div>
            <div className="col-span-4 h-full">

            </div>
          </div>
        </div>
      </aside>
    </div>
  );
}