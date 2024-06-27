
#  Suno API

# 好消息
我提供 SUNO API，不需要任何部署，不需要购买 SUNO 账户。更低的价格，更好的体验。
网址: https://api.bltcy.ai

## Disclaimer
- This project is only released on GitHub under the MIT license, free and open-source for learning purposes.

## Supported Features
- [x] Suno API 支持灵感模式、自定义模式、续写
- [x] 提供符合 OpenAI 格式的接口标准化服务，支持流式、非流输出内容
- [x] 填写账号信息程序自动维护与保活
- [x] 持久化保存任务记录，程序中断重启也能恢复任务
- [x] 支持自定义 OpenAI Chat 返回内容格式，基于 Go Template 语法
- [x] 适配 chat-next-web 等前端项目
- [x] 简化部署流程，支持 docker-compose、docker


## API 文档

http://localhost:8000/swagger/index.html

## Deployment

### Configuration
这些先从浏览器中获取，后期会自动保活。
![cookie](./docs/images/image1.png)

### Env Environment Variables
| 环境变量 | 说明                            | 默认值                        |
| --- |-------------------------------|----------------------------|
| SESSION_ID | 上图获取的 suno seesion_id         | 空                          |
| COOKIE | 上图获取的 suno cookie             | 空                          | 
| BASE_URL | SUNO 官方请求 API URL<br/>        | https://studio-api.suno.ai | 
| PROXY | Http 代理                       | 空                          | 
| SQL_DSN | Mysql DSN，为空时使用sqlite         |   空                         | 
| PORT | 开放端口                          | 8000                       | 
| ROTATE_LOGS | 日志是否按天轮转                      | 是                          | 
| LOG_DIR | 日志输出路径                        | ./logs                     | 
| DEBUG | 是否开启 Debug 日志                 | 否                          | 
| PPROF | 是否开启 Pprof 性能分析，开启后 8005 端口使用 | 否                          |
| CHAT_OPENAI_BASE | OpenAI API 接口地址               | https://api.openai.com     |
| CHAT_OPENAI_KEY | OpenAI API 密钥                 | sk-xxxxx                   |
| CHAT_OPENAI_MODEL | 默认模型                          | gpt-4o                     |
| CHAT_TIME_OUT | Chat 请求超时时间                   | 600 秒                      |
| CHAT_TEMPLATE_DIR | chat 模板读取路径                   | ./template                 |

### Docker Deployment
本教程提供如何使用特定的环境变量及端口映射来运行一个Docker容器的分步指导。为了本指南的目的，敏感信息如SQL名称、密码和IP地址将被替换为占位符。

```bash
# 使用 MySQL 的部署命令，在上面的基础上添加 `-e SQL_DSN="root:123456@tcp(localhost:3306)/sunoapi"`，请自行修改数据库连接参数。

docker run --name suno-api -d -p 8000:8000 \
-e SESSION_ID=xxxx \
-e COOKIE=xxxx  \
-e CHAT_OPENAI_BASE=https://api.openai.com  \
-e CHAT_OPENAI_KEY=sk-xxxxx  \
 sunoapigo/suno-api
```

docker-compose deployment
```bash
docker-compose pull && docker-compose up -d
```

docker-compose.yml
```bash
version: '3.2'

services:
  sunoapi:
    image: sunoapigo/suno-api:latest
    container_name: sunoapi
    restart: always
    ports:
      - "8000:8000"
    volumes:
      - ./logs:/logs
      - ./template:/template
    environment:
      - PORT=8000
      - SQL_DSN=root:123456@tcp(localhost:3306)/sunoapi
      - TZ=Asia/Shanghai
      - ROTATE_LOGS=false
      - PPROF=false
      - DEBUG=false
      - CHAT_TEMPLATE_DIR=./template
      - CHAT_OPENAI_MODEL=gpt-4o
      - CHAT_OPENAI_BASE=https://one-api.bltcy.top
      - CHAT_OPENAI_KEY=sk-
```


## 自定义 OpenAI Chat 返回内容格式
编辑 ./template 中的 suno.yaml
使用 Go template 语法

chat_stream_submit 流式时，任务提交成功时输出  
chat_stream_tick 流式时，任务每次进度查询时输出  
chat_resp 流式时，完成输出的格式  

## 参考
- Suno AI 官网: https://suno.com
- Suno-API: https://github.com/SunoAI-API/Suno-API


## License
MIT © [Suno API](./license)


## 给我买瓶可乐
![zanshangcode.jpg](./docs/images/zanshangcode.jpg)

此项目开源于GitHub ，基于MIT协议且免费，没有任何形式的付费行为！如果你觉得此项目对你有帮助，请帮我点个Star并转发扩散，在此感谢你！
