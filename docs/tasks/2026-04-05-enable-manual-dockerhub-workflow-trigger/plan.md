# 为 Docker Hub 构建增加手动触发入口计划

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**目标：** 为现有 GitHub Actions Docker Hub 构建流程增加手动触发入口，并补充配套说明，方便在不新提交代码时手动发起一次镜像构建与推送。

**非目标：**
- 不修改应用业务代码。
- 不调整 Docker 镜像标签规则。
- 不变更 Secrets / Variables 的命名。

**方案摘要：**
直接在现有的 [build-release.yml](/e:/pyWorkSpace/sublinkPro/.github/workflows/build-release.yml) 上增加 `workflow_dispatch` 触发器，保持原来的 `push` / `tag` 触发逻辑不变。随后同步更新 Docker Hub 部署说明，新增“如何手动触发一次”的操作步骤，确保用户在 GitHub UI 中可以看到对应按钮并知道如何验证结果。

**执行步骤：**
1. 创建本次任务目录并记录计划与任务清单。
2. 修改现有 GitHub Actions workflow，增加 `workflow_dispatch`。
3. 更新 Docker Hub 配置文档，补充手动触发说明。
4. 执行 YAML 与 diff 相关验证。
5. 提交并推送到 GitHub，确保远端可见“Run workflow”按钮。

**风险点：**
- 若 GitHub 仓库未启用 Actions，手动触发按钮仍不会显示。
- 若 Secrets / Variables 未配置正确，手动触发后仍会在推送镜像阶段失败。

**回滚思路：**
- 删除 `workflow_dispatch` 配置即可撤销手动触发入口。
- 若文档描述不准确，可单独回滚说明文档，不影响构建流程。
