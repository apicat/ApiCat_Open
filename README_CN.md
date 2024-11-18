<div align="center">
    <img alt="ApiCat" width="350px" src="https://cdn.apicat.ai/uploads/github-logo.png"/>
</div>

<p align="center">
  <a href="./README.md">English</a> |
  简体中文
</p>

<p align="center">
    <a href="https://apicat.ai" target="_blank">
        <img alt="Static Badge" src="https://img.shields.io/badge/ai-apicat?logo=ai&logoColor=red&label=apicat&labelColor=4894FF&color=EAECF0">
    </a>
    <a href="https://discord.gg/BdF8Cd3G" target="_blank">
        <img alt="Static Badge" src="https://img.shields.io/badge/chat-Discord-4E5AF0?logo=Discord">
    </a>
    <a href="https://github.com/apicat/apicat/blob/main/LICENSE">
        <img alt="GitHub License" src="https://img.shields.io/github/license/apicat/apicat">
    </a>
</p>

ApiCat 是一个 API 文档管理工具，完全兼容 OpenAPI 规范。在 ApiCat 上您可以自由高效的管理您的 API，它结合了 LLM 的能力，不仅可以帮您自动生成 API 文档以及数据模型，还可以根据 API 内容帮您生成相应的测试用例。使用 ApiCat，您可以快速完成代码之外的任何事情，让精力聚焦在代码本身。

## 使用云端服务

使用 [ApiCat](https://apicat.ai) 提供开源版本的所有功能。

## 本地部署

本地部署 ApiCat 最简单方法是运行我们的 [docker-compose.yaml](./docker/docker-compose.yaml) 文件。在运行安装命令之前，请确保您的机器上安装了 [Docker](https://docs.docker.com/get-docker/) 和 [Docker Compose](https://docs.docker.com/compose/install/)：

```bash
cd docker
docker compose up -d
```

运行后，可以在浏览器上访问 [http://localhost:8000](http://localhost:8000) 即可开始使用 ApiCat。

如果您需要自定义配置，请参考我们的 [docker-compose.yaml](./docker/docker-compose.yaml) 文件，并手动设置环境配置。更改后，请再次运行 `docker-compose up -d`。

## 交流

如果您有任何想和我们交流讨论的内容，欢迎加入我们的社区。

- [Discord](https://discord.gg/BdF8Cd3G)
- 微信群

![Wechat Group](https://cdn.apicat.ai/uploads/wechat-group-qrcode.png)

## License

[Apache-2.0](https://github.com/apicat/apicat/blob/main/LICENSE)