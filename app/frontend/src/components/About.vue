<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2>{{ t('about.title') }}</h2>
    </n-flex>
  </n-flex>
  <n-flex align="center">

    <n-form label-placement="top" style="text-align: left;">
      <n-form-item :label="t('about.projectHomepage')">
        <n-button @click="BrowserOpenURL(kafka_home_url)" :render-icon="renderIcon(HouseTwotone)">
          {{ t('about.kafkaKing') }} Github
        </n-button>
      </n-form-item>
      <n-form-item :label="t('about.esClient')">
        <n-button @click="BrowserOpenURL(es_home_url)" :render-icon="renderIcon(HouseTwotone)">
          {{ t('about.esKing') }} Github
        </n-button>
      </n-form-item>
      <n-form-item :label="t('about.technicalGroup')">
        <n-button :focusable="false" @click="openUrl(qq_url)">
          {{ t('about.qqGroup') }}
        </n-button>
      </n-form-item>
      <n-form-item :label="t('about.translate')">
        <n-button :focusable="false" @click="openUrl(es_translation_url)">
          Github
        </n-button>
      </n-form-item>
    </n-form>
  </n-flex>
</template>

<script setup>
import {onMounted} from 'vue'
import {NButton, NForm, NFormItem, useMessage,} from 'naive-ui'
import {BrowserOpenURL} from "../../wailsjs/runtime";
import {openUrl, renderIcon} from "../utils/common";
import {HouseTwotone} from '@vicons/material'
import {useI18n} from "vue-i18n";

const { t } = useI18n()

const es_home_url = "https://github.com/Bronya0/ES-King"
const es_translation_url = "https://github.com/Bronya0/Kafka-King/issues/51"
const kafka_home_url = "https://github.com/Bronya0/kafka-King"
const qq_url = "https://qm.qq.com/cgi-bin/qm/qr?k=pDqlVFyLMYEEw8DPJlRSBN27lF8qHV2v&jump_from=webapi&authKey=Wle/K0ARM1YQWlpn6vvfiZuMedy2tT9BI73mUvXVvCuktvi0fNfmNR19Jhyrf2Nz"

const message = useMessage()

onMounted(async () => {

})



</script>