hkbuildings
===========

Small program to fetch info of private buildings of Hong Kong.

民政事務總署[香港私人大廈電腦資料庫](https://bmis2.buildingmgt.gov.hk/bd_hadbiex/content/searchbuilding/building_search.jsf?renderedValue=true)

```
go build hkbuildings.go
./hkbuildings 0 100 | tee -a data.json
```
