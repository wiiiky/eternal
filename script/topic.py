import requests

cookies = {
    'q_c1': '22842d8a164f40b9b812661d62af7e51|1526477109000|1498311384000',
    '_zap': '8d9061d7-5dc8-43ef-a050-6f84479d7b99',
    'd_c0': 'AAACfOZP_guPTjIr916gHo22XYeiWmBkjLE=|1498825157',
    'r_cap_id': 'ZjQ0MmVmYzRiNDY4NGZkMGJlYWJjODE2ZmU4MDI5YTU=|1525279707|8af37037016355bb9158afff5935337d94beac8b',
    'cap_id': 'ZDM4YzkzOGI0NzlhNGMwYzlmZTYyMGM3MzY3NjFmNTM=|1525279707|a72dad3d98c4f7094253b8ab8dc8239db497330c',
    'l_cap_id': 'MDBlNjc4MDdmNWU5NGFhOWI2NzliZGZmNzFhYjQ2OWE=|1525279708|fd2d57343dd765f148a280af0fcb7dca156d3d47',
    'capsion_ticket': '2|1:0|10:1525362246|14:capsion_ticket|44:NjVmMjVkOGExMzM1NDAxM2IxZGQ0NGY4MWI4NWQyM2M=|fe46915ad49df9378bfe47dc83f34ccd1b4edaba82f33efc83d87863c37978be',
    'z_c0': '2|1:0|10:1525362252|4:z_c0|92:Mi4xcmIwc0FBQUFBQUFBQUFKODVrXy1DeVlBQUFCZ0FsVk5USHpZV3dCazMxbnJQenJvbmlEdTdfWXFXcksyU3FSa2FB|0bdcd3ed9cdfe1c959dd9f56e4f47402b00390c863ecbb60b55b20b7d86ba33a',
    '_xsrf': 'a13e6a31-302e-4ff1-ba77-2e72f6c6c9ef',
    'l_n_c': '1',
    'n_c': '1',
    'aliyungf_tc': 'AQAAAPtRlnRJPQwA3n/lZV8NZ0qBXe/O',
}

headers = {
    'Host': 'api.zhihu.com',
    'User-Agent': 'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:60.0) Gecko/20100101 Firefox/60.0',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
    'Accept-Language': 'en,en-US;q=0.8,zh-CN;q=0.5,zh;q=0.3',
    'DNT': '1',
    'Connection': 'keep-alive',
    'Upgrade-Insecure-Requests': '1',
    'Cache-Control': 'max-age=0',
}


def get_topic(pk):
    response = requests.get(
        'https://api.zhihu.com/topics/' + pk, headers=headers, cookies=cookies)
    return response.json()
