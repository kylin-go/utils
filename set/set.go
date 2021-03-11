package set

import "sync"

type setMember map[interface{}]bool

type SetBody struct {
	member setMember
	lock   sync.Mutex
}

func CreateSet(member ...interface{}) SetBody {
	var s = SetBody{setMember{}, sync.Mutex{}}
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, m := range member {
		s.member[m] = true
	}
	return s
}

// 获取所有元素
func (s SetBody) Get() setMember {
	return s.member
}

// 判断元素是否存在
func (s SetBody) IsExist(member interface{}) bool {
	_, ok := s.member[member]
	return ok
}

// 统计元素数量
func (s SetBody) Count() int {
	return len(s.member)
}

// 添加元素
func (s SetBody) Add(member ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, m := range member {
		s.member[m] = true
	}
}

// 删除元素
func (s SetBody) Del(member ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, m := range member {
		if _, ok := s.member[m]; ok {
			delete(s.member, m)
		}
	}
}

//取交集
func (s SetBody) Intersection(body SetBody) SetBody {
	result := CreateSet()
	result.lock.Lock()
	defer result.lock.Unlock()
	mem := body.Get()
	for m, _ := range s.Get() {
		if _, ok := mem[m]; ok {
			result.Add(m)
		}
	}
	return result
}

// 获取并集
func (s SetBody) Union(body SetBody) SetBody {
	s.lock.Lock()
	defer s.lock.Unlock()
	for m, _ := range body.Get() {
		if ok := s.IsExist(m); !ok {
			s.Add(m)
		}
	}
	return s
}

// 获取差集
func (s SetBody) Difference(body SetBody) SetBody {
	s.lock.Lock()
	defer s.lock.Unlock()
	for m, _ := range body.Get() {
		if ok := body.IsExist(m); ok {
			s.Del(m)
		}
	}
	return s
}
