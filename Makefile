cmd/step0_repl/step0_repl: cmd/step0_repl/step0_repl.go
	cd cmd/step0_repl; go build

cmd/step1_read_print/step1_read_print: cmd/step1_read_print/step1_read_print.go pkg/reader/reader.go
	cd cmd/step1_read_print; go build

cmd/step2_eval/step2_eval: cmd/step2_eval/step2_eval.go pkg/reader/reader.go
	cd cmd/step2_eval; go build

cmd/step3_env/step3_env: cmd/step3_env/step3_env.go pkg/reader/reader.go pkg/env/env.go
	cd cmd/step3_env; go build

cmd/step4_if_fn_do/step4_if_fn_do: cmd/step4_if_fn_do/step4_if_fn_do.go pkg/reader/reader.go pkg/env/env.go pkg/core/core.go
	cd cmd/step4_if_fn_do; go build

cmd/step5_tco/step5_tco: cmd/step5_tco/step5_tco.go pkg/reader/reader.go pkg/env/env.go pkg/core/core.go
	cd cmd/step5_tco; go build
