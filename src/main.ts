import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// Vue.prototype.$apiEndpoint = 'http://localhost:3333/api'

import PrimeVue from 'primevue/config'
import 'primevue/resources/themes/aura-light-green/theme.css'

const app = createApp(App)
app.use(PrimeVue)

app.use(router)

app.mount('#app')
