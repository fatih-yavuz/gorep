test:
	rm test_files/test.tar.gz > /dev/null 2>&1; tar zcvf test_files/test.tar.gz test_resources/*;
	rm test_files/test.gz > /dev/null 2>&1; gzip --keep -f test_resources/file1.html && mv test_resources/file1.html.gz test_files/test.gz;
	go test

go:
	make tar && go build zgrep.go && ./zgrep