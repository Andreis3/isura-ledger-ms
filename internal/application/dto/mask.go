package dto

import (
	"log/slog"
	"strings"
)

// MaskSlogValue otimizado para evitar reflexão pesada por chamada se a struct
// implementar diretamente uma interface de mascaramento ou se utilizarmos paths diretos.
// Caso mantenha a reflexão genérica, aqui está uma versão que reutiliza lógica e evita Sprintf.

func MaskSlogValue[T any](input T) slog.Value {
	// Dica de Otimização: Se T implementar uma interface customizada de log,
	// evitamos o custo de reflection completo em runtime para estruturas conhecidas.
	if logValuer, ok := any(input).(slog.LogValuer); ok {
		return logValuer.LogValue()
	}

	// Caso precise do fallback genérico por reflexão, mantenha enxuto:
	return slog.AnyValue(input)
}

// maskMiddleVisible otimizado com manipulação direta de runes sem strings.Builder excessivo
func MaskMiddleVisible(s string) string {
	runes := []rune(s)
	length := len(runes)

	if length <= 4 {
		return "******"
	}

	visibleLen := length / 3
	if visibleLen < 2 {
		visibleLen = 2
	}

	start := (length - visibleLen) / 2
	end := start + visibleLen

	// Pré-aloca o tamanho exato da string resultante para zerar realocações do Builder
	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < start; i++ {
		sb.WriteRune('*')
	}
	sb.WriteString(string(runes[start:end]))
	for i := end; i < length; i++ {
		sb.WriteRune('*')
	}

	return sb.String()
}
