import { Geist, Geist_Mono } from "next/font/google";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import Header from "@/components/layout/Header";
import Sidebar from "@/components/layout/Sidebar";
import { Button } from "@/components/ui/button";
import { HeaderNavItem, HeaderNav } from "@/components/headerNav";
import { IoPeopleCircleOutline, IoKeyOutline } from "react-icons/io5";
import UserDialog from "@/components/dialogs/UserDialog";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function Access() {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [users, setUsers] = useState<any[]>([]);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<any>(null);
  const [error, setError] = useState("");

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
      fetchUsers();
    } catch (error) {
      console.error("Failed to parse user data", error);
      router.push("/auth/login");
      return;
    }

    setLoading(false);
  }, [router]);

  const fetchUsers = async () => {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/users`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch users");
      }

      const data = await response.json();
      setUsers(data);
    } catch (error) {
      console.error("Error fetching users:", error);
      setError("获取用户列表失败");
    }
  };

  const handleAddUser = async (userData: any) => {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/users`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(userData),
      });

      if (!response.ok) {
        throw new Error("Failed to create user");
      }

      await fetchUsers();
      setIsDialogOpen(false);
    } catch (error) {
      console.error("Error creating user:", error);
      setError("创建用户失败");
    }
  };

  const handleEditUser = async (userData: any) => {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/users/${editingUser.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(userData),
      });

      if (!response.ok) {
        throw new Error("Failed to update user");
      }

      await fetchUsers();
      setIsDialogOpen(false);
      setEditingUser(null);
    } catch (error) {
      console.error("Error updating user:", error);
      setError("更新用户失败");
    }
  };

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

  const items: HeaderNavItem[] = [
    {
      title: "Admin User Management",
      icon: <IoPeopleCircleOutline />,
      component: <></>
    },
    {
      title: "API Role Management",
      icon: <IoKeyOutline />,
      component: <></>
    }
  ]

  return (
    <div
      className={`${geistSans.variable} ${geistMono.variable} w-full h-screen bg-white dark:bg-black flex flex-col`}
    >
      <Header username={user?.username} onLogout={handleLogout} />
      <aside className="w-full h-full grid grid-cols-6">
        <Sidebar />
        <div className="col-span-5 flex flex-col justify-start items-center">
          <div className="w-full h-max flex items-center gap-4 border-b-[1px] border-gray-200 border-dotted px-5 py-3">
            <h2 className="text-base w-max font-semibold text-black dark:text-white">
              Access Control
            </h2>

            <HeaderNav items={items} onChange={
              (component: any, index: number) => {
                console.log(component, index);
              }
            } />
            
          </div>
          
          <div className="max-w-4xl w-full flex flex-col h-full pt-5 px-5 border-l-[1px] border-r-[1px] border-gray-200 border-dotted">
            <div className="flex justify-between items-center mb-4">
                <div className="flex flex-col">
                    <h3 className="text-base font-semibold text-black dark:text-white">
                        List of all users with access to the admin panel
                    </h3>
                    <p className="text-sm mt-2 text-gray-600 dark:text-gray-400">
                        You can add or remove users from this list, or grant them access to the admin panel.
                    </p>
                </div>
                <Button onClick={() => {
                    setEditingUser(null);
                    setIsDialogOpen(true);
                }}>
                    Add User
                </Button>
            </div>

            <div className="w-full border-[1px] border-gray-200 rounded-lg overflow-hidden">
                <div className="flex w-full bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 py-3 px-4">
                    <div className="w-1/4 font-semibold text-sm text-gray-600 dark:text-gray-300">Username</div>
                    <div className="w-1/3 font-semibold text-sm text-gray-600 dark:text-gray-300">Email</div>
                    <div className="w-1/5 font-semibold text-sm text-gray-600 dark:text-gray-300">Status</div>
                    <div className="w-1/5 font-semibold text-sm text-gray-600 dark:text-gray-300 text-right">Action</div>
                </div>

                {users.map((user) => (
                  <div key={user.id} className="flex w-full text-base items-center border-b font-mono border-gray-100 dark:border-gray-800 py-2 px-4 hover:bg-gray-50 dark:hover:bg-gray-900/30">
                    <div className="w-1/4 font-medium text-black dark:text-white">{user.username}</div>
                    <div className="w-1/3 text-gray-600 dark:text-gray-400">{user.email}</div>
                    <div className="w-1/5">
                      <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${user.active ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'}`}>
                        {user.active ? "Active" : "Inactive"}
                      </span>
                    </div>
                    <div className="w-1/5 text-right">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 px-2 text-blue-600 dark:text-blue-400"
                        onClick={() => {
                          setEditingUser(user);
                          setIsDialogOpen(true);
                        }}
                      >
                        Edit
                      </Button>
                    </div>
                  </div>
                ))}
                
            </div>

            <UserDialog 
              isOpen={isDialogOpen}
              onClose={() => {
                setIsDialogOpen(false);
                setEditingUser(null);
              }}
              onSubmit={editingUser ? handleEditUser : handleAddUser}
            />
          </div>
        </div>
      </aside>
    </div>
  );
}