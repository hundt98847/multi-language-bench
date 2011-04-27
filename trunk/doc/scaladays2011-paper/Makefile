all: main.pdf

pdf: main.dvi
	dvipdf main.dvi

ps: main.ps

main.pdf: *.tex *.bib 
	pdflatex main.tex
	bibtex main || echo "Bibtex failed"
	pdflatex main.tex
	pdflatex main.tex

main.ps: main.dvi
	dvips -o main.ps main.dvi

main.dvi: *.tex *.bib
	latex main.tex
	bibtex main || echo "Bibtex failed"
	latex main.tex
	latex main.tex

spellcheck:
	for i in *.tex; do aspell --personal=./personal.wl -c $$i; done

test: pdf
	acroread main.pdf &
clean:
	rm -f *.dvi *.ps *.pdf *.bbl *.aux *.log

xclean: clean
	rm -f *~ */*~ */*/*~
