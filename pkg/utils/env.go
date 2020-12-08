package utils

import(
	"os"
	"log"
	"strconv"
	_ "github.com/joho/godotenv/autoload"
)

//MustGet will return ENV variable, panic if not exists
func MustGet(envName string) string{
	v := os.Getenv(envName)
	if v == ""{
		log.Panicln("ENV variable missing: " + envName)
	}
	return v
}

func MustGetBool(envName string) bool{
	v := MustGet(envName)
	b, err := strconv.ParseBool(v)

	if err != nil{
		log.Panicln("ENV Bool parse error: " + envName + "\n" + err.Error())
	}
	return b

}