# revai
AI Code Review Assistant

[English version below](#english-version)

## Обзор

Это инструмент для полуавтоматического код-ревью с использованием AI (OpenAI/DeepSeek и др.). Анализирует git diff или отдельные файлы, генерирует отчет в Markdown.

## Установка

1. Убедитесь, что установлен Go (версия 1.20+)
2. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/revai.git
   cd revai
   ```
3. Соберите бинарник:
   ```bash
   go build -o revai main.go
   ```

## Конфигурация

Создайте файл `config/config.json`:
```json
{
  "role": "user",
  "prompt": "Проведи строгое код-ревью. Укажи на: 1) Баги 2) Стиль кода 3) Оптимизации",
  "ai": {
    "model": "gpt-4",
    "provider": "openai",
    "apiKey": "your_api_key",
    "apiBase": "https://api.openai.com/v1/chat/completions"
  },
  "excludeDirs": ["vendor", "node_modules"],
  "excludeFiles": ["*.md"]
}
```

## Использование

### Анализ git diff:
```bash
./revai -key YOUR_API_KEY
```
apikey можно указывать в конфиг файле либо при каждом запуске

### Анализ конкретного файла:
```bash
./revai -file path/to/file.go
```
если указать конкретный файл, будет проведено ревью именно его, иначе будет использован `git diff`

### Специфичный конфиг:
```bash
./revai -config path/to/config.json
```

## Пример вывода

Отчет сохраняется в файл формата `crYYYYMMDD_HHMM.md` с содержимым:
```markdown
## AI Code Review Report

### 1. Потенциальные баги
- В строке 42: возможен nil pointer dereference...
```

<a name="english-version"></a>
# AI Code Review Assistant

## Overview

Tool for semi-automated code reviews using AI (OpenAI/DeepSeek etc.). Analyzes git diff or specific files, generates Markdown reports.

## Installation

1. Ensure Go is installed (version 1.20+)
2. Clone repository:
   ```bash
   git clone https://github.com/yourusername/revai.git
   cd revai
   ```
3. Build binary:
   ```bash
   go build -o revai main.go
   ```

## Configuration

Create `config/config.json`:
```json
{
  "role": "user",
  "prompt": "Conduct strict code review. Highlight: 1) Bugs 2) Code style 3) Optimizations",
  "ai": {
    "model": "gpt-4",
    "provider": "openai",
    "apiKey": "your_api_key",
    "apiBase": "https://api.openai.com/v1/chat/completions"
  },
  "excludeDirs": ["vendor", "node_modules"],
  "excludeFiles": ["*.md"]
}
```

## Usage

### Analyze git diff:
```bash
./revai -key YOUR_API_KEY
```

### Analyze specific file:
```bash
./revai -file path/to/file.go
```

### Custom config:
```bash
./revai -config path/to/config.json
```

## Sample Output

Report is saved as `crYYYYMMDD_HHMM.md`:
```markdown
## AI Code Review Report

### 1. Potential Bugs
- Line 42: possible nil pointer dereference...
```
