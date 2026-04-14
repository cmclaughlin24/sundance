import { test, expect } from '@playwright/test';

test('page loads without errors', async ({ page }) => {
  await page.goto('/');

  const consoleErrors: string[] = [];
  page.on('console', (msg) => {
    if (msg.type() === 'error') {
      consoleErrors.push(msg.text());
    }
  });

  await page.waitForLoadState('networkidle');

  expect(consoleErrors).toHaveLength(0);
});
