import React from 'react'

import { Footer, Layout, Navbar } from 'nextra-theme-docs'
import { useRouter } from 'next/router'

const footer = <Footer>
  <div className='flex flex-col w-full items-center justify-center gap-2'>
    <div className='flex w-full my-5 items-center justify-center gap-2'>
      <img src="/contentive_logo_white.svg" alt="Contentive Logo" className='h-10 hidden dark:block' />
      <img src="/contentive_logo_black.svg" alt="Contentive Logo" className='h-10 dark:hidden' />
      <p className='font-bold text-3xl text-black dark:text-white font-mono'>
        contentive_
      </p>
    </div>
    <p className='text-xs dark:text-white'>
      MIT © Contentive.
    </p>
  </div>
</Footer>

export default footer