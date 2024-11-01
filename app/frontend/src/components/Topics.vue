<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="width: 42px;">主题</h2>
      <n-button @click="getData" text :render-icon="renderIcon(RefreshOutlined)">refresh</n-button>
      <n-text>共计{{ data.length }}个</n-text>
      <n-button @click="downloadAllDataCsv" :render-icon="renderIcon(DriveFileMoveTwotone)">导出为csv</n-button>

    </n-flex>
    <n-spin :show="loading" description="Connecting...">
      <n-tabs type="segment" animated v-model:value="activeTab">
        <n-tab-pane name="主题" tab="主题">
          <n-flex vertical>
            <!--          搜索框、新增按钮-->
            <n-flex align="center">
              <n-input v-model:value="searchText" placeholder="输入主题名称" clearable style="width: 300px" />
              <n-button @click="getData" :render-icon="renderIcon(SearchOutlined)">search</n-button>
              <n-button @click="showDrawer=true" :render-icon="renderIcon(SearchOutlined)">创建主题</n-button>
            </n-flex>
            <n-data-table
                ref="tableRef"
                :columns="columns"
                :data="data"
                size="small"
                :bordered="false"
                striped
                :pagination="pagination"

            />
          </n-flex>

        </n-tab-pane>

        <n-tab-pane name="详情" tab="详情">
          <n-data-table
              :columns="partitions_columns"
              :data="partitions_data"
              size="small"
              :bordered="false"
              striped
              :max-height="600"
              virtual-scroll
          />
        </n-tab-pane>

        <n-tab-pane name="配置" tab="配置">
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

  <n-drawer v-model:show="showDrawer" :width="500">
    <n-drawer-content title="创建Topic配置">
      <n-form
          ref="formRef"
          :model="topic_add"
          label-placement="left"
          label-width="120"
          require-mark-placement="right-hanging"
      >
        <n-form-item label="Topic名称" path="topics">
          <n-dynamic-tags
              v-model:value="topic_add.topics"
              :max="10"
              placeholder="请输入Topic名称后回车"
          />
        </n-form-item>

        <n-form-item label="分区数" path="partitions">
          <n-input-number
              v-model:value="topic_add.partitions"
              :min="1"
              :max="100"
              placeholder="请输入分区数"
          />
        </n-form-item>

        <n-form-item label="副本因子" path="replicationFactor">
          <n-input-number
              v-model:value="topic_add.replication_factor"
              :min="1"
              :max="5"
              placeholder="请输入副本因子"
          />
        </n-form-item>

        <n-form-item label="Topic配置" path="config">
          <n-input
              v-model:value="topic_add.configs"
              type="textarea"
              placeholder="请输入JSON格式的Topic配置"
              :rows="8"
          />
        </n-form-item>

      </n-form>

      <template #footer>
        <n-space>
          <n-button @click="showDrawer = false">取消</n-button>
          <n-button type="primary" @click="addTopic">确认</n-button>
        </n-space>
      </template>
    </n-drawer-content>
  </n-drawer>

</template>
<script setup>
import {h, onMounted, ref} from "vue";
import emitter from "../utils/eventBus";
import {NButton, NButtonGroup, NDataTable, NIcon, NPopconfirm, NTag, NText, useMessage} from 'naive-ui'
import {createCsvContent, download_file, renderIcon} from "../utils/common";
import {
  DeleteForeverTwotone,
  DriveFileMoveTwotone,
  InfoOutlined,
  RefreshOutlined,
  SearchOutlined,
  SettingsTwotone
} from "@vicons/material";
import {CreateTopics, DeleteTopic, GetTopicConfig, GetTopics} from "../../wailsjs/go/service/Service";

const config_data = ref([])
const partitions_data = ref([])
const activeTab = ref('主题');
const loading = ref(false)
const data = ref([])
const topic_add = ref({
  topics: [],
  partitions: 1,
  replication_factor: 1,
  configs: ""
})
const message = useMessage()
const tableRef = ref();
const searchText = ref("");
const showDrawer = ref(false)

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
    const res = await GetTopics()
    if (res.err !== "") {
      message.error(res.err)
    } else {
      console.log(res)
      data.value = res.results
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
  itemCount: data.value.length
})


const downloadAllDataCsv = async () => {
  const csvContent = createCsvContent(
      activeTab.value === "broker" ? data.value : config_data.value,
      activeTab.value === "broker" ? columns : config_columns
  )
  download_file(csvContent, '导出.csv', 'text/csv;charset=utf-8;')
}

const getType = (value) => {
  const type = {
    true: "success",
    false: "warning",
  }
  return type[value] || 'error'
}

const columns = [
  {title: 'ID', key: 'ID', sorter: 'default', width: 40, resizable: true, ellipsis: {tooltip: true},},
  {title: 'topic', key: 'topic', sorter: 'default', width: 80, resizable: true, ellipsis: {tooltip: true},},
  {title: '主题故障', key: 'Err', sorter: 'default', width: 40, resizable: true, ellipsis: {tooltip: true},},
  {title: '分区数', key: 'partition_count', sorter: 'default', width: 20, resizable: true},
  {title: '副本因子', key: 'replication_factor', sorter: 'default', width: 30, resizable: true},
  {
    title: '内部主题',
    key: 'IsInternal',
    width: 30,
    resizable: true,
    sorter: (row1, row2) => Number(row1['IsInternal']) - Number(row2['IsInternal']),
    render: (row) => h(NTag, {type: getType(row['IsInternal'])}, {default: () => row['IsInternal'] === true ? "是" : "否"}),
  },
  {
    title: '操作',
    key: 'actions',
    width: 140,  // 调整宽度以适应两个按钮
    resizable: true,
    render: (row) => h(
        NButtonGroup,
        {
          vertical: false,
        },
        {
          default: () => [
            h(
                NButton,
                {
                  strong: true,
                  secondary: true,
                  onClick: async () => {
                    await getTopicConfig(row["topic"])
                  }
                },
                {
                  default: () => '配置',
                  icon: () => h(NIcon, null, {default: () => h(SettingsTwotone)})
                }
            ),
            h(
                NButton,
                {
                  strong: true,
                  secondary: true,
                  onClick: async () => {
                    await getTopicDetail(row["topic"])
                  }
                },
                {
                  default: () => '详情',
                  icon: () => h(NIcon, null, {default: () => h(InfoOutlined)})
                }
            ),
            h(
                NPopconfirm,
                {
                  onPositiveClick: () => deleteTopic(row["topic"])
                },
                {
                  trigger: () =>
                      h(
                          NButton,
                          {
                            strong: true,
                            secondary: true,
                            type: 'error'
                          },
                          {default: () => '删除',
                            icon: () => h(NIcon, null, {default: () => h(DeleteForeverTwotone)})
                          }
                      ),
                  default: () => `确认删除${row["topic"]}?`
                }
            ),
          ]
        }
    )
  },
]

const partitions_columns = [
  {title: 'ID', key: 'partition', sorter: 'default', width: 20, resizable: true},
  {title: '分区故障', key: 'err', sorter: 'default', width: 30, resizable: true, ellipsis: {tooltip: true},},
  {title: 'Leader ID', key: 'leader', sorter: 'default', width: 15, resizable: true},
  {title: 'LeaderEpoch', key: 'LeaderEpoch', sorter: 'default', width: 15, resizable: true},
  {title: '托管此分区的副本ID集', key: 'replicas', sorter: 'default', width: 15, resizable: true},
  {title: 'ISR副本ID集', key: 'isr', sorter: 'default', width: 15, resizable: true},
  {title: '离线副本ID集', key: 'OfflineReplicas', sorter: 'default', width: 15, resizable: true},
]
const config_columns = [
  {title: '配置名', key: 'Name', sorter: 'default', width: 100, resizable: true},
  {title: '值', key: 'Value', sorter: 'default', width: 140, resizable: true},
  {
    title: '是否只读',
    key: 'ReadOnly',
    width: 20,
    resizable: true,
    sorter: (row1, row2) => Number(row1['ReadOnly']) - Number(row2['ReadOnly']),
    render: (row) => h(NTag, {type: getType(row['ReadOnly'])}, {default: () => row['ReadOnly'] === true ? "是" : "否"}),
  },
  {
    title: '是否默认',
    key: 'Default',
    width: 20,
    resizable: true,
    sorter: (row1, row2) => Number(row1['Default']) - Number(row2['Default']),
    render: (row) => h(NTag, {type: getType(row['ReadOnly'])}, {default: () => row['Default'] === true ? "是" : "否"}),
  },
  {
    title: '是否敏感',
    key: 'Sensitive',
    width: 20,
    resizable: true,
    sorter: (row1, row2) => Number(row1['Sensitive']) - Number(row2['Sensitive']),
    render: (row) => h(NTag, {type: getType(row['Sensitive'])}, {default: () => row['Sensitive'] === true ? "是" : "否"}),
  },

]


const getTopicConfig = async (topic) => {
  loading.value = true
  try {
    const res = await GetTopicConfig(topic)
    console.log(res)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      // 排序
      res.results.sort((a, b) => a["Name"] > b["Name"] ? 1 : -1)
      config_data.value = res.results
      activeTab.value = "配置"
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}
const getTopicDetail = async (topic) => {
  loading.value = true
  try {
    partitions_data.value = data.value.find(item => item.topic === topic).partitions
    activeTab.value = "详情"
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}
const deleteTopic = async (topic) => {
  loading.value = true
  try {
    const res = await DeleteTopic([topic])
    console.log(res)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success("删除成功")
      await getData()
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}

const addTopic = async () => {
  loading.value = true
  try {
    const res = await CreateTopics(topic_add.value.topics, topic_add.value.partitions, topic_add.value.replication_factor, topic_add.value.configs)
    console.log(res)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success("删除成功")
      await getData()
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}

</script>


<style scoped>

</style>