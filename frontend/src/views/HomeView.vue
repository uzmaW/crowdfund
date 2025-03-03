<template>
  <div>
    <!-- Hero Section with Background Image -->
    <section class="relative bg-gradient-to-r from-blue-600 to-blue-800 text-white overflow-hidden">
      <!-- Background Pattern -->
      <div class="absolute inset-0 opacity-10">
        <div class="absolute inset-0" style="background-image: url('data:image/svg+xml,%3Csvg width=\'60\' height=\'60\' viewBox=\'0 0 60 60\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cg fill=\'none\' fill-rule=\'evenodd\'%3E%3Cg fill=\'%23ffffff\' fill-opacity=\'1\'%3E%3Cpath d=\'M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z\'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E')"></div>
      </div>
      
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="py-20 md:py-28">
          <h1 class="text-4xl md:text-6xl font-bold mb-6 leading-tight">
            Fund Your Dreams<br/>
            <span class="text-blue-200">Support Amazing Projects</span>
          </h1>
          <p class="text-xl md:text-2xl text-blue-100 mb-8 max-w-2xl">
            Join our community of innovators and backers. Discover projects that are changing the world or start your own campaign today.
          </p>
          <div class="flex flex-wrap gap-4">
            <router-link 
              to="/projects/create" 
              class="btn-primary text-lg px-8 py-3"
            >
              Start Your Project
            </router-link>
            <router-link 
              to="/projects" 
              class="btn bg-white text-blue-600 hover:bg-blue-50 text-lg px-8 py-3"
            >
              Explore Projects
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Featured Projects Section -->
    <section class="py-16 bg-gray-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">Featured Projects</h2>
          <p class="text-xl text-gray-600 max-w-2xl mx-auto">
            Discover innovative projects from creators around the world
          </p>
        </div>

        <!-- Loading State -->
        <div v-if="projectStore.loadingProjects" class="flex justify-center items-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-4 border-blue-600 border-t-transparent"></div>
        </div>

        <!-- Projects Grid -->
        <div v-else-if="projectStore.projects.length > 0" 
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8"
        >
          <div 
            v-for="project in projectStore.projects" 
            :key="project.id" 
            class="card group hover:shadow-lg transition-shadow duration-300"
          >
            <!-- Project Image/Placeholder -->
            <div class="h-48 bg-gradient-to-br from-blue-500 to-blue-600 relative overflow-hidden">
              <div class="absolute inset-0 bg-black bg-opacity-20"></div>
              <div class="absolute bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-black/50">
                <h3 class="text-xl font-bold text-white truncate">{{ project.title }}</h3>
              </div>
            </div>
            
            <div class="p-6">
              <p class="text-gray-600 mb-4 line-clamp-2">{{ project.description }}</p>
              
              <!-- Progress Bar -->
              <div class="mb-4">
                <div class="flex justify-between text-sm text-gray-600 mb-1">
                  <span class="font-medium">${{ project.currentFunding.toLocaleString() }} raised</span>
                  <span>{{ Math.round((project.currentFunding / project.fundingGoal) * 100) }}%</span>
                </div>
                <div class="progress-bar">
                  <div 
                    class="progress-bar-fill"
                    :style="{ 
                      width: `${Math.min((project.currentFunding / project.fundingGoal) * 100, 100)}%`,
                      background: 'linear-gradient(90deg, var(--color-primary-light) 0%, var(--color-primary) 100%)'
                    }"
                  ></div>
                </div>
              </div>

              <div class="flex items-center justify-between">
                <span 
                  :class="{
                    'badge-success': project.status === 'active',
                    'badge-warning': project.status === 'completed',
                    'badge-danger': project.status === 'canceled'
                  }"
                  class="badge"
                >
                  {{ project.status.charAt(0).toUpperCase() + project.status.slice(1) }}
                </span>
                <router-link 
                  :to="`/projects/${project.id}`" 
                  class="text-blue-600 hover:text-blue-700 font-medium group-hover:translate-x-1 transition-transform duration-200"
                >
                  View Details →
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="text-center py-12">
          <div class="max-w-md mx-auto">
            <svg class="mx-auto h-16 w-16 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
            </svg>
            <h3 class="mt-4 text-lg font-medium text-gray-900">No Projects Yet</h3>
            <p class="mt-2 text-gray-600">Be the first to create an amazing project and start your crowdfunding journey.</p>
            <div class="mt-6">
              <router-link 
                to="/projects/create" 
                class="btn-primary"
              >
                Create Project
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="py-16 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">How It Works</h2>
          <p class="text-xl text-gray-600 max-w-2xl mx-auto">
            Start your journey with Crowdfund in three simple steps
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <!-- Step 1 -->
          <div class="text-center">
            <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2">Create</h3>
            <p class="text-gray-600">Share your story and set your funding goal</p>
          </div>

          <!-- Step 2 -->
          <div class="text-center">
            <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2">Fund</h3>
            <p class="text-gray-600">Get support from backers around the world</p>
          </div>

          <!-- Step 3 -->
          <div class="text-center">
            <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold mb-2">Achieve</h3>
            <p class="text-gray-600">Make your project a reality</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useProjectStore } from '../store/project';

const projectStore = useProjectStore();

onMounted(() => {
  projectStore.fetchProjects();
});
</script>