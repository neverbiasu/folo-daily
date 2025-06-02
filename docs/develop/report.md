# 日报生成逻辑

## 功能描述
- 从`data/<日期>/`目录中读取Webhook保存的Markdown文件。
- 生成Markdown格式的日报，存储为`data/<日期>/daily_report.md`。

## 实现细节
- 文件路径：`handlers/report_generator.go`
- 文件结构：
  ```markdown
  # Daily Report - YYYY-MM-DD

  ## webhook_<时间>.md
  <文件内容>
  ```

## 文件存储结构
- 所有数据存储在`data/<日期>/`目录中。
- 日报文件命名为`daily_report.md`。

---
