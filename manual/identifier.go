package main

type Identifier struct {
	Package string
	Name    string
}

func (instance Identifier) String() string {
	result := ""
	if len(instance.Package) > 0 {
		result += instance.Package + "."
	}
	result += instance.Name
	return result
}

func NewIdentifier(packageName string, name string) Identifier {
	return Identifier{
		Package: packageName,
		Name: name,
	}
}
