module gitlab.itestor.com/helei/utils.git/ftp

go 1.21

replace (
	gitlab.itestor.com/helei/utils.git => ../
	gitlab.itestor.com/helei/utils.git/crypto/sha256 => ../crypto/sha256
)

require github.com/jlaffaye/ftp v0.2.0

require (
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
)
