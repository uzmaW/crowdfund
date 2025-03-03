<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useProjectStore } from '../store/project';
import type { Project } from '../types/project';
import { useDebounceFn } from '@vueuse/core';

const projectStore = useProjectStore();
const searchQuery = ref('');
const currentPage = ref(1);
const itemsPerPage = ref(10);
const error = ref<string | null>(null);

interface Header {
  text: string;
  value: string;
  sortable?: boolean;
  width?: number;
}

const headers: Header[] = [
  { text: 'Title', value: 'title', sortable: true, width: 200 },
  { text: 'Goal', value: 'fundingGoal', sortable: true },
  { text: 'Progress', value: 'progressPercentage', sortable: true },
  { text: 'Deadline', value: 'deadline', sortable: true },
  { text: 'Status', value: 'status', sortable: true }
];

const loadProjects = useDebounceFn(async () => {
  try {
    error.value = null;
    await projectStore.fetchProjects({
      search: searchQuery.value,
      page: currentPage.value,
      limit: itemsPerPage.value
    });
  } catch (err) {
    error.value = 'Failed to load projects. Please try again later.';
    console.error('Project loading error:', err);
  }
}, 300);

onMounted(() => {
  loadProjects();
});

const projects = computed<Project[]>(() => projectStore.projects);
const totalItems = computed(() => projectStore.totalProjects);
const loading = computed(() => projectStore.loadingProjects);
</script>

<template>
  <div class="project-list-container">
    <div class="project-list-controls">
      <h2 class="project-list-title">Crowdfunding Projects</h2>
      <div class="project-list-actions">
        <input
          v-model="searchQuery"
          type="search"
          placeholder="Search projects..."
          class="project-search-input"
          @input="loadProjects"
          aria-label="Search projects"
        >
        <button
          class="refresh-button"
          @click="loadProjects"
          aria-label="Refresh projects"
        >
          ↻
        </button>
      </div>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
      <button @click="loadProjects" class="retry-button">Retry</button>
    </div>

    <EasyDataTable
      :headers="headers"
      :items="projects"
      :loading="loading"
      :table-class-name="'project-table'"
      :rows-items="[10, 20, 50]"
      :rows-per-page="itemsPerPage"
      :current-page="currentPage"
      :server-current-page="currentPage"
      :server-items-length="totalItems"
      @update:server-options="(options) => {
        currentPage = options.page;
        itemsPerPage = options.rowsPerPage;
        loadProjects();
      }"
    >
      <template #loading>
        <div class="loading-skeleton">
          <div v-for="i in 5" :key="i" class="skeleton-row"></div>
        </div>
      </template>

      <template #item-title="{ title, item }">
        <router-link 
          :to="`/projects/${item.id}`" 
          class="project-title-link"
          aria-label="View project details"
        >
          {{ title }}
        </router-link>
      </template>

      <template #item-progressPercentage="{ progressPercentage }">
        <div class="progress-bar-container">
          <div 
            class="progress-bar" 
            :style="{ width: `${progressPercentage}%` }"
            role="progressbar"
            :aria-valuenow="progressPercentage"
            aria-valuemin="0"
            aria-valuemax="100"
          ></div>
          <span class="progress-text">{{ progressPercentage }}%</span>
        </div>
      </template>

      <template #item-status="{ status }">
        <span :class="`status-badge status-${status.toLowerCase()}`">
          {{ status }}
        </span>
      </template>

      <template #empty-message>
        <div class="empty-state">
          <p>No projects found matching your criteria.</p>
          <button @click="() => {
            searchQuery = '';
            loadProjects();
          }" class="clear-filters-button">
            Clear Filters
          </button>
        </div>
      </template>
    </EasyDataTable>
  </div>
</template>

<style scoped>
.project-list-container {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.project-list-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.project-list-actions {
  display: flex;
  gap: 1rem;
}

.project-search-input {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 300px;
}

.refresh-button {
  padding: 0.5rem 1rem;
  background: #f0f0f0;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.2s;
}

.refresh-button:hover {
  background: #e0e0e0;
}

.error-message {
  color: #dc3545;
  background: #f8d7da;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.loading-skeleton {
  padding: 1rem;
}

.skeleton-row {
  height: 50px;
  background: #f0f0f0;
  margin-bottom: 0.5rem;
  border-radius: 4px;
  animation: pulse 1.5s infinite;
}

.progress-bar-container {
  position: relative;
  height: 24px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: #4CAF50;
  transition: width 0.3s ease;
}

.progress-text {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-weight: bold;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.875rem;
  display: inline-block;
}

.status-active {
  background: #4CAF50;
  color: white;
}

.status-completed {
  background: #2196F3;
  color: white;
}

.status-canceled {
  background: #f44336;
  color: white;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #666;
}

@keyframes pulse {
  0% { opacity: 1 }
  50% { opacity: 0.5 }
  100% { opacity: 1 }
}
</style>