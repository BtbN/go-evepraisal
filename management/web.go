package management

import (
	"encoding/json"
	"expvar"
	"net/http"

	"github.com/evepraisal/go-evepraisal"
	boltdb "github.com/evepraisal/go-evepraisal/bolt"
	"github.com/husobee/vestigo"
)

// Context is the context that the web management needs
type Context struct {
	App                *evepraisal.App
	AppraisalBackupDir string
}

// HandleRestore is the handler for /restore, this is only used for partial restores
func (ctx *Context) HandleRestore(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var appraisal evepraisal.Appraisal
	err := json.NewDecoder(r.Body).Decode(&appraisal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ctx.App.AppraisalDB.PutNewAppraisal(&appraisal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleBackup is the handler for /backup/appraisals
// func (ctx *Context) HandleBackup(w http.ResponseWriter, req *http.Request) {
// 	db, ok := ctx.App.AppraisalDB.(*boltdb.AppraisalDB)
// 	if !ok {
// 		http.Error(w, "backup not supported for this database", http.StatusInternalServerError)
// 		return
// 	}
// 	err := db.DB.View(func(tx *bolt.Tx) error {
// 		w.Header().Set("Content-Type", "application/octet-stream")
// 		w.Header().Set("Content-Disposition", `attachment; filename="appraisals"`)
// 		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
// 		_, err := tx.WriteTo(w)
// 		return err
// 	})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// HandleBackup is the handler for /backup/appraisals
func (ctx *Context) HandleBackup(w http.ResponseWriter, req *http.Request) {
	db, ok := ctx.App.AppraisalDB.(*boltdb.AppraisalDB)
	if !ok {
		http.Error(w, "backup not supported for this database", http.StatusInternalServerError)
		return
	}

	err := db.Backup(ctx.AppraisalBackupDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HTTPHandler returns the http.Handler for the web management api
func HTTPHandler(app *evepraisal.App, appraisalBackupDir string) http.Handler {
	ctx := Context{App: app, AppraisalBackupDir: appraisalBackupDir}
	router := vestigo.NewRouter()
	router.Get("/backup/appraisals", ctx.HandleBackup)
	router.Post("/restore", ctx.HandleRestore)
	router.Handle("/expvar", expvar.Handler())
	return router
}
