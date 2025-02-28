package bolt

import (
	"context"
	"encoding/json"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/influxdata/platform"
)

var (
	sessionBucket = []byte("sessionsv1")
)

var _ platform.SessionService = (*Client)(nil)

func (c *Client) initializeSessions(ctx context.Context, tx *bolt.Tx) error {
	if _, err := tx.CreateBucketIfNotExists([]byte(sessionBucket)); err != nil {
		return err
	}
	return nil
}

// FindSession retrieves the session found at the provided key.
func (c *Client) FindSession(ctx context.Context, key string) (*platform.Session, error) {
	op := getOp(platform.OpFindSession)
	var sess *platform.Session
	err := c.db.View(func(tx *bolt.Tx) error {
		s, err := c.findSession(ctx, tx, key)
		if err != nil {
			return err
		}

		sess = s
		return nil
	})

	if err != nil {
		return nil, &platform.Error{
			Err: err,
			Op:  op,
		}
	}

	if err := sess.Expired(); err != nil {
		// todo(leodido) > do we want to return session also if expired?
		return sess, &platform.Error{
			Err: err,
			Op:  op,
		}
	}
	return sess, nil
}

func (c *Client) findSession(ctx context.Context, tx *bolt.Tx, key string) (*platform.Session, *platform.Error) {
	v := tx.Bucket(sessionBucket).Get([]byte(key))
	if len(v) == 0 {
		return nil, &platform.Error{
			Code: platform.ENotFound,
			Msg:  platform.ErrSessionNotFound,
		}
	}

	s := &platform.Session{}
	if err := json.Unmarshal(v, s); err != nil {
		return nil, &platform.Error{
			Err: err,
		}
	}

	// TODO(desa): these values should be cached so it's not so expensive to lookup each time.
	f := platform.UserResourceMappingFilter{UserID: s.UserID}
	mappings, err := c.findUserResourceMappings(ctx, tx, f)
	if err != nil {
		return nil, &platform.Error{
			Err: err,
		}
	}

	ps := make([]platform.Permission, 0, len(mappings))
	for _, m := range mappings {
		ps = append(ps, m.ToPermissions()...)
	}
	s.Permissions = ps
	return s, nil
}

// PutSession puts the session at key.
func (c *Client) PutSession(ctx context.Context, s *platform.Session) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		if err := c.putSession(ctx, tx, s); err != nil {
			return err
		}
		return nil
	})
}

func (c *Client) putSession(ctx context.Context, tx *bolt.Tx, s *platform.Session) *platform.Error {
	v, err := json.Marshal(s)
	if err != nil {
		return &platform.Error{
			Err: err,
		}
	}
	if err := tx.Bucket(sessionBucket).Put([]byte(s.Key), v); err != nil {
		return &platform.Error{
			Err: err,
		}
	}
	return nil
}

// ExpireSession expires the session at the provided key.
func (c *Client) ExpireSession(ctx context.Context, key string) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		s, err := c.findSession(ctx, tx, key)
		if err != nil {
			return err
		}

		s.ExpiresAt = time.Now()

		if err := c.putSession(ctx, tx, s); err != nil {
			return err
		}
		return nil
	})
}

// CreateSession creates a session for a user with the users maximal privileges.
func (c *Client) CreateSession(ctx context.Context, user string) (*platform.Session, error) {
	var sess *platform.Session
	err := c.db.Update(func(tx *bolt.Tx) error {
		s, err := c.createSession(ctx, tx, user)
		if err != nil {
			return err
		}

		sess = s

		return nil
	})

	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (c *Client) createSession(ctx context.Context, tx *bolt.Tx, user string) (*platform.Session, *platform.Error) {
	u, pe := c.findUserByName(ctx, tx, user)
	if pe != nil {
		return nil, pe
	}

	s := &platform.Session{}
	s.ID = c.IDGenerator.ID()
	k, err := c.TokenGenerator.Token()
	if err != nil {
		return nil, &platform.Error{
			Err: err,
		}
	}
	s.Key = k
	s.UserID = u.ID
	s.CreatedAt = time.Now()
	// TODO(desa): make this configurable
	s.ExpiresAt = s.CreatedAt.Add(time.Hour)
	// TODO(desa): not totally sure what to do here. Possibly we should have a maximal privilege permission.
	s.Permissions = []platform.Permission{}

	if err := c.putSession(ctx, tx, s); err != nil {
		return nil, err
	}

	return s, nil
}
