package util



const (
	UnitW      string = "w"
	Unitpf     string = "pf"
	UnitkWh    string = "kwh"
)


func GetUInt32FromIndex(index uint32, message []byte) uint32{
	flag := ByteToUint32(message, int(index))
	return flag
}