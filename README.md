To build this project please install Golang and dep, then on mac/linux run build.sh

Golang: https://golang.org/doc/install\
Dep: https://github.com/golang/dep#installation

###Running the program
Please run the right version of the binary for the os.
Under `output` folder there are 3 file
1. cps842-linux
2. cps842-mac
3. cps842-win

####Run invert
Program help, using the -h flag

```
./cps842 invert -h
Take a collection of documents and generate its inverted index.

Usage:
  cps842 invert [flags]

Flags:
  -f, --file string            The location of CACM collection (required)
  -h, --help                   help for invert
  -l, --lower-case             lower case each term
  -o, --output-folder string   The location to output the final files (required)
  -p, --porter                 enable Porter's Stemming algorithm
  -s, --stop-word string       add stop word removal
```

To run the program without porter stemming/ stop word and assuming the input file is `cacm.all`

`./cps842 invert -f cacm.all -o "./data"`

With stop word

`./cps842 invert -f cacm.all -o "./data" -s common_words`

For Porter's Stemming algorithm and to lower case every word add either `-p` or `-l` flags to the run.

Ex of all the flags, `./cps842 invert -f cacm.all -o "./data" -s common_words -l -p`

####Run test

Program help, using the -h flag

```
Take the posting file generate from invert and test it.
   
   Usage:
     cps842 test [flags]
   
   Flags:
     -h, --help             help for test
     -p, --posting string   The location of Posting File(required)

```

Given that the posting file is in the folder `data`. Run the following command to execute the program

```
 ./cps842 test -p ./data/postings
```