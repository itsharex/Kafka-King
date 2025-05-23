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
        enter: "확인",
        cancel: "취소",
        name: "닉네임",
        save: "저장",
        check: "선택해주세요",
        add: "생성",
        edit: "편집",
        delete: "삭제",
        deleteFinish: "삭제 성공",
        deleteOk: "정말로 삭제하시겠습니까?",
        count: "수량",
        refresh: "새로 고침",
        read: "읽기",
        config: "설정",
        csv: "CSV로 내보내기",
        connecting: "연결 중...",
        action: "동작",
        compress: '압축',
        decompress: '압축 해제',
    },
    aside: {
        cluster: "클러스터",
        node: "노드",
        topic: "토픽",
        producer: "프로듀서",
        consumer: "컨슈머",
        group: "컨슈머 그룹",
        monitor: "점검",
        settings: "설정"
    },
    conn: {
        title: "클러스터",
        add: "클러스터 추가",
        edit: "연결 편집",
        test: "연결 테스트",
        add_link: "연결 추가",
        please_add_name: "닉네임을 입력해 주세요",
        please_add_link: "연결 주소를 입력해 주세요",
        input_name: "이름을 입력해 주세요",
        bootstrap_servers: "연결 주소",
        tip: "참고: 로컬에서 Kafka 설정의 advertised.listeners 주소에 접근할 수 있어야 합니다 (특히 도메인 해석, IP를 사용하더라도 로컬에서 hosts를 설정해야 합니다)",
        tls: "TLS 사용",
        skipTLSVerify: "TLS 검증 생략",
        tls_cert_file: "PEM 인증서 경로 입력",
        tls_key_file: "키 개인 경로 입력",
        tls_ca_file: "CA 인증서 경로 입력",
        use_sasl: "SASL 사용",
        sasl_mechanism: "SASL 메커니즘",
        sasl_user: "SASL 사용자 이름",
        sasl_pwd: "SASL 비밀번호",
        kerberos_user_keytab: "Kerberos keytab 파일 경로",
        kerberos_krb5_conf: "Kerberos krb5.conf 경로",
        Kerberos_user: "Kerberos 사용자",
        Kerberos_realm: "Kerberos 도메인",
        kerberos_service_name: "Kerberos 서비스 이름"
    },
    node: {
        title: "노드",
        source: "출처",
        value: "값 (더블 클릭하여 편집 가능)",
        sensitive: "민감한가요?",
        ok_message: "편집 성공, 설정 새로 고침"
    },
    topic: {
        title: "토픽",
        add: "토픽 생성",
        add_name: "토픽 이름 입력",
        selectedGroup: "그룹 선택 및 오프셋 읽기",
        partition: "파티션",
        add_partition: "파티션 추가",
        add_partition_count: "추가할 파티션 수",
        replication_factor: "복제 수",
        err: "토픽 오류",
        lag: "지연",
        viewProduce: "메시지 생산",
        viewConsumer: "메시지 보기/소비",
        viewOffset: "오프셋 읽기",
        viewConfig: "토픽 설정",
        viewPartition: "파티션 보기",
        deleteTopic: "토픽 삭제"
    },
    producer: {
        title: "프로듀서",
        desc: "지정된 토픽으로 메시지를 전송하는 프로듀서 클라이언트입니다.",
        selectTopic: "토픽 선택",
        topicPlaceholder: "필수: 토픽 선택 또는 검색",
        optionalMessageKey: "선택 사항: 메시지 키 입력",
        keyPlaceholder: "선택 사항: 메시지 키 입력",
        specifyPartition: "선택 사항: 파티션 번호 지정",
        messageContentPlaceholder: "필수: 메시지 내용, 문자열 형식, JSON 지원",
        headersTitle: "메시지 헤더:",
        addHeader: "헤더 추가",
        removeHeader: "삭제",
        headerKeyPlaceholder: "헤더 키",
        headerValuePlaceholder: "헤더 값",
        sendTimes: "전송 횟수",
        sendTimesPlaceholder: "전송 횟수",
        sendMessage: "메시지 전송"
    },
    consumer: {
        title: "컨슈머",
        desc: "토픽의 메시지를 확인하는 간단한 컨슈머 클라이언트입니다.",
        requiredTopic: "필수: 토픽",
        topicPlaceholder: "토픽 선택 또는 검색",
        requiredMessagesCount: "필수: 소비할 메시지 수",
        messagesCountPlaceholder: "소비할 메시지 수",
        pollTimeoutDescription: "poll 타임아웃: 기본 10초. 예외나 소비할 메시지가 없으면 타임아웃됩니다",
        pollTimeoutPlaceholder: "poll 타임아웃",
        optionalGroup: "선택 사항: 그룹 (한번 선택하면 소비 시 오프셋이 자동으로 제출됩니다. 새로운 그룹 생성도 지원)",
        groupPlaceholder: "컨슈머 그룹 선택 또는 생성",
        commitOffsetTooltip: "소비 후 오프셋 제출 여부",
        firstConsumeTip: "첫 번째 소비는 오프셋을 밸런싱하는 데 시간이 걸릴 수 있습니다",
        consumeMessage: "메시지 소비"
    },
    group: {
        title: "컨슈머 그룹",
        member: "멤버",
        warn: "구체적인 토픽에서 이 페이지로 전환해 주세요"
    },
    inspection: {
        title: "점검",
        desc: "Kafka 지연 상황을 점검합니다.",
        topicsLabel: "토픽",
        topicPlaceholder: "토픽 선택 또는 검색",
        groupLabel: "그룹",
        groupPlaceholder: "컨슈머 그룹 선택 또는 생성",
        startInspection: "점검 시작",
        autoFetch: "데이터를 5분마다 자동으로 가져오기",
        lagFormula: "지연 = 최종 오프셋 - 제출된 오프셋.",
        lag: '지연',
        commit: '커밋된 오프셋',
        end: '종료 오프셋',
    },
    settings: {
        title: "설정",
        width: "창 너비",
        height: "창 높이",
        lang: "언어",
        theme: "테마"
    },
    message: {
        noMemberFound: "멤버를 찾을 수 없습니다",
        saveSuccess: "저장되었습니다",
        connectSuccess: "연결되었습니다",
        fetchSuccess: "조회되었습니다",
        sendSuccess: "전송되었습니다",
        selectGroupFirst: "먼저 그룹을 선택해 주세요",
        selectTopic: "토픽을 선택해 주세요",
        selectTopicGroup: "토픽과 그룹을 선택해 주세요",
        connectErr: "연결 실패",
        addOk: "추가되었습니다",
        editOk: "편집 완료, 설정 새로 고침",
        mustFill: "모든 필수 필드를 채워주세요",
        saveErr: "저장 실패",
        pleaseInput: "메시지 내용을 입력해 주세요",
    },
    about: {
        title: "정보",
        projectHomepage: "프로젝트 홈페이지",
        kafkaKing: "Kafka-King",
        esClient: "동일한 ES 클라이언트",
        esKing: "ES-King",
        technicalGroup: "기술 교류 그룹",
        qqGroup: "QQ 교류 그룹",
        translate: "번역: '번역에 문제가 있나요? 신고 또는 번역 참여'",
    },
    header: {
        desc: "더 사용자 친화적인 Kafka GUI",
        c_node: "현재 클러스터",
        netErr: "GitHub에 연결할 수 없습니다. 네트워크를 확인해 주세요",
        newVersion: "새 버전 발견",
        down: "즉시 다운로드"
    }
}
