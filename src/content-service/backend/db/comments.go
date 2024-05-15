package db

import (
	"cid.com/content-service/common/retry"
	"context"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"time"
)

const (
	tableComments     = "comments"
	columnCommentId   = "comment_id"
	columnUserId      = "user_id"
	columnMediaItemId = "media_item_id"
	columnContent     = "content"
	columnDate        = "date"
)

type Comment struct {
	CommentId   gocql.UUID
	UserId      gocql.UUID
	MediaItemId gocql.UUID
	Content     string
	Date        time.Time
}

type CommentFilter struct {
	CommentId   *gocql.UUID
	MediaItemId *gocql.UUID
}

func (cf CommentFilter) Where() []qb.Cmp {
	cmp := make([]qb.Cmp, 0)
	if cf.CommentId != nil {
		cmp = append(cmp, qb.Eq(columnCommentId))
	}
	if cf.MediaItemId != nil {
		cmp = append(cmp, qb.Eq(columnMediaItemId))
	}
	return cmp
}

type CommentsTableAccessor struct {
	accessor
	table *table.Table
}

func NewCommentsTable(cfg *gocql.ClusterConfig, retry retry.Settings, keyspace string) CommentsTableAccessor {
	metaData := table.Metadata{
		Name: fmt.Sprintf("%s.%s", keyspace, tableComments),
		//Name:    tableComments,
		Columns: []string{columnCommentId, columnUserId, columnMediaItemId, columnContent, columnDate},
		PartKey: []string{columnMediaItemId},
		SortKey: []string{columnDate},
	}
	tbl := table.New(metaData)
	return CommentsTableAccessor{
		accessor: newAccessor(cfg, retry, keyspace),
		table:    tbl,
	}
}

func (ct *CommentsTableAccessor) Create() error {
	return ct.withSession(func(session gocqlx.Session) error {
		return session.ExecStmt(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.comments ( 
		comment_id uuid,
		media_item_id uuid,
		user_id uuid,
		"content" text,
		"date" timestamp,
		PRIMARY KEY(media_item_id, "date")
    )
    WITH CLUSTERING ORDER BY ("date" DESC)`, ct.keyspace))
	})
}

func (ct *CommentsTableAccessor) Insert(ctx context.Context, comment Comment) (Comment, error) {
	return comment, ct.withSession(func(session gocqlx.Session) error {
		// convert the time to UTC
		if comment.Date.Location() != time.UTC {
			comment.Date = comment.Date.UTC()
		}
		comment.Date = comment.Date.Round(time.Second)

		qry := ct.table.InsertQueryContext(ctx, session).
			BindStruct(comment).
			RetryPolicy(ct.retry)

		return qry.ExecRelease()
	})
}

func (ct *CommentsTableAccessor) Delete(ctx context.Context, id gocql.UUID) (err error) {
	item, err := ct.Get(ctx, id)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("not found")
	}

	return ct.withSession(func(session gocqlx.Session) error {
		qry := ct.table.DeleteQueryContext(ctx, session).
			BindStruct(item).
			RetryPolicy(ct.retry)
		return qry.ExecRelease()
	})
}

func (ct *CommentsTableAccessor) Get(ctx context.Context, id gocql.UUID) (*Comment, error) {
	cql := qb.Select(ct.table.Metadata().Name).
		Columns(ct.table.Metadata().Columns...).
		AllowFiltering().
		Where(qb.EqNamed(columnCommentId, "id")).
		Limit(1)

	comment := new(Comment)
	return comment, ct.withSession(func(session gocqlx.Session) error {
		qry := cql.QueryContext(ctx, session).
			BindMap(map[string]interface{}{"id": id}).
			RetryPolicy(ct.retry)

		err := qry.GetRelease(comment)
		if err == nil || !errors.Is(err, gocql.ErrNotFound) { // ignore not found error
			return err
		}
		comment = nil
		return nil
	})
}

func (ct *CommentsTableAccessor) ByMediaItem(ctx context.Context, mediaItemId gocql.UUID, pageSize int, anchor []byte) (list []Comment, nextAnchor []byte, err error) {
	return ct.filter(ctx, CommentFilter{MediaItemId: &mediaItemId}, pageSize, anchor)
}

func (ct *CommentsTableAccessor) Count(ctx context.Context, mediaItemId gocql.UUID) (count int, err error) {
	return count, ct.withSession(func(session gocqlx.Session) error {
		qb.Select(ct.table.Metadata().Name).
			Where(qb.EqNamed(columnMediaItemId, "id")).
			CountAll().QueryContext(ctx, session).
			BindMap(map[string]interface{}{"id": mediaItemId}).
			Iter().Scan(&count)
		return nil
	})
}

func (ct *CommentsTableAccessor) filter(ctx context.Context, filter CommentFilter, pageSize int, anchor []byte) (list []Comment, nextAnchor []byte, err error) {
	return list, nextAnchor, ct.withSession(func(session gocqlx.Session) error {
		qry := ct.table.SelectQueryContext(ctx, session).
			RetryPolicy(ct.retry).
			BindStruct(&filter).
			PageSize(pageSize).
			PageState(anchor)

		defer qry.Release()

		iter := qry.Iter()

		nextAnchor = iter.PageState()
		return iter.Select(&list)
	})
}
