<template>
  <div id="app" class="min-h-screen flex flex-col">
    <header class="bg-gray-800 text-white p-4">
      <nav class="container mx-auto flex justify-between items-center">
        <router-link to="/" class="font-bold text-xl">Crowdfund</router-link>
        <div class="flex items-center">
          <router-link v-if="!userStore.token" to="/login" class="mr-4 hover:text-gray-300">Login</router-link>
          <router-link v-if="!userStore.token" to="/register" class="mr-4 hover:text-gray-300">Register</router-link>
          <router-link v-if="userStore.token" to="/profile" class="mr-4 hover:text-gray-300">Profile</router-link>
          <button v-if="userStore.token" @click="logout" class="hover:text-gray-300">Logout</button>
        </div>
      </nav>
    </header>

    <main class="flex-grow">
      <router-view />
    </main>

    <footer class="bg-gray-800 text-white p-4 text-center">
      <p>&copy; {{ new Date().getFullYear() }} Crowdfund</p>
    </footer>
  </div>
</template>

<script setup>
import { useUserStore } from './store/user';
import { useRouter } from 'vue-router';

const userStore = useUserStore();
const router = useRouter();

const logout = () => {
  userStore.clearUser();
  router.push('/login');
};
</script>

<style scoped>
/* Add any global styles here */
</style>