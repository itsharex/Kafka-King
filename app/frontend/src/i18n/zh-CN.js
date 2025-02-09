/*
 * Copyright 2025 Bronya0 <tangssst@163.com>.
 * Author Github: https://github.com/Bronya0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export default {
    common:{
        enter: '确认',
        cancel: '取消',
        name: '昵称',
        save: '保存',
        check: '请选择',
        add: '创建',
        edit: '编辑',
        delete: '删除',
        deleteFinish: '删除成功',
        deleteOk: '确定删除吗？',
        count: '数量',
        refresh: '刷新',
        read: '读取',
        config: '配置',
        csv: '导出为csv',
        connecting: '连接中...',
        action: '操作',
    },
    aside: {
        cluster: '集群',
        node: '节点',
        topic: '主题',
        producer: '生产者',
        consumer: '消费者',
        group: '消费者组',
        monitor: '巡检',
        settings: '设置'
    },
    conn:{
        title: '集群',
        add: '添加集群',
        edit: '编辑连接',
        test: '连接测试',
        add_link: '添加连接',
        please_add_name: '请输入昵称',
        please_add_link: '请输入连接地址',
        input_name: '请输入名称',
        bootstrap_servers: '连接地址',
        tip: '注意：必须保证本地能够访问 kafka 配置的 advertised.listeners 地址 （特别是域名解析，即使你填的是ip，也需要在本地配置好hosts）',
        tls: '使用 TLS',
        skipTLSVerify: '跳过 TLS 验证',
        tls_cert_file: '输入 pem 证书路径',
        tls_key_file: '输入 key 私钥路径',
        tls_ca_file: '输入 CA 证书路径',
        use_sasl: '使用 SASL',
        sasl_mechanism: 'SASL 机制',
        sasl_user: 'SASL 用户名',
        sasl_pwd: 'SASL 密码',
        kerberos_user_keytab: 'kerberos keytab 文件路径',
        kerberos_krb5_conf: 'kerberos krb5.conf 路径',
        Kerberos_user: 'Kerberos_user 用户名',
        Kerberos_realm: 'Kerberos_realm 领域域名',
        kerberos_service_name: 'kerberos_service_name 服务名',
    },
    node: {
        title: '节点',
        source: '来源',
        value: 'Value（双击可编辑）',
        sensitive: '是否敏感',
        ok_message: '编辑成功，刷新配置',
    },
    topic: {
        title: '主题',
        add: '创建主题',
        add_name: '输入主题名称',
        selectedGroup: '选择Group并读取Offsets',
        partition: '分区',
        add_partition: '添加分区',
        add_partition_count: '添加的额外的分区数',
        replication_factor: '副本数',
        err: '主题故障',
        lag: '积压',
        viewProduce: '生产消息',
        viewConsumer: '查看/消费消息',
        viewOffset: '读取 offset',
        viewConfig: '主题配置',
        viewPartition: '查看分区',
        deleteTopic: '删除Topic',
    },
    producer: {
        title: "Producer",
        desc: "一个生产者客户端，将消息推送到指定的Topic。",
        selectTopic: "选择Topic",
        topicPlaceholder: "必选：选择或搜索Kafka Topic",
        optionalMessageKey: "可选：输入消息Key",
        keyPlaceholder: "可选：输入消息Key",
        specifyPartition: "可选：指定分区号",
        messageContentPlaceholder: "必填：消息内容，字符串格式，支持JSON",
        headersTitle: "消息Headers:",
        addHeader: "添加Header",
        removeHeader: "删除",
        headerKeyPlaceholder: "Header Key",
        headerValuePlaceholder: "Header Value",
        sendTimes: "发送次数",
        sendTimesPlaceholder: "发送次数",
        sendMessage: "发送消息"
    },
    consumer: {
        title: "Consumer",
        desc: "一个简单消费者客户端，查看Topic消息。",
        requiredTopic: "必选：Topic",
        topicPlaceholder: "选择或搜索Kafka Topic",
        requiredMessagesCount: "必选：消费数量",
        messagesCountPlaceholder: "消费消息数量",
        pollTimeoutDescription: "poll超时时间：默认10s。如异常或无可消费消息，则会超时",
        pollTimeoutPlaceholder: "poll超时时间",
        optionalGroup: "可选：Group（一旦选择，消费时会自动提交Offset。支持创建新Group）",
        groupPlaceholder: "选择或创建Consumer Group",
        consumeMessage: "消费消息"
    },
    group:{
        title: "消费者组",
        member: '成员',
        warn: '请先从具体的topic切换到本页',
    },
    inspection: {
        title: "巡检",
        desc: "巡检Kafka积压情况。",
        topicsLabel: "Topics",
        topicPlaceholder: "选择或搜索Kafka Topic",
        groupLabel: "Group",
        groupPlaceholder: "选择或创建Consumer Group",
        startInspection: "开始巡检",
        autoFetch: "每5分钟自动抓取一次数据",
        lagFormula: "积压 = 终末offset - 提交offset。",
        lag: '积压量',
        commit: '提交 offset',
        end: '终末 offset',
    },
    settings:{
        title: '设置',
        width: '窗口宽度',
        height: '窗口高度',
        lang: '语言（需要重启）',
        theme: '主题',
    },
    message:{
        noMemberFound: "没有找到成员",
        saveSuccess: "保存成功",
        connectSuccess: "连接成功",
        fetchSuccess: "获取成功",
        sendSuccess: '发送成功',
        selectGroupFirst: "请先选择 Group",
        selectTopic: "请选择Topic",
        selectTopicGroup: "请选择Topic和Group",
        connectErr: "连接失败",
        addOk: "添加成功",
        editOk: "编辑成功，刷新配置",
        mustFill: "请填写所有必填字段",
        saveErr: "保存失败",
        pleaseInput: "请输入消息内容",
    },
    about: {
        title: "关于",
        projectHomepage: "项目主页",
        kafkaKing: "Kafka-King",
        esClient: "同款 ES 客户端",
        esKing: "ES-King",
        technicalGroup: "技术交流群",
        qqGroup: "QQ交流群",
        translate: '翻译有问题？报告或参与翻译',
    },
    header:{
        desc: '更人性化的 Kafka GUI',
        c_node: '当前集群',
        netErr: '无法连接github，请检查网络',
        newVersion: '发现新版本',
        down: '立即下载',
    }
};