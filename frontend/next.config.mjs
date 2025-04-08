/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "export", // Generate static files for deployment
  images: {
    unoptimized: true, // Required for static export
  },
  trailingSlash: true, // Add trailing slashes to all URLs

  // Development-only configuration
  ...(process.env.NODE_ENV === "development" && {
    async rewrites() {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081";

      return [
        {
          source: "/health",
          destination: `${apiUrl}/health`,
        },
        {
          source: "/process/:path*",
          destination: `${apiUrl}/process/:path*`,
        },
      ];
    },
  }),
};

export default nextConfig;
