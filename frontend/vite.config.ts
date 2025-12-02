import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react-swc'
import tailwindcss from '@tailwindcss/vite'


// https://vite.dev/config/
export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), '')
    return {
        // vite config
        plugins: [
            {
                name: 'check-if-env-variables-are-set',
                config() {
                    if (!env.VITE_API_BASE_URL) {
                        throw new Error("VITE_API_BASE_URL not set")
                    }
                }
            },
            react(),
            tailwindcss()
        ],
        server: {
            port: 3000,
            host: true
        }
    }
})