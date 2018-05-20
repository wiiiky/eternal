-- DBTYPE: POSTGRESQL

CREATE DATABASE eternal WITH ENCODING='UTF8';
\c eternal;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

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
CREATE INDEX account__mobile ON account(mobile);

-- Token
CREATE TABLE token(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  user_id UUID NOT NULL UNIQUE,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 性别
CREATE TYPE GenderType AS enum('MALE', 'FEMALE' ,'');
-- 用户信息
CREATE TABLE user_profile(
  user_id UUID PRIMARY KEY,
  name VARCHAR(32) NOT NULL DEFAULT '', -- '昵称'
  gender GenderType DEFAULT '', -- '性别'
  birthday TIMESTAMP WITH TIME ZONE,
  avatar VARCHAR(64) NOT NULL DEFAULT '', -- 头像的图片ID
  cover VARCHAR(64) NOT NULL DEFAULT '',
  description VARCHAR(256) NOT NULL DEFAULT '', -- 自我描述
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/* 话题 */
CREATE TABLE topic(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  name VARCHAR(32) NOT NULL, -- 话题名
  icon VARCHAR(64) NOT NULL DEFAULT '', -- 图片ID
  introduction TEXT NOT NULL DEFAULT '', -- 描述
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 问题
CREATE TABLE question (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  title VARCHAR(64) NOT NULL, -- 问题标题
  content TEXT NOT NULL, -- 问题详细描述
  user_id UUID NOT NULL,
  view_count INTEGER NOT NULL DEFAULT 0, -- view count 查看数 一个用户只会计一次
  answer_count INTEGER NOT NULL DEFAULT 0, -- answer count 回答数
  answer_index FLOAT NOT NULL DEFAULT 0, -- 回答指数，在特定时间内收获的回答数
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/* 问题和话题的关联表 */
CREATE TABLE question_topic(
  qid UUID NOT NULL, -- 问题ID
  tid UUID NOT NULL, -- 话题ID
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(qid, tid)
);

CREATE TABLE answer (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  content TEXT NOT NULL, -- 回答正文
  question_id UUID NOT NULL, -- 问题ID
  user_id UUID NOT NULL,
  view_count INTEGER NOT NULL DEFAULT 0, -- view count 查看数 一个用户只会计一次
  like_count INTEGER NOT NULL DEFAULT 0, -- like count 喜欢数
  dislike_count INTEGER NOT NULL DEFAULT 0, -- dislike count 不喜欢数
  like_index FLOAT NOT NULL DEFAULT 0,  -- 喜欢指数，在特定时间内收获的点赞数
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX answer__question_id ON answer(question_id);
CREATE INDEX answer__user_id ON answer(user_id);

CREATE TABLE answer_like(
  user_id UUID NOT NULL,
  answer_id UUID NOT NULL,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, answer_id)
);

CREATE TABLE answer_dislike(
  user_id UUID NOT NULL,
  answer_id UUID NOT NULL,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, answer_id)
);

CREATE TABLE file(
  id VARCHAR(64) NOT NULL,
  content_type VARCHAR(128) NOT NULL, -- 文件类型
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(id)
);

/* 
 * 热门回答
 * 在一定时间内点赞数达到一定程度的回答
 * 使用定时任务计算热门回答
 */
CREATE TABLE hot_answer(
	id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
	answer_id UUID NOT NULL, -- 回答ID
	question_id UUID NOT NULL, -- 所属问题
	topic_id UUID NOT NULL, -- 所属话题,
	ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX hot_answer__answer_id ON hot_answer(answer_id);
CREATE INDEX hot_answer__question_id ON hot_answer(question_id);
CREATE INDEX hot_answer__topic_id ON hot_answer(topic_id);
CREATE INDEX hot_answer__ctime ON hot_answer(ctime);
