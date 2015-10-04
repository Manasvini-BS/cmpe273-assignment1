# cmpe273-assignment1
Assumptions :
-> User should always enter correct input in specified format
->When purchase price is equal to current price of the stock,it is considered as profit and hence "+" sign is added in the portfolio response.
Buy Input format : <"symbol:percentage,symbol:percentage"><space><budget>
ex:$ go run client.go  "YHOO:50%,AAPL:50%" 5000
Portfolio Input format: <tradeid>
ex : $ go run client.go 845
