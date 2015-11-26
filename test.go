package main
import("fmt"
       "github.com/roypur/puntastic-lib/lib/src/puntastic")
func main(){
    puntastic.Load("dictionary.gob")

    fmt.Println(puntastic.Get("pun"))
}
