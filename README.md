# microservice architecture demo-app

## 起動

```
docker compose up -d
```

## 停止

```
docker compose down
```

## docker 内の不要なデータを削除

```
docker system prune
```

## サービス単体の動作確認

```
cd service-1
docker build -t service-1 .
docker run -p 8080:8080 service-1
```

```
curl http://localhost:8080
```

## コンテナ内で操作を行うコマンド

```
docker exec -it microservice-architecture-demo sh
```

## サービスが起動しているか確認

```
docker-compose ps
```

## アクセス

```
▼ サービス1
http://localhost:8080/api/service-1/

▼ サービス2
http://localhost:8080/api/service-2/

▼ サービス3
http://localhost:8080/api/service-3/
```

## ディレクトリの基本構成の説明

```
microservice-architecture-demo/
├── conf/
│   └── localhost/               # Nginxの設定などが入る（ルーティングやプロキシ設定）
├── service-1/
│   ├── dockerfile               # service-1のDockerイメージを作成するためのファイル
│   └── main.go                  # service-1のアプリケーションコード
├── service-2/
│   ├── dockerfile               # service-2のDockerイメージを作成するためのファイル
│   └── main.go                  # service-2のアプリケーションコード
├── service-3/
│   ├── dockerfile               # service-3のDockerイメージを作成するためのファイル
│   └── main.go                  # service-3のアプリケーションコード
├── docker-compose.yml           # サービス全体を定義するComposeファイル
└── README.md
```

## service-1 のディレクトリ構成

```
service-1/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/
│   │   └── models/
│   │       └── stock.go
│   ├── usecase/
│   │   ├── stock_interactor.go
│   │   └── stock_repository.go
│   ├── interfaces/
│   │   ├── handlers/
│   │   │   └── stock_handler.go
│   │   └── repositories/
│   │       └── stock_repository.go
│   └── infrastructure/
│       └── database/
│           └── postgres.go
├── dockerfile
├── go.mod
├── go.sum
└── main.go
```

## 簡易在庫管理システム

マイクロサービスアーキテクチャを用いた簡易的な在庫管理システム

## システム構成

### service-1 (Stock Service)

在庫の基本的な CRUD 操作
在庫数の管理
最小在庫数の監視

### service-2 (Transaction Service)

入出庫履歴の記録
トランザクションログの管理
在庫移動の追跡

### service-3 (Alert Service)

在庫アラートの管理
在庫レポートの生成

### セットアップ

```shell
docker compose up -d
```

### 起動確認

```shell
docker compose ps
```

### コンテナとボリュームの停止・削除

```shell
docker compose down -v
```

### コンテナの再起動

```shell
docker compose up -d
```

# API 動作確認手順

## Service-1 (Stock Service)

### 商品登録

```shell
   curl -X POST http://localhost:8080/api/service-1/stocks \
    -H "Content-Type: application/json" \
    -d '{
   "id": "123",
   "name": "テスト商品",
   "description": "テスト用の商品です",
   "quantity": 50,
   "unit": "個",
   "min_quantity": 10
   }'
```

### 登録した商品の確認

```shell
curl http://localhost:8080/api/service-1/stocks/123
```

### 全商品の一覧取得

```shell
curl http://localhost:8080/api/service-1/stocks
```

### 特定商品の取得

```shell
curl http://localhost:8080/api/service-1/stocks/123
```

### 商品情報の更新

```shell
curl -X PUT http://localhost:8080/api/service-1/stocks/123 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新商品",
    "description": "更新後の商品です",
    "quantity": 50,
    "unit": "個",
    "min_quantity": 15
  }'
```

### 在庫数の直接更新

```shell
curl -X PATCH http://localhost:8080/api/service-1/stocks/123/quantity \
  -H "Content-Type: application/json" \
  -d '{
    "adjustment": 10,
    "note": "在庫調整"
  }'
```

### 商品の削除

```shell
curl -X DELETE http://localhost:8080/api/service-1/stocks/123
```

# Service-2 (Transaction Service)

## 入出庫処理

```shell
curl -X POST http://localhost:8080/api/service-2/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "stock_id": "123",
    "type": "in",
    "quantity": 20,
    "note": "入庫処理"
  }'
```

### 出庫処理

```shell
curl -X POST http://localhost:8080/api/service-2/transactions \
 -H "Content-Type: application/json" \
 -d '{
"stock_id": "123",
"type": "out",
"quantity": 30,
"note": "出庫テスト"
}'
```

### 在庫・トランザクション確認

```shell
   curl http://localhost:8080/api/service-1/stocks/123
```

### トランザクション履歴の確認

```shell
curl http://localhost:8080/api/service-2/stocks/123/transactions
```

### 特定のトランザクションの取得

```shell
curl http://localhost:8080/api/service-2/transactions/{transaction_id}
```

### 在庫サマリーの確認

```shell
curl http://localhost:8080/api/service-2/stocks/123/summary
```

### 日付範囲指定でのトランザクション取得

```shell
curl "http://localhost:8080/api/service-2/transactions?start_date=2024-01-01&end_date=2024-12-31"
```

## エラーケースのテスト

### アラート設定の作成

```shell
curl -X POST http://localhost:8080/api/service-3/configs \
  -H "Content-Type: application/json" \
  -d '{
    "stock_id": "123",
    "min_quantity": 20,
    "max_quantity": 100
  }'
```

### アラート一覧の取得

```shell
curl http://localhost:8080/api/service-3/alerts
```

### アクティブなアラートの取得

```shell
curl "http://localhost:8080/api/service-3/alerts?resolved=false"
```

### 特定商品のアラート取得

```shell
curl "http://localhost:8080/api/service-3/alerts?stock_id=123"
```

### アラート解決

```shell
curl -X POST http://localhost:8080/api/service-3/alerts/{alert_id}/resolve
```

### 在庫レポートの生成

```shell
curl http://localhost:8080/api/service-3/reports/stocks
```

### アラートレポートの生成

```shell
curl http://localhost:8080/api/service-3/reports/alerts
```

### 最小在庫数を下回った商品の確認

```shell
curl "http://localhost:8080/api/service-1/stocks?low_stock=true"
```

### エラーハンドリング

在庫不足の場合は 400 エラー
存在しない商品 ID の場合は 404 エラー
サーバーエラーの場合は 500 エラー

### 注意事項

ローカル環境での実行を前提としています
データベースは PostgreSQL を使用しています
初回起動時はデータベースのマイグレーションが自動で実行されます
