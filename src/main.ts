import './assets/main.css'

import App from './App.vue'
import router from './router'

// Vue.prototype.$apiEndpoint = 'http://localhost:3333/api'

import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'
import Vue, { createApp } from '@vue/compat'

// Import Bootstrap and BootstrapVue CSS files (order is important)
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

// Make BootstrapVue available throughout your project
// Optionally install the BootstrapVue icon components plugin

Vue.use(BootstrapVue)

const app = createApp(App)

app.use(router)

app.mount('#app')
