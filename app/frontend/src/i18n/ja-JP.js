export default {
    common: {
        enter: '確認', // Confirm
        cancel: 'キャンセル', // Cancel
        name: '名前', // Name
        save: '保存', // Save
        check: '選択してください', // Please select
        edit: '編集',
        delete: '削除',
        deleteOk: '本当に削除してもよろしいですか?',
    },
    aside: {
        cluster: 'クラスタ', // Cluster
        node: 'ノード', // Node
        topic: 'トピック', // Topic
        producer: 'プロデューサー', // Producer
        consumer: 'コンシューマー', // Consumer
        group: 'コンシューマグループ', // Consumer Group
        monitor: '巡回検査', // Inspection
        settings: '設定' // Settings
    },
    conn: {
        title: 'クラスタ', // Cluster
        add: 'クラスタ追加', // Add Cluster
        edit: '接続編集', // Edit Connection
        test: '接続テスト', // Test Connection
        add_link: '接続追加', // Add Connection
        please_add_name: 'ニックネームを入力してください', // Please enter a nickname
        please_add_link: '接続アドレスを入力してください', // Please enter the connection address
        input_name: '名称を入力', // Enter name
        bootstrap_servers: '接続アドレス', // Connection Address
        tip: '注意：Kafkaで設定されたadvertised.listenersアドレスがローカルからアクセス可能であることを確認してください（特にDNS解決、IPを入力しても、ローカルのhostsファイルに適切に設定する必要があります）',
        tls: 'TLSを使用', // Use TLS
        skipTLSVerify: 'TLS検証をスキップ', // Skip TLS Verification
        tls_cert_file: 'PEM証明書のパスを入力', // Enter PEM certificate path
        tls_key_file: '秘密鍵ファイルのパスを入力', // Enter key private key path
        tls_ca_file: 'CA証明書のパスを入力', // Enter CA certificate path
        use_sasl: 'SASLを使用', // Use SASL
        sasl_mechanism: 'SASLメカニズム', // SASL Mechanism
        sasl_user: 'SASLユーザー名', // SASL Username
        sasl_pwd: 'SASLパスワード', // SASL Password
        kerberos_user_keytab: 'Kerberos keytabファイルのパス', // Kerberos keytab file path
        kerberos_krb5_conf: 'Kerberos krb5.confのパス', // Kerberos krb5.conf path
        Kerberos_user: 'Kerberosユーザー名', // Kerberos Username
        Kerberos_realm: 'Kerberosドメイン', // Kerberos Realm
        kerberos_service_name: 'Kerberosサービス名', // Kerberos Service Name
    },
    settings: {
        title: '設定',
        width: 'ウィンドウ幅',
        height: 'ウィンドウ高さ',
        lang: '言語',
        theme: 'テーマ',
    }
};