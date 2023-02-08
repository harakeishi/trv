package trv

import "fmt"

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
func (s Source) setDbData() (DB, error) {
	var db DB

	db.loadData(s.Repo, s.Path)

	client, ctx, err := s.NewClient()
	if err != nil {
		return DB{}, err
	}

	if err := db.fetchDBInfo(client, ctx, s); err != nil {
		return DB{}, fmt.Errorf("set DB data fail:%w", err)
	}
	db.saveData(s.Repo, s.Path)
	return db, nil
}
