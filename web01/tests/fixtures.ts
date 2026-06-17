import type { Page, Route } from '@playwright/test';

export type Channel = 'facebook' | 'line' | 'instagram';
export type Status = 'pending' | 'replied';

export interface InboxMessage {
  id: string;
  customer_id: string;
  sender_name: string;
  avatar_url: string;
  channel: Channel;
  preview: string;
  status: Status;
  unread: boolean;
  created_at: string;
  updated_at: string;
}

const AVATAR = 'https://example.com/avatar.jpg';

function msg(
  id: number,
  sender_name: string,
  channel: Channel,
  status: Status,
  preview: string
): InboxMessage {
  const ts = new Date(2026, 5, 17, 10, 30, 0).toISOString();
  return {
    id: String(id),
    customer_id: `cust_${id}`,
    sender_name,
    avatar_url: AVATAR,
    channel,
    preview,
    status,
    unread: status === 'pending',
    created_at: ts,
    updated_at: ts
  };
}

/**
 * Curated + generated dataset. The first four entries mirror the original
 * dashboard mock so layout tests keep working; the rest pad the set to 25 so
 * pagination kicks in at page_size 20.
 */
export function buildMessages(): InboxMessage[] {
  const curated: InboxMessage[] = [
    msg(1, 'Marcus Watanabe', 'line', 'pending', 'Can you confirm the tracking number for order #8812?'),
    msg(2, 'Sarah Jenkins', 'facebook', 'replied', 'Thank you for the quick resolution!'),
    msg(3, 'David Wilson', 'instagram', 'pending', 'Do you have the new winter collection in stock yet?'),
    msg(4, 'Li Na', 'line', 'replied', "I'd like to update my shipping address for the recent order.")
  ];

  const channels: Channel[] = ['facebook', 'line', 'instagram'];
  const generated: InboxMessage[] = [];
  for (let i = 5; i <= 25; i++) {
    const channel = channels[i % channels.length]!;
    const status: Status = i % 2 === 0 ? 'replied' : 'pending';
    generated.push(msg(i, `Customer ${i}`, channel, status, `Generic message number ${i}`));
  }

  return [...curated, ...generated];
}

interface MockOptions {
  /** Force every request to fail with this HTTP status (e.g. 401). */
  forceStatus?: number;
}

/**
 * Installs a route handler for the inbox API that filters/paginates an
 * in-memory dataset and shapes the response exactly like the Go backend.
 */
export async function mockInboxApi(page: Page, options: MockOptions = {}): Promise<void> {
  const all = buildMessages();

  await page.route('**/api/v1/inbox/messages*', async (route: Route) => {
    const request = route.request();
    const auth = request.headers()['authorization'];

    if (options.forceStatus === 401 || !auth || !auth.startsWith('Bearer ')) {
      await route.fulfill({
        status: 401,
        contentType: 'application/json',
        body: JSON.stringify({
          success: false,
          error: { code: 'UNAUTHORIZED', message: 'Missing or invalid JWT' }
        })
      });
      return;
    }

    const url = new URL(request.url());
    const channel = url.searchParams.get('channel') ?? '';
    const status = url.searchParams.get('status') ?? '';
    const q = (url.searchParams.get('q') ?? '').toLowerCase();
    const page_ = Number(url.searchParams.get('page') ?? '1');
    const pageSize = Number(url.searchParams.get('page_size') ?? '20');

    let filtered = all;
    if (channel) filtered = filtered.filter(m => m.channel === channel);
    if (status) filtered = filtered.filter(m => m.status === status);
    if (q) {
      filtered = filtered.filter(
        m => m.sender_name.toLowerCase().includes(q) || m.preview.toLowerCase().includes(q)
      );
    }

    const totalItems = filtered.length;
    const totalPages = Math.ceil(totalItems / pageSize);
    const start = (page_ - 1) * pageSize;
    const data = filtered.slice(start, start + pageSize);

    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        success: true,
        data,
        pagination: {
          page: page_,
          page_size: pageSize,
          total_items: totalItems,
          total_pages: totalPages
        }
      })
    });
  });
}
