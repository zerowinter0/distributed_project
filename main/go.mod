module main

require (
	example_pkg v0.0.0
)

require google.golang.org/protobuf v1.36.5 // indirect

replace example_pkg => ../example_pkg

go 1.24
