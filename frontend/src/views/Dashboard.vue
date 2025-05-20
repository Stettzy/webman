<script setup lang="ts">
import { ref } from "vue";

// Types
type BodyType = "JSON" | "Text" | "FormData";
type KeyValuePair = { key: string, value: string };
type Method = "POST" | "GET" | "PATCH" | "PUT" | "DELETE" | "HEAD" | "OPTIONS";

interface Request {
  method: Method;
  requestUrl: string;
  params: KeyValuePair[];
  headers: KeyValuePair[];
  body: {
    type: BodyType;
    content: string;
  }
}

interface Response {
  body: string;
  headers: KeyValuePair[];
}

// Constants
const OPTIONS = ["Params", "Headers", "Body"];
const BODY_TYPES = ["JSON", "Text", "FormData"];
const METHODS = ["POST", "GET", "PATCH", "PUT", "DELETE", "HEAD", "OPTIONS"];

// State
const requestUrl = ref<string>("");
const bodyContent = ref<string>("");
const params = ref<KeyValuePair[]>([]);
const headers = ref<KeyValuePair[]>([]);
const selectedOption = ref<string>(OPTIONS[0]);
const selectedMethod = ref<Method>(METHODS[0]);
const selectedBodyType = ref<BodyType>(BODY_TYPES[0]);
const paramPairs = ref<KeyValuePair>({ key: "", value: "" });

// Functions
const addParam = () => {
  params.value.push(paramPairs.value);
  // Reset the state
  paramPairs.value = { key: "", value: "" };
}

const addHeader = () => {
  headers.value.push(paramPairs.value);
  // Reset the state
  paramPairs.value = { key: "", value: "" };
}

const sendRequest = async () => {
  const headersMap = {};
  headers.value.forEach(h => {
    headersMap[h.key] = h.value;
  });

  const request: Request = {
    method: selectedMethod.value,
    url: requestUrl.value,
    headers: headersMap,
    body: bodyContent.value,
  }

  const url = "http://localhost:8080";

  try {
    const response = await fetch(url, {
      method: "POST",
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(request)
    }).then(response => response.json())
      .then(data => { console.log(data); });

    if (!response.ok) {
      throw new Error('Response status: ', response.status);
    }

    console.log(response.body);
  } catch (error) {
    console.log("ERROR", error);
  }
}
</script>

<template>
  <div class="space-y-6">
    <!-- Top section -->
    <div class="flex items-center justify-start gap-10 bg-gray-100 p-4">
      <select v-model="selectedMethod" class="border border-gray-900 px-4 py-2">
        <option v-for="method in METHODS"> {{ method }}</option>
      </select>
      <input type="text" v-model="requestUrl" class="bg-gray-200 border border-gray-900 px-4 py-2" />
      <button class="bg-gray-900 text-white px-4 py-2" @click="sendRequest">submit</button>
    </div>
    <!-- Body section -->
    <div class="flex flex-col gap-10 p-4 space-y-6">
      <div>
        <button @click="() => selectedOption = option" v-for="option in OPTIONS"
          :class="{ 'bg-gray-900 text-white': selectedOption === option }" class="px-4 py-2">
          {{ option }}
        </button>
      </div>
      <div class="flex space-x-5">
        <template v-if="selectedOption === OPTIONS[0]">
          <div>
            <label for="key">Key</label>
            <input id="key" type="text" v-model="paramPairs.key" class="border border-gray-900 px-4 py-2" />
          </div>
          <div>
            <label for="value">Value</label>
            <input id="value" type="text" v-model="paramPairs.value" class="border border-gray-900 px-4 py-2" />
          </div>
          <button @click="addParam" class="bg-gray-900 text-white px-4 py-2">Save</button>
        </template>
        <template v-else-if="selectedOption === OPTIONS[1]">
          <div>
            <label for="key">Key</label>
            <input id="key" type="text" v-model="paramPairs.key" class="border border-gray-900 px-4 py-2" />
          </div>
          <div>
            <label for="value">Value</label>
            <input id="value" type="text" v-model="paramPairs.value" class="border border-gray-900 px-4 py-2" />
          </div>
          <button @click="addHeader" class="bg-gray-900 text-white px-4 py-2">Save</button>
        </template>
        <template v-else>
          <div class="flex flex-col space-y-4 w-1/2">
            <label for="body">Body</label>
            <select v-model="selectedBodyType" class="border border-gray-900 px-4 py-2">
              <option v-for="type in BODY_TYPES"> {{ type }}</option>
            </select>
            <textarea id="body" rows="10" class="border border-gray-900 px-4 py-2"></textarea>
          </div>
        </template>

        <template v-if="selectedOption === OPTIONS[0]">
          <div :key="index" v-for="(pair, index) in params">
            <input v-model="pair.key" />
            <input v-model="pair.value" />
          </div>
        </template>

        <template v-else-if="selectedOption === OPTIONS[1]">
          <div :key="index" v-for="(pair, index) in headers">
            <input v-model="pair.key" />
            <input v-model="pair.value" />
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style>
button {
  cursor: pointer;
}
</style>