# sacloud/iaas-api-go

- URL: https://github.com/sacloud/iaas-api-go/pull/2
- Parent: https://github.com/sacloud/iaas-service-go/pull/1
- Author: @yamamoto-febc

## 概要

[iaas-service-goの基本方針](https://github.com/sacloud/iaas-service-go/pull/1)に従い、sacloud/libsacloud v2からIaaS部分を切り出す。

## やること/やらないこと

### やること

- libsacloudからのIaaS部分の切り出し
- iaas-api-go v1としてリリース
- libsacloudの`sacloud`パッケージ配下の整理
  - ~typesやostypeといったパッケージ構成の再考/整理~
   
### やらないこと

- libsacloud v2の独自DSLを含むlibsacloudの実装の改善
  基本的にlibsacloud v2をそのまま移植する。ただし、前述の`sacloud`パッケージ配下の整理などのリファクタレベルの修正は行う。
  従来[libsacloud v3として検討されてきた内容](https://github.com/sacloud/libsacloud/issues/791)はiaas-service-goやiaas-api-go v2で実現する。

## 実装

### libsacloudとiaas-api-goの並列開発

当面はlibsacloudの修正を継続する。libsacloudに対して行われた修正は手作業でiaas-api-goに取り込む。  
iaas-api-goへの移植は[libsacloud v2.32.2](https://github.com/sacloud/libsacloud/tree/v2.32.2)を元にする。  

### 方針

`sacloud`パッケージについて、libsacloudのクライアント側での修正が容易に行える程度の改修をしつつ移植する。
(容易 == 機械的に置き換えできる、という程度)

### 移植対象/対応

#### リポジトリ運用

[libsacloud v2.32.2](https://github.com/sacloud/libsacloud/tree/v2.32.2)を基点にソースコード類をコピーして移行する。  
libsacloudからのforkは行わず新たなリポジトリで開発していく。

#### libsacloudのパッケージ構成/移行対象

```console
- examples: otel利用例
- helper: 高レベルAPI群(一部をiaas-service-goへ)
- internal: 独自DSL
- pkg: libsacloudに依存しないユーティリティなど => packages-goへ
- sacloud
  - accessor
  - fake
  - naked
  - ostype
  - pointer => packages-goへ
  - profile => api-client-goへ
  - search
  - stub
  - test
  - testutil => 一部をpackages-goへ
  - trace
  - types
  - sacloud直下
```

- `profile`はapi-client-goで実装する  
- testutilは整理してから切り出し/分割などの対応が必要  
- ~typesは整理してからsacloud直下へ移動などの対応が必要~  

#### iaas-api-goのパッケージ構成

従来はsacloudパッケージ配下だったものをiaas-api-goの配下にする。  
パッケージ名は`iaas`とする。

```console
- accessor
- defaults => libsacloudのhelper/defaultsの移植
- fake
- helper
  - cleanup  
  - plans    
  - power    
  - query    
  - wait     
- internal => libsacloudの独自DSL実装など
- naked
- ostype
- search
- stub
- test
- testutil
- trace
- types
- sacloud直下
```

## 改訂履歴

- 2022/3/4: 初版作成
- 2022/3/7: libsacloud/v2直下のパッケージについて追記
- 2022/3/10: typesの改修を当面延期
- 2022/3/29: sacloud-goの分割/リネームを反映