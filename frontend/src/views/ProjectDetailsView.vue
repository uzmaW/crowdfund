<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="container mx-auto px-4">
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-12">
        <p class="text-red-600">{{ error }}</p>
        <button 
          @click="fetchProject" 
          class="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          Retry
        </button>
      </div>

      <!-- Project Details -->
      <div v-else-if="project" class="max-w-4xl mx-auto">
        <ProjectDetails 
          :project="project" 
          @donation-success="fetchProject"
        />
      </div>

      <!-- Not Found State -->
      <div v-else class="text-center py-12">
        <p class="text-gray-600">Project not found.</p>
        <router-link 
          to="/" 
          class="mt-4 inline-block px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          Back to Home
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import ProjectDetails from '@/components/ProjectDetails.vue';
import type { Project } from '@/types/project';

const route = useRoute();
const loading = ref(true);
const error = ref('');
const project = ref<Project | null>(null);

const fetchProject = async () => {
  loading.value = true;
  error.value = '';

  try {
    const response = await axios.get<Project>(`/api/projects/${route.params.id}`);
    project.value = response.data;
  } catch (err) {
    error.value = 'Failed to load project details. Please try again.';
    console.error('Project loading error:', err);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchProject);
</script> 