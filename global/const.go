package global

const (
	dburi = "mongodb+srv://standard:standard@cluster0.gxrc1.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	dbname = "blog-application"
	performance = 100
)

var (
	JwtSecret = []byte("blogSecret")
)