package api

type SaveLoad interface {
	SaveGame(autosave bool) (string, error)
}
