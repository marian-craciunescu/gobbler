package store

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	sqliteMaxIdleConns = 2
	sqliteMaxOpenConns = 5
)

var writeTestFilename = "db_testfile"

var sqliteLogger = log.WithField("module", "kv-sqlite")

type SqliteKVStore struct {
	gormKVStore
	filename    string
	syncOnWrite bool
}

func NewSqliteKVStore(filename string, syncOnWrite bool) *SqliteKVStore {
	kvStore := &SqliteKVStore{}
	kvStore.filename = filename
	kvStore.syncOnWrite = syncOnWrite
	return kvStore
}

// Open opens the database file. If the directory does not exist, it will be created.
func (kvStore *SqliteKVStore) Open() error {
	directoryPath := filepath.Dir(kvStore.filename)
	if err := ensureWriteableDirectory(directoryPath); err != nil {
		sqliteLogger.WithFields(log.Fields{
			"dbFilename": kvStore.filename,
			"err":        err,
		}).Error("DB Directory is not writeable")
		return err
	}

	sqliteLogger.WithField("dbFilename", kvStore.filename).Info("Opening database")

	gormdb, err := gorm.Open("sqlite3", kvStore.filename)
	if err != nil {
		sqliteLogger.WithFields(log.Fields{
			"dbFilename": kvStore.filename,
			"err":        err,
		}).Error("Error opening database")
		return err
	}

	if err := gormdb.DB().Ping(); err != nil {
		sqliteLogger.WithFields(log.Fields{
			"dbFilename": kvStore.filename,
			"err":        err,
		}).Error("Error pinging database")
	} else {
		sqliteLogger.WithField("dbFilename", kvStore.filename).Info("Ping reply from database")
	}

	gormdb.LogMode(gormLogMode)
	gormdb.SingularTable(true)
	gormdb.DB().SetMaxIdleConns(sqliteMaxIdleConns)
	gormdb.DB().SetMaxOpenConns(sqliteMaxOpenConns)

	if err := gormdb.AutoMigrate(&kvEntry{}).Error; err != nil {
		sqliteLogger.WithField("err", err).Error("Error in schema migration")
		return err
	}

	sqliteLogger.Info("Ensured database schema")

	if !kvStore.syncOnWrite {
		sqliteLogger.Info("Setting db: PRAGMA synchronous = OFF")
		if err := gormdb.Exec("PRAGMA synchronous = OFF").Error; err != nil {
			sqliteLogger.WithField("err", err).Error("Error setting PRAGMA synchronous = OFF")
			return err
		}
	}
	kvStore.gormKVStore = gormKVStore{gormdb}
	return nil
}

func ensureWriteableDirectory(dir string) error {
	dirInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		dirInfo, err = os.Stat(dir)
	}
	if err != nil || !dirInfo.IsDir() {
		return fmt.Errorf("kv-sqlite: not a directory %v", dir)
	}
	writeTest := path.Join(dir, writeTestFilename)
	if err := ioutil.WriteFile(writeTest, []byte("writeTest"), 0644); err != nil {
		return err
	}
	if err := os.Remove(writeTest); err != nil {
		return err
	}
	return nil
}
