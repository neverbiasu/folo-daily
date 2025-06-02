#new # Folo订阅源日报自动总结与推送 产品设计文档（PRD）_v2

## 目录

1. [产品背景](#一产品背景)
2. [产品目标](#二产品目标)
3. [整体架构及可行性](#三整体架构及可行性)
    - 关键环节拆解
    - 关键能力支持
4. [功能需求](#四功能需求)
    - Folo端
    - GitHub Actions端
    - 日报格式
    - 管理与配置
5. [接口与集成细节](#五接口与集成细节)
    - Folo Webhook → GitHub Actions
    - Actions定时拉取与数据补全
    - 日报推送
    - 配置文件(config.yml)结构优化
6. [异常处理机制](#六异常处理机制)
7. [安全与配额风控](#七安全与配额风控)
8. [上线策略与版本规划](#八上线策略与版本规划)
9. [可行性与风险提示（优化）](#九可行性与风险提示优化)
10. [附录](#十附录)

---

## 一、产品背景

随着信息爆炸，用户常常订阅大量RSS/Atom源，但难以及时、高效地消化全部内容。市面上已有的RSS聚合工具多为内容展示，缺乏智能摘要、日报推送等自动化能力。  
Folo作为新一代可自部署RSS平台，内置AI摘要和Webhook能力。GitHub Actions作为云端Serverless自动化平台，具备定时任务、数据聚合与推送能力。  
本产品旨在**零服务器、零运维**的前提下，帮助用户实现RSS内容的智能聚合、自动总结、每日定时推送，提升信息获取效率。

---

## 二、产品目标

- 用户可自定义订阅源，通过Folo自动聚合内容
- 每日自动对新内容生成AI摘要，并聚合为日报
- 日报通过指定渠道（如邮件、GitHub Pages、Telegram等）自动推送给用户
- 全流程基于Folo和GitHub Actions，无需额外服务器
- 支持多用户、多渠道、个性化定制

---

## 三、整体架构及可行性

### 1. 关键环节拆解
- **Folo 侧**：负责订阅、AI 摘要、Webhook 推送。
- **GitHub Actions 侧**：负责定时聚合、格式生成、多渠道推送。

**分析结论**：两部分职责明确，互为解耦，可行性高。

### 2. 关键能力支持
- Folo：已支持AI摘要、Webhook、分组/过滤、多用户多feed。
- GitHub Actions：支持repository_dispatch、schedule、内容聚合、推送、Secrets安全管理。
- config.yml：支持灵活配置。

---

## 四、功能需求

### 4.1 Folo端

- 支持AI摘要（自动在新条目顶部生成）
- 支持Webhook推送（POST Json，支持自定义Header和Payload）
- 支持多用户/多feed配置
- 可选过滤、屏蔽、分组

### 4.2 GitHub Actions端

- 支持repository_dispatch事件：可被Folo Webhook实时触发
- 支持schedule定时任务：每天定时聚合内容（补拉/去重，防漏数据）
- 支持内容存储与去重（如以data/目录存储日报原文和原始feed数据）
- 支持多渠道推送（邮件、Pages、Bot等，优先支持稳定渠道）
- 支持内容聚合、分组、排序、分页
- 支持异常报警（如推送失败自动邮件报警）

### 4.3 日报格式

- 标题（自动带日期）
- 每条摘要含：订阅源、标题、AI摘要、原文链接、发布时间
- 支持多feed分组、按时间排序
- 支持HTML和Markdown两种输出
- 支持附带统计（今日新条数、总订阅源数等）

### 4.4 管理与配置

- 支持通过仓库配置文件（如config.yml）自定义推送渠道、订阅源分组、推送频率、内容筛选规则
- 支持多用户/多频道配置，层级分明，便于扩展和权限隔离
- 日报、原始feed数据归档和清理策略
- 支持推送日志审计与问题定位

---

## 五、接口与集成细节

### 5.1 Folo Webhook → GitHub Actions

- Webhook指向：GitHub repository_dispatch API  
  URL: `https://api.github.com/repos/<owner>/<repo>/dispatches`
- Header:  
  - `Content-Type: application/json`
  - `Authorization: token <PAT>`
- Body:  
  - `event_type`（如folo_new_entry）
  - `client_payload`（映射Folo的entry/feed结构，见下方示例）

```json
{
  "event_type": "folo_new_entry",
  "client_payload": {
    "entry": {...},
    "feed": {...}
  }
}
```

> 风险点：Folo Webhook需支持自定义Header和Body结构。如果Folo暂未支持，可通过Serverless中转（如Cloudflare Worker）做一次转换，但建议推动Folo原生适配。

### 5.2 Actions定时拉取与数据补全

- 每天定时拉取Folo API最近24小时内容，增量补全，防止Webhook漏数据。
- 采用唯一id去重，聚合后归档。

### 5.3 日报推送

- 邮件：推荐 [actions-send-mail](https://github.com/marketplace/actions/send-email)
- Pages：自动push静态日报到gh-pages分支
- Bot/社交平台：插件/可选支持，推荐用成熟第三方Actions
- 所有敏感配置均用GitHub Secrets管理

### 5.4 配置文件(config.yml)结构优化

- 支持用户-频道-feed三级嵌套
- 支持推送频率、分组、内容筛选、推送方式定制
- 自动校验与回退机制

```yaml
users:
  - id: alice
    channels:
      - type: email
        to: alice@example.com
        time: "08:00"
      - type: telegram
        bot_token: "xxx"
        chat_id: "yyy"
    feeds:
      - id: "tech"
        group: "科技"
        filter: "AI|大模型"
        enable: true
      - id: "finance"
        group: "财经"
        enable: true
format: "markdown"
```

---

## 六、异常处理机制

- Webhook推送失败/被拒绝：自动降级为定时拉取
- 日报生成失败：重试3次，仍失败邮件报警
- 推送失败：记录日志，支持重发，报警
- 配置文件错误：自动检测并邮件提醒
- 数据重复/丢失：唯一id校验，缺失时定时补全

---

## 七、安全与配额风控

- PAT/SMTP等敏感信息仅存GitHub Secrets
- 配置文件权限隔离，防止越权推送
- 强调用户自建私有仓库，规避公开泄露风险
- GitHub Actions免费额度监控，推送频率/并发量可控

---

## 八、上线策略与版本规划

1. **MVP**（单用户、单渠道、邮件/Pages推送、聚合Markdown）
2. 支持多用户/多渠道
3. 配置化分组/筛选/模板
4. 异常补偿、推送日志、报警
5. 插件化推送渠道，支持Bot/第三方平台
6. 社区模板/插件市场

---

## 九、可行性与风险提示（优化）

- **Folo端完全可行，建议推动Webhook直接适配GitHub API，减少中转。**
- **GitHub Actions侧如需长期存储，建议日报/原始数据本地归档。**
- **多用户多渠道时，配置结构应清晰、易维护，并增加校验与错误回退。**
- **推送Bot类渠道作为选配，主链路优先支持邮件/Pages。**
- **注意GitHub Actions免费额度，建议面向核心用户自建私有仓库。**

---

## 十、附录

- Folo官方文档：https://github.com/RSSNext/Folo.wiki.git
- GitHub Actions文档：https://docs.github.com/en/actions
- 邮件Action市场：https://github.com/marketplace/actions/send-email
- Pages部署文档：https://pages.github.com/

---

**最终目标**：以最少的运维和开发门槛，为RSS深度用户和团队提供高效、智能、低成本的信息流聚合与推送工具。  
本产品设计兼顾可扩展性、健壮性和安全性，适合大厂级别的工程与产品实施落地。