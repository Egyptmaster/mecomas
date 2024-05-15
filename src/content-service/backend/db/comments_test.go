package db

import (
	"cid.com/content-service/common/retry"
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestComments(t *testing.T) {
	cluster := gocql.NewCluster("localhost:9042")
	a := &accessor{conn: cluster, keyspace: "tests"}
	require.NoError(t, a.createKeyspace())

	cluster = gocql.NewCluster("localhost:9042")
	//cluster.Keyspace = "tests"
	comments := NewCommentsTable(cluster, retry.Settings{Retry: 1, Delay: time.Millisecond * 150}, "tests")

	ctx := context.Background()

	t.Cleanup(func() {
		require.NoError(t, a.withSession(func(session gocqlx.Session) error {
			return session.ExecStmt(fmt.Sprintf("DROP KEYSPACE IF EXISTS %s", "tests"))
		}))
	})
	require.NoError(t, comments.Create())

	media := gocql.MustRandomUUID()
	data := make([]Comment, 10)

	for i := range len(data) {
		comment := Comment{
			CommentId:   gocql.MustRandomUUID(),
			UserId:      gocql.MustRandomUUID(),
			MediaItemId: media,
			Content:     fmt.Sprintf("%d. comment on this", i),
			Date:        time.Now().Add(time.Duration(-100+i) * time.Minute),
		}
		var err error
		data[i], err = comments.Insert(ctx, comment)
		require.NoError(t, err)
		require.NotNil(t, data[i].Date.Location())
		require.Equal(t, time.UTC, data[i].Date.Location())
		require.Equal(t, 0, data[i].Date.Nanosecond())

		c, err := comments.Get(ctx, comment.CommentId)
		require.NoError(t, err)
		require.NotNil(t, c)
		require.Equal(t, c.CommentId, data[i].CommentId)
		require.Equal(t, c.UserId, data[i].UserId)
		require.Equal(t, c.MediaItemId, data[i].MediaItemId)
		require.Equal(t, c.Content, data[i].Content)
		require.Equal(t, c.Date, data[i].Date)
	}

	cnt, err := comments.Count(ctx, media)
	require.NoError(t, err)
	require.Equal(t, len(data), cnt)

	list, anchor, err := comments.ByMediaItem(ctx, media, 4, nil)
	require.NoError(t, err)
	require.Len(t, list, 4)
	require.Equal(t, data[9].CommentId, list[0].CommentId)
	require.Equal(t, data[8].CommentId, list[1].CommentId)
	require.Equal(t, data[7].CommentId, list[2].CommentId)
	require.Equal(t, data[6].CommentId, list[3].CommentId)
	require.NotNil(t, anchor)

	list, anchor, err = comments.ByMediaItem(ctx, media, 4, anchor)
	require.NoError(t, err)
	require.Len(t, list, 4)
	require.Equal(t, data[5].CommentId, list[0].CommentId)
	require.Equal(t, data[4].CommentId, list[1].CommentId)
	require.Equal(t, data[3].CommentId, list[2].CommentId)
	require.Equal(t, data[2].CommentId, list[3].CommentId)
	require.NotNil(t, anchor)

	list, anchor, err = comments.ByMediaItem(ctx, media, 4, anchor)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, data[1].CommentId, list[0].CommentId)
	require.Equal(t, data[0].CommentId, list[1].CommentId)
	require.Len(t, anchor, 0)

	for i, d := range data {
		err = comments.Delete(ctx, d.CommentId)
		require.NoError(t, err)

		cnt, err = comments.Count(ctx, media)
		require.NoError(t, err)
		require.Equal(t, len(data)-(i+1), cnt)

		var c *Comment
		c, err = comments.Get(ctx, d.CommentId)
		require.NoError(t, err)
		require.Nil(t, c)
	}

	cnt, err = comments.Count(ctx, media)
	require.NoError(t, err)
	require.Equal(t, 0, cnt)

	list, anchor, err = comments.ByMediaItem(ctx, media, 4, nil)
	require.NoError(t, err)
	require.Len(t, list, 0)
	require.Len(t, anchor, 0)
}
