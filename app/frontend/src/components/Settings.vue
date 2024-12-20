<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2>{{ t('settings.title') }}</h2>
    </n-flex>
  </n-flex>
  <n-flex align="center">

    <n-form :model="config" label-placement="top" style="text-align: left;">

      <n-form-item :label="t('settings.width')" >
        <n-input-number v-model:value="config.width" :min="800" :max="1920" :style="{ maxWidth: '120px' }"/>
      </n-form-item>
      <n-form-item :label="t('settings.height')">
        <n-input-number v-model:value="config.height" :min="600" :max="1080" :style="{ maxWidth: '120px' }"/>
      </n-form-item>
      <n-form-item :label="t('settings.lang')">
        <n-select v-model:value="config.language" :options="languageOptions"
                  :style="{ maxWidth: '120px' }"/>
      </n-form-item>

      <n-form-item :label="t('settings.theme')" >
        <n-button circle :focusable="false" @click="changeTheme" :render-icon="renderIcon(MoonOrSunnyOutline)"/>
      </n-form-item>

      <n-form-item>
        <n-button @click="saveConfig" strong type="primary">{{ t('common.save')}}</n-button>
      </n-form-item>

      <n-form-item label="">
      </n-form-item>
      <n-form-item label="">
      </n-form-item>

    </n-form>
  </n-flex>
</template>

<script setup>
import {onMounted, ref, shallowRef} from 'vue'
import {darkTheme, lightTheme, NButton, NForm, NFormItem, NInputNumber, NSelect, useMessage,} from 'naive-ui'
import {GetConfig, SaveConfig} from '../../wailsjs/go/config/AppConfig'
import {WindowSetSize} from "../../wailsjs/runtime";
import {renderIcon} from "../utils/common";
import {NightlightRoundFilled, WbSunnyOutlined} from '@vicons/material'
import emitter from "../utils/eventBus";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const message = useMessage()
let theme = lightTheme
let MoonOrSunnyOutline = shallowRef(WbSunnyOutlined)

const config = ref({
  width: 1248,
  height: 768,
  language: 'en-US',
  theme: theme.name,
})
const languageOptions = [
  {label: '中文', value: 'zh-CN'},
  {label: 'English', value: 'en-US'},
  {label: '日本語', value: 'ja-JP'},
  {label: '한국인', value: 'ko-KR'},
  {label: 'рускі', value: 'ru-RU'},
]

onMounted(async () => {
  console.info("初始化settings……")

  // 从后端加载配置
  const loadedConfig = await GetConfig()
  console.log(loadedConfig)
  config.value = loadedConfig
  theme = loadedConfig.theme === lightTheme.name ? lightTheme : darkTheme
  MoonOrSunnyOutline.value = loadedConfig.theme === lightTheme.name ? WbSunnyOutlined : NightlightRoundFilled;
})


const saveConfig = async () => {
  config.value.theme = theme.name
  const err = await SaveConfig(config.value)
  if (err !== "") {
    message.error(t('message.saveErr') + "：" + err)
    return
  }

  await WindowSetSize(config.value.width, config.value.height)

  emitter.emit('update_theme', theme)
  // 可以添加保存成功的提示
  message.success(t('message.saveSuccess'))
  config.value = await GetConfig()

}

// 语言变更
// const changeLang = (value) => {
//   emitter.emit('language_change', value)
// }

const changeTheme = () => {
  MoonOrSunnyOutline.value = MoonOrSunnyOutline.value === NightlightRoundFilled ? WbSunnyOutlined : NightlightRoundFilled;
  theme = MoonOrSunnyOutline.value === NightlightRoundFilled ? darkTheme : lightTheme
  emitter.emit('update_theme', theme)
}

</script>