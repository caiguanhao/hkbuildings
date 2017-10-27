hkbuildings
===========

Small program to fetch info of private buildings of Hong Kong.

民政事務總署[香港私人大廈電腦資料庫](https://bmis2.buildingmgt.gov.hk/bd_hadbiex/content/searchbuilding/building_search.jsf?renderedValue=true)

File `data.json` was generated on 2017-10-26.

The server is slow and buggy. It takes time to fetch all data.

Example
-------

```
go build hkbuildings.go
./hkbuildings 0 100 | tee -a data.json
gem install -V write_xlsx
ruby csv.rb
ruby xlsx.rb
```
