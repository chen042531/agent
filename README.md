# Gemma 4 AI Agent

這是一個基於 **純 Python** 與 **Ollama API** 實作的終端機 AI Agent。
這個 Agent 能夠連結強大的 **Gemma 4** 模型，具備自主思考與呼叫工具（Tool Calling）的能力，幫助您搜尋網路、操作檔案、甚至執行系統指令。最重要的是，它內建了**隨時可調控的權限管理系統**，讓您可以安全地掌控 Agent 的行為。

## 🌟 核心特色

- **純 Python 實作**：無須學習龐大複雜的 Agent 框架（如 LangChain），輕鬆理解底層邏輯。
- **本地端執行**：透過 Ollama 在本地運行 Gemma 4，保障您的隱私與資料安全，且完全免費。
- **工具呼叫 (Tool Calling)**：內建多種實用工具，並可輕易擴充。
- **細粒度權限控制**：您可以在終端機中隨時開啟或關閉任意工具的權限，防止 Agent 執行危險操作。

---

## 🚀 快速上手 (Quick Start)

### 1. 系統需求
- **Python 3.8+**
- **Ollama**（請至 [Ollama 官方網站](https://ollama.com/) 下載並安裝）

### 2. 安裝與準備模型
首先，請打開終端機，下載並啟動 Gemma 4 模型。根據您的電腦記憶體大小，可以選擇 `4b` 或是 `26b` 版本：

```bash
# 如果您的 RAM 約 8GB ~ 16GB
ollama run gemma4:4b

# 如果您的 RAM 有 32GB 以上，可以嘗試更強大的版本
ollama run gemma4:26b
```
*(第一次執行時會自動下載模型，可能需要一些時間。確保看到可以對話的提示字元後，即可使用 `Ctrl+D` 退出)*

### 3. 安裝 Python 套件
進入本專案資料夾，安裝所需要的套件：

```bash
cd c:\Users\USER\Desktop\ai_agent
pip install -r requirements.txt
```

### 4. 啟動 Agent
安裝完畢後，執行以下指令啟動對話介面：

```bash
python main.py
```

---

## 🛠️ 可用指令 (Commands)

進入對話介面後，您可以像平常與 AI 聊天一樣用自然語言輸入問題。此外，本程式還支援以下特殊指令：

| 指令 | 說明 |
| :--- | :--- |
| `/help` | 顯示指令說明表。 |
| `/tools` | 檢視所有可用的工具，以及它們目前的權限狀態（開啟/關閉）。 |
| `/toggle <工具>` | 開啟或關閉某個工具的權限。例如：`/toggle execute_system_command`。 |
| `/clear` | 清除目前的對話歷史紀錄（當對話太長，模型開始有點混亂時可使用）。 |
| `/exit` | 關閉並退出程式。 |

---

## 🛡️ 安全性與權限管理

為了安全考量，當 Agent 試圖使用工具時，程式會先檢查該工具的權限狀態。
如果您不希望 Agent 亂改檔案或執行系統指令，請在對話前先輸入：
```
/toggle write_file
/toggle execute_system_command
```
當 Agent 嘗試呼叫被禁用的工具時，它會收到「Permission denied」的警告，並在對話中告知您它缺乏權限。

---
**現在，盡情享受與您的 Gemma 4 Agent 結伴工作的時光吧！**
