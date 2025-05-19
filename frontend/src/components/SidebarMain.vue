<script setup lang="ts">
  const navigation = [
    {
      icon: 'pi pi-home',
      name: 'Dashboard',
      path: '/',
    },
    {
      icon: 'pi pi-briefcase',
      name: 'Collections',
      path: '/collections',
    }
  ];
</script>

<template>
  <ul>
    <li :key="index" v-for="(item, index ) in navigation">
      <i :class="item.icon" class="mr-2"></i>
      <span class="inline-block">{{ item.name}}</span>
    </li>
  </ul>
</template>
