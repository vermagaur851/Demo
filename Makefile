BINDIR=build
LIBNAME=libamantyametrics
WRAPPER=./lang_wrapper/lang_wrapper.go
CFLAGS=-I$(BINDIR) -L$(BINDIR) -l:$(LIBNAME).so

build:
	mkdir -p $(BINDIR)
	go build -buildmode=c-shared -o $(BINDIR)/$(LIBNAME).so $(WRAPPER)

test-c: build
	gcc test/test.c -o $(BINDIR)/test_c $(CFLAGS)
	LD_LIBRARY_PATH=$(BINDIR) ./$(BINDIR)/test_c

test-cpp:
	g++ test/test.cpp -o build/test_cpp -Ibuild -Lbuild -l:libamantyametrics.so -std=c++11
	LD_LIBRARY_PATH=build ./build/test_cpp

test-python: build
	python3 test/test.py

clean:
	rm -rf $(BINDIR)
