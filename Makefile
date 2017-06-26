.PHONY: clean ingestor/ingestor stopwordfilter/stopwordfilter wordcounter/wordcounter searches/searches windower/windower

all: ingestor/ingestor stopwordfilter/stopwordfilter wordcounter/wordcounter searches/searches windower/windower

clean:
	cd ingestor && $(MAKE) clean
	cd stopwordfilter && $(MAKE) clean
	cd wordcounter && $(MAKE) clean
	cd searches && $(MAKE) clean
	cd windower && $(MAKE) clean

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
