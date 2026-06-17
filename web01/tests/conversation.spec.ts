import { test, expect } from '@playwright/test';
import { mockInboxApi } from './fixtures';

test.describe('Thread navigation', () => {
  test('each thread row links to its conversation route', async ({ page }) => {
    await mockInboxApi(page);
    await page.goto('/');

    await expect(page.getByTestId('thread-link-2')).toHaveAttribute(
      'href',
      '/conversations/2'
    );
  });

  test('clicking a thread navigates to the conversation view', async ({ page }) => {
    await mockInboxApi(page);
    await page.goto('/');
    await expect(page.getByTestId('thread-link-3')).toBeVisible();

    await page.getByTestId('thread-link-3').click();

    await expect(page).toHaveURL(/\/conversations\/3$/);

    // The inbox list persists in the shared layout and highlights the open thread.
    await expect(page.getByTestId('message-item-3')).toBeVisible();
    await expect(page.getByTestId('message-item-3')).toHaveClass(/message-item-active/);
    await expect(page.getByTestId('message-item-1')).not.toHaveClass(/message-item-active/);
  });
});
