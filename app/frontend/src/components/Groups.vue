<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2 style="max-width: 200px;">消费者组</h2>
      <n-button @click="getData" text :render-icon="renderIcon(RefreshOutlined)">refresh</n-button>
      <n-text>共计{{ group_data.length }}个</n-text>
      <n-button @click="downloadAllDataCsv" :render-icon="renderIcon(DriveFileMoveTwotone)">导出为csv</n-button>
    </n-flex>
    <n-spin :show="loading" description="Connecting...">
      <n-data-table
          ref="tableRef"
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
    <n-drawer-content title="创建Topic配置">
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
import {DeleteGroup, DeleteTopic, GetGroupMembers, GetGroups} from "../../wailsjs/go/service/Service";

const group_data = ref([])
const members_data = ref([])
const loading = ref(false)
const showDrawer = ref(false)
const message = useMessage()
const tableRef = ref();

const selectNode = async (node) => {
  const data_lst = [members_data, group_data]
  for (const k in data_lst) {
    data_lst[k].value = []
  }
}

onMounted(async () => {
  emitter.on('selectNode', selectNode)
  await getData()
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
      }
    }
  } catch (e) {
    message.error(e)
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
        message.warning("没有找到成员")
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
  const csvContent = createCsvContent(
      group_data.value, columns
  )
  download_file(csvContent, '导出.csv', 'text/csv;charset=utf-8;')
}


const columns = [
  {title: 'Group', key: 'Group', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: true}},
  {
    title: 'Coordinator', key: 'Coordinator', sorter: 'default', width: 50, resizable: true, ellipsis: {tooltip: true},
    render: (row) => h(NTag, {type: "info"}, {default: () => row['Coordinator']}),
  },
  {
    title: 'State', key: 'State', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: true},
    render: (row) => h(NTag, {type: "success"}, {default: () => row['State']}),
  },
  {
    title: 'ProtocolType',
    key: 'ProtocolType',
    sorter: 'default',
    width: 20,
    resizable: true,
    ellipsis: {tooltip: true}
  },
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
                NButton,
                {
                  strong: true,
                  secondary: true,
                  onClick: async () => {
                    await getMembers(row['Group'])
                    showDrawer.value = true
                  }
                },
                {default: () => '查看成员', icon: () => h(NIcon, null, {default: () => h(SettingsTwotone)})}
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
                                default: () => '删除',
                                icon: () => h(NIcon, null, {default: () => h(DeleteForeverTwotone)})
                              }
                          )
                  ,
                  default: () => `确认删除${row["Group"]}?`
                }
            ),
          ]
        }
    )
  },
]

const members_columns = [
  {title: 'MemberID', key: 'MemberID', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: true}},
  {title: 'InstanceID', key: 'InstanceID', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: true}},
  {title: 'ClientID', key: 'ClientID', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: true}},
  {title: 'ClientHost', key: 'ClientHost', sorter: 'default', width: 20, resizable: true, ellipsis: {tooltip: true}},

]

const deleteGroups = async (group) => {
  loading.value = true
  try {
    const res = await DeleteGroup([group])
    if (res.err !== "") {
      message.error(res.err)
    } else {
      message.success("删除成功")
      await getData()
    }
  } catch (e) {
    message.error(e)
  }
  loading.value = false

}
</script>


<style scoped>

</style>