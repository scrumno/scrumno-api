package get_menu

import "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"

type Fetcher struct {
	provider *interfaces.MenuProvider
	builder  *interfaces.MenuBuilder
}

func NewFetcher(provider *interfaces.MenuProvider, builder *interfaces.MenuBuilder) *Fetcher {
	return &Fetcher{provider: provider, builder: builder}
}

func (f *Fetcher) Fetch() {

}
