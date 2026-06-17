import { test, expect } from '@playwright/test';
import { mockInboxApi } from './fixtures';

test.beforeEach(async ({ page }) => {
  await mockInboxApi(page);
  await page.goto('/');
});

test('has the dashboard page title', async ({ page }) => {
  await expect(page).toHaveTitle('OmniChat Inbox Dashboard');
});

test('renders the sidebar brand and navigation', async ({ page }) => {
  await expect(page.getByRole('heading', { name: 'OmniChat', exact: true })).toBeVisible();
  await expect(page.getByText('Aggregator')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Compose' })).toBeVisible();

  for (const item of ['Inbox', 'Sent', 'Archived', 'Settings']) {
    await expect(page.getByRole('link', { name: item })).toBeVisible();
  }

  await expect(page.getByText('Alex Chen')).toBeVisible();
});

test('renders the top navigation bar', async ({ page }) => {
  await expect(page.getByRole('heading', { name: 'OmniChat Dashboard' })).toBeVisible();
  await expect(page.getByPlaceholder('Search across all channels...')).toBeVisible();
});

test('renders the inbox list with channel filters', async ({ page }) => {
  await expect(page.getByRole('heading', { name: /Inbox/ })).toBeVisible();
  await expect(page.getByText('(25)')).toBeVisible();

  for (const channel of ['All', 'Facebook', 'LINE', 'Instagram']) {
    await expect(page.getByRole('button', { name: channel, exact: true })).toBeVisible();
  }
});

test('renders all conversation items', async ({ page }) => {
  const names = ['Marcus Watanabe', 'Sarah Jenkins', 'David Wilson', 'Li Na'];
  for (const name of names) {
    await expect(page.getByTestId(/message-item-/).filter({ hasText: name })).toBeVisible();
  }

  await expect(
    page.getByTestId('message-item-1').getByText('Can you confirm the tracking number for order #8812?')
  ).toBeVisible();
  await expect(page.getByTestId('message-item-1').getByText('Pending')).toBeVisible();
  await expect(page.getByTestId('message-item-2').getByText('Replied')).toBeVisible();
});

test('first conversation is active by default', async ({ page }) => {
  const firstItem = page.getByTestId('message-item-1');
  await expect(firstItem).toHaveClass(/message-item-active/);
});

test('clicking a conversation makes it active', async ({ page }) => {
  const target = page.getByTestId('message-item-3');
  await expect(target).not.toHaveClass(/message-item-active/);

  await target.click();

  await expect(target).toHaveClass(/message-item-active/);
  await expect(page.getByTestId('message-item-1')).not.toHaveClass(/message-item-active/);
});

test('selecting a channel filter highlights it', async ({ page }) => {
  const lineFilter = page.getByRole('button', { name: 'LINE', exact: true });
  await lineFilter.click();
  await expect(lineFilter).toHaveClass(/bg-primary/);
});

test('renders the chat panel', async ({ page }) => {
  const chatHeader = page.locator('section').filter({ hasText: 'Online via LINE' });
  await expect(chatHeader.getByRole('heading', { name: 'Marcus Watanabe' })).toBeVisible();
  await expect(page.getByText('Online via LINE')).toBeVisible();
  await expect(page.getByPlaceholder('Type your message here...')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Send' })).toBeVisible();
});

test('renders the contextual detail pane', async ({ page }) => {
  await expect(page.getByText('Tokyo, Japan (GMT+9)')).toBeVisible();
  await expect(page.getByText('Orders')).toBeVisible();
  await expect(page.getByText('LTV')).toBeVisible();
  await expect(page.getByText('$1.4k')).toBeVisible();

  for (const tag of ['VIP_CUSTOMER', 'TECH_EARLY_ADOPTER', 'URGENT_TICKET']) {
    await expect(page.getByText(tag)).toBeVisible();
  }

  await expect(page.getByRole('button', { name: 'View Full Profile' })).toBeVisible();
});
