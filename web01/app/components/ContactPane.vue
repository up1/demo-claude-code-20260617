<script setup lang="ts">
const stats = [
  { label: 'Orders', value: '12', highlight: false },
  { label: 'LTV', value: '$1.4k', highlight: true }
]

const contactInfo = [
  { icon: 'mail', value: 'marcus.w@example.com' },
  { icon: 'phone', value: '+81 90-1234-5678' },
  { icon: 'alternate_email', value: '@marcus_dev' }
]

const tags = [
  { label: 'VIP_CUSTOMER', variant: 'secondary' },
  { label: 'TECH_EARLY_ADOPTER', variant: 'secondary' },
  { label: 'URGENT_TICKET', variant: 'error' }
]
</script>

<template>
  <aside class="w-sidebar-width bg-surface-container-lowest border-l border-outline-variant p-lg hidden xl:flex flex-col gap-lg overflow-y-auto scrollbar-hide">
    <div class="text-center py-md border-b border-outline-variant/50">
      <img
        alt="Profile Large"
        class="w-24 h-24 rounded-full mx-auto mb-md border-4 border-surface shadow-lg"
        src="https://lh3.googleusercontent.com/aida-public/AB6AXuCP-XGEOaM99C9Hsnn8J8L1QDWnff40Bq9Q5lzwZ4VEm6osmzF0ewh2aXpC8hkY9YFg7VxmedZItDYhHVcATuUN4GCNtgBO8ip-yKGlWtOeZau2JLpQjKZUpexyDSE0gbltjAskDGJ067p-_5n1hgPGs7wMgaOKrX6jlu0kGUo6d6PGoJ6slIZiCi3QCOZAoWi5xMjHMSK0YW9BwG3TF8OtOOHwurkEwlOQ7dhr5VtpkgpWYbO_8grGz-62SgSzi6QQM_q4zvpwtbI"
      >
      <h4 class="font-headline-md text-headline-md font-bold text-on-surface">Marcus Watanabe</h4>
      <p class="font-body-sm text-body-sm text-on-surface-variant">Tokyo, Japan (GMT+9)</p>
    </div>
    <div class="grid grid-cols-2 gap-sm">
      <div
        v-for="stat in stats"
        :key="stat.label"
        class="bg-surface-container p-md rounded-xl border border-outline-variant/30"
      >
        <p class="font-label-caps text-label-caps text-on-surface-variant mb-xs uppercase">{{ stat.label }}</p>
        <p
          class="font-headline-md text-headline-md font-bold"
          :class="stat.highlight ? 'text-primary' : ''"
        >{{ stat.value }}</p>
      </div>
    </div>
    <div class="space-y-md">
      <h5 class="font-label-caps text-label-caps text-on-surface-variant uppercase tracking-widest border-b border-outline-variant pb-xs">Contact Info</h5>
      <div class="space-y-sm">
        <div
          v-for="info in contactInfo"
          :key="info.icon"
          class="flex items-center gap-md"
        >
          <span class="material-symbols-outlined text-outline">{{ info.icon }}</span>
          <span class="text-body-sm text-on-surface truncate">{{ info.value }}</span>
        </div>
      </div>
    </div>
    <div class="space-y-md">
      <h5 class="font-label-caps text-label-caps text-on-surface-variant uppercase tracking-widest border-b border-outline-variant pb-xs">Tags</h5>
      <div class="flex flex-wrap gap-xs">
        <span
          v-for="tag in tags"
          :key="tag.label"
          :class="tag.variant === 'error'
            ? 'px-sm py-1 rounded bg-error-container text-on-error-container text-[10px] font-bold'
            : 'px-sm py-1 rounded bg-secondary-container text-on-secondary-fixed-variant text-[10px] font-bold'"
        >{{ tag.label }}</span>
      </div>
    </div>
    <div class="mt-auto pt-lg">
      <button class="w-full py-md border border-outline text-on-surface-variant rounded-xl font-headline-md text-headline-md hover:bg-surface-container transition-all">
        View Full Profile
      </button>
    </div>
  </aside>
</template>
