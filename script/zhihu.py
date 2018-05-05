# encoding=utf8

from bs4 import BeautifulSoup
import requests
import json

cookies = {
    'q_c1': '22842d8a164f40b9b812661d62af7e51|1525190991000|1498311384000',
    '_zap': '8d9061d7-5dc8-43ef-a050-6f84479d7b99',
    'd_c0': 'AAACfOZP_guPTjIr916gHo22XYeiWmBkjLE=|1498825157',
    '__DAYU_PP': 'fYVARJfYbRMEBnQ2vAvF3ddbe1600ada',
    'r_cap_id': 'OTBkMTgwMjg2NWI4NDVmNThjZTBlNmFmNWIxMTlhOTE=|1525190972|151514e1585af3e99fcd3a75fdf4de9c9ac78c01',
    'cap_id': 'NGUzYjY3MzBlNzBhNGJiM2JjZmQ4NzJlYmJmYTg3YjI=|1525190972|2ba612aff225ccdc5441f973929466ef43f8bbdb',
    'l_cap_id': 'Y2ZmOTc5ZTBmOTcxNGUyN2I5MGUwY2Y5NmYxODQ4MGM=|1525190972|ad5985f8ecb1adbe1d122838413dab5b611118fb',
    'capsion_ticket': '2|1:0|10:1525265114|14:capsion_ticket|44:NThlYjY4MWYxYWI4NDYyMGFjMTY4OGE2NmE1NjUzMWU=|c33f1d049dbd7803fcc54db9268d01452e7bc4082f06e09afc3df06dfcdbd926',
    'aliyungf_tc': 'AQAAAKqesgS/YAwA3n/lZYVlpq40/bCF',
    '_xsrf': 'a13e6a31-302e-4ff1-ba77-2e72f6c6c9ef',
    'l_n_c': '1',
    'n_c': '1',
}

headers = {
    'Accept': 'application/json, text/plain, */*',
    'Accept-Language': 'en,en-US;q=0.8,zh-CN;q=0.5,zh;q=0.3',
    'authorization': 'oauth c3cef7c66a1843f8b3a9e6a1e3160e20',
    'Connection': 'keep-alive',
    'DNT': '1',
    'Host': 'www.zhihu.com',
    'origin': 'https://www.zhihu.com',
    'Referer': 'https://www.zhihu.com/search?q=%E6%94%BF%E6%B2%BB&type=content',
    'User-Agent': 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:59.0) Gecko/20100101 Firefox/59.0',
    'x-api-version': '3.0.91',
    'x-app-za': 'OS=Web',
    'x-udid': 'AAACfOZP_guPTjIr916gHo22XYeiWmBkjLE=',
}

params = (
    ('t', 'general'),
    ('q', '\u653F\u6CBB'),
    ('correction', '1'),
    ('offset', '0'),
    ('limit', '100'),
    ('search_hash_id', '96db19f0705ffcf6c932690166b08c3e'),
)

response = requests.get('https://www.zhihu.com/api/v4/search_v3', headers=headers, params=params, cookies=cookies)
data = response.json()

data = data['data']
for d in data:
    if 'object' in d and 'question' in d['object']:
        question = d['object']['question']
        name = question['name']
        pk = question['id']
        url = question['url']
        name = name.replace('<em>', '').replace('</em>','')
        response = requests.get('https://www.zhihu.com/question/' + pk, headers=headers, cookies=cookies)
        print(pk,name,url)
        print(response.content.decode('utf8'))
#print(data)


#doc = BeautifulSoup(response.content, 'lxml')
#for item in doc.findAll('div', {'class':'List-item'}):
#    answerItem = item.find('div', {'class':'AnswerItem'})
#    if not answerItem:
#        continue
#    titleItem = answerItem.find('h2', {'class':'ContentItem-title'})
#    if not titleItem:
#        continue
#    print(titleItem.text)
