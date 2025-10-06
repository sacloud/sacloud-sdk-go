# sacloud/go-template

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/go-template.svg)](https://pkg.go.dev/github.com/sacloud/go-template)
[![Tests](https://github.com/sacloud/go-template/workflows/Tests/badge.svg)](https://github.com/sacloud/go-template/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/go-template)](https://goreportcard.com/report/github.com/sacloud/go-template)

さくらのクラウド向けOSSプロダクトでのプロジェクトテンプレート(Go)

## 概要

さくらのクラウド向けOSSプロダクトでGo言語を中心に用いるプロジェクトのためのテンプレート

## 使い方

GitHubでリポジトリを作成する際にテンプレートとしてsacloud/go-templateを選択して作成します。  
![テンプレートの選択](docs/new_repo.png)

次に`go-teplate`という文字列を自身のプロジェクトのものに置き換えてください。

例: exampleという名前のプロジェクトを作成する場合

```bash
# 作成したプロジェクトのディレクトリに移動
cd example
# 置き換え
find . -type f | xargs sed -i '' -e "s/go-template/example/g"
```

### DockerイメージをGitHub Container Registryで公開する際の注意点

デフォルトでは`CR_PAT`が渡されないためGitHub Actionsでのイメージのビルド/プッシュに失敗します。
また、パッケージを公開したい場合は初回のみ手作業が必要です。

このためDockerイメージをGitHub Container Registryで公開したい場合はオーガニゼーション管理者にご相談ください。

## License

`go-template` Copyright (C) 2022-2025 The sacloud/go-template authors.
This project is published under [Apache 2.0 License](LICENSE).