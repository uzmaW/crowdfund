import { createApp } from 'vue';
import App from '@/App.vue';
import router from '@/router';
import pinia from '@/store';
import Surely from '@surely-vue/table';
import '@/assets/styles/index.css'; // Import Tailwind CSS

const app = createApp(App);

app.use(router);
app.use(pinia);
app.use(Surely);

app.mount('#app');