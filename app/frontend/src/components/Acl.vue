<template>
  <n-flex vertical>
    <n-flex align="center">
      <h2>{{ t('acl.title') }}</h2>
      <n-button @click="getAcls" text :render-icon="renderIcon(RefreshOutlined)">{{ t('common.refresh') }}</n-button>
      <n-text>{{ t('common.count') }}：{{ acls ? acls.length : 0 }}</n-text>
      <n-button @click="showAddModal = true" :render-icon="renderIcon(AddFilled)">{{ t('acl.add') }}</n-button>
    </n-flex>
    <n-spin :show="loading" :description="t('common.connecting')">
      <n-data-table
          :columns="columns"
          :data="acls"
          size="small"
          :bordered="false"
          striped
          :pagination="pagination"
      />
    </n-spin>
  </n-flex>

  <!-- 添加/编辑 ACL 模态框 -->
  <n-modal v-model:show="showAddModal" preset="dialog" :title="t('acl.add')">
    <n-form
        ref="formRef"
        :model="aclForm"
        label-placement="top"
        style="text-align: left;"
        require-mark-placement="right-hanging"
    >
      <n-form-item :label="t('acl.principal')" path="principal">
        <n-input v-model:value="aclForm.principal" :placeholder="t('acl.principalPlaceholder')"/>
      </n-form-item>
      <n-form-item :label="t('acl.resourceType')" path="resourceType">
        <n-select
            v-model:value="aclForm.resourceType"
            :options="resourceTypeOptions"
            :placeholder="t('acl.resourceTypePlaceholder')"
        />
      </n-form-item>
      <n-form-item :label="t('acl.resourceName')" path="resourceName">
        <n-input v-model:value="aclForm.resourceName" :placeholder="t('acl.resourceNamePlaceholder')"/>
      </n-form-item>
      <n-form-item :label="t('acl.operation')" path="operation">
        <n-select
            v-model:value="aclForm.operation"
            :options="operationOptions"
            :placeholder="t('acl.operationPlaceholder')"
        />
      </n-form-item>
      <n-form-item :label="t('acl.permissionType')" path="permissionType">
        <n-select
            v-model:value="aclForm.permissionType"
            :options="permissionTypeOptions"
            :placeholder="t('acl.permissionTypePlaceholder')"
        />
      </n-form-item>
    </n-form>
      <n-flex>
        <n-button @click="showAddModal = false">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" @click="addAcl">{{ t('common.enter') }}</n-button>
      </n-flex>
  </n-modal>
</template>

<script setup>
import {h, onMounted, ref} from "vue";
import {NButton, NDataTable, NForm, NFormItem, NInput, NSelect, NSpace, useMessage} from 'naive-ui';
import {AddFilled, RefreshOutlined} from '@vicons/material';
import {CreateAcl, DeleteAcl, GetAcls} from "../../wailsjs/go/service/Service";
import {useI18n} from "vue-i18n";
import {renderIcon} from "../utils/common";

const {t} = useI18n();

const message = useMessage();

const acls = ref([]);
const loading = ref(false);
const showAddModal = ref(false);

const aclForm = ref({
  principal: '',
  resourceType: '',
  resourceName: '',
  operation: '',
  permissionType: ''
});

const resourceTypeOptions = [
  {label: 'Topic', value: 'TOPIC'},
  {label: 'Group', value: 'GROUP'},
  {label: 'Cluster', value: 'CLUSTER'},
  {label: 'Transactional ID', value: 'TRANSACTIONAL_ID'},
  {label: 'Delegation Token', value: 'DELEGATION_TOKEN'}
];

const operationOptions = [
  {label: 'Read', value: 'READ'},
  {label: 'Write', value: 'WRITE'},
  {label: 'Create', value: 'CREATE'},
  {label: 'Delete', value: 'DELETE'},
  {label: 'Alter', value: 'ALTER'},
  {label: 'Describe', value: 'DESCRIBE'},
  {label: 'Cluster Action', value: 'CLUSTER_ACTION'},
  {label: 'Describe Configs', value: 'DESCRIBE_CONFIGS'},
  {label: 'Alter Configs', value: 'ALTER_CONFIGS'},
  {label: 'Idempotent Write', value: 'IDEMPOTENT_WRITE'}
];

const permissionTypeOptions = [
  {label: 'Allow', value: 'ALLOW'},
  {label: 'Deny', value: 'DENY'}
];

const columns = [
  {title: 'Principal', key: 'principal', sorter: 'default', width: 100},
  {title: 'Host', key: 'host', sorter: 'default', width: 100},
  {title: 'Resource Type', key: 'resourceType', sorter: 'default', width: 100},
  {title: 'Resource Name', key: 'resourceName', sorter: 'default', width: 100},
  {title: 'Operation', key: 'operation', sorter: 'default', width: 100},
  {title: 'Pattern Type', key: 'patternType', sorter: 'default', width: 100},
  {title: 'Permission Type', key: 'permissionType', sorter: 'default', width: 100},
  {
    title: t('common.action'),
    key: 'actions',
    width: 100,
    render: (row) => {
      return h(NSpace, null, {
        default: () => [
          h(NButton, {
            type: 'error',
            onClick: () => deleteAcl(row)
          }, {default: () => t('common.delete')})
        ]
      });
    }
  }
];

const pagination = ref({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [5, 10, 20, 30, 40],
  onChange: (page) => {
    pagination.value.page = page;
  },
  onUpdatePageSize: (pageSize) => {
    pagination.value.pageSize = pageSize;
    pagination.value.page = 1;
  },
});


const getAcls = async () => {
  loading.value = true;
  try {
    const res = await GetAcls();
    if (res.err !== "") {
      message.error(res.err, {duration:  5000});
    } else {
      console.log(res)
      acls.value = res.results;
    }
  } catch (e) {
    message.error(e.message, {duration:  5000});
  }
  loading.value = false;
};

const addAcl = async () => {
  try {
    const res = await CreateAcl(aclForm.value);
    if (res.err !== "") {
      message.error(res.err, {duration:  5000});
    } else {
      message.success(t('message.addOk'));
      showAddModal.value = false;
      await getAcls();
    }
  } catch (e) {
    message.error(e.message, {duration:  5000});
  }
};

const deleteAcl = async (acl) => {
  try {
    const res = await DeleteAcl(acl.id);
    if (res.err !== "") {
      message.error(res.err, {duration:  5000});
    } else {
      message.success(t('message.deleteFinish'));
      await getAcls();
    }
  } catch (e) {
    message.error(e.message, {duration:  5000});
  }
};

onMounted(() => {
  getAcls();
});
</script>

<style scoped>
</style>