<script setup lang="ts">
interface Message {
  id: number
  name: string
  avatar: string
  channelColor: string
  time: string
  timeHighlighted: boolean
  preview: string
  status: 'Pending' | 'Replied'
}

const channels = [
  { label: 'All', color: '' },
  { label: 'Facebook', color: '#1877F2' },
  { label: 'LINE', color: '#06C755' },
  { label: 'Instagram', color: '#E4405F' }
]

const messages: Message[] = [
  {
    id: 1,
    name: 'Marcus Watanabe',
    avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuAwPp5iHNxM-YtHR_qdfKPxxqu4m0gad5q2rXW_LB86zEQ-hxHGB7sjYpuC5QnykUBJeC5iBX7Z5-ExZdVgSFDKoqGe_1IOnfDu4Z19smwLwAOKqwQr1Demseztpb09ZIsUn_VxBThlgFRrPG0EfkFEhE8w-egl0dUA-rUL1DyH0t7t_0JnTGRZc_LNWJnSSz19C6NNMmR-MwY8sAaqi4HsjsNtkDwClLCZNRVFRPYbo7YzD-hsreR9b4xjLfMMi1pg2vK1EvGACsc',
    channelColor: '#06C755',
    time: '10:30 AM',
    timeHighlighted: true,
    preview: 'Can you confirm the tracking number for order #8812?',
    status: 'Pending'
  },
  {
    id: 2,
    name: 'Sarah Jenkins',
    avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuDGuT3Z3u_JOCf3Nlp1xIPLet0oDLymZkoNEBYiN91tpz1cWvqgLq46i7mdsLtqHVXNBr1w22GFljaDmdZ0x_vUpdBrHnLOR0EgZDVuVUKow3c61Cf_TtnLsRVqsihuDuBRX3kwYveJMXytBPcC39n_nwS4lnNLmYDd4OsXRWkH0WiJNdc1aZON5-v-G4ENs_5MK1w1_hhJQz8jx1dC9ag540SAE40KzfyNrZLb_Mwg3kGNIpATzZwgrZkcnqYv_BxDyezwCJ1K4tM',
    channelColor: '#1877F2',
    time: 'Yesterday',
    timeHighlighted: false,
    preview: 'Thank you for the quick resolution!',
    status: 'Replied'
  },
  {
    id: 3,
    name: 'David Wilson',
    avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuCOsS7cbdAa6-xRklsbbgouQLT46LSkBRFkzbiTMWDJRYps1cab4WVMfe8b9eH41rPKv5dPAcP0g9q1tmMOdnU2KTecD7_Sm8kHBmQEep1iwqP_93h0KRv5OPnez9SAsZSkVgL6bO7JJAdZmImK9cKLaT0TKE36O7_2YjCtKg4CUWTpe9rTHNRTHCsn-V6KAOAz9DONS6CmpJhOdfLKHps-3U-bCIcPfHGkwsx6P1sXF1U1i_EVTlGbFRBcq2VBiOIhhYAr-IVzvEc',
    channelColor: '#E4405F',
    time: '2:15 PM',
    timeHighlighted: true,
    preview: 'Do you have the new winter collection in stock yet?',
    status: 'Pending'
  },
  {
    id: 4,
    name: 'Li Na',
    avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuC0-d54lHHaH-FX3ThH0MFOzRW0OLgoe9_wLynLHAu9XIjzWxqzCswbU8ut7DyKlm98UvstJ95nkFsmM4RTrV8dSlXJ0y4l_rZqSx4SxN65-Fzoj5SUIrm2pR3t9db9Wppxz3ZAjVwoN1TKj6czqsyJecKkKRyBq1J1mH8bt5oQQkDF-KB1EShyMzCJBjOWbY1cZsqvs5U-UYBl3uqKdVUcDGS_c9R_cLElJ5J1ZlI4YP5pZKLIL3IEObUnDbFIYLwTj6vSBW8ztbA',
    channelColor: '#06C755',
    time: 'Monday',
    timeHighlighted: false,
    preview: "I'd like to update my shipping address for the recent order.",
    status: 'Replied'
  }
]

const activeId = ref(messages[0]!.id)
const activeChannel = ref('All')

function selectMessage(id: number) {
  activeId.value = id
}
</script>

<template>
  <section class="w-chat-list-width flex flex-col border-r border-outline-variant bg-surface-container-lowest">
    <div class="p-lg space-y-md">
      <div class="flex items-center justify-between">
        <h3 class="font-display-lg text-display-lg">
          Inbox <span class="text-on-surface-variant font-normal text-headline-md ml-xs">(128)</span>
        </h3>
        <button class="p-unit rounded hover:bg-surface-container transition-colors">
          <span class="material-symbols-outlined text-on-surface-variant">filter_list</span>
        </button>
      </div>
      <div class="flex gap-sm overflow-x-auto scrollbar-hide pb-unit">
        <button
          v-for="channel in channels"
          :key="channel.label"
          :class="activeChannel === channel.label
            ? 'px-md py-1.5 rounded-full bg-primary text-on-primary font-label-caps text-label-caps flex items-center gap-xs'
            : 'px-md py-1.5 rounded-full bg-surface-container text-on-surface-variant hover:bg-surface-container-high transition-colors font-label-caps text-label-caps flex items-center gap-xs'"
          @click="activeChannel = channel.label"
        >
          <span
            v-if="channel.color"
            class="w-2 h-2 rounded-full"
            :style="{ backgroundColor: channel.color }"
          />
          {{ channel.label }}
        </button>
      </div>
      <div class="flex gap-sm">
        <div class="flex-1 relative">
          <span class="absolute left-2.5 top-1/2 -translate-y-1/2 material-symbols-outlined text-[18px] text-outline">calendar_today</span>
          <select class="w-full bg-surface-container-low border-none rounded-lg pl-9 pr-md py-2 text-body-sm appearance-none focus:ring-1 focus:ring-primary">
            <option>Last 7 days</option>
            <option>Last 30 days</option>
            <option>Custom Range</option>
          </select>
        </div>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto scrollbar-hide px-sm space-y-xs pb-lg">
      <div
        v-for="message in messages"
        :key="message.id"
        :data-testid="`message-item-${message.id}`"
        :class="activeId === message.id
          ? 'message-item-active group p-md rounded-xl cursor-pointer transition-all duration-200 hover:shadow-sm'
          : 'group p-md rounded-xl cursor-pointer transition-all duration-200 hover:bg-surface-container border-l-4 border-transparent'"
        @click="selectMessage(message.id)"
      >
        <div class="flex gap-md">
          <div class="relative flex-shrink-0">
            <img alt="Avatar" class="w-12 h-12 rounded-full object-cover" :src="message.avatar">
            <div class="absolute -bottom-0.5 -right-0.5 w-4 h-4 bg-white rounded-full flex items-center justify-center border border-outline-variant">
              <div class="w-2.5 h-2.5 rounded-full" :style="{ backgroundColor: message.channelColor }" />
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-start mb-0.5">
              <h4 class="font-body-lg text-body-lg font-bold text-on-surface truncate">{{ message.name }}</h4>
              <span
                class="font-label-caps text-[10px] whitespace-nowrap"
                :class="message.timeHighlighted ? 'text-primary' : 'text-on-surface-variant'"
              >{{ message.time }}</span>
            </div>
            <p
              class="font-body-sm text-body-sm truncate mb-sm"
              :class="message.status === 'Pending' ? 'text-on-secondary-container' : 'text-on-surface-variant'"
            >{{ message.preview }}</p>
            <div class="flex items-center justify-between">
              <span
                v-if="message.status === 'Pending'"
                class="px-sm py-0.5 rounded-md bg-primary-container/10 text-primary font-status-label text-status-label border border-primary/20"
              >Pending</span>
              <span
                v-else
                class="px-sm py-0.5 rounded-md bg-tertiary-fixed text-tertiary font-status-label text-status-label border border-outline"
              >Replied</span>
              <div v-if="message.status === 'Pending'" class="w-2.5 h-2.5 rounded-full bg-primary" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
