import zhCN from "./zh-CN";
import enUS from "./en-US";
import jaJP from "./ja-JP";
import koKR from "./ko-KR";
import ruRU from "./ru-RU";
import {createI18n} from "vue-i18n";

const languageMap = {
    'zh-CN': zhCN,
    'en-US': enUS,
    'ja-JP': jaJP,
    'ko-KR': koKR,
    'ru-RU': ruRU,
};
const default_language = 'en-US';

const i18n = createI18n({
    legacy: false, // 使用Composition API
    locale: default_language, // 使用检测到的语言
    fallbackLocale: default_language, // 备用语言
    messages: languageMap, // 语言包
})

export default i18n;