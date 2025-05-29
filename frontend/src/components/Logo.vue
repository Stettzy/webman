# Logo.vue
<script setup lang="ts">
defineProps<{
    size?: 'sm' | 'md' | 'lg'
    color?: string
}>()
</script>

<template>
    <div class="logo" :class="size || 'md'">
        <svg viewBox="0 0 120 60" fill="none" xmlns="http://www.w3.org/2000/svg">
            <!-- Curvy W with wider proportions -->
            <path :stroke="color || 'currentColor'"
                d="M10 45 C20 45, 25 15, 30 45 C35 45, 40 25, 60 25 C80 25, 85 45, 90 45 C95 45, 100 15, 110 45"
                stroke-width="6" stroke-linecap="round" stroke-linejoin="round" fill="none" />
        </svg>
    </div>
</template>

<style scoped>
.logo {
    display: inline-block;
}

.logo.sm {
    width: 32px;
    height: 20px;
}

.logo.md {
    width: 48px;
    height: 30px;
}

.logo.lg {
    width: 64px;
    height: 40px;
}
</style>