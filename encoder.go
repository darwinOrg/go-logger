package dglogger

import (
	"strings"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// customEncoder 包装原始 encoder，重写 EncodeEntry 来定制堆栈输出
type customEncoder struct {
	zapcore.Encoder
}

func (e *customEncoder) Clone() zapcore.Encoder {
	return &customEncoder{Encoder: e.Encoder.Clone()}
}

func (e *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 关键修改：直接操作 entry.Stack，堆栈信息存在这里
	if entry.Stack != "" {
		entry.Stack = formatStack(entry.Stack)
	}
	// 同时处理 fields 中的 stacktrace 字段（兼容不同版本的 zap）
	for i, f := range fields {
		if f.Key == "stacktrace" {
			if f.Type == zapcore.StringType && f.String != "" {
				fields[i].String = formatStack(f.String)
			}
		}
	}
	return e.Encoder.EncodeEntry(entry, fields)
}

// formatStack 将多行堆栈转为单行 <= 分隔，路径只保留最后 2 段
func formatStack(raw string) string {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var parts []string

	for i := 0; i < len(lines); i++ {
		funcLine := strings.TrimSpace(lines[i])
		if funcLine == "" || strings.HasPrefix(funcLine, "goroutine") {
			continue
		}

		// 下一行是文件路径
		if i+1 < len(lines) {
			fileLine := strings.TrimSpace(lines[i+1])
			if fileLine != "" && strings.Contains(fileLine, ".go:") {
				// 提取函数名：取最后一个 . 后面的部分
				funcName := funcLine
				if idx := strings.LastIndex(funcName, "."); idx >= 0 {
					funcName = funcName[idx+1:]
				}
				// 精简路径：只保留最后 2 段
				shortFile := shortenPath(fileLine)
				// 拼接：文件名(函数名):行号
				parts = append(parts, shortFile+"("+funcName+")")
				i++ // 跳过已处理的文件路径行
				continue
			}
		}

		// 兜底：没有配对的文件路径行，只输出函数名
		parts = append(parts, funcLine)
	}

	return strings.Join(parts, " <= ")
}

// shortenPath 保留路径最后 2 段
func shortenPath(path string) string {
	path = strings.TrimSpace(path)
	parts := strings.Split(path, "/")
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return path
}
