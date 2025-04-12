<!--
  - Copyright 2025 Bronya0 <tangssst@163.com>.
  - Author Github: https://github.com/Bronya0
  -
  - Licensed under the Apache License, Version 2.0 (the "License");
  - you may not use this file except in compliance with the License.
  - You may obtain a copy of the License at
  -
  -     https://www.apache.org/licenses/LICENSE-2.0
  -
  - Unless required by applicable law or agreed to in writing, software
  - distributed under the License is distributed on an "AS IS" BASIS,
  - WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  - See the License for the specific language governing permissions and
  - limitations under the License.
  -->

<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2>{{ t('node.title') }}</h2>
      <n-button @click="getData" text :render-icon="renderIcon(RefreshOutlined)">{{ t('common.refresh') }}</n-button>
      <n-text>{{ t('common.count') }}：{{ data.length }}</n-text>
      <n-button @click="downloadAllDataCsv" :render-icon="renderIcon(DriveFileMoveTwotone)">{{ t('common.csv') }}
      </n-button>
    </n-flex>
    <n-spin :show="loading" :description="t('common.connecting')">
      <n-tabs type="line" animated v-model:value="activeTab">
        <n-tab-pane name="Broker">
          <template #tab>
            {{ t('node.title') }}
          </template>
          <n-data-table
              :columns="columns"
              :data="data"
              size="small"
              :bordered="false"
              striped
              :pagination="pagination"
          />
        </n-tab-pane>
        <n-tab-pane name="Config">
          <template #tab>
            {{ t('common.config') }}
          </template>
          <n-flex vertical>

            <n-flex align="center">
              <n-input :disabled='activeConfigNode' placeholder="search" v-model:value="configSearchText" clearable style="width: 300px"/>
              <n-button :disabled='activeConfigNode' @click="getBrokerConfig(activeConfigNode)" :render-icon="renderIcon(RefreshOutlined)">
                {{ t('common.refresh') }}
              </n-button>
            </n-flex>
            <n-data-table
                :columns="config_columns"
                :data="config_data"
                :bordered="false"
                :pagination="pagination"

            />
          </n-flex>

        </n-tab-pane>
      </n-tabs>

    </n-spin>
  </n-flex>


</template>
<script setup>
import {h, onMounted, ref} from "vue";
import emitter from "../utils/eventBus";
import {NButton, NDataTable, NIcon, NInput, NTag, NText, useMessage} from 'naive-ui'
import {createCsvContent, download_file, getCurrentDateTime, renderIcon} from "../utils/common";
import {DriveFileMoveTwotone, RefreshOutlined, SettingsTwotone} from "@vicons/material";
import {AlterNodeConfig, GetBrokerConfig, GetBrokers} from "../../wailsjs/go/service/Service";
import ShowOrEdit from "../common/ShowOrEdit.vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const config_data = ref([])
const data = ref([])
const configSearchText = ref("")
// 当前活动的 TabPane 名称
const activeTab = ref('Broker');
const activeConfigNode = ref('');
const loading = ref(false)
const message = useMessage()

const selectNode = async (node) => {
  config_data.value = []
  data.value = []
  activeConfigNode.value = ''
  configSearchText.value = ''
  loading.value = false

  await getData()
}

onMounted(() => {
  emitter.on('selectNode', selectNode)
  getData()
})


const getData = async () => {
  loading.value = true
  try {
    const res = await GetBrokers()
    if (res.err !== "") {
      message.error(res.err)
    } else {
      const result = res.result
      data.value = result.brokers
    }
  } catch (e) {
    message.error(e)
  }

  loading.value = false

}

const pagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [5, 10, 20, 30, 40],
  onChange: (page) => {
    pagination.value.page = page
  },
  onUpdatePageSize: (pageSize) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1
  },
})

const downloadAllDataCsv = async () => {
  const csvContent = createCsvContent(
      activeTab.value === "Broker" ? data.value : config_data.value,
      activeTab.value === "Broker" ? columns : config_columns
  )
  download_file(csvContent, `${getCurrentDateTime()}.csv`, 'text/csv;charset=utf-8;')
}


const columns = [
  {title: 'node_id', key: 'node_id', sorter: 'default', width: 20, resizable: true},
  {
    title: 'host', key: 'host', sorter: 'default', width: 50, resizable: true,
    render: (row) => h(NTag, {type: "info"}, {default: () => row['host']}),
  },
  {
    title: 'port', key: 'port', sorter: 'default', width: 20, resizable: true,
    render: (row) => h(NTag, {type: "success"}, {default: () => row['port']}),
  },
  {title: 'rack', key: 'rack', sorter: 'default', width: 20, resizable: true},
  {
    title: 'config', key: 'config', width: 30, resizable: true, ellipsis: {tooltip: {style: {maxWidth: '800px'},}},
    render: (row) => h(
        NButton,
        {
          strong: true,
          secondary: true,
          onClick: async () => {
            await getBrokerConfig(row["node_id"])
            activeConfigNode.value = row["node_id"]
          }
        },
        {default: () => t('common.config'), icon: () => h(NIcon, null, {default: () => h(SettingsTwotone)})}
    )
  },
]

const config_columns = [
  {
    title: 'Name', key: 'Name', sorter: 'default', width: 80, resizable: true,
  },
  {
    title: t('node.value'), key: 'Value', sorter: 'default', width: 140, resizable: true,
    render: (row) => {
      return h(ShowOrEdit, {
        value: row['Value'],
        onUpdateValue(v) {
          alterNodeConfig(activeConfigNode.value, row['Name'], v)
        }
      })
    }
  },
  {title: t('node.source'), key: 'Source', sorter: 'default', width: 50, resizable: true,},
  {
    title: t('node.sensitive'),
    key: 'Sensitive',
    width: 20,
    resizable: true,
    sorter: (row1, row2) => Number(row1['Sensitive']) - Number(row2['Sensitive']),
    render: (row) => h(NTag, {type: row['Sensitive'] === true ? "error" : "info"}, {default: () => row['Sensitive'] === true ? "yes" : "no"}),
  },

]

const getBrokerConfig = async (node_id) => {
  loading.value = true
  try {
    const res = await GetBrokerConfig(node_id)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      // 排序
      res.results.sort((a, b) => a["Name"] > b["Name"] ? 1 : -1)
      if (configSearchText.value){
        res.results = res.results.filter(item => item['Name'].includes(configSearchText.value))
      }
      config_data.value = res.results
      activeTab.value = "Config"
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}


const alterNodeConfig = async (node_id, name, value) => {
  loading.value = true
  try {
    const res = await AlterNodeConfig(node_id, name, value)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success(t('node.ok_message'))
      await getBrokerConfig(activeConfigNode.value)
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}

</script>


<style scoped>

</style>