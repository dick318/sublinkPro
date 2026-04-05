# 任务清单：为 Docker Hub 构建增加手动触发入口

- [x] 任务 1：确认现有 workflow 触发方式。
验收点：已确认当前构建流程只有 `push` / `tag` 触发，没有 `workflow_dispatch`。

- [x] 任务 2：为现有 workflow 增加手动触发入口。
验收点：`build-release.yml` 已支持在 GitHub Actions 页面手动运行。

- [x] 任务 3：补充手动触发说明文档。
验收点：已写明 Secrets / Variables 如何填写，以及如何在 GitHub 页面手动点击运行。

- [x] 任务 4：完成本次验证并推送。
验收点：已完成 YAML / diff 验证，并将本次改动推送到 GitHub。

## 验证记录

- 已执行：在 [build-release.yml](/e:/pyWorkSpace/sublinkPro/.github/workflows/build-release.yml) 中增加 `workflow_dispatch`，保留原有 `push` / `tag` 触发方式。
- 已执行：更新 [github-actions-dockerhub.md](/e:/pyWorkSpace/sublinkPro/docs/deployments/github-actions-dockerhub.md)，补充“如何手动触发一次”的页面操作说明。
- 已执行：`git diff --check -- ".github/workflows/build-release.yml" "docs/deployments/github-actions-dockerhub.md" "docs/tasks/2026-04-05-enable-manual-dockerhub-workflow-trigger/plan.md" "docs/tasks/2026-04-05-enable-manual-dockerhub-workflow-trigger/task.md"`。
- 已执行：`D:\Anaconda3\envs\python3.10.9\python.exe -c "import yaml; yaml.safe_load(open('.github/workflows/build-release.yml', encoding='utf-8').read()); print('workflow-ok')"`，结果为 `workflow-ok`。
