package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/half-blood-prince-2710/blog-aggregator/internal/config"
)

func cleanInput(str string) ([]string) {
		return strings.Fields(strings.ToLower(str))
}



func main(){
	
	cfg, err := config.Read()
	if err!=nil {
		fmt.Print("err",err.Error()," \n")
	}
	fmt.Print(cfg.DbUrl,"\n")
	err=cfg.SetUser("Manish")
		if err!=nil {
		fmt.Print("err",err.Error()," \n")
	}
	cfg, err = config.Read()
	if err!=nil {
		fmt.Print("err",err.Error()," \n")
	}

	scanner:=bufio.NewScanner(os.Stdin)

	for{
		str :=scanner.Text()
		cmd :=cleanInput(str)
		
		if len(cmd)==0 {

		}
	}
	// fmt.Print("db_url: ",cfg.DbUrl,"\n Username: ",cfg.Username,"\n")

}