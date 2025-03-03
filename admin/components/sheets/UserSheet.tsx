import { Button } from "@/components/ui/button";
import React, { useEffect, useState } from "react";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select";
import { UserType, CreateUserType, UpdateUserType } from "@/types/user";
import { RoleType } from "@/types/role";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

type UserSheetProps = {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (userData: CreateUserType | UpdateUserType) => Promise<void>;
  initialData?: UserType | null;
  isSubmitting?: boolean;
};

export default function UserSheet({ 
  isOpen, 
  onClose, 
  onSubmit, 
  initialData,
  isSubmitting = false 
}: UserSheetProps) {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [active, setActive] = useState(true);
  const [roleId, setRoleId] = useState("");
  const [password, setPassword] = useState("");
  const [roles, setRoles] = useState<RoleType[]>([]);
  const [loading, setLoading] = useState(false);
  const [currentUser, setCurrentUser] = useState<UserType>();

  // Get current user from localStorage
  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      setCurrentUser(JSON.parse(storedUser));
    }
  }, []);

  // Fetch role list
  useEffect(() => {
    const fetchRoles = async () => {
      try {
        setLoading(true);
        const token = localStorage.getItem("token");
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/admin/roles`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch role list");
        }

        const data = await response.json();
        setRoles(data);
      } catch (error) {
        console.error("Error fetching roles:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchRoles();
  }, []);

  // Set initial data when editing
  useEffect(() => {
    if (initialData) {
      setUsername(initialData.username);
      setEmail(initialData.email);
      setActive(initialData.active);
      setRoleId(initialData.role_id || initialData.role.ID);
      setPassword(""); // Don't show password when editing
    } else {
      // Reset form for new user
      setUsername("");
      setEmail("");
      setActive(true);
      setRoleId(roles.length > 0 ? roles[0].ID : "");
      setPassword("");
    }
  }, [initialData, roles]);

  const handleSubmit = async () => {
    // Validate required fields
    if (!username || !email || (!initialData && !password) || !roleId) {
      toast.error("Please fill in all required fields");
      return;
    }

    // Check super admin operations
    if (initialData?.role?.Type === 'super_admin') {
      if (currentUser?.id !== initialData.id) {
        toast.error("Cannot modify other super admin users");
        return;
      }
    }

    try {
      if (initialData) {
        // Update existing user
        const updateData: UpdateUserType = {
          username,
          email,
          role_id: roleId,
          active,
        };
        
        // Only include password if it has been changed
        if (password.trim()) {
          updateData.password = password;
        }
        
        await onSubmit(updateData);
      } else {
        // Create new user
        const createData: CreateUserType = {
          username,
          email,
          password,
          role_id: roleId,
          active,
        };
        
        await onSubmit(createData);
      }
    } catch (error) {
      console.error("Error submitting user data:", error);
      toast.error("Failed to save user data");
    }
  };

  const isSuperAdmin = initialData?.role?.Type === 'super_admin';
  const isCurrentUser = currentUser?.id === initialData?.id;
  // const canModifyRole = !isSuperAdmin || (isSuperAdmin && isCurrentUser);

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent className="bg-white dark:bg-black text-black dark:text-white">
        <SheetHeader>
          <SheetTitle className="text-black dark:text-white">
            {initialData ? "Edit User" : "Add New User"}
          </SheetTitle>
          <SheetDescription className="text-gray-600 dark:text-gray-400">
            {initialData
              ? "Make changes to the user account here."
              : "Add a new user to the admin panel here."}
          </SheetDescription>
        </SheetHeader>

        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="username" className="text-black dark:text-white">
              Username <span className="text-red-500">*</span>
            </Label>
            <Input
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter username"
              className="bg-transparent border-gray-200 dark:border-gray-800"
              disabled={isSubmitting}
            />
          </div>
          
          <div className="grid gap-2">
            <Label htmlFor="email" className="text-black dark:text-white">
              Email <span className="text-red-500">*</span>
            </Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter email"
              className="bg-transparent border-gray-200 dark:border-gray-800"
              disabled={isSubmitting}
            />
          </div>
          
          <div className="grid gap-2">
            <Label htmlFor="password" className="text-black dark:text-white">
              {initialData ? "Password (leave empty to keep current)" : "Password *"}
            </Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder={initialData ? "Leave empty to keep current password" : "Enter password"}
              className="bg-transparent border-gray-200 dark:border-gray-800"
              disabled={isSubmitting}
            />
          </div>
          
          <div className="grid gap-2">
            <Label htmlFor="role" className="text-black dark:text-white">
              Role <span className="text-red-500">*</span>
            </Label>
            <Select 
              value={roleId} 
              onValueChange={setRoleId}
              disabled={loading || isSubmitting || isSuperAdmin}
            >
              <SelectTrigger className="bg-transparent border-gray-200 dark:border-gray-800 text-black dark:text-white">
                <SelectValue placeholder="Select a role" />
              </SelectTrigger>
              <SelectContent className="bg-white dark:bg-gray-900">
                {roles.map((role) => (
                  <SelectItem 
                    key={role.ID} 
                    value={role.ID} 
                    className="text-black dark:text-white"
                  >
                    {role.Name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {isSuperAdmin && !isCurrentUser && (
              <p className="text-sm text-yellow-600 dark:text-yellow-400 mt-1">
                Cannot modify super admin role
              </p>
            )}
          </div>

          <div className="flex items-center justify-between">
            <Label htmlFor="active" className="text-black dark:text-white">Active Status</Label>
            <div className="flex items-center space-x-2">
              <Switch
                id="active"
                checked={active}
                onCheckedChange={setActive}
                disabled={isSubmitting || isSuperAdmin}
              />
            </div>
          </div>

          <div className="flex justify-end gap-4 pt-4">
            <Button variant="outline" onClick={onClose} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button onClick={handleSubmit} disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  {initialData ? "Saving..." : "Creating..."}
                </>
              ) : (
                initialData ? "Save Changes" : "Create User"
              )}
            </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}