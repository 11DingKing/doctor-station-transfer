# 博士工作站课题协作与成果转化

纯后端 Go 服务，围绕课题申报、专家评审、合同预算、里程碑、样品数据交接和审计形成可恢复状态流。SQLite 使用版本化 migration，服务启动自动建库；默认管理员 `admin@doctor.local` / `admin` 仅用于本地开发。

运行：`go run ./cmd/server`。健康检查：`GET /healthz`，就绪检查：`GET /readyz`。登录后携带 `Authorization: Bearer <token>` 调用 `/v1/projects`、`/v1/reviews` 等接口。
