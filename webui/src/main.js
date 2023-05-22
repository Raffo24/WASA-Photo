import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import {axios, axiosUpdate as axiosUpdate, getCurrentSession} from './services/axios.js'
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import PostCard from './components/PostCard.vue'
import UserCard from './components/UserCard.vue'
import Modal from './components/Modal.vue'
import 'bootstrap-icons/font/bootstrap-icons.css'
import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.config.globalProperties.$axiosUpdate = axiosUpdate;
app.config.globalProperties.$currentSession = getCurrentSession;

app.component("UserCard", UserCard);
app.component("Modal", Modal);
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("PostCard", PostCard);

app.use(router)
app.mount('#app')
