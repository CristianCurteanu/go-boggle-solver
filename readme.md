# Golang Boggle Solver

This project is a implementation of boggle solving algorithm in Golang

## Usage

In order to build make sure that you have `make` tool installed, and once cloned this repository run:

```
$ make build
```

Then you can run the app by using following command:

```
$ ./main
```

You can see the options of this app by calling:

```
$ ./main -h
```

which should produce an output like following:

```
Usage of ./main:
  -dictionary string
        Path to dictionary file (default "./dictionary.txt")
  -scheme string
        Path to boggle board schema definition file (default "./scheme")
```

These are the flags that may be applied when some dictionary or board definition files will be defined. By default it will take the files from current directory.

## Board file definition

It should be a text format file, where every columns are separated by space, and a single character will be defined. It will throw an error if for a cell there are multiple characters string defined. As an example of such file definition, you can check `scheme` file in the repository
