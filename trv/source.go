package trv

type Source struct {
	Owner        string `json:"owner"`
	Repo         string `json:"repo"`
	Path         string `json:"path"`
	IsEnterprise bool   `json:"isEnterprise"`
	Token        string `json:"token"`
	BaseURL      string `json:"baseURL"`
	UploadURL    string `json:"uploadURL"`
}

/*
If there is DB data locally, load it and return it.
If not, retrieve it from a remote location.
*/
func (s Source) setDbData() Db {
	var db Db

	db.loadData(s.Repo, s.Path)

	client, ctx := s.NewClient()

	if len(db.tables) != 0 {
		return db
	}

	db.tables = fetchDbInfo(client, ctx, s)
	db.saveData(s.Repo, s.Path)
	return db
}
