import { useState } from 'react';
import React from 'react';

export type HeaderNavItem = {
    title: string;
    icon: React.ReactNode;
    component: React.ReactNode;
}

export type HeaderNavProps = {
    items: HeaderNavItem[];
    onChange: (component: React.ReactNode, index: number) => void;
}

export function HeaderNav({
    items,
    onChange,
}: HeaderNavProps) {
    const [activeIndex, setActiveIndex] = useState(0);
    const handleChange = (component: React.ReactNode, index: number) => {
        setActiveIndex(index);
        onChange(component, index);
    }

    return (
        <div className="flex w-max justify-start items-center gap-6">
            {
                items.map((item, index) => (
                    <div
                        key={index}
                        className={`flex w-max justify-start items-center gap-2 cursor-pointer transition-all duration-300 ease-in-out px-2 py-1 ${activeIndex === index ? 'text-gray-800 dark:text-gray-400 bg-gray-100 dark:bg-gray-800 rounded-lg': 'text-gray-400 dark:text-gray-600 bg-transparent rounded-lg hover:text-gray-800 dark:hover:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'}`}
                        onClick={() => handleChange(item.component, index)}
                    >
                        <div className="flex w-max justify-start items-center gap-2">
                            <div className="flex w-max justify-start items-center gap-2">
                                {item.icon}
                            </div>
                            <div className="flex w-max justify-start items-center gap-2">
                                <p className="text-sm font-medium">{item.title}</p>
                            </div>
                        </div>
                    </div>
                ))
            }
        </div>
    )
}