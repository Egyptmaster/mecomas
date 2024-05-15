package db

import (
	"cid.com/content-service/common/retry"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type accessor struct {
	conn     *gocql.ClusterConfig
	retry    gocql.RetryPolicy
	keyspace string
}

func newAccessor(conn *gocql.ClusterConfig, retry retry.Settings, keyspace string) accessor {
	return accessor{
		conn:     conn,
		retry:    &gocql.ExponentialBackoffRetryPolicy{NumRetries: int(retry.Retry), Max: retry.Delay, Min: retry.Delay / 2},
		keyspace: keyspace,
	}
}

func (t *accessor) createKeyspace() error {
	return t.withSession(func(session gocqlx.Session) error {
		return session.ExecStmt(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`, t.keyspace))
	})
}

func (t *accessor) withSession(fn func(gocqlx.Session) error) error {
	session, err := gocqlx.WrapSession(t.conn.CreateSession())
	if err != nil {
		return err
	}
	defer session.Close()

	err = fn(session)
	return err
}
