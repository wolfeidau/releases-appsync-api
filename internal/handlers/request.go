package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/realworld-appsync-ddb/internal/graph"
	"github.com/wolfeidau/realworld-appsync-ddb/internal/graph/model"
)

var (
	ErrFailedToProcess = errors.New("failed to process request")
)

type appsyncEvent struct {
	TypeName  string                 `json:"typeName"`
	Field     string                 `json:"field"`
	Headers   map[string]interface{} `json:"headers"`
	Identity  map[string]interface{} `json:"identity"`
	Arguments map[string]interface{} `json:"arguments"`
}

func NewRequestProcessor(resolvers *graph.Resolvers) *RequestProcessor {
	return &RequestProcessor{resolvers: resolvers}
}

type RequestProcessor struct {
	resolvers *graph.Resolvers
}

func (rp *RequestProcessor) Handler(ctx context.Context, payload []byte) ([]byte, error) {
	log.Ctx(ctx).Info().Msg("handling request")

	ae := new(appsyncEvent)

	err := json.Unmarshal(payload, ae)
	if err != nil {
		log.Ctx(ctx).Err(err).Stack().Msg("failed to process request")
		return nil, ErrFailedToProcess
	}

	action := fmt.Sprintf("%s.%s", ae.TypeName, ae.Field)

	log.Ctx(ctx).Info().Str("action", action).Msg("process resolver")

	switch action {
	case "Query.releaseList":
		input := new(struct {
			NextToken *string `mapstructure:"nextToken"`
			Limit     *int    `mapstructure:"limit"`
		})
		err := mapstructure.Decode(ae.Arguments, &input)
		if err != nil {
			log.Ctx(ctx).Err(err).Stack().Msg("failed to Decode input")
			return nil, ErrFailedToProcess
		}

		page, err := rp.resolvers.Query().ReleaseList(ctx, input.NextToken, input.Limit)
		if err != nil {
			log.Ctx(ctx).Err(err).Stack().Msg("failed to process ReleaseList")
			return nil, ErrFailedToProcess
		}

		return json.Marshal(page)
	case "Query.release":
		input := new(struct {
			ID string `mapstructure:"id"`
		})
		err := mapstructure.Decode(ae.Arguments, &input)
		if err != nil {
			log.Ctx(ctx).Err(err).Stack().Msg("failed to Decode input")
			return nil, ErrFailedToProcess
		}

		page, err := rp.resolvers.Query().Release(ctx, input.ID)
		if err != nil {
			log.Ctx(ctx).Err(err).Stack().Msg("failed to process ReleaseList")
			return nil, ErrFailedToProcess
		}

		return json.Marshal(page)
	case "Mutation.releaseCreate":
		var input model.CreateReleaseInput

		fields, ok := ae.Arguments["input"]
		if !ok {
			log.Ctx(ctx).Err(err).Stack().Msg("failed to locate input")
			return nil, ErrFailedToProcess
		}

		err := mapstructure.Decode(fields, &input)
		if err != nil {
			log.Ctx(ctx).Err(err).Stack().Msg("failed to Decode input")
			return nil, ErrFailedToProcess
		}

		rel, err := rp.resolvers.Mutation().ReleaseCreate(ctx, input)
		if err != nil {
			log.Ctx(ctx).Err(err).Stack().Msg("failed to process ReleaseCreate")
			return nil, ErrFailedToProcess
		}

		return json.Marshal(rel)
	}

	return []byte(`{}`), nil
}
