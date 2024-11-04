<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="max-width: 200px;">Producer</h2>
      <p>一个生产者客户端，将消息推送到指定的Topic。</p>
    </n-flex>

    <n-flex align="center">
      选择Topic：
      <n-select
          v-model:value="selectedTopic"
          :options="data"
          placeholder="必选：选择或搜索Kafka Topic"
          filterable
          clearable
          style="width: 300px"
      />
      可选：输入消息Key：
      <n-input v-model:value="messageKey" placeholder="可选：输入消息Key" style="width: 300px"/>
      可选：指定分区号：
      <n-input-number v-model:value="partition" placeholder="" style="width: 120px"/>
    </n-flex>

    <n-input
        v-model:value="messageContent"
        type="textarea"
        placeholder="必填：消息内容，字符串格式，支持JSON"
        :rows="12"
        style="text-align: left;"
    />

    <!--    动态添加header-->
    <n-flex vertical>
      <n-flex align="center">
        <span>消息Headers:</span>
        <n-button size="small" @click="addHeader">
          添加Header
        </n-button>
      </n-flex>

      <n-flex v-for="(header, index) in headers" :key="index" :wrap="false">
        <n-input v-model:value="header.key" placeholder="Header Key" style="width: 200px"/>
        <n-input v-model:value="header.value" placeholder="Header Value" style="width: 200px"/>
        <n-button size="small" type="error" @click="removeHeader(index)">
          删除
        </n-button>
      </n-flex>
    </n-flex>

    <n-flex align="center">
      发送次数
      <n-input-number
          v-model:value="nums"
          :min="1"
          placeholder="发送次数"
          style="width: 160px"
      />

    </n-flex>
    <n-flex align="center">
      <n-button @click="produce" :loading="loading" :render-icon="renderIcon(SendTwotone)">发送消息</n-button>
    </n-flex>

  </n-flex>
</template>

<script setup>

import {onMounted, ref} from 'vue'
import emitter from "../utils/eventBus";
import {useMessage} from "naive-ui";
import {renderIcon} from "../utils/common";
import {SendTwotone} from "@vicons/material";
import {GetTopics, Produce} from "../../wailsjs/go/service/Service";

const message = useMessage()
const data = ref([]);
// 表单数据
const selectedTopic = ref()
const messageKey = ref()
const messageContent = ref('')
const headers = ref([])
const nums = ref(1)
const partition = ref(0)
const loading = ref(false)

const selectNode = async (node) => {
}

onMounted(async () => {
  emitter.on('selectNode', selectNode)
  await getData()
})

// 读取topic及分区信息
const getData = async () => {
  try {
    const res = await GetTopics()
    if (res.err !== "") {
      message.error(res.err)
    } else {
      // 排序
      let data_lst = []
      if (res.results) {
        res.results.sort((a, b) => a['topic'] > b['topic'] ? 1 : -1)
        for (const k in res.results) {
          data_lst.push({
            "label": res.results[k]['topic'],
            "value": res.results[k]['topic']
          })
        }
      }
      data.value = data_lst
    }
  } catch (e) {
    message.error(e)
  }

}


// 添加header
const addHeader = () => {
  headers.value.push({
    key: '',
    value: ''
  })
}

// 删除header
const removeHeader = (index) => {
  headers.value.splice(index, 1)
}

// 发送消息
const produce = async () => {
  if (!selectedTopic.value) {
    message.error('请选择Topic')
    return
  }

  if (!messageContent.value) {
    message.error('请输入消息内容')
    return
  }
  loading.value = true
  try {
    const res = await Produce(selectedTopic.value, messageKey.value, messageContent.value, partition.value, nums.value, headers.value)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success('发送成功')
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}


</script>
