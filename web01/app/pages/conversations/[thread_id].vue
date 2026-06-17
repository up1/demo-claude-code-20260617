<script setup lang="ts">
// Conversation view for a single thread. Navigation-only for now: the thread id
// is available from the route so a future iteration can load conversation
// details via GET /api/v1/conversations/{thread_id}.
const route = useRoute()
const threadId = computed(() => String(route.params.thread_id))
</script>

<template>
  <div class="contents">
    <ChatPanel :data-thread-id="threadId" />
    <ContactPane />
  </div>
</template>
