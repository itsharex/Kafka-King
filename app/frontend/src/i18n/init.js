/*
 *
 *  * Copyright (c) 2025 Bronya0 <tangssst@163.com>. All rights reserved.
 *  * Original source: https://github.com/Bronya0
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *    http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

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