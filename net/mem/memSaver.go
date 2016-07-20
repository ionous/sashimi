package mem

// future: change to io.reader?
type SaveGameBlob []byte

type MemSaver interface {
	SaveBlob(slot string, blob SaveGameBlob) (string, error)
	LoadBlob(slot string) (SaveGameBlob, error)
}
