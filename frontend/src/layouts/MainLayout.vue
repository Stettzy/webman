<script setup lang="ts">
  import SidebarMain from '@/components/SidebarMain.vue'
</script>

<template>
  <div class="overflow-hidden flex">
    <aside class="p-10 flex-1 w-full flex flex-col border-r border-gray:50 min-h-screen overflow-y-auto">
      <SidebarMain />
    </aside>
    <main class="flex-10">
      <div class="flex flex-col">
        <slot name="main"></slot>
      </div>
    </main>
  </div>
</template>
