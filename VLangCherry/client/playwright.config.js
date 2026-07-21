import { defineConfig, devices } from '@playwright/test';

/* E2E de VLangCherry: levanta el servidor Go (ANTLR + interprete) y el
   cliente React (Vite) y corre las pruebas contra ambos juntos. */
export default defineConfig({
  testDir: './e2e',
  timeout: 30000,
  fullyParallel: false,
  reporter: 'list',
  use: {
    baseURL: 'http://localhost:5174',
    trace: 'retain-on-failure',
  },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  webServer: [
    {
      command: 'go run ./cmd/servidor',
      cwd: '../server',
      env: { PORT: '4100' },
      url: 'http://localhost:4100/salud',
      timeout: 30000,
      reuseExistingServer: true,
    },
    {
      command: 'npm run dev -- --port 5174 --strictPort',
      cwd: '.',
      url: 'http://localhost:5174',
      timeout: 20000,
      reuseExistingServer: true,
    },
  ],
});
