テストで使うためのコンテナ用Dockerfileとコマンド例

```
# コンテナレジストリはlinux/amd64しかサポートしていないため、arm64環境では--platformで指定する
$ docker build --platform linux/amd64 -t test .
$ docker login sakura-oss-dev.sakuracr.jp
$ docker tag test:latest sakura-oss-dev.sakuracr.jp/test:latest
$ docker push sakura-oss-dev.sakuracr.jp/test:latest
```
