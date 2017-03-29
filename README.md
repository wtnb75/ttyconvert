# ttyconvert: convert ttyrec format from/to json

## Install

- go get github.com/wtnb75/ttyconvert

## Usage

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
