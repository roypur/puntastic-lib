package main
import("fmt"
       "puntastic")
func main(){
    Load("dictionary.gob")

    fmt.Println(puntastic.GeneratePun("pun"))
}
