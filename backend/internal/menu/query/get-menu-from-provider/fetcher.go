package get_menu_from_provider

import "github.com/scrumno/scrumno-api/infra/integration-system/shared/interfaces"

type Fetcher struct {
	Provider interfaces.MenuProvider
}

type Query struct {
	OrgID       string
	TermGroupID string
	Token       string
}

func NewFetcher(
	provider interfaces.MenuProvider,
) *Fetcher {
	return &Fetcher{
		Provider: provider,
	}
}

func (f *Fetcher) Fetch(query Query) any {
	return f.Provider.GetMenu(query.OrgID)
}
