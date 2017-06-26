.PHONY: clean ingestor/ingestor stopwordfilter/stopwordfilter wordcounter/wordcounter searches/searches

all: ingestor/ingestor stopwordfilter/stopwordfilter wordcounter/wordcounter searches/searches

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
