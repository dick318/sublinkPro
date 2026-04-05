# 任务清单：回滚订阅修复并接入 GitHub Actions 推送 Docker Hub

- [x] 任务 1：确认回滚边界与现状。
验收点：已确认最新提交 `808bbe2` 的改动文件范围，并区分哪些内容需要回滚、哪些内容需要保留为后续参考记录。

- [x] 任务 2：回滚订阅修复相关代码与当前部署语义。
验收点：订阅下载链路代码回到 `808bbe2` 之前的状态，且不会继续暴露“当前分支已包含该修复”的部署说明。

- [x] 任务 3：补充回滚记录，保留后续再次修复所需线索。
验收点：已新增文档记录“改过什么 / 为什么回滚 / 下次从哪里重新开始”。

- [x] 任务 4：接入每次提交自动推送 Docker Hub 的 GitHub Actions。
验收点：现有 workflow 已支持在 `push` 时构建并推送镜像，镜像仓库名与凭据改为通过仓库变量 / secrets 注入。

- [x] 任务 5：补充手动配置说明并完成本地验证。
验收点：已写明 GitHub Secrets / Variables 的手动配置步骤，并完成至少 1 次与 workflow 或相关文件直接关联的验证。

## 验证记录

- 已执行：`git restore --source=198fd67 --worktree -- "node/sub.go"`，将订阅下载链路恢复到实验性修复前的版本。
- 已执行：删除 [node/sub_test.go](/e:/pyWorkSpace/sublinkPro/node/sub_test.go) 与历史 Hugging Face 部署说明文件，撤下仅服务于该次修复的测试与部署说明。
- 已执行：新增 [rollback-notes.md](/e:/pyWorkSpace/sublinkPro/docs/tasks/2026-04-05-rollback-subscription-fix-and-setup-dockerhub-actions/rollback-notes.md)，保留实验性修复的背景、根因判断与后续重做入口。
- 已执行：更新 [build-release.yml](/e:/pyWorkSpace/sublinkPro/.github/workflows/build-release.yml)，把 Docker 镜像标签策略调整为“分支标签 + sha 标签 + main/latest”，并改为优先读取仓库变量 `DOCKERHUB_IMAGE`。
- 已执行：新增 [github-actions-dockerhub.md](/e:/pyWorkSpace/sublinkPro/docs/deployments/github-actions-dockerhub.md)，记录 GitHub Secrets / Variables 的手动配置步骤。
- 已执行：`git diff --check`，未发现空白符错误。
- 已执行：`K:\Program\go\bin\go.exe test ./node -count=1`，结果为 `?   	sublink/node	[no test files]`，说明回滚后的 `node` 包可正常完成基础测试扫描。
- 已执行：`D:\Anaconda3\envs\python3.10.9\python.exe -c "import yaml; yaml.safe_load(open('.github/workflows/build-release.yml', encoding='utf-8').read()); print('workflow-ok')"`，结果为 `workflow-ok`。

## 风险提示

- 回滚后如果旧 bug 仍然存在，需要基于记录重新修复。
- 自动推送 Docker Hub 依赖 GitHub 仓库配置正确，否则 workflow 会在登录或打标签阶段失败。
