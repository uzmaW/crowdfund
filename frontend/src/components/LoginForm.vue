<template>
  <form @submit.prevent="login" class="space-y-6">
    <div>
      <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
      <input
        v-model="username"
        type="text"
        id="username"
        required
        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        placeholder="Enter your username"
      >
    </div>

    <div>
      <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
      <input
        v-model="password"
        type="password"
        id="password"
        required
        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        placeholder="Enter your password"
      >
    </div>

    <div>
      <button
        type="submit"
        :disabled="loading"
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
      >
        {{ loading ? 'Logging in...' : 'Login' }}
      </button>
    </div>

    <p v-if="error" class="text-red-600 text-sm text-center">{{ error }}</p>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import axios from 'axios';

const router = useRouter();
const userStore = useUserStore();

const username = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');

const login = async () => {
  loading.value = true;
  error.value = '';

  try {
    const response = await axios.post('/api/users/login', {
      username: username.value,
      password: password.value
    });

    userStore.$patch({
      token: response.data.token,
      user: response.data.user
    });
    router.push('/');
  } catch (err) {
    error.value = 'Invalid username or password';
    console.error('Login error:', err);
  } finally {
    loading.value = false;
  }
};
</script>
