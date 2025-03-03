<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="container mx-auto px-4">
      <div class="max-w-4xl mx-auto">
        <!-- User Profile Section -->
        <div class="bg-white rounded-lg shadow-md p-6 mb-8">
          <UserProfile />
        </div>

        <!-- User's Projects Section -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-2xl font-semibold text-gray-900">My Projects</h2>
            <router-link
              to="/projects/create"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            >
              Create New Project
            </router-link>
          </div>

          <!-- Loading State -->
          <div v-if="loading" class="flex justify-center items-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>

          <!-- Error State -->
          <div v-else-if="error" class="text-center py-12">
            <p class="text-red-600">{{ error }}</p>
            <button 
              @click="fetchUserProjects" 
              class="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Retry
            </button>
          </div>

          <!-- Projects List -->
          <div v-else-if="projects.length > 0" class="space-y-4">
            <div 
              v-for="project in projects" 
              :key="project.id"
              class="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors"
            >
              <div class="flex justify-between items-start">
                <div>
                  <h3 class="text-lg font-medium text-gray-900">
                    <router-link 
                      :to="`/projects/${project.id}`"
                      class="hover:text-blue-600"
                    >
                      {{ project.title }}
                    </router-link>
                  </h3>
                  <p class="text-gray-600 mt-1">{{ project.description.substring(0, 100) }}...</p>
                </div>
                <span 
                  class="inline-block px-3 py-1 rounded-full text-sm"
                  :class="{
                    'bg-green-100 text-green-800': project.status === 'active',
                    'bg-blue-100 text-blue-800': project.status === 'completed',
                    'bg-red-100 text-red-800': project.status === 'canceled'
                  }"
                >
                  {{ project.status.charAt(0).toUpperCase() + project.status.slice(1) }}
                </span>
              </div>

              <div class="mt-4">
                <div class="flex justify-between text-sm text-gray-600 mb-1">
                  <span>${{ project.currentFunding.toLocaleString() }} raised</span>
                  <span>of ${{ project.fundingGoal.toLocaleString() }} goal</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    class="bg-blue-600 h-2 rounded-full" 
                    :style="{ width: `${Math.min((project.currentFunding / project.fundingGoal) * 100, 100)}%` }"
                  ></div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="text-center py-12">
            <p class="text-gray-600 mb-4">You haven't created any projects yet.</p>
            <router-link
              to="/projects/create"
              class="inline-block px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            >
              Create Your First Project
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useUserStore } from '@/store/user';
import UserProfile from '@/components/UserProfile.vue';
import axios from 'axios';
import type { Project } from '@/types/project';

const userStore = useUserStore();
const loading = ref(true);
const error = ref('');
const projects = ref<Project[]>([]);

const fetchUserProjects = async () => {
  loading.value = true;
  error.value = '';

  try {
    const response = await axios.get<Project[]>('/api/projects/user', {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
    projects.value = response.data;
  } catch (err) {
    error.value = 'Failed to load your projects. Please try again.';
    console.error('User projects loading error:', err);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchUserProjects);
</script> 