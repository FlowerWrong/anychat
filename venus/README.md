# venus

## 技术栈

* [gin](https://github.com/gin-gonic/gin)
* [xorm](https://github.com/xormplus/xorm)
* [viper](https://github.com/spf13/viper)

## hiject

```golang
func(table *Table) ModelName() string {
	return flect.New(table.Name).Singularize().Pascalize().String()
}

func(table *Table) FileName() string {
	return strings.ToLower(table.ModelName())
}
```
