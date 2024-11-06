<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="max-width: 200px;">主题</h2>
      <n-button @click="getData" text :render-icon="renderIcon(RefreshOutlined)">刷新</n-button>
      <n-text>共计{{ data.length }}个</n-text>
      <n-button @click="downloadAllDataCsv" :render-icon="renderIcon(DriveFileMoveTwotone)">导出为csv</n-button>

    </n-flex>
    <n-spin :show="loading" description="Connecting...">
      <n-tabs type="line" animated v-model:value="activeTab">

        <n-tab-pane name="主题">
          <template #tab>
            <n-icon>
              <LibraryBooksOutlined/>
            </n-icon>
            主题
          </template>
          <n-flex vertical>
            <!--          搜索框、新增按钮-->
            <n-flex align="center">
              <n-input v-model:value="searchText" @keydown.enter="getData" placeholder="输入主题名称" clearable style="width: 300px"/>
              <n-button @click="getData" :render-icon="renderIcon(SearchOutlined)"></n-button>
              <n-button @click="showDrawer=true" :render-icon="renderIcon(AddFilled)">创建主题</n-button>
              <n-button @click="getData" :render-icon="renderIcon(RefreshOutlined)">刷新 Topic</n-button>
<!--              <n-dropdown :options="group_data"  @select="getTopicsOffsets"><n-button :render-icon="renderIcon(RefreshOutlined)">刷新 Offsets</n-button></n-dropdown>-->
              <n-select
                  v-model:value="selectedGroup"
                  :options="group_data"
                  placeholder="选择Group并读取Offsets"
                  filterable
                  clearable
                  style="width: 250px"
              />
              <n-button @click="getTopicsOffsets" :render-icon="renderIcon(RefreshOutlined)">刷新 Offsets</n-button>
            </n-flex>
            <n-data-table
                :columns="columns"
                :data="data"
                size="small"
                :bordered="false"
                striped
                :pagination="pagination"
                :row-key="rowKey"
                v-model:checked-row-keys="selectedRowKeys"
            />
          </n-flex>

        </n-tab-pane>

        <n-tab-pane name="分区">
          <template #tab>
            <n-icon>
              <AddRoadOutlined/>
            </n-icon>
            分区
          </template>

          <n-flex vertical>
            <n-flex align="center">
              <n-text>topic: </n-text>
              <n-tag type="success">
                {{ activeDetailTopic }}
              </n-tag>
              <n-button @click="showModal=true" :render-icon="renderIcon(AddFilled)">添加分区</n-button>
              <n-select
                  v-model:value="selectedGroup"
                  :options="group_data"
                  placeholder="选择Group并读取Offsets"
                  filterable
                  clearable
                  style="width: 250px"
              />
              <n-button @click="getPartitionOffsets" :render-icon="renderIcon(RefreshOutlined)">刷新 Offsets</n-button>
            </n-flex>
            <n-data-table
                :columns="partitions_columns"
                :data="partitions_data"
                :bordered="false"
                :pagination="pagination"
            />
          </n-flex>

        </n-tab-pane>

        <n-tab-pane name="配置">
          <template #tab>
            <n-icon>
              <SettingsRound/>
            </n-icon>
            配置
          </template>

          <n-flex vertical>
            <n-flex align="center">
              <n-text>topic: </n-text>

              <n-tag type="success">
                {{ activeConfigTopic }}
              </n-tag>

              <n-button @click="getTopicConfig(activeConfigTopic)" :render-icon="renderIcon(RefreshOutlined)">刷新
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

  <n-drawer v-model:show="showDrawer" :width="500">
    <n-drawer-content title="创建Topic配置">
      <n-form
          ref="formRef"
          :model="topic_add"
          label-placement="top"
          style="text-align: left;"
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

  <n-modal v-model:show="showModal" preset="dialog" title="添加分区">
    <n-form
        label-placement="top"
        style="text-align: left;"
    >
      <n-form-item label="添加的额外的分区数" path="addPartitionNum">
        <n-input-number v-model:value="addPartitionNum" :min="1" placeholder="添加的额外的分区数"
                        :style="{ maxWidth: '120px' }"/>
      </n-form-item>
      <n-flex>

        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" @click="addTopicPartition">确定</n-button>
      </n-flex>
    </n-form>

  </n-modal>

</template>
<script setup>
import {h, onMounted, ref} from "vue";
import emitter from "../utils/eventBus";
import {NButton, NButtonGroup, NDataTable, NIcon, NInput, NPopconfirm, NTag, NText, useMessage} from 'naive-ui'
import {
  AddFilled,
  AddRoadOutlined,
  DeleteForeverTwotone,
  DriveFileMoveTwotone,
  LibraryBooksOutlined,
  RefreshOutlined,
  SearchOutlined,
  SettingsRound
} from '@vicons/material'
import {createCsvContent, download_file, isValidJson, renderIcon} from "../utils/common";
import {
  AlterTopicConfig,
  CreatePartitions,
  CreateTopics,
  DeleteTopic,
  GetGroups,
  GetTopicConfig,
  GetTopicOffsets,
  GetTopics
} from "../../wailsjs/go/service/Service";
import ShowOrEdit from "../common/ShowOrEdit.vue";

const config_data = ref([])
const partitions_data = ref([])
const group_data = ref([])
const offsets = ref({
  start_map: {},
  end_map: {},
  commit_map: {},
})
const activeTab = ref('主题');
const selectedGroup = ref();
const loading = ref(false)
const data = ref([])
const topic_add = ref({
  topics: [],
  partitions: 1,
  replication_factor: 1,
  configs: ""
})
const rowKey = (row) => row['topic']
const selectedRowKeys = ref([]);

const message = useMessage()
const searchText = ref("");
const activeDetailTopic = ref("");
const activeConfigTopic = ref("");
const showDrawer = ref(false)
const showModal = ref(false)
const addPartitionNum = ref(1)

const selectNode = async (node) => {
  config_data.value = []
  partitions_data.value = []
  group_data.value = []
  selectedRowKeys.value = []
  data.value = []
  offsets.value = {
    start_map: {},
    end_map: {},
    commit_map: {},
  }
  topic_add.value = {
    topics: [],
    partitions: 1,
    replication_factor: 1,
    configs: ""
  }

  selectedGroup.value = null
  activeDetailTopic.value = ''
  activeConfigTopic.value = ''
  loading.value = false
  showDrawer.value = false
  showModal.value = false
  addPartitionNum.value = 1

  await getData()
  await getGroups()
}

onMounted(async () => {
  emitter.on('selectNode', selectNode)
  await getData()
  await getGroups()
})

// 读取topic及分区信息
const getData = async () => {
  loading.value = true
  try {
    const res = await GetTopics()
    if (res.err !== "") {
      message.error(res.err)
    } else {
      // 排序
      if (res.results) {
        res.results.sort((a, b) => a['topic'] > b['topic'] ? 1 : -1)
        if (searchText.value !== "") {
          // 模糊搜索
          data.value = res.results.filter(item => item['topic'].includes(searchText.value))
        } else {
          data.value = res.results
        }
      } else {
        data.value = []
      }
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
  let datas = {
    "配置":  [config_data,config_columns],
    "主题": [data,columns],
    "分区": [partitions_data,partitions_columns]
  }
  const csvContent = createCsvContent(
      datas[activeTab.value][0].value,
      datas[activeTab.value][1],
  )
  download_file(csvContent, `${activeTab.value}.csv`, 'text/csv;charset=utf-8;')
}

const columns = [
  {type: "selection",},
  // {title: 'ID', key: 'ID', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},},
  {
    title: 'topic', key: 'topic', sorter: 'default', width: 80, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
    render: (row) => h(NButton, {
      tertiary : true,
      type: "info",
      onClick: async () => {
        await getTopicConfig(row["topic"])
        activeConfigTopic.value = row["topic"]
      }
    }, {default: () => row['topic']}),
  },
  {title: '分区', key: 'partition_count', sorter: 'default', width: 30, resizable: true,
    render: (row) => h(NButton, {
      tertiary : true,
      type: "info",
      onClick: async () => {
        await getTopicDetail(row["topic"])
        activeDetailTopic.value = row["topic"]
      }
    }, {default: () => row['partition_count']}),
  },
  {title: '副本', key: 'replication_factor', sorter: 'default', width: 30, resizable: true},
  {
    title: '主题故障', key: 'Err', sorter: 'default', width: 40, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
    render: (row) => h(NTag, {type: row['Err'] === "" ? "success" : 'error'}, {default: () => row['Err'] === "" ? "健康" : row['Err']}),
  },
  {title: 'StartOffset', key: 'StartOffset', sorter: 'default', width: 50, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},},
  {title: 'CommittedOffset', key: 'CommittedOffset', sorter: 'default', width: 60, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},},
  {title: 'EndOffset', key: 'EndOffset', sorter: 'default', width: 50, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},},
  {
    title: '操作',
    key: 'actions',
    width: 80,  // 调整宽度以适应两个按钮
    resizable: true,
    render: (row) => h(
        NButtonGroup,
        {
          vertical: false,
        },
        {
          default: () => [
            h(
                NPopconfirm,
                {
                  onPositiveClick: () => deleteTopic(row["topic"])
                },
                {
                  trigger: () =>
                      row['IsInternal'] === false ?
                          h(
                              NButton,
                              {
                                strong: true,
                                secondary: true,
                                type: 'error'
                              },
                              {
                                default: () => '删除',
                                icon: () => h(NIcon, null, {default: () => h(DeleteForeverTwotone)})
                              }
                          )
                          : h(NButton, {disabled: true}, {default: () => '不可删除'})
                  ,
                  default: () => `确认删除${row["topic"]}?`
                }
            ),
          ]
        }
    )
  },
]

const partitions_columns = [
  {title: 'ID', key: 'partition', sorter: 'default', width: 10, resizable: true},
  {
    title: '分区故障', key: 'err', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
    render: (row) => h(NTag, {type: row['err'] === "" ? "success" : 'error'}, {default: () => row['err'] === "" ? "健康" : row['err']}),
  },
  {title: 'StartOffset', key: 'StartOffset', sorter: 'default', width: 15, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},},
  {title: 'CommittedOffset', key: 'CommittedOffset', sorter: 'default', width: 16, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},},
  {title: 'EndOffset', key: 'EndOffset', sorter: 'default', width: 15, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},},
  {title: 'Leader ID', key: 'leader', sorter: 'default', width: 15, resizable: true},
  {title: 'LeaderEpoch', key: 'LeaderEpoch', sorter: 'default', width: 15, resizable: true},
  {title: '托管此分区的副本ID集', key: 'replicas', sorter: 'default', width: 15, resizable: true},
  {title: 'ISR副本ID集', key: 'isr', sorter: 'default', width: 15, resizable: true},
  {title: '离线副本ID集', key: 'OfflineReplicas', sorter: 'default', width: 15, resizable: true},
]

const config_columns = [
  {title: 'Name', key: 'Name', sorter: 'default', width: 100, resizable: true},
  {
    title: 'Value（双击可编辑）', key: 'Value', sorter: 'default', width: 140, resizable: true,
    render: (row) => {
      return h(ShowOrEdit, {
        value: row['Value'],
        onUpdateValue(v) {
          alterTopicConfig(activeConfigTopic.value, row['Name'], v)
        }
      })
    }
  },
  {title: '来源', key: 'Source', sorter: 'default', width: 50, resizable: true,},
  {
    title: '是否敏感',
    key: 'Sensitive',
    width: 20,
    resizable: true,
    sorter: (row1, row2) => Number(row1['Sensitive']) - Number(row2['Sensitive']),
    render: (row) => h(NTag, {type: row['Sensitive'] === true ? "error" : "info"}, {default: () => row['Sensitive'] === true ? "是" : "否"}),
  },

]


const getTopicConfig = async (topic) => {
  loading.value = true
  try {
    const res = await GetTopicConfig(topic)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      // 排序
      if (res.results) {
        res.results.sort((a, b) => a["Name"] > b["Name"] ? 1 : -1)
        config_data.value = res.results
      } else {
        config_data.value = []
      }
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
    let partitions = data.value.find(item => item['topic'] === topic)['partitions']
    partitions.sort((a, b) => a['partition'] > b['partition'] ? 1 : -1)
    // 给每个item添加topic属性，后面匹配会用到
    partitions.forEach(item => item['topic'] = topic)
    partitions_data.value = partitions
    // 获取offsets
    if (selectedRowKeys.value.length > 0) {
      await getTopicsOffsets()
    }else {
      mergeOffsets()
    }
    activeTab.value = "分区"
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}

// 获取当页数据
const getPage =  (data_lst) => {
  const start = (pagination.value.page - 1) * pagination.value.pageSize;
  const end = start + pagination.value.pageSize;
  return data_lst.slice(start, end)
}

const getTopicsOffsets =  async () => {
  if (!selectedGroup.value){
    message.warning("请先选择 Group")
    return
  }
  const page_data = getPage(data.value)
  const topics = page_data.map(item => item['topic'])
  await getOffsets(topics, selectedGroup.value)
}

const getPartitionOffsets =  async () => {
  if (activeDetailTopic.value === ""){
    message.warning("请先从具体的topic切换到本页")
    return
  }
  if (!selectedGroup.value){
    message.warning("请先选择 Group")
    return
  }

  await getOffsets([activeDetailTopic.value], selectedGroup.value)
}

const getOffsets = async (topics, key) => {
  try {
    loading.value = true
    const res = await GetTopicOffsets(topics, key)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      offsets.value.start_map = {...res.result.start_map}
      offsets.value.end_map = {...res.result.end_map}
      offsets.value.commit_map = {...res.result.commit_map}
      mergeOffsets()
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}

// 刷新topic和分区的offset
const mergeOffsets = () => {
  // 刷新topic 列表data，把offsets塞进去
  for (const k in data.value){
    const v = data.value[k]
    const topic = v['topic']
    if (topic in offsets.value.start_map){v['StartOffset'] = addOffsets(offsets.value.start_map[topic])}
    if (topic in offsets.value.end_map){v['EndOffset'] = addOffsets(offsets.value.end_map[topic])}
    if (topic in offsets.value.commit_map){v['CommittedOffset'] = addOffsets(offsets.value.commit_map[topic])}
  }

  for (const k in partitions_data.value){
    const v = partitions_data.value[k]
    const topic = v['topic']
    const partitions_num = v['partition']
    if (topic in offsets.value.start_map){v['StartOffset'] = offsets.value.start_map[topic][partitions_num]['At']}
    if (topic in offsets.value.end_map){v['EndOffset'] = offsets.value.end_map[topic][partitions_num]['At']}
    if (topic in offsets.value.commit_map){v['CommittedOffset'] = offsets.value.commit_map[topic][partitions_num]['At']}

  }

}

const addOffsets = (item) => {

  let count = 0;
  for (const k in item) {
    const at = item[k]['At']
    if (at > 0) {
      count += at
    }
  }
  return count
}
const getGroups = async () => {
  loading.value = true
  try {
    const res = await GetGroups()
    if (res.err !== "") {
      message.error(res.err)
    } else {
      if (res.results){
        let groups = []
        for (const k in res.results) {
          const data = res.results[k]
          groups.push({
            label: data['Group'],
            value: data['Group'],
            State: data['State'],
            ProtocolType: data['ProtocolType'],
            Coordinator: data['Coordinator'],
          })
        }
        groups.sort((a, b) => a['label'] > b['label'] ? 1 : -1)
        group_data.value = groups
      }
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}

const deleteTopic = async (topic) => {
  loading.value = true
  try {
    const res = await DeleteTopic([topic])
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success("删除成功")
      await getData()
      emitter.emit('refreshTopic')
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}

const addTopic = async () => {
  loading.value = true
  try {
    let configs = {}
    if (isValidJson(topic_add.value.configs)) {
      configs = JSON.parse(topic_add.value.configs)
    }
    const res = await CreateTopics(topic_add.value.topics, topic_add.value.partitions, topic_add.value.replication_factor, configs)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success("创建成功")
      showDrawer.value = false
      topic_add.value.topics = []
      await getData()
      emitter.emit('refreshTopic')
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}
const addTopicPartition = async () => {
  loading.value = true
  try {
    const res = await CreatePartitions([activeDetailTopic.value], addPartitionNum.value)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success("添加成功")
      await getData()
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false
  showModal.value = false
  await getTopicDetail(activeDetailTopic.value)

}

const alterTopicConfig = async (topic, name, value) => {
  loading.value = true
  try {
    const res = await AlterTopicConfig(topic, name, value)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success("编辑成功，刷新配置")
      await getTopicConfig(activeConfigTopic.value)
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}
</script>


<style scoped>

</style>