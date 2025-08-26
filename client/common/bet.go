
type Bet struct {
	name string 
	lastname string
	document string 
	birthdate string
	number string 
}


func newBet() *Bet {
	name string := os.Getenv("NAME")
	lastname string := os.Getenv("LASTNAME")
	document string := os.Getenv("DOCUMENT")
	birthdate string := os.Getenv("BIRTHDATE")
	number string := os.Getenv("NUMBER")

	return &Bet{
		name: name,
		lastname: lastname,
		document: document,
		birthdate: birthdate,
		number: number,
	}
}

func (bet Bet) serialize() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s", bet.name, bet.lastname, bet.document, bet.birthdate, bet.number)
}