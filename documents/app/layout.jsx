import { Footer, Layout, Navbar } from 'nextra-theme-docs'
import { Banner, Head } from 'nextra/components'
import { getPageMap } from 'nextra/page-map'
import 'nextra-theme-docs/style.css'
import '../globals.css'

export const metadata = {
  // Define your metadata here
  // For more information on metadata API, see: https://nextjs.org/docs/app/building-your-application/optimizing/metadata
}
 
const banner = <Banner storageKey="some-key">Nextra 4.0 is released 🎉</Banner>
const navbar = (
  <Navbar
    logo={
      <div className='flex w-full items-center justify-center gap-2'>
          <img src="/contentive_logo_white.svg" alt="Contentive Logo" className='h-5 hidden dark:block' />
          <img src="/contentive_logo_black.svg" alt="Contentive Logo" className='h-5 dark:hidden' />
          <p className='font-bold x:text-black x:dark:text-white font-mono'>
            contentive_
          </p>
      </div>
    }
  />
)
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
      MIT {new Date().getFullYear()} © Contentive.
    </p>
  </div>
</Footer>
 
export default async function RootLayout({ children }) {
  return (
    <html
      // Not required, but good for SEO
      lang="en"
      // Required to be set
      dir="ltr"
      // Suggested by `next-themes` package https://github.com/pacocoursey/next-themes#with-app
      suppressHydrationWarning
    >
      <Head
      // ... Your additional head options
      >
        {/* Your additional tags should be passed as `children` of `<Head>` element */}
      </Head>
      <body>
        <Layout
          // banner={banner}
          navbar={navbar}
          pageMap={await getPageMap()}
          docsRepositoryBase="https://github.com/inannan423/contentive/tree/main/documents"
          footer={footer}
          // ... Your additional layout options
        >
          {children}
        </Layout>
      </body>
    </html>
  )
}