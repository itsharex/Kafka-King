import {createApp} from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import i18n from './i18n/init'


const app = createApp(App)

app.use(i18n)
app.use(naive)

app.mount('#app')