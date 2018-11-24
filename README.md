# splitfile
Splits a file into multiple files all populated with a set number of lines or with an even number of lines per file.

Command line parameters

-i Name of the file to be split (string)  
-l Maximum lines file to be split into (int)  
-m Maximum number of files to be output (0 for all) (int)  
-e Split the file evenly over the maximum number of files (bool)  
-o Destination file name (string)  

e.g. - Split the input file over 5 files with no more than 68 lines in each file  
./splitfile -i test.txt -l 68 -m 5 -o testout.txt

e.g. - Split the input file over 5 files with an even number of lines in each file  
./splitfile -i test.txt -e -m 5 -o testout.txt

