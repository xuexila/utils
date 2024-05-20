module gitlab.itestor.com/helei/utils.git/rsa

go 1.21.1

toolchain go1.21.4

replace (
	gitlab.itestor.com/helei/utils.git => ../
	gitlab.itestor.com/helei/utils.git/crypto/sha256 => ../crypto/sha256
)

require (
	gitlab.itestor.com/helei/utils.git v0.0.0-00010101000000-000000000000
	gitlab.itestor.com/helei/utils.git/crypto/sha256 v0.0.0-00010101000000-000000000000
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	golang.org/x/net v0.17.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/gorm v1.25.10 // indirect
)
