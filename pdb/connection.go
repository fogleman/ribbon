package pdb

type Connection struct {
	Serial1 int
	Serial2 int
}

func ParseConnections(line string) []*Connection {
	var connections []*Connection
	a := parseInt(line[6:11])
	b1 := parseInt(line[11:16])
	b2 := parseInt(line[16:21])
	b3 := parseInt(line[21:26])
	b4 := parseInt(line[26:31])
	for _, b := range []int{b1, b2, b3, b4} {
		if b != 0 {
			connections = append(connections, &Connection{a, b})
		}
	}
	return connections
}
