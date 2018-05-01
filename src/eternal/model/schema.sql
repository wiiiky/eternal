CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE DATABASE eternal WITH ENCODING=UTF8;


CREATE TABLE supported_country(
  code VARCHAR(8) NOT NULL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  sort INT NOT NULL
);
CREATE INDEX supported_country_sort ON supported_country(sort);

INSERT INTO supported_country(code,name,sort) VALUES('86','中国',0),('1','美国',1),('81','日本',2);

-- 密码加密类型
CREATE TYPE PasswordType AS enum('MD5','SHA1', 'SHA256');
-- 账号
CREATE TABLE account(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  country_code VARCHAR(8) NOT NULL DEFAULT '86',
  mobile VARCHAR(32) NOT NULL,
  salt VARCHAR(32) NOT NULL,
  passwd VARCHAR(256) NOT NULL,
  ptype PasswordType NOT NULL,
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(mobile)
);
CREATE INDEX account_mobile ON account(mobile);

-- Token
CREATE TABLE token(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  user_id UUID NOT NULL UNIQUE,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 性别
CREATE TYPE GenderType AS enum('MALE', 'FEMALE');
-- 用户信息
CREATE TABLE user_profile(
  user_id UUID PRIMARY KEY,
  name VARCHAR(32) NOT NULL DEFAULT '', -- '昵称'
  gender GenderType, -- '性别'
  birthday TIMESTAMP WITH TIME ZONE,
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 问题
CREATE TABLE question (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  title VARCHAR(32) NOT NULL, -- 问题标题 
  description TEXT NOT NULL, -- 问题详细描述
  user_id UUID NOT NULL UNIQUE,
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE answer (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  content TEXT NOT NULL, -- 回答
  user_id UUID NOT NULL,
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
