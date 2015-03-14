package collections

type PersistentMap struct{
	mappings map[interface{}][]interface{}
}

func ExtendPersistentMap(orig map[interface{}]interface{}) PersistentMap {
	pm := NewPersistentMap()
	for k, v := range orig {
		pm.Set(k, v)
	}
	return pm
}

func NewPersistentMap() PersistentMap {
	return PersistentMap{make(map[interface{}][]interface{})}
}

func (m *PersistentMap) Contains(key interface{}) bool {
	_, mappingExists := m.mappings[key]
	return mappingExists
}

func (m *PersistentMap) Set(key interface{}, value interface{}) {
	if !m.Contains(key) {
		m.mappings[key] = []interface{}{}
	}
	m.mappings[key] = append(m.mappings[key], value)
}

func (m *PersistentMap) Get(key interface{}, defaultValue interface{}) interface{} {
	value, mappingExists := m.mappings[key]
	if !mappingExists {
		return defaultValue
	} else {
		return value[len(value) - 1]
	}
}

func (m *PersistentMap) GetOrNil(key interface{}) interface{} {
	return m.Get(key, nil)
}

func (m *PersistentMap) Freeze() map[interface{}]interface{} {
	frozen := make(map[interface{}]interface{})
	for k, v := range m.mappings {
		frozen[k] = v[0]
	}
	return frozen
}
