<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="width: 42px;">节点</h2>
      <n-button @click="getData" text :render-icon="renderIcon(RefreshOutlined)">refresh</n-button>
      <n-text>共计{{ data.length }}个</n-text>
      <n-button @click="downloadAllDataCsv" :render-icon="renderIcon(DriveFileMoveTwotone)">导出为csv</n-button>

    </n-flex>
    <n-spin :show="loading" description="Connecting...">
      <n-tabs type="segment" animated  v-model:value="activeTab">
        <n-tab-pane name="broker" tab="broker" display-directive="show">
          <n-data-table
              ref="tableRef"
              :columns="columns"
              :data="data"
              size="small"
              :bordered="false"
              striped
          />
        </n-tab-pane>
        <n-tab-pane name="配置" tab="配置" display-directive="show">
          <n-data-table
              :columns="config_columns"
              :data="config_data"
              size="small"
              :bordered="false"
              striped
              :max-height="600"
              virtual-scroll
          />
        </n-tab-pane>
      </n-tabs>

    </n-spin>
  </n-flex>



</template>
<script setup>
import {h, onMounted, ref} from "vue";
import emitter from "../utils/eventBus";
import {NButton, NDataTable, NIcon, NTag, NText, useMessage} from 'naive-ui'
import {createCsvContent, download_file, formattedJson, renderIcon} from "../utils/common";
import {DriveFileMoveTwotone, RefreshOutlined, SettingsTwotone} from "@vicons/material";
import {GetBrokerConfig, GetBrokers} from "../../wailsjs/go/service/Service";

const config_data = ref([])
// 当前活动的 TabPane 名称
const activeTab = ref('broker');
const loading = ref(false)
const data = ref([])
const message = useMessage()
const tableRef = ref();

const selectNode = async (node) => {
  await getData()
}

onMounted(async () => {
  emitter.on('selectNode', selectNode)
  await getData()
})


const getData = async () => {
  loading.value = true
  try {
    const res = await GetBrokers()
    if (res.err !== "") {
      message.error(res.err)
    } else {
      console.log(res)
      const result = res.result
      data.value = result.brokers
    }
  }catch (e) {
    message.error(e)
  }

  loading.value = false

}


const downloadAllDataCsv = async () => {
  const csvContent = createCsvContent(
      activeTab.value === "broker" ? data.value : config_data.value ,
      activeTab.value === "broker" ? columns : config_columns
  )
  download_file(csvContent, '导出.csv', 'text/csv;charset=utf-8;')
}


const columns = [
  { title: 'node_id', key: 'node_id', sorter: 'default',width: 20,resizable: true },
  { title: 'host', key: 'host', sorter: 'default',width: 50,resizable: true },
  { title: 'rack', key: 'rack', sorter: 'default',width: 20,resizable: true },
  {
    title: '配置', key: 'config', width: 30, resizable: true, ellipsis: {tooltip: true},
    render: (row) => h(
        NButton,
        {
          strong: true,
          secondary: true,
          onClick: async () => {
            await getBrokerConfig(row["node_id"])
          }
        },
        {default: () => '配置', icon: () => h(NIcon, null, { default: () => h(SettingsTwotone) })}
    )
  },
]
const getType = (value) => {
  const type = {
    true: "success",
    false: "warning",
  }
  return type[value] || 'error'
}

const config_columns = [
  { title: '配置名', key: 'Name', sorter: 'default',width: 100,resizable: true },
  { title: '值', key: 'Value', sorter: 'default',width: 140,resizable: true },
  { title: '是否只读', key: 'ReadOnly', sorter: 'default',width: 20,resizable: true,
    render: (row) => h(NTag, {type: getType(row['ReadOnly'])}, {default: () => row['ReadOnly'] === true ? "是": "否"}),
  },
  { title: '是否默认', key: 'Default', width: 20,resizable: true,
    render: (row) => h(NTag, {type: getType(row['ReadOnly'])}, {default: () => row['Default'] === true ? "是": "否"}),
  },
  { title: '是否敏感', key: 'Sensitive', width: 20,resizable: true,
    render: (row) => h(NTag, {type: getType(row['Sensitive'])}, {default: () => row['Sensitive'] === true ? "是": "否"}),
  },

]

const getBrokerConfig = async (node_id) => {
  loading.value = true
  try {
    const res = await GetBrokerConfig(node_id)
    console.log(res)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      // 排序
      res.results.sort((a, b) => a["Name"] > b["Name"] ? 1 : -1)
      config_data.value = res.results
      activeTab.value = "配置"
    }
  }catch (e) {
    message.error(e)
  }
  loading.value = false

}


</script>



<style scoped>

</style>