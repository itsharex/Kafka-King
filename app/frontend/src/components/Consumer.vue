<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="max-width: 200px;">Consumer</h2>
      <p>一个简单消费者客户端，查看Topic消息。</p>
    </n-flex>
    <!-- 查询条件区域 -->
    <n-flex align="center">
      必选：Topic：
      <n-select
          v-model:value="selectedTopic"
          :options="topic_data"
          placeholder="选择或搜索Kafka Topic"
          filterable
          clearable
          style="width: 300px"
      />

      必选：消费数量：
      <n-input-number
          v-model:value="maxMessages"
          :min="1"
          placeholder="消费消息数量"
          style="width: 160px"
      />
      poll超时时间：默认10s。如异常或无可消费消息，则会超时
      <n-input-number
          v-model:value="timeout"
          :min="1"
          placeholder="poll超时时间"
          style="width: 160px"
      />
    </n-flex>
    <n-flex align="center">
      可选：Group（一旦选择，消费时会自动提交Offset。支持创建新Group）
      <n-select
          v-model:value="selectedGroup"
          :options="group_data"
          placeholder="选择或创建Consumer Group"
          filterable
          clearable
          tag
          style="width: 300px"
      />
      <n-button @click="consume" :loading="loading" :render-icon="renderIcon(MessageOutlined)">消费消息</n-button>
    </n-flex>

    <!-- 消息列表 -->
    <n-data-table
        :columns="columns"
        :data="messages"
        :pagination="pagination"
        :bordered="true"
        striped
    />
  </n-flex>
</template>

<script setup>
import {onMounted, ref} from 'vue'
import emitter from "../utils/eventBus";
import {renderIcon} from "../utils/common";
import {MessageOutlined} from "@vicons/material";
import {NButton, NDataTable, NFlex, useMessage} from 'naive-ui'
import {Consumer, GetGroups, GetTopics} from "../../wailsjs/go/service/Service";

const message = useMessage()
const topic_data = ref([]);
const group_data = ref([]);
const messages = ref([])

// 表单数据
const selectedTopic = ref()
const selectedGroup = ref()
const maxMessages = ref(10)
const timeout = ref(10)
const loading = ref(false)

const selectNode = async (node) => {
  topic_data.value = []
  group_data.value = []
  messages.value = []
  selectedTopic.value = null
  selectedGroup.value = null
  loading.value = false
  await getData()
}

onMounted(async () => {
  emitter.on('selectNode', selectNode)
  await getData()
})


const getData = async () => {
  console.log('初始化消费者数据')
  try {
    const res = await GetTopics()
    const res2 = await GetGroups()
    if (res.err !== "" || res2.err !== "") {
      message.error(res.err === res2.err? res.err : res.err + res2.err)
    } else {
      let topic_data_lst = []
      if (res.results) {
        res.results.sort((a, b) => a['topic'] > b['topic'] ? 1 : -1)
        for (const k in res.results) {
          topic_data_lst.push({
            "label": res.results[k]['topic'],
            "value": res.results[k]['topic']
          })
        }
      }
      topic_data.value = topic_data_lst

      let groups = []
      for (const k in res2.results) {
        const g_data = res2.results[k]
        groups.push({
          label: g_data['Group'],
          value: g_data['Group'],
          State: g_data['State'],
          ProtocolType: g_data['ProtocolType'],
          Coordinator: g_data['Coordinator'],
        })
      }
      groups.sort((a, b) => a['label'] > b['label'] ? 1 : -1)
      group_data.value = groups

    }
  } catch (e) {
    message.error(e)
  }
}

// 分页配置
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


// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'ID',
    width: 10,
  },
  {
    title: 'Offset',
    key: 'Offset',
    width: 20,
    ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
  },
  {
    title: 'Key',
    key: 'Key',
    width: 20,
    ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
  },
  {
    title: '消息内容',
    key: 'Value',
    width: 40,
    ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
  },
  {
    title: '时间戳',
    key: 'Timestamp',
    width: 20,
    ellipsis: {
      tooltip: true
    },
  },
  {
    title: 'Partition',
    key: 'Partition',
    width: 20,
  },
  {
    title: 'Topic',
    key: 'Topic',
    width: 20,
    ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
  },
  {
    title: 'Headers',
    key: 'Headers',
    width: 20,
    ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
  },
  {
    title: 'LeaderEpoch',
    key: 'LeaderEpoch',
    width: 20,
  },
  {
    title: 'ProducerEpoch',
    key: 'ProducerEpoch',
    width: 20,
  },
  {
    title: 'ProducerID',
    key: 'ProducerID',
    width: 20,
  }
]


// 获取消息
const consume = async () => {
  if (!selectedTopic.value) {
    message.error('请选择Topic')
    return
  }

  loading.value = true
  try {
    const result = await Consumer(selectedTopic.value, selectedGroup.value, maxMessages.value, timeout.value)
    if (result.err !== "") {
      message.error(result.err)
    } else {
      messages.value = result.results
      message.success('获取成功')
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 保存为文本文件
const saveMessageAsString = (message) => {
  const content = message.content
  const blob = new Blob([content], {type: 'text/plain;charset=utf-8'})
  // saveAs(blob, `message-${message.offset}.txt`)
}

// 保存为二进制文件
const saveMessageAsBinary = (message) => {
  // 这里假设message.content是base64编码的二进制数据
  const content = atob(message.content)
  const bytes = new Uint8Array(content.length)
  for (let i = 0; i < content.length; i++) {
    bytes[i] = content.charCodeAt(i)
  }
  const blob = new Blob([bytes], {type: 'application/octet-stream'})
  // saveAs(blob, `message-${message.offset}.bin`)
}
</script>