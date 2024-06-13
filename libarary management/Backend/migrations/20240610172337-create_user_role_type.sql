
-- +migrate Up
CREATE DOMAIN USER_ROLE AS VARCHAR
    CHECK (VALUE IN ('reader', 'admin','super')); 

-- +migrate Down

DROP DOMAIN IF EXISTS USER_ROLE ;
