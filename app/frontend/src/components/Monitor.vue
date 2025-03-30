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
      <h2>{{ t('inspection.title') }}</h2>
      <p>{{ t('inspection.desc') }}</p>
    </n-flex>
    <n-flex align="center">
      {{ t('inspection.topicsLabel') }}:
      <n-select
          v-model:value="selectedTopics"
          @update:value="clear_offset"
          :options="topic_data"
          :placeholder="t('inspection.topicPlaceholder')"
          filterable
          clearable
          multiple
          style="width: 600px"
      />
      {{ t('inspection.groupLabel') }}：
      <n-select
          v-model:value="selectedGroup"
          @update:value="clear_offset"
          :options="group_data"
          :placeholder="t('inspection.groupPlaceholder')"
          filterable
          clearable
          tag
          style="width: 300px"
      />
      <n-button @click="fetchData" :loading="loading" :render-icon="renderIcon(MessageOutlined)">
        {{ t('inspection.startInspection') }}
      </n-button>
      {{ t('inspection.autoFetch') }}
    </n-flex>

    <n-flex vertical>
      <n-flex align="center">
        <div ref="lag_chartRef" style="width: 48%; height: 440px"></div>
        <div ref="commit_chartRef" style="width: 48%; height: 440px"></div>
      </n-flex>
      <n-flex align="center">
        <div ref="end_chartRef" style="width: 48%; height: 440px"></div>
        <div ref="productionSpeed_chartRef" style="width: 48%; height: 440px"></div>
      </n-flex>
      <n-flex align="center">
        <div ref="consumptionSpeed_chartRef" style="width: 48%; height: 440px"></div>
      </n-flex>
    </n-flex>
  </n-flex>
</template>

<script setup>
import { onMounted, ref, shallowRef } from 'vue';
import * as echarts from 'echarts/core';
import { GridComponent, LegendComponent, TitleComponent, ToolboxComponent, TooltipComponent } from 'echarts/components';
import { LineChart } from 'echarts/charts';
import { UniversalTransition } from 'echarts/features';
import { CanvasRenderer } from 'echarts/renderers';
import { lightTheme, NButton, NFlex, useMessage } from 'naive-ui';
import { GetGroups, GetTopicOffsets, GetTopics } from '../../wailsjs/go/service/Service';
import emitter from '../utils/eventBus';
import { renderIcon } from '../utils/common';
import { MessageOutlined } from '@vicons/material';
import { useI18n } from 'vue-i18n';
import { GetConfig } from '../../wailsjs/go/config/AppConfig';

const { t } = useI18n();
const message = useMessage();
const topic_data = ref([]);
const group_data = ref([]);
const selectedTopics = ref([]);
const selectedGroup = ref(null);

const lag_chartRef = shallowRef(null);
const commit_chartRef = shallowRef(null);
const end_chartRef = shallowRef(null);
const productionSpeed_chartRef = shallowRef(null);
const consumptionSpeed_chartRef = shallowRef(null);

const lag_chart = shallowRef(null);
const commit_chart = shallowRef(null);
const end_chart = shallowRef(null);
const productionSpeed_chart = shallowRef(null);
const consumptionSpeed_chart = shallowRef(null);

const offsetData = ref({
  lag: {},
  commit: {},
  end: {},
  productionSpeed: {},
  consumptionSpeed: {},
});
const loading = ref(false);

let echarts_theme = 'dark';

onMounted(async () => {
  emitter.on('selectNode', selectNode);
  emitter.on('refreshTopic', refreshTopic);

  const loadedConfig = await GetConfig();
  echarts_theme = loadedConfig.theme === lightTheme.name ? 'light' : 'dark';

  await getData();
  initChart();
  setInterval(fetchData, 5 * 60 * 1000); // Update every 5 minutes

  window.addEventListener('resize', handleResize);
});

const refreshTopic = async () => {
  await getData();
};

const handleResize = () => {
  [lag_chart, commit_chart, end_chart, productionSpeed_chart, consumptionSpeed_chart].forEach(chart => {
    if (chart.value) chart.value.resize();
  });
};

const initChart = () => {
  echarts.use([
    TitleComponent,
    TooltipComponent,
    GridComponent,
    LegendComponent,
    LineChart,
    CanvasRenderer,
    UniversalTransition,
    ToolboxComponent,
  ]);

  const option = {
    backgroundColor: 'transparent',
    title: { text: '' },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', boundaryGap: false, splitLine: { show: true }, data: [] },
    yAxis: { type: 'value' },
    legend: { data: [] },
    series: [],
  };

  lag_chart.value = echarts.init(lag_chartRef.value, echarts_theme);
  lag_chart.value.setOption({ ...option, title: { text: t('inspection.lag') } });

  commit_chart.value = echarts.init(commit_chartRef.value, echarts_theme);
  commit_chart.value.setOption({ ...option, title: { text: t('inspection.commit') } });

  end_chart.value = echarts.init(end_chartRef.value, echarts_theme);
  end_chart.value.setOption({ ...option, title: { text: t('inspection.end') } });

  productionSpeed_chart.value = echarts.init(productionSpeed_chartRef.value, echarts_theme);
  productionSpeed_chart.value.setOption({ ...option, title: { text: 'ProduceSpeed(msg/s)' } });

  consumptionSpeed_chart.value = echarts.init(consumptionSpeed_chartRef.value, echarts_theme);
  consumptionSpeed_chart.value.setOption({ ...option, title: { text: 'ConsumeSpeed(msg/s)' } });
};

const clear_offset = () => {
  offsetData.value = { lag: {}, commit: {}, end: {}, productionSpeed: {}, consumptionSpeed: {} };
};

const updateChart = () => {
  const chart_map = {
    lag: lag_chart.value,
    commit: commit_chart.value,
    end: end_chart.value,
    productionSpeed: productionSpeed_chart.value,
    consumptionSpeed: consumptionSpeed_chart.value,
  };

  for (const k in offsetData.value) {
    let series = [];
    let legendData = [];
    let xs = [];

    Object.entries(offsetData.value[k]).forEach(([topic, data]) => {
      legendData.push(topic);
      series.push({
        name: topic,
        type: 'line',
        symbol: 'circle',
        data: data.map(item => (k === 'productionSpeed' || k === 'consumptionSpeed' ? item.speed : item.offset)),
      });
      xs = data.map(item => {
        const date = new Date(item.time);
        return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`;
      });
    });

    chart_map[k].setOption({
      xAxis: { data: xs },
      legend: { data: legendData },
      series,
    });
  }
};

const fetchData = async () => {
  if (selectedTopics.value.length === 0 || !selectedGroup.value) {
    message.warning(t('message.selectTopicGroup'));
    return;
  }
  loading.value = true;

  try {
    const res = await GetTopicOffsets(selectedTopics.value, selectedGroup.value);
    if (res.err !== '') {
      message.error(res.err);
    } else {
      const now = new Date();
      const time = now.getTime(); // Use timestamp for precise time difference

      selectedTopics.value.forEach(topic => {
        const endOffset = addOffsets(res.result.end_map[topic]) || 0;
        const commitOffset = addOffsets(res.result.commit_map[topic]) || 0;

        // Lag
        if (!offsetData.value.lag[topic]) offsetData.value.lag[topic] = [];
        offsetData.value.lag[topic].push({ time, offset: endOffset - commitOffset });

        // Commit
        if (!offsetData.value.commit[topic]) offsetData.value.commit[topic] = [];
        offsetData.value.commit[topic].push({ time, offset: commitOffset });

        // End
        if (!offsetData.value.end[topic]) offsetData.value.end[topic] = [];
        offsetData.value.end[topic].push({ time, offset: endOffset });

        // Production Speed
        if (!offsetData.value.productionSpeed[topic]) offsetData.value.productionSpeed[topic] = [];
        if (offsetData.value.end[topic].length >= 2) {
          const lastEnd = offsetData.value.end[topic][offsetData.value.end[topic].length - 2];
          const deltaEnd = endOffset - lastEnd.offset;
          const deltaTime = (time - lastEnd.time) / 1000; // Convert to seconds
          const speed = deltaTime > 0 ? deltaEnd / deltaTime : 0;
          offsetData.value.productionSpeed[topic].push({ time, speed });
        }

        // Consumption Speed
        if (!offsetData.value.consumptionSpeed[topic]) offsetData.value.consumptionSpeed[topic] = [];
        if (offsetData.value.commit[topic].length >= 2) {
          const lastCommit = offsetData.value.commit[topic][offsetData.value.commit[topic].length - 2];
          const deltaCommit = commitOffset - lastCommit.offset;
          const deltaTime = (time - lastCommit.time) / 1000; // Convert to seconds
          const speed = deltaTime > 0 ? deltaCommit / deltaTime : 0;
          offsetData.value.consumptionSpeed[topic].push({ time, speed });
        }

        // Keep only the last 100 data points
        for (const key of ['lag', 'commit', 'end', 'productionSpeed', 'consumptionSpeed']) {
          if (offsetData.value[key][topic]?.length > 100) offsetData.value[key][topic].shift();
        }
      });

      updateChart();
    }
  } catch (e) {
    message.error(e);
  } finally {
    loading.value = false;
  }
};

const getData = async () => {
  console.log('Initializing consumer data');
  try {
    const [res, res2] = await Promise.all([GetTopics(), GetGroups()]);
    if (res.err !== '' || res2.err !== '') {
      message.error(res.err || res2.err);
    } else {
      topic_data.value = res.results
          ?.sort((a, b) => (a.topic > b.topic ? 1 : -1))
          .map(r => ({ label: r.topic, value: r.topic })) || [];

      group_data.value = res2.results
          ?.map(g => ({
            label: g.Group,
            value: g.Group,
            State: g.State,
            ProtocolType: g.ProtocolType,
            Coordinator: g.Coordinator,
          }))
          .sort((a, b) => (a.label > b.label ? 1 : -1)) || [];
    }
  } catch (e) {
    message.error(e);
  }
};

const addOffsets = (item) => {
  return Object.values(item || {}).reduce((sum, { At }) => sum + (At > 0 ? At : 0), 0);
};

const selectNode = async () => {
  topic_data.value = [];
  group_data.value = [];
  selectedTopics.value = [];
  selectedGroup.value = null;
  offsetData.value = { lag: {}, commit: {}, end: {}, productionSpeed: {}, consumptionSpeed: {} };
  loading.value = false;
  await getData();
};
</script>

<style scoped>
</style>