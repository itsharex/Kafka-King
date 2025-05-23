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
    common: {
        enter: "Confirm",
        cancel: "Cancel",
        name: "Nickname",
        save: "Save",
        check: "Please select",
        add: "Create",
        edit: "Edit",
        delete: "Delete",
        deleteFinish: "Delete successful",
        deleteOk: "Are you sure you want to delete?",
        count: "Quantity",
        refresh: "Refresh",
        read: "Read",
        config: "Configuration",
        csv: "Export as CSV",
        connecting: "Connecting...",
        action: "Action",
        compress: 'Compress',
        decompress: 'Decompress',
    },
    aside: {
        cluster: "Cluster",
        node: "Node",
        topic: "Topic",
        producer: "Producer",
        consumer: "Consumer",
        group: "Consumer Group",
        monitor: "Inspection",
        settings: "Settings"
    },
    conn: {
        title: "Cluster",
        add: "Add Cluster",
        edit: "Edit Connection",
        test: "Test Connection",
        add_link: "Add Connection",
        please_add_name: "Please enter a nickname",
        please_add_link: "Please enter the connection address",
        input_name: "Please enter the name",
        bootstrap_servers: "Connection Address",
        tip: "Note: Ensure that your local environment can access the Kafka advertised.listeners address (especially domain name resolution, even if you fill in an IP, you need to configure hosts properly on your local machine)",
        tls: "Use TLS",
        skipTLSVerify: "Skip TLS Verification",
        tls_cert_file: "Enter PEM certificate path",
        tls_key_file: "Enter key private path",
        tls_ca_file: "Enter CA certificate path",
        use_sasl: "Use SASL",
        sasl_mechanism: "SASL Mechanism",
        sasl_user: "SASL Username",
        sasl_pwd: "SASL Password",
        kerberos_user_keytab: "Kerberos keytab file path",
        kerberos_krb5_conf: "Kerberos krb5.conf path",
        Kerberos_user: "Kerberos User",
        Kerberos_realm: "Kerberos Realm",
        kerberos_service_name: "Kerberos Service Name"
    },
    node: {
        title: "Node",
        source: "Source",
        value: "Value (double-click to edit)",
        sensitive: "Is it sensitive?",
        ok_message: "Edit successful, refreshing configuration"
    },
    topic: {
        title: "Topic",
        add: "Create Topic",
        add_name: "Enter topic name",
        selectedGroup: "Select Group and read Offsets",
        partition: "Partition",
        add_partition: "Add Partition",
        add_partition_count: "Number of additional partitions to add",
        replication_factor: "Replication Factor",
        err: "Topic Error",
        lag: "Lag",
        viewProduce: "Produce Messages",
        viewConsumer: "View/Consume Messages",
        viewOffset: "Read Offset",
        viewConfig: "Topic Configuration",
        viewPartition: "View Partitions",
        deleteTopic: "Delete Topic"
    },
    producer: {
        title: "Producer",
        desc: "A producer client that pushes messages to the specified Topic.",
        selectTopic: "Select Topic",
        topicPlaceholder: "Required: Select or search for a Kafka Topic",
        optionalMessageKey: "Optional: Enter message Key",
        keyPlaceholder: "Optional: Enter message Key",
        specifyPartition: "Optional: Specify partition number",
        messageContentPlaceholder: "Required: Message content, string format, JSON supported",
        headersTitle: "Message Headers:",
        addHeader: "Add Header",
        removeHeader: "Remove",
        headerKeyPlaceholder: "Header Key",
        headerValuePlaceholder: "Header Value",
        sendTimes: "Send Times",
        sendTimesPlaceholder: "Send times",
        sendMessage: "Send Message"
    },
    consumer: {
        title: "Consumer",
        desc: "A simple consumer client to view messages from a Topic.",
        requiredTopic: "Required: Topic",
        topicPlaceholder: "Select or search for a Kafka Topic",
        requiredMessagesCount: "Required: Number of messages to consume",
        messagesCountPlaceholder: "Number of messages to consume",
        pollTimeoutDescription: "Poll timeout: Default is 10s. If there's an exception or no consumable messages, it will time out",
        pollTimeoutPlaceholder: "Poll timeout",
        optionalGroup: "Optional: Group (Once selected, offsets will be committed automatically during consumption. Supports creating new groups)",
        groupPlaceholder: "Select or create Consumer Group",
        commitOffsetTooltip: "Whether to commit Offset after consumption",
        isLatest: "Default consume location",
        onlyTip: "Note: Only effective when Group is first consumed, subsequent changes are invalid",
        consumeMessage: "Consume Messages",
        firstConsumeTip: "First consume may take some time to rebalance",
        startTimestamp: "Default consume start time",
    },
    group: {
        title: "Consumer Group",
        member: "Members",
        warn: "Please switch from a specific topic to this page first"
    },
    inspection: {
        title: "Inspection",
        desc: "Inspect Kafka lag conditions.",
        topicsLabel: "Topics",
        topicPlaceholder: "Select or search for a Kafka Topic",
        groupLabel: "Group",
        groupPlaceholder: "Select or create Consumer Group",
        startInspection: "Start Inspection",
        stopInspection: "Stop Inspection",
        autoFetch: "Automatically fetch data every {interval} minutes",
        lagFormula: "Lag = End offset - Committed offset.",
        lag: 'Lag',
        commit: 'Committed offset',
        end: 'End offset',
    },
    settings: {
        title: "Settings",
        width: "Window Width",
        height: "Window Height",
        lang: "Language",
        theme: "Theme"
    },
    message: {
        noMemberFound: "No members found",
        saveSuccess: "Saved successfully",
        connectSuccess: "Connected successfully",
        fetchSuccess: "Fetched successfully",
        sendSuccess: "Sent successfully",
        selectGroupFirst: "Please select a Group first",
        selectTopic: "Please select a Topic",
        selectTopicGroup: "Please select a Topic and Group",
        connectErr: "Connection failed",
        addOk: "Added successfully",
        editOk: "Edited successfully, refreshing configuration",
        mustFill: "Please fill in all required fields",
        saveErr: "Failed to save",
        pleaseInput: "Please enter message content"
    },
    about: {
        title: "About",
        projectHomepage: "Project Homepage",
        kafkaKing: "Kafka-King",
        esClient: "Matching ES Client",
        esKing: "ES-King",
        technicalGroup: "Technical Exchange Group",
        qqGroup: "QQ Exchange Group",
        translate: 'Are there any problems with the translation? Report it or participate in the translation',
    },
    header: {
        desc: "A more user-friendly Kafka GUI",
        c_node: "Current Cluster",
        netErr: "Unable to connect to GitHub, please check your network",
        newVersion: "New version available",
        down: "Download now"
    }
};