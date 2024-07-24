<div align="center">
    <img alt="ApiCat" width="350px" src="https://cdn.apicat.net/uploads/2d02ff2f6b19d3d6d3f134c1872484aa.png"/>
</div>

<p align="center">
  English |
  <a href="./README_CN.md">简体中文</a>
</p>

<p align="center">
    <a href="https://apicat.ai" target="_blank">
        <img alt="Static Badge" src="https://img.shields.io/badge/ai-apicat?logo=ai&logoColor=red&label=apicat&labelColor=4894FF&color=EAECF0">
    </a>
    <a href="https://discord.gg/6UFBGhNu" target="_blank">
        <img alt="Static Badge" src="https://img.shields.io/badge/chat-Discord-4E5AF0?logo=Discord">
    </a>
    <a href="https://github.com/apicat/apicat/blob/main/LICENSE">
        <img alt="Static Badge" src="https://img.shields.io/badge/license-MIT-green">
    </a>
</p>

ApiCat is an API documentation management tool that is fully compatible with the OpenAPI specification. With ApiCat, you can freely and efficiently manage your APIs. It integrates the capabilities of LLM, which not only helps you automatically generate API documentation and data models but also creates corresponding test cases based on the API content. Using ApiCat, you can quickly accomplish anything outside of coding, allowing you to focus your energy on the code itself.

## Using our Cloud Services

You can try out [ApiCat](https://apicat.ai) now. It provides all the capabilities of the self-deployed version.

## Local Installation

The easiest way to start the ApiCat is to run our docker-compose.yaml file. Before running the installation command, make sure that [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) are installed on your machine:

```bash
cd docker
docker compose up -d
```

After running, you can visit [http://localhost:8000](http://localhost:8000) on your browser to start using ApiCat.

If you need to customize the configuration, please refer to our [docker-compose.yaml](./docker/docker-compose.yaml) file and manually set the environment configuration. After making the changes, please run `docker-compose up -d` again.

## Community

If you have anything you would like to discuss with us, please join our community.

- [Discord](https://discord.gg/6UFBGhNu)

## License

[MIT](https://github.com/apicat/apicat/blob/main/LICENSE)