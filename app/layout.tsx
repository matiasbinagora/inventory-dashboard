import type { Metadata } from 'next'
import './styles.css'

export const metadata: Metadata = {
  title: 'Inventory control room',
  description: 'Local editorial inventory dashboard',
}

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="en"><body>{children}</body></html>
}
