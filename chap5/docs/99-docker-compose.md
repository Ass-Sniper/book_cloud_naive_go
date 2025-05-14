

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
