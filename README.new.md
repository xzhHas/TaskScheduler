# TimeService
定时微服务

```yaml
version: '3.0'

networks:
  etcd-net: # 网 络
    driver: bridge # 桥 接 模 式

volumes:
  etcd1_data: # 挂 载 到 本 地 的 数 据 卷 名
    driver: local
  etcd2_data:
    driver: local
  etcd3_data:
    driver: local
### etcd 其 他 环 境 配 置 见 ： https://doczhcn.gitbook.io/etcd/index/index-1/configuration
services:
  etcd1:
    image: bitnami/etcd:latest # 镜 像
    container_name: etcd1 # 容 器 名  --name
    restart: always # 总 是 重 启
    networks:
      - etcd-net # 使 用 的 网 络  --network
    ports: # 端 口 映 射  -p
      - '20000:2379'
      - '20001:2380'
    environment: # 环 境 变 量  --env
      - ALLOW_NONE_AUTHENTICATION=yes # 允 许 不 用 密 码 登 录
      - ETCD_NAME=etcd1 # etcd 的 名 字
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd1:2380 # 列 出 这 个 成 员 的 伙 伴  URL 以 便 通 告 给 集 群 的 其 他 成 员
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380 # 用 于 监 听 伙 伴 通 讯 的 URL列 表
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 # 用 于 监 听 客 户 端 通 讯 的 URL列 表
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd1:2379 # 列 出 这 个 成 员 的 客 户 端 URL， 通 告 给 集 群 中 的 其 他 成 员
      - ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster # 在 启 动 期 间 用 于  etcd 集 群 的 初 始 化 集 群 记 号
      - ETCD_INITIAL_CLUSTER=etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380 # 为 启 动 初 始 化 集 群 配 置
      - ETCD_INITIAL_CLUSTER_STATE=new # 初 始 化 集 群 状 态
    volumes:
      - etcd1_data:/bitnami/etcd # 挂 载 的 数 据 卷

  etcd2:
    image: bitnami/etcd:latest
    container_name: etcd2
    restart: always
    networks:
      - etcd-net
    ports:
      - '20002:2379'
      - '20003:2380'
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_NAME=etcd2
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd2:2380
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd2:2379
      - ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster
      - ETCD_INITIAL_CLUSTER=etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380
      - ETCD_INITIAL_CLUSTER_STATE=new
    volumes:
      - etcd2_data:/bitnami/etcd

  etcd3:
    image: bitnami/etcd:latest
    container_name: etcd3
    restart: always
    networks:
      - etcd-net
    ports:
      - '20004:2379'
      - '20005:2380'
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_NAME=etcd3
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd3:2380
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd3:2379
      - ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster
      - ETCD_INITIAL_CLUSTER=etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380
      - ETCD_INITIAL_CLUSTER_STATE=new
    volumes:
      - etcd3_data:/bitnami/etcd
```

