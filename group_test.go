package main

import (
	"testing"
)

func TestGroup_Add(t *testing.T) {
	groups := []*Group{
		{

			Name: "0-3岁",
			Children: []*Group{
				{
					Name: "益智",
					Children: []*Group{
						{
							Name: "游戏",
						}, {
							Name: "语文",
						},
					},
				},
			},
		},
		{
			Name: "3~6岁",
		},
		{
			Name: "6岁以上",
		},
	}

	AddGroup(groups)
}
func AddGroup(groups []*Group, pids ...int) {
	pid := 0
	if len(pids) > 0 {
		pid = pids[0]
	}
	for _, g := range groups {
		g.Pid = pid
		g.Add()
		AddGroup(g.Children, g.Id)
	}
}
