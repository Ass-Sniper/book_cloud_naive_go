
# Introduction to Docker Compose

---

## Docker Compose 如何工作？

使用 Docker Compose 时，你需要编写一个 YAML 配置文件（称为 Compose 文件）来配置应用的各个服务，然后通过 Compose CLI（命令行界面）来创建并启动这些服务。

Compose 文件（通常命名为 `compose.yaml`）遵循 [Compose 规范（Compose Specification）](https://github.com/compose-spec/compose-spec)，用于描述多容器应用的结构。Docker Compose 是该规范的一个具体实现。

---

### 🧩 Compose 应用模型

#### 📄 Compose 文件

默认的 Compose 文件路径为 `compose.yaml`（推荐），或兼容旧版本的 `compose.yml`。此外，也支持 `docker-compose.yaml` 和 `docker-compose.yml` 作为历史兼容。如果多个文件同时存在，Compose 优先使用 `compose.yaml`。

为了提高可维护性和复用性，Compose 支持使用片段（fragments）和扩展（extensions）。

你还可以合并多个 Compose 文件来定义应用模型。多个 YAML 文件将根据你设定的顺序合并：

* 简单属性和映射会被顺序中优先的文件覆盖；
* 列表会通过追加合并；
* 相对路径将以第一个 Compose 文件的父目录为基础进行解析。

对于那些既可以是字符串也可以是复杂对象的配置项，Compose 会先将其扩展为统一形式再执行合并。

🔗 参见：[Working with multiple Compose files](https://docs.docker.com/compose/extends/)

此外，如果你想复用其他 Compose 文件，或将大型应用拆分为多个 Compose 文件模块，还可以使用 `include` 功能。这对于跨团队协作或共享配置特别有用。

---

### 🧰 CLI：命令行工具

Docker 提供了 `docker compose` 命令行工具，用于管理通过 Compose 定义的多容器应用。

你可以通过 CLI 来：

* 启动服务
* 停止服务
* 查看状态
* 管理日志
* 调试服务等

#### 关键命令

```bash
# 启动 compose.yaml 中定义的所有服务
docker compose up

# 停止并移除正在运行的服务及其资源
docker compose down

# 查看容器输出日志，便于排查问题
docker compose logs

# 列出所有服务及其当前状态
docker compose ps
```

完整命令参考文档见：[Compose CLI 命令文档](https://docs.docker.com/compose/reference/)

---

### 💡 示例说明

以下示例展示了 Compose 的核心概念。设想你有一个包含前端和后端的 Web 应用：

* 前端服务部署在 HTTPS 上，配置文件与证书由平台基础设施注入；
* 后端服务保存数据在一个持久化卷中；
* 两者通过私有网络通信，前端同时暴露端口 443 用于外部访问。

![Docker Compose 应用结构图](https://docs.docker.com/compose/images/compose-application.webp)

#### 🧱 Compose 应用模型示例

```yaml
services:
  frontend:
    image: example/webapp
    ports:
      - "443:8043"
    networks:
      - front-tier
      - back-tier
    configs:
      - httpd-config
    secrets:
      - server-certificate

  backend:
    image: example/database
    volumes:
      - db-data:/etc/data
    networks:
      - back-tier

volumes:
  db-data:
    driver: flocker
    driver_opts:
      size: "10GiB"

configs:
  httpd-config:
    external: true

secrets:
  server-certificate:
    external: true

networks:
  front-tier: {}
  back-tier: {}
```

`docker compose up` 会完成以下操作：

* 启动 `frontend` 和 `backend` 服务；
* 创建所需的网络与数据卷；
* 将配置和证书注入前端服务。

运行状态检查示例：

```bash
$ docker compose ps

NAME                IMAGE                COMMAND                 SERVICE   CREATED          STATUS          PORTS
frontend-1          example/webapp       "nginx ..."             frontend  2 minutes ago    Up 2 minutes    0.0.0.0:443->8043/tcp
backend-1           example/database     "docker-entrypoint..."  backend   2 minutes ago    Up 2 minutes
```

---

## 接下来你可以做什么？

* 试试 [Compose 快速上手](https://docs.docker.com/compose/gettingstarted/)
* 浏览 [示例应用](https://github.com/docker/awesome-compose)
* 熟悉 [Compose 规范](https://compose-spec.io/)

---

## 为什么使用 Docker Compose？

### Docker Compose 的主要优势

使用 Docker Compose 可以简化容器化应用的开发、部署和管理流程，主要优点包括：

#### ✅ 简化控制

Docker Compose 允许你通过一个 YAML 文件来定义和管理多容器应用。它简化了多个服务之间的编排与协调，使得应用环境的管理和复制变得更加容易。

#### ✅ 高效协作

Docker Compose 的配置文件易于共享，有助于开发人员、运维团队及其他相关人员之间的协作。这种共享配置带来了更顺畅的工作流程、更快的问题解决效率以及整体生产效率的提升。

#### ✅ 快速的应用开发

Compose 会缓存创建容器所用的配置。当你重启一个未被修改的服务时，Compose 会复用现有容器，而不是重新创建。通过容器重用，你可以非常快速地对环境进行更改。

#### ✅ 跨环境可移植性

Compose 文件中支持使用变量。你可以借助变量来为不同的环境或用户自定义 Compose 配置，从而实现环境之间的便捷迁移。

#### ✅ 丰富的社区与支持

Docker Compose 拥有一个活跃的开源社区，因此你可以获得大量的学习资源、教程和技术支持。社区的活跃也推动了 Compose 的持续改进，同时帮助用户高效排查问题。

---

### Docker Compose 的常见使用场景

Docker Compose 可以在多种场景中使用，以下是一些典型例子：

#### 🛠️ 开发环境

在开发过程中，能在隔离环境中运行和交互应用是非常重要的。Compose 命令行工具可以帮助你快速创建这样的环境。

Compose 文件用于记录和配置应用的所有服务依赖（如数据库、消息队列、缓存、Web 服务 API 等）。你只需一条命令：

```
docker compose up
```

就可以为每个依赖启动一个或多个容器。

通过这种方式，你可以将原本需要几页说明文档才能完成的「开发环境配置」，简化为一个 Compose 文件和几条命令。

---

#### 🧪 自动化测试环境

在持续部署（CD）或持续集成（CI）流程中，自动化测试是至关重要的一环。Compose 提供了一种简单的方法，可以为你的测试套件创建和销毁隔离的测试环境。

只需如下几条命令即可完成自动化测试流程：

```bash
docker compose up -d
./run_tests
docker compose down
```

这样你可以在测试前快速部署环境，测试结束后自动清理资源。

---

#### 🧩 单机部署

虽然 Compose 最初是为开发和测试场景设计的，但现在越来越多地引入了适用于生产环境的特性。

若想了解 Compose 在生产环境中的使用方法，请参考官方文档的：[Compose in production](https://docs.docker.com/compose/production/)。

---

## **Docker Compose 的历史与发展**

本页面提供了：

* Docker Compose CLI 的发展简史
* 组成 Compose v1 和 Compose v2 的主要版本和文件格式的清晰解释
* Compose V1 和 Compose v2 之间的主要区别

### **简介**

![Docker Compose CLI 版本管理](https://docs.docker.com/compose/images/v1-versus-v2.png)

上一张图展示了 Compose v1 和 Compose v2 之间的主要区别。当前支持的 Docker Compose CLI 版本是 Compose v2，它由 Compose 规范定义。

它还快速展示了文件格式、命令行语法和顶级元素的差异。以下部分将更详细地介绍这些内容。

### **Docker Compose CLI 版本管理**

Docker Compose 命令行二进制文件的第一个版本于 2014 年发布。它是用 Python 编写的，通过 `docker-compose` 调用。通常，Compose V1 项目在 `compose.yaml` 文件中包含一个顶级版本元素，版本号范围从 2.0 到 3.8，这些版本号指代特定的文件格式。

Docker Compose 命令行二进制文件的第二个版本于 2020 年发布，是用 Go 编写的，通过 `docker compose` 调用。Compose v2 会忽略 `compose.yaml` 文件中的版本顶级元素。

### **Compose 文件格式版本管理**

Docker Compose CLI 是由特定的文件格式定义的。

Compose V1 有三个主要的文件格式版本发布：

* Compose 文件格式 1，随 Compose 1.0.0 于 2014 年发布
* Compose 文件格式 2.x，随 Compose 1.6.0 于 2016 年发布
* Compose 文件格式 3.x，随 Compose 1.10.0 于 2017 年发布

Compose 文件格式 1 与后续所有格式有显著不同，因为它缺少顶级的 `services` 键。它的使用历史悠久，使用这种格式编写的文件无法在 Compose v2 上运行。

Compose 文件格式 2.x 和 3.x 非常相似，但后者引入了许多针对 Swarm 部署的新选项。

为了处理 Compose CLI 版本管理、Compose 文件格式版本管理和根据是否使用 Swarm 模式而产生的功能差异，文件格式 2.x 和 3.x 被合并到 Compose 规范中。

**Compose v2 使用 Compose 规范进行项目定义。** 与之前的文件格式不同，Compose 规范是滚动更新的，并使得版本顶级元素成为可选项。Compose v2 还利用了可选的规范 — 部署、开发和构建。

为了使迁移更容易，Compose v2 对 Compose 文件格式 2.x/3.x 和 Compose 规范之间已弃用或更改的某些元素提供了向后兼容性。


# How-tos
在技术文档（如 Docker 官方手册）中，**How-tos** 是 **“How to do something”** 的缩写，中文常译为 **“操作指南”** 或 **“实践教程”**。这类内容专注于 **具体问题的解决方法** 和 **分步骤的操作指导**，通常包含：

---

## 📖 核心含义
- **字面意义**：*“如何做某事”* 的复数形式，即一系列具体任务的实现指南。
- **文档定位**：介于理论概念（如 *Introduction*）与命令参考手册（如 *CLI Reference*）之间的**实践性内容**。
- **中文对应**：类似“使用技巧”“实战教程”“操作手册”等表述。

---

## 使用 Compose Watch

**要求：**
Docker Compose 2.22.0 或更高版本

`watch` 属性会在你编辑并保存代码时，**自动更新并预览**正在运行的 Compose 服务。对于许多项目来说，一旦 Compose 启动，这就实现了一个几乎**无需手动干预**的开发流程——服务会在你保存工作时自动更新自己。

watch 遵循以下文件路径规则：

* 所有路径相对于项目目录（ignore 文件模式除外）
* 目录会被递归监听
* **不支持**通配符（glob patterns）
* 遵循 `.dockerignore` 中的规则
* 可通过 `ignore` 选项定义额外忽略的路径（语法相同）
* 常见 IDE（如 Vim、Emacs、JetBrains 等）生成的临时/备份文件会被自动忽略
* `.git` 目录会被自动忽略

你**不需要**对 Compose 项目中的所有服务都启用 `watch`。在一些场景中，比如仅需要前端（如 JavaScript）支持自动更新即可。

> ⚠️ Compose Watch 仅适用于使用 `build` 属性构建本地源代码的服务。它不会追踪那些通过 `image` 属性使用预构建镜像的服务。

---

### Compose Watch 与绑定挂载（bind mount）的区别

Compose 支持将主机目录挂载到容器中。`watch` 并**不是替代**这项功能，而是作为一种**补充机制**，更适合在容器中进行开发。

比起 bind mount，`watch` 提供了**更精细的控制能力**，允许你忽略特定文件或整个目录。

例如，在 JavaScript 项目中忽略 `node_modules/` 有两个好处：

* **性能提升**：包含大量小文件的文件树在某些配置中会产生很高的 I/O 负载
* **跨平台兼容**：如果宿主机和容器的系统架构不同，则无法共享编译产物

例如在 Node.js 项目中，不建议同步 `node_modules/`，即便 JS 是解释型语言，npm 包中也可能包含跨平台不兼容的本地代码。

---

### 配置说明

`watch` 属性用于定义规则，基于本地文件的更改来自动更新服务。

每条规则必须指定：

* 一个路径（path）
* 一个操作类型（action）：可选值为：

  * `sync`：同步文件
  * `rebuild`：重新构建镜像并替换服务
  * `sync+restart`：同步后重启容器

> 这些规则适用于多种语言和框架，具体路径和行为按项目而定，但原理通用。

#### 前提条件

容器镜像中应包含以下常用可执行文件：

* `stat`
* `mkdir`
* `rmdir`

此外，容器的 `USER` 用户需要有**写入权限**，通常你可以使用 Dockerfile 中的 `COPY --chown` 指令来设置初始文件的所有者，例如：

```dockerfile
# 以非 root 用户运行
FROM node:18
RUN useradd -ms /bin/sh -u 1001 app
USER app

# 安装依赖
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install

# 拷贝源代码并设置所有者
COPY --chown=app:app . /app
```

---

#### 各种操作说明

##### 1. `sync`

* 作用：自动将宿主机上的文件变更同步到容器内对应路径
* 适用场景：支持“热重载”（Hot Reload）的框架，如前端开发

##### 2. `rebuild`

* 作用：自动使用 BuildKit 重新构建镜像，并替换正在运行的容器
* 等价于命令：`docker compose up --build <服务名>`
* 适合场景：编译型语言或需完整构建的文件变动（如 `package.json`）

##### 3. `sync+restart`

* 作用：同步文件后重启服务
* 适合场景：配置文件更改（如数据库或 nginx 配置），无需重建镜像，只需重启服务主进程

---

#### `path` 和 `target`

* `path`：监听的本地路径
* `target`：容器中的目标路径

例如：

```yaml
path: ./app/html
target: /app/html
```

当你修改了 `./app/html/index.html`，则文件会被同步到 `/app/html/index.html`。

---

#### `ignore` 忽略规则

`ignore` 的路径是**相对于当前 `watch` 动作中的 `path` 路径**，而不是项目根目录。

---

### 示例 1：Node.js 项目

项目结构：

```
myproject/
├── web/
│   ├── App.jsx
│   ├── index.js
│   └── node_modules/
├── Dockerfile
├── compose.yaml
└── package.json
```

Compose 配置：

```yaml
services:
  web:
    build: .
    command: npm start
    develop:
      watch:
        - action: sync
          path: ./web
          target: /src/web
          ignore:
            - node_modules/
        - action: rebuild
          path: package.json
```

运行 `docker compose up --watch` 后：

* 启动服务，运行 `npm start`
* 启用文件监听：当你编辑 `web/` 中的文件，Compose 会将其同步到容器内 `/src/web/`
* 忽略了 `node_modules/` 目录的变动
* 如果 `package.json` 发生更改，会触发重新构建镜像

---

### 示例 2：加上 `sync+restart`

```yaml
services:
  web:
    build: .
    command: npm start
    develop:
      watch:
        - action: sync
          path: ./web
          target: /app/web
          ignore:
            - node_modules/
        - action: sync+restart
          path: ./proxy/nginx.conf
          target: /etc/nginx/conf.d/default.conf

  backend:
    build:
      context: backend
      target: builder
```

这个配置展示了如何使用 `sync+restart` 来更新 nginx 配置并重启服务，同时前端代码使用 `sync` 实时热更新，后端使用构建目标分离。

---

### 如何使用 watch

1. 在 `compose.yaml` 中为一个或多个服务添加 `watch` 配置
2. 执行命令：

   ```bash
   docker compose up --watch
   ```
3. 使用你的编辑器编辑服务源代码

> 💡 如果你不想把应用日志和同步日志混合在一起，也可以使用专门的命令：
> `docker compose watch`

---

### 示例项目与反馈

你可以参考官方示例项目：

* [dockersamples/avatars](https://github.com/dockersamples/avatars)
* Docker 文档的本地开发配置示例

欢迎前往 [Compose Specification 仓库](https://github.com/compose-spec/compose-spec) 提交反馈或报告问题。

---

如需进一步配置说明，可查阅：[Compose Develop 规范（Compose Develop Specification）](https://github.com/compose-spec/compose-spec/blob/main/develop.md)




