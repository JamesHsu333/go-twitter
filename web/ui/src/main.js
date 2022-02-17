import { createApp } from 'vue'
import store from './store'
import routers from "./routers";
import './routers/guardian'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'

import './css/index.css'

const app = createApp(App)

app.use(store)
app.use(routers)
app.use(ElementPlus)
app.mount('#app')