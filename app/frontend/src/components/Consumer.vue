<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="max-width: 200px;">Consumer</h2>
      <p>一个简单消费者客户端，查看Topic消息。</p>
    </n-flex>
    <!-- 查询条件区域 -->
    <n-flex align="center">
      选择 Topic：
      <n-select
          v-model:value="selectedTopic"
          :options="topic_data"
          placeholder="选择或搜索Kafka Topic"
          filterable
          clearable
          style="width: 300px"
      />
      选择 Group：
      <n-select
          v-model:value="selectedGroup"
          :options="group_data"
          placeholder="选择或搜索Consumer Group"
          filterable
          clearable
          style="width: 300px"
      />
      消费消息数量：
      <n-input-number
          v-model:value="maxMessages"
          :min="1"
          placeholder="消费消息数量"
          style="width: 160px"
      />
      起始Offset：
      <n-input-number
          v-model:value="startOffset"
          :min="0"
          placeholder="起始Offset"
          style="width: 160px"
      />
    </n-flex>
    <n-flex align="center">
      <n-button @click="consume" :loading="loading" :render-icon="renderIcon(SendTwotone)">拉取消息</n-button>
    </n-flex>

    <!-- 进度显示 -->
    <n-flex vertical v-if="messages.length">
      <n-progress
          type="line"
          :percentage="progressPercentage"
          :indicator-placement="'inside'"
      >
        已获取 {{ messages.length }}/{{ maxMessages }} 条消息
      </n-progress>
    </n-flex>

    <!-- 消息列表 -->
    <n-data-table
        :columns="columns"
        :data="messages"
        :pagination="pagination"
        :bordered="true"
        striped
    />
    <!-- 消息详情抽屉 -->
    <n-drawer v-model:show="showDrawer" :width="500">
      <n-drawer-content>
        <template #header>
          消息详情
        </template>
        <n-flex vertical>
          <n-descriptions bordered>
            <n-descriptions-item label="Offset">
              {{ selectedMessage?.offset }}
            </n-descriptions-item>
            <n-descriptions-item label="Key">
              {{ selectedMessage?.key }}
            </n-descriptions-item>
            <n-descriptions-item label="Timestamp">
              {{ formatTimestamp(selectedMessage?.timestamp) }}
            </n-descriptions-item>
          </n-descriptions>
          <n-card title="Headers">
            <n-descriptions bordered>
              <n-descriptions-item
                  v-for="(value, key) in selectedMessage?.headers"
                  :key="key"
                  :label="key"
              >
                {{ value }}
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
          <n-card title="Content">
            <pre style="white-space: pre-wrap; word-wrap: break-word;">{{ selectedMessage?.content }}</pre>
          </n-card>
        </n-flex>
      </n-drawer-content>
    </n-drawer>

  </n-flex>
</template>

<script setup>
import {onMounted, ref, computed, h} from 'vue'
import emitter from "../utils/eventBus";
import {renderIcon} from "../utils/common";
import {SendTwotone} from "@vicons/material";
import {NButton, NButtonGroup, NDataTable, NFlex, NIcon, NPopconfirm, NTag, NText, useMessage} from 'naive-ui'
import {GetGroups, GetTopics} from "../../wailsjs/go/service/Service";
import Consumer from "./Consumer.vue";

const message = useMessage()
const topic_data = ref([]);
const group_data = ref([]);

// 表单数据
const selectedTopic = ref('')
const selectedGroup = ref('')
const maxMessages = ref(100)
const startOffset = ref(null)
const loading = ref(false)
const messages = ref([])
const showDrawer = ref(false)
const selectedMessage = ref(null)

const selectNode = async (node) => {
}

onMounted(async () => {
  emitter.on('selectNode', selectNode)
  await getData()
})


const getData = async () => {
  try {
    const res = await GetTopics()
    const res2 = await GetGroups()
    if (res.err !== "" || res2.err !== "") {
      message.error(res.err + res2.err)
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
          key: g_data['Group'],
          State: g_data['State'],
          ProtocolType: g_data['ProtocolType'],
          Coordinator: g_data['Coordinator'],
        })
      }
      groups.sort((a, b) => a['label'] > b['label'] ? 1 : -1)
      group_data.value = groups

    }
  }catch (e) {
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

// 进度百分比
const progressPercentage = computed(() => {
  return Math.floor((messages.value.length / maxMessages.value) * 100)
})

// 格式化时间戳
const formatTimestamp = (timestamp) => {
  if (!timestamp) return ''
  return new Date(timestamp).toLocaleString()
}

// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'ID',
    width: 10,
  },
  {
    title: 'Offset',
    key: 'offset',
    width: 20,
  },
  {
    title: 'Key',
    key: 'Key',
    width: 20,
  },
  {
    title: '消息内容',
    key: 'Value',
    width: 20,
    ellipsis: {
      tooltip: true
    },
    render(row) {
      const content = row.content || ''
      return h(
          'div',
          {
            style: {
              cursor: 'pointer',
            },
            onClick: () => showMessageDetail(row)
          },
          content.length > 100 ? content.substring(0, 100) + '...' : content
      )
    }
  },
  {
    title: '时间戳',
    key: 'Timestamp',
    width: 20,
    render(row) {
      return formatTimestamp(row.timestamp)
    }
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
  },
  {
    title: 'Headers',
    key: 'Headers',
    width: 20,
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
  },
  {
    title: '操作',
    key: 'actions',
    width: 20,
    render(row) {
      return h(
          NFlex,
          null,
          {
            default: () => [
              h(
                  NButton,
                  {
                    size: 'small',
                    onClick: () => saveMessageAsString(row)
                  },
                  {default: () => '保存为文本'}
              ),
              h(
                  NButton,
                  {
                    size: 'small',
                    onClick: () => saveMessageAsBinary(row)
                  },
                  {default: () => '保存为二进制'}
              )
            ]
          }
      )
    }
  }
]

// 显示消息详情
const showMessageDetail = (message) => {
  selectedMessage.value = message
  showDrawer.value = true
}

// 获取消息
const consume = async () => {
  if (!selectedTopic.value) {
    message.error('请选择Topic')
    return
  }

  loading.value = true
  try {
    const result = await Consumer(selectedTopic.value,selectedGroup.value,maxMessages.value,)
    if (result.err !== "") {
      message.error(result.err)
    }else {
      messages.value = result
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