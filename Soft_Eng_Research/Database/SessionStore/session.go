package sessionstore

import "github.com/gorilla/sessions"

// Session store
var StoreStuProf = sessions.NewCookieStore([]byte("studentsessionkey"))
