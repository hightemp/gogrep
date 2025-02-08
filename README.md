# gogrep

Multithreaded variation of grep made with golang.


```console
hightemp@computer-01:~/Projects/gogrep$ time grep -r "Rackham" /home/hightemp/
/home/hightemp/.config/Code/User/History/-18c78d46/PYrF.txt:1914 Английский перевод Harris Rackham гласит:
/home/hightemp/.config/Code/User/History/-18c78d46/PNUh.txt:1914 Английский перевод Harris Rackham гласит:
grep: /home/hightemp/.config/Code/User/globalStorage/state.vscdb: binary file matches
/home/hightemp/Projects/gogrep/test/file2.txt:1914 Английский перевод Harris Rackham гласит:
/home/hightemp/android-studio/plugins/textmate/lib/bundles/adoc/README.md:* [AsciiDoc](http://asciidoc.org/) by Stuart Rackham
/home/hightemp/.local/share/JetBrains/Toolbox/apps/intellij-idea-community-edition/plugins/textmate/lib/bundles/adoc/README.md:* [AsciiDoc](http://asciidoc.org/) by Stuart Rackham

real    8m46.712s
user    1m59.208s
sys     2m1.003s
hightemp@computer-01:~/Projects/gogrep$ time ./gogrep -r "Rackham" /home/hightemp/
/home/hightemp/.config/Code/User/History/-18c78d46/PNUh.txt:6:1914 Английский перевод Harris Rackham гласит:
/home/hightemp/.config/Code/User/History/-18c78d46/PYrF.txt:6:1914 Английский перевод Harris Rackham гласит:
/home/hightemp/.local/share/JetBrains/Toolbox/apps/intellij-idea-community-edition/plugins/textmate/lib/bundles/adoc/README.md:256:* [AsciiDoc](http://asciidoc.org/) by Stuart Rackham
/home/hightemp/Projects/gogrep/test/file2.txt:6:1914 Английский перевод Harris Rackham гласит:
/home/hightemp/android-studio/plugins/textmate/lib/bundles/adoc/README.md:256:* [AsciiDoc](http://asciidoc.org/) by Stuart Rackham

real    0m36.250s
user    4m16.947s
sys     1m0.827s
```