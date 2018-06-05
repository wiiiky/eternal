# encoding=utf8
import requests
import psycopg2
import uuid
import sys
import random
from topic import get_topic
from answer import get_answers

cookies = {
    'q_c1': '22842d8a164f40b9b812661d62af7e51|1526477104000|1498311384000',
    '_zap': '8d9061d7-5dc8-43ef-a050-6f84479d7b99',
    'd_c0': 'AAACfOZP_guPTjIr916gHo22XYeiWmBkjLE=|1498825157',
    '__DAYU_PP': 'fYVARJfYbRMEBnQ2vAvF3ddbe1600ada',
    'r_cap_id': 'ZjQ0MmVmYzRiNDY4NGZkMGJlYWJjODE2ZmU4MDI5YTU=|1525279707|8af37037016355bb9158afff5935337d94beac8b',
    'cap_id': 'ZDM4YzkzOGI0NzlhNGMwYzlmZTYyMGM3MzY3NjFmNTM=|1525279707|a72dad3d98c4f7094253b8ab8dc8239db497330c',
    'l_cap_id': 'MDBlNjc4MDdmNWU5NGFhOWI2NzliZGZmNzFhYjQ2OWE=|1525279708|fd2d57343dd765f148a280af0fcb7dca156d3d47',
    'capsion_ticket': '2|1:0|10:1525362246|14:capsion_ticket|44:NjVmMjVkOGExMzM1NDAxM2IxZGQ0NGY4MWI4NWQyM2M=|fe46915ad49df9378bfe47dc83f34ccd1b4edaba82f33efc83d87863c37978be',
    'z_c0': '2|1:0|10:1525362252|4:z_c0|92:Mi4xcmIwc0FBQUFBQUFBQUFKODVrXy1DeVlBQUFCZ0FsVk5USHpZV3dCazMxbnJQenJvbmlEdTdfWXFXcksyU3FSa2FB|0bdcd3ed9cdfe1c959dd9f56e4f47402b00390c863ecbb60b55b20b7d86ba33a',
    'aliyungf_tc': 'AQAAAKqesgS/YAwA3n/lZYVlpq40/bCF',
    '_xsrf': 'a13e6a31-302e-4ff1-ba77-2e72f6c6c9ef',
    'l_n_c': '1',
    'n_c': '1',
    'tgw_l7_route': '27a99ac9a31c20b25b182fd9e44378b8',
}

headers = {
    'Host': 'www.zhihu.com',
    'User-Agent': 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0',
    'Accept': 'application/json, text/plain, */*',
    'Accept-Language': 'en,en-US;q=0.8,zh-CN;q=0.5,zh;q=0.3',
    'Referer': 'https://www.zhihu.com/',
    'x-udid': 'AAACfOZP_guPTjIr916gHo22XYeiWmBkjLE=',
    'x-api-version': '3.0.53',
    'origin': 'https://www.zhihu.com',
    'DNT': '1',
    'Connection': 'keep-alive',
}


def get_feeds(page):
    params = (
        ('action_feed', 'True'),
        ('limit', '5'),
        ('session_token', '5bc73951ee8ebb9e6f29a39c39e10d05'),
        ('action', 'down'),
        ('after_id', str(page * 6)),
        ('desktop', 'true'),
    )
    response = requests.get('https://www.zhihu.com/api/v3/feed/topstory',
                            headers=headers, params=params, cookies=cookies)
    return response.json()['data']


db = psycopg2.connect('dbname=eternal user=postgres')
cur = db.cursor()


def new_user():
    while True:
        try:
            pk = str(uuid.uuid4())
            mobile = str(random.randint(18600000000, 18621578815))
            cur.execute(
                '''INSERT INTO account(id,country_code,mobile,salt,passwd,ptype) VALUES(%s,%s,%s,'','','MD5')''', (pk, '86', mobile))
            cur.execute(
                '''INSERT INTO user_profile(user_id,description,name) VALUES(%s,%s,%s)''', (pk, '测试用户的简介', '测试用户'))
            return pk
        except:
            pass


def save_topic(name, introduction):
    '''保存话题'''
    cur.execute('''SELECT id FROM topic WHERE name=%s''', (name,))
    r = cur.fetchone()
    if r:
        return r[0]
    pk = str(uuid.uuid4())
    cur.execute('''INSERT INTO topic(id, name,introduction) VALUES(%s,%s,%s)''',
                (pk, name, introduction))
    return pk


def save_question(title, description, topics):
    '''保存问题'''
    cur.execute('''SELECT id FROM question WHERE title=%s''', (title,))
    r = cur.fetchone()
    if r:
        return r[0], False
    userID = new_user()
    pk = str(uuid.uuid4())
    cur.execute('''INSERT INTO question(id,title,content,user_id) VALUES(%s,%s,%s,%s)''',
                (pk, title, description, userID))
    for t in topics:
        cur.execute(
            '''INSERT INTO question_topic(question_id,topic_id) VALUES(%s,%s)''', (pk, t))
    return pk, True


def save_answer(qid, content, excerpt):
    userID = new_user()
    pk = str(uuid.uuid4())
    cur.execute('''INSERT INTO answer(id,user_id,question_id,content,excerpt) VALUES(%s,%s,%s,%s,%s)''',
                (pk, userID, qid, content, excerpt))
    return pk


for i in range(10):
    data = get_feeds(i)
    for d in data:
        if 'target' not in d:
            continue
        if 'question' not in d['target']:
            continue
        question = d['target']['question']
        qid = question['id']
        title = question['title']
        excerpt = question['excerpt']
        topics = []
        actors = d['actors']
        for actor in actors:
            topic = get_topic(actor['id'])
            if not topic:
                continue
            pk = save_topic(topic['name'], topic['introduction'])
            topics.append(pk)
        pk, new = save_question(title, excerpt, topics)
        if not new:
            continue
        answers = get_answers(qid)
        for a in answers:
            save_answer(pk, a['content'], a['excerpt'])
        print(title)

# 添加热门回答
cur.execute('''
INSERT INTO hot_answer (question_id, answer_id, topic_id)
    SELECT DISTINCT question.id, answer.id, question_topic.topic_id
        FROM question
        INNER JOIN answer ON answer.question_id = question.id
        INNER JOIN question_topic ON question_topic.question_id = question.id
            WHERE answer.id NOT IN (SELECT answer_id FROM hot_answer)
''')

db.commit()
