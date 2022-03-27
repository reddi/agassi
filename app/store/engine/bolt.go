package engine

import (
	"encoding/json"
	"fmt"
	log "github.com/go-pkgz/lgr"
	"github.com/reddi/agassi/app/store"
	bolt "go.etcd.io/bbolt"
)

const (
	// top level buckets
	playersBucketName = "players"
	coachesBucketName = "coaches"
)

type BoltDB struct {
	db *bolt.DB
}

// NewBoltDB makes persistent boltdb-based store. For each site new boltdb file created
func NewBoltDB(options bolt.Options, dbFileName string) (*BoltDB, error) {
	log.Printf("[INFO] bolt store for file %+v, options %+v", dbFileName, options)
	db, err := bolt.Open(dbFileName, 0o600, &options) //nolint:gocritic //octalLiteral is OK as FileMode
	if err != nil {
		return nil, fmt.Errorf("failed to make boltdb for %s: %w", dbFileName, err)
	}
	result := &BoltDB{
		db: db,
	}
	// make top-level buckets
	topBuckets := []string{playersBucketName, coachesBucketName}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bktName := range topBuckets {
			if _, e := tx.CreateBucketIfNotExists([]byte(bktName)); e != nil {
				return fmt.Errorf("failed to create top level bucket %s: %w", bktName, e)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create top level bucket): %w", err)
	}

	log.Printf("[DEBUG] bolt store created for file %s", dbFileName)
	return result, nil
}

// save marshaled value to key for bucket. Should run in update tx
func (b *BoltDB) saveIfNotExists(bkt *bolt.Bucket, key []byte, value interface{}) (err error) {
	if value == nil {
		return fmt.Errorf("can't save nil value for %s", key)
	}
	if bkt.Get(key) != nil {
		return fmt.Errorf("key %s already in store", key)
	}
	jdata, jerr := json.Marshal(value)
	if jerr != nil {
		return fmt.Errorf("can't marshal comment: %w", jerr)
	}
	if err = bkt.Put(key, jdata); err != nil {
		return fmt.Errorf("failed to save key %s: %w", key, err)
	}
	return nil
}

func (b *BoltDB) createUser(bucketName string, userID string, user interface{}) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("no bucket %s", bucketName)
		}
		if err = b.saveIfNotExists(bucket, []byte(userID), user); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (b *BoltDB) CreatePlayer(player store.Player) (playerID string, err error) {
	err = b.createUser(playersBucketName, player.ID, player)
	return player.ID, err
}

func (b *BoltDB) CreateCoach(coach store.Coach) (coachID string, err error) {
	err = b.createUser(coachesBucketName, coach.ID, coach)
	return coach.ID, err
}

func (b *BoltDB) ListPlayers() (players []store.Player, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(playersBucketName))
		if bucket == nil {
			return fmt.Errorf("no bucket %s", playersBucketName)
		}
		return bucket.ForEach(func(k, v []byte) error {
			player := store.Player{}
			if e := json.Unmarshal(v, &player); e != nil {
				return fmt.Errorf("failed to unmarshal: %w", e)
			}
			players = append(players, player)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (b *BoltDB) ListCoaches() (coaches []store.Coach, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(coachesBucketName))
		if bucket == nil {
			return fmt.Errorf("no bucket %s", coachesBucketName)
		}
		return bucket.ForEach(func(k, v []byte) error {
			coach := store.Coach{}
			if e := json.Unmarshal(v, &coach); e != nil {
				return fmt.Errorf("failed to unmarshal: %w", e)
			}
			coaches = append(coaches, coach)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return coaches, nil
}

func (b *BoltDB) AddReview(playerID string, review store.Review) (err error) {
	return nil
}

func (b *BoltDB) ListReviews(playerID string) (reviews []store.Review, err error) {
	return make([]store.Review, 0), nil
}

func (b *BoltDB) Close() error {
	if err := b.db.Close(); err != nil {
		return fmt.Errorf("can't close store: %w", err)
	}
	return nil
}
