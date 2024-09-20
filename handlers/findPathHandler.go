package handlers

import (
	"github.com/GuohaoMa/tucowDemo/common/response"
	"github.com/GuohaoMa/tucowDemo/model"
	"github.com/gin-gonic/gin"
)

type PathRq struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type CheapestPathRq struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type Query struct {
	Paths    PathRq         `json:"paths,omitempty"`
	Cheapest CheapestPathRq `json:"cheapest,omitempty"`
}

type FindPathRq struct {
	Queries []Query `json:"queries,omitempty"`
}

type PathRs struct {
	From     string     `json:"from,omitempty"`
	To       string     `json:"to,omitempty"`
	AllPaths [][]string `json:"paths,omitempty"`
}

type CheapestPathRs struct {
	From string      `json:"from,omitempty"`
	To   string      `json:"to,omitempty"`
	Path interface{} `json:"paths,omitempty"`
}

type Answer struct {
	Paths    *PathRs         `json:"paths,omitempty"`
	Cheapest *CheapestPathRs `json:"cheapest,omitempty"`
}
type FindPathRs struct {
	Answers []Answer `json:"answers,omitempty"`
}

type EdgeCost struct {
	To   string
	Cost float64
}

func FindPathHandler(graph *model.Graph) gin.HandlerFunc {
	return func(c *gin.Context) {
		g := *&model.Graph{Db: graph.Db, Id: graph.Id}
		findPathRq := FindPathRq{}
		if err := c.ShouldBindJSON(&findPathRq); err != nil {
			response.InteralErrorWithMessage("Json binding failure.", c)
			return
		}

		if findPathRq.Queries == nil {
			response.ValidationFailureWithMessage("Invalid params.", c)
			return
		}

		g.Get()

		cy, _ := g.FindCycles()
		if len(cy) > 0 {
			response.ValidationFailureWithMessage("Cycle detected.", c)
			return
		}

		graphMap := make(map[string][]EdgeCost)
		for _, edge := range g.Edges {
			if edge.FromIdentity != edge.ToIdentity {
				graphMap[edge.FromIdentity] = append(graphMap[edge.FromIdentity], EdgeCost{To: edge.ToIdentity, Cost: edge.Cost})
			}
		}

		findPathRs := FindPathRs{}
		for _, q := range findPathRq.Queries {
			if q.Paths != (PathRq{}) {
				result := [][]string{}
				start, end := q.Paths.Start, q.Paths.End
				findAllPaths(start, start, end, []string{start}, &result, graphMap)
				findPathRs.Answers = append(findPathRs.Answers, Answer{Paths: &PathRs{From: start, To: end, AllPaths: result}})
			}
			if q.Cheapest != (CheapestPathRq{}) {
				result := []string{}
				start, end := q.Cheapest.Start, q.Cheapest.End
				_, path := findCheapestPath(start, start, end, 0, 100, []string{start}, result, graphMap)
				a := Answer{Cheapest: &CheapestPathRs{From: start, To: end, Path: path}}
				if len(path) == 0 {
					a.Cheapest.Path = false
				}
				findPathRs.Answers = append(findPathRs.Answers, a)
			}
		}

		c.IndentedJSON(200, findPathRs)
		return
	}
}

func findAllPaths(cur string, from string, end string, path []string, result *[][]string, graphMap map[string][]EdgeCost) {
	if cur == end {
		pathCopy := make([]string, len(path))
		copy(pathCopy, path)
		*result = append(*result, pathCopy)
		return
	}
	for _, next := range graphMap[cur] {
		if next.To != from {
			findAllPaths(next.To, cur, end, append(path, next.To), result, graphMap)
		}
	}
}

func findCheapestPath(cur string, from string, end string, curCost float64, minCost float64, path []string, result []string, graphMap map[string][]EdgeCost) (float64, []string) {
	if cur == end && curCost < minCost {
		return curCost, path
	}
	for _, next := range graphMap[cur] {
		if next.To != from {
			minCost, result = findCheapestPath(next.To, cur, end, curCost+next.Cost, minCost, append(path, next.To), result, graphMap)
		}
	}
	return minCost, result
}
