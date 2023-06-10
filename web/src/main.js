import { createApp } from 'vue';
import App from './App.vue';
import router from './pages/_routes';
import './styles.scss';

createApp(App).use(router).mount('#app');
