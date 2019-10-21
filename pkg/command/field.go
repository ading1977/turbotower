package command

import "github.com/turbonomic/turbo-go-sdk/pkg/proto"

type commType int

const (
	boughtType commType = 0
	soldType   commType = 1
)

type display struct {
	header   string
	commName string
	commType commType
}

var (
	entitiesToTopCommoditiesMap = map[proto.EntityDTO_EntityType][]display{
		proto.EntityDTO_VIRTUAL_APPLICATION: {
			{
				header:   "QPS",
				commName: "TRANSACTION",
				commType: boughtType,
			},
			{
				header:   "LATENCY",
				commName: "RESPONSE_TIME",
				commType: boughtType,
			},
		},
		proto.EntityDTO_APPLICATION: {
			{
				header:   "VCPU",
				commName: "VCPU",
				commType: boughtType,
			},
			{
				header:   "VMEM",
				commName: "VMEM",
				commType: boughtType,
			},
		},
		proto.EntityDTO_CONTAINER: {
			{
				header:   "VCPU",
				commName: "VCPU",
				commType: soldType,
			},
			{
				header:   "VMEM",
				commName: "VMEM",
				commType: soldType,
			},
		},
		proto.EntityDTO_CONTAINER_POD: {
			{
				header:   "VCPU",
				commName: "VCPU",
				commType: soldType,
			},
			{
				header:   "VMEM",
				commName: "VMEM",
				commType: soldType,
			},
		},
		proto.EntityDTO_VIRTUAL_MACHINE: {
			{
				header:   "VCPU",
				commName: "VCPU",
				commType: soldType,
			},
			{
				header:   "VMEM",
				commName: "VMEM",
				commType: soldType,
			},
			{
				header:   "VCPUREQUEST",
				commName: "VCPU_REQUEST",
				commType: soldType,
			},
			{
				header:   "VMEMREQUEST",
				commName: "VMEM_REQUEST",
				commType: soldType,
			},
		},
		proto.EntityDTO_PHYSICAL_MACHINE: {
			{
				header:   "CPU",
				commName: "CPU",
				commType: soldType,
			},
			{
				header:   "MEM",
				commName: "MEM",
				commType: soldType,
			},
		},
		proto.EntityDTO_STORAGE: {
			{
				header:   "AMOUNT",
				commName: "STORAGE_AMOUNT",
				commType: soldType,
			},
			{
				header:   "LATENCY",
				commName: "STORAGE_LATENCY",
				commType: soldType,
			},
		},
		proto.EntityDTO_DATACENTER: {
			{
				header:   "CPU",
				commName: "CPU",
				commType: soldType,
			},
			{
				header:   "MEM",
				commName: "MEM",
				commType: soldType,
			},
		},
	}
)
