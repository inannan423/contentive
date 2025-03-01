import Image from "next/image";
import { Button } from "@/components/ui/button";

type HeaderProps = {
  username: string;
  onLogout: () => void;
};

export default function Header({ username, onLogout }: HeaderProps) {
  return (
    <div className="header w-full border-b-[1px] border-gray-200 h-max py-2 px-6">
      <div className='flex w-full items-center justify-between'>
        <div className='flex items-center gap-2'>
          <Image src="/contentive_logo_white.svg" alt="Contentive Logo" className='h-5 hidden dark:block' width={20} height={20} />
          <Image src="/contentive_logo_black.svg" alt="Contentive Logo" className='h-5 dark:hidden' width={20} height={20} />
          <p className='font-bold text-black dark:text-white font-mono'>
            contentive_
          </p>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Welcome, <span className="font-medium text-black dark:text-white">{username}</span>
            </p>
            <Button variant="outline" className="text-black dark:text-white" size="sm" onClick={onLogout}>
              Logout
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}