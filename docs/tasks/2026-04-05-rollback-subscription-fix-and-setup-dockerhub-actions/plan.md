# 回滚订阅修复并接入 GitHub Actions 推送 Docker Hub 计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**目标：** 回滚 2026-04-05 的订阅修复代码，保留可追溯的修改记录，并为仓库补齐“每次提交自动构建并推送到 Docker Hub”的 GitHub Actions 工作流与配置说明。

**非目标：**
- 不继续扩展订阅解析逻辑，也不在本轮重新修复订阅下载问题。
- 不修改现有业务配置语义、数据库结构或前端页面行为。
- 不在本轮执行 `git commit`、`git push` 或重新部署远端服务。

**方案摘要：**
先以 `808bbe2` 的父提交为基准回退订阅链路代码，仅移除当前分支中那次实验性修复带来的运行时行为变化，并同步撤下会误导当前部署语义的 Hugging Face 部署文档。为了方便后续再次排查，再新增一份任务记录，明确写出当时改过什么、为什么先回滚、后续如果要重做应从哪些文件和日志入手。随后直接调整现有的 GitHub Actions workflow，让它在每次 `push` 后构建前端、构建多架构镜像并推送到 Docker Hub；镜像仓库名改为通过仓库变量配置，登录凭据通过 Actions secrets 配置，避免继续写死在 workflow 里。

**执行步骤：**
1. 创建本次任务目录，写入 `plan.md` 与 `task.md`，明确回滚和自动化交付边界。
2. 回滚 `node/sub.go` 与相关当前部署说明，确保运行时行为回到 `808bbe2` 之前。
3. 删除仅为该次修复新增的测试文件，并新增一份回滚记录文档，保留后续重做所需线索。
4. 调整 GitHub Actions Docker Hub 自动推送 workflow，并补充手动配置 secrets / variables 的部署说明文档。
5. 执行与本次改动直接相关的校验，更新 `task.md` 状态，并整理变更风险与回滚方式。

**风险点：**
- 回滚后，若“附件响应解析失败”问题真实存在，订阅问题会恢复到修复前状态。
- 新增 workflow 若与现有发布流水线标签策略冲突，可能造成镜像标签混淆。
- Docker Hub 推送依赖用户自行配置 secrets / variables，缺少任一项都会导致 workflow 失败。

**回滚思路：**
- 若要恢复本轮回滚前的订阅修复，可重新 cherry-pick `808bbe2`，或根据本次新增回滚记录重新实现。
- 若新增的 GitHub Actions workflow 不合适，可单独删除对应 workflow 与说明文档，不影响当前应用运行。
- 若文档说明有误，可只回滚文档文件，不需要动业务代码。
