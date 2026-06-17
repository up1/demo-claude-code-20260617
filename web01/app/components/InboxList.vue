<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useInboxStore, type Channel, type InboxMessage } from '~/stores/inbox'

const store = useInboxStore()
const { messages, channel, status, page, totalItems, totalPages, loading, error, isEmpty } = storeToRefs(store)

const channels: { label: string, value: '' | Channel, color: string }[] = [
  { label: 'All', value: '', color: '' },
  { label: 'Facebook', value: 'facebook', color: '#1877F2' },
  { label: 'LINE', value: 'line', color: '#06C755' },
  { label: 'Instagram', value: 'instagram', color: '#E4405F' }
]

const channelColors: Record<Channel, string> = {
  facebook: '#1877F2',
  line: '#06C755',
  instagram: '#E4405F'
}

const route = useRoute()
const searchInput = ref('')
let searchTimer: ReturnType<typeof setTimeout> | null = null

// The active thread is driven by the route (/conversations/{id}); on the
// dashboard (`/`) the first conversation is highlighted by default.
const activeId = computed<string | null>(() => {
  const param = route.params.thread_id
  if (typeof param === 'string' && param) return param
  return messages.value[0]?.id ?? null
})

// Fetch on the client so Playwright can intercept the request.
onMounted(() => {
  store.fetchMessages()
})

function onSearchInput() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    store.setSearch(searchInput.value.trim())
  }, 250)
}

function formatTime(message: InboxMessage): string {
  const date = new Date(message.updated_at)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
}
</script>

<template>
  <section class="w-chat-list-width flex flex-col border-r border-outline-variant bg-surface-container-lowest">
    <div class="p-lg space-y-md">
      <div class="flex items-center justify-between">
        <h3 class="font-display-lg text-display-lg">
          Inbox <span class="text-on-surface-variant font-normal text-headline-md ml-xs">({{ totalItems }})</span>
        </h3>
        <button class="p-unit rounded hover:bg-surface-container transition-colors">
          <span class="material-symbols-outlined text-on-surface-variant">filter_list</span>
        </button>
      </div>

      <!-- Channel Filters -->
      <div data-testid="filter-channel" class="flex gap-sm overflow-x-auto scrollbar-hide pb-unit">
        <button
          v-for="ch in channels"
          :key="ch.label"
          :data-testid="`filter-channel-${ch.value || 'all'}`"
          :class="channel === ch.value
            ? 'px-md py-1.5 rounded-full bg-primary text-on-primary font-label-caps text-label-caps flex items-center gap-xs'
            : 'px-md py-1.5 rounded-full bg-surface-container text-on-surface-variant hover:bg-surface-container-high transition-colors font-label-caps text-label-caps flex items-center gap-xs'"
          @click="store.setChannel(ch.value)"
        >
          <span
            v-if="ch.color"
            class="w-2 h-2 rounded-full"
            :style="{ backgroundColor: ch.color }"
          />
          {{ ch.label }}
        </button>
      </div>

      <!-- Status filter & Search -->
      <div class="flex gap-sm">
        <div class="flex-1 relative">
          <span class="absolute left-2.5 top-1/2 -translate-y-1/2 material-symbols-outlined text-[18px] text-outline">tune</span>
          <select
            data-testid="filter-status"
            class="w-full bg-surface-container-low border-none rounded-lg pl-9 pr-md py-2 text-body-sm appearance-none focus:ring-1 focus:ring-primary"
            :value="status"
            @change="store.setStatus(($event.target as HTMLSelectElement).value as '' | 'pending' | 'replied')"
          >
            <option value="">All statuses</option>
            <option value="pending">Pending</option>
            <option value="replied">Replied</option>
          </select>
        </div>
      </div>
      <div class="relative">
        <span class="absolute left-2.5 top-1/2 -translate-y-1/2 material-symbols-outlined text-[18px] text-outline">search</span>
        <input
          v-model="searchInput"
          data-testid="search-input"
          type="text"
          placeholder="Search sender or message..."
          class="w-full bg-surface-container-low border-none rounded-lg pl-9 pr-md py-2 text-body-sm focus:ring-1 focus:ring-primary"
          @input="onSearchInput"
        >
      </div>
    </div>

    <!-- Scrollable Message List -->
    <div class="flex-1 overflow-y-auto scrollbar-hide px-sm space-y-xs pb-lg">
      <!-- Loading -->
      <div v-if="loading" data-testid="inbox-loading" class="p-lg text-center text-on-surface-variant font-body-sm">
        Loading messages...
      </div>

      <!-- Unauthorized / Error -->
      <div v-else-if="error" data-testid="inbox-error" class="p-lg text-center text-error font-body-sm">
        {{ error }}
      </div>

      <!-- Empty -->
      <div v-else-if="isEmpty" data-testid="inbox-empty" class="p-lg text-center text-on-surface-variant font-body-sm">
        No messages found
      </div>

      <!-- Messages -->
      <template v-else>
        <NuxtLink
          v-for="message in messages"
          :key="message.id"
          :to="`/conversations/${message.id}`"
          :data-testid="`thread-link-${message.id}`"
          class="block"
        >
          <div
            :data-testid="`message-item-${message.id}`"
            :class="activeId === message.id
              ? 'message-item-active group p-md rounded-xl cursor-pointer transition-all duration-200 hover:shadow-sm'
              : 'group p-md rounded-xl cursor-pointer transition-all duration-200 hover:bg-surface-container border-l-4 border-transparent'"
          >
            <div class="flex gap-md">
              <div class="relative flex-shrink-0">
                <img alt="Avatar" class="w-12 h-12 rounded-full object-cover" :src="message.avatar_url">
                <div class="absolute -bottom-0.5 -right-0.5 w-4 h-4 bg-white rounded-full flex items-center justify-center border border-outline-variant">
                  <div class="w-2.5 h-2.5 rounded-full" :style="{ backgroundColor: channelColors[message.channel] }" />
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex justify-between items-start mb-0.5">
                  <h4 class="font-body-lg text-body-lg font-bold text-on-surface truncate">{{ message.sender_name }}</h4>
                  <span
                    class="font-label-caps text-[10px] whitespace-nowrap"
                    :class="message.unread ? 'text-primary' : 'text-on-surface-variant'"
                  >{{ formatTime(message) }}</span>
                </div>
                <p
                  class="font-body-sm text-body-sm truncate mb-sm"
                  :class="message.status === 'pending' ? 'text-on-secondary-container' : 'text-on-surface-variant'"
                >{{ message.preview }}</p>
                <div class="flex items-center justify-between">
                  <span
                    v-if="message.status === 'pending'"
                    class="px-sm py-0.5 rounded-md bg-primary-container/10 text-primary font-status-label text-status-label border border-primary/20"
                  >Pending</span>
                  <span
                    v-else
                    class="px-sm py-0.5 rounded-md bg-tertiary-fixed text-tertiary font-status-label text-status-label border border-outline"
                  >Replied</span>
                  <div v-if="message.unread" class="w-2.5 h-2.5 rounded-full bg-primary" />
                </div>
              </div>
            </div>
          </div>
        </NuxtLink>
      </template>
    </div>

    <!-- Pagination -->
    <div
      v-if="totalPages > 1"
      data-testid="pagination-controls"
      class="px-lg py-md border-t border-outline-variant flex items-center justify-between gap-sm"
    >
      <button
        data-testid="pagination-prev"
        class="p-unit rounded hover:bg-surface-container transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
        :disabled="page <= 1"
        @click="store.goToPage(page - 1)"
      >
        <span class="material-symbols-outlined text-on-surface-variant">chevron_left</span>
      </button>
      <div class="flex items-center gap-xs">
        <button
          v-for="p in totalPages"
          :key="p"
          :data-testid="`pagination-page-${p}`"
          :class="page === p
            ? 'w-8 h-8 rounded-full bg-primary text-on-primary font-label-caps text-label-caps'
            : 'w-8 h-8 rounded-full text-on-surface-variant hover:bg-surface-container transition-colors font-label-caps text-label-caps'"
          @click="store.goToPage(p)"
        >
          {{ p }}
        </button>
      </div>
      <button
        data-testid="pagination-next"
        class="p-unit rounded hover:bg-surface-container transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
        :disabled="page >= totalPages"
        @click="store.goToPage(page + 1)"
      >
        <span class="material-symbols-outlined text-on-surface-variant">chevron_right</span>
      </button>
    </div>
  </section>
</template>
