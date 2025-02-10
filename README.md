# Beauty Rating Server

> 一个基于Go语言开发的美颜评分服务后端系统，提供图片上传、美颜分析、用户管理等功能。

## ✨ 功能特性

### 👤 用户系统
- 微信登录
- 用户信息管理(头像、昵称、性别等)
- 基于Token的用户认证
- 多角色支持(普通用户、管理员、专业用户)

### 📸 图片分析系统
- 图片上传与管理
- 美颜评分与分析
  - 五官评分
  - 气质评分
  - 妆容评分
  - 发型评分
- 分析结果管理(收藏/取消收藏)

## 🛠 技术栈

- Go 1.23.4
- Gin Web框架
- GORM
- MinIO 对象存储
- Redis 缓存
- gRPC
- MySQL

## 🚀 快速开始

### 环境要求

- Go 1.23.4+
- MySQL 8.0+
- Redis 6.0+
- MinIO

### 配置

```bash
service: beauty-server

redisAuth:
  server: localhost:6379
  db: 10
mysqlAuth:
  instance: localhost:3306
  database: your_database
  username: your_username
  password: your_password
minioAuth:
  endpoint: localhost:9000 
  accessKey: your_access_key 
  secretKey: your_secret_key
  bucket: your_bucket
```

## 📚 API文档

### 用户相关

| 接口 | 方法 | 路径 |
|------|------|------|
| 获取用户信息 | GET | `/api/v1/user/info` |
| 更新用户名 | PUT | `/api/v1/user/nickname/update` |
| 更新性别 | PUT | `/api/v1/user/gender/update` |
| 上传头像 | POST | `/api/v1/user/avatar/upload` |
| 获取头像 | GET | `/api/v1/user/avatar/:avatar_id` |

### 分析相关

| 接口 | 方法 | 路径 |
|------|------|------|
| 上传图片 | POST | `/api/v1/analysis/image/upload` |
| 获取图片 | GET | `/api/v1/analysis/image/:image_id` |
| 获取分析结果 | POST | `/api/v1/analysis` |
| 收藏分析结果 | POST | `/api/v1/analysis/favorite/:repord_id` |
| 取消收藏分析结果 | POST | `/api/v1/analysis/unfavorite/:repord_id` |
| 删除分析结果 | DELETE | `/api/v1/analysis/:repord_id` |

## 📄 许可证

本项目采用 MIT 许可证，详情请参见 [LICENSE](LICENSE) 文件。
