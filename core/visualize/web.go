package visualize

import (
	"brain/adapters/store/sqlite"
	"brain/core/object"
	"encoding/json"
	"fmt"
	"os"
)

type WebVisualizer struct{}

func init() {
	Register("web", &WebVisualizer{})
}

func (w *WebVisualizer) Visualize(objects []*object.Object, links []sqlite.Link) error {
	type Node struct {
		ID    string `json:"id"`
		Label string `json:"label"`
		Type  string `json:"type"`
	}
	type Link struct {
		Source string `json:"source"`
		Target string `json:"target"`
	}
	type GraphData struct {
		Nodes []Node `json:"nodes"`
		Links []Link `json:"links"`
	}

	data := GraphData{}
	for _, o := range objects {
		data.Nodes = append(data.Nodes, Node{
			ID:    fmt.Sprintf("%x", o.ID),
			Label: fmt.Sprintf("%s (%x)", o.Type, o.ID[:4]),
			Type:  o.Type,
		})
	}
	for _, l := range links {
		data.Links = append(data.Links, Link{
			Source: fmt.Sprintf("%x", l.Parent),
			Target: fmt.Sprintf("%x", l.Child),
		})
	}

	jsonData, _ := json.Marshal(data)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Brain Graph</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <style>
        body { margin: 0; background: #1e1e1e; color: #fff; font-family: sans-serif; overflow: hidden; }
        .node circle { stroke: #fff; stroke-width: 1.5px; }
        .node text { font-size: 10px; fill: #ccc; pointer-events: none; }
        .link { stroke: #999; stroke-opacity: 0.6; stroke-width: 1px; }
        #controls { position: absolute; top: 10px; left: 10px; background: rgba(0,0,0,0.5); padding: 10px; border-radius: 5px; }
    </style>
</head>
<body>
    <div id="controls"><h1>Brain Graph</h1><p>Drag nodes to explore</p></div>
    <svg id="graph"></svg>
    <script>
        const data = %s;
        const width = window.innerWidth;
        const height = window.innerHeight;

        const svg = d3.select("#graph")
            .attr("width", width)
            .attr("height", height);

        const simulation = d3.forceSimulation(data.nodes)
            .force("link", d3.forceLink(data.links).id(d => d.id).distance(100))
            .force("charge", d3.forceManyBody().strength(-300))
            .force("center", d3.forceCenter(width / 2, height / 2));

        const link = svg.append("g")
            .attr("class", "links")
            .selectAll("line")
            .data(data.links)
            .enter().append("line")
            .attr("class", "link");

        const node = svg.append("g")
            .attr("class", "nodes")
            .selectAll("g")
            .data(data.nodes)
            .enter().append("g")
            .call(d3.drag()
                .on("start", dragstarted)
                .on("drag", dragged)
                .on("end", dragended));

        node.append("circle")
            .attr("r", 8)
            .attr("fill", d => d.type === 'task' ? '#ff9800' : '#2196f3');

        node.append("text")
            .attr("dx", 12)
            .attr("dy", ".35 font-size")
            .text(d => d.label);

        simulation.on("tick", () => {
            link
                .attr("x1", d => d.source.x)
                .attr("y1", d => d.source.y)
                .attr("x2", d => d.target.x)
                .attr("y2", d => d.target.y);

            node
                .attr("transform", d => "translate(" + d.x + "," + d.y + ")");
        });

        function dragstarted(event) {
            if (!event.active) simulation.alphaTarget(0.3).restart();
            event.subject.fx = event.subject.x;
            event.subject.fy = event.subject.y;
        }

        function dragged(event) {
            event.subject.fx = event.x;
            event.subject.fy = event.y;
        }

        function dragended(event) {
            if (!event.active) simulation.alphaTarget(0);
            event.subject.fx = null;
            event.subject.fy = null;
        }
    </script>
</body>
</html>
`, jsonData)

	err := os.WriteFile("graph.html", []byte(html), 0644)
	if err != nil {
		return err
	}
	fmt.Println("Visualization generated: graph.html")
	return nil
}
