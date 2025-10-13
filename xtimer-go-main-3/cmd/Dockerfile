FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 将项目文件复制到容器中
COPY cmd/xtimer/xtimer /app/cmd/xtimer/xtimer
COPY configs/config.yaml /app/configs/config.yaml

# 赋予可执行权限
RUN chmod +x /app/cmd/xtimer/xtimer
# 暴露端口（如果 xtimer 使用了端口）
EXPOSE 9001 9002

# 设置入口点
ENTRYPOINT ["/app/cmd/xtimer/xtimer"]

# 传递命令行参数到 ENTRYPOINT（使用 shell form 以确保环境变量展开）
CMD ["-conf", "/app/configs/config.yaml"]
