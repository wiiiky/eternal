# encoding=utf8

from bs4 import BeautifulSoup
import requests

cookies = {
    'q_c1': '22842d8a164f40b9b812661d62af7e51|1525190972000|1498311384000',
    '_zap': '8d9061d7-5dc8-43ef-a050-6f84479d7b99',
    'd_c0': 'AAACfOZP_guPTjIr916gHo22XYeiWmBkjLE=|1498825157',
    '__DAYU_PP': 'fYVARJfYbRMEBnQ2vAvF3ddbe1600ada',
    'r_cap_id': 'OTBkMTgwMjg2NWI4NDVmNThjZTBlNmFmNWIxMTlhOTE=|1525190972|151514e1585af3e99fcd3a75fdf4de9c9ac78c01',
    'cap_id': 'NGUzYjY3MzBlNzBhNGJiM2JjZmQ4NzJlYmJmYTg3YjI=|1525190972|2ba612aff225ccdc5441f973929466ef43f8bbdb',
    'l_cap_id': 'Y2ZmOTc5ZTBmOTcxNGUyN2I5MGUwY2Y5NmYxODQ4MGM=|1525190972|ad5985f8ecb1adbe1d122838413dab5b611118fb',
    'capsion_ticket': '2|1:0|10:1525154848|14:capsion_ticket|44:MWY5MDYwNjQ4NTAyNDQ4ZDk1ZDNiMzA1M2M4ZmVjYTM=|d1817eaac2b29a69d76217d8fe92c36e3b25b10ede76d91e1c2ab6b6d2cd2e75',
    'aliyungf_tc': 'AQAAAKqesgS/YAwA3n/lZYVlpq40/bCF',
    '_xsrf': 'a13e6a31-302e-4ff1-ba77-2e72f6c6c9ef',
    'l_n_c': '1',
    'n_c': '1',
}

headers = {
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
    'Accept-Language': 'en,en-US;q=0.8,zh-CN;q=0.5,zh;q=0.3',
    'Cache-Control': 'max-age=0',
    'Connection': 'keep-alive',
    'DNT': '1',
    'Host': 'www.zhihu.com',
    'Upgrade-Insecure-Requests': '1',
    'User-Agent': 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:59.0) Gecko/20100101 Firefox/59.0',
}

params = (
    ('type', 'content'),
    ('q', '\u653F\u6CBB'),
)

response = requests.get('https://www.zhihu.com/search', headers=headers, params=params, cookies=cookies)
doc = BeautifulSoup(response.content, 'lxml')
for item in doc.findAll('div', {'class':'List-item'}):
    answerItem = item.find('div', {'class':'AnswerItem'})
    if not answerItem:
        continue
    titleItem = answerItem.find('h2', {'class':'ContentItem-title'})
    if not titleItem:
        continue
    print(titleItem.text)