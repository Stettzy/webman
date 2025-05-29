import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Dashboard from '../Dashboard.vue'
import MonacoEditor from 'monaco-editor-vue3'

// Mock MonacoEditor component
vi.mock('monaco-editor-vue3', () => ({
    default: {
        name: 'MonacoEditor',
        template: '<div class="monaco-editor"></div>'
    }
}))

describe('Dashboard.vue', () => {
    let wrapper: any

    beforeEach(() => {
        // Reset all mocks before each test
        vi.clearAllMocks()

        wrapper = mount(Dashboard, {
            global: {
                components: {
                    MonacoEditor
                }
            }
        })
    })

    it('initializes with default values', () => {
        expect(wrapper.vm.requestUrl).toBe('')
        expect(wrapper.vm.bodyContent).toContain('"key": "value"')
        expect(wrapper.vm.params).toEqual([])
        expect(wrapper.vm.headers).toEqual([])
        expect(wrapper.vm.selectedOption).toBe('Params')
        expect(wrapper.vm.selectedMethod).toBe('POST')
        expect(wrapper.vm.selectedBodyType).toBe('JSON')
    })

    it('adds and removes params correctly', async () => {
        const paramKey = 'testKey'
        const paramValue = 'testValue'

        // Set param values directly
        wrapper.vm.paramPairs = { key: paramKey, value: paramValue }
        await wrapper.vm.$nextTick()

        // Find and click the Add button in the Params section
        const addButton = wrapper.find('button.px-4.py-2\\.5')
        await addButton.trigger('click')

        // Verify param was added
        expect(wrapper.vm.params).toContainEqual({ key: paramKey, value: paramValue })

        // Find and click the remove button
        const removeButton = wrapper.find('button.text-red-600')
        await removeButton.trigger('click')

        // Verify param was removed
        expect(wrapper.vm.params).toHaveLength(0)
    })

    it('adds and removes headers correctly', async () => {
        // Switch to Headers tab directly
        wrapper.vm.selectedOption = 'Headers'
        await wrapper.vm.$nextTick()

        const headerKey = 'Content-Type'
        const headerValue = 'application/json'

        // Set header values directly
        wrapper.vm.paramPairs = { key: headerKey, value: headerValue }
        await wrapper.vm.$nextTick()

        // Find and click the Add button
        const addButton = wrapper.find('button.px-4.py-2\\.5')
        await addButton.trigger('click')

        // Verify header was added
        expect(wrapper.vm.headers).toContainEqual({ key: headerKey, value: headerValue })

        // Find and click the remove button
        const removeButton = wrapper.find('button.text-red-600')
        await removeButton.trigger('click')

        // Verify header was removed
        expect(wrapper.vm.headers).toHaveLength(0)
    })

    it('formats JSON correctly', async () => {
        const unformattedJson = '{"key":"value","nested":{"inner":"value"}}'
        const formattedJson = `{
  "key": "value",
  "nested": {
    "inner": "value"
  }
}`
        // Set values directly
        wrapper.vm.selectedOption = 'Body'
        wrapper.vm.selectedBodyType = 'JSON'
        wrapper.vm.bodyContent = unformattedJson
        await wrapper.vm.$nextTick()

        // Find and click the Format JSON button
        const formatButton = wrapper.find('button.px-4.py-2')
        await formatButton.trigger('click')

        // Verify JSON was formatted
        expect(wrapper.vm.bodyContent).toBe(formattedJson)
    })

    it('sends request and handles response correctly', async () => {
        const responseBody = '{"success":true}'
        const formattedResponseBody = `{
  "success": true
}`
        const mockResponse = {
            statusCode: 200,
            body: btoa(responseBody),
            headers: { 'Content-Type': 'application/json' }
        }

        // Mock fetch
        global.fetch = vi.fn().mockResolvedValue({
            json: () => Promise.resolve(mockResponse)
        })

        // Set request data directly
        wrapper.vm.requestUrl = 'http://api.example.com'
        wrapper.vm.bodyContent = '{"test":"data"}'
        wrapper.vm.headers = [{ key: 'Content-Type', value: 'application/json' }]
        await wrapper.vm.$nextTick()

        // Find and click the Send button
        const sendButton = wrapper.find('button.bg-blue-600')
        await sendButton.trigger('click')

        // Verify request was sent
        expect(global.fetch).toHaveBeenCalledWith(
            'http://localhost:9090',
            expect.objectContaining({
                method: 'POST',
                headers: { 'Content-Type': 'application/json' }
            })
        )

        // Wait for response processing
        await wrapper.vm.$nextTick()

        // Verify response handling
        expect(wrapper.vm.responseData).toBeTruthy()
        expect(wrapper.vm.responseData.status).toBe(200)
        expect(wrapper.vm.responseData.body).toBe(formattedResponseBody)
    })

    it('handles request errors gracefully', async () => {
        // Mock console.error to prevent test output noise
        const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => { })

        // Mock fetch to reject
        global.fetch = vi.fn().mockRejectedValue(new Error('Network error'))

        // Set request data directly
        wrapper.vm.requestUrl = 'http://api.example.com'
        await wrapper.vm.$nextTick()

        // Find and click the Send button
        const sendButton = wrapper.find('button.bg-blue-600')
        await sendButton.trigger('click')

        // Wait for error handling
        await wrapper.vm.$nextTick()

        // Verify error handling
        expect(consoleSpy).toHaveBeenCalled()
        expect(wrapper.vm.responseData).toEqual({
            body: 'Request failed: Network error',
            headers: [],
            status: 500
        })

        consoleSpy.mockRestore()
    })

    it('sends request with URL parameters correctly', async () => {
        const mockResponse = {
            statusCode: 200,
            body: btoa('{"success":true}'),
            headers: { 'Content-Type': 'application/json' }
        }

        // Mock fetch
        global.fetch = vi.fn().mockResolvedValue({
            json: () => Promise.resolve(mockResponse)
        })

        // Add URL parameters
        wrapper.vm.params = [
            { key: 'page', value: '1' },
            { key: 'limit', value: '10' }
        ]
        wrapper.vm.requestUrl = 'http://api.example.com/users'
        await wrapper.vm.$nextTick()

        // Send request
        const sendButton = wrapper.find('button.bg-blue-600')
        await sendButton.trigger('click')

        // Verify URL includes parameters
        expect(global.fetch).toHaveBeenCalledWith(
            'http://localhost:9090',
            expect.objectContaining({
                body: expect.stringContaining('http://api.example.com/users?page=1&limit=10')
            })
        )
    })

    it('handles FormData requests correctly', async () => {
        const mockResponse = {
            statusCode: 200,
            body: btoa('{"success":true}'),
            headers: { 'Content-Type': 'application/json' }
        }

        // Mock fetch
        global.fetch = vi.fn().mockResolvedValue({
            json: () => Promise.resolve(mockResponse)
        })

        // Set up FormData request
        wrapper.vm.selectedBodyType = 'FormData'
        wrapper.vm.bodyContent = 'name=John Doe\nemail=john@example.com'
        wrapper.vm.requestUrl = 'http://api.example.com/submit'
        await wrapper.vm.$nextTick()

        // Send request
        const sendButton = wrapper.find('button.bg-blue-600')
        await sendButton.trigger('click')

        // Verify request was sent with correct content type and body
        expect(global.fetch).toHaveBeenCalledWith(
            'http://localhost:9090',
            expect.objectContaining({
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: expect.stringContaining('application/x-www-form-urlencoded')
            })
        )
    })

    it('handles invalid URLs gracefully', async () => {
        // Set invalid URL
        wrapper.vm.requestUrl = 'not-a-valid-url'
        await wrapper.vm.$nextTick()

        // Send request
        const sendButton = wrapper.find('button.bg-blue-600')
        await sendButton.trigger('click')

        // Verify error handling
        expect(wrapper.vm.responseData).toBeTruthy()
        expect(wrapper.vm.responseData.status).toBe(400)
        expect(wrapper.vm.responseData.body).toContain('Invalid URL')
    })

    it('updates placeholder when body type changes', async () => {
        // Change to Text type first (which has no placeholder)
        wrapper.vm.selectedBodyType = 'Text'
        wrapper.vm.bodyContent = ''
        await wrapper.vm.$nextTick()

        // Change to JSON type and verify placeholder
        wrapper.vm.selectedBodyType = 'JSON'
        await wrapper.vm.$nextTick()
        expect(wrapper.vm.bodyContent).toContain('"key": "value"')

        // Change to FormData type and verify placeholder
        wrapper.vm.bodyContent = ''
        await wrapper.vm.$nextTick()
        wrapper.vm.selectedBodyType = 'FormData'
        await wrapper.vm.$nextTick()
        expect(wrapper.vm.bodyContent).toContain('name=John Doe')
    })
}) 