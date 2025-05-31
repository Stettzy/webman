import { defineStore } from 'pinia'
import { ref } from 'vue'

interface Request {
    id: string
    name: string
    method: string
    url: string
    headers: Record<string, string>
    body: string
    bodyType: string
    createdAt: string
    updatedAt: string
}

interface Collection {
    id: string
    name: string
    description: string
    requests: Request[]
    createdAt: string
    updatedAt: string
}

export const useCollectionsStore = defineStore('collections', () => {
    const collections = ref<Collection[]>([])
    const currentCollection = ref<Collection | null>(null)
    const currentRequest = ref<Request | null>(null)

    async function fetchCollections() {
        try {
            const response = await fetch('/api/collections')
            const data = await response.json()
            collections.value = data
        } catch (error) {
            console.error('Failed to fetch collections:', error)
        }
    }

    async function createCollection(name: string, description: string) {
        try {
            const response = await fetch('/api/collections', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name, description }),
            })
            const data = await response.json()
            collections.value.push(data)
            return data
        } catch (error) {
            console.error('Failed to create collection:', error)
            throw error
        }
    }

    async function updateCollection(id: string, name: string, description: string) {
        try {
            const response = await fetch(`/api/collections/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name, description }),
            })
            const data = await response.json()
            const index = collections.value.findIndex(c => c.id === id)
            if (index !== -1) {
                collections.value[index] = data
            }
            return data
        } catch (error) {
            console.error('Failed to update collection:', error)
            throw error
        }
    }

    async function deleteCollection(id: string) {
        try {
            await fetch(`/api/collections/${id}`, {
                method: 'DELETE',
            })
            collections.value = collections.value.filter(c => c.id !== id)
            if (currentCollection.value?.id === id) {
                currentCollection.value = null
                currentRequest.value = null
            }
        } catch (error) {
            console.error('Failed to delete collection:', error)
            throw error
        }
    }

    async function addRequest(collectionId: string, request: Omit<Request, 'id' | 'createdAt' | 'updatedAt'>) {
        try {
            const response = await fetch(`/api/collections/${collectionId}/requests`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(request),
            })
            if (!response.ok) {
                throw new Error('Failed to add request')
            }
            await fetchCollections() // Refresh collections to get the updated data
        } catch (error) {
            console.error('Failed to add request:', error)
            throw error
        }
    }

    async function updateRequest(collectionId: string, request: Request) {
        try {
            const response = await fetch(`/api/collections/${collectionId}/requests/${request.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(request),
            })
            if (!response.ok) {
                throw new Error('Failed to update request')
            }
            await fetchCollections() // Refresh collections to get the updated data
            if (currentRequest.value?.id === request.id) {
                currentRequest.value = request
            }
        } catch (error) {
            console.error('Failed to update request:', error)
            throw error
        }
    }

    async function deleteRequest(collectionId: string, requestId: string) {
        try {
            const response = await fetch(`/api/collections/${collectionId}/requests/${requestId}`, {
                method: 'DELETE',
            })
            if (!response.ok) {
                throw new Error('Failed to delete request')
            }
            await fetchCollections() // Refresh collections to get the updated data
            if (currentRequest.value?.id === requestId) {
                currentRequest.value = null
            }
        } catch (error) {
            console.error('Failed to delete request:', error)
            throw error
        }
    }

    return {
        collections,
        currentCollection,
        currentRequest,
        fetchCollections,
        createCollection,
        updateCollection,
        deleteCollection,
        addRequest,
        updateRequest,
        deleteRequest,
    }
}) 