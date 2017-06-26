.PHONY: clean ingestor/ingestor stopwordfilter/stopwordfilter wordcounter/wordcounter searches/searches windower/windower

all: ingestor/ingestor stopwordfilter/stopwordfilter wordcounter/wordcounter searches/searches windower/windower

clean:
	rm -f ingestor/ingestor &&\
	rm -f stopwordfilter/stopwordfilter &&\
	rm -f wordcounter/wordcounter

ingestor/ingestor:
	cd ingestor && $(MAKE)

stopwordfilter/stopwordfilter:
	cd stopwordfilter && $(MAKE)

wordcounter/wordcounter:
	cd wordcounter && $(MAKE)

searches/searches:
	cd searches && $(MAKE)

windower/windower:
	cd windower && $(MAKE)
