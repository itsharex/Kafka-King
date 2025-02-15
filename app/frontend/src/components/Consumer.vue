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
      <h2>{{ t('consumer.title') }}</h2>
      <p>{{ t('consumer.desc') }}</p>
      <n-button @click="downloadAllDataCsv" :render-icon="renderIcon(DriveFileMoveTwotone)">{{ t('common.csv') }}
      </n-button>
    </n-flex>
    <!-- 查询条件区域 -->
    <n-form
        ref="formRef"
        :model="select"
        inline
        label-placement="top"
        style="text-align: left;"
        label-width="auto"
        :rules="{
              selectedTopic: {required: true, trigger: 'blur'},
              maxMessages: {required: true, type: 'number',trigger: 'blur'},
            }"
    >

      <n-form-item label="Topic" path="selectedTopic">
        <n-select
            v-model:value="select.selectedTopic"
            :options="topic_data"
            :placeholder="t('consumer.requiredTopic')"
            filterable
            clearable
        />
      </n-form-item>
      <n-form-item label="Num" path="maxMessages" style="width: 100px">
        <n-tooltip>
          <template #trigger>
            <n-input-number
                v-model:value="select.maxMessages"
                :min="1"
            />
          </template>
          {{ t('consumer.messagesCountPlaceholder') }}
        </n-tooltip>
      </n-form-item>

      <n-form-item label="ConsumerGroup" path="selectedGroup">
        <n-tooltip>
          <template #trigger>
            <n-select
                v-model:value="select.selectedGroup"
                :options="group_data"
                filterable
                clearable
                tag
            />
          </template>
          {{ t('consumer.optionalGroup') }}
        </n-tooltip>
      </n-form-item>

      <n-form-item label="Timeout">
        <n-tooltip>
          <template #trigger>
            <n-input-number
                v-model:value="select.timeout"
                :min="1"
                style="max-width: 100px"
            />
          </template>
          {{ t('consumer.pollTimeoutPlaceholder') }}
        </n-tooltip>
      </n-form-item>

      <n-form-item label="Decompress" path="decompress">
        <n-select
            v-model:value="select.decompress"
            :options="[
              {label: 'gzip', value: 'gzip'},
            ]"
            filterable
            clearable
            style="width: 100px"
        />
      </n-form-item>

      <n-form-item label="CommitOffset" path="isCommit">
        <n-tooltip>
          <template #trigger>
            <n-switch :round="false" :checked-value=true :unchecked-value=false v-model:value="select.isCommit"/>
          </template>
          {{ t('consumer.commitOffsetTooltip') }}
        </n-tooltip>
      </n-form-item>


      <n-form-item>

        <n-button tertiary type="primary" @click="consume" :loading="loading"
                  :render-icon="renderIcon(MessageOutlined)">
          {{ t('consumer.consumeMessage') }}
        </n-button>
      </n-form-item>


    </n-form>
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
import {createCsvContent, download_file, renderIcon} from "../utils/common";
import {DriveFileMoveTwotone, MessageOutlined} from "@vicons/material";
import {NButton, NDataTable, NFlex, useMessage} from 'naive-ui'
import {Consumer, GetGroups, GetTopics} from "../../wailsjs/go/service/Service";
import {useI18n} from "vue-i18n";

const {t} = useI18n()
const formRef = ref(null)

const message = useMessage()
const topic_data = ref([]);
const group_data = ref([]);
const messages = ref([])

// 表单数据
const select = ref({
  selectedTopic: null,
  selectedGroup: null,
  maxMessages: 10,
  timeout: 10,
  isCommit: false,
  decompress: null,
})

const loading = ref(false)

const refreshTopic = async () => {
  await getData()
}
const selectNode = async (node) => {
  topic_data.value = []
  group_data.value = []
  messages.value = []
  select.value.selectedTopic = null
  select.value.selectedGroup = null
  loading.value = false
  await getData()
}

onMounted(async () => {
  emitter.on('selectNode', selectNode)
  emitter.on('refreshTopic', refreshTopic)
  await getData()
})


const getData = async () => {
  console.log('初始化消费者数据')
  try {
    const res = await GetTopics()
    const res2 = await GetGroups()
    if (res.err !== "" || res2.err !== "") {
      message.error(res.err === res2.err ? res.err : res.err + res2.err)
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
    title: 'Offset',
    key: 'Offset',
    width: 20,
    ellipsis: {tooltip: {style: {maxWidth: '800px'},}},
    sorter: 'default'
  },
  {
    title: 'Key',
    key: 'Key',
    width: 15,
    resizable: true,
    ellipsis: {tooltip: {style: {maxWidth: '800px'},}},
    sorter: 'default'
  },
  {
    title: 'Value',
    key: 'Value',
    width: 40,
    resizable: true,
    ellipsis: {tooltip: {style: {maxWidth: '800px'},}},
    sorter: 'default'
  },
  {
    title: 'Timestamp',
    key: 'Timestamp',
    width: 20,
    resizable: true,
    ellipsis: {tooltip: {style: {maxWidth: '800px'},}},
    sorter: (rowA, rowB) => {
      const dateA = new Date(rowA['Timestamp']);
      const dateB = new Date(rowB['Timestamp']);
      return dateA - dateB;
    }
  },
  {
    title: 'Topic',
    key: 'Topic',
    width: 20,
    ellipsis: {tooltip: {style: {maxWidth: '800px'},}},
    resizable: true,
    sorter: 'default'
  },
  {
    title: 'Partition',
    key: 'Partition',
    width: 10,
    resizable: true,
    sorter: 'default'
  },
  {
    title: 'Headers',
    key: 'Headers',
    width: 20,
    resizable: true,
    ellipsis: {tooltip: {style: {maxWidth: '800px'},}},
    sorter: 'default'
  },
  // {
  //   title: 'LeaderEpoch',
  //   key: 'LeaderEpoch',
  //   width: 10,
  // },
  // {
  //   title: 'ProducerEpoch',
  //   key: 'ProducerEpoch',
  //   width: 10,
  // },
  {
    title: 'ProducerID',
    key: 'ProducerID',
    width: 10,
    resizable: true,
    sorter: 'default'
  }
]


// 获取消息
const consume = async () => {
  if (!select.value.selectedTopic) {
    message.error(t('message.selectTopic'))
    return
  }

  loading.value = true
  try {
    const result = await Consumer(select.value.selectedTopic, select.value.selectedGroup,
        select.value.maxMessages, select.value.timeout, select.value.decompress,
        select.value.isCommit)
    if (result.err !== "") {
      message.error(result.err)
    } else {
      messages.value = result.results
      message.success(t('message.fetchSuccess'))
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const downloadAllDataCsv = async () => {
  const csvContent = createCsvContent(
      messages.value, columns
  )
  download_file(csvContent, 'messages.csv', 'text/csv;charset=utf-8;')
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