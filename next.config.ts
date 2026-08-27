import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      { source: '/api/:path*', destination: 'http://127.0.0.1:8080/api/:path*' },
      { source: '/media/:path*', destination: 'http://127.0.0.1:8080/media/:path*' },
    ]
  },
}

export default nextConfig
