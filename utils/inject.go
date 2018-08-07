package utils

import (
	"reflect"
	"strings"

	"github.com/facebookgo/inject"
	"github.com/facebookgo/structtag"
)

//Provider provider interface
type Provider interface {
	Provide() []*inject.Object
}

// ProviderFunc type
type ProviderFunc func() []*inject.Object

//Provide 执行 ProviderFunc 方法
func (f ProviderFunc) Provide() []*inject.Object {
	return f()
}

func Inject(objs ...interface{}) {
	// 递推遍历所有字元素，将携带 `inject` 字眼的成员提取出来
	names := extractNamedInject(objs)

	// 对于db、cfg等常规项，msf框架按约定的命名规则自动将其初始化
	commonProvider := newCommonProvider(names)

	// 仿 inject.Populate() 对带有部分初始化节点的DAG图，全部铺开
	var g inject.Graph
	g.Provide(commonProvider.Provide()...)
	for _, obj := range objs {
		if p, _ := obj.(Provider); p != nil {
			g.Provide(p.Provide()...)
		} else {
			g.Provide(&inject.Object{Value: obj})
		}
	}
	g.Populate()
}

func newCommonProvider(names map[string]reflect.Type) Provider {
	return ProviderFunc(func() []*inject.Object {
		objs := []*inject.Object{}
		for name, _ := range names {
			var v interface{}
			switch {
			// TODO 这里有待优化，需要根据类型判断
			case strings.HasPrefix(name, "ldap"):
				//v = NewLdap()
			}
			if v != nil {
				objs = append(objs, &inject.Object{Name: name, Value: v})
			}
		}
		return objs
	})
}

func extractNamedInject(objs []interface{}) map[string]reflect.Type {
	names := map[string]reflect.Type{}
	for _, o := range objs {
		reflectType := reflect.TypeOf(o)
		_extractNamedInject(reflectType, names)
	}
	return names
}

func _extractNamedInject(reflectType reflect.Type, names map[string]reflect.Type) {
	if !_isPtrStruct(reflectType) {
		return
	}

	for i := 0; i < reflectType.Elem().NumField(); i++ {
		field := reflectType.Elem().Field(i)
		fieldType := field.Type
		fieldTag := field.Tag
		found, value, err := structtag.Extract("inject", string(fieldTag))
		if err != nil || found == false {
			continue
		}
		if value != "" {
			names[value] = fieldType
		}
		if _isPtrStruct(fieldType) {
			_extractNamedInject(fieldType, names)
		}
	}
}

func _isPtrStruct(reflectType reflect.Type) bool {
	if reflectType.Kind() != reflect.Ptr {
		return false
	}
	t := reflectType.Elem()
	if t.Kind() != reflect.Struct {
		return false
	}
	return true
}
