# 修复 HF 订阅附件型响应兼容性与重新部署计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**目标：** 修复机场拉取时对 `application/octet-stream + Content-Disposition: attachment` 类型订阅响应的兼容性问题，并完成 Hugging Face 的重新部署准备。

**非目标：**
- 不修改机场订阅业务配置语义。
- 不重构任务系统、通知系统或前端机场管理交互。
- 不变更 `/c/*` 下已有订阅输出协议。

**方案摘要：**
先为 `node/sub.go` 增加回归测试，覆盖“附件风格 Base64 订阅正文”和“带 BOM 的原始多行节点文本”两类输入，再抽取统一的订阅载荷归一化/解析函数，确保解析逻辑不依赖响应是否被浏览器视为下载文件。随后同步更新 Hugging Face 部署文档，从“依赖官方镜像”改为“直接构建当前仓库代码”，避免线上继续使用旧逻辑。

**执行步骤：**
1. 创建本次任务目录并写入 `plan.md`、`task.md`。
2. 先编写订阅解析回归测试，覆盖附件风格 Base64 和 BOM 原始文本场景。
3. 在 `node/sub.go` 中实现统一的订阅载荷归一化与解析函数，并替换现有分散解析逻辑。
4. 更新 Hugging Face 部署文档，明确使用当前仓库源码构建部署。
5. 尝试执行与本次改动直接相关的验证，并整理重新部署步骤。

**风险点：**
- 订阅解析逻辑修改后，可能影响既有 YAML / Base64 / 原始多行文本兼容性。
- 当前环境缺少 Go 工具链时，无法本地执行 `go test` 做完整验证。
- Hugging Face 重新部署若仍沿用旧镜像，将无法体现本次修复。

**回滚思路：**
- 回滚 `node/sub.go` 与新增测试文件。
- 若线上部署异常，回退 Hugging Face Space 到上一版已知可用提交或上一份 Docker 配置。
- 若仅文档策略不合适，单独回滚 `docs/deployments/sublinkpro-huggingface-spaces.md`。
