package gasp

type Map struct {
	m     *Map
	key   Object
	value Object
}

func NewMap() Map {
	return Map{}
}

func (m Map) Get(key Object) (ret Object) {
	if m.key == nil {
		return
	}
	if m.key.String() == key.String() {
		return m.value
	}
	if m.m != nil {
		return m.m.Get(key)
	}
	return nil
}

func (m Map) Assoc(key Object, value Object) Map {
	return Map{
		m:     &m,
		key:   key,
		value: value,
	}
}
