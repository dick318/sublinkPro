# 回滚记录：2026-04-05 订阅修复尝试

## 1. 本次回滚对应的实验性修复

- 原实验提交：`808bbe2 fix: 修复机场订阅附件解析与403兼容问题`
- 本次回滚目标：先撤下这次修复带来的运行时行为变化，避免在根因未确认前长期保留一组带诊断和回退逻辑的实验代码。

## 2. 当时做过的修改

主要改动集中在 [node/sub.go](/e:/pyWorkSpace/sublinkPro/node/sub.go) 和新增测试 [node/sub_test.go](/e:/pyWorkSpace/sublinkPro/node/sub_test.go)：

1. 统一订阅载荷解析流程，尝试兼容：
   - `Content-Disposition: attachment`
   - Base64 订阅正文
   - 带 UTF-8 BOM 的原始多行节点文本
2. 下载订阅时默认把 `User-Agent` 设为 `Clash`。
3. 加入大量请求诊断日志，记录：
   - HTTP 状态码
   - `Content-Type`
   - `Content-Disposition`
   - 是否使用代理
   - 最终 URL
   - 响应体预览
4. 代理下载失败时尝试直连回退。
5. 直连返回 `403 + OUTDATED_TOKEN` 时尝试代理重试。
6. 透传上游 JSON 错误信息。

## 3. 为什么现在先回滚

后续线上日志显示，真正稳定复现的错误已经不是“附件响应导致解析失败”，而是上游直接返回：

```text
HTTP 403
{"error":"TOKEN过期, 请重新复制","code":"OUTDATED_TOKEN"}
```

同时，本机对同一条订阅 URL 直连抓取时拿到的是 `200` 和正常正文。这意味着当时的线上失败更像是：

- 服务器出口环境差异
- 代理链路差异
- 上游对请求来源的限制

而不是已经被确认的“附件下载型响应解析 bug”。

在根因还不够确定时，保留那组补丁会让当前分支承担额外复杂度，所以先回滚，只保留记录。

## 4. 如果后续要重新修

建议按下面顺序重新开始：

1. 先在目标部署环境复现实验日志，确认拿到的是 `200 + 正常正文`，还是 `403 + OUTDATED_TOKEN`。
2. 同时比较：
   - 直连请求
   - 代理请求
   - `User-Agent`
   - 返回头和正文预览
3. 如果再次确认“附件响应解析失败”真实存在，再从 [node/sub.go](/e:/pyWorkSpace/sublinkPro/node/sub.go) 的订阅下载入口重新补解析逻辑。
4. 重新补回针对订阅正文解析的测试，重点覆盖：
   - 附件风格 Base64 正文
   - 带 BOM 的原始节点文本
   - 代理失败后直连回退
   - 上游 JSON 错误透传

## 5. 参考线索

- 原任务记录：[task.md](/e:/pyWorkSpace/sublinkPro/docs/tasks/2026-04-05-fix-hf-subscription-attachment/task.md)
- 原计划文档：[plan.md](/e:/pyWorkSpace/sublinkPro/docs/tasks/2026-04-05-fix-hf-subscription-attachment/plan.md)
- 本次回滚任务：[task.md](/e:/pyWorkSpace/sublinkPro/docs/tasks/2026-04-05-rollback-subscription-fix-and-setup-dockerhub-actions/task.md)

当前分支已经撤下实验性代码，但保留这些记录，后续如果要再次尝试，可以直接从这里回溯。
