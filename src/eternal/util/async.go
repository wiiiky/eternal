package util

import (
	"fmt"
	"reflect"
	"sync"
)

/* 并发处理数据，保持数组顺序不变 */
func GoAsyncSlice(s interface{}, f func(interface{}) interface{}) []interface{} {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	datas := reflect.ValueOf(s)
	if datas.Kind() != reflect.Slice {
		return nil
	}

	tmp := make(map[string]interface{})
	for i := 0; i < datas.Len(); i++ {
		d := datas.Index(i).Interface()
		wg.Add(1)
		go func(d interface{}) {
			defer wg.Done()
			v := f(d)
			k := fmt.Sprintf("%p", d)
			defer mutex.Unlock()
			mutex.Lock()
			tmp[k] = v
		}(d)
	}

	wg.Wait()
	results := make([]interface{}, datas.Len())
	for i := 0; i < datas.Len(); i++ {
		d := datas.Index(i).Interface()
		k := fmt.Sprintf("%p", d)
		results[i] = tmp[k]
	}
	return results
}
