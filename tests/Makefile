test:
	rm -Rf out
	../venom run --html-report --output-dir=out --var-from-file cfg/vars.yaml -v .
	go tool covdata percent -i out/coverage
	go tool covdata textfmt -i out/coverage -o out/coverage.txt
	go tool cover -html out/coverage.txt -o out/coverage.html