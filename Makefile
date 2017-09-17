all: $(foreach res,1920x1080 3840x2160,outputs/v1-$(res).png)

%/render: %/*.go
	go build -o $@ $<

outputs/v1-%.png: v1/render inputs/logo-plain-with-glow-%.png
	./v1/render inputs/logo-plain-with-glow-$*.png | convert pgm:- $@
	optipng $@
