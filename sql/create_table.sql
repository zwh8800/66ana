-- DROP DATABASE IF EXISTS "66ana";
-- CREATE DATABASE "66ana" ENCODING 'UTF8';
-- CREATE USER "66ana" PASSWORD '123456';

DROP TABLE IF EXISTS dy_user;
DROP TABLE IF EXISTS dy_room;
DROP TABLE IF EXISTS dy_user_room;
DROP TABLE IF EXISTS dy_cate;
DROP TABLE IF EXISTS dy_gift;
DROP TABLE IF EXISTS dy_danmu;
DROP TABLE IF EXISTS dy_gift_history;
DROP TABLE IF EXISTS dy_deserve;
DROP TABLE IF EXISTS dy_user_enter;
DROP TABLE IF EXISTS dy_super_danmu;

CREATE TABLE dy_user (
  id                     BIGSERIAL,
  uid                    BIGINT      NOT NULL,
  nickname               VARCHAR(120) NOT NULL DEFAULT '',
  level                  INTEGER      NOT NULL DEFAULT 0,
  strength               INTEGER      NOT NULL DEFAULT 0,
  gift_rank              INTEGER      NOT NULL DEFAULT 0,
  platform_privilege     INTEGER      NOT NULL DEFAULT 1,
  deserve_level          INTEGER      NOT NULL DEFAULT 0,
  deserve_count          INTEGER      NOT NULL DEFAULT 0,
  bdeserve_level         INTEGER      NOT NULL DEFAULT 0,
  first_appeared_room_id BIGINT      NOT NULL DEFAULT 0,
  last_appeared_room_id  BIGINT      NOT NULL DEFAULT 0,
  created_at             TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at             TIMESTAMPTZ  NOT NULL DEFAULT now(),
  deleted_at             TIMESTAMPTZ           DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX uni_idx_dy_user_uid ON dy_user(uid);

CREATE TABLE dy_room (
  id             BIGSERIAL,
  rid            BIGINT      NOT NULL,
  cate_id        BIGINT      NOT NULL DEFAULT 0,
  name           VARCHAR(255) NOT NULL DEFAULT '',
  status         INTEGER      NOT NULL DEFAULT 2,
  thumb          VARCHAR(255) NOT NULL DEFAULT '',
  avatar         VARCHAR(255) NOT NULL DEFAULT '',
  fans_count     INTEGER      NOT NULL DEFAULT 0,
  online_count   INTEGER      NOT NULL DEFAULT 0,
  owner_name     VARCHAR(120) NOT NULL DEFAULT '',
  weight         INTEGER      NOT NULL DEFAULT 0,
  last_live_time TIMESTAMPTZ  NOT NULL DEFAULT now(),
  created_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ  NOT NULL DEFAULT now(),
  deleted_at     TIMESTAMPTZ           DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX uni_idx_dy_room_rid ON dy_room(rid);

CREATE TABLE dy_user_room (
  id             BIGSERIAL,
  user_id        BIGINT     NOT NULL DEFAULT 0,
  room_id        BIGINT     NOT NULL DEFAULT 0,
  room_privilege INTEGER     NOT NULL DEFAULT 1,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at     TIMESTAMPTZ          DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE INDEX idx_dy_user_room_user_id ON dy_user_room(user_id);
CREATE INDEX idx_dy_user_room_room_id ON dy_user_room(room_id);

CREATE TABLE dy_cate (
  id         BIGSERIAL,
  cid        BIGINT      NOT NULL DEFAULT 0,
  game_name  VARCHAR(32)  NOT NULL DEFAULT '',
  short_name VARCHAR(16)  NOT NULL DEFAULT '',
  game_url   VARCHAR(120) NOT NULL DEFAULT '',
  game_src   VARCHAR(255) NOT NULL DEFAULT '',
  game_icon  VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ           DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE UNIQUE INDEX uni_idx_dy_cate_cid ON dy_cate(cid);

CREATE TABLE dy_gift (
  id           BIGSERIAL,
  room_id      BIGINT       NOT NULL DEFAULT 0,
  gid          BIGINT       NOT NULL DEFAULT 0,
  name         VARCHAR(120)  NOT NULL DEFAULT '',
  gift_type    INTEGER       NOT NULL DEFAULT 1,
  price        DECIMAL(7, 2) NOT NULL DEFAULT 0,
  contribution INTEGER       NOT NULL DEFAULT 0,
  intro        VARCHAR(120)  NOT NULL DEFAULT '',
  description  VARCHAR(120)  NOT NULL DEFAULT '',
  himg         VARCHAR(255)  NOT NULL DEFAULT '',
  mimg         VARCHAR(255)  NOT NULL DEFAULT '',
  created_at   TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ   NOT NULL DEFAULT now(),
  deleted_at   TIMESTAMPTZ            DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE INDEX idx_dy_gift_room_id ON dy_gift(room_id);
CREATE INDEX idx_dy_gift_gid ON dy_gift(gid);

CREATE TABLE dy_danmu (
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
CREATE INDEX idx_dy_danmu_cid ON dy_danmu(cid);
CREATE INDEX idx_dy_danmu_user_id ON dy_danmu(user_id);
CREATE INDEX idx_dy_danmu_room_id ON dy_danmu(room_id);

CREATE TABLE dy_gift_history (
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
CREATE INDEX idx_dy_gift_history_user_id ON dy_gift_history(user_id);
CREATE INDEX idx_dy_gift_history_room_id ON dy_gift_history(room_id);
CREATE INDEX idx_dy_gift_history_gift_id ON dy_gift_history(gift_id);

CREATE TABLE dy_deserve (
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
CREATE INDEX idx_dy_deserve_user_id ON dy_deserve(user_id);
CREATE INDEX idx_dy_deserve_room_id ON dy_deserve(room_id);

CREATE TABLE dy_user_enter (
  id         BIGSERIAL,
  user_id    BIGINT     NOT NULL DEFAULT 0,
  room_id    BIGINT     NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ          DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE INDEX idx_dy_user_enter_user_id ON dy_user_enter(user_id);
CREATE INDEX idx_dy_user_enter_room_id ON dy_user_enter(room_id);

CREATE TABLE dy_super_danmu (
  id           BIGSERIAL,
  sdid         VARCHAR(255) NOT NULL DEFAULT '',
  room_id      BIGINT      NOT NULL DEFAULT 0,
  jump_room_id BIGINT      NOT NULL DEFAULT 0,
  content      VARCHAR(255) NOT NULL DEFAULT '',
  created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
  deleted_at   TIMESTAMPTZ           DEFAULT NULL,

  PRIMARY KEY (id)
);
CREATE INDEX idx_dy_super_danmu_sdid ON dy_super_danmu(sdid);
CREATE INDEX idx_dy_super_danmu_room_id ON dy_super_danmu(room_id);
CREATE INDEX idx_dy_super_danmu_jump_room_id ON dy_super_danmu(jump_room_id);

GRANT ALL PRIVILEGES ON DATABASE "66ana" to "66ana";
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public to "66ana";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public to "66ana";
