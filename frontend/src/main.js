import { createApp } from 'vue'
// import { createPinia } from 'pinia'

import App from './App.vue'
import './assets/main.css'

import './assets/vis/vis-timeline.css'
import Multiselect from 'vue-multiselect'
// import router from './router'

createApp(App)
    //.use(router)
    .component('Multiselect', Multiselect)
    .mount('#app')
