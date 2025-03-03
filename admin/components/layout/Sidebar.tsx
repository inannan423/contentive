import { IoSettingsOutline, IoDocumentTextOutline, IoGridOutline, IoKeyOutline, IoLayersOutline } from "react-icons/io5";
import Link from "next/link";
import { useRouter } from "next/router";
import React from "react";
import { useAuth } from "@/contexts/AuthContext";

type NavItemProps = {
  icon: React.ReactNode;
  label: string;
  path: string;
  isActive: boolean;
};

const NavItem = ({ icon, label, path, isActive }: NavItemProps) => (
  <Link href={path} passHref className="flex flex-col justify-start items-start w-full gap-1">
    <div
      className={`flex ${isActive ? "bg-gray-100" : "hover:bg-gray-50"} h-8 justify-start items-center gap-3 w-full px-4 py-1 rounded-sm cursor-pointer`}
    >
      <div className="flex items-center justify-center w-5 text-black dark:text-white">
        {icon}
      </div>
      <p className="text-black font-semibold text-sm dark:text-white">
        {label}
      </p>
    </div>
  </Link>
);

export default function Sidebar() {
  const router = useRouter();
  const { user } = useAuth();
  const currentPath = router.pathname;
  
  const navItems = [
    {
      path: "/content",
      label: "Content Management",
      icon: <IoDocumentTextOutline size={16} />,
      requiredRole: ["super_admin", "content_admin", "editor", "viewer"],
    },
    {
      path: "/schema",
      label: "Schema Builder",
      icon: <IoLayersOutline size={16} />,
      requiredRole: ["super_admin", "content_admin"],
    },
    {
      path: "/media",
      label: "Media Library",
      icon: <IoGridOutline size={16} />,
      requiredRole: ["super_admin", "content_admin", "editor", "viewer"],
    },
    {
      path: "/access",
      label: "Access",
      icon: <IoKeyOutline size={16} />,
      requiredRole: ["super_admin"],
    },
    {
      path: "/settings",
      label: "Settings",
      icon: <IoSettingsOutline size={16} />,
      requiredRole: ["super_admin", "content_admin"],
    },
  ];

  // Filter nav items based on user role
  const authorizedNavItems = navItems.filter(item => 
    item.requiredRole.includes(user?.role || '')
  );

  console.log(authorizedNavItems);

  return (
    <div className="col-span-1 border-r border-gray-200 flex flex-col justify-start items-center px-3 py-3 w-full gap-2">
      {authorizedNavItems.map((item) => (
        <NavItem
          key={item.path}
          icon={item.icon}
          label={item.label}
          path={item.path}
          isActive={currentPath.startsWith(item.path)}
        />
      ))}
    </div>
  );
}