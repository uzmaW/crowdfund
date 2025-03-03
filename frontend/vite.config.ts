import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import path from 'path';
import { defineConfig, loadEnv, UserConfig } from 'vite';
import VueDevTools from 'vite-plugin-vue-devtools';
import tailwindcss from '@tailwindcss/vite'

// https://vitejs.dev/config/
export default defineConfig(({ mode }: UserConfig) => {

  // https://vitejs.dev/config/#environment-variables
  const env = loadEnv(mode ?? 'development', process.cwd(), '');

  const BASE_PATH: string = env.APP_BASE_PATH.endsWith('/')
    ? env.APP_BASE_PATH.slice(0, -1)
    : env.APP_BASE_PATH + '/';

  const proxy: Record<string, string> = {
    api: BASE_PATH + '/api',
  };

  return {
    base: BASE_PATH,
    plugins: [
      vue(),
      vueJsx(),
      VueDevTools(),
      tailwindcss(),
    ],
    define: {
      'process.env': env,
    },
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      host: true,
      proxy: {
        [proxy.api]: {
          target: env.BASE_API_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(BASE_PATH, ''),
        },
      },
    },
  };
});
