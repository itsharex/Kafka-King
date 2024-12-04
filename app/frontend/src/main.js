import {createApp} from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import i18n from './i18n'


const app = createApp(App)

app.use(naive)
app.use(i18n)

app.mount('#app')