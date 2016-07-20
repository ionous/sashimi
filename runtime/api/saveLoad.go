package api

type SaveLoad interface {
	SaveGame() (string, error)
}
