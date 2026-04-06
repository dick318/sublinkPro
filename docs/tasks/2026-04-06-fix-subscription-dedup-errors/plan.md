# 修复订阅管理报错、节点去重开关与保存性能计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**目标：** 修复订阅管理中因重复节点关系和订阅日志外键导致的常见报错，新增“订阅输出同名节点去重”全局开关默认关闭，并优化新增/编辑订阅时的保存性能。

**非目标：**
- 不修改机场拉取阶段现有的跨机场去重、按协议去重或 `DeduplicationRule` 规则语义。
- 不新增数据库迁移，也不调整现有表结构。
- 不重构订阅编辑页、节点选择器或设置页的整体交互结构。
- 不在本轮自动执行 `git commit`、`git push`、`git merge` 或分支切换。

**方案摘要：**
本次修复按“先测后改、后端兜底、前端减噪、跨层同步”的最小闭环推进。先补充后端回归测试，固定三类高频问题：一是删除订阅时未先清理 `sub_logs`，导致 PostgreSQL 在硬删除 `subcriptions` 时触发外键错误；二是订阅保存链路允许重复 `node_id` 落入 `subcription_nodes`，导致复合主键冲突；三是主人明确指出的“订阅管理 -> 添加订阅弹窗 -> 右下角确定”这一条保存链路，在前端 `webs/src/views/subscriptions/component/SubscriptionFormDialog.jsx` 触发 `webs/src/views/subscriptions/index.jsx` 的 `handleSubmit` 后，后端 `api/sub.go` 逐个循环 `nodeIds` 并反复 `GetByID()`，再逐条写入节点关系，选中节点较多时会形成明显的 N+1 查询和多次小写入。随后在 `models/subcription.go` 中为节点关系写入增加稳定去重，并把订阅删除改成事务化清理子表后再删除主记录；在 `api/sub.go` 中优先复用现有的 `models.GetNodesByIDs()` 批量取节点，减少保存时的查询次数。

前后端需要同步收口。后端在 `api/setting.go` 扩展节点去重设置接口，返回并保存“跨机场节点去重”和“订阅输出同名节点去重”两个开关；前端在 `webs/src/views/settings/components/NodeDedupSettings.jsx` 增加第二个开关，并调整提示文案，明确两个开关分别作用于“节点入库”和“订阅输出”。同时在 `webs/src/views/subscriptions/index.jsx` 的节点选择、批量添加和提交前逻辑中补一层基于节点 ID 的前端去重，减少重复请求数据进入后端；若后端保存链路仍有明显耗时，再把订阅关联写入从逐条 `Create` 收敛为批量写入。最后同步更新相关功能文档，说明新开关默认关闭，以及它与现有跨机场去重的区别。

**执行步骤：**
1. 在 `docs/tasks/2026-04-06-fix-subscription-dedup-errors/` 下创建 `plan.md` 与 `task.md`，固定范围、风险和验收标准。
2. 为订阅模型与设置接口补充失败优先的回归测试，覆盖“删除带日志的订阅”“重复节点 ID 保存”“节点去重配置默认值/保存值”三条路径。
3. 修改 `api/sub.go`，复用 `models.GetNodesByIDs()` 批量读取节点，并梳理新增/编辑订阅保存时的节点、分组、脚本参数处理，减少逐条查询带来的耗时。
4. 修改 `models/subcription.go`，为 `AddNode()` / `UpdateNodes()` 增加稳定去重与批量写入保护，并在 `Del()` 中先清理 `sub_logs` 后再硬删除订阅。
5. 修改 `models/subcription.go`、`api/setting.go` 与设置缓存读取逻辑，为订阅输出按名称去重增加全局开关，默认关闭，同时保持旧配置兼容。
6. 修改 `webs/src/views/settings/components/NodeDedupSettings.jsx` 与 `webs/src/views/subscriptions/index.jsx`，新增设置开关并在前端节点选择链路做 ID 去重。
7. 同步更新 `docs/features/airport.md` 或其他已存在且最贴近该功能的文档，说明两个去重开关的区别与默认行为，并在任务记录中补充“保存性能优化”说明。
8. 执行与本次改动直接相关的验证，更新 `task.md` 状态，并记录限制项、风险和最小回滚方式。

**影响文件（预估）：**
- 修改：`models/subcription.go`
- 修改：`api/sub.go`
- 修改：`api/setting.go`
- 修改：`webs/src/views/subscriptions/component/SubscriptionFormDialog.jsx`
- 修改：`webs/src/views/settings/components/NodeDedupSettings.jsx`
- 修改：`webs/src/views/subscriptions/index.jsx`
- 修改：`docs/features/airport.md`
- 可能新增：`models/subcription_test.go`
- 可能新增：`api/setting_test.go`

**建议验证命令：**
- `go test ./models -run "TestSubcription" -count=1`
- `go test ./api -run "TestSub(Add|Update|NodeDedupConfig)" -count=1`
- `go test ./api -run "TestNodeDedupConfig" -count=1`
- `yarn run lint`
- `yarn run build`

**风险点：**
- `GetSub()` 输出逻辑变更后，可能影响依赖“同名节点自动去重”的旧订阅结果；需要用默认值和提示文案明确新语义。
- 订阅保存链路改为批量取节点和批量写入后，需要确认节点顺序、分组顺序和脚本顺序不被破坏。
- 删除订阅逻辑改为先删子表后删主表后，若事务处理不完整，可能留下部分关联数据未清理。
- 前端补去重后，已有依赖“重复选择可视化显示次数”的极少数使用方式会改变，但该变化与当前报错修复方向一致。
- 当前工作区已有未提交改动，执行阶段必须避免覆盖与本次任务无关的用户修改。

**回滚思路：**
- 若后端修复引发回归，优先回滚 `models/subcription.go` 与对应测试改动，恢复原有保存/删除链路。
- 若保存性能优化引发排序或关联异常，可单独回滚 `api/sub.go` 与订阅关联写入逻辑，保留外键和去重修复。
- 若新开关语义不符合预期，可单独回滚 `api/setting.go`、`webs/src/views/settings/components/NodeDedupSettings.jsx` 和文档说明，不影响已有订阅基础功能。
- 若仅前端去重策略需要调整，可单独回滚 `webs/src/views/subscriptions/index.jsx`，保留后端兜底修复。
