<template>
  <div class="flex h-screen overflow-hidden bg-surface text-on-surface">
    <AppSidebar />
    <main class="ml-sidebar-width flex-1 flex flex-col h-full bg-surface">
      <AppTopbar />
      <div class="flex-1 flex overflow-hidden">
        <InboxList />
        <slot />
      </div>
    </main>
  </div>
</template>
