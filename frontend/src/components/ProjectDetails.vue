<template>
  <div class="space-y-8">
    <!-- Project Header -->
    <div class="bg-white rounded-lg shadow-md p-6">
      <div class="flex justify-between items-start">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 mb-2">{{ project.title }}</h1>
          <p class="text-gray-500">Created by {{ project.creatorId }}</p>
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

      <div class="mt-6">
        <p class="text-gray-700 whitespace-pre-line">{{ project.description }}</p>
      </div>

      <!-- Funding Progress -->
      <div class="mt-8">
        <div class="flex justify-between text-sm text-gray-600 mb-2">
          <span>${{ project.currentFunding.toLocaleString() }} raised</span>
          <span>of ${{ project.fundingGoal.toLocaleString() }} goal</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div 
            class="bg-blue-600 h-2 rounded-full" 
            :style="{ width: `${Math.min((project.currentFunding / project.fundingGoal) * 100, 100)}%` }"
          ></div>
        </div>
        <div class="flex justify-between items-center mt-2">
          <span class="text-sm text-gray-600">
            {{ Math.round((project.currentFunding / project.fundingGoal) * 100) }}% funded
          </span>
          <span class="text-sm text-gray-600">
            {{ daysLeft }} days left
          </span>
        </div>
      </div>
    </div>

    <!-- Donation Section -->
    <div v-if="project.status === 'active'" class="bg-white rounded-lg shadow-md p-6">
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Support this project</h2>
      <DonationForm :project-id="project.id" @donation-success="onDonationSuccess" />
    </div>

    <!-- Project is not active message -->
    <div 
      v-else 
      class="bg-gray-50 border border-gray-200 rounded-lg p-6 text-center"
    >
      <p class="text-gray-600">
        {{ 
          project.status === 'completed' 
            ? 'This project has been successfully funded!' 
            : 'This project is no longer accepting donations.'
        }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import DonationForm from './DonationForm.vue';
import type { Project } from '@/types/project';

interface Props {
  project: Project;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: 'donation-success'): void;
}>();

const daysLeft = computed(() => {
  const deadline = new Date(props.project.deadline);
  const now = new Date();
  const diff = deadline.getTime() - now.getTime();
  return Math.max(0, Math.ceil(diff / (1000 * 60 * 60 * 24)));
});

const onDonationSuccess = () => {
  emit('donation-success');
};
</script>