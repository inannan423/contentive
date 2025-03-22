import { Layout, Navbar } from 'nextra-theme-docs'
import { Banner, Head } from 'nextra/components'
import { getPageMap } from 'nextra/page-map'
import 'nextra-theme-docs/style.css'
import '../globals.css'
import footer from '../components/footer'

export const metadata = {
  title: 'Contentive Documentation',
  description: 'Documentation for Contentive',
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

 
export default async function RootLayout({ children }) {
  return (
    <html
      lang="en"
      dir="ltr"
      suppressHydrationWarning
    >
      <Head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <body suppressHydrationWarning>
        <Layout
          // banner={banner}
          navbar={navbar}
          pageMap={await getPageMap()}
          docsRepositoryBase="https://github.com/inannan423/contentive/tree/main/documents"
          footer={footer}
        >
          {children}
        </Layout>
      </body>
    </html>
  )
}