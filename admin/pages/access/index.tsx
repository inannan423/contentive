import React from "react";
import { Geist, Geist_Mono } from "next/font/google";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import Header from "@/components/layout/Header";
import Sidebar from "@/components/layout/Sidebar";
import { Button } from "@/components/ui/button";
import { HeaderNavItem, HeaderNav } from "@/components/headerNav";
import { IoPeopleCircleOutline, IoKeyOutline } from "react-icons/io5";
import {
  Table,
  TableHeader,
  TableBody,
  TableHead,
  TableRow,
  TableCell,
} from "@/components/ui/table";
import UserSheet from "@/components/sheets/UserSheet";
import { UserType } from "@/types/user";
import { CreateUserType, UpdateUserType } from "@/types/user";
import { APIRoleType, CreateAPIRoleType, UpdateAPIRoleType } from "@/types/api_role";
import APIRoleSheet from "@/components/sheets/APIRoleSheet";
import { toast } from "sonner";
import { useAuth } from "@/contexts/AuthContext";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import PermissionGuard from "@/components/auth/PermissionGuard";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const isSuperAdmin = (user: UserType) => {
  return user.role.Type === "super_admin";
};

export default function Access() {
  return (
    <PermissionGuard requiredRole="super_admin">
      <AccessContent />
    </PermissionGuard>
  );
}

function AccessContent() {
  const router = useRouter();
  const { user, token, loading, logout } = useAuth();
  const [users, setUsers] = useState<UserType[]>([]);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<UserType | null>();
  const [isSubmitting, setIsSubmitting] = useState(false);
//   const [deletingUserId, setDeletingUserId] = useState<string | null>(null);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [apiRoles, setApiRoles] = useState<APIRoleType[]>([]);
  const [isApiRoleDialogOpen, setIsApiRoleDialogOpen] = useState(false);
  const [editingApiRole, setEditingApiRole] = useState<APIRoleType | null>(null);
  const [isApiRoleSubmitting, setIsApiRoleSubmitting] = useState(false);

  useEffect(() => {
    if (!loading && !token) {
      router.push("/auth/login");
      return;
    }

    if (!loading) {
      fetchUsers();
      fetchApiRoles();
    }
  }, [loading, token, router]);

  const fetchUsers = async () => {
    try {
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
      toast.dismiss();
      toast.error("Failed to fetch users");
    }
  };

  const handleAddUser = async (userData: CreateUserType | UpdateUserType) => {
    try {
      setIsSubmitting(true);
      const createData = userData as CreateUserType;
      
      toast.loading("Creating user...");
      
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/users`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(createData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        // console.log(errorData)
        // throw new Error(errorData.message || "Failed to create user");
        toast.dismiss();
        toast.error(errorData.error || "Failed to create user");
      } else {
        await fetchUsers();
        setIsDialogOpen(false);
        toast.dismiss();
        toast.success("User created successfully");
      }
    } catch (error) {
      console.error("Error creating user:", error);
    //   setError("Failed to create user");
      toast.dismiss();
      toast.error("Failed to create user");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEditUser = async (userData: CreateUserType | UpdateUserType) => {
    try {
      setIsSubmitting(true);
      const updateData = userData as UpdateUserType;

      console.log("Updating user:", updateData);
      
      toast.loading("Saving changes...");
      
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/users/${editingUser?.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          username: updateData.username,
          email: updateData.email,
          role_id: updateData.role_id,
          active: updateData.active,
          ...(updateData.password ? { password: updateData.password } : {})
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        toast.dismiss();
        toast.error(errorData.error || "Failed to update user");
      } else {
        await fetchUsers();
        setIsDialogOpen(false);
        setEditingUser(null);
        toast.dismiss();
        toast.success("User updated successfully");
      }
    } catch (error) {
      console.error("Error updating user:", error);
    //   setError("Failed to update user");
      toast.dismiss();
      toast.error(error instanceof Error ? error.message : "Failed to update user");
    } finally {
      setIsSubmitting(false);
    }
  };

  const fetchApiRoles = async () => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/api-roles`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch API roles");
      }

      const data = await response.json();
      setApiRoles(data);
      console.log(data)
    } catch (error) {
      console.error("Error fetching API roles:", error);
      toast.dismiss();
      toast.error("Failed to fetch API roles");
    }
  };

  const handleAddApiRole = async (roleData: CreateAPIRoleType | UpdateAPIRoleType) => {
    try {
      setIsApiRoleSubmitting(true);
      const createData = roleData as CreateAPIRoleType;
      
      toast.loading("Creating API role...");
      
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/api-roles`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(createData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        toast.dismiss();
        toast.error(errorData.error || "Failed to create API role");
      } else {
        await fetchApiRoles();
        setIsApiRoleDialogOpen(false);
        toast.dismiss();
        toast.success("API role created successfully");
      }
    } catch (error) {
      console.error("Error creating API role:", error);
      toast.dismiss();
      toast.error("Failed to create API role");
    } finally {
      setIsApiRoleSubmitting(false);
    }
  };

  const handleEditApiRole = async (roleData: CreateAPIRoleType | UpdateAPIRoleType) => {
    try {
      setIsApiRoleSubmitting(true);
      const updateData = roleData as UpdateAPIRoleType;
      
      toast.loading("Saving changes...");
      
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/api-roles/${editingApiRole?.ID}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(updateData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        toast.dismiss();
        toast.error(errorData.error || "Failed to update API role");
      } else {
        await fetchApiRoles();
        setIsApiRoleDialogOpen(false);
        setEditingApiRole(null);
        toast.dismiss();
        toast.success("API role updated successfully");
      }
    } catch (error) {
      console.error("Error updating API role:", error);
      toast.dismiss();
      toast.error(error instanceof Error ? error.message : "Failed to update API role");
    } finally {
      setIsApiRoleSubmitting(false);
    }
  };

  const handleDeleteApiRole = async (roleId: string) => {
    try {
      toast.loading("Deleting API role...");
      
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/api-roles/${roleId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to delete API role");
      }

      await fetchApiRoles();
      toast.dismiss();
      toast.success("API role deleted successfully");
    } catch (error) {
      console.error("Error deleting API role:", error);
      toast.dismiss();
      toast.error(error instanceof Error ? error.message : "Failed to delete API role");
    }
  };

  const handleRegenerateApiKey = async (roleId: string) => {
    try {
      toast.loading("Regenerating API key...");
      
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/api-roles/${roleId}/regenerate-key`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to regenerate API key");
      }

      await fetchApiRoles();
      toast.dismiss();
      toast.success("API key regenerated successfully");
    } catch (error) {
      console.error("Error regenerating API key:", error);
      toast.dismiss();
      toast.error(error instanceof Error ? error.message : "Failed to regenerate API key");
    }
  };

  const handleDeleteUser = async (userId: string) => {
    try {
      toast.loading("Deleting user...");
      
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/users/${userId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to delete user");
      }

      await fetchUsers();
    //   setDeletingUserId(null);
      toast.dismiss();
      toast.success("User deleted successfully");
    } catch (error) {
      console.error("Error deleting user:", error);
      toast.dismiss();
      toast.error(error instanceof Error ? error.message : "Failed to delete user");
    }
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
      component: (
        <div className="w-full">
          <div className="flex justify-between items-center mb-4">
            <div className="flex flex-col">
              <h3 className="text-base font-semibold text-black dark:text-white">
                List of all users with access to the admin panel
              </h3>
              <p className="text-sm mt-2 text-gray-600 dark:text-gray-400">
                You can add or remove users from this list, or grant them access to the admin panel.
              </p>
            </div>
            <Button
              onClick={() => {
                setEditingUser(null);
                setIsDialogOpen(true);
              }}
            >
              Add User
            </Button>
          </div>

          <Table>
            <TableHeader className="text-black dark:text-white">
              <TableRow className="bg-gray-50 dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800">
                <TableHead className="w-1/5 font-semibold">Username</TableHead>
                <TableHead className="w-1/4 font-semibold">Email</TableHead>
                <TableHead className="w-1/6 font-semibold">Role</TableHead>
                <TableHead className="w-1/6 font-semibold">Status</TableHead>
                <TableHead className="w-1/6 text-right font-semibold">Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="text-black dark:text-white">
              {users.map((user) => (
                <TableRow key={user.id} className="hover:bg-gray-50 dark:hover:bg-gray-900/30">
                  <TableCell className="font-medium text-sm">{user.username}</TableCell>
                  <TableCell className="text-gray-600 dark:text-gray-400 text-sm">{user.email}</TableCell>
                  <TableCell className="text-gray-600 dark:text-gray-400 text-sm">{user.role.Name || 'User'}</TableCell>
                  <TableCell>
                    <span
                      className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                        user.active
                          ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                          : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
                      }`}
                    >
                      {user.active ? "Active" : "Inactive"}
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end space-x-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 px-2 text-sm text-blue-600 dark:text-blue-400"
                        onClick={() => {
                          setEditingUser(user);
                          setIsDialogOpen(true);
                        }}
                      >
                        Edit
                      </Button>

                      <Popover>
                        <PopoverTrigger asChild>
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-8 px-2 text-sm text-red-600 dark:text-red-400"
                            disabled={isSuperAdmin(user)}
                          >
                            Delete
                          </Button>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-4">
                          <div className="space-y-4">
                            <h4 className="font-medium">Confirm Deletion</h4>
                            <p className="text-sm text-gray-500 dark:text-gray-400">
                              {isSuperAdmin(user) ? (
                                "Super admin user cannot be deleted."
                              ) : (
                                <>
                                  Are you sure you want to delete user <span className="font-semibold">{user.username}</span>?
                                  This action cannot be undone.
                                </>
                              )}
                            </p>
                            <div className="flex justify-end space-x-2">
                              <Button
                                variant="destructive"
                                size="sm"
                                onClick={() => handleDeleteUser(user.id)}
                                disabled={isSuperAdmin(user)}
                              >
                                Delete
                              </Button>
                            </div>
                          </div>
                          
                        </PopoverContent>
                      </Popover>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          <div className="h-20 w-full"></div>

          <UserSheet
            isOpen={isDialogOpen}
            onClose={() => {
              setIsDialogOpen(false);
              setEditingUser(null);
            }}
            onSubmit={editingUser ? handleEditUser : handleAddUser}
            initialData={editingUser}
            isSubmitting={isSubmitting}
          />
        </div>
      )
    },
    {
      title: "API Role Management",
      icon: <IoKeyOutline />,
      component: (
        <div className="w-full">
          <div className="flex justify-between items-center mb-4">
            <div className="flex flex-col">
              <h3 className="text-base font-semibold text-black dark:text-white">
                List of all API keys with access to the API
              </h3>
              <p className="text-sm mt-2 text-gray-600 dark:text-gray-400">
                Generate API keys for your users to access the API.
              </p>
            </div>
            <Button
              onClick={() => {
                setEditingApiRole(null);
                setIsApiRoleDialogOpen(true);
              }}
            >
              Add API Role
            </Button>
          </div>

          <Table>
            <TableHeader className="text-black dark:text-white">
              <TableRow className="bg-gray-50 dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800">
                <TableHead className="w-1/5 font-semibold">Name</TableHead>
                <TableHead className="w-1/2 font-semibold">Description</TableHead>
                <TableHead className="w-1/6 font-semibold">Type</TableHead>
                <TableHead className="w-1/6 text-right font-semibold">Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="text-black dark:text-white">
              {apiRoles.map((role) => (
                <TableRow key={role.ID} className="hover:bg-gray-50 dark:hover:bg-gray-900/30">
                  <TableCell className="font-medium text-sm">{role.Name}</TableCell>
                  <TableCell className="text-gray-600 dark:text-gray-400 text-sm">
                    {role.Description}
                  </TableCell>
                  <TableCell className="text-gray-600 dark:text-gray-400 text-sm">
                    <span
                      className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                        role.IsSystem
                          ? 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400'
                          : 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
                      }`}
                    >
                      {role.IsSystem ? "System" : "Custom"}
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 text-sm text-blue-600 dark:text-blue-400"
                        onClick={() => {
                          setEditingApiRole(role);
                          setIsApiRoleDialogOpen(true);
                        }}
                        disabled={role.Type == "public_user"}
                      >
                        View Details
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          <div className="mt-4 mb-2 text-sm text-gray-400 dark:text-gray-400">
          The initial permissions for the aforementioned API Roles are empty. Please set the specific permissions for each API Role in the Schema Builder.
          </div>

          <div className="h-20 w-full"></div>

          <APIRoleSheet
            isOpen={isApiRoleDialogOpen}
            onClose={() => {
              setIsApiRoleDialogOpen(false);
              setEditingApiRole(null);
            }}
            onSubmit={editingApiRole ? handleEditApiRole : handleAddApiRole}
            initialData={editingApiRole}
            isSubmitting={isApiRoleSubmitting}
            onDelete={handleDeleteApiRole}
            onRegenerateToken={handleRegenerateApiKey}
          />
        </div>
      )
    }
  ];

  return (
    <div className={`${geistSans.variable} ${geistMono.variable} w-full h-screen bg-white dark:bg-black flex flex-col overflow-hidden`}>
      <Header username={user?.username} onLogout={logout} />
      <aside className="w-full flex-1 grid grid-cols-6">
        <Sidebar />
        <div className="col-span-5 flex flex-col justify-start items-center overflow-hidden">
          <div className="w-full flex items-center gap-4 border-b-[1px] border-gray-200 border-dotted px-5 py-3">
            <h2 className="text-base w-max font-semibold text-black dark:text-white">
              Access Control
            </h2>

            <HeaderNav
              items={items}
              onChange={(component: React.ReactNode, index: number) => {
                setSelectedIndex(index);
              }}
            />
          </div>

          <div className="max-w-4xl w-full flex flex-col h-[calc(100vh-4rem)] overflow-auto pt-5 px-5 border-l-[1px] border-r-[1px] border-gray-200 border-dotted">
            {items[selectedIndex].component}
          </div>
        </div>
      </aside>
    </div>
  );
}