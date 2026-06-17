<script setup lang="ts">
const navItems = [
  { icon: 'inbox', label: 'Inbox', active: true },
  { icon: 'send', label: 'Sent', active: false },
  { icon: 'archive', label: 'Archived', active: false },
  { icon: 'settings', label: 'Settings', active: false }
]

const bottomItems = [
  { icon: 'help', label: 'Help' },
  { icon: 'logout', label: 'Logout' }
]
</script>

<template>
  <aside class="fixed h-full w-sidebar-width left-0 top-0 border-r border-outline-variant bg-surface flex flex-col z-50">
    <div class="px-lg py-xl flex flex-col gap-sm">
      <div class="flex items-center gap-md mb-lg">
        <div class="w-10 h-10 bg-primary rounded-lg flex items-center justify-center text-on-primary">
          <span class="material-symbols-outlined" style="font-variation-settings: 'FILL' 1;">hub</span>
        </div>
        <div>
          <h1 class="font-headline-md text-headline-md font-bold text-primary">OmniChat</h1>
          <p class="font-label-caps text-label-caps text-on-surface-variant uppercase tracking-widest">Aggregator</p>
        </div>
      </div>
      <button class="w-full bg-primary text-on-primary py-md px-lg rounded-xl font-headline-md text-headline-md mb-lg flex items-center justify-center gap-sm shadow-md active:scale-95 transition-transform">
        <span class="material-symbols-outlined">edit_square</span>
        Compose
      </button>
      <nav class="flex flex-col gap-xs">
        <a
          v-for="item in navItems"
          :key="item.label"
          href="#"
          :class="item.active
            ? 'flex items-center gap-md px-lg py-md border-l-4 border-primary bg-secondary-container text-on-secondary-container'
            : 'flex items-center gap-md px-lg py-md text-on-surface-variant hover:bg-surface-container transition-colors'"
        >
          <span
            class="material-symbols-outlined"
            :style="item.active ? 'font-variation-settings: \'FILL\' 1;' : ''"
          >{{ item.icon }}</span>
          <span class="font-body-lg text-body-lg">{{ item.label }}</span>
        </a>
      </nav>
    </div>
    <div class="mt-auto px-lg py-lg border-t border-outline-variant flex flex-col gap-xs">
      <a
        v-for="item in bottomItems"
        :key="item.label"
        href="#"
        class="flex items-center gap-md px-lg py-md text-on-surface-variant hover:bg-surface-container transition-colors"
      >
        <span class="material-symbols-outlined">{{ item.icon }}</span>
        <span class="font-body-lg text-body-lg">{{ item.label }}</span>
      </a>
      <div class="flex items-center gap-md p-md mt-md">
        <img
          alt="User profile avatar"
          class="w-10 h-10 rounded-full border-2 border-primary-container object-cover"
          src="https://lh3.googleusercontent.com/aida-public/AB6AXuDx8JKscxrdzrvOsJD4vWuqN-umX1-r2fpJ_hU0j2LJSz7kZksyG_85Y6pVZlVB2mPzYg4Qtfi1PEMZYUQS4lMAAbanshwz9jKTCNIB33Q3tCC2yy1HFcGJOyBMgBpI1rX9-mOGTUk41SQjZnfR9Yemrzp--yNXjXeZPOWLY3FFsvkk4eSDofNacosZ_9BS_8DFkQgPkeSqSaFyg30D7ESeUZSQPGqIjF2SsgovwXaSWxF80G56byUPZI4bcK4M22lQ7zUlaLqZX54"
        >
        <div class="overflow-hidden">
          <p class="font-label-caps text-label-caps text-on-surface font-bold truncate">Alex Chen</p>
          <p class="text-[10px] text-on-surface-variant truncate">Admin Access</p>
        </div>
      </div>
    </div>
  </aside>
</template>
