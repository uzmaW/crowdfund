<template>
  <form @submit.prevent="donate">
    <div class="mb-4">
      <label class="block text-gray-700 text-sm font-bold mb-2" for="amount">Amount</label>
      <input v-model="amount" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="amount" type="number" placeholder="Amount">
    </div>
    <button class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">Donate</button>
  </form>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios';
import { useRoute } from 'vue-router';

const route = useRoute();
const projectId = route.params.id;
const amount = ref(0);

const donate = async () => {
  try {
    await axios.post(`/api/projects/${projectId}/donations`, { amount });
    // Handle success (e.g., display success message)
  } catch (error) {
    console.error('Donation failed:', error);
    // Handle error (e.g., display error message)
  }
};
</script>