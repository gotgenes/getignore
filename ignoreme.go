package ignoreme

type ignoreFetcher struct {
	baseUrl string
}

func (fetcher ignoreFetcher) NameToUrl(name string) string {
	return fetcher.baseUrl + "/" + name + ".gitignore"
}
