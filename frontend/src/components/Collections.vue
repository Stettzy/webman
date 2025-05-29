<!-- Collections.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useCollectionsStore } from '../stores/collections'

const collectionsStore = useCollectionsStore()
const showNewCollectionModal = ref(false)
const newCollection = ref({ name: '', description: '' })
const expandedCollections = ref<Set<string>>(new Set())

// Expose methods for parent components
defineExpose({
    openNewCollectionModal: () => {
        showNewCollectionModal.value = true
    }
})

onMounted(() => {
    collectionsStore.fetchCollections()
})

const createCollection = async () => {
    try {
        await collectionsStore.createCollection(newCollection.value.name, newCollection.value.description)
        newCollection.value = { name: '', description: '' }
        showNewCollectionModal.value = false
    } catch (error) {
        console.error('Failed to create collection:', error)
    }
}

const toggleCollection = (id: string) => {
    if (expandedCollections.value.has(id)) {
        expandedCollections.value.delete(id)
    } else {
        expandedCollections.value.add(id)
    }
}

const loadRequest = (collectionId: string, request: any) => {
    collectionsStore.currentCollection = collectionsStore.collections.find(c => c.id === collectionId) || null
    if (collectionsStore.currentCollection) {
        collectionsStore.currentRequest = request
    }
}

const deleteRequest = async (collectionId: string, requestId: string) => {
    if (!confirm('Are you sure you want to delete this request?')) return
    try {
        await collectionsStore.deleteRequest(collectionId, requestId)
    } catch (error) {
        console.error('Failed to delete request:', error)
    }
}

const deleteCollection = async (id: string) => {
    if (!confirm('Are you sure you want to delete this collection?')) return
    try {
        await collectionsStore.deleteCollection(id)
    } catch (error) {
        console.error('Failed to delete collection:', error)
    }
}
</script>

<template>
    <div class="h-full">
        <!-- Collections List -->
        <div class="space-y-1">
            <div v-for="collection in collectionsStore.collections" :key="collection.id"
                class="rounded-md overflow-hidden">
                <!-- Collection Header -->
                <div class="flex items-center px-2 py-1.5 hover:bg-gray-50 cursor-pointer group"
                    @click="toggleCollection(collection.id)">
                    <svg xmlns="http://www.w3.org/2000/svg"
                        class="h-3.5 w-3.5 mr-1.5 text-gray-400 transition-transform"
                        :class="{ 'rotate-90': expandedCollections.has(collection.id) }" viewBox="0 0 20 20"
                        fill="currentColor">
                        <path fill-rule="evenodd"
                            d="M7.293 4.293a1 1 0 011.414 0L14.414 10l-5.707 5.707a1 1 0 01-1.414-1.414L11.586 10 7.293 5.707a1 1 0 010-1.414z"
                            clip-rule="evenodd" />
                    </svg>
                    <span class="flex-1 text-sm text-gray-600">{{ collection.name }}</span>
                    <button @click.stop="deleteCollection(collection.id)"
                        class="p-1 text-gray-400 hover:text-red-500 rounded opacity-0 group-hover:opacity-100 transition-opacity">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20"
                            fill="currentColor">
                            <path fill-rule="evenodd"
                                d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                                clip-rule="evenodd" />
                        </svg>
                    </button>
                </div>

                <!-- Requests List -->
                <div v-if="expandedCollections.has(collection.id)" class="ml-5 space-y-0.5">
                    <div v-for="request in collection.requests" :key="request.id"
                        class="flex items-center px-2 py-1.5 hover:bg-gray-50 cursor-pointer group rounded-md"
                        @click="loadRequest(collection.id, request)">
                        <span class="w-10 text-xs font-medium text-gray-400">{{ request.method }}</span>
                        <span class="flex-1 text-sm text-gray-600 truncate">{{ request.name }}</span>
                        <button @click.stop="deleteRequest(collection.id, request.id)"
                            class="p-1 text-gray-400 hover:text-red-500 rounded opacity-0 group-hover:opacity-100 transition-opacity">
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" viewBox="0 0 20 20"
                                fill="currentColor">
                                <path fill-rule="evenodd"
                                    d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
                                    clip-rule="evenodd" />
                            </svg>
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- New Collection Modal -->
        <Transition enter-active-class="ease-out duration-200" enter-from-class="opacity-0" enter-to-class="opacity-100"
            leave-active-class="ease-in duration-150" leave-from-class="opacity-100" leave-to-class="opacity-0">
            <div v-if="showNewCollectionModal" class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        </Transition>

        <Transition enter-active-class="ease-out duration-300"
            enter-from-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
            enter-to-class="opacity-100 translate-y-0 sm:scale-100" leave-active-class="ease-in duration-200"
            leave-from-class="opacity-100 translate-y-0 sm:scale-100"
            leave-to-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
            <div v-if="showNewCollectionModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-0">
                <div
                    class="relative transform overflow-hidden rounded-lg bg-white shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg">
                    <div class="absolute right-0 top-0 pr-4 pt-4">
                        <button @click="showNewCollectionModal = false"
                            class="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none">
                            <span class="sr-only">Close</span>
                            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd"
                                    d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                                    clip-rule="evenodd" />
                            </svg>
                        </button>
                    </div>

                    <div class="px-6 pt-6 pb-4">
                        <h3 class="text-xl font-semibold text-gray-900 mb-6">Create New Collection</h3>
                        <div class="space-y-5">
                            <div>
                                <label for="collection-name"
                                    class="block text-sm font-medium text-gray-700 mb-1">Name</label>
                                <input id="collection-name" v-model="newCollection.name" type="text"
                                    class="block w-full rounded-md border-0 py-2.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm"
                                    placeholder="My Collection" />
                            </div>
                            <div>
                                <label for="collection-description"
                                    class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                                <textarea id="collection-description" v-model="newCollection.description" rows="3"
                                    class="block w-full rounded-md border-0 py-2.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm"
                                    placeholder="Optional description for your collection"></textarea>
                            </div>
                        </div>
                    </div>

                    <div class="bg-gray-50 px-6 py-4 flex items-center justify-end space-x-3">
                        <button @click="showNewCollectionModal = false"
                            class="inline-flex justify-center rounded-md bg-white px-4 py-2.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
                            Cancel
                        </button>
                        <button @click="createCollection"
                            class="inline-flex justify-center rounded-md bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
                            Create Collection
                        </button>
                    </div>
                </div>
            </div>
        </Transition>
    </div>
</template>