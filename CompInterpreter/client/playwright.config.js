import { defineConfig, devices } from '@playwright/test';

/* E2E de CompInterpreter: levanta el servidor Express (Jison + intérprete)
   y el cliente React (Vite) y corre las pruebas contra ambos juntos. */
export default defineConfig({
  testDir: './e2e',
  timeout: 30000,
  fullyParallel: false,
  reporter: 'list',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'retain-on-failure',
  },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  webServer: [
    {
      command: 'node server.js',
      cwd: '../server',
      url: 'http://localhost:4000/salud',
      timeout: 20000,
      reuseExistingServer: true,
    },
    {
      command: 'npm run dev -- --port 5173 --strictPort',
      cwd: '.',
      url: 'http://localhost:5173',
      timeout: 20000,
      reuseExistingServer: true,
    },
  ],
});
