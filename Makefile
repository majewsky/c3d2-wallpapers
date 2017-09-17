all: $(foreach res,1920x1080 3840x2160,$(foreach variant,v1 v2,outputs/$(variant)-$(res).png))

%/render: %/*.go
	go build -o $@ $<

outputs/v1-%.png: v1/render inputs/logo-plain-with-glow-%.png
	./v1/render inputs/logo-plain-with-glow-$*.png | convert pgm:- $@
	optipng $@

outputs/v2.svg: v2/render
	./$< > $@

outputs/v2-%.png: outputs/v2.svg
	inkscape -z -e $@ -w $(shell echo $* | cut -dx -f1) -h $(shell echo $* | cut -dx -f2) $<
	optipng $@
