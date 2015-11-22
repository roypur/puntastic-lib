package main
import("fmt"
       "puntastic")
func main(){
    Load("out.gob")

    fmt.Println(puntastic.GeneratePun("pun"))
}
