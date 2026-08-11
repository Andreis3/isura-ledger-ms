package decoder

import "github.com/bytedance/sonic"

// Frozen global instance, thread-safe and optimized for high performance
var jsonEngine = sonic.ConfigFastest