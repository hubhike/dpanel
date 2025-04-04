package event

import (
	"context"
	"github.com/docker/docker/api/types/container"
)

var ContainerCreateEvent = "container_create"

type ContainerCreate struct {
	InspectInfo *container.InspectResponse
	ContainerId string
	Ctx         context.Context
}
