import {createApp} from 'vue'
import naive from 'naive-ui'
import App from './App.vue'
import {createI18n} from "vue-i18n";
import zhCN from "./i18n/zh-CN";
import enUS from "./i18n/en-US";
import jaJP from "./i18n/ja-JP";
import koKR from "./i18n/ko-KR";
import ruRU from "./i18n/ru-RU";

const i18n = createI18n({
    legacy: false, // 使用Composition API
    locale: 'zh-CN', // 默认语言
    fallbackLocale: 'en-US', // 备用语言
    messages: {
        'zh-CN': zhCN,
        'en-US': enUS,
        'ja-JP': jaJP,
        'ko-KR': koKR,
        'ru-RU': ruRU,
    }
})

const app = createApp(App)

app.use(i18n)
app.use(naive)

app.mount('#app')