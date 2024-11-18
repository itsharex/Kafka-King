<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="max-width: 200px;">巡检</h2>
    </n-flex>

  <div class="kafka-monitor">
    <n-card title="Kafka Offset Monitor">
      <n-space vertical>
        <!-- Topic选择器 -->
        <n-select
            v-model:value="selectedTopics"
            multiple
            placeholder="请选择Topics"
            :options="topicOptions"
            @update:value="handleTopicChange"
        />

        <!-- 消费组选择器 -->
        <n-select
            v-model:value="selectedGroups"
            multiple
            placeholder="请选择消费组"
            :options="groupOptions"
            @update:value="handleGroupChange"
        />

        <!-- 图表容器 -->
        <div ref="chartRef" style="width: 100%; height: 400px"></div>
      </n-space>
    </n-card>
  </div>
  </n-flex>

</template>

<script setup>
import {onMounted, onUnmounted, ref} from 'vue'
import * as echarts from 'echarts'
// import { FetchTopics, FetchGroups, FetchOffsets } from '../wailsjs/go/main/App'

const chartRef = ref(null)
const chart = ref(null)
const selectedTopics = ref([])
const selectedGroups = ref([])
const topicOptions = ref([])
const groupOptions = ref([])
const offsetData = ref({})

// 初始化主题和消费组选项
const initOptions = async () => {
  try {
    const topics = await FetchTopics()
    const groups = await FetchGroups()

    topicOptions.value = topics.map(topic => ({
      label: topic,
      value: topic
    }))

    groupOptions.value = groups.map(group => ({
      label: group,
      value: group
    }))
  } catch (error) {
    console.error('Failed to fetch options:', error)
  }
}

// 初始化图表
const initChart = () => {
  if (chart.value) {
    chart.value.dispose()
  }

  chart.value = echarts.init(chartRef.value)
  const option = {
    title: {
      text: 'Kafka Offset监控',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      formatter: function(params) {
        let result = params[0].axisValue + '<br/>'
        params.forEach(param => {
          result += `${param.seriesName}: ${param.value}<br/>`
        })
        return result
      }
    },
    legend: {
      data: [],
      top: 30
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'time',
      boundaryGap: false
    },
    yAxis: {
      type: 'value',
      name: 'Offset'
    },
    series: []
  }

  chart.value.setOption(option)
}

// 更新图表数据
const updateChart = () => {
  if (!chart.value) return

  const series = []
  const legendData = []

  Object.entries(offsetData.value).forEach(([key, data]) => {
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
      }
    })
  })

  chart.value.setOption({
    legend: {
      data: legendData
    },
    series: series
  })
}

// 定时获取数据
let timer = null
const fetchData = async () => {
  if (selectedTopics.value.length === 0 || selectedGroups.value.length === 0) return

  try {
    const offsets = await FetchOffsets(selectedTopics.value, selectedGroups.value)
    const timestamp = new Date().getTime()

    // 更新数据结构
    selectedTopics.value.forEach(topic => {
      selectedGroups.value.forEach(group => {
        const key = `${topic}-${group}`
        if (!offsetData.value[key]) {
          offsetData.value[key] = []
        }

        offsetData.value[key].push({
          timestamp,
          offset: offsets[topic]?.[group] || 0
        })

        // 只保留最近30个数据点
        if (offsetData.value[key].length > 30) {
          offsetData.value[key].shift()
        }
      })
    })

    updateChart()
  } catch (error) {
    console.error('Failed to fetch offsets:', error)
  }
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

// 生命周期钩子
onMounted(async () => {
  // await initOptions()
  // initChart()
  // fetchData()
  // timer = setInterval(fetchData, 5 * 60 * 1000) // 每5分钟更新一次

  // 响应窗口大小变化
  window.addEventListener('resize', () => {
    chart.value?.resize()
  })
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
  chart.value?.dispose()
})
</script>

<style scoped>
.kafka-monitor {
  padding: 20px;
}
</style>