package annals

import "time"

type CompilationMetadata struct {
	Duration time.Duration `json:"duration"`
	Args     []string      `json:"args"`
}
