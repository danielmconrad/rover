package rover

import "log"

func logSuccess(params ...interface{}) {
	log.Println(append([]interface{}{"[SUCCESS] "}, params...)...)
}

func logInfo(params ...interface{}) {
	log.Println(append([]interface{}{"[INFO]    "}, params...)...)
}

func logWarning(params ...interface{}) {
	log.Println(append([]interface{}{"[WARN]    "}, params...)...)
}

func logError(params ...interface{}) {
	log.Println(append([]interface{}{"[ERROR]   "}, params...)...)
}

func logPanic(params ...interface{}) {
	log.Println(append([]interface{}{"[PANIC]   "}, params...)...)
}
