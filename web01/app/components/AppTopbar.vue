<template>
  <header class="h-16 flex items-center justify-between px-lg w-full border-b border-outline-variant bg-surface sticky top-0 z-40">
    <div class="flex items-center gap-md">
      <h2 class="font-headline-md text-headline-md font-bold text-primary">OmniChat Dashboard</h2>
    </div>
    <div class="flex items-center gap-lg">
      <div class="relative group">
        <span class="absolute left-3 top-1/2 -translate-y-1/2 material-symbols-outlined text-outline text-sm">search</span>
        <input
          class="bg-surface-container-low border-none rounded-full pl-10 pr-lg py-1.5 text-body-sm focus:ring-2 focus:ring-primary w-64 transition-all"
          placeholder="Search across all channels..."
          type="text"
        >
      </div>
      <div class="flex items-center gap-md">
        <button class="text-on-surface-variant hover:text-primary transition-colors relative">
          <span class="material-symbols-outlined">notifications</span>
          <span class="absolute top-0 right-0 w-2 h-2 bg-error rounded-full border-2 border-surface" />
        </button>
        <button class="text-on-surface-variant hover:text-primary transition-colors">
          <span class="material-symbols-outlined">apps</span>
        </button>
        <div class="h-8 w-8 rounded-full overflow-hidden border border-outline-variant">
          <img
            alt="Administrator Profile"
            class="w-full h-full object-cover"
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuAOCFZhrUYwc_jE3G-Dr6pQImEp1KMk3bWs7dsiYRNF7ctNFpEpS_kQDLXrIWq2GVZ_Ga6JQQdD8cRLpbjLHyisESxqrveDqzncuFuRsZ6zGmL7f-bJap20rVSCiFtd5ASHOqMntYIPWHYA-U2dbOHg4JQF08uO25877oZ3OIbfnf5IOMTTUd8LrjjYOTE8oZ-mPOzXBFlMcmNuCLFgZSHFmBNmrLLuaF6BaCaYZuo_1q0eZt2em0Vr4ytDjEiqln4Z9AZTXPCip_M"
          >
        </div>
      </div>
    </div>
  </header>
</template>
