<template>
  <div class="container mx-auto p-4">
    <h2 class="text-2xl font-semibold mb-4">Profile</h2>
    <div v-if="userStore.user">
      <p>Username: {{ userStore.user.username }}</p>
      <p>Email: {{ userStore.user.email }}</p>
      </div>
    <div v-else>
      <p>Loading user profile...</p>
    </div>
  </div>
</template>

<script setup>
import { useUserStore } from '../store/user';
import { onMounted } from 'vue';
import axios from 'axios';

const userStore = useUserStore();

onMounted(async () => {
  if (!userStore.user && userStore.token) {
    try {
      const response = await axios.get('/api/users/profile', {
        headers: {
          Authorization: `Bearer ${userStore.token}`,
        },
      });
      userStore.setUser(response.data);
    } catch (error) {
      console.error('Error fetching user profile:', error);
      // Handle error (e.g., redirect to login)
    }
  }
});
</script>