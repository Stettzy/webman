<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import Collections from './Collections.vue'
import Logo from './Logo.vue'

const route = useRoute()
const collections = ref<InstanceType<typeof Collections>>()

const navigation = [
  {
    icon: 'pi pi-compass',
    name: 'Explorer',
    path: '/',
  }
]
</script>

<template>
  <div class="flex flex-col h-full bg-white">
    <!-- App Logo -->
    <div class="flex items-center h-14 px-5">
      <div class="flex items-center space-x-2.5">
        <Logo size="sm" color="#4F46E5" />
        <span class="text-sm font-medium text-gray-800">Webman</span>
      </div>
    </div>

    <!-- Navigation -->
    <div class="px-3 py-3">
      <RouterLink v-for="(item, index) in navigation" :key="index" :to="item.path"
        class="flex items-center px-2 py-1.5 text-sm text-gray-600 rounded-md hover:bg-gray-50 transition-colors duration-150"
        :class="{ 'bg-gray-50 text-gray-900': route.path === item.path }">
        <i :class="[item.icon, 'mr-2 text-base']"></i>
        <span class="font-medium">{{ item.name }}</span>
      </RouterLink>
    </div>

    <!-- Collections -->
    <div class="flex-1 overflow-y-auto">
      <div class="px-5 py-2">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-gray-400 uppercase tracking-wide">Collections</span>
          <button @click="collections?.openNewCollectionModal()"
            class="p-1 text-gray-400 hover:text-gray-600 rounded transition-colors">
            <i class="pi pi-plus text-xs"></i>
          </button>
        </div>
      </div>
      <div class="px-3">
        <Collections ref="collections" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.router-link-active {
  background-color: rgb(249 250 251);
  color: rgb(17 24 39);
}
</style>
