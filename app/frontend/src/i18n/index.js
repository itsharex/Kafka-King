// src/i18n/index.js
import {createI18n} from 'vue-i18n'
import zhCN from './zh-CN'
import enUS from './en-US'
import jaJP from './ja-JP'
import koKR from './ko-KR'
import ruRU from './ru-RU'

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

export default i18n