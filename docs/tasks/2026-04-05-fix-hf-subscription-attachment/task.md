# 任务清单：修复 HF 订阅附件型响应兼容性

- [x] 任务 1：完成问题定位与样本确认。
验收点：已确认样本订阅返回 `application/octet-stream`、`Content-Disposition: attachment`，正文实际为可解析的 Base64 订阅文本。

- [x] 任务 2：补充订阅解析回归测试。
验收点：测试覆盖附件风格 Base64 订阅正文，以及带 UTF-8 BOM 的原始多行节点文本。

- [x] 任务 3：实现并接入统一的订阅载荷解析逻辑。
验收点：`node/sub.go` 不再依赖分散的字符串分支，解析逻辑可稳定处理附件型订阅响应。

- [x] 任务 4：同步更新 Hugging Face 部署文档。
验收点：部署文档明确改为使用当前仓库源码构建，而不是继续依赖可能滞后的官方镜像。

- [x] 任务 5：完成本次可执行验证并整理部署结果。
验收点：已记录实际执行的验证命令、限制项与重新部署所需后续动作。

## 验证记录

- 已执行：将 Go `1.26.1` 安装到 `K:\Program\go`，并将 `K:\Program\go\bin` 加入用户级 `PATH`。
- 已执行：`K:\Program\go\bin\go.exe version`，结果为 `go version go1.26.1 windows/amd64`。
- 已执行：对样本订阅 `https://cdn2.xmrth.one/s/3kQWEwaMj4Po-41748746?sub=2` 的手工抓取确认，直连返回 `200`、`Content-Disposition: attachment` 和可解析的 Base64 正文。
- 已确认：线上报错与手工直连样本不一致，当前更可疑的是“代理下载链路与直连链路拿到不同响应内容”。
- 已执行：`git diff --check -- node/sub.go node/sub_test.go`，完成本地静态检查，未发现空白符错误。
- 已执行：`git diff -- node/sub.go node/sub_test.go`，完成本地静态自查。
- 已执行：`K:\Program\go\bin\go.exe test ./node -run 'TestParseSubscriptionPayload|TestParseSubscriptionWithFallback|TestFetchSubscriptionResponse' -count=1`，结果通过。
- 已确认：同一条订阅 URL 在本机直连仍返回 `200`，但 Hugging Face 线上环境会收到 `403` 和 `{"code":"OUTDATED_TOKEN"}`，说明当前问题已从“解析兼容性”转为“来源环境差异导致的上游拒绝”。
- 已执行：`git push origin main`，当前 GitHub 提交为 `630abdf8483adc638d8387f54b7ee07aa0030907`。
- 已执行：将 Hugging Face Space `liu4318/sublinkpro-postgres` 的 `SUBLINK_REF` 更新到 `630abdf8483adc638d8387f54b7ee07aa0030907`。
- 已执行：轮询 Hugging Face Space API，确认状态从 `RUNNING_BUILDING` -> `RUNNING_APP_STARTING` -> `RUNNING`，最终 `runtime.sha` 已切换为 `aa36c82d6b68f4c6d08be5b4c1aa0a549118e181`。
- 已执行：访问 `https://liu4318-sublinkpro-postgres.hf.space/login`，返回 `200`。

## 限制项

- 当前终端已具备可复用 Go 环境，但 Hugging Face 线上 403 是否完全消失，仍需再次以真实机场更新任务做回归确认。

## 重新部署后续动作

1. 将当前修复提交推送到 GitHub，并把 Hugging Face Space 的源码构建引用更新到新提交。
2. 触发 `liu4318/sublinkpro-postgres` 重新构建，确认日志中完成前端构建与 `go build -tags=prod`。
3. 构建完成后，重新触发一次该机场订阅更新任务。
4. 核对新日志是否出现“代理响应解析失败，已回退直连成功”或直接解析成功，并确认节点已正常入库。

## 本轮补充说明

- 已检查前端机场表单与相关 API 字段：本次改动仅发生在后端订阅下载链路，前端交互与字段语义无需同步修改。
- 已检查部署文档语义：本轮新增的是运行时容错与诊断日志，不改变部署命令和配置含义，因此无需额外修改部署文档。
