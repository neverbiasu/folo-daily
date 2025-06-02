# Webhook处理逻辑

## 功能描述
- 接收来自Folo的Webhook请求。
- 将数据保存为Markdown格式，存储在`data/<日期>/webhook_<时间>.md`。

## 实现细节
- 文件路径：`handlers/webhook_handler.go`
- 请求方法：POST
- 数据结构：
  ```json
  {
    "entry": { ... },
    "feed": { ... }
  }
  ```
- 文件结构：
  ```markdown
  # Webhook Data

  ## Entry
  <entry内容>

  ## Feed
  <feed内容>
  ```

## 文件存储结构
- 所有数据存储在`data/<日期>/`目录中。
- 文件命名为`webhook_<时间>.md`。

---
