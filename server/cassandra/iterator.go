package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/mock"
)

// Iterator iterates CQL query result rows.
type Iterator interface {
	// Close closes the Iterator.
	Close() error

	// Scan puts the current result row in results and returns whether there are
	// more result rows.
	Scan(results ...interface{}) bool

}

var (
	_ Iterator = IteratorMock{}
	_ Iterator = iterator{}
)

// IteratorMock is a mock Iterator. See github.com/maraino/go-mock.
type IteratorMock struct {
	mock.Mock
}

// Close implements Iterator.
func (m IteratorMock) Close() error {
	return nil
}

type testData struct {
	data int64
	x    int64
}

// Scan implements Iterator.
func (m IteratorMock) Scan(results ...interface{}) bool {
	args := m.Called(results...)
	return args.Bool(0)
}


type iterator struct {
	i *gocql.Iter
}

func (i iterator) Close() error {
	return i.i.Close()
}

func (i iterator) Scan(results ...interface{}) bool {
	return i.i.Scan(results...)
}
