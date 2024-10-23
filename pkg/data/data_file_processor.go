package data

/*
	This type decouples the read/write from operations on the data
*/

// This implements the read/write behaviour for datafile
type DataFileReadWriter struct {
	dataFilePath string
}

// Reader interface
func (h *DataFileHandler) Read(p []byte) (n int, err error) {

}
