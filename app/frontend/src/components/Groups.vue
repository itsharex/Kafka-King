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
      <h2 >{{t('group.title')}}</h2>
      <n-button @click="getData" text :render-icon="renderIcon(RefreshOutlined)">{{t('common.refresh')}}</n-button>
      <n-text>{{t('common.count')}}：{{ group_data.length }}</n-text>
      <n-button @click="downloadAllDataCsv" :render-icon="renderIcon(DriveFileMoveTwotone)">{{t('common.csv')}}</n-button>
    </n-flex>
    <n-spin :show="loading" :description="t('common.connecting')">
      <n-data-table
          :columns="columns"
          :data="group_data"
          size="small"
          :bordered="false"
          striped
          :pagination="pagination"
      />
    </n-spin>
  </n-flex>

  <n-drawer v-model:show="showDrawer" :width="800">
    <n-drawer-content :title="t('group.member')">
      <n-data-table
          :columns="members_columns"
          :data="members_data"
          :bordered="false"
          :pagination="pagination"
      />
    </n-drawer-content>
  </n-drawer>

</template>
<script setup>
import {h, onMounted, ref} from "vue";
import emitter from "../utils/eventBus";
import {NButton, NButtonGroup, NDataTable, NIcon, NPopconfirm, NTag, NText, useMessage} from 'naive-ui'
import {createCsvContent, download_file, renderIcon} from "../utils/common";
import {DeleteForeverTwotone, DriveFileMoveTwotone, RefreshOutlined, SettingsTwotone} from "@vicons/material";
import {DeleteGroup, GetGroupMembers, GetGroups} from "../../wailsjs/go/service/Service";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const group_data = ref([])
const members_data = ref([])
const loading = ref(false)
const showDrawer = ref(false)
const message = useMessage()

const selectNode = async (node) => {
  group_data.value = []
  members_data.value = []

  loading.value = false
  await getData()
}

onMounted(() => {
  emitter.on('selectNode', selectNode)

  getData()
})


const getData = async () => {
  loading.value = true
  try {
    const res = await GetGroups()
    if (res.err !== "") {
      message.error(res.err)
    } else {
      if (res.results) {
        res.results.sort((a, b) => a['Group'] > b['Group'] ? 1 : -1)
        group_data.value = res.results
        console.log(res.results)
      }
    }
  } catch (e) {
    message.error(e.message)
  }
  loading.value = false
}

const getMembers = async (group) => {
  loading.value = true
  try {
    const res = await GetGroupMembers([group])
    if (res.err !== "") {
      message.error(res.err)
    } else {
      if (res.results[0]) {
        let data0 = res.results[0]
        members_data.value = data0['Members']
      }else {
        message.warning(t('message.noMemberFound'))
      }
    }

  } catch (e) {
    message.error(e.message)
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
  const csvContent = createCsvContent(
      group_data.value, columns
  )
  download_file(csvContent, 'kafka-group.csv', 'text/csv;charset=utf-8;')
}


const columns = [
  {title: 'Group', key: 'Group', sorter: 'default', width: 60, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}}},
  {
    title: 'Coordinator', key: 'Coordinator', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
    render: (row) => h(NTag, {type: "info"}, {default: () => row['Coordinator']}),
  },
  {
    title: 'State', key: 'State', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}},
    render: (row) => h(NTag, {type: "success"}, {default: () => row['State'] || 'unknown'}),
  },
  {
    title: 'ProtocolType',
    key: 'ProtocolType',
    sorter: 'default',
    width: 20,
    resizable: true,
    ellipsis: {tooltip: {style: { maxWidth: '800px' },}}
  },
  {
    title: t('common.action'),
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
                NButton,
                {
                  strong: true,
                  secondary: true,
                  onClick: async () => {
                    await getMembers(row['Group'])
                    showDrawer.value = true
                  }
                },
                {default: () => t('group.member'), icon: () => h(NIcon, null, {default: () => h(SettingsTwotone)})}
            ),
            h(
                NPopconfirm,
                {
                  onPositiveClick: () => deleteGroups(row["Group"])
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
                              {
                                default: () => t('common.delete'),
                                icon: () => h(NIcon, null, {default: () => h(DeleteForeverTwotone)})
                              }
                          )
                  ,
                  default: () => `${t('common.deleteOk')} Group: ${row["Group"]}?`
                }
            ),
          ]
        }
    )
  },
]

const members_columns = [
  {title: 'MemberID', key: 'MemberID', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}}},
  {title: 'InstanceID', key: 'InstanceID', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}}},
  {title: 'ClientID', key: 'ClientID', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}}},
  {title: 'ClientHost', key: 'ClientHost', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: {style: { maxWidth: '800px' },}}},

]

const deleteGroups = async (group) => {
  loading.value = true
  try {
    const res = await DeleteGroup(group)
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success(t('common.deleteFinish'))
      await getData()
    }
  } catch (e) {
    message.error(e.message)
  }
  loading.value = false

}
</script>


<style scoped>

</style>