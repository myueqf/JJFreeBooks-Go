# Changelog

## [0.4.0](https://github.com/MEMLTS/JJFreeBooks-Go/compare/v0.3.1...v0.4.0) (2025-08-29)


### ✨ 新功能

* 优化程序结构和输出信息 ([3dedb7a](https://github.com/MEMLTS/JJFreeBooks-Go/commit/3dedb7a0db1a6b02b479c69424d38b46b08fa031))

## [0.3.1](https://github.com/MEMLTS/JJFreeBooks-Go/compare/v0.3.0...v0.3.1) (2025-08-29)


### 📦️ 构建系统

* **release:** 优化构建流程并添加版本信息 ([e9c65e6](https://github.com/MEMLTS/JJFreeBooks-Go/commit/e9c65e6f5bd7842528d61afcd863142a6ffa250b))

## [0.3.0](https://github.com/MEMLTS/JJFreeBooks-Go/compare/v0.2.0...v0.3.0) (2025-08-29)


### ✨ 新功能

* 添加软件版本号 ([e99578f](https://github.com/MEMLTS/JJFreeBooks-Go/commit/e99578f84f8d7b7b4048731de83ac57e359bb8ac))


### 🐛 错误修复

* 移除主界面中的版本号显示 ([8729b6e](https://github.com/MEMLTS/JJFreeBooks-Go/commit/8729b6ef9b2de43d7e8c1aa63c41c27b6eee7b43))


### 📦️ 构建系统

* 优化构建过程并添加 LDFLAGS 参数 ([9e42720](https://github.com/MEMLTS/JJFreeBooks-Go/commit/9e4272016367105f6dfd3e1c27afb6cc200627af))

## [0.2.0](https://github.com/MEMLTS/JJFreeBooks-Go/compare/v0.1.0...v0.2.0) (2025-08-29)


### ✨ 新功能

* **crypto:** 重构 DES 加密相关函数并添加动态密钥解密支持 ([8a23991](https://github.com/MEMLTS/JJFreeBooks-Go/commit/8a23991a01ce748c0e1f5f0986f6646127054c52))
* **main:** 添加定时任务功能- 使用 robfig/cron 包实现定时任务 ([99d02e6](https://github.com/MEMLTS/JJFreeBooks-Go/commit/99d02e63040d7c6fa4f26f3b365861e03c9d098e))


### 📝 文档更新

* **README:** 更新配置说明- 将配置文件从 config.json 改为 config.yaml ([46fbd54](https://github.com/MEMLTS/JJFreeBooks-Go/commit/46fbd5498b987d68cb95655d6777e07efedf8d18))
* **README:** 添加项目说明和使用指南 ([13b450b](https://github.com/MEMLTS/JJFreeBooks-Go/commit/13b450ba1e8e57a0dc3ce19f44214d08d60868bc))


### ♻️ 代码重构

* **api:** 重构 API调用并添加错误处理- 重命名结构体和函数以提高代码可读性 ([1c45ae7](https://github.com/MEMLTS/JJFreeBooks-Go/commit/1c45ae7b8a3367e65bca059e2e178d7b95f7c970))
* **config:** 移除 config.yaml 并修改相关配置逻辑- 删除了 config.yaml 文件，使用空字符串作为 Token 默认值 ([bf9d522](https://github.com/MEMLTS/JJFreeBooks-Go/commit/bf9d52242fe529431c76ac0c91b9fb1235c0d68e))


### 🎡 持续集成

* 添加发布发行版的 GitHub Actions 工作流 ([492b5c5](https://github.com/MEMLTS/JJFreeBooks-Go/commit/492b5c5430d402e41a047735e81ed3f7d8d1a4c0))
