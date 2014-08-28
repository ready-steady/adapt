packages := $(dir $(shell find . -name main.go))
tests := $(addprefix test_,$(packages))

test: $(tests)

$(tests): test_%:
	cd $* && go test

.PHONY: test $(tests)
