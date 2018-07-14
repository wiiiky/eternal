-- DBTYPE: POSTGRESQL

CREATE DATABASE eternal WITH ENCODING='UTF8';
\c eternal;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- 客户端配置
CREATE TABLE "client" (
  id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v1mc(), -- 客户端ID
  name VARCHAR(32) NOT NULL,						-- 客户端名称
  token_max_age INTEGER NOT NULL,					-- 登录有效时长，单位秒
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO "client"(id, name, token_max_age)
  VALUES('137ff912-7106-11e8-9430-bb0f063260f6', 'Web', 3600 * 24 * 15),
        ('137ff913-7106-11e8-9430-bb9af99e7bb7', 'Android', 3600 * 24 * 365),
        ('137ff914-7106-11e8-9430-e3546a325cfb','IOS', 3600 * 24 * 365);

/* 注册帐号所支持的国家 */
CREATE TABLE "supported_country"(
  code VARCHAR(8) NOT NULL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  sort INT NOT NULL
);
CREATE INDEX supported_country_sort ON supported_country(sort);

INSERT INTO supported_country(code,name,sort) VALUES('86','中国',0),('1','美国',1),('81','日本',2);

-- 密码加密类型
CREATE TYPE PasswordType AS enum('MD5','SHA1', 'SHA256');
-- 账号
CREATE TABLE "account"(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  country_code VARCHAR(8) NOT NULL DEFAULT '86',
  phone_number VARCHAR(32) NOT NULL,
  salt VARCHAR(32) NOT NULL,
  passwd VARCHAR(256) NOT NULL,
  ptype PasswordType NOT NULL,
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(phone_number)
);
CREATE INDEX account__phone_number ON account(phone_number);

-- Token
CREATE TABLE "token"(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  user_id UUID NOT NULL,
  client_id UUID NOT NULL,
  etime TIMESTAMP WITH TIME ZONE,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  UNIQUE("user_id", "client_id")
);
CREATE INDEX token__user_id__client_id ON "token"(user_id, client_id);

-- 性别
CREATE TYPE GenderType AS enum('MALE', 'FEMALE' ,'');
-- 用户信息
CREATE TABLE "user_profile"(
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
CREATE TABLE "topic"(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  name VARCHAR(32) NOT NULL, -- 话题名
  icon VARCHAR(64) NOT NULL DEFAULT '', -- 图片ID
  introduction TEXT NOT NULL DEFAULT '', -- 描述
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX topic__name ON topic(name);

-- 问题
CREATE TABLE "question" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  title VARCHAR(64) NOT NULL, -- 问题标题
  content TEXT NOT NULL DEFAULT '', -- 问题详细描述
  user_id UUID NOT NULL,
  follow_count INTEGER NOT NULL DEFAULT 0, -- 关注数量
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/* 问题和话题的关联表 */
CREATE TABLE "question_topic" (
  question_id UUID NOT NULL, -- 问题ID
  topic_id UUID NOT NULL, -- 话题ID
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(question_id, topic_id)
);

/* 关注问题 */
CREATE TABLE "question_follow" (
  question_id UUID NOT NULL, -- 问题ID
  user_id UUID NOT NULL, -- 用户ID
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(question_id, user_id)
);

CREATE TABLE "answer" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
  content TEXT NOT NULL, -- 回答正文
  excerpt TEXT NOT NULL, -- 回答摘录
  question_id UUID NOT NULL, -- 问题ID
  user_id UUID NOT NULL,
  view_count INTEGER NOT NULL DEFAULT 0, -- view count 查看数 一个用户只会计一次
  upvote_count INTEGER NOT NULL DEFAULT 0, -- 点赞数
  downvote_count INTEGER NOT NULL DEFAULT 0, -- 不喜欢数
  utime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX answer__question_id ON answer(question_id);
CREATE INDEX answer__user_id ON answer(user_id);
CREATE INDEX answer__upvote_count ON answer(upvote_count);
CREATE INDEX answer__downvote_count ON answer(downvote_count);

CREATE TABLE "answer_upvote"(
  user_id UUID NOT NULL,
  answer_id UUID NOT NULL,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, answer_id)
);
CREATE INDEX answer_upvote__ctime ON answer_upvote(ctime);

CREATE TABLE "answer_downvote"(
  user_id UUID NOT NULL,
  answer_id UUID NOT NULL,
  ctime TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(user_id, answer_id)
);
CREATE INDEX answer_downvote__ctime ON answer_downvote(ctime);

CREATE TABLE "file"(
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
CREATE TABLE "hot_answer"(
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
