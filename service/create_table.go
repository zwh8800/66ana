package service

import (
	"time"

	"strings"

	"github.com/zwh8800/66ana/model"
)

func CreateTable(sql string) error {
	return dbConn.Exec(sql).Error
}

const (
	createDyDanmuSql = `CREATE TABLE dy_danmu_:day: (
  id         BIGSERIAL,
  cid        VARCHAR(255) NOT NULL DEFAULT '',
  user_id    BIGINT      NOT NULL DEFAULT 0,
  room_id    BIGINT      NOT NULL DEFAULT 0,
  content    VARCHAR(255) NOT NULL DEFAULT '',
  color      INTEGER      NOT NULL DEFAULT 0,
  client     INTEGER      NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ           DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE INDEX idx_dy_danmu_:day:_cid ON dy_danmu_:day:(cid);
CREATE INDEX idx_dy_danmu_:day:_user_id ON dy_danmu_:day:(user_id);
CREATE INDEX idx_dy_danmu_:day:_room_id ON dy_danmu_:day:(room_id);
CREATE INDEX idx_dy_danmu_:day:_color ON dy_danmu_:day:(color);
CREATE INDEX idx_dy_danmu_:day:_client ON dy_danmu_:day:(client);
CREATE INDEX idx_dy_danmu_:day:_created_at ON dy_danmu_:day:(created_at);
CREATE INDEX idx_dy_danmu_:day:_updated_at ON dy_danmu_:day:(updated_at);
CREATE INDEX idx_dy_danmu_:day:_deleted_at ON dy_danmu_:day:(deleted_at);`

	createDyGiftHistorySql = `CREATE TABLE dy_gift_history_:day: (
  id         BIGSERIAL,
  user_id    BIGINT      NOT NULL DEFAULT 0,
  room_id    BIGINT      NOT NULL DEFAULT 0,
  gift_id    BIGINT      NOT NULL DEFAULT 0,
  count      INTEGER      NOT NULL DEFAULT 1,
  hits       INTEGER      NOT NULL DEFAULT 1,
  gift_style VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ           DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE INDEX idx_dy_gift_history_:day:_user_id ON dy_gift_history_:day:(user_id);
CREATE INDEX idx_dy_gift_history_:day:_room_id ON dy_gift_history_:day:(room_id);
CREATE INDEX idx_dy_gift_history_:day:_gift_id ON dy_gift_history_:day:(gift_id);
CREATE INDEX idx_dy_gift_history_:day:_created_at ON dy_gift_history_:day:(created_at);
CREATE INDEX idx_dy_gift_history_:day:_updated_at ON dy_gift_history_:day:(updated_at);
CREATE INDEX idx_dy_gift_history_:day:_deleted_at ON dy_gift_history_:day:(deleted_at);`

	createDyDeserveSql = `CREATE TABLE dy_deserve_:day: (
  id         BIGSERIAL,
  user_id    BIGINT     NOT NULL DEFAULT 0,
  room_id    BIGINT     NOT NULL DEFAULT 0,
  level      INTEGER     NOT NULL DEFAULT 1,
  count      INTEGER     NOT NULL DEFAULT 1,
  hits       INTEGER     NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ          DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE INDEX idx_dy_deserve_:day:_user_id ON dy_deserve_:day:(user_id);
CREATE INDEX idx_dy_deserve_:day:_room_id ON dy_deserve_:day:(room_id);
CREATE INDEX idx_dy_deserve_:day:_level ON dy_deserve_:day:(level);
CREATE INDEX idx_dy_deserve_:day:_created_at ON dy_deserve_:day:(created_at);
CREATE INDEX idx_dy_deserve_:day:_updated_at ON dy_deserve_:day:(updated_at);
CREATE INDEX idx_dy_deserve_:day:_deleted_at ON dy_deserve_:day:(deleted_at);`
)

var (
	createSqlList = []string{
		createDyDanmuSql,
		createDyGiftHistorySql,
		createDyDeserveSql,
	}
	createTableName = []string{
		model.DyDanmuTableName,
		model.DyGiftHistoryTableName,
		model.DyDeserveTableName,
	}
)

func CreateFurtherTable(count int64) error {
	day := time.Now()
	for i := int64(0); i < count; i++ {
		for j, sql := range createSqlList {
			dayStr := day.Format("20060102")
			tableName := createTableName[j] + "_" + dayStr
			if !dbConn.HasTable(tableName) {
				sql = strings.Replace(sql, ":day:", dayStr, -1)
				if err := CreateTable(sql); err != nil {
					return err
				}
			}
		}
		day = day.Add(time.Duration(1 * 24 * int64(time.Hour)))
	}

	return nil
}
