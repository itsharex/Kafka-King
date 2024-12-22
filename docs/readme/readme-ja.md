<p align="center">
  <img src="../snap/icon.ico" alt="画像タイトル">
</p>
<h1 align="center">Kafka King </h1>

<div align="center">

![License](https://img.shields.io/github/license/Bronya0/Kafka-King)
![GitHub release](https://img.shields.io/github/release/Bronya0/Kafka-King)
![GitHub All Releases](https://img.shields.io/github/downloads/Bronya0/Kafka-King/total)
![GitHub stars](https://img.shields.io/github/stars/Bronya0/Kafka-King)
![GitHub forks](https://img.shields.io/github/forks/Bronya0/Kafka-King)

<h3 align="center">現代的で実用的なKafka GUIクライアント</h3>

</div>

Kafkaの使い勝手を向上させる。

このプロジェクトは、複数のプラットフォームに対応したKafka GUIクライアントです。オープンソースへの支援としてスターマークをお願いします。ありがとうございます！

> 同様に便利なElasticsearchクライアント `ES-King` もご確認ください：https://github.com/Bronya0/ES-King


# Kafka-Kingの機能一覧
- [x] クラスタノードリストの表示、BrokerとTopic設定の動的構成サポート。
- [x] 消費者クライアントのサポート、指定されたグループによる特定のトピック、サイズ、タイムアウトでの消費、メッセージの詳細情報を表形式で表示。
- [x] PLAIN、SSL、SASL、Kerberos、sasl_plaintextなどへの対応。
- [x] トピックの作成（バッチ処理もサポート）、削除、レプリカとパーティションの指定。
- [x] 各コンシューマーグループごとに各トピックの全メッセージ数、コミット数、遅延量の統計。
- [x] トピックのパーティションの詳細情報（オフセット）の表示、追加のパーティションの追加サポート。
- [x] 生産者のシミュレーション、ヘッダーとパーティション指定でのメッセージの一括送信。
- [x] トピックとパーティションのヘルスチェック（完了済み）。
- [x] コンシューマーグループと個々の消費者の表示。
- [x] オフセット検査レポート。

# ダウンロード
右側からダウンロードするか、[ダウンロードページ](https://github.com/Bronya0/Kafka-King/releases) を開いて 【Assets】 を展開し、自分のプラットフォーム向けのバージョンを選択してください。Windows、macOS、Linuxをサポートしています。

`重要な注意点:`

1. **使用前に、Kafkaクラスタの`advertised.listeners`設定が正しく行われていることを確認してください。未設定またはドメイン名を使用している場合、接続先のドメイン名解像度エントリをローカルマシンのhostsファイルに追加して、ドメイン名の解像度により引き起こされる接続問題を避けてください。**
2. **SSLが必要な接続の場合、TLSを有効にして認証を無視してください（証明書がある場合は、TLS認証を有効にして証明書パスを入力）。**
3. **SASLユーザーはSASLを有効にし、適切なSASLプロトコル（通常はplain）を選択し、ユーザー名とパスワードを入力してください。**
4. **webview2ランタイムエラーが発生した場合は、Microsoftの公式ウェブサイトから最新のランタイムをダウンロードして再インストールしてください：https://developer.microsoft.com/ja-jp/microsoft-edge/webview2**



# ビルド
ソースコードを研究する場合のみ、手動でビルドが必要です。

cd app

wails dev


# Star履歴
[![Stargazers over time](https://starchart.cc/Bronya0/Kafka-King.svg)](https://starchart.cc/Bronya0/Kafka-King)


# ライセンス
Apache-2.0ライセンス

# 感謝
- wails: https://wails.io/docs/gettingstarted/installation
- naive ui: https://www.naiveui.com/
- franz-go: https://github.com/twmb/franz-go/
- xicons: https://xicons.org/#/

# 翻訳する
中国語、日本語、英語、韓国語、ロシア語などの言語をサポート

新しい言語の修正または追加：https://github.com/Bronya0/Kafka-King/issues/51