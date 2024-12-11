# local-microservice

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
cd service-a
docker build -t service-a .
docker run -p 8080:8080 service-a
```

```
curl http://localhost:8080
```

## Nginx が設定を正しく読み込んでいるか確認

```
docker exec -it local-microservice-nginx-proxy-1 sh
cat /etc/nginx/vhost.d/localhost
```

## サービスが起動しているか確認

```
docker-compose ps
```

## アクセス

```
▼ サービスA
http://localhost:8080/api/service-1/

▼ サービスB
http://localhost:8080/api/service-2/

▼ サービスC
http://localhost:8080/api/service-3/
```

## ディレクトリ構成の説明

```
LOCAL-MICROSERVICE/
├── conf/
│   └── localhost/               # Nginxの設定などが入る（ルーティングやプロキシ設定）
├── service-1/
│   ├── dockerfile               # service-aのDockerイメージを作成するためのファイル
│   └── main.go                  # service-aのアプリケーションコード
├── service-2/
│   ├── dockerfile               # service-bのDockerイメージを作成するためのファイル
│   └── main.go                  # service-bのアプリケーションコード
├── service-3/
│   ├── dockerfile               # service-cのDockerイメージを作成するためのファイル
│   └── main.go                  # service-cのアプリケーションコード
├── docker-compose.yml           # サービス全体を定義するComposeファイル
└── README.md                    # プロジェクトの説明
```
