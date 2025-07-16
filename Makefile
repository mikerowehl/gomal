cmd/step0_repl/step0_repl: cmd/step0_repl/step0_repl.go
	cd cmd/step0_repl; go build

cmd/step1_read_print/step1_read_print: cmd/step1_read_print/step1_read_print.go pkg/reader/reader.go
	cd cmd/step1_read_print; go build

cmd/step2_eval/step2_eval: cmd/step2_eval/step2_eval.go pkg/reader/reader.go
	cd cmd/step2_eval; go build
