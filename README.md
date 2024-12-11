# microservice architecture demo

## 起動

```
docker-compose up -d
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

## ディレクトリ構成の説明

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
