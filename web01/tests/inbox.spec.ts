import { test, expect } from '@playwright/test';
import { mockInboxApi } from './fixtures';

test.describe('Inbox message list', () => {
  test('loads inbox messages on visiting /', async ({ page }) => {
    await mockInboxApi(page);
    await page.goto('/');

    await expect(page.getByTestId('message-item-1')).toBeVisible();
    await expect(
      page.getByTestId('message-item-1').getByText('Can you confirm the tracking number for order #8812?')
    ).toBeVisible();
    await expect(page.getByText('(25)')).toBeVisible();
  });

  test('filters by channel', async ({ page }) => {
    await mockInboxApi(page);
    await page.goto('/');
    await expect(page.getByTestId('message-item-2')).toBeVisible(); // Sarah / facebook

    await page.getByTestId('filter-channel-line').click();

    // Marcus (line) stays, Sarah (facebook) is gone.
    await expect(page.getByTestId('message-item-1')).toBeVisible();
    await expect(page.getByTestId('message-item-2')).toHaveCount(0);
    await expect(page.getByTestId('filter-channel-line')).toHaveClass(/bg-primary/);
  });

  test('filters by status', async ({ page }) => {
    await mockInboxApi(page);
    await page.goto('/');
    await expect(page.getByTestId('message-item-1')).toBeVisible(); // pending
    await expect(page.getByTestId('message-item-2')).toBeVisible(); // replied

    await page.getByTestId('filter-status').selectOption('pending');

    await expect(page.getByTestId('message-item-1')).toBeVisible(); // pending stays
    await expect(page.getByTestId('message-item-2')).toHaveCount(0); // replied gone
  });

  test('combined channel filter + search', async ({ page }) => {
    await mockInboxApi(page);
    await page.goto('/');
    await expect(page.getByTestId('message-item-1')).toBeVisible();

    await page.getByTestId('filter-channel-line').click();
    await page.getByTestId('search-input').fill('tracking');

    // Only Marcus matches line + "tracking".
    await expect(page.getByTestId('message-item-1')).toBeVisible();
    await expect(page.getByTestId('message-item-4')).toHaveCount(0); // Li Na (line, no tracking)
    await expect(page.getByTestId(/message-item-/)).toHaveCount(1);
  });

  test('paginates through results', async ({ page }) => {
    await mockInboxApi(page);
    await page.goto('/');
    await expect(page.getByTestId('message-item-1')).toBeVisible();

    // 25 items / 20 per page => 2 pages.
    await expect(page.getByTestId('pagination-controls')).toBeVisible();

    await page.getByTestId('pagination-page-2').click();

    // Page 2 holds items 21-25; page-1 items are gone.
    await expect(page.getByTestId('message-item-21')).toBeVisible();
    await expect(page.getByTestId('message-item-1')).toHaveCount(0);
    await expect(page.getByTestId('pagination-page-2')).toHaveClass(/bg-primary/);
  });

  test('shows empty state when nothing matches', async ({ page }) => {
    await mockInboxApi(page);
    await page.goto('/');
    await expect(page.getByTestId('message-item-1')).toBeVisible();

    await page.getByTestId('search-input').fill('no-such-message-zzz');

    await expect(page.getByTestId('inbox-empty')).toBeVisible();
    await expect(page.getByText('No messages found')).toBeVisible();
    await expect(page.getByTestId(/message-item-/)).toHaveCount(0);
  });

  test('shows an error state on 401 unauthorized', async ({ page }) => {
    await mockInboxApi(page, { forceStatus: 401 });
    await page.goto('/');

    await expect(page.getByTestId('inbox-error')).toBeVisible();
    await expect(page.getByTestId(/message-item-/)).toHaveCount(0);
  });
});
