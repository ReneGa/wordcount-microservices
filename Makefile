.PHONY: clean recorder ingestor stopwordfilter wordcounter searches windower

all: ingestor recorder stopwordfilter wordcounter searches windower

clean:
	cd ingestor && $(MAKE) clean
	cd recorder && $(MAKE) clean
	cd stopwordfilter && $(MAKE) clean
	cd wordcounter && $(MAKE) clean
	cd searches && $(MAKE) clean
	cd windower && $(MAKE) clean

ingestor:
	cd ingestor && $(MAKE)

recorder:
	cd recorder && $(MAKE)

stopwordfilter:
	cd stopwordfilter && $(MAKE)

wordcounter:
	cd wordcounter && $(MAKE)

searches:
	cd searches && $(MAKE)

windower:
	cd windower && $(MAKE)
