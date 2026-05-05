# Quick Start: Go Code Analyzer MCP for Claude

快速設置並開始使用 Go 代碼分析 MCP server。

## 📋 前提條件

- ✅ Go 1.21+ 已安裝
- ✅ Claude Desktop 已安裝
- ✅ 有 Go 代碼庫需要分析

## 🚀 5 分鐘快速設置

### 步驟 1: 確認構建

```bash
cd /Users/ycchen/Desktop/agent/mcp-server
ls -lh go-code-mcp
# 應該看到 6.1M 的可執行文件
```

### 步驟 2: 配置 Claude Desktop

**macOS**:
```bash
# 打開配置文件
open ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

**配置內容**:
```json
{
  "mcpServers": {
    "go-code-analyzer": {
      "command": "/Users/ycchen/Desktop/agent/mcp-server/go-code-mcp",
      "env": {
        "WORKSPACE_ROOT": "/Users/ycchen/Desktop/agent"
      }
    }
  }
}
```

**重要**: 修改路徑為你的實際路徑！

### 步驟 3: 重啟 Claude Desktop

完全退出並重新啟動 Claude Desktop。

### 步驟 4: 測試

在 Claude Desktop 中：

```
You: List all packages in this project

Claude: [自動調用 list_packages 工具]
```

## 🎯 核心功能

### 1. 探索式分析

```
You: "Help me understand this codebase"

Claude 的工作流程:
1. list_packages (level: summary)
   → 看到項目結構
2. analyze_package (internal/agent)
   → 理解核心組件
3. get_type_details (Agent)
   → 了解主要類型
```

### 2. 漸進式上下文

**Level 1 - 最小化** (只有名稱):
```json
{"level": "minimal"}
→ ["agent", "cli", "tools"]
```

**Level 2 - 摘要** (包含計數):
```json
{"level": "summary"}
→ [
  {"path": "agent", "type_count": 3, "function_count": 12},
  {"path": "tools", "type_count": 8, "function_count": 25}
]
```

**Level 3 - 詳細** (包含簽名):
```json
{"include": ["types", "functions"], "exported_only": true}
→ 完整的類型和函數簽名
```

**Level 4 - 完整** (包含源碼):
```json
{"include_body": true}
→ 函數實現代碼
```

### 3. 智能搜索

```
You: "Find all functions related to permissions"

Claude: [調用 search_code]
→ 找到 PermissionManager, Toggle, IsAllowed 等
```

### 4. 深入分析

```
You: "How does the Chat method work?"

Claude的工作流程:
1. search_code("Chat") → 找到位置
2. get_function_details(include_body: true) → 看實現
3. get_type_details("Agent") → 理解上下文
4. 解釋給你聽
```

## 📊 可用工具

| 工具 | 用途 | 參數 |
|------|------|------|
| `list_packages` | 列出所有包 | path, level |
| `analyze_package` | 深入分析包 | package_path, include, exported_only |
| `get_function_details` | 函數詳情 | file_path, function_name, include_body |
| `get_type_details` | 類型詳情 | file_path, type_name, include_methods |
| `search_code` | 搜索符號 | query, kind, max_results |
| `get_file_overview` | 文件概覽 | file_path, level |
| `find_dependencies` | 依賴分析 | target, direction |

## 💡 使用範例

### 範例 1: 新專案入門

```
You: "I'm new to this codebase, give me an overview"

Claude 會:
1. 列出所有包
2. 分析主要包
3. 展示核心類型
4. 解釋架構
```

### 範例 2: 添加新功能

```
You: "How do I add a new tool to the agent?"

Claude 會:
1. 搜索 Tool 接口
2. 獲取接口定義
3. 查看現有工具範例
4. 指導你實現
```

### 範例 3: 調試問題

```
You: "Why isn't my tool being called?"

Claude 會:
1. 搜索工具執行相關代碼
2. 查看 Chat 方法實現
3. 檢查權限系統
4. 找出問題
```

## 🎓 動態上下文調整

### Claude 如何智能使用 MCP:

**初始查詢** (寬泛):
```
"What's in this project?"
→ list_packages(level: "minimal")
```

**進一步了解**:
```
"Tell me more about the agent package"
→ analyze_package(include: ["types", "functions"], exported_only: true)
```

**深入細節**:
```
"Show me how Chat works"
→ get_function_details(include_body: true)
```

**查看完整上下文**:
```
"What are all the fields and methods of Agent?"
→ get_type_details(include_methods: true)
→ analyze_package(include: ["types", "methods", "fields"], exported_only: false)
```

## ⚙️ 高級配置

### 多專案配置

```json
{
  "mcpServers": {
    "agent-project": {
      "command": "/path/to/go-code-mcp",
      "env": {
        "WORKSPACE_ROOT": "/Users/ycchen/Desktop/agent"
      }
    },
    "another-project": {
      "command": "/path/to/go-code-mcp",
      "env": {
        "WORKSPACE_ROOT": "/path/to/another/project"
      }
    }
  }
}
```

### 啟用日誌

```json
{
  "mcpServers": {
    "go-code-analyzer": {
      "command": "/path/to/go-code-mcp",
      "env": {
        "WORKSPACE_ROOT": "/Users/ycchen/Desktop/agent",
        "MCP_LOG_LEVEL": "debug"
      }
    }
  }
}
```

查看日誌:
```bash
tail -f ~/Library/Logs/Claude/mcp-server-go-code-analyzer.log
```

## 🐛 故障排除

### 問題: 工具未出現

**檢查 1**: 配置文件語法
```bash
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json | jq .
```

**檢查 2**: 可執行文件
```bash
/path/to/go-code-mcp
# 應該啟動並等待輸入
```

**檢查 3**: Claude Desktop 日誌
```bash
cat ~/Library/Logs/Claude/mcp-server-go-code-analyzer.log
```

### 問題: Server 崩潰

**檢查**: WORKSPACE_ROOT 是否存在且包含 Go 代碼
```bash
ls -la $WORKSPACE_ROOT
```

### 問題: 回應緩慢

**解決方案**:
1. 使用 `exported_only: true` 過濾
2. 限制 `max_results`
3. 從 minimal level 開始

## ✨ 最佳實踐

### 1. 從寬到窄

```
✅ 好的方式:
list_packages → analyze_package → get_type_details

❌ 不好的方式:
直接要求分析整個代碼庫
```

### 2. 使用適當的詳細程度

```
✅ 探索階段: level: "minimal" 或 "summary"
✅ 理解階段: level: "detailed" 或 exported_only: true
✅ 實現階段: include_body: true
```

### 3. 利用搜索

```
✅ 先搜索找到位置
✅ 再用 get_function_details 查看詳情
```

## 🎉 開始使用

現在你已經準備好了！在 Claude Desktop 中試試這些:

```
1. "List all packages in this project"
2. "Analyze the internal/agent package"
3. "Show me the Agent struct and its methods"
4. "Find all functions related to 'Tool'"
5. "How does the permission system work?"
```

Claude 會智能地使用 MCP 工具來回答這些問題，並根據需要動態調整上下文深度！

---

**需要更多幫助？** 查看:
- `README.md` - 完整文檔
- `INTEGRATION_GUIDE.md` - 詳細集成指南
- [MCP 文檔](https://modelcontextprotocol.io/)
