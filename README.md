# ttyconvert: convert ttyrec format from/to json

## Install

- go get github.com/wtnb75/ttyconvert

## Usage

```
NAME:
   ttyconvert - ttyrec <-> json converter

USAGE:
   ttyconvert [tojson|fromjson] < input > output

VERSION:
   0.0.1

COMMANDS:
     tojson    convert to json from ttyrec
     fromjson  convert from json to ttyrec
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

- ttyrec output.ttyrec
  - do something
- ttyconvert tojson < output.ttyrec > to-edit.json
- vi to-edit.json
  - edit timing, string, etc...
- ttyconvert fromjson < to-edit.json > new-output.ttyrec
- ttyplay new-output.ttyrec

## tojson

one json per line

- Time
    - microseconds from epoch
- TimeDelta
    - microseconds from previous data
- Data
    - string

## fromjson

- fromjson use TimeDelta value to set timing info
- when "Filter" exists in each json object, call filter command to convert "Data"
  - `{"Filter":"lolcat"}`
  - `{"Filter":"cowsay"}`
  - etc...
