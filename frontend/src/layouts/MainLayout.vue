<script setup lang="ts">
  import SidebarMain from '@/components/SidebarMain.vue'
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="flex h-screen">
      <!-- Sidebar -->
      <aside class="w-64 bg-white shadow-lg">
        <div class="h-full px-3 py-4 overflow-y-auto">
          <SidebarMain />
        </div>
      </aside>
      
      <!-- Main content -->
      <main class="flex-1 overflow-y-auto bg-gray-50">
        <div class="container mx-auto px-6 py-8">
          <slot name="main"></slot>
        </div>
      </main>
    </div>
  </div>
</template>
