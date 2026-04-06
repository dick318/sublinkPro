# 任务清单：修复订阅管理报错、节点去重开关与保存性能

- [x] 任务 1：确认问题根因与跨层影响范围。
验收点：已确认 `sub_logs` 外键删除错误、`subcription_nodes` 复合主键冲突、订阅输出阶段按名称二次去重，以及“添加订阅弹窗右下角确定”触发的逐条查节点慢链路这四条真实触发链路。

- [x] 任务 2：补充后端回归测试。
验收点：已有测试覆盖“删除带日志的订阅不会触发外键错误”“重复节点 ID 不会再写出重复关联”“节点去重设置接口默认返回订阅输出去重关闭”，并为订阅保存辅助逻辑提供最小回归保护。

- [x] 任务 3：优化新增/编辑订阅的保存链路性能。
验收点：`api/sub.go` 不再逐个 `GetByID()` 查节点，“添加订阅弹窗右下角确定”对应的保存链路优先使用批量节点读取，节点较多时请求耗时明显下降且节点顺序保持稳定。

- [x] 任务 4：修复订阅模型的删除与节点关系写入逻辑。
验收点：`models/subcription.go` 在删除订阅时会先清理关联日志，保存节点关系时会稳定去重且不再触发 `subcription_nodes_pkey` 冲突。

- [x] 任务 5：新增订阅输出同名节点去重开关并同步前后端。
验收点：设置页可同时配置“跨机场节点去重”和“订阅输出同名节点去重”，后者默认关闭，且订阅实际输出结果与开关状态一致。

- [x] 任务 6：同步更新文档并完成最小验证。
验收点：相关文档已说明两个去重开关的作用边界和保存性能优化范围，且已完成至少 1 次后端测试与 1 次前端校验或明确记录无法执行的原因。

## 验证记录

- 已执行：`go test ./models -run "TestSubcription(AddNodeSkipsDuplicateNodeIDs|UpdateNodesReplacesRelationsWithDedupedIDs|DelRemovesAssociatedSubLogs|GetSubKeepsSameNameNodesWhenOutputDedupDisabled)" -count=1`
结果：通过。
- 已执行：`go test ./api -run "Test(SubAddKeepsDistinctSameNameNodesAndDedupesRepeatedIDs|GetNodeDedupConfigReturnsSubscriptionOutputFlagDisabledByDefault|UpdateNodeDedupConfigPersistsSubscriptionOutputFlag)" -count=1`
结果：通过。
- 已执行：`corepack yarn run build`
结果：通过。
- 已执行：`corepack yarn run lint`
结果：命令执行完成，当前输出仍包含仓库既有的全仓 `prettier/prettier` CRLF 警告，本次改动未额外引入新的逻辑错误。

## 当前已知风险

- 订阅输出同名节点去重默认改为关闭后，部分历史订阅结果条数会增加，需要通过设置文案和文档明确说明。
- 保存性能优化如果改动到节点读取/关系写入顺序，需要额外关注生成订阅时的节点顺序是否仍与用户选择一致。
- 工作区当前存在其他未提交修改，执行阶段必须逐文件核对，不能误覆盖主人已有改动。
