package handlers

import (
	"context"

	"github.com/rs/zerolog/log"
)

type RequestProcessor struct {
}

func (rp *RequestProcessor) Handler(ctx context.Context, payload []byte) ([]byte, error) {
	log.Ctx(ctx).Info().Msg("handling request")

	return nil, nil
}
