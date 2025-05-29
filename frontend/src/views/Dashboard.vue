<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, computed } from "vue";
import MonacoEditor from 'monaco-editor-vue3'
import { useCollectionsStore } from '../stores/collections'
import Collections from '../components/Collections.vue'
import Header from '../components/Header.vue'

// Types
type BodyType = "JSON" | "Text" | "FormData";
type KeyValuePair = { key: string; value: string };
type Method = "POST" | "GET" | "PATCH" | "PUT" | "DELETE" | "HEAD" | "OPTIONS";

interface Request {
  method: Method;
  url: string;
  headers: Record<string, string>;
  body: number[];
}

interface Response {
  body: string;
  headers: KeyValuePair[];
  status: number;
}

// Constants
const OPTIONS = ["Params", "Headers", "Body"] as const;
const BODY_TYPES = ["JSON", "Text", "FormData"] as const;
const METHODS = ["POST", "GET", "PATCH", "PUT", "DELETE", "HEAD", "OPTIONS"] as const;

// State
const requestUrl = ref<string>("");
const bodyContent = ref<string>("");
const params = ref<KeyValuePair[]>([]);
const headers = ref<KeyValuePair[]>([]);
const selectedOption = ref<string>(OPTIONS[0]);
const selectedMethod = ref<Method>(METHODS[0] as Method);
const selectedBodyType = ref<BodyType>("JSON");
const paramPairs = ref<KeyValuePair>({ key: "", value: "" });
const responseData = ref<Response | null>(null);
const isLoading = ref(false);
const showSaveModal = ref(false)
const selectedCollectionId = ref('')

const collectionsStore = useCollectionsStore()

// Watch for current request changes
watch(() => collectionsStore.currentRequest, (request) => {
  if (request) {
    requestUrl.value = request.url
    selectedMethod.value = request.method as Method
    selectedBodyType.value = request.bodyType as BodyType || "JSON"
    bodyContent.value = request.body || getEditorPlaceholder()
    headers.value = Object.entries(request.headers).map(([key, value]) => ({ key, value }))
  } else {
    selectedBodyType.value = "JSON"
    bodyContent.value = getEditorPlaceholder()
  }
}, { immediate: true })

// Set initial placeholder on mount
onMounted(() => {
  if (!bodyContent.value.trim()) {
    bodyContent.value = getEditorPlaceholder();
  }
});

const addParam = () => {
  if (paramPairs.value.key.trim() === "") return;
  params.value.push({ ...paramPairs.value });
  paramPairs.value = { key: "", value: "" };
};

const addHeader = () => {
  if (paramPairs.value.key.trim() === "") return;
  headers.value.push({ ...paramPairs.value });
  paramPairs.value = { key: "", value: "" };
};

const removeParam = (index: number) => {
  params.value.splice(index, 1);
};

const removeHeader = (index: number) => {
  headers.value.splice(index, 1);
};

const formatJSON = () => {
  if (selectedBodyType.value === 'JSON' && bodyContent.value) {
    try {
      const parsed = JSON.parse(bodyContent.value);
      bodyContent.value = JSON.stringify(parsed, null, 2);
    } catch (e) {
      console.error('Invalid JSON');
    }
  }
};

const decodeBase64 = (base64String: string): string => {
  try {
    return atob(base64String);
  } catch (e) {
    console.error('Failed to decode base64:', e);
    return base64String;
  }
};

const getEditorPlaceholder = () => {
  switch (selectedBodyType.value) {
    case 'JSON':
      return `{
  "key": "value",
  "nested": {
    "array": [1, 2, 3],
    "string": "text"
  }
}`;
    case 'FormData':
      return `# You can use either JSON format:
{
  "name": "John Doe",
  "email": "john@example.com"
}

# Or key=value format:
name=John Doe
email=john@example.com`;
    default:
      return ''; // Plain text has no placeholder
  }
};

// Watch for body type changes to update placeholder
watch(selectedBodyType, (newType) => {
  if (!bodyContent.value.trim()) {
    bodyContent.value = getEditorPlaceholder();
  }
});

const saveCurrentRequest = async () => {
  if (!selectedCollectionId.value) {
    alert('Please select a collection')
    return
  }

  const request = {
    name: `${selectedMethod.value} ${requestUrl.value}`,
    method: selectedMethod.value,
    url: requestUrl.value,
    headers: headers.value.reduce((acc, { key, value }) => {
      if (key && value) acc[key] = value
      return acc
    }, {} as Record<string, string>),
    body: bodyContent.value,
    bodyType: selectedBodyType.value,
  }

  try {
    await collectionsStore.addRequest(selectedCollectionId.value, request)
    showSaveModal.value = false
    selectedCollectionId.value = ''
    alert('Request saved successfully')
  } catch (error) {
    console.error('Failed to save request:', error)
    alert('Failed to save request')
  }
}

const sendRequest = async () => {
  isLoading.value = true;
  responseData.value = null;

  // Validate URL
  try {
    // Build URL with parameters
    let url = new URL(requestUrl.value);
    params.value.forEach((param) => {
      url.searchParams.append(param.key, param.value);
    });

    // Convert headers array to object
    const headersMap: Record<string, string> = {};
    headers.value.forEach((h) => {
      headersMap[h.key] = h.value;
    });

    // Prepare request body based on type
    let bodyBytes: Uint8Array;
    if (selectedBodyType.value === 'FormData') {
      const formData = new FormData();
      try {
        // Try to parse as JSON first to support nested form data
        const formDataObj = JSON.parse(bodyContent.value);
        Object.entries(formDataObj).forEach(([key, value]) => {
          formData.append(key, String(value));
        });
      } catch {
        // If not valid JSON, try to parse as key=value pairs
        bodyContent.value.split('\n').forEach(line => {
          const [key, value] = line.split('=').map(s => s.trim());
          if (key && value) {
            formData.append(key, value);
          }
        });
      }
      // Convert FormData to URLSearchParams
      const params = new URLSearchParams();
      formData.forEach((value, key) => {
        params.append(key, value.toString());
      });
      bodyBytes = new TextEncoder().encode(params.toString());
      // Add proper content type header
      headersMap['Content-Type'] = 'application/x-www-form-urlencoded';
    } else {
      bodyBytes = new TextEncoder().encode(bodyContent.value);
      if (selectedBodyType.value === 'JSON' && !headersMap['Content-Type']) {
        headersMap['Content-Type'] = 'application/json';
      }
    }

    const request: Request = {
      method: selectedMethod.value,
      url: url.toString(),
      headers: headersMap,
      body: Array.from(bodyBytes)
    };

    const proxyUrl = "http://localhost:9090";

    try {
      const response = await fetch(proxyUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
      });

      const data = await response.json();

      if (data.message === "fail") {
        throw new Error("Request failed");
      }

      // Handle the proxied response
      let bodyContent = "";
      if (data.body) {
        // First decode the base64 string
        const decodedBody = decodeBase64(data.body);
        try {
          // Check if the decoded content is JSON
          if (decodedBody.trim().startsWith('{') || decodedBody.trim().startsWith('[')) {
            const parsed = JSON.parse(decodedBody);
            bodyContent = JSON.stringify(parsed, null, 2);
          } else {
            bodyContent = decodedBody;
          }
        } catch {
          // If JSON parsing fails, use the decoded content as is
          bodyContent = decodedBody;
        }
      }

      responseData.value = {
        body: bodyContent,
        headers: Object.entries(data.headers || {}).map(([key, value]) => ({ key, value: String(value) })),
        status: data.statusCode
      };
    } catch (error) {
      console.error("Request failed:", error);
      responseData.value = {
        body: "Request failed: " + (error instanceof Error ? error.message : String(error)),
        headers: [],
        status: 500
      };
    }
  } catch (error) {
    console.error("Invalid URL:", error);
    responseData.value = {
      body: "Invalid URL: Please enter a valid URL including the protocol (e.g., http:// or https://)",
      headers: [],
      status: 400
    };
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <div class="flex flex-col h-screen bg-gray-50">
    <!-- Top Bar -->
    <div class="h-16 bg-white border-b border-gray-200">
      <div class="h-full px-6">
        <div class="flex items-center h-full space-x-4">
          <select v-model="selectedMethod"
            class="h-9 w-28 px-3 text-sm font-medium bg-gray-50 border border-gray-200 rounded-md focus:outline-none focus:ring-1 focus:ring-blue-500">
            <option v-for="method in METHODS" :key="method" :value="method">{{ method }}</option>
          </select>
          <div class="flex-1">
            <input v-model="requestUrl" type="text" placeholder="Enter request URL"
              class="w-full h-9 px-3 text-sm bg-gray-50 border border-gray-200 rounded-md focus:outline-none focus:ring-1 focus:ring-blue-500" />
          </div>

          <div class="flex items-center space-x-3">
            <button @click="showSaveModal = true"
              class="h-9 px-3 text-gray-600 hover:text-gray-800 bg-gray-50 border border-gray-200 rounded-md hover:bg-gray-100 transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                <path
                  d="M7.707 10.293a1 1 0 10-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 11.586V6h1a2 2 0 012 2v7a2 2 0 01-2 2H7a2 2 0 01-2-2V8a2 2 0 012-2h1v5.586l-1.293-1.293zM9 5a1 1 0 012 0v3a1 1 0 11-2 0V5z" />
              </svg>
            </button>

            <button @click="sendRequest" :disabled="isLoading"
              class="h-9 px-4 text-sm font-medium bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
              <span v-if="isLoading">Sending...</span>
              <span v-else>Send</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 overflow-y-auto">
      <div class="h-full">
        <!-- Options Tabs -->
        <div class="bg-white border-b border-gray-200">
          <div class="flex">
            <button v-for="option in OPTIONS" :key="option" @click="selectedOption = option"
              class="px-6 py-3 text-sm font-medium border-b-2 transition-colors" :class="{
                'text-blue-600 border-blue-600': selectedOption === option,
                'text-gray-500 hover:text-gray-700 border-transparent': selectedOption !== option
              }">
              {{ option }}
            </button>
          </div>

          <div class="px-6 py-6">
            <!-- Params Section -->
            <div v-if="selectedOption === OPTIONS[0]" class="space-y-4">
              <div class="flex gap-4">
                <input type="text" v-model="paramPairs.key" placeholder="Parameter Key"
                  class="flex-1 h-9 px-3 text-sm bg-white border border-gray-300 rounded-md focus:ring-1 focus:ring-blue-500 focus:border-blue-500" />
                <input type="text" v-model="paramPairs.value" placeholder="Parameter Value"
                  class="flex-1 h-9 px-3 text-sm bg-white border border-gray-300 rounded-md focus:ring-1 focus:ring-blue-500 focus:border-blue-500" />
                <button @click="addParam"
                  class="px-4 h-9 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:ring-2 focus:ring-gray-200">
                  Add
                </button>
              </div>

              <!-- Parameters List -->
              <div class="space-y-2">
                <div v-for="(pair, index) in params" :key="index"
                  class="flex items-center gap-4 px-4 py-2 bg-gray-50 rounded-md">
                  <span class="flex-1 font-medium text-sm text-gray-700">{{ pair.key }}</span>
                  <span class="flex-1 text-sm text-gray-600">{{ pair.value }}</span>
                  <button @click="removeParam(index)" class="text-gray-400 hover:text-red-500 transition-colors">
                    <i class="pi pi-times"></i>
                  </button>
                </div>
              </div>
            </div>

            <!-- Headers Section -->
            <div v-if="selectedOption === OPTIONS[1]" class="space-y-4">
              <div class="flex gap-4">
                <input type="text" v-model="paramPairs.key" placeholder="Header Key"
                  class="flex-1 h-9 px-3 text-sm bg-white border border-gray-300 rounded-md focus:ring-1 focus:ring-blue-500 focus:border-blue-500" />
                <input type="text" v-model="paramPairs.value" placeholder="Header Value"
                  class="flex-1 h-9 px-3 text-sm bg-white border border-gray-300 rounded-md focus:ring-1 focus:ring-blue-500 focus:border-blue-500" />
                <button @click="addHeader"
                  class="px-4 h-9 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:ring-2 focus:ring-gray-200">
                  Add
                </button>
              </div>

              <!-- Headers List -->
              <div class="space-y-2">
                <div v-for="(pair, index) in headers" :key="index"
                  class="flex items-center gap-4 px-4 py-2 bg-gray-50 rounded-md">
                  <span class="flex-1 font-medium text-sm text-gray-700">{{ pair.key }}</span>
                  <span class="flex-1 text-sm text-gray-600">{{ pair.value }}</span>
                  <button @click="removeHeader(index)" class="text-gray-400 hover:text-red-500 transition-colors">
                    <i class="pi pi-times"></i>
                  </button>
                </div>
              </div>
            </div>

            <!-- Body Section -->
            <div v-if="selectedOption === OPTIONS[2]" class="space-y-4">
              <div class="flex items-center gap-4">
                <select v-model="selectedBodyType"
                  class="h-9 px-3 text-sm bg-white border border-gray-300 rounded-md focus:ring-1 focus:ring-blue-500 focus:border-blue-500">
                  <option v-for="type in BODY_TYPES" :key="type" :value="type">
                    {{ type }}
                  </option>
                </select>
                <button v-if="selectedBodyType === 'JSON'" @click="formatJSON"
                  class="h-9 px-4 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:ring-2 focus:ring-gray-200">
                  Format JSON
                </button>
              </div>
              <MonacoEditor v-model="bodyContent" :options="{
                theme: 'vs-dark',
                language: selectedBodyType === 'JSON' ? 'json' : 'plaintext',
                automaticLayout: true,
                minimap: { enabled: false },
                lineNumbers: 'on',
                scrollBeyondLastLine: false,
                fontSize: 13,
                tabSize: 2,
                wordWrap: 'on',
                wrappingIndent: 'indent',
                padding: { top: 12, bottom: 12 }
              }" class="border border-gray-300 rounded-lg overflow-hidden" height="400" />
            </div>
          </div>
        </div>

        <!-- Response Section -->
        <div v-if="responseData" class="bg-white border-t border-gray-200">
          <div class="px-6 py-3 border-b border-gray-200">
            <div class="flex items-center space-x-3">
              <h2 class="text-sm font-medium text-gray-900">Response</h2>
              <span :class="{
                'text-green-600': responseData.status >= 200 && responseData.status < 300,
                'text-yellow-600': responseData.status >= 300 && responseData.status < 400,
                'text-red-600': responseData.status >= 400
              }" class="text-sm font-medium">
                Status: {{ responseData.status }}
              </span>
            </div>
          </div>

          <div class="px-6 py-6">
            <!-- Response Headers -->
            <div v-if="responseData.headers.length > 0">
              <h3 class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-3">Headers</h3>
              <div class="space-y-2">
                <div v-for="(header, index) in responseData.headers" :key="index"
                  class="flex items-center px-4 py-2 bg-gray-50 rounded-md">
                  <span class="flex-1 text-sm font-medium text-gray-700">{{ header.key }}</span>
                  <span class="flex-1 text-sm text-gray-600">{{ header.value }}</span>
                </div>
              </div>
            </div>

            <!-- Response Body -->
            <div>
              <h3 class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-3">Body</h3>
              <pre
                class="p-4 bg-gray-50 rounded-md font-mono text-sm text-gray-700 overflow-x-auto">{{ responseData.body }}</pre>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Save Request Modal -->
    <Transition enter-active-class="ease-out duration-200" enter-from-class="opacity-0" enter-to-class="opacity-100"
      leave-active-class="ease-in duration-150" leave-from-class="opacity-100" leave-to-class="opacity-0">
      <div v-if="showSaveModal" class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
    </Transition>

    <Transition enter-active-class="ease-out duration-300"
      enter-from-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
      enter-to-class="opacity-100 translate-y-0 sm:scale-100" leave-active-class="ease-in duration-200"
      leave-from-class="opacity-100 translate-y-0 sm:scale-100"
      leave-to-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95">
      <div v-if="showSaveModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-0">
        <div
          class="relative transform overflow-hidden rounded-lg bg-white shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-lg">
          <div class="absolute right-0 top-0 pr-4 pt-4">
            <button @click="showSaveModal = false"
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
            <h3 class="text-xl font-semibold text-gray-900 mb-6">Save Request to Collection</h3>
            <div class="space-y-5">
              <div>
                <label for="collection-select" class="block text-sm font-medium text-gray-700 mb-1">Collection</label>
                <select id="collection-select" v-model="selectedCollectionId"
                  class="block w-full rounded-md border-0 py-2.5 px-3 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm">
                  <option value="">Choose a collection...</option>
                  <option v-for="collection in collectionsStore.collections" :key="collection.id"
                    :value="collection.id">
                    {{ collection.name }}
                  </option>
                </select>
              </div>
              <div>
                <label for="request-name" class="block text-sm font-medium text-gray-700 mb-1">Request Name</label>
                <div class="mt-1 flex rounded-md shadow-sm ring-1 ring-inset ring-gray-300">
                  <span
                    class="flex select-none items-center px-3 text-gray-500 sm:text-sm bg-gray-50 rounded-l-md border-r border-gray-300">{{
                      selectedMethod }}</span>
                  <input type="text" id="request-name" :value="`${requestUrl}`" readonly
                    class="block flex-1 border-0 bg-transparent py-2.5 px-3 text-gray-900 placeholder:text-gray-400 focus:ring-0 sm:text-sm" />
                </div>
                <p class="mt-1 text-sm text-gray-500">The request will be saved with this name based on the URL</p>
              </div>
            </div>
          </div>

          <div class="bg-gray-50 px-6 py-4 flex items-center justify-end space-x-3">
            <button @click="showSaveModal = false"
              class="inline-flex justify-center rounded-md bg-white px-4 py-2.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50">
              Cancel
            </button>
            <button @click="saveCurrentRequest"
              class="inline-flex justify-center rounded-md bg-blue-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
              :disabled="!selectedCollectionId">
              Save Request
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style>
button {
  cursor: pointer;
}

button:disabled {
  cursor: not-allowed;
}

.monaco-editor {
  padding-top: 8px;
}

.monaco-editor .overflow-guard {
  border-radius: 8px;
}
</style>