<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-800 mb-8">Create New Project</h1>
    
    <form @submit.prevent="createProject" class="max-w-2xl mx-auto bg-white rounded-lg shadow-md p-6">
      <div class="mb-6">
        <label for="title" class="block text-sm font-medium text-gray-700 mb-2">Project Title</label>
        <input
          v-model="form.title"
          type="text"
          id="title"
          required
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Enter project title"
        >
      </div>

      <div class="mb-6">
        <label for="description" class="block text-sm font-medium text-gray-700 mb-2">Description</label>
        <textarea
          v-model="form.description"
          id="description"
          required
          rows="4"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Describe your project"
        ></textarea>
      </div>

      <div class="mb-6">
        <label for="fundingGoal" class="block text-sm font-medium text-gray-700 mb-2">Funding Goal ($)</label>
        <input
          v-model.number="form.fundingGoal"
          type="number"
          id="fundingGoal"
          required
          min="1"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Enter funding goal"
        >
      </div>

      <div class="mb-6">
        <label for="deadline" class="block text-sm font-medium text-gray-700 mb-2">Deadline</label>
        <input
          v-model="form.deadline"
          type="date"
          id="deadline"
          required
          :min="minDate"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
      </div>

      <div class="flex justify-end gap-4">
        <router-link
          to="/"
          class="px-4 py-2 text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
        >
          Cancel
        </router-link>
        <button
          type="submit"
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
          :disabled="loading"
        >
          {{ loading ? 'Creating...' : 'Create Project' }}
        </button>
      </div>

      <p v-if="error" class="mt-4 text-red-600">{{ error }}</p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import type { Project } from '@/types/project';

const router = useRouter();
const loading = ref(false);
const error = ref('');

const form = ref({
  title: '',
  description: '',
  fundingGoal: 0,
  deadline: ''
});

const minDate = computed(() => {
  const today = new Date();
  today.setDate(today.getDate() + 1); // Minimum deadline is tomorrow
  return today.toISOString().split('T')[0];
});

const createProject = async () => {
  loading.value = true;
  error.value = '';

  try {
    const response = await axios.post<Project>('/api/projects', form.value);
    router.push(`/projects/${response.data.id}`);
  } catch (err) {
    error.value = 'Failed to create project. Please try again.';
    console.error('Project creation error:', err);
  } finally {
    loading.value = false;
  }
};
</script> 