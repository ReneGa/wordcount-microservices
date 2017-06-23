.PHONY: clean ingestor/ingestor stopwordfilter/stopwordfilter wordcounter/wordcounter

all: ingestor/ingestor stopwordfilter/stopwordfilter wordcounter/wordcounter

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
