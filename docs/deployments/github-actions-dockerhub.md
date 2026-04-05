# GitHub Actions 自动构建并推送 Docker Hub

本文档说明如何让当前仓库在每次 `push` 后，自动通过 GitHub Actions 构建并推送 Docker 镜像到你自己的 Docker Hub 仓库。

## 1. 当前工作流行为

本仓库使用现有的 [build-release.yml](/e:/pyWorkSpace/sublinkPro/.github/workflows/build-release.yml) 完成前端构建、Go 二进制构建和 Docker 镜像推送。

调整后，Docker 镜像标签策略如下：

- 每次分支 `push`：推送 `<分支名>` 和 `sha-<短提交哈希>` 两个标签。
- `main` 分支 `push`：额外推送 `latest`。
- Git tag 推送：推送 `<tag>`、`sha-<短提交哈希>` 和 `latest`。

镜像名不再写死给 fork 使用，而是优先读取 GitHub 仓库变量 `DOCKERHUB_IMAGE`。如果不配置该变量，工作流仍会回退到仓库默认值 `zerodeng/sublink-pro`。

## 2. 先在 Docker Hub 准备仓库与令牌

1. 登录 Docker Hub。
2. 创建你自己的镜像仓库，例如 `liu4318/sublinkpro`。
3. 在 Docker Hub 账号安全页面创建一个 **Access Token**。
4. 保存好两项信息：
   - Docker Hub 用户名，例如 `liu4318`
   - 刚创建的 Access Token

建议把 Access Token 作为 GitHub Secret 使用，不要直接使用 Docker Hub 登录密码。

## 3. 在 GitHub 手动配置 Secrets

进入你的 GitHub 仓库：

`Settings -> Secrets and variables -> Actions`

在 `Secrets` 页签里新增以下两个仓库级 secret：

- `DOCKER_USERNAME`
  - 值填你的 Docker Hub 用户名。
- `DOCKER_PASSWORD`
  - 值填 Docker Hub Access Token。

## 4. 在 GitHub 手动配置 Variables

仍然在：

`Settings -> Secrets and variables -> Actions`

切换到 `Variables` 页签，新增：

- `DOCKERHUB_IMAGE`
  - 值填完整镜像名，例如 `liu4318/sublinkpro`

这样 workflow 在推送镜像时，就会自动把标签打到你的 Docker Hub 仓库下，而不是推默认镜像名。

## 5. 触发方式

当前 workflow 配置为：

- 可以在 GitHub Actions 页面手动触发一次。
- 任意分支 `push` 都会触发。
- 以 `v` 开头的 tag 推送也会触发。

如果你只想在 `main` 分支提交时触发，可以把 [build-release.yml](/e:/pyWorkSpace/sublinkPro/.github/workflows/build-release.yml) 顶部的：

```yaml
on:
  push:
    tags:
      - 'v*'
    branches:
      - '**'
```

改成：

```yaml
on:
  push:
    tags:
      - 'v*'
    branches:
      - main
```

## 6. 第一次验证

Secrets 和 Variables 配好后，执行一次正常提交并推送：

```bash
git add .
git commit -m "chore: 接入 Docker Hub 自动构建"
git push origin main
```

然后到 GitHub 仓库的 `Actions` 页面观察 `Build and Release Multi-Platform` 工作流。

成功后可以在 Docker Hub 看到类似标签：

- `main`
- `sha-abc1234`
- `latest`（仅 `main`）

## 7. 手动触发一次

当你不想再额外提交代码，但又想马上打一次镜像时，可以直接在 GitHub 页面手动触发：

1. 打开你的仓库 GitHub 页面。
2. 进入 `Actions`。
3. 左侧选择 `Build and Release Multi-Platform`。
4. 右上角点击 `Run workflow`。
5. `Use workflow from` 选择你要打包的分支，通常选 `main`。
6. 点击绿色的 `Run workflow` 按钮。
7. 等待任务开始后，点进这次运行记录。
8. 重点看 `build-docker` 任务里的这几步：
   - `Resolve Docker image tags`
   - `Log in to Docker Hub`
   - `Build and Push Multi-Arch Docker Image`

如果你从 `main` 手动触发，正常会推这几个标签：

- `main`
- `sha-<短哈希>`
- `latest`

如果你从其他分支手动触发，通常会推：

- `<分支名>`
- `sha-<短哈希>`

如果页面里看不到 `Run workflow` 按钮，通常有 3 个原因：

1. 新 workflow 还没推到 GitHub。
2. 仓库没启用 Actions。
3. 当前打开的不是带 `workflow_dispatch` 的那一条 workflow。

## 8. 常见失败点

- `unauthorized` / `requested access to the resource is denied`
  - 通常是 `DOCKER_USERNAME`、`DOCKER_PASSWORD` 或 `DOCKERHUB_IMAGE` 填错了。
- workflow 能登录但推送到错误仓库
  - 检查 `DOCKERHUB_IMAGE` 是否写成了完整形式：`用户名/仓库名`。
- 只想保留一个稳定标签
  - 可以保留 `main` + `latest`，删除 workflow 里的 `sha-<短哈希>` 标签逻辑。
- 分支名里带 `/`
  - 当前 workflow 已自动把 `/` 转成 `-`，避免 Docker tag 非法。

## 9. 与本次回滚的关系

这份自动化配置只影响构建和推送链路，不改变应用的数据库、业务逻辑或运行时配置。即使后续你要重新修订阅问题，也不需要重新设计 Docker Hub 自动推送方案。
