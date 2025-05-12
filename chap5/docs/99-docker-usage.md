
# **1. 拉取镜像**

## **1.1. 清理所有未使用的镜像和层**

```shell
kay@kay-vm:docker$
kay@kay-vm:docker$ docker system prune -a
WARNING! This will remove:
  - all stopped containers
  - all networks not used by at least one container
  - all images without at least one container associated to them
  - all build cache

Are you sure you want to continue? [y/N] y
Deleted Images:
untagged: ubuntu:20.04
untagged: swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu:20.04
untagged: swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu@sha256:a3465b987ac01e4f028c6659d77f13648db791d6026d89c35e45cee16de6c045
deleted: sha256:9df6d6105df2788299e5cbbf3dfd09aa6b3b5a10b784b214c35a035f76deb4ba
deleted: sha256:3ec3ded77c0ce89e931f92aed086b2a2c774a6fbd51617853decc8afa4e1087a

Total reclaimed space: 72.81MB
kay@kay-vm:docker$
```

## **1.2. 拉取镜像时看到 560c024910be: Already exists**

在拉取镜像时看到 `560c024910be: Already exists`，这表示 Docker 检测到该镜像的某些层已经存在于本地缓存中，因此无需重复下载。这是 Docker 的正常行为，**无需解决**，也不会影响镜像的功能。以下是详细解释和扩展建议：

---

### **为什么会显示 `Already exists`？**
1. **Docker 镜像的分层机制**  
   Docker 镜像由多个只读层（Layer）组成。每个层对应 Dockerfile 中的一条指令（如 `RUN`、`COPY`` 等）。  
   - 如果某个层已经存在于本地，Docker 会直接复用该层，避免重复下载。
   - 例如，你之前可能拉取过其他基于 `ubuntu:20.04` 的镜像（如其他应用镜像），其基础层与本镜像共享。

2. **华为云镜像仓库的镜像来源**  
   你拉取的镜像地址为：  
   ```bash
   swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu:20.04
   ```  
   该镜像实际是华为云仓库对 Docker Hub 官方镜像 `ubuntu:20.04` 的代理或缓存。  
   - 如果你之前拉取过 `ubuntu:20.04` 或其他基于它的镜像，基础层 `560c024910be` 可能已被缓存。

---

### **如何验证镜像完整性？**
如果担心镜像被篡改或下载不完整，可以检查镜像的 `Digest`：
1. 查看拉取日志中的 `Digest` 值：
   ```bash
   Digest: sha256:a3465b987ac01e4f028c6659d77f13648db791d6026d89c35e45cee16de6c045
   ```
2. 与 Docker Hub 官方镜像的 Digest 对比：
   ```bash
   docker pull ubuntu:20.04
   docker inspect ubuntu:20.04 | grep -A 2 "RepoDigests"
   ```
   - 如果 Digest 一致，说明镜像内容完全相同，只是通过不同仓库分发。

---

### **如何强制重新下载镜像？**
如果希望完全重新下载镜像（即使本地已有缓存），可以使用 `--no-cache` 或删除旧镜像：
1. **方法 1：强制拉取最新版本**  
   添加 `--no-cache` 参数（但实际对 `docker pull` 无效，正确方法见下文）：
   ```bash
   docker pull --no-cache swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu:20.04
   ```
   - **注意**：`--no-cache` 对 `docker pull` 无效，正确方法是先删除本地镜像再拉取：
     ```bash
     docker rmi swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu:20.04
     docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu:20.04
     ```

2. **方法 2：清理所有未使用的镜像和层**  
   ```bash
   docker system prune -a
   ```
   - 这会删除所有未被容器使用的镜像和缓存层，之后重新拉取镜像将完全从远程仓库下载。

---

### **为何华为云仓库的镜像与 Docker Hub 的层相同？**
华为云的 `swr` 仓库可能直接代理或同步了 Docker Hub 的官方镜像。  
- 此类镜像的层与 Docker Hub 完全一致，只是通过华为云仓库加速访问（尤其在中国大陆地区访问更稳定）。

---

### **总结**
1. **无需操作**：`Already exists` 是正常提示，表示 Docker 复用了本地已有的层，可节省时间和流量。
2. **强制更新**：若需彻底重新下载，先删除本地镜像再拉取。
3. **镜像一致性**：通过对比 `Digest` 确保镜像内容一致。


## **1.3. 拉取指定镜像**

[Docker 渡渡鸟镜像同步站](https://docker.aityp.com/)

根据搜索结果，**https://docker.aityp.com/** 是 **Docker 渡渡鸟镜像同步站**的官方网站，这是一个专注于为国内开发者提供 Docker 镜像加速和同步服务的平台。以下是关于该网站的详细信息：

---

### **1. 核心功能**
1. **镜像同步与查询**  
   - 用户可通过搜索功能查找 Docker 镜像（如 `nginx`、`python` 等）。  
   - 若镜像未同步至平台，提交请求后系统会在 **1 小时内完成同步**，并提供国内高速下载地址。

2. **镜像加速器**  
   - 支持通过修改 Docker 的 `daemon.json` 配置文件，将默认镜像源替换为 `https://docker.aityp.com`，从而提升拉取速度。

3. **API 接口支持**  
   - 提供公共 API 用于查询镜像同步状态、获取版本信息等。例如，通过 `https://docker.aityp.com/api/v1/image?search=nginx` 可查询 `nginx` 镜像的详细信息。

4. **镜像版本管理**  
   - 支持查看镜像的历史版本，帮助用户选择适合的版本进行部署。

---

### **2. 适用场景**
- **国内网络环境优化**：针对中国大陆地区访问 Docker Hub、Google Container Registry（GCR）等官方仓库速度慢的问题，提供稳定的镜像加速服务。
- **开发和部署效率提升**：通过快速同步镜像和配置加速器，减少镜像拉取时间，优化 CI/CD 流程。

---

### **3. 使用注意事项**
- **加速器配置建议**：虽然支持配置为 Docker 加速器，但官方建议优先使用其他推荐地址（如阿里云、网易等）以获得更佳性能。
- **免费服务限制**：镜像同步服务为免费提供，但可能存在资源限制，建议合理使用并支持平台运营。

---

### **4. 与其他国内镜像源的对比**  
国内常见的 Docker 镜像加速器还包括：
- **阿里云加速器**：需注册并获取专属加速地址。
- **网易加速器**：直接使用 `http://hub-mirror.c.163.com`。
- **USTC 镜像**：由中国科学技术大学维护，地址为 `https://docker.mirrors.ustc.edu.cn`。  
相比而言，渡渡鸟镜像同步站的特色在于 **快速同步非官方镜像**（如 GCR 等）和 **API 集成支持**。

---

### **5. 使用教程示例**  
若需配置加速器，可按以下步骤操作：
1. 编辑 Docker 配置文件 `/etc/docker/daemon.json`，添加：
   ```json
   {
     "registry-mirrors": ["https://docker.aityp.com"]
   }
   ```
2. 重启 Docker 服务：
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart docker
   ```
3. 验证配置是否生效：
   ```bash
   docker info | grep "Registry Mirrors"
   ```
   若输出中包含 `https://docker.aityp.com`，则配置成功。

---

### **总结**  
**docker.aityp.com** 是一个针对国内用户的 Docker 镜像加速与同步平台，尤其适合需要快速访问海外镜像的场景。其核心优势在于镜像同步的时效性和 API 功能的集成，但需注意合理使用以避免资源滥用。如需更多国内镜像源信息，可参考其他平台如阿里云、网易等。


### **镜像拉取日志**

```shell
kay@kay-vm:docker$
kay@kay-vm:docker$
kay@kay-vm:docker$ docker pull swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu:20.04
20.04: Pulling from ddn-k8s/docker.io/library/ubuntu
560c024910be: Pull complete
Digest: sha256:a3465b987ac01e4f028c6659d77f13648db791d6026d89c35e45cee16de6c045
Status: Downloaded newer image for swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu:20.04
swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu:20.04
kay@kay-vm:docker$
kay@kay-vm:docker$
kay@kay-vm:docker$ docker tag  swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu:20.04  docker.io/library/ubuntu:20.04
kay@kay-vm:docker$
kay@kay-vm:docker$
kay@kay-vm:docker$
kay@kay-vm:docker$ docker images
REPOSITORY                                                          TAG       IMAGE ID       CREATED        SIZE
ubuntu                                                              20.04     9df6d6105df2   9 months ago   72.8MB
swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/ubuntu   20.04     9df6d6105df2   9 months ago   72.8MB
kay@kay-vm:docker$
kay@kay-vm:docker$
```


# **2. 重启docker**

```shell
sudo systemctl restart docker
```


