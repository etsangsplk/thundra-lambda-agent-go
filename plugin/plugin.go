package plugin

import (
	"context"
	"encoding/json"
	"sync"
)

type Plugin interface {
	BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup)
	AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string)
	OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string)
}
