package utils

type CustomSlice interface {
	Contain(interface{}) bool
	DataReader([]interface{})
}

type StringSlice struct {
	lstData []string
}

func (w *StringSlice) DataReader(data []interface{}) {
	for _, v := range data {
		w.lstData = append(w.lstData, v.(string))
	}
}
func (w *StringSlice) Contain(searchValue interface{}) bool {
	for _, v := range w.lstData {
		if searchValue.(string) == v {
			return true
		}
	}
	return false
}

type Filter struct {
}

func (f *Filter) Contain(CSlice CustomSlice, data interface{}) bool {
	return CSlice.Contain(data)
}
