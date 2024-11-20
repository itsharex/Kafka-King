<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="max-width: 200px;">巡检</h2>
      <p>巡检Kafka积压情况。</p>
    </n-flex>

    <n-flex align="center">
      必选：Topic，支持多选：
      <n-select
          v-model:value="selectedTopics"
          :options="topic_data"
          placeholder="选择或搜索Kafka Topic"
          filterable
          clearable
          multiple
          style="width: 600px"
      />
      可选：Group
      <n-select
          v-model:value="selectedGroup"
          :options="group_data"
          placeholder="选择或创建Consumer Group"
          filterable
          clearable
          tag
          style="width: 300px"
      />
      <n-button @click="fetchData" :loading="loading" :render-icon="renderIcon(MessageOutlined)">开始巡检</n-button>

    </n-flex>

    <n-flex vertical>

      <n-flex align="center">
      </n-flex>
      <!-- 图表容器 -->
      <div ref="start_chartRef" style="width: 100%; height: 500px"></div>
      <div ref="commit_chartRef" style="width: 100%; height: 500px"></div>
      <div ref="end_chartRef" style="width: 100%; height: 500px"></div>

      <n-flex align="center">
      </n-flex>

    </n-flex>

  </n-flex>

</template>

<script setup>
import {onMounted, ref} from 'vue'
import * as echarts from 'echarts'
import {NButton, NFlex, useMessage} from "naive-ui";
import {GetGroups, GetTopicOffsets, GetTopics} from "../../wailsjs/go/service/Service";
import emitter from "../utils/eventBus";
import {renderIcon} from "../utils/common";
import {MessageOutlined} from "@vicons/material";
import {EventsOn} from "../../wailsjs/runtime";

const message = useMessage()
const topic_data = ref([]);
const group_data = ref([]);
const selectedTopics = ref([])
const selectedGroup = ref(null)

const start_chartRef = ref(null)
const commit_chartRef = ref(null)
const end_chartRef = ref(null)
const start_chart = ref(null)
const commit_chart = ref(null)
const end_chart = ref(null)

const offsetData = ref({
  start: {},
  commit: {},
  end: {},
})
const loading = ref(false)

const selectNode = async (node) => {
  topic_data.value = []
  group_data.value = []
  selectedTopics.value = []
  selectedGroup.value = null

  start_chartRef.value = null
  commit_chartRef.value = null
  end_chartRef.value = null
  start_chart.value = null
  commit_chart.value = null
  end_chart.value = null

  offsetData.value = {
    start: {},
    commit: {},
    end: {},
  }
  loading.value = false

  await getData()
}

onMounted(async () => {
  emitter.on('selectNode', selectNode)
  emitter.on('refreshTopic', refreshTopic)

  await getData()
  initChart()
  await fetchData()
  timer = setInterval(fetchData, 2 * 60 * 1000) // 定时更新一次
  EventsOn('resize', () => {
    start_chart.value?.resize()
    commit_chart.value?.resize()
    end_chart.value?.resize()
  })
})

const refreshTopic = async () => {
  await getData()
}


// 初始化图表
const initChart = () => {
  if (start_chart.value) {
    start_chart.value.dispose()
  }

  start_chart.value = echarts.init(start_chartRef.value, 'dark')
  const option = {
    title: {
      text: 'Kafka Offset监控',
      left: 'left'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      },
      formatter: function(params) {
        const time = new Date(params[0].axisValue).toLocaleString()
        let result = `${time}<br/>`
        params.forEach(param => {
          const color = param.color
          result += `${param.seriesName}: ${param.value}<br/>`
        })
        return result
      }
    },
    xAxis: {
      type: 'time',
      boundaryGap: false,
    },
    grid: {
      left: '10%',
      right: '20%',  // 为图例留出空间
      bottom: '10%',
      containLabel: true
    },
    yAxis: {
      type: 'value',
      boundaryGap: [0, '100%']
    },
    legend: {
      orient: 'vertical',  // 图例垂直显示
      right: '10%',        // 图例位置在右侧
      top: 'center'        // 图例垂直居中
    },
    series: []
  }
  option.title.text='Start Offset监控'
  start_chart.value.setOption(option)

  option.title.text='Committed Offset监控'
  commit_chart.value = echarts.init(commit_chartRef.value, 'dark')
  commit_chart.value.setOption(option)

  option.title.text='End Offset监控'
  end_chart.value = echarts.init(end_chartRef.value, 'dark')
  end_chart.value.setOption(option)

}

// 更新图表数据
const updateChart = () => {
  if (!start_chart.value) return

  const series = []
  const legendData = []

  // start
  Object.entries(offsetData.value.start).forEach(([key, data]) => {
    legendData.push(key)
    series.push({
      name: key,
      type: 'line',
      symbol: 'none',
      sampling: 'lttb',
      itemStyle: {
        color: 'rgb(255, 70, 131)'
      },
      data: data.map(item => [item.timestamp, item.offset]),
    })
  })

  start_chart.value.setOption({
    legend: {
      data: legendData
    },
    series: series
  })

  // commit
  Object.entries(offsetData.value.commit).forEach(([key, data]) => {
    legendData.push(key)
    series.push({
      name: key,
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 8,
      data: data.map(item => [item.timestamp, item.offset]),
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowColor: 'rgba(0,0,0,0.3)'
        }
      },
    })
  })
  commit_chart.value.setOption({
    legend: {
      data: legendData
    },
    series: series
  })

  // end
  Object.entries(offsetData.value.end).forEach(([key, data]) => {
    legendData.push(key)
    series.push({
      name: key,
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 8,
      data: data.map(item => [item.timestamp, item.offset]),
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowColor: 'rgba(0,0,0,0.3)'
        }
      },
    })
  })
  end_chart.value.setOption({
    legend: {
      data: legendData
    },
    series: series
  })
}

// 定时获取数据
let timer = null
const fetchData = async () => {
  if (selectedTopics.value.length === 0 || !selectedGroup.value) {
    message.warning('请选择Topic和Group')
    return
  }
  loading.value = true

  try {
    const res = await GetTopicOffsets(selectedTopics.value, selectedGroup.value)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      const timestamp = new Date().getTime()

      // 更新数据结构 start_map
      selectedTopics.value.forEach(topic => {

        if (!offsetData.value.start[topic]) {
          offsetData.value.start[topic] = []
        }
        offsetData.value.start[topic].push({
          timestamp,
          offset: addOffsets(res.result.start_map[topic])|| 0
        })

        if (!offsetData.value.commit[topic]) {
          offsetData.value.commit[topic] = []
        }
        offsetData.value.commit[topic].push({
          timestamp,
          offset: addOffsets(res.result.commit_map[topic])|| 0
        })

        if (!offsetData.value.end[topic]) {
          offsetData.value.end[topic] = []
        }
        offsetData.value.end[topic].push({
          timestamp,
          offset: addOffsets(res.result.end_map[topic])|| 0
        })


        // 只保留最近30个数据点
        if (offsetData.value.start[topic].length > 30) {
          offsetData.value.start[topic].shift()
        }
        // 只保留最近30个数据点
        if (offsetData.value.commit[topic].length > 30) {
          offsetData.value.commit[topic].shift()
        }
        // 只保留最近30个数据点
        if (offsetData.value.end[topic].length > 30) {
          offsetData.value.end[topic].shift()
        }

      })

      console.log(offsetData.value)

      updateChart()
    }

  } catch (e) {
      message.error(e)
  }finally {
    loading.value = false
  }
}


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

// 处理选择变化
const handleTopicChange = () => {
  offsetData.value = {}
  updateChart()
}

const handleGroupChange = () => {
  offsetData.value = {}
  updateChart()
}

</script>

<style scoped>
.kafka-monitor {
  padding: 20px;
}
</style>