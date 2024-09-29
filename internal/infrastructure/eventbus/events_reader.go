package eventbus

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sgatu/ezmail/internal/domain/models/events"
)

type RedisEventsReader struct {
	conn       *redis.Client
	stream     string
	name       string
	autocommit bool
}

func NewRedisEventsReader(conn *redis.Client, stream string, readerName string, autocommit bool) *RedisEventsReader {
	return &RedisEventsReader{
		conn:       conn,
		stream:     stream,
		name:       readerName,
		autocommit: autocommit,
	}
}

func (re *RedisEventsReader) getCommitKey() string {
	return "ci_" + re.stream + "_" + re.name
}

func (re *RedisEventsReader) getLastMessageReadId(ctx context.Context) string {
	result := re.conn.Get(ctx, re.getCommitKey())
	data, err := result.Result()
	if err != nil {
		return "0"
	}
	return data
}

func (re *RedisEventsReader) Read(ctx context.Context, limit int32) ([]events.EventWrapper, error) {
	lastId := re.getLastMessageReadId(ctx)
	result := re.conn.XRead(ctx, &redis.XReadArgs{
		Streams: []string{re.stream, lastId},
		Count:   int64(limit),
		Block:   0,
	})

	resultData, err := result.Result()
	if err != nil {
		return nil, err
	}
	eventsList := make([]events.EventWrapper, 0)
	if len(resultData) == 0 {
		return eventsList, nil
	}
	streamData := resultData[0]
	for _, msg := range streamData.Messages {
		eventData, ok := msg.Values["payload"]
		lastId = msg.ID
		if ok {
			typedEvent, err := events.RetrieveTypedEvent([]byte(eventData.(string)))
			if err != nil {
				return nil, err
			}
			eventsList = append(eventsList, events.EventWrapper{Event: typedEvent, Id: msg.ID})
		}
	}
	if re.autocommit {
		re.Commit(ctx, lastId)
	}
	return eventsList, nil
}

func (re *RedisEventsReader) Commit(ctx context.Context, commitInfo interface{}) error {
	result := re.conn.Set(ctx, re.getCommitKey(), commitInfo, 0)
	return result.Err()
}
