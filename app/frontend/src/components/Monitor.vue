<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="max-width: 200px;">巡检</h2>
      <p>巡检Kafka积压情况。</p>
    </n-flex>
    <n-flex align="center">
      Topics:
      <n-select
          v-model:value="selectedTopics"
          @update:value="clear_offset"
          :options="topic_data"
          placeholder="选择或搜索Kafka Topic"
          filterable
          clearable
          multiple
          style="width: 600px"
      />
      Group：
      <n-select
          v-model:value="selectedGroup"
          @update:value="clear_offset"
          :options="group_data"
          placeholder="选择或创建Consumer Group"
          filterable
          clearable
          tag
          style="width: 300px"
      />
      <n-button @click="fetchData" :loading="loading" :render-icon="renderIcon(MessageOutlined)">开始巡检</n-button>
      每5分钟自动抓取一次数据

    </n-flex>

    <n-flex vertical>
<!--      <n-flex align="center">-->
<!--        <div ref="commit_chartRef" style="width: 48%; height: 500px"></div>-->
<!--        <div ref="end_chartRef" style="width: 48%; height: 500px"></div>-->
<!--      </n-flex>-->
      积压 = 终末offset - 提交offset。
      <div ref="lag_chartRef" style="width: 98%; height: 440px"></div>
      <div ref="commit_chartRef" style="width: 98%; height: 440px"></div>
      <div ref="end_chartRef" style="width: 98%; height: 440px"></div>
    </n-flex>

  </n-flex>

</template>

<script setup>
import {onMounted, ref, shallowRef} from 'vue'
import * as echarts from 'echarts/core';
import {GridComponent, LegendComponent, TitleComponent, ToolboxComponent, TooltipComponent} from 'echarts/components';
import {LineChart} from 'echarts/charts';
import {UniversalTransition} from 'echarts/features';
import {CanvasRenderer} from 'echarts/renderers';

import {NButton, NFlex, useMessage} from "naive-ui";
import {GetGroups, GetTopicOffsets, GetTopics} from "../../wailsjs/go/service/Service";
import emitter from "../utils/eventBus";
import {renderIcon} from "../utils/common";
import {MessageOutlined} from "@vicons/material";

const message = useMessage()
const topic_data = ref([]);
const group_data = ref([]);
const selectedTopics = ref([])
const selectedGroup = ref(null)

const commit_chartRef = shallowRef(null)
const lag_chartRef = shallowRef(null)
const end_chartRef = shallowRef(null)

const lag_chart = shallowRef(null)
const commit_chart = shallowRef(null)
const end_chart = shallowRef(null)

const offsetData = ref({
  lag: {},
  commit: {},
  end: {},
})
const loading = ref(false)

const selectNode = async (node) => {
  topic_data.value = []
  group_data.value = []
  selectedTopics.value = []
  selectedGroup.value = null

  lag_chartRef.value = null
  commit_chartRef.value = null
  end_chartRef.value = null

  lag_chart.value = null
  commit_chart.value = null
  end_chart.value = null

  offsetData.value = {
    lag: {},
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
  // await fetchData()
  timer = setInterval(fetchData, 5 * 60 * 1000) // 定时更新一次

  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)
})

const refreshTopic = async () => {
  await getData()
}

const handleResize = () => {
  for (const chartRef of [lag_chartRef.value, commit_chart.value, end_chart.value]) {
    if (chartRef) {
      chartRef.resize()
    }
  }
}

// 初始化图表
const initChart = () => {

  echarts.use([
    TitleComponent,
    TooltipComponent,
    GridComponent,
    LegendComponent,
    LineChart,
    CanvasRenderer,
    UniversalTransition,
    ToolboxComponent
  ]);

  const option = {
    backgroundColor: 'transparent', // 设置背景色为透明
    title: {
      text: 'Kafka Offset监控',
    },
    tooltip: {
      trigger: 'axis',
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      splitLine: {
        show: true
      },
      data: []
    },
    yAxis: {
      type: 'value'
    },
    legend: {
      data: []
    },
    series: []
  }

  const lag_option = {...option}
  lag_option.title.text = '积压量'
  lag_chart.value = echarts.init(lag_chartRef.value, 'dark')
  lag_chart.value.setOption({...lag_option})

  const commit_option = {...option}
  commit_option.title.text = '提交 offset'
  commit_chart.value = echarts.init(commit_chartRef.value, 'dark')
  commit_chart.value.setOption({...commit_option})

  const end_option = {...option}
  end_option.title.text = '终末 offset'
  end_chart.value = echarts.init(end_chartRef.value, 'dark')
  end_chart.value.setOption({...end_option})

}

// 当选择了新的topic或group后，清空之前的数据
const clear_offset = (value, option) => {
  offsetData.value = {
    lag: {},
    commit: {},
    end: {},
  }
}

// 更新图表数据
const updateChart = () => {

  const chart_map = {
    lag: lag_chart.value,
    commit: commit_chart.value,
    end: end_chart.value
  }

  // 把offsetData的数据渲染到图表上
  for (const k in offsetData.value) {
    let series = []
    let legendData = []

    // 迭代每个topic
    let xs = []
    Object.entries(offsetData.value[k]).forEach(([topic, data]) => {
      // data是对象列表，item是time和offset
      legendData.push(topic)
      series.push({
        name: topic,
        type: 'line',
        symbol: 'circle',
        data: data.map(item => item.offset)
      })
      xs = data.map(item => item.time)
    })

    chart_map[k].setOption({
      xAxis: {
        data: xs
      },
      legend: {
        data: legendData
      },
      series: series
    })
  }

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
      const now = new Date()
      const time = now.getHours() + ":" + now.getMinutes() + ":" + now.getSeconds()

      // 更新数据结构 start_map
      selectedTopics.value.forEach(topic => {

        // lag，是end-commit
        if (!offsetData.value.lag[topic]) {
          offsetData.value.lag[topic] = []
        }
        offsetData.value.lag[topic].push({
          time: time,
          offset: addOffsets(res.result.end_map[topic]) - addOffsets(res.result.commit_map[topic]) || 0
        })

        // commit
        if (!offsetData.value.commit[topic]) {
          offsetData.value.commit[topic] = []
        }
        offsetData.value.commit[topic].push({
          time: time,
          offset: addOffsets(res.result.commit_map[topic]) || 0
        })

        // end
        if (!offsetData.value.end[topic]) {
          offsetData.value.end[topic] = []
        }
        offsetData.value.end[topic].push({
          time: time,
          offset: addOffsets(res.result.end_map[topic]) || 0
        })

        // 只保留最近30个数据点
        for (const key of ["lag", "commit", "end"]){
          if (offsetData.value[key][topic].length > 100) {
            offsetData.value[key][topic].shift()
          }
        }

      })


      updateChart()
    }

  } catch (e) {
    message.error(e)
  } finally {
    loading.value = false
  }
}


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

</script>

<style scoped>
</style>