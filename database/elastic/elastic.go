package elastic

type Elastic struct {
	baseURI  string
	username string
	password string
}

func New() *Elastic {
	// uri, ok := os.LookupEnv("ELASTIC_URI")
	// if !ok {
	//     log.Fatal("elastic uri is not specified")
	// }
	//
	// password, ok := os.LookupEnv("ELASTIC_PASSWORD")
	// if !ok {
	// 	log.Fatal("elastic password is not specified")
	// }

	return &Elastic{
		//	baseURI:  uri,
		//	username: "elastic",
		//	password: password,
	}
}
