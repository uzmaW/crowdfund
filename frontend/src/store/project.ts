// store/project.ts
import { defineStore } from 'pinia';
import axios from 'axios';
import type { Project } from '@/types/project';

interface ProjectState {
  projects: Project[];
  currentProject: Project | null;
  totalProjects: number;
  loadingProjects: boolean;
  loadingCurrentProject: boolean;
  error: string | null;
}

interface FetchProjectsParams {
  search?: string;
  page?: number;
  limit?: number;
  sortBy?: string;
  sortDirection?: 'asc' | 'desc';
}

export const useProjectStore = defineStore('project', {
  state: (): ProjectState => ({
    projects: [],
    currentProject: null,
    totalProjects: 0,
    loadingProjects: false,
    loadingCurrentProject: false,
    error: null,
  }),

  actions: {
    async fetchProjects(params?: FetchProjectsParams) {
      this.loadingProjects = true;
      this.error = null;
      
      try {
        const response = await axios.get('/api/projects', {
          params: {
            search: params?.search,
            page: params?.page,
            limit: params?.limit,
            sort_by: params?.sortBy,
            sort_dir: params?.sortDirection,
          }
        });

        this.projects = response.data.items;
        this.totalProjects = response.data.total;
      } catch (error) {
        this.handleError(error, 'Failed to fetch projects');
      } finally {
        this.loadingProjects = false;
      }
    },

    async fetchProject(id: string) {
      this.loadingCurrentProject = true;
      this.error = null;
      
      try {
        const response = await axios.get<Project>(`/api/projects/${id}`);
        this.currentProject = response.data;
      } catch (error) {
        this.handleError(error, 'Failed to fetch project details');
      } finally {
        this.loadingCurrentProject = false;
      }
    },

    clearCurrentProject() {
      this.currentProject = null;
    },

    handleError(error: unknown, defaultMessage: string) {
      if (axios.isAxiosError(error)) {
        this.error = error.response?.data?.message || defaultMessage;
      } else {
        this.error = defaultMessage;
      }
      console.error(error);
    }
  },

  getters: {
    featuredProjects: (state) => state.projects.slice(0, 3),
    activeProjects: (state) => state.projects.filter(p => p.status === 'active'),
  },
});