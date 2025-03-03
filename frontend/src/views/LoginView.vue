<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center py-8">
    <div class="container">
       <div class="w-2/4 mx-auto border-b-2 border-blue-600">
          <LoginForm />
      </div>
    </div>
  </div>
</template>

<script setup>
import LoginForm from '@/components/LoginForm.vue';
</script>