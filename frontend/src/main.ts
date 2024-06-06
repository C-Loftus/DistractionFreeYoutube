import './assets/main.css'

import App from './App.vue'
import router from './router'
import { createApp } from 'vue'

// Vue.prototype.$apiEndpoint = 'http://localhost:3333/api'

const app = createApp(App)
app.config.globalProperties.$apiEndpoint = 'http://localhost:3333/api'

app.use(router)

app.mount('#app')
