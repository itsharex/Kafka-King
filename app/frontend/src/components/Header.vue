<template>
  <n-page-header style="padding: 4px;--wails-draggable:drag">
    <template #avatar>
      <n-avatar :src="logo"/>
    </template>
    <template #title>
      <div style="font-weight: 800">{{app_name}}</div>
    </template>
    <template #subtitle>
      <n-tooltip>
        <template #trigger>
          <n-tag :type=title_tag v-if="subtitle">{{subtitle}}</n-tag>
          <n-p v-else>{{t('header.desc')}}</n-p>
        </template>
      </n-tooltip>
    </template>
    <template #extra>
      <n-flex justify="flex-end" style="--wails-draggable:no-drag" class="right-section">
<!--        <n-button quaternary :focusable="false" @click="openUrl(qq_url)">交流群</n-button>-->
<!--        <n-button quaternary :focusable="false" @click="changeTheme" :render-icon="renderIcon(MoonOrSunnyOutline)"/>-->
        <n-tooltip placement="bottom" trigger="hover">
          <template #trigger>
            <n-button quaternary @click="openUrl(update_url)" :render-icon="renderIcon(CodeFilled)"/>
          </template>
          <span>Source Code</span>
        </n-tooltip>

        <n-tooltip placement="bottom" trigger="hover">
          <template #trigger>
            <n-button quaternary :focusable="false" :loading="update_loading" @click="checkForUpdates"
                      :render-icon="renderIcon(SystemUpdateAltSharp)"/>
          </template>
          <span> Check Version：{{ version.tag_name }} {{ check_msg }}</span>
        </n-tooltip>
        <n-button quaternary :focusable="false" @click="minimizeWindow" :render-icon="renderIcon(RemoveOutlined)"/>
        <n-button quaternary :focusable="false" @click="resizeWindow" :render-icon="renderIcon(MaxMinIcon)"/>
        <n-button quaternary style="font-size: 22px" :focusable="false" @click="closeWindow">
          <n-icon>
            <CloseFilled/>
          </n-icon>
        </n-button>
      </n-flex>
    </template>
  </n-page-header>
</template>

<script setup>
import {NAvatar, NButton, NFlex, useNotification} from 'naive-ui'
import {
  CloseFilled,
  CodeFilled,
  ContentCopyFilled,
  CropSquareFilled,
  RemoveOutlined,
  SystemUpdateAltSharp
} from '@vicons/material'
import logo from '../assets/images/appicon.png'
import {h, onMounted, ref, shallowRef} from "vue";
import {BrowserOpenURL, Quit, WindowMaximise, WindowMinimise, WindowUnmaximise} from "../../wailsjs/runtime";
import {CheckUpdate} from '../../wailsjs/go/system/Update'
import {openUrl, renderIcon} from "../utils/common";
import {GetAppName, GetVersion} from "../../wailsjs/go/config/AppConfig";
import emitter from "../utils/eventBus";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

// defineProps(['options', 'value']);

// const MoonOrSunnyOutline = shallowRef(WbSunnyOutlined)
const isMaximized = ref(false);
const check_msg = ref("");
const app_name = ref("");
const title_tag = ref("success");
const MaxMinIcon = shallowRef(CropSquareFilled)
const update_url = "https://github.com/Bronya0/Kafka-King/releases"
const update_loading = ref(false)
// let theme = lightTheme

let version = ref({
  tag_name: "",
  body: "",
})

const subtitle = ref("")

const notification = useNotification()

onMounted(async () => {
  emitter.on('selectNode', selectNode)
  // emitter.on('changeTitleType', changeTitleType)

  app_name.value = await GetAppName()

  // const config = await GetConfig()
  // MoonOrSunnyOutline.value = config.theme === lightTheme.name ? WbSunnyOutlined : NightlightRoundFilled
  const v = await GetVersion()
  version.value.tag_name = v
  subtitle.value = t('header.desc') + " " + v
  await checkForUpdates()
})

const selectNode = (node) => {
  subtitle.value = `${t('header.c_node')}：【` + node.name + "】"
}

const checkForUpdates = async () => {
  update_loading.value = true
  try {
    const v = await GetVersion()
    const resp = await CheckUpdate()
    if (!resp) {
      check_msg.value = `${t('header.netErr')}`
    } else if (resp.tag_name !== v) {
      check_msg.value = `${t('header.newVersion')} ` + resp.tag_name
      version.value.body = resp.body
      const n = notification.success({
        title: check_msg.value,
        content: resp.body,
        action: () =>
              h(NFlex, {justify: "flex-end" }, () => [
                h(
                    NButton,
                    {
                      type: 'primary',
                      secondary: true,
                      onClick: () => BrowserOpenURL(update_url),
                    },
                    () => t('header.down'),
                ),
                h(
                    NButton,
                    {
                      secondary: true,
                      onClick: () => {
                        n.destroy()
                      },
                    },
                    () => t('common.cancel'),
                ),
            ]),
        onPositiveClick: () => BrowserOpenURL(update_url),
      })
    }
  } finally {
    update_loading.value = false
  }
}

const minimizeWindow = () => {
  WindowMinimise()
}

const resizeWindow = () => {
  isMaximized.value = !isMaximized.value;
  if (isMaximized.value) {
    WindowMaximise();
    MaxMinIcon.value = ContentCopyFilled;
  } else {
    WindowUnmaximise();
    MaxMinIcon.value = CropSquareFilled;
  }
  console.log(isMaximized.value)

}

const closeWindow = () => {
  Quit()
}
// const changeTheme = () => {
//   MoonOrSunnyOutline.value = MoonOrSunnyOutline.value === NightlightRoundFilled ? WbSunnyOutlined : NightlightRoundFilled;
//   theme = MoonOrSunnyOutline.value === NightlightRoundFilled ? darkTheme : lightTheme
//   emitter.emit('update_theme', theme)
// }
</script>

<style scoped>


.right-section .n-button {
  padding: 0 8px;
}
</style>