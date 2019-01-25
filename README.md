
# HOW DO I RUN THE TESTS?

# GoCQLSockets

Save kWh, wattage & power factor using websockets and GoCQL

# Setup
 ## Install Go
https://golang.org/doc/install 

### Setup go path 
https://github.com/golang/go/wiki/SettingGOPATH
### Get libraries

```
go get github.com/gocql/gocql 
```
```
go get github.com/stretchr/testify
```
When using old golang version
```
go get github.com/stretchr/testify
```
### Change imports
If you change the directory name, also change the imports in the go files of this project. 
 ## Install Cassandra
http://cassandra.apache.org/download/
 ### Run Cassandra
In your cassandra folder run file 'bin/cassandra.exe'
 ### Run cqlsh
In your cassandra folder run file 'bin/cqlsh.exe'
 #### Setup Cassandra tables
```
 run setup.sh
```
 # Run server
 
## Setup config file
Edit conf.json to your setup
 ## Two options
 
 
```
 go run main.go
```
 
```
 go build main.go
 .\main.exe
```


## Standard Message Header
In general, each message consists of a standard message header followed by request-specific data. The standard message header is structured as follows:

Type | Name | Description
------------ | ------------- | -------------
uint32 |messageLength | The total size of the message in bytes. This total includes the 4 bytes that holds the message length
uint32 |requestID | A client or database-generated identifier that uniquely identifies this message.
uint32 |responseTo | Clients can use the requestID and the responseTo fields to associate query responses with the originating query.
uint32 |opCode | Type of message. See Request Opcodes for details.

## Request Opcodes
Opcode Name | Value | Comment
------------ | ------------- | -------------
[OP_REPLY](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#op_reply) | 1 | Reply to a client request. responseTo is set.
[OP_QUERY](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#op_query)| 200 | Query measurements by stone_id(s), fields, time & interval
[OP_INSERT](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#op_insert) | 100 | Insert measurements by stone_id, time & type + value
[OP_DELETE](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#op_delete) | 500 | Delete measurements by stone_id, time & types

## Client Request Messages
###  OP_QUERY


type | Name | Description
------------ | ------------ | -------------
16 byte | header | Message header, as described in [Standard Message Header](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#standard-message-header).
uint32 | flag | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **1** Energy Usage  - [Request](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#query_request_type_1) - [Response](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#query_response_type_1)<br> **2** Power Usage - [Request](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#query_request_type_1) -  [Response](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#query_response_type_2)


#### Query request type 1
```json 
{  
   "stoneIDs":[  
      "5b8d0018acc9bc3124af2cc2"
   ],
   "startTime":"1548252675000", //Epoch in ms
   "endTime":"1548252905000", //Epoch in ms
   "interval":0, //optional
}
```

#### Query response type 1
```json 
{  
   "startTime":"1548252675000", //Epoch in ms
   "endTime":"1548252905000", //Epoch in ms
   "interval":0, //optional
   "stones":[  
      {  
         "stoneID":"bf82e78d-24a2-470d-abb8-9e0a2720619f",
         "Data":[  
            {  
               "time":"2018-11-12T13:19:55.148Z",
               "value":
               {
                    "kWh" : 3.0233014
               }
            },
            {  
               "time":"2018-11-12T13:19:55.149Z",
               "value":
               {
                   "kWh" : 2.188571
               }
            }
         ]
      }
   ]
}
```
#### Query response type 2
```json 
{  
   "startTime":"1548252675000", //Epoch in ms
   "endTime":"1548252905000", //Epoch in ms
   "interval":0, //optional
   "stones":[  
      {  
         "stoneID":"bf82e78d-24a2-470d-abb8-9e0a2720619f",
         "Data":[  
            {  
               "time":"2018-11-12T13:19:55.148Z",
               "value": {"w": 3, "pf" : 1}
            },
            {  
               "time":"2018-11-12T13:19:55.149Z",
               "value": {"w": 2, "pf" : 1}
            }
         ]
      }
   ]
}
```
###  OP_INSERT 

type | Name | Description
------------ | ------------ | -------------
16 byte | header | Message header, as described in [Standard Message Header](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#standard-message-header).
uint32 | flag | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **1** [Energy History](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#insert-request-type-1) <br>  **2** [Power History](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#insert-request-type-2)

#### Insert request type 1
```json 
{  
   "stoneID":"bf82e78d-24a2-470d-abb8-9e0a2720619f", 
   "data":[  
      {  
         "time":"1548252675000", //Epoch in ms
         "kWh":3.3228004
       },
       {  
          "time":"1548252679000", //Epoch in ms
          "kWh":3.4341154
       }
    ]
 }
```

#### Insert request type 2
```json 
{  
   "stoneID":"bf82e78d-24a2-470d-abb8-9e0a2720619f", 
   "data":[  
      {  
          "time":"1548252675000", //Epoch in ms
          "watt":3.0233014,
          "pf":4.702545 
       },
       {  
          "time":"1548252679000", //Epoch in ms
          "kWh":3.4341154
       }
    ]
}
```

The database will respond to an OP_QUERY message with an [OP_REPLY](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#op_reply) message.
	
###  OP_DELETE

type | Name | Description
------------ | ------------ | -------------
16 byte | header | Message header, as described in [Standard Message Header](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#standard-message-header).
uint32 | flag | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **1** Energy History  - [Request](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#delete_request_type_1)<br>  **2** Power History  - [Request](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#delete_request_type_1)

#### Delete request type 1
```json 
{  
   "stoneID":"bf82e78d-24a2-470d-abb8-9e0a2720619f",
   "startTime":"1548252679000", //Epoch in ms
   "endTime":"1548252680000" //optional
} 
```
The database will respond to an OP_QUERY message with an [OP_REPLY](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#op_reply) message.

##  Database Response Messages
###  OP_REPLY

type | Name | Description
------------ | ------------ | -------------
16 byte | header | Message header, as described in [Standard Message Header](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#standard-message-header).
uint32 | flag | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **1** Ok (Payload depends of [Client request message](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#client-request-messages ))<br> **2** No Content (No payload)<br> **100** [Error](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#error-codes)


#### Error codes

```json
{  
   "code":100,
   "message":"StoneID is missing"
}
```

see errors in datatypes/error.go 

# Run tests

```
go test ./...
   ```