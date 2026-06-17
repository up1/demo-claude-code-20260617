<script setup lang="ts">
const inputActions = ['attach_file', 'sentiment_satisfied', 'image']
</script>

<template>
  <section class="flex-1 flex flex-col bg-surface relative">
    <div class="absolute inset-0 pointer-events-none">
      <div class="absolute top-[-10%] right-[-10%] w-[40%] h-[40%] bg-primary/5 blur-[120px] rounded-full" />
      <div class="absolute bottom-[-10%] left-[-10%] w-[30%] h-[30%] bg-secondary/5 blur-[100px] rounded-full" />
    </div>

    <div class="h-16 px-lg flex items-center justify-between border-b border-outline-variant relative z-10">
      <div class="flex items-center gap-md">
        <img
          alt="Profile"
          class="w-9 h-9 rounded-full object-cover"
          src="https://lh3.googleusercontent.com/aida-public/AB6AXuCR_lvb2SPRATSTvfjTT6PyL5mhgOBuuB5Kp1-3wdHD1Pu3PNRNoTIqyZUbm65V2I0btszMpHxuL78d7EMhCmLvUK5Msr5xwroV2LTip-5MEiVaZf_jGk38ByO-Vrk54Z7qkjh7svxowJHu-OXvdpJBx8HUZ-3_oejrDjCcsOM9BwlmpiqAsQBStUNL5Z7N1l321ZCQwuiEWpK57DuS8TetiT8B9ACNZsh3msfvnXp7OTUeZbQblDNZ9zeFKPhuAD7ZRhgtdyJx77k"
        >
        <div>
          <h3 class="font-body-lg text-body-lg font-bold text-on-surface">Marcus Watanabe</h3>
          <div class="flex items-center gap-xs">
            <span class="w-2 h-2 rounded-full bg-green-500" />
            <span class="text-label-caps font-label-caps text-on-surface-variant">Online via LINE</span>
          </div>
        </div>
      </div>
      <div class="flex items-center gap-md">
        <button class="w-10 h-10 rounded-full border border-outline-variant flex items-center justify-center text-on-surface-variant hover:bg-surface-container transition-all active:scale-95">
          <span class="material-symbols-outlined">call</span>
        </button>
        <button class="w-10 h-10 rounded-full border border-outline-variant flex items-center justify-center text-on-surface-variant hover:bg-surface-container transition-all active:scale-95">
          <span class="material-symbols-outlined">videocam</span>
        </button>
        <button class="px-md py-2 rounded-lg border border-outline-variant flex items-center gap-xs font-label-caps text-label-caps text-on-surface-variant hover:bg-surface-container transition-all">
          <span class="material-symbols-outlined text-[18px]">more_vert</span>
        </button>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto p-lg space-y-lg relative z-10">
      <div class="flex justify-center my-md">
        <span class="px-lg py-1 rounded-full bg-surface-container-high text-on-surface-variant font-label-caps text-label-caps uppercase">Today, 10:15 AM</span>
      </div>
      <div class="flex gap-md max-w-[70%]">
        <img
          alt="User"
          class="w-8 h-8 rounded-full mt-auto"
          src="https://lh3.googleusercontent.com/aida-public/AB6AXuChO25ysKRsKwCCVxMjOhAcv_Q-raNko7a8F3MtNwpHLYO_WXg1N5QSbnb06e76JU3W7cmdLj-YwgYmvHDnFvRLhv0OcMz7Sx09vrGzb-qpopYArVRdGzFCGQ0D5q0NoHt-U9AXdsBiwrZvudORRtAgJtFLzl6DYeRboHRd4Sn7nIXRIYIzfGXJIOdiatcO0PqWslMhp-24vCnqg9vVeAodXUdXyIv5KS4JJiROJYi_Vguhyz4GsNPymcYwLQ4n6oytH2ueSIsrQig"
        >
        <div class="bg-surface-container p-lg rounded-2xl rounded-bl-none shadow-sm border border-outline-variant/30">
          <p class="font-body-lg text-body-lg text-on-surface">Hello! I'm checking on order #8812. It was supposed to arrive yesterday but I haven't received it yet.</p>
        </div>
      </div>
      <div class="flex gap-md max-w-[70%] ml-10">
        <div class="bg-surface-container p-lg rounded-2xl rounded-bl-none shadow-sm border border-outline-variant/30">
          <p class="font-body-lg text-body-lg text-on-surface">Can you confirm the tracking number for order #8812?</p>
        </div>
      </div>
      <div class="flex flex-row-reverse gap-md max-w-[70%] ml-auto">
        <div class="bg-primary p-lg rounded-2xl rounded-br-none shadow-md">
          <p class="font-body-lg text-body-lg text-on-primary">Hi Marcus! Let me check that for you right away. One moment please while I access the warehouse logs.</p>
        </div>
      </div>
      <div class="flex gap-md max-w-[70%] items-center animate-pulse">
        <div class="bg-surface-container-high px-md py-2 rounded-full border border-outline-variant/20 flex gap-1">
          <div class="w-1.5 h-1.5 bg-on-surface-variant/40 rounded-full" />
          <div class="w-1.5 h-1.5 bg-on-surface-variant/40 rounded-full" />
          <div class="w-1.5 h-1.5 bg-on-surface-variant/40 rounded-full" />
        </div>
      </div>
    </div>

    <div class="p-lg bg-surface border-t border-outline-variant relative z-10">
      <div class="bg-surface-container-low border border-outline-variant rounded-2xl p-md focus-within:ring-2 focus-within:ring-primary focus-within:border-primary transition-all">
        <textarea
          class="w-full bg-transparent border-none focus:ring-0 text-body-lg scrollbar-hide min-h-[60px] resize-none"
          placeholder="Type your message here..."
        />
        <div class="flex items-center justify-between mt-md border-t border-outline-variant/30 pt-md">
          <div class="flex items-center gap-md">
            <button
              v-for="action in inputActions"
              :key="action"
              class="text-on-surface-variant hover:text-primary transition-colors"
            >
              <span class="material-symbols-outlined">{{ action }}</span>
            </button>
            <div class="h-4 w-[1px] bg-outline-variant" />
            <button class="flex items-center gap-xs px-sm py-1 rounded-md hover:bg-surface-container transition-colors text-on-surface-variant">
              <span class="material-symbols-outlined text-[18px]">bolt</span>
              <span class="font-label-caps text-label-caps">Templates</span>
            </button>
          </div>
          <button class="bg-primary text-on-primary px-lg py-2 rounded-xl font-headline-md text-headline-md flex items-center gap-xs shadow-md active:scale-95 transition-transform">
            Send
            <span class="material-symbols-outlined" style="font-variation-settings: 'FILL' 1;">send</span>
          </button>
        </div>
      </div>
    </div>
  </section>
</template>
