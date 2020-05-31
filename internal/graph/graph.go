package graph

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/dynastore"
	"github.com/wolfeidau/realworld-appsync-ddb/internal/graph/exec"
	"github.com/wolfeidau/realworld-appsync-ddb/internal/graph/model"
)

func NewResolvers(release dynastore.Partition) *Resolvers {
	return &Resolvers{release: release}
}

type Resolvers struct {
	release dynastore.Partition
}

func (r *Resolvers) Query() exec.QueryResolver {
	return (*QueryResolver)(r)
}

func (r *Resolvers) Mutation() exec.MutationResolver {
	return (*MutationResolver)(r)
}

type QueryResolver Resolvers

func (q *QueryResolver) Release(ctx context.Context, id string) (*model.Release, error) {
	log.Ctx(ctx).Info().Msg("Release")

	kv, err := q.release.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get record")
	}

	rel := new(model.Release)
	err = json.Unmarshal(kv.BytesValue(), rel)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode input")
	}

	return rel, nil
}

func (q *QueryResolver) ReleaseList(ctx context.Context, nextToken *string, limit *int) (*model.ReleasePage, error) {
	log.Ctx(ctx).Info().Msg("ReleaseList")

	var ropts []dynastore.ReadOption

	if nextToken != nil {
		ropts = append(ropts, dynastore.ReadWithStartKey(aws.StringValue(nextToken)))
	}

	if limit != nil {
		ropts = append(ropts, dynastore.ReadWithLimit(int64(*limit)))
	}

	kvpage, err := q.release.ListPage("", ropts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list page")
	}

	page := &model.ReleasePage{
		Items: []*model.Release{},
	}

	for _, kv := range kvpage.Keys {
		rel := new(model.Release)
		err = json.Unmarshal(kv.BytesValue(), rel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode json input")
		}

		page.Items = append(page.Items, rel)
	}

	if kvpage.LastKey != "" {
		page.NextToken = aws.String(kvpage.LastKey)
	}

	return page, nil
}

type MutationResolver Resolvers

func (m *MutationResolver) ReleaseCreate(ctx context.Context, input model.CreateReleaseInput) (*model.Release, error) {
	rel := new(model.Release)

	err := mapstructure.Decode(&input, &rel)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode input")
	}

	rel.ID = uuid.Must(uuid.NewUUID()).String()
	rel.Created = time.Now().Format(time.RFC3339Nano)

	log.Ctx(ctx).Info().Str("ID", rel.ID).Msg("release")

	data, err := json.Marshal(rel)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal attributes")
	}

	_, kv, err := m.release.AtomicPut(rel.ID, dynastore.WriteWithBytes(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to put record")
	}

	err = json.Unmarshal(kv.BytesValue(), rel)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal json")
	}

	return rel, nil
}

type SubscriptionResolver Resolvers

func (r *Resolvers) Subscription() exec.SubscriptionResolver {
	return (*SubscriptionResolver)(r)
}

func (s *SubscriptionResolver) NewRelease(ctx context.Context) (<-chan *model.Release, error) {
	panic("implement me")
}
