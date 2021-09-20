# domain-registrations
Generate csv file containing every domain registered on specified date/date range
## Installation
```
go install github.com/lbirchler/domain-registrations@latest
```
## Examples 
### All domains registered on 9/15
```
$ ./domain-registrations -d 2021-09-15
```
result:
```
$ wc -l domains.csv
91657 domains.csv
```
```
$ head domains.csv
2021-09-15,5daystoexplosiveprofits.com
2021-09-15,5dcopywrite.com
2021-09-15,5dczzk.com
2021-09-15,5dm3xs.com
2021-09-15,5dollarsmusic.com
2021-09-15,5dsmm.com
2021-09-15,5dtarget.com
2021-09-15,5dtarget.online
2021-09-15,5ee2.com
2021-09-15,5eplay-liansai.com
```
### Domains registered from 9/1 to 9/7 that match regex "\ .(online|xyz)"**
```
$ ./domain-registrations -d 2021-09-01,2021-09-07 -r "\.(online|xyz)"
```
result:
```
$ wc -l domains.csv
111849 domains.csv
```
```
$ head domains.csv
2021-09-02,active-biy.xyz
2021-09-02,activesenior-supplement.online
2021-09-02,acxve.xyz
2021-09-02,ad1.xyz
2021-09-02,ad7tech.xyz
2021-09-02,ada2021.online
2021-09-02,adacardano.xyz
2021-09-02,adairiaq.xyz
2021-09-02,adamheitzman.xyz
2021-09-02,adaptshortly.xyz
```
## Command-line Options
```
  -d string
        domain registration date e.g. 2021-01-01 
        date range e.g. 2021-01-01,2021-01-15
  -o string
        csv output path (default "domains.csv")
  -r string
        only return domains that match provided regex 
        e.g. "^[a-zA-Z]\-[a-zA-Z0-9]{2,3}\.(xyz|club|shop|online)"
```
