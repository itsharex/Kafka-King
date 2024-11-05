<template>
  <div>
    <n-flex vertical>
      <n-flex align="center">
        <h2 style="max-width: 200px;">集群</h2>
        <n-text>共有 {{ Nodes.length }} 个</n-text>
        <n-button @click="addNewNode" :render-icon="renderIcon(AddFilled)">添加集群</n-button>
      </n-flex>
      <n-spin :show="spin_loading" description="Connecting...">

        <n-grid :x-gap="12" :y-gap="12" :cols="4">
          <n-gi v-for="node in Nodes" :key="node.id">
            <n-card :title="node.name" @click="selectNode(node)" hoverable class="conn_card">

              <template #header-extra>
                <n-space>
                  <n-button @click.stop="editNode(node)" size="small">
                    编辑
                  </n-button>
                  <n-popconfirm @positive-click="deleteNode(node.id)" negative-text="取消" positive-text="确定">
                    <template #trigger>
                      <n-button @click.stop size="small">
                        删除
                      </n-button>
                    </template>
                    确定删除吗？
                  </n-popconfirm>
                </n-space>
              </template>
              <n-descriptions :column="1" label-placement="left">
                <n-descriptions-item label="地址">
                  {{ node.bootstrap_servers }}
                </n-descriptions-item>
              </n-descriptions>
            </n-card>
          </n-gi>
        </n-grid>
      </n-spin>
    </n-flex>

    <n-drawer v-model:show="showEditDrawer" :width="500" placement="right">
      <n-drawer-content :title="drawerTitle">
        <n-form
            ref="formRef"
            :model="currentNode"
            :rules="{
              name: {required: true, message: '请输入昵称', trigger: 'blur'},
              bootstrap_servers: {required: true, message: '请输入连接地址', trigger: 'blur'},
            }"
            label-placement="top"
            style="text-align: left;"
        >
          <n-form-item label="昵称" path="name">
            <n-input v-model:value="currentNode.name" placeholder="输入名称"/>
          </n-form-item>

          <n-form-item label="连接地址" path="bootstrap_servers">
            <n-input v-model:value="currentNode.bootstrap_servers" placeholder="127.0.0.1:9092,127.0.0.1:9093"/>
          </n-form-item>
          注意：必须保证本地能够访问 kafka 配置的 advertised.listeners 地址 （特别是域名解析，即使你填的是ip，也需要在本地配置好hosts）
          <n-form-item label="使用 TLS" path="tls">
            <n-switch checked-value="enable" unchecked-value="disable" v-model:value="currentNode.tls"/>
          </n-form-item>

          <n-form-item label="跳过 TLS 验证" path="skipTLSVerify">
            <n-switch checked-value="enable" unchecked-value="disable" value="enable"
                      v-model:value="currentNode.skipTLSVerify"/>
          </n-form-item>

          <n-form-item label="TLS certFile" path="tls_cert_file">
            <n-input v-model:value="currentNode.tls_cert_file" placeholder="输入 pem 证书路径"/>
          </n-form-item>

          <n-form-item label="TLS keyFile" path="tls_key_file">
            <n-input v-model:value="currentNode.tls_key_file" placeholder="输入 key 私钥路径"/>
          </n-form-item>

          <n-form-item label="TLS CA 证书" path="tls_ca_file">
            <n-input v-model:value="currentNode.tls_ca_file" placeholder="输入 CA 证书路径"/>
          </n-form-item>

          <n-form-item label="使用 SASL" path="sasl">
            <n-switch checked-value="enable" unchecked-value="disable" v-model:value="currentNode.sasl"/>
          </n-form-item>

          <n-form-item label="SASL 机制" path="sasl_mechanism">
            <n-select
                v-model:value="currentNode.sasl_mechanism"
                :options="sasl_mechanism_options"
                placeholder="请选择"
                filterable
                clearable
                style="width: 200px"
            />
          </n-form-item>

          <n-form-item label="SASL 用户名" path="sasl_user">
            <n-input v-model:value="currentNode.sasl_user" placeholder="输入用户名"/>
          </n-form-item>

          <n-form-item label="SASL 密码" path="sasl_pwd">
            <n-input
                v-model:value="currentNode.sasl_pwd"
                type="password"
                placeholder="输入密码"
            />
          </n-form-item>

          <n-form-item label="kerberos keytab 路径" path="kerberos_user_keytab">
            <n-input v-model:value="currentNode.kerberos_user_keytab" placeholder="输入keytab文件路径"/>
          </n-form-item>

          <n-form-item label="kerberos krb5.conf 路径" path="kerberos_krb5_conf">
            <n-input v-model:value="currentNode.kerberos_krb5_conf" placeholder="输入krb5.conf文件路径"/>
          </n-form-item>

          <n-form-item label="Kerberos_user 用户名" path="sasl_user">
            <n-input v-model:value="currentNode.Kerberos_user" placeholder="输入Kerberos_user用户名"/>
          </n-form-item>

          <n-form-item label="Kerberos_realm 领域域名" path="Kerberos_realm">
            <n-input v-model:value="currentNode.Kerberos_realm" placeholder="输入Kerberos领域域名"/>
          </n-form-item>

          <n-form-item label="kerberos_service_name 服务名" path="sasl_user">
            <n-input v-model:value="currentNode.kerberos_service_name" placeholder="输入配置的kerberos_service_name"/>
          </n-form-item>

        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button @click="test_connect" :loading="test_connect_loading">连接测试</n-button>
            <n-button @click="showEditDrawer = false">取消</n-button>
            <n-button type="primary" @click="saveNode">保存</n-button>
          </n-space>
        </template>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup>
import {computed, onMounted, ref} from 'vue'
import {useMessage} from 'naive-ui'
import {renderIcon} from "../utils/common";
import {AddFilled} from "@vicons/material";
import emitter from "../utils/eventBus";
import {SetConnect, TestClient} from "../../wailsjs/go/service/Service";
import {GetConfig, SaveConfig} from "../../wailsjs/go/config/AppConfig";


const message = useMessage()

const Nodes = ref([])

const showEditDrawer = ref(false)
const currentNode = ref({
  id: 0,
  name: '',
  bootstrap_servers: '',
  tls: 'disable',
  skipTLSVerify: 'true',
  tls_cert_file: '',
  tls_key_file: '',
  tls_ca_file: '',
  sasl: 'disable',
  sasl_mechanism: "PLAIN",
  sasl_user: '',
  sasl_pwd: '',
  kerberos_user_keytab: '',
  kerberos_krb5_conf: '',
  Kerberos_user: '',
  Kerberos_realm: '',
  kerberos_service_name: '',
})
const isEditing = ref(false)
const spin_loading = ref(false)
const test_connect_loading = ref(false)
const sasl_mechanism_options = [
  {
    label: 'PLAIN',
    value: 'PLAIN'
  },
  {
    label: 'SCRAM-SHA-256',
    value: 'SCRAM-SHA-256'
  },
  {
    label: 'SCRAM-SHA-512',
    value: 'SCRAM-SHA-512'
  },
  {
    label: 'GSSAPI',
    value: 'GSSAPI'
  },
]

const drawerTitle = computed(() => isEditing.value ? '编辑连接' : '添加连接')

const formRef = ref(null)

onMounted(async () => {
  await refreshNodeList()
})

const refreshNodeList = async () => {
  spin_loading.value = true
  const config = await GetConfig()
  Nodes.value = config.connects
  spin_loading.value = false
}

function editNode(node) {
  currentNode.value = {...node}
  isEditing.value = true
  showEditDrawer.value = true
}

const addNewNode = async () => {
  currentNode.value = {}
  isEditing.value = false
  showEditDrawer.value = true
}

const saveNode = async () => {
  formRef.value?.validate(async (errors) => {
    if (!errors) {

      const config = await GetConfig()
      // edit
      if (isEditing.value) {
        const index = Nodes.value.findIndex(node => node.id === currentNode.value.id)
        if (index !== -1) {
          Nodes.value[index] = {...currentNode.value}
        }
      } else {
        // add
        const newId = Math.max(...Nodes.value.map(node => node.id), 0) + 1
        Nodes.value.push({...currentNode.value, id: newId})
      }
      console.log(config)

      // 保存
      config.connects = Nodes.value
      await SaveConfig(config)
      showEditDrawer.value = false

      await refreshNodeList()
      message.success('保存成功')
    } else {
      message.error('请填写所有必填字段')
    }
  })
}

const deleteNode = async (id) => {
  Nodes.value = Nodes.value.filter(node => node.id !== id)
  const config = await GetConfig()
  config.connects = Nodes.value
  await SaveConfig(config)
  await refreshNodeList()
  message.success('删除成功')
}

const test_connect = async () => {
  formRef.value?.validate(async (errors) => {
    if (!errors) {

      test_connect_loading.value = true
      try {
        const res = await TestClient(currentNode.value.name, currentNode.value)
        if (res.err !== "") {
          message.error("连接失败：" + res.err)
        } else {
          message.success('连接成功')
        }
      } catch (e) {
        message.error(e)
      }
      test_connect_loading.value = false
    }else {
      message.error('请填写所有必填字段')
    }
  })
}

const selectNode = async (node) => {
  // 这里实现切换菜单的逻辑
  console.log('选中节点:', node)
  spin_loading.value = true
  try {
    const res = await SetConnect(node.name, node, false)
    if (res.err !== "") {
      message.error("连接失败：" + res.err)
    } else {
      emitter.emit('menu_select', "节点")
      emitter.emit('selectNode', node)
      message.success('连接成功')
    }
  } catch (e) {
    message.error(e)
  }
  spin_loading.value = false
}

const handleSelect = (key) => {
  currentNode.value.sasl_mechanism = key
}
</script>

<style>

.lightTheme .conn_card {
  background-color: #fafafc
}
</style>