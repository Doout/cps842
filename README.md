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

####Run search

Program help, using the -h flag

```
./cps842 search -h
Search for the best documents match using cos-sim and output the best match

Usage:
  cps842 search [flags]

Flags:
  -f, --folder string   Folder location where posting/doc files are (required)
  -h, --help            help for search
  -i, --interact-mode   Interact with the user
  -s, --search string   What to search for

```

Given that the posting file is in the folder `data`. Run the following command to execute the program

```
 ./cps842 search -f ./data -s "<term/sent>"
```
Or to use the interact mode
```
 ./cps842 search -f ./data -i
```

####Run eval
Program help, using the -h flag
```
./cps842 eval -h
Evaluate the performance of the IR system, output will be the average MAP and R-Precision values over all queries

Usage:
  cps842 eval [flags]

Flags:
  -f, --folder string   Folder location where the posting/doc files are (required)
  -h, --help            help for eval
  -r, --qrels string    The qrels.text file
  -q, --query string    The query file

```

Given that the input folder is`data`, qrels `./input/qrels.text` and query is `./input/query.text` . Run the following command to execute the program

```
./cps842 eval -f data -r ./input/qrels.text -q ./input/query.text
```
