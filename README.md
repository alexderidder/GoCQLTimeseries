# GoCQLSockets

Save kWh, wattage & power factor using websockets and GoCQL

# Setup
 # Install Go
https://golang.org/doc/install 

Setup go path 
https://github.com/golang/go/wiki/SettingGOPATH
## Get libraries
```
go get github.com/gocql/gocql 
```
```
go get github.com/stretchr/testify
```
 # Install Cassandra
http://cassandra.apache.org/download/
 ## Run Cassandra
In your cassandra folder run file 'bin/cassandra.exe'
 ## Run cqlsh
In your cassandra folder run file 'bin/cqlsh.exe'
 ### Setup Cassandra tables
```
 run setup.sh
```
 # Run sockets
 
## Setup config file
Edit conf.json to your setup
 ## Run server
go run main.go --mode server

Send messages according protocol below

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
uint32 | flag | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **0** [Select measurement payload](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#select-measurement-payload)


#### Select measurement payload
```json 
{  
   "stoneIDs":[  
      "bf82e78d-24a2-470d-abb8-9e0a2720619f"
   ],
   "types":[  
      "w",
      "pf",
      "kwh"
   ],
   "startTime":"2018-11-12T14:01:59.1708508+01:00", //optional
   "endTime":"2018-11-12T14:31:59.1708508+01:00", //optional
   "interval":0 //optional
}
```

#### Response select measurement payload
```json 
{  
   "startTime":"2018-11-12T14:01:59.1708508+01:00", //optional
   "endTime":"2018-11-12T14:31:59.1708508+01:00", //optional
   "interval":0, //optional
   "stones":[  
      {  
         "stoneID":"bf82e78d-24a2-470d-abb8-9e0a2720619f",
         "fields":[  
            {  
               "field":"w",
               "Data":[  
                  {  
                     "time":"2018-11-12T13:19:55.148Z",
                     "value":3.0233014
                  },
                  {  
                     "time":"2018-11-12T13:19:55.149Z",
                     "value":2.188571
                  }
               ]
            },
            {  
               "field":"pf",
               "Data":[  
                  {  
                     "time":"2018-11-12T13:19:55.148Z",
                     "value":4.702545
                  },
                  {  
                     "time":"2018-11-12T13:19:55.149Z",
                     "value":2.1231875
                  }
               ]
            },
            {  
               "field":"kwh",
               "Data":[  
               ]
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
uint32 | flag | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **0** [Insert measurement payload](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#insert-measurement-payload)


#### Insert measurement payload
```json 
{  
   "stoneID":"bf82e78d-24a2-470d-abb8-9e0a2720619f", 
   "data":[  
      {  
         "time":"2018-11-12T13:54:38.5078751+01:00",
         "kWh":3.3228004, 
         "watt":3.0233014,
         "pf":4.702545 
      },
      {  
         "time":"2018-11-12T13:54:39.5078751+01:00",
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
uint32 | flag | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **0** [Delete measurement payload](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#delete-measurement-payload)


#### Delete measurement payload
```json 
{  
   "stoneID":"bf82e78d-24a2-470d-abb8-9e0a2720619f",
   "types":[  
      "w",
      "pf",
      "kWh"
   ],
   "startTime":"0001-01-01T00:00:00Z", //optional
   "endTime":"0001-01-01T00:00:00Z" //optional
}
```
The database will respond to an OP_QUERY message with an [OP_REPLY](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#op_reply) message.

##  Database Response Messages
###  OP_REPLY

type | Name | Description
------------ | ------------ | -------------
16 byte | header | Message header, as described in [Standard Message Header](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#standard-message-header).
uint32 | flag | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **1** OK (Payload depends of [Client request message](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#client-request-messages)<br> **2** No Content (No payload)<br> **100** [Error](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/#error-codes)


#### Error codes

```json
{  
   "code":100,
   "message":"StoneID is missing"
}
```

see errors in model/error.go (work in progress)

# Run tests

to be continued
