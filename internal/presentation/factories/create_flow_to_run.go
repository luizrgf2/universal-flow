package factories

import (
	"github.com/luizrgf2/universal-flow/internal/core/usecases"
	"github.com/luizrgf2/universal-flow/internal/infra"
)

func CreateFlowToRunFactory() *usecases.CreateFlowToRunUseCase {
	flowStateManagerService, err := infra.NewFlowStateManagerSqlite("flow_state.db")
	if err != nil {
		panic(err)
	}

	usecase := usecases.MakeCreateFlowToRun(flowStateManagerService)
	return usecase
}
