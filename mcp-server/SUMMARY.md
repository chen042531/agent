# Go Code Analyzer MCP Server - 實現總結

## 🎯 項目概述

成功創建了一個 **Model Context Protocol (MCP) server**，為 Claude 提供智能的 Go 代碼分析能力。

## ✨ 核心特性

### 1. 動態上下文調整

MCP server 能夠根據 LLM 的需求動態提供不同深度的上下文：

- **Minimal**: 僅符號名稱
- **Summary**: 符號名稱 + 計數統計  
- **Detailed**: 完整簽名 + 文檔
- **Full**: 包含源碼實現

### 2. 漸進式探索

支持從高層概覽逐步深入到具體實現：

```
Level 1: list_packages (minimal)
  ↓
Level 2: analyze_package (summary) 
  ↓
Level 3: get_type_details (detailed)
  ↓
Level 4: get_function_details (include_body: true)
```

### 3. 智能工具選擇

7 個專門設計的工具，讓 Claude 能夠:

- 📋 **list_packages** - 發現代碼結構
- 📦 **analyze_package** - 深入分析包
- 🔍 **search_code** - 快速查找符號
- 📄 **get_file_overview** - 文件概覽
- 🎯 **get_function_details** - 函數詳情
- 🏗️ **get_type_details** - 類型詳情
- 🔗 **find_dependencies** - 依賴分析

## 📊 技術架構

```
┌─────────────────┐
│  Claude Desktop │
│   或 Claude Code │
└────────┬────────┘
         │ MCP Protocol
         │ (stdio)
┌────────▼────────┐
│   MCP Server    │
│  (Go Binary)    │
├─────────────────┤
│ 7 Tools:        │
│ - list_packages │
│ - analyze_pkg   │
│ - search_code   │
│ - get_function  │
│ - get_type      │
│ - get_file      │
│ - find_deps     │
├─────────────────┤
│ 2 Resources:    │
│ - code://package│
│ - code://file   │
└────────┬────────┘
         │
┌────────▼────────┐
│ Go AST Analyzer │
│ (go/ast, go/    │
│  parser, go/    │
│  token)         │
└────────┬────────┘
         │
┌────────▼────────┐
│  Go Codebase    │
│  (Workspace)    │
└─────────────────┘
```

## 📁 項目文件

```
mcp-server/
├── go-code-mcp (6.1MB)          # 可執行文件
├── main.go (209 lines)          # MCP server 主程序
├── analyzer.go (838 lines)      # AST 分析邏輯
├── go.mod                       # Go 依賴
├── go.sum                       # 依賴校驗
├── README.md                    # 完整文檔
├── INTEGRATION_GUIDE.md         # 集成指南
├── QUICK_START.md               # 快速開始
├── SUMMARY.md                   # 本文檔
└── claude_desktop_config.json  # 配置範例
```

## 🚀 使用流程

### 配置 Claude Desktop

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

### Claude 的智能工作流

當用戶問: **"幫我理解這個代碼庫"**

Claude 會自動:

1. 🔍 `list_packages(level: "summary")` 
   → 看到項目有 agent, cli, ollama, tools 四個主要包

2. 📦 `analyze_package("internal/agent", include: ["types", "functions"])`
   → 發現 Agent, PermissionManager, ConversationManager 三個核心類型

3. 🎯 `get_type_details("Agent", include_methods: true)`
   → 了解 Agent 有哪些字段和方法

4. 💡 綜合分析並解釋給用戶

**關鍵**: 整個過程中，Claude 根據需要動態調整請求的詳細程度，避免一次性載入過多上下文。

## 🎓 動態上下文範例

### 場景 1: 探索階段

```
User: "What packages exist?"
→ list_packages(level: "minimal")
→ Response: ["agent", "cli", "ollama", "tools"]
```

只返回包名，極少 token 消耗。

### 場景 2: 理解階段

```
User: "Tell me about the tools package"
→ analyze_package(package_path: "internal/tools", 
                  include: ["types", "functions"],
                  exported_only: true)
→ Response: {
    types: [Tool, Registry, ...],
    functions: [NewRegistry, ...]
  }
```

返回結構化信息，適中的 token 消耗。

### 場景 3: 實現階段

```
User: "Show me how WebSearchTool works"
→ get_function_details(file_path: "internal/tools/web_search.go",
                        function_name: "WebSearchTool.Execute",
                        include_body: true)
→ Response: {
    signature: "...",
    parameters: [...],
    returns: [...],
    body: "完整源碼..."
  }
```

返回完整實現，高 token 消耗但必要。

## 🎯 與 LLM 的緊密互動

### 互動模式 1: 詢問-細化

```
User: "Find permission-related code"
Claude → search_code(query: "permission")
Result → 找到 10+ 個結果

User: "Show me the PermissionManager"
Claude → get_type_details("PermissionManager", include_methods: true)
Result → 完整的類型定義和方法

User: "How does Toggle work?"
Claude → get_function_details("Toggle", include_body: true)
Result → 完整實現
```

### 互動模式 2: 層次探索

```
User: "Explain the agent architecture"

Claude 的思考過程:
1. 先看包結構 → list_packages
2. 分析核心包 → analyze_package("internal/agent")
3. 理解主類型 → get_type_details("Agent")
4. 查看關鍵方法 → get_function_details("Chat", include_body: true)
5. 綜合解釋架構
```

每一步都基於前一步的結果，動態決定下一步需要什麼深度的信息。

### 互動模式 3: 問題驅動

```
User: "Why doesn't my tool get called?"

Claude的調試流程:
1. 搜索工具調用相關代碼 → search_code("Execute")
2. 查看工具註冊邏輯 → get_function_details("RegisterToolRegistry")
3. 檢查權限系統 → get_type_details("PermissionManager")
4. 查看調用點 → get_function_details("Chat", include_body: true)
5. 識別問題並解釋
```

## 💡 關鍵設計決策

### 1. 參數化深度控制

每個工具都支持控制返回信息的詳細程度:

- `level`: minimal, summary, detailed, full
- `include`: 選擇包含哪些部分
- `exported_only`: 過濾非導出符號
- `include_body`: 是否包含源碼

### 2. 標準化輸出

所有工具返回結構化 JSON，便於 Claude 解析和理解。

### 3. 上下文感知

工具設計考慮了典型的代碼探索工作流:

```
發現 → 概覽 → 理解 → 深入 → 實現
  ↓       ↓       ↓      ↓       ↓
list  analyze  search  type  function
                              (body)
```

## 🔮 使用場景

### 1. 代碼庫入門

新開發者使用 Claude + MCP 快速理解項目結構和設計。

### 2. 功能實現

開發時詢問 Claude 如何實現某功能，Claude 通過 MCP 分析現有代碼並給出建議。

### 3. 代碼審查

審查時讓 Claude 分析變更影響範圍。

### 4. 重構規劃

討論重構方案時，Claude 分析依賴關係和影響。

### 5. 調試協助

遇到問題時，Claude 分析相關代碼邏輯幫助定位。

## 📈 性能特點

- ✅ **快速解析**: 使用 Go 標準庫 `go/ast`
- ✅ **按需加載**: 只解析請求的文件/包
- ✅ **低開銷**: 二進制文件僅 6.1MB
- ✅ **並發安全**: 支持多個並發請求

## 🎉 總結

成功實現了一個智能的 MCP server，能夠:

✅ **動態調整上下文** - 根據 LLM 需求提供不同深度的信息  
✅ **漸進式探索** - 從概覽到細節的自然工作流  
✅ **緊密互動循環** - Claude 可以基於結果決定下一步查詢  
✅ **最小化 token 使用** - 只在需要時請求詳細信息  
✅ **高效準確** - 基於 Go AST 的精確分析  

這個 MCP server 將 Claude 從"需要完整代碼上下文"的助手，轉變為"能夠智能探索和理解代碼"的協作夥伴！

---

**文件位置**: `/Users/ycchen/Desktop/agent/mcp-server/`  
**可執行文件**: `go-code-mcp` (6.1MB)  
**狀態**: ✅ 完成並可用  
**創建日期**: 2026-05-05
