package common

import (
	"os"
	"time"
)

var Version = "v0.0.0"
var StartTime = time.Now().Unix() // unit: second

var PProfEnabled = os.Getenv("PPROF") == "true"
var DebugEnabled = os.Getenv("DEBUG") == "true"
var LogDir = GetOrDefaultString("LOG_DIR", "./logs")
var RotateLogs = os.Getenv("ROTATE_LOGS") == "true"

var Port = GetOrDefaultString("PORT", "8000")

var SQLitePath = GetOrDefaultString("SQLITE_PATH", "api.db?_busy_timeout=5000")

var Proxy = GetOrDefaultString("PROXY", "")
var TemplateDir = GetOrDefaultString("TEMPLATE_DIR", "./template")

var BaseUrl = GetOrDefaultString("BASE_URL", "https://studio-api.suno.ai")
var SessionID = GetOrDefaultString("SESSION_ID", "")
var COOKIE = GetOrDefaultString("COOKIE", "")
var SunoChatOpenaiModel = GetOrDefaultString("SUNO_CHAT_OPENAI_MODEL", "gpt-4o")
var SunoChatOpenaiApiBASE = GetOrDefaultString("SUNO_CHAT_OPENAI_BASE", "https://api.openai.com")
var SunoChatOpenaiApiKey = GetOrDefaultString("SUNO_CHAT_OPENAI_KEY", "")

var TimeOut = GetOrDefault("TIME_OUT", 600) // 任务超时时间
