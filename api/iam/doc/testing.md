# 結合テストに必要な権限まとめ

- IAMのテストを実行するにはIAMのAPIを叩くための権限が必要。
- "owner"は全部入りの権限ではない。"owner"だけ付与すればいいという話ではない。

## ./permit.go

- 手作業ではあまりにも難しすぎるのでサービスプリンシパルに十分な権限を付与する `./permit.go` を作成した。
- これは「十分な権限」なのであって「必要な権限」あるいは「必要十分な権限」ではない。注意が必要。

## IAM role

- `owner`
- `organization-admin`
- `servicepolicy-admin`
- `folder-admin`
- `project-creator`

があれば(付与しすぎの可能性は排除できないものの)テストが動きそうです

## ID role

- `identity-admin`

があれば(付与しすぎの可能性は排除できないものの)テストが動きそうです
