# Foundation 规格冻结表

- foundation_profile: compact_10；task_count: 10；project_variant: backend。
- 业务边界：课题征集→意向申报→专家评审/回避→合同预算→阶段拨款→样品与数据交接→里程碑→延期/暂停→结题验收→知识产权许可收益；参与者为卫健委/科协、山西中医药大学负责人、羊头山项目经理、外部专家、财务审计。
- 持久化：SQLite 真实 SQL，版本化 migrations，users/projects/reviews/contracts/milestones/transfers/audits/sessions/idempotency_jobs 共 10 张关联表；外键、唯一约束、索引、时间戳和版本号。
- 事务：合同创建同时写项目、合同、审计；拨款同时更新预算与流水；交接同时写 transfer 与审计；失败回滚有测试。
- 状态机：project draft/submitted/reviewing/approved/contracted/active/suspended/completed/rejected；review pending/accepted/recused/submitted；milestone pending/accepted/overdue。
- 并发：项目 version 乐观锁，预算条件更新，评审名额唯一约束；并发测试使用同步屏障并以 race 检查。
- context：HTTP→service→repository 全链路传递 context，超时和取消保留 errors.Is 链。
- worker：后台处理过期里程碑，重试/退避/永久失败写 jobs 表，支持优雅停止与重启恢复。
- 错误传播：稳定 code/message/request_id JSON，包装底层错误，不吞错。
- HTTP：/healthz、/readyz、/v1/auth/login、/logout、/v1/projects、/reviews、/milestones，鉴权中间件和角色权限。
- 身份：可撤销会话、过期检查、退出撤销；roles admin/researcher/reviewer/auditor，服务层强制角色差异。
- Docker：真实 Dockerfile 从 cmd/server 构建，linux/amd64 与 linux/arm64 镜像、health/ready 检查。
- 测试：领域、service、真实 SQLite migration/重启/事务、HTTP、worker、并发、分页过滤、鉴权；目标测试 Go ≥1500 行。
- 规模：非测试生产 Go ≥2000 行、≥20 文件、≥10 package；禁止空壳、重复或注释凑数。
- 后续容量：10 个独立边界（事务回滚、状态迁移、预算并发、评审回避、幂等、context 取消、worker 重试、审计一致性、会话撤销、分页组合查询）；本阶段不植入 Bug、不创建题目分支/私测/intake。
