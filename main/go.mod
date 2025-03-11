module main

require example_pkg v0.0.0

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/linxGnu/grocksdb v1.9.8 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

replace example_pkg => ../example_pkg

go 1.24
