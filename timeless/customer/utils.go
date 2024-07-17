package customer

func Paginate(people []Customer, chunkSize int) [][]Customer {
	var chunks [][]Customer
	for i := 0; i < len(people); i += chunkSize {
		end := i + chunkSize
		if end > len(people) {
			end = len(people)
		}
		chunks = append(chunks, people[i:end])
	}
	return chunks
}
