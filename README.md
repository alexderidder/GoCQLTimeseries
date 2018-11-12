
# GoCQLSockets

Save kWh, wattage & power factor using websockets and GoCQL

## Usage

Send messages according protocol below

## Standard Message Header
In general, each message consists of a standard message header followed by request-specific data. The standard message header is structured as follows:

```golang
type Header struct {
	MessageLength, RequestID, ResponseID, OpCode uint32
}
```

Field | Description
------------ | -------------
messageLength | The total size of the message in bytes. This total includes the 4 bytes that holds the message length
requestID | A client or database-generated identifier that uniquely identifies this message.
responseTo | Clients can use the requestID and the responseTo fields to associate query responses with the originating query.
opCode | Type of message. See Request Opcodes for details.

## Request Opcodes
Opcode Name | Value | Comment
------------ | ------------- | -------------
OP_REPLY | 1 | Reply to a client request. responseTo is set.
OP_QUERY | 100 | Query measurements by stone_id(s), fields, time and interval
OP_UPDATE | 200 | Na
OP_INSERT | 300 | Na
OP_DELETE | 400 | Na

## Client Request Messages
###  OP_QUERY
```golang
struct OP_QUERY {
    MsgHeader header,
    int32     flags,
    json   queryDetails
}
```

Field | Description
------------ | -------------
header | Message header, as described in Standard Message Header.
flags | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **0** - default <br>  **1** -  ...
queryDetails| **0** - <br>type OP_QUERY struct { <br> &nbsp;&nbsp;&nbsp;&nbsp; StoneIDs &nbsp;&nbsp; []gocql.UUID &nbsp;  `json:"stoneIDs"` <br> &nbsp;&nbsp;&nbsp;&nbsp; Types &nbsp;&nbsp; &nbsp;&nbsp; []string &nbsp; &nbsp; &nbsp; `json:"types""` <br> &nbsp;&nbsp;&nbsp;&nbsp; StartTime &nbsp;&nbsp;time.Time &nbsp; &nbsp;&nbsp; `json:"startTime"` <br> &nbsp;&nbsp;&nbsp;&nbsp; EndTime &nbsp;&nbsp;&nbsp; time.Time &nbsp; &nbsp; &nbsp;`json:"endTime"` <br> &nbsp;&nbsp;&nbsp;&nbsp; Interval &nbsp;&nbsp; uint32 &nbsp; &nbsp; &nbsp; &nbsp; `json:"interval"` <br>} 

The database will respond to an OP_QUERY message with an [OP_REPLY](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/###OP_REPLY) message.
			
###  OP_INSERT 
```golang
struct OP_INSERT {
    MsgHeader header,
    int32     flags,
    json   queryDetails
}
```

Field | Description
------------ | -------------
header | Message header, as described in Standard Message Header.
flags | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **0** - default <br>  **1** -  ...
insertDetails| **0** - <br>type OP_INSERT struct { <br> &nbsp;&nbsp;&nbsp;&nbsp; StoneID &nbsp; gocql.UUID &nbsp;`json:"stoneID"` <br> &nbsp;&nbsp;&nbsp;&nbsp; Data []struct{ <br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Time&nbsp; &nbsp; &nbsp; &nbsp;  time.Time &nbsp;&nbsp; `json:"time"`<br> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; KWH&nbsp; &nbsp; &nbsp; &nbsp; &nbsp;float32 &nbsp; &nbsp;&nbsp; `json:"kWh"` <br> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Watt &nbsp; &nbsp; &nbsp; &nbsp;float32 &nbsp; &nbsp;&nbsp; `json:"watt"` <br> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; PowerFactor float32 &nbsp; &nbsp; &nbsp;`json:"pf"`<br> &nbsp;&nbsp;&nbsp;&nbsp; } <br>	} 
			
Atm flag 2(No Content) without JSON (Maybe: There is no response to an OP_INSERT message. ) 

###  OP_DELETE
```golang
struct OP_DELETE{
    MsgHeader header,
    int32     flags,
    json   queryDetails
}
```

Field | Description
------------ | -------------
header | Message header, as described in Standard Message Header.
flags | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **0** - default <br>  **1** -  ...
deleteDetails| **0** - <br>type OP_DELETE struct { <br> &nbsp;&nbsp;&nbsp;&nbsp; StoneID &nbsp;&nbsp;&nbsp; gocql.UUID &nbsp;&nbsp;&nbsp;  `json:"stoneID"` <br> &nbsp;&nbsp;&nbsp;&nbsp; Types &nbsp;&nbsp; &nbsp;&nbsp; []string &nbsp; &nbsp; &nbsp; `json:"types""` <br> &nbsp;&nbsp;&nbsp;&nbsp; StartTime &nbsp;&nbsp;time.Time &nbsp; &nbsp;&nbsp; `json:"startTime"` <br> &nbsp;&nbsp;&nbsp;&nbsp; EndTime &nbsp;&nbsp;&nbsp; time.Time &nbsp; &nbsp; &nbsp;`json:"endTime"`  <br>} 

Atm flag 2(No Content) without JSON (Maybe: There is no response to an OP_DELETE message. )

##  Database Response Messages
###  OP_REPLY

```golang
struct OP_REPLY{
    MsgHeader header,
    int32     flags,
    json   replyDetails
}
```
Field | Description
------------ | -------------
header | Message header, as described in Standard Message Header.
flags | (Bit vector to specify flags for the operation. The bit values correspond to the following: <br>  **1** -  OK <br>  **2** -  No Content <br>  **100** -  Error
queryDetails|  **1** - Response to an OP_QUERY  <br> type OP_INSERT struct { <br> &nbsp;&nbsp;&nbsp;&nbsp; StartTime &nbsp;&nbsp;time.Time &nbsp;`json:"startTime"` <br> &nbsp;&nbsp;&nbsp;&nbsp; EndTime &nbsp;&nbsp;&nbsp; time.Time &nbsp;`json:"endTime"` <br> &nbsp;&nbsp;&nbsp;&nbsp; Interval &nbsp; &nbsp;uint32 &nbsp; &nbsp; `json:"interval"` <br>&nbsp;&nbsp;&nbsp;&nbsp; Stones []struct { <br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; StoneID&nbsp; gocql.UUID &nbsp;&nbsp;`json:"stoneID"`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Fields &nbsp; []struct { <br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;Field&nbsp;&nbsp;string &nbsp; &nbsp;&nbsp; `json:"field"`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp; Measurements[]struct{  <br> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Time&nbsp; &nbsp; &nbsp; &nbsp;time.Time &nbsp; &nbsp;&nbsp; `json:"time"` <br> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Value &nbsp;&nbsp;&nbsp;&nbsp; float32 &nbsp; &nbsp; &nbsp;&nbsp;&nbsp;`json:"value"`<br> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp; } `json:"data"`<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; } `json:"fields"`<br> &nbsp;&nbsp;&nbsp; } `json:"stones"` <br>	}  <br> **2** - No JSON <br> **3** - See  [Error code](https://github.com/alexderidder/GoCQLTimeseries/blob/master/README.md/####Error_codes)

#### Error_codes

```golang
struct Error{
    uin32 code,
    string message
}
```

see errors in model/error.go (to be continued)
