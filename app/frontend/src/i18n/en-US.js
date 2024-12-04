export default {
    common: {
        enter: 'Confirm',
        cancel: 'Cancel',
        name: 'Name',
        save: 'Save',
        check: 'Please select',
    },
    aside: {
        cluster: 'Cluster',
        node: 'Node',
        topic: 'Topic',
        producer: 'Producer',
        consumer: 'Consumer',
        group: 'Consumer Group',
        monitor: 'Inspection',
        settings: 'Settings'
    },
    conn: {
        title: 'Cluster',
        add: 'Add Cluster',
        edit: 'Edit Connection',
        test: 'Test Connection',
        add_link: 'Add Connection',
        please_add_name: 'Please enter a nickname',
        please_add_link: 'Please enter the connection address',
        input_name: 'Enter name',
        bootstrap_servers: 'Connection Address',
        tip: 'Note: Ensure that the advertised.listeners address configured in Kafka is accessible locally (especially DNS resolution, even if you fill in an IP, it must be properly configured in the local hosts file)',
        tls: 'Use TLS',
        skipTLSVerify: 'Skip TLS Verification',
        tls_cert_file: 'Enter PEM certificate path',
        tls_key_file: 'Enter key private key path',
        tls_ca_file: 'Enter CA certificate path',
        use_sasl: 'Use SASL',
        sasl_mechanism: 'SASL Mechanism',
        sasl_user: 'SASL Username',
        sasl_pwd: 'SASL Password',
        kerberos_user_keytab: 'Kerberos keytab file path',
        kerberos_krb5_conf: 'Kerberos krb5.conf path',
        Kerberos_user: 'Kerberos Username',
        Kerberos_realm: 'Kerberos Realm',
        kerberos_service_name: 'Kerberos Service Name',
    },
    settings: {
        title: 'Settings'
    }
};