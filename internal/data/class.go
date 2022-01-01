package data

var ClassIndex = make(map[string]PHPClass)

type PHPClass struct {
	Name    string
	Methods []PHPClassMethod
}

func (c PHPClass) GetMethod(methodName string) (PHPClassMethod, bool) {
	for _, method := range c.Methods {
		if method.Name == methodName {
			return method, true
		}
	}
	return PHPClassMethod{}, false
}
