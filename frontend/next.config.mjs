/** @type {import('next').NextConfig} */
const nextConfig = {
    // output: 'export', // Removed static export to support dynamic routes
    // Optional: Add other configurations here if needed
    images: {
        unoptimized: true,
    },
    // Recommended: Disable trailing slash for consistency with API routes
    trailingSlash: false,
};

export default nextConfig; 