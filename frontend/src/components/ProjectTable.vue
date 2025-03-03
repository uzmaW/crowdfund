<template>
  <div class="card">
    <s-table
      :columns="columns"
      :data-source="projects"
      :loading="loading"
      :pagination="{
        total: totalItems,
        current: serverOptions.page,
        pageSize: serverOptions.rowsPerPage,
        showSizeChanger: true,
        showTotal: (total: number) => `Total ${total} projects`
      }"
      :scroll="{ x: 'max-content' }"
      @change="onChange"
    >
      <!-- Title Column -->
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'title'">
          <router-link
            :to="`/projects/${record.id}`"
            class="text-primary-600 hover:text-primary-700 font-medium"
          >
            {{ record.title }}
          </router-link>
        </template>

        <template v-else-if="column.dataIndex === 'progress'">
          <div class="w-full">
            <div class="flex justify-between text-sm text-gray-600 mb-1">
              <span class="font-medium">${{ record.currentFunding.toLocaleString() }}</span>
              <span>{{ Math.round((record.currentFunding / record.fundingGoal) * 100) }}%</span>
            </div>
            <div class="progress-bar">
              <div 
                class="progress-bar-fill bg-primary-600"
                :style="{ 
                  width: `${Math.min((record.currentFunding / record.fundingGoal) * 100, 100)}%`,
                  background: `linear-gradient(90deg, var(--color-primary-light) 0%, var(--color-primary) 100%)`
                }"
              ></div>
            </div>
          </div>
        </template>

        <template v-else-if="column.dataIndex === 'fundingGoal'">
          <span class="font-medium">${{ record.fundingGoal.toLocaleString() }}</span>
        </template>

        <template v-else-if="column.dataIndex === 'deadline'">
          <span class="text-gray-600">
            {{ new Date(record.deadline).toLocaleDateString(undefined, {
              year: 'numeric',
              month: 'short',
              day: 'numeric'
            }) }}
          </span>
        </template>

        <template v-else-if="column.dataIndex === 'status'">
          <span 
            :class="{
              'badge-success': record.status === 'active',
              'badge-warning': record.status === 'completed',
              'badge-danger': record.status === 'canceled'
            }"
            class="badge"
          >
            {{ record.status.charAt(0).toUpperCase() + record.status.slice(1) }}
          </span>
        </template>
      </template>
    </s-table>
  </div>
</template>

<script setup lang="ts">
import type { Project } from '@/types/project';
import type { STableColumnsType } from '@surely-vue/table';

interface Props {
  projects: Project[];
  loading: boolean;
  totalItems: number;
  serverOptions: {
    page: number;
    rowsPerPage: number;
    sortField?: string;
    sortOrder?: 'ascend' | 'descend';
  };
}

interface Emits {
  (e: 'update:options', options: {
    page: number;
    rowsPerPage: number;
    sortField?: string;
    sortOrder?: 'ascend' | 'descend';
  }): void;
}

const props = withDefaults(defineProps<Props>(), {
  totalItems: 0,
  serverOptions: () => ({
    page: 1,
    rowsPerPage: 10
  })
});

const emit = defineEmits<Emits>();

const columns: STableColumnsType = [
  {
    title: 'Title',
    dataIndex: 'title',
    key: 'title',
    sorter: true,
    width: 250,
    fixed: 'left'
  },
  {
    title: 'Progress',
    dataIndex: 'progress',
    key: 'progress',
    width: 250
  },
  {
    title: 'Goal',
    dataIndex: 'fundingGoal',
    key: 'fundingGoal',
    sorter: true,
    width: 150,
    align: 'right'
  },
  {
    title: 'Deadline',
    dataIndex: 'deadline',
    key: 'deadline',
    sorter: true,
    width: 150
  },
  {
    title: 'Status',
    dataIndex: 'status',
    key: 'status',
    sorter: true,
    width: 120,
    align: 'center'
  }
];

const onChange = (pagination: any, filters: any, sorter: any) => {
  emit('update:options', {
    page: pagination.current,
    rowsPerPage: pagination.pageSize,
    sortField: sorter.field,
    sortOrder: sorter.order
  });
};
</script> 