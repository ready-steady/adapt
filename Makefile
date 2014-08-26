packages := basis/newtoncotes interp/local
tests := $(addprefix test_,$(packages))

test: $(tests)

$(tests): test_%:
	cd $* && go test

.PHONY: test
