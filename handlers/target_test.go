package handlers

import (
	"fmt"
	"testing"

	"github.com/devops-rob/boundary-census/clients/boundary"
	"github.com/devops-rob/boundary-census/clients/boundary/mocks"

	"github.com/hashicorp/boundary/api/targets"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

func setupTarget(t *testing.T) (*Target, *mocks.Client) {
	m := &mocks.Client{}
	l := hclog.NewNullLogger()

	// setup mocks for happy path
	m.On("FindProjectIDByName", "myscope", "myproject").Return("123abc", nil)
	m.On("CreateTarget", "mytarget_9090", uint32(9090), "123abc").Return(&targets.Target{Id: "target1"}, nil)
	m.On("CreateTarget", "mytarget_9091", uint32(9091), "123abc").Return(&targets.Target{Id: "target2"}, nil)

	return NewTarget(l, m), m
}

func TestCreateReturnsErrorWhenProjectDoesNotExist(t *testing.T) {
	tgt, mc := setupTarget(t)

	si := &ServiceInstance{
		Location: "127.0.0.1",
		Ports:    []uint32{9090, 9091},
	}

	mc.On("FindProjectIDByName", "broken", "myproject").Return("", boundary.ProjectNotFoundError)

	_, err := tgt.Create(si, "mytarget", "broken", "myproject")
	require.Error(t, err)
}

func TestCreateReturnsErrorWhenTargetCreateFails(t *testing.T) {
	tgt, m := setupTarget(t)

	si := &ServiceInstance{
		Location: "127.0.0.1",
		Ports:    []uint32{9002, 9091},
	}

	m.On("CreateTarget", "mytarget_9002", uint32(9002), "123abc").Return(nil, fmt.Errorf("boom"))

	_, err := tgt.Create(si, "mytarget", "myscope", "myproject")
	require.Error(t, err)
}

func TestCreatesTargets(t *testing.T) {
	tgt, m := setupTarget(t)

	si := &ServiceInstance{
		Location: "127.0.0.1",
		Ports:    []uint32{9090, 9091},
	}

	_, err := tgt.Create(si, "mytarget", "myscope", "myproject")
	require.NoError(t, err)

	m.AssertCalled(t, "CreateTarget", "mytarget_9090", uint32(9090), "123abc")
	m.AssertCalled(t, "CreateTarget", "mytarget_9091", uint32(9091), "123abc")
}
