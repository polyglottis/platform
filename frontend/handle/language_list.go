package handle

import (
	"log"
	"time"
)

func (w *Worker) fetchLanguageListPeriodically() {
	ok := w.fetchLanguageList(true)
	if !ok {
		ticker := time.NewTicker(time.Minute)
		for _ = range ticker.C {
			if w.fetchLanguageList(false) {
				break
			}
		}
		ticker.Stop()
	}
	for _ = range time.Tick(24 * time.Hour) {
		w.fetchLanguageList(true)
	}
}

// fetchLanguageList returns true only if the language list has been set and can be used.
func (w *Worker) fetchLanguageList(force bool) bool {
	if force || !w.languageListSet {
		list, err := w.Language.List()
		if err != nil {
			log.Println("Failed to fetch language list:", err)
			return false
		}

		err = w.Server.SetLanguageList(list)
		if err != nil {
			log.Println("Failed to set language list:", err)
			return false
		}
		w.languageListSet = true
		log.Println("Language list has been set.")
	}
	return true
}
